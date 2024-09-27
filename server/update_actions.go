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

// @Summary Duplicate a specific update
// @Tags updates
// @id duplicateUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} createdUpdateResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference}/duplicate [post]
func (u *updatesHandler) duplicate(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	user := getUserFromContext(ctx)

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Duplicating update")

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

	newUpdate := &malak.Update{
		Content:     update.Content,
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeUpdate),
		CreatedBy:   user.ID,
		Status:      malak.UpdateStatusDraft,
		WorkspaceID: update.WorkspaceID,
	}

	if err := u.updateRepo.Create(ctx, newUpdate); err != nil {
		logger.Error("could not create updates", zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError,
			"could not create updates"), StatusFailed
	}

	return createdUpdateResponse{
		Update:    util.DeRef(newUpdate),
		APIStatus: newAPIStatus(http.StatusCreated, "update has been duplicated"),
	}, StatusSuccess
}

// @Tags updates
// @Summary Delete a specific update
// @id deleteUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference} [delete]
func (u *updatesHandler) delete(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Deleting update")

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

	if err := u.updateRepo.Delete(ctx, update); err != nil {
		logger.Error("could not create updates", zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError,
			"could not create updates"), StatusFailed
	}

	return newAPIStatus(http.StatusOK,
			"update has been deleted"),
		StatusSuccess
}

// @Tags updates
// @Summary Toggle pinned status a specific update
// @id toggleUpdatePin
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} createdUpdateResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference}/pin [post]
func (u *updatesHandler) togglePinned(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Toggling update's pin status")

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

	// simplistic, we are not really working with realtime data where you want
	// to offload to db directly with NOT
	update.IsPinned = !update.IsPinned

	if err := u.updateRepo.Update(ctx, update); err != nil {
		logger.Error("could not toggle pinned state of update", zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError,
			"could not toggle pinned status"), StatusFailed
	}

	var msg = "update has been pinned"
	if !update.IsPinned {
		msg = "update has been unpinned"
	}

	return createdUpdateResponse{
		Update:    util.DeRef(update),
		APIStatus: newAPIStatus(http.StatusOK, msg),
	}, StatusSuccess
}
