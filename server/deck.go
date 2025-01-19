package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"time"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type deckHandler struct {
	deckRepo           malak.DeckRepository
	referenceGenerator malak.ReferenceGeneratorOperation
	cache              cache.Cache
	cfg                config.Config
}

func hashURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	normalizedURL := parsedURL.String()

	hasher := fnv.New64a()
	hasher.Write([]byte(normalizedURL))

	return "deck-" + fmt.Sprintf("%x", hasher.Sum64()), nil
}

// @Summary Upload a deck
// @Tags decks
// @id uploadDeck
// @Accept  json
// @Produce  json
// @Param image_body formData file true "image body"
// @Success 200 {object} uploadImageResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /uploads/decks [post]
func (u *deckHandler) uploadImage(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("file uploaded using Gulter")

	files, err := gulter.FilesFromContextWithKey(r, "image_body")
	if err != nil {
		logger.Error("could not fetch gulter uploaded files", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal failure while fetching file from storage"), StatusFailed
	}

	// only one file we are expecting at a time
	file := files[0]

	uploadedURL := fmt.Sprintf("%s/%s/%s",
		u.cfg.Uploader.S3.Endpoint,
		file.FolderDestination,
		file.UploadedFileName)

	// Add the details of the file to redis
	// So when it is being addded to a deck, you
	// can just fetch the size directly and move on
	// This is definitely:
	// - better than relying on client to send the size
	// - simpler than fetching/downloading the file when creating the deck
	// in other to fetch the size
	//
	// Simple but works
	cacheKey, err := hashURL(uploadedURL)
	if err != nil {
		logger.Error("could not create hash key", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not create cache key"), StatusFailed
	}

	var f = struct {
		Size int64
	}{
		Size: file.Size,
	}

	var b = bytes.NewBuffer(nil)

	if err := json.NewEncoder(b).Encode(&f); err != nil {
		logger.Error("could not encode size", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal error"), StatusFailed
	}

	if err := u.cache.Add(ctx, cacheKey, b.Bytes(), time.Hour*4); err != nil {
		logger.Error("could not add to cache", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal error"), StatusFailed
	}

	return uploadImageResponse{
		URL:       uploadedURL,
		APIStatus: newAPIStatus(http.StatusOK, "deck was uploaded"),
	}, StatusSuccess
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

	// get the file size details
	cacheKey, err := hashURL(req.DeckURL)
	if err != nil {
		logger.Error("could not fetch file details from redis", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal error"), StatusFailed
	}

	data, err := d.cache.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("could not fetch cache details from redis", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not fetch size of file. Reupload file"), StatusFailed
	}

	var file struct {
		Size int64
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&file); err != nil {
		logger.Error("could not decode file size from Redis", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal error while getting size of file"), StatusFailed
	}

	if file.Size <= 0 {
		logger.Error("file size is negative", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"file size is invalid. Try uploading another file"), StatusFailed
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
		DeckSize:    file.Size,
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
		Reference:   ref,
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
	})
	if err != nil {
		logger.Error("could not fetch deck", zap.Error(err))
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

// @Summary fetch a deck
// @Tags decks
// @Accept  json
// @Produce  json
// @Param reference path string required "deck unique reference.. e.g deck_"
// @Success 200 {object} fetchDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks/{reference} [get]
func (d *deckHandler) fetch(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetching deck")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	deck, err := d.deckRepo.Get(ctx, malak.FetchDeckOptions{
		Reference:   ref,
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
	})
	if err != nil {
		logger.Error("could not fetch deck", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching deck"

		if errors.Is(err, malak.ErrDeckNotFound) {
			status = http.StatusNotFound
			msg = "deck does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	return fetchDeckResponse{
		APIStatus: newAPIStatus(http.StatusOK, "fetched deck details"),
		Deck:      hermes.DeRef(deck),
	}, StatusSuccess
}

type updateDeckPreferencesRequest struct {
	GenericRequest

	EnableDownloading  bool `json:"enable_downloading,omitempty"`
	RequireEmail       bool `json:"require_email,omitempty"`
	PasswordProtection struct {
		Enabled bool           `json:"enabled,omitempty"`
		Value   malak.Password `json:"value,omitempty"`
	} `json:"password_protection,omitempty"`
}

func (u *updateDeckPreferencesRequest) Validate() error {

	if u.PasswordProtection.Enabled {
		if u.PasswordProtection.Value.IsZero() {
			return errors.New("please provide your password")
		}
	}

	return nil
}

// @Summary update a deck preferences
// @Tags decks
// @Accept  json
// @Produce  json
// @Param reference path string required "deck unique reference.. e.g deck_"
// @Param message body updateDeckPreferencesRequest true "deck preferences request body"
// @Success 200 {object} fetchDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks/{reference}/preferences [put]
// TODO: make this a PATCH?
func (d *deckHandler) updatePreferences(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("updating deck preferences")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	req := new(updateDeckPreferencesRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	deck, err := d.deckRepo.Get(ctx, malak.FetchDeckOptions{
		Reference:   ref,
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
	})
	if err != nil {
		logger.Error("could not fetch deck", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching deck"

		if errors.Is(err, malak.ErrDeckNotFound) {
			status = http.StatusNotFound
			msg = "deck does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	deck.DeckPreference.Password = malak.PasswordDeckPreferences{
		Enabled:  req.PasswordProtection.Enabled,
		Password: req.PasswordProtection.Value,
	}

	deck.DeckPreference.EnableDownloading = req.EnableDownloading
	deck.DeckPreference.RequireEmail = req.RequireEmail

	if err := d.deckRepo.UpdatePreferences(ctx, deck); err != nil {
		logger.Error("could not update preferences", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not update preferences"),
			StatusFailed
	}

	return fetchDeckResponse{
		APIStatus: newAPIStatus(http.StatusOK, "Updated deck preferences"),
		Deck:      hermes.DeRef(deck),
	}, StatusSuccess
}

// @Summary toggle archive status of a deck
// @Tags decks
// @Accept  json
// @Produce  json
// @id toggleArchive
// @Param reference path string required "deck unique reference.. e.g deck_"
// @Success 200 {object} fetchDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /decks/{reference}/archive [post]
func (d *deckHandler) toggleArchive(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("toggling archive status of a deck")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	deck, err := d.deckRepo.Get(ctx, malak.FetchDeckOptions{
		Reference:   ref,
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
	})
	if err != nil {
		logger.Error("could not fetch deck", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching deck"

		if errors.Is(err, malak.ErrDeckNotFound) {
			status = http.StatusNotFound
			msg = "deck does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	if err := d.deckRepo.ToggleArchive(ctx, deck); err != nil {
		logger.Error("could not toggle archive status", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not toggle archival status"),
			StatusFailed
	}

	return fetchDeckResponse{
		APIStatus: newAPIStatus(http.StatusOK, "Updated deck archival status"),
		Deck:      hermes.DeRef(deck),
	}, StatusSuccess
}
