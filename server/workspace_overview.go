package server

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// @Description fetch workspace overview
// @Tags workspaces
// @Accept  json
// @Produce  json
// @Success 200 {object} workspaceOverviewResponse
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

	workspace := getWorkspaceFromContext(ctx)

	g, gctx := errgroup.WithContext(ctx)

	var (
		updates  *malak.UpdateOverview
		decks    *malak.DeckOverview
		contacts *malak.ContactOverview
		shares   *malak.ShareOverview
	)

	g.Go(func() error {
		var err error
		updates, err = wo.updateRepo.Overview(gctx, workspace.ID)
		if err != nil {
			logger.Error("could not fetch updates overview", zap.Error(err))
		}
		return err
	})

	g.Go(func() error {
		var err error
		decks, err = wo.deckRepo.Overview(gctx, workspace.ID)
		if err != nil {
			logger.Error("could not fetch decks overview", zap.Error(err))
		}
		return err
	})

	g.Go(func() error {
		var err error
		contacts, err = wo.contactRepo.Overview(gctx, workspace.ID)
		if err != nil {
			logger.Error("could not fetch contacts overview", zap.Error(err))
		}
		return err
	})

	g.Go(func() error {
		var err error
		shares, err = wo.shareRepo.Overview(gctx, workspace.ID)
		if err != nil {
			logger.Error("could not fetch shares overview", zap.Error(err))
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return newAPIStatus(http.StatusInternalServerError, "could not fetch workspace overview"), StatusFailed
	}

	return workspaceOverviewResponse{
		Updates:   updates,
		Decks:     decks,
		Contacts:  contacts,
		Shares:    shares,
		APIStatus: newAPIStatus(http.StatusOK, "workspace overview fetched"),
	}, StatusSuccess
}
