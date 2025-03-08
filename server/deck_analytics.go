package server

import (
	"context"
	"errors"
	"net"
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
	IPAddress  string         `json:"ip_address,omitempty" validate:"required"`
	Email      malak.Email    `json:"email,omitempty" validate:"required"`
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

	if !hermes.IsStringEmpty(c.Email.String()) {
		_, err := mail.ParseAddress(c.Email.String())
		if err != nil {
			return err
		}
	}

	ip := net.ParseIP(c.IPAddress)
	if ip != nil {
		return errors.New("please provide a valid IP")
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

	var objectLink string
	var country, city string
	var contactID uuid.UUID

	var g errgroup.Group

	g.Go(func() error {
		var err error

		objectLink, err = d.gulterStore.Path(ctx, gulter.PathOptions{
			Bucket:         d.cfg.Uploader.S3.DeckBucket,
			Key:            deck.ObjectKey,
			ExpirationTime: time.Minute * 15,
			IsSecure:       true,
		})

		return err
	})

	g.Go(func() error {
		var err error

		ip, err := netip.ParseAddr(req.IPAddress)
		if err != nil {
			return err
		}

		country, city, err = d.geolocationService.FindByIP(ctx, ip)
		return err
	})

	g.Go(func() error {
		var err error

		contact, err := d.contactRepo.Get(ctx, malak.FetchContactOptions{
			Email:       req.Email,
			WorkspaceID: deck.WorkspaceID,
		})

		if err != nil {
			contactID = uuid.Nil
		} else {
			contactID = contact.ID
		}

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
		ContactID:  contactID,
		SessionID:  malak.Reference(d.referenceGenerator.Generate(malak.EntityTypeSession)),
		DeviceInfo: req.DeviceInfo,
		OS:         req.OS,
		Browser:    req.Browser,
		IPAddress:  req.IPAddress,
		Country:    country,
		City:       city,
	}

	if err := d.deckRepo.CreateDeckSession(ctx, sessionReq); err != nil {
		logger.Error("could not create deck session", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create deck session"),
			StatusFailed
	}

	return fetchPublicDeckResponse{
		APIStatus: newAPIStatus(http.StatusOK, "fetched deck details"),
		Deck: malak.PublicDeck{
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
