package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Summary fetch workspace preferences
// @Tags workspace
// @Accept  json
// @Produce  json
// @Success 200 {object} listIntegrationResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations [get]
func (wo *workspaceHandler) getIntegrations(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request,
) (render.Renderer, Status) {

	logger.Debug("fetching workspace integrations")

	workspace := getWorkspaceFromContext(ctx)

	integrations, err := wo.integrationRepo.List(ctx, workspace)
	if err != nil {
		logger.Error("could not list integrations", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not list integrations"), StatusFailed
	}

	return &listIntegrationResponse{
		Integrations: integrations,
		APIStatus:    newAPIStatus(http.StatusOK, "workspace integrations retrieved"),
	}, StatusSuccess
}
