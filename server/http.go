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
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Entry,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	googleAuthProvider socialauth.SocialAuthProvider,
	mid *httplimit.Middleware) (*http.Server, func()) {

	srv := &http.Server{
		Handler: buildRoutes(logger, cfg, jwtTokenManager,
			userRepo, workspaceRepo, planRepo,
			googleAuthProvider, mid),
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
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

func buildRoutes(
	logger *logrus.Entry,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	googleAuthProvider socialauth.SocialAuthProvider,
	mid *httplimit.Middleware) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(otelchi.Middleware("malak.server", otelchi.WithChiRoutes(router)))
	router.Use(jsonResponse)
	router.Use(mid.Handle)

	auth := &authHandler{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		googleCfg:     googleAuthProvider,
		tokenManager:  jwtTokenManager,
	}

	workspaceHandler := &workspaceHandler{
		workspaceRepo:           workspaceRepo,
		cfg:                     cfg,
		userRepo:                userRepo,
		planRepo:                planRepo,
		referenceGenerationFunc: malak.GenerateReference,
	}

	contactHandler := &contactHandler{
		cfg: cfg,
	}

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/connect/{provider}", WrapMalakHTTPHandler(auth.Login, cfg, "Auth.Login"))
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Get("/", WrapMalakHTTPHandler(auth.fetchCurrentUser, cfg, "Auth.fetchCurrentUser"))
		})

		r.Route("/workspaces", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Post("/", WrapMalakHTTPHandler(workspaceHandler.createWorkspace, cfg, "workspaces.new"))
		})

		r.Route("/contacts", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Post("/", WrapMalakHTTPHandler(contactHandler.Create, cfg, "contacts.create"))
		})
	})

	return cors.AllowAll().
		Handler(router)
}
