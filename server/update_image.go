package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Summary Upload an image
// @Tags updates
// @id uploadImage
// @Accept  json
// @Produce  json
// @Success 200 {object} uploadImageResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/image [post]
func (u *updatesHandler) uploadImage(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating a new update")

	user := getUserFromContext(r.Context())
	workspace := getWorkspaceFromContext(r.Context())

	return uploadImageResponse{
		URL:       "",
		APIStatus: newAPIStatus(http.StatusCreated, "image was uploaded"),
	}, StatusSuccess
}
