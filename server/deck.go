package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type deckHandler struct {
	cfg                config.Config
	deckRepo           malak.DeckRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createDeckRequest struct {
	GenericRequest

	Title             string `json:"title,omitempty"`
	DeckURL           string `json:"deck_url,omitempty"`
	RequireEmail      bool   `json:"require_email,omitempty"`
	EnableDownloading bool   `json:"enable_downloading,omitempty"`

	Password struct {
		Enabled  bool            `json:"enabled,omitempty" validate:"required"`
		Password *malak.Password `json:"password,omitempty" validate:"required"`
	} `json:"password,omitempty" validate:"required"`

	ExpiresAt *time.Time `json:"expires_at,omitempty"`
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

	if c.Password.Enabled {
		password := hermes.DeRef(c.Password.Password)
		if len(string(password)) < 5 {
			return errors.New("password must be 5 characters or more")
		}
	}

	if c.ExpiresAt != nil {
		exp := hermes.DeRef(c.ExpiresAt)
		if exp.Before(time.Now()) {
			return errors.New("expiration date cannot be in the past")
		}
	}

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

	var pass string

	if len(string(hermes.DeRef(req.Password.Password))) >= 0 {
		pass = string(hermes.DeRef(req.Password.Password))
	}

	opts := &malak.CreateDeckOptions{
		RequireEmail:      req.RequireEmail,
		EnableDownloading: req.EnableDownloading,
		Password: struct {
			Enabled  bool           "json:\"enabled,omitempty\" validate:\"required\""
			Password malak.Password "json:\"password,omitempty\" validate:\"required\""
		}{
			Enabled:  req.Password.Enabled,
			Password: malak.Password(pass),
		},
		ExpiresAt: req.ExpiresAt,
	}

	deck := &malak.Deck{
		Title:       req.Title,
		ShortLink:   malak.ShortLink(),
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
