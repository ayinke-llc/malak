package server

import (
	"net/http"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
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

		// ctx, span, rid := getTracer(r.Context(), r, "resend.webhookHandler", we.cfg.Otel.IsEnabled)
		// defer span.End()
		//
		// logger = logger.With(zap.String("request_id", rid))
		//
		// logger.Debug("Process resend webhook")
		//
		// var req = new(resendWebhookRequest)
		//
		// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		// 	_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "invalid/unexpected resend body"))
		// 	return
		// }
		//
		// contactEmail := req.Data.To[0]
		// if hermes.IsStringEmpty(contactEmail) {
		// 	_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "no email provided"))
		// 	return
		// }
		//
		// contact, err := we.contactRepo.Get(ctx, malak.FetchContactOptions{
		// 	Email: malak.Email(contactEmail),
		// })
		// if err != nil {
		// 	_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "could not fetch contact"))
		// 	return
		// }
		//
		// switch req.Type {
		// case "opened":
		// 	logger.Info("Handling 'opened' event", zap.String("email_id", req.Data.EmailID))
		// 	// Process opened event logic here
		// case "bounced":
		// 	logger.Info("Handling 'bounced' event", zap.String("email_id", req.Data.EmailID))
		// 	// Process bounced event logic here
		// case "delayed":
		// 	logger.Info("Handling 'delayed' event", zap.String("email_id", req.Data.EmailID))
		// 	// Process delayed event logic here
		// case "delivered":
		// 	logger.Info("Handling 'delivered' event", zap.String("email_id", req.Data.EmailID))
		// 	// Process delivered event logic here
		// default:
		// 	logger.Warn("Unhandled event type", zap.String("type", req.Type))
		// 	_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "unsupported event type"))
		// 	return
		// }
	}
}
