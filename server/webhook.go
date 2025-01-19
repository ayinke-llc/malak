package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type webhookHandler struct {
	cfg                config.Config
	userRepo           malak.UserRepository
	workspaceRepo      malak.WorkspaceRepository
	planRepo           malak.PlanRepository
	updateRepo         malak.UpdateRepository
	contactRepo        malak.ContactRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type resendWebhookRequest struct {
	CreatedAt time.Time `json:"created_at"`
	Data      struct {
		CreatedAt string   `json:"created_at"`
		EmailID   string   `json:"email_id"`
		From      string   `json:"from"`
		Subject   string   `json:"subject"`
		To        []string `json:"to"`
	} `json:"data"`
	Type string `json:"type"`
}

func (we *webhookHandler) handleResend(
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span, rid := getTracer(r.Context(), r, "resend.webhookHandler", we.cfg.Otel.IsEnabled)
		defer span.End()

		logger = logger.With(zap.String("request_id", rid))

		logger.Debug("Process resend webhook")

		var req = new(resendWebhookRequest)

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "invalid/unexpected resend body"))
			return
		}

		_, recipientStat, err := we.updateRepo.GetStatByEmailID(ctx, req.Data.EmailID, malak.UpdateRecipientLogProviderResend)
		if err != nil {
			logger.Error("could not fetch recipient by id", zap.Error(err),
				zap.String("email_reference", req.Data.EmailID))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find recipient"))
			return
		}

		update := &malak.Update{
			ID: recipientStat.Recipient.UpdateID,
		}

		updateStat, err := we.updateRepo.Stat(ctx, update)
		if err != nil {
			logger.Error("could not fetch update stats by id", zap.Error(err),
				zap.String("update_id", update.ID.String()))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find update stat"))
			return
		}

		switch req.Type {
		case "email.opened":
			updateStat.TotalOpens++
			if recipientStat.LastOpenedAt == nil {
				updateStat.UniqueOpens++
			}

			recipientStat.LastOpenedAt = hermes.Ref(time.Now())

		case "email.bounced":
			recipientStat.IsBounced = true

		case "email.delivered":
			recipientStat.IsDelivered = true

		default:
			_ = render.Render(w, r, newAPIStatus(http.StatusOK, "unsupported event type"))
			return
		}

		if err := we.updateRepo.UpdateStat(ctx, updateStat, recipientStat); err != nil {
			logger.Error("could not update stat", zap.Error(err),
				zap.String("recipient_stat", recipientStat.ID.String()),
				zap.String("update_stat_id", update.ID.String()))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not update stat"))
			return
		}

		_ = render.Render(w, r, newAPIStatus(http.StatusOK, "processed stat"))
	}
}
