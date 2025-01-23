package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// @Tags updates
// @Summary Fetch analytics for a specific update
// @id fetchUpdateAnalytics
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} fetchUpdateAnalyticsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference}/analytics [get]
func (u *updatesHandler) fetchUpdateAnalytics(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Fetching analytics update")

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

	var stat *malak.UpdateStat
	var recipients []malak.UpdateRecipient

	var g errgroup.Group

	g.Go(func() error {
		var err error

		recipients, err = u.updateRepo.RecipientStat(ctx, update)
		if err != nil {
			logger.Error("could not fetch recipient stats", zap.Error(err))
			return err
		}

		return nil
	})

	g.Go(func() error {
		var err error

		stat, err = u.updateRepo.Stat(ctx, update)
		if err != nil {
			logger.Error("could not fetch update stat", zap.Error(err))
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return newAPIStatus(http.StatusInternalServerError, "could not fetch analytics"),
			StatusFailed
	}

	return fetchUpdateAnalyticsResponse{
		APIStatus:  newAPIStatus(http.StatusOK, "update fetched"),
		Recipients: recipients,
		Update:     hermes.DeRef(stat),
	}, StatusSuccess
}
