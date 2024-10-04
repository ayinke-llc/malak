package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Tags updates
// @Summary Send preview of an update
// @id previewUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference}/preview [post]
func (u *updatesHandler) previewUpdate(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Fetching update")

	update, err := u.updateRepo.Get(ctx, malak.FetchUpdateOptions{
		Reference: malak.Reference(ref),
	})
	if errors.Is(err, malak.ErrUpdateNotFound) {
		return newAPIStatus(http.StatusNotFound,
			"update does not exists"), StatusFailed
	}

	if err != nil {
		logger.Error("could not fetch update", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"an error occurred while fetching update"), StatusFailed
	}

	return fetchUpdateReponse{
		APIStatus: newAPIStatus(http.StatusOK, "update fetched"),
		Update:    util.DeRef(update),
	}, StatusSuccess
}
