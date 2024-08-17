package server

import (
	"errors"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
)

type authHandler struct {
	logger    *logrus.Entry
	googleCfg socialauth.SocialAuthProvider
	cfg       config.Config
	userRepo  malak.UserRepository
}

type authenticateUserRequest struct {
	genericRequest

	Code string `json:"code,omitempty"`
}

// @Summary Sign in with a social login provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param message body authenticateUserRequest true "auth exchange data"
// @Success 200 {object} createdUserResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /auth/login [post]
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx, span, rid := getTracer(r.Context(), r, "Login", a.cfg.Otel.IsEnabled)
	defer span.End()

	provider := chi.URLParam(r, "provider")

	logger := a.logger.WithField("method", "login").
		WithField("request_id", rid).
		WithField("provider", provider)

	span.SetAttributes(attribute.String("auth_provider", provider))

	logger.Debug("Authenticating user")

	if provider != "google" {
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "unspported provider"))
		return
	}

	req := new(authenticateUserRequest)

	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "invalid request body"))
		return
	}

	token, err := a.googleCfg.Validate(ctx, socialauth.ValidateOptions{
		Code: req.Code,
	})
	if err != nil {
		logger.WithError(err).Error("could not exchange token")
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "could not verify your sign in with Google"))
		return
	}

	u, err := a.googleCfg.User(ctx, token)
	if err != nil {
		logger.WithError(err).Error("could not fetch user details")
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, "could not fetch user details from oauth2 provider"))
		return
	}

	user := &malak.User{
		Email:    malak.Email(u.Email),
		FullName: u.Name,
		Metadata: &malak.UserMetadata{},
	}

	err = a.userRepo.Create(ctx, user)
	if errors.Is(err, malak.ErrUserExists) {
		_ = render.Render(w, r, newAPIStatus(http.StatusBadRequest, err.Error()))
		return
	}

	if err != nil {
		logger.WithError(err).Error("an error occurred while creating user")
		_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "an error occurred while creating user"))
		return
	}

	_ = render.Render(w, r, createdUserResponse{
		User:      user,
		APIStatus: newAPIStatus(http.StatusOK, "user Successfully created"),
	})
}
