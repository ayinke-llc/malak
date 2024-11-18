package server

import (
	"net/http"

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

func (we *webhookHandler) handleResend(
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger.Debug("Process resend webhook")
	}
}
