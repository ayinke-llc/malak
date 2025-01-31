package stripe

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/billing"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

var tracer = otel.Tracer("billing.stripe")

type stripeClient struct {
	stripeClient *client.API
	cfg          config.Config
}

func New(cfg config.Config) (billing.Client, error) {

	stripeLib := &client.API{}

	config := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(0), // Zero retries
		EnableTelemetry:   hermes.Ref(false),
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}

	stripeLib.Init(cfg.Billing.Stripe.APIKey, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, config),
	})

	return &stripeClient{
		stripeClient: stripeLib,
		cfg:          cfg,
	}, nil
}

func (s *stripeClient) CreateCustomer(ctx context.Context,
	opts *billing.CreateCustomerOptions) (string, error) {

	ctx, span := tracer.Start(ctx, "customer.create")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(http.Header{}))

	customer, err := s.stripeClient.Customers.New(&stripe.CustomerParams{
		Email: hermes.Ref(opts.Email.String()),
		Name:  hermes.Ref(opts.Name),
	})
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "could not create customers")

		return "", err
	}

	span.SetStatus(codes.Ok, "created customer")

	return customer.ID, nil
}

func (s *stripeClient) Portal(ctx context.Context,
	opts *billing.CreateBillingPortalOptions) (string, error) {

	ctx, span := tracer.Start(ctx, "billing.portal")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(http.Header{}))

	portal, err := s.stripeClient.BillingPortalSessions.New(&stripe.BillingPortalSessionParams{
		Customer: hermes.Ref(opts.CustomerID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not create billing portal")

		return "", err
	}

	span.SetStatus(codes.Ok, "created billing portal")

	return portal.URL, nil
}

func (s *stripeClient) AddPlanToCustomer(ctx context.Context,
	opts *billing.AddPlanToCustomerOptions) (string, error) {

	ctx, span := tracer.Start(ctx, "subscription.create")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(http.Header{}))

	createSubscriptionOptions := &stripe.SubscriptionParams{
		Customer: hermes.Ref(opts.Workspace.StripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: hermes.Ref(opts.Workspace.Plan.DefaultPriceID),
			},
		},
		TrialPeriodDays: hermes.Ref(int64(s.cfg.Billing.TrialDays)),
	}

	sub, err := s.stripeClient.Subscriptions.New(createSubscriptionOptions)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not create subscriptions")

		return "", err
	}

	span.SetStatus(codes.Ok, "created subscription")
	return sub.ID, nil
}
