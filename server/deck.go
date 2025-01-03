package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type deckHandler struct {
	deckRepo           malak.DeckRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createDeckRequest struct {
	GenericRequest

	Title   string `json:"title,omitempty"`
	DeckURL string `json:"deck_url,omitempty"`
}

func (c *createDeckRequest) Validate() error {
	if hermes.IsStringEmpty(c.DeckURL) {
		return errors.New("please provide the deck url")
	}

	if hermes.IsStringEmpty(c.Title) {
		return errors.New("please provide the title of the deck")
	}

	p := bluemonday.StrictPolicy()

	c.Title = p.Sanitize(c.Title)

	return nil
}

// @Summary Creates a new deck
// @Tags decks
// @Accept  json
// @Produce  json
// @Param message body createDeckRequest true "deck request body"
// @Success 200 {object} fetchDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks [post]
func (d *deckHandler) Create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating deck")

	req := new(createDeckRequest)

	workspace := getWorkspaceFromContext(r.Context())

	user := getUserFromContext(r.Context())

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	opts := &malak.CreateDeckOptions{
		RequireEmail:      true,
		EnableDownloading: false,
		Password: struct {
			Enabled  bool           "json:\"enabled,omitempty\" validate:\"required\""
			Password malak.Password "json:\"password,omitempty\" validate:\"required\""
		}{
			Enabled: false,
		},
		Reference: d.referenceGenerator.Generate(malak.EntityTypeDeckPreference),
	}

	deck := &malak.Deck{
		Title:       req.Title,
		ShortLink:   d.referenceGenerator.ShortLink(),
		Reference:   d.referenceGenerator.Generate(malak.EntityTypeDeck),
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
	}

	if err := d.deckRepo.Create(ctx, deck, opts); err != nil {
		logger.Error("could not create deck",
			zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create deck"),
			StatusFailed
	}

	return fetchDeckResponse{
		Deck:      hermes.DeRef(deck),
		APIStatus: newAPIStatus(http.StatusOK, "deck created"),
	}, StatusSuccess
}

// @Summary list all decks. No pagination
// @Tags decks
// @Accept  json
// @Produce  json
// @Success 200 {object} fetchDecksResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks [get]
func (d *deckHandler) List(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("listing deck")

	workspace := getWorkspaceFromContext(r.Context())

	decks, err := d.deckRepo.List(ctx, workspace)
	if err != nil {
		logger.Error("could not list decks", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not list decks"),
			StatusFailed
	}

	return fetchDecksResponse{
		Decks:     decks,
		APIStatus: newAPIStatus(http.StatusOK, "fetched your decks"),
	}, StatusSuccess
}

// @Summary delete a deck
// @Tags decks
// @Accept  json
// @Produce  json
// @Param reference path string required "deck unique reference.. e.g deck_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks/{reference} [delete]
func (d *deckHandler) Delete(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("deleting deck")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	deck, err := d.deckRepo.Get(ctx, malak.FetchDeckOptions{
		Reference: ref,
	})
	if err != nil {
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching deck"

		if errors.Is(err, malak.ErrDeckNotFound) {
			status = http.StatusNotFound
			msg = "deck does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	if err := d.deckRepo.Delete(ctx, deck); err != nil {
		logger.Error("error occurred while deleting the deck",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "error occurred while deleting deck"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "deleted your deck"), StatusSuccess
}
