package server

import (
	"fmt"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Entry,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	googleAuthProvider socialauth.SocialAuthProvider) (*http.Server, func()) {

	srv := &http.Server{
		Handler: buildRoutes(logger, cfg, jwtTokenManager, userRepo, workspaceRepo, googleAuthProvider),
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
	}

	return srv, initOTELCapabilities(cfg, logger)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func buildRoutes(logger *logrus.Entry,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	googleAuthProvider socialauth.SocialAuthProvider) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(otelchi.Middleware("malak.server", otelchi.WithChiRoutes(router)))
	router.Use(jsonResponse)

	auth := &authHandler{
		logger:        logger,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		googleCfg:     googleAuthProvider,
		tokenManager:  jwtTokenManager,
	}

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/connect/{provider}", auth.Login)
		})
	})

	return cors.AllowAll().
		Handler(router)
}
