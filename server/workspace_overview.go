package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Description fetch overview
// @Tags integrations
// @Accept  json
// @Produce  json
// @Param message body createChartRequest true "request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/overview [get]
func (wo *workspaceHandler) overview(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetching overview data")

	return newAPIStatus(http.StatusCreated, "created chart for integration"),
		StatusSuccess
}
