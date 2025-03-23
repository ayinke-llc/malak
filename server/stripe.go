package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/billing"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/render"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type stripeHandler struct {
	user            malak.UserRepository
	planRepo        malak.PlanRepository
	logger          *zap.Logger
	billingClient   billing.Client
	workRepo        malak.WorkspaceRepository
	preferencesRepo malak.PreferenceRepository
	taskQueue       queue.QueueHandler
	cfg             config.Config
}

func (s *stripeHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx, span, rid := getTracer(r.Context(), r, "stripe-webhook", s.cfg.Otel.IsEnabled)
	defer span.End()

	logger := s.logger.With(
		zap.String("request_id", rid),
		zap.String("method", "handleWebhook"))

	logger.Debug("handling stripe webhook")

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
		return
	}

	opts := webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	}

	ev, err := webhook.ConstructEventWithOptions(payload, r.Header.Get("Stripe-Signature"),
		s.cfg.Billing.Stripe.WebhookSecret, opts)
	if err != nil {
		logger.Error("could not verify webhook", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	switch ev.Type {
	case "customer.subscription.deleted":
		req := new(SubscriptionDeletedRequest)

		if err := json.Unmarshal(ev.Data.Raw, req); err != nil {
			_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
			return
		}

		s.handleExpiredSubscription(ctx, span, w, r, req, logger)

	case "customer.subscription.trial_will_end":
		req := new(TrialWillEnd)

		if err := json.Unmarshal(ev.Data.Raw, req); err != nil {
			_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
			return
		}

		s.sendTrialExpiringEmail(ctx, span, w, r, req, logger)

	case "invoice.paid":
		// if invoice paid, activate the subscription
		req := new(Invoice)

		if err := json.Unmarshal(ev.Data.Raw, req); err != nil {
			_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
			return
		}

		s.addInvoice(ctx, span, w, r, req, logger)
	case "customer.created":
		// creating a new customer will always add a free trial immediately

		req := new(CustomerCreatedEvent)
		if err := json.Unmarshal(ev.Data.Raw, req); err != nil {
			_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
			return
		}

		s.createFreeTrialSubscription(ctx, span, w, r, req, logger)

	default:
		_ = render.Render(w, r, newAPIStatus(http.StatusOK, "skipping this webhook"))
	}
}

func (s *stripeHandler) handleExpiredSubscription(ctx context.Context,
	_ trace.Span, w http.ResponseWriter,
	r *http.Request, req *SubscriptionDeletedRequest,
	logger *zap.Logger,
) {

	ctx, span := tracer.Start(ctx, "handleExpiredSubscription")
	defer span.End()

	logger = logger.With(zap.String("method", "handleExpiredSubscription"))

	logger.Debug("handling expired subscription")

	workspace, err := s.workRepo.Get(ctx, &malak.FindWorkspaceOptions{
		StripeCustomerID: req.Customer,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.Error("could not find workspace by stripe customer ID",
			zap.Error(err),
			zap.String("stripe_customer_id", req.Customer))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find workspace"))
		return
	}

	prefs, err := s.preferencesRepo.Get(ctx, workspace)
	if err != nil {
		logger.Error("could not fetch preferences", zap.Error(err))
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not fetch preferences"))
		return
	}

	workspace.IsSubscriptionActive = false
	workspace.SubscriptionID = ""

	// keep the old plan, when you resubscribe, the value updates

	if err := s.workRepo.Update(ctx, workspace); err != nil {
		logger.Error("could not update workspace to remove subscription", zap.Error(err))
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not update workspace susbcription"))
		return
	}

	err = s.taskQueue.Add(ctx, queue.QueueTopicSubscriptionExpired, &queue.SubscriptionExpiredOptions{
		Workspace: workspace,
		Recipient: prefs.Billing.FinanceEmail,
	})
	if err != nil {
		logger.Error("could not encode queue data", zap.Error(err))
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not encode queue data"))
		return
	}

	_ = render.Render(w, r, newAPIStatus(http.StatusOK, ""))
}

func (s *stripeHandler) sendTrialExpiringEmail(ctx context.Context,
	_ trace.Span, w http.ResponseWriter,
	r *http.Request, req *TrialWillEnd,
	logger *zap.Logger,
) {

	ctx, span := tracer.Start(ctx, "sendTrialExpiringEmail")
	defer span.End()

	logger = logger.With(zap.String("method", "sendTrialExpiringEmail"))

	logger.Debug("sending trial expiration email")

	workspace, err := s.workRepo.Get(ctx, &malak.FindWorkspaceOptions{
		StripeCustomerID: req.Customer,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.Error("could not find workspace by stripe customer ID",
			zap.Error(err),
			zap.String("stripe_customer_id", req.Customer))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find workspace"))
		return
	}

	prefs, err := s.preferencesRepo.Get(ctx, workspace)
	if err != nil {
		logger.Error("could not fetch preferences", zap.Error(err))
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not fetch preferences"))
		return
	}

	err = s.taskQueue.Add(ctx, queue.QueueTopicBillingTrialEnding, &queue.SendBillingTrialEmailOptions{
		Workspace:  workspace,
		Expiration: humanize.Time(time.Unix(req.TrialEnd, 0)),
		Recipient:  prefs.Billing.FinanceEmail,
	})
	if err != nil {
		logger.Error("could not encode queue data", zap.Error(err))
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not encode queue data"))
		return
	}

	_ = render.Render(w, r, newAPIStatus(http.StatusOK, ""))
}

func (s *stripeHandler) createFreeTrialSubscription(ctx context.Context,
	_ trace.Span, w http.ResponseWriter,
	r *http.Request, req *CustomerCreatedEvent,
	logger *zap.Logger,
) {

	ctx, span := tracer.Start(ctx, "createFreeTrialSubscription")
	defer span.End()

	logger = logger.With(zap.String("method", "createSubcriptionForCustomer"))

	logger.Debug("creating free subscription for user")

	workspace, err := s.workRepo.Get(ctx, &malak.FindWorkspaceOptions{
		StripeCustomerID: req.ID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.Error("could not find workspace by stripe customer ID",
			zap.Error(err),
			zap.String("stripe_customer_id", req.Data.Object.ID))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find workspace"))
		return
	}

	opts := &billing.AddPlanToCustomerOptions{
		Workspace: workspace,
	}

	subID, err := s.billingClient.AddPlanToCustomer(ctx, opts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.Error("could not create plan on stripe",
			zap.Error(err),
			zap.String("stripe_customer_id", workspace.StripeCustomerID))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not create subscription on Stripe"))
		return
	}

	workspace.SubscriptionID = subID

	if err := s.workRepo.Update(ctx, workspace); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.
			Error("could not update workspace stripe sub id",
				zap.String("workspace_id", workspace.ID.String()),
				zap.String("subscription_id", subID))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not update subscription on workspace"))
		return
	}

	_ = render.Render(w, r, newAPIStatus(http.StatusOK, ""))
}

func (s *stripeHandler) addInvoice(
	ctx context.Context, _ trace.Span, w http.ResponseWriter,
	r *http.Request, req *Invoice, logger *zap.Logger,
) {

	ctx, span := tracer.Start(ctx, "addInvoice")
	defer span.End()

	logger = logger.With(zap.String("method", "addInvoice"))

	logger.Debug("handling paid invoice")

	workspace, err := s.workRepo.Get(ctx, &malak.FindWorkspaceOptions{
		StripeCustomerID: req.Customer,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.
			Error("could not find workspace by stripe customer ID",
				zap.Error(err), zap.String("stripe_customer_id", req.Customer))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find workspace"))
		return
	}

	plan, err := s.planRepo.Get(ctx, &malak.FetchPlanOptions{
		Reference: req.Lines.Data[0].Plan.Product,
	})
	if err != nil {
		span.RecordError(errors.New("plan reference not match"))
		span.SetStatus(codes.Error, "plan reference not match")

		logger.
			Error("could not fetch plan",
				zap.Error(err),
				zap.String("workspace_id", workspace.ID.String()),
				zap.String("plan_id", workspace.PlanID.String()))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not fetch plan"))
		return
	}

	workspace.PlanID = plan.ID
	workspace.IsSubscriptionActive = true
	workspace.SubscriptionID = req.Subscription

	if err := s.workRepo.Update(ctx, workspace); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		logger.
			Error("could not mark subscription as active",
				zap.Error(err),
				zap.String("workspace_id", workspace.ID.String()),
				zap.String("subscription_id", workspace.SubscriptionID))

		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not update subscription"))
		return
	}

	_ = render.Render(w, r, newAPIStatus(http.StatusOK, ""))
}
