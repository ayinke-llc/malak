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

// @Description Duplicate a specific update
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
		Reference:   malak.Reference(ref),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
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

	opts := &malak.TemplateCreateUpdateOptions{
		IsSystemTemplate: false,
	}

	if err := u.updateRepo.Create(ctx, newUpdate, opts); err != nil {
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
// @Description Delete a specific update
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
		Reference:   malak.Reference(ref),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
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
// @Description Toggle pinned status a specific update
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
		Reference:   malak.Reference(ref),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
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

	if err := u.updateRepo.TogglePinned(ctx, update); err != nil {
		logger.Error("could not toggle pinned state of update", zap.Error(err))

		var msg = "could not toggle pinned status"
		status := http.StatusInternalServerError
		if errors.Is(err, malak.ErrPinnedUpdateCapacityExceeded) {
			msg = err.Error()
			status = http.StatusBadRequest
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	return createdUpdateResponse{
		Update:    util.DeRef(update),
		APIStatus: newAPIStatus(http.StatusOK, "Pinned status updated"),
	}, StatusSuccess
}

// @Tags updates
// @Description Fetch a specific update
// @id fetchUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} fetchUpdateReponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference} [get]
func (u *updatesHandler) fetchUpdate(
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
		Reference:   malak.Reference(ref),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
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

// @Tags updates
// @Description Fetch a specific update
// @Id reactPost
// @Accept  json
// @Produce  json
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Param provider query string true "provider type"
// @Param email_id query string true "email id"
// @Router /updates/react [get]
// ?provider=resend&email_id=xxxx
func (u *updatesHandler) handleReaction(
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span, rid := getTracer(r.Context(), r, "updates.reaction.handler", u.cfg.Otel.IsEnabled)
		defer span.End()

		logger = logger.With(zap.String("request_id", rid))

		ref := chi.URLParam(r, "reference")

		provider := malak.UpdateRecipientLogProvider(r.URL.Query().Get("provider"))
		emailID := r.URL.Query().Get("email_id")

		span.SetAttributes(attribute.String("id", ref))

		logger = logger.With(zap.String("reference", ref))

		logger.Debug("reacting to update")

		_, recipientStat, err := u.updateRepo.GetStatByEmailID(ctx, emailID, provider)
		if err != nil {
			logger.Error("could not fetch recipient by id", zap.Error(err),
				zap.String("email_reference", emailID))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find recipient"))
			return
		}

		update := &malak.Update{
			ID: recipientStat.Recipient.UpdateID,
		}

		updateStat, err := u.updateRepo.Stat(ctx, update)
		if err != nil {
			logger.Error("could not fetch update stats by id", zap.Error(err),
				zap.String("update_id", update.ID.String()))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not find update stat"))
			return
		}

		recipientStat.HasReaction = true
		updateStat.TotalReactions += 1

		if err := u.updateRepo.UpdateStat(ctx, updateStat, recipientStat); err != nil {
			logger.Error("could not update stat", zap.Error(err),
				zap.String("recipient_stat", recipientStat.ID.String()))
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "could not update reaction"))
			return
		}

		_ = render.Render(w, r, newAPIStatus(http.StatusOK, "added reaction"))
	}
}
