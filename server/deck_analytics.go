package server

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"net/netip"
	"time"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type createDeckViewerSession struct {
	OS         string         `json:"os,omitempty" validate:"required"`
	DeviceInfo string         `json:"device_info,omitempty" validate:"required"`
	Password   malak.Password `json:"password,omitempty" validate:"required"`
	Browser    string         `json:"browser,omitempty" validate:"required"`

	GenericRequest
}

func (c *createDeckViewerSession) Validate() error {
	if hermes.IsStringEmpty(c.OS) {
		return errors.New("provide operating system of viewer")
	}

	if hermes.IsStringEmpty(c.DeviceInfo) {
		return errors.New("provide device information")
	}

	if hermes.IsStringEmpty(c.DeviceInfo) {
		return errors.New("provide browser information")
	}

	return nil
}

// @Description public api to fetch a deck
// @Tags decks-viewer
// @Accept  json
// @Produce  json
// @Param reference path string required "deck unique reference.. "
// @Param message body createDeckViewerSession true "deck session request body"
// @Success 200 {object} fetchPublicDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /public/decks/{reference} [post]
func (d *deckHandler) publicDeckDetails(
	ctx context.Context,
	_ trace.Span,
	logger *zap.Logger,
	_ http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetching deck public resource")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	logger = logger.With(zap.String("reference", ref))

	req := new(createDeckViewerSession)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	ipAddr := hermes.GetIP(r)

	deck, err := d.deckRepo.PublicDetails(ctx, malak.Reference(ref))
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

	if deck.IsArchived {
		return newAPIStatus(http.StatusBadRequest, "deck not available as it is archived"),
			StatusFailed
	}

	var objectLink string
	var country, city string
	var contactID uuid.UUID

	var g errgroup.Group

	g.Go(func() error {
		var err error

		objectLink, err = d.gulterStore.Path(ctx, gulter.PathOptions{
			Key:            deck.ObjectKey,
			ExpirationTime: time.Minute * 15,
			IsSecure:       true,
		})

		return err
	})

	g.Go(func() error {
		var err error

		ip, err := netip.ParseAddr(ipAddr.String())
		if err != nil {
			return err
		}

		country, city, err = d.geolocationService.FindByIP(ctx, ip)
		return err
	})

	if err := g.Wait(); err != nil {
		logger.Error("could not process deck details", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not find path to deck"),
			StatusFailed
	}

	sessionReq := &malak.DeckViewerSession{
		DeckID:     deck.ID,
		Reference:  d.referenceGenerator.Generate(malak.EntityTypeDeckViewerSession),
		SessionID:  malak.Reference(d.referenceGenerator.Generate(malak.EntityTypeSession)),
		DeviceInfo: req.DeviceInfo,
		OS:         req.OS,
		Browser:    req.Browser,
		IPAddress:  ipAddr.String(),
		Country:    country,
		City:       city,
	}

	if contactID != uuid.Nil {
		sessionReq.ContactID = contactID
	}

	if err := d.deckRepo.CreateDeckSession(ctx, sessionReq); err != nil {
		logger.Error("could not create deck session", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create deck session"),
			StatusFailed
	}

	return fetchPublicDeckResponse{
		APIStatus: newAPIStatus(http.StatusOK, "fetched deck details"),
		Deck: malak.PublicDeck{
			Session:     hermes.DeRef(sessionReq),
			Reference:   malak.Reference(ref),
			WorkspaceID: deck.WorkspaceID,
			Title:       deck.Title,
			ShortLink:   deck.ShortLink,
			DeckSize:    deck.DeckSize,
			IsArchived:  deck.IsArchived,
			CreatedAt:   deck.CreatedAt,
			UpdatedAt:   deck.UpdatedAt,
			DeckPreference: &malak.PublicDeckPreference{
				EnableDownloading: deck.DeckPreference.EnableDownloading,
				RequireEmail:      deck.DeckPreference.RequireEmail,
				HasPassword:       deck.DeckPreference.Password.Enabled,
			},
			ObjectLink: objectLink,
		},
	}, StatusSuccess
}

type updateDeckViewerSession struct {
	Email     malak.Email    `json:"email,omitempty" validate:"optional"`
	Password  malak.Password `json:"password,omitempty" validate:"optional"`
	TimeSpent int64          `json:"time_spent,omitempty" validate:"optional"`
	SessionID string         `json:"session_id,omitempty" validate:"optional"`

	GenericRequest
}

func (c *updateDeckViewerSession) Validate() error {
	if !hermes.IsStringEmpty(c.Email.String()) {
		_, err := mail.ParseAddress(c.Email.String())
		if err != nil {
			return errors.New("please provide a valid email")
		}
	}

	return nil
}

// @Description update the session details
// @Tags decks-viewer
// @Accept  json
// @Produce  json
// @Param reference path string required "session unique reference.. "
// @Param message body createDeckViewerSession true "deck session request body"
// @Success 200 {object} fetchPublicDeckResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /public/decks/{reference} [put]
func (d *deckHandler) updateDeckViewerSession(
	ctx context.Context,
	_ trace.Span,
	logger *zap.Logger,
	_ http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("updating deck session viewing")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	logger = logger.With(zap.String("reference", ref))

	req := new(updateDeckViewerSession)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	logger = logger.With(zap.String("session_id", req.SessionID))

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	deck, err := d.deckRepo.PublicDetails(ctx, malak.Reference(ref))
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

	if deck.IsArchived {
		return newAPIStatus(http.StatusBadRequest, "deck not available as it is archived"),
			StatusFailed
	}

	opts := &malak.UpdateDeckSessionOptions{}

	if !hermes.IsStringEmpty(req.Email.String()) {
		contact, err := d.contactRepo.Get(ctx, malak.FetchContactOptions{
			Email:       req.Email,
			WorkspaceID: deck.WorkspaceID,
		})

		if err != nil {
			if errors.Is(err, malak.ErrContactNotFound) {
				opts.CreateContact = true
				opts.Contact = &malak.Contact{
					Email:       req.Email,
					FirstName:   req.Email.String(),
					Metadata:    make(malak.CustomContactMetadata),
					WorkspaceID: deck.WorkspaceID,
					Reference:   d.referenceGenerator.Generate(malak.EntityTypeContact),
				}

			} else {
				logger.Error("could not fetch contat from dtabase", zap.Error(err))
				return newAPIStatus(http.StatusInternalServerError, "an error occurred while finding contact by email"),
					StatusFailed
			}
		} else {
			opts.Contact = contact
		}
	}

	if !hermes.IsStringEmpty(string(req.Password)) {
		if !malak.VerifyPassword(string(deck.DeckPreference.Password.Password), string(req.Password)) {
			return newAPIStatus(http.StatusBadRequest, "deck password not correct"), StatusFailed
		}
	}

	session, err := d.deckRepo.FindDeckSession(ctx, req.SessionID)
	if err != nil {
		logger.Error("could not find deck session")
		return newAPIStatus(http.StatusInternalServerError, "deck session not found"), StatusFailed
	}

	session.TimeSpentSeconds = req.TimeSpent

	opts.Session = session

	if err := d.deckRepo.UpdateDeckSession(ctx, opts); err != nil {
		logger.Error("could not create deck session", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create deck session"),
			StatusFailed
	}

	return newAPIStatus(http.StatusOK, "fetched deck details"), StatusSuccess
}
