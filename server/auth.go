package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
)

type authHandler struct {
	logger *logrus.Entry
	// googleCfg
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
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /auth/login [post]
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx, span, rid := getTracer(r.Context(), r, "Login")
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

	_ = ctx
}
