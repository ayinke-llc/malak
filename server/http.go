package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/uploader"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/rs/cors"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func New(logger *zap.Logger,
	cfg config.Config,
	db *bun.DB,
	jwtTokenManager jwttoken.JWTokenManager,
	googleAuthProvider socialauth.SocialAuthProvider,
	mid *httplimit.Middleware, uploadImpl uploader.Uploader) (*http.Server, func()) {

	srv := &http.Server{
		Handler: buildRoutes(logger, db, cfg, jwtTokenManager,
			googleAuthProvider, mid, uploadImpl),
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
	}

	cleanupOtelResources := initOTELCapabilities(cfg, logger)

	return srv, func() {
		cleanupOtelResources()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("could not shut down server gracefully", zap.Error(err))
		}
	}
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
	logger *zap.Logger,
	db *bun.DB,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	googleAuthProvider socialauth.SocialAuthProvider,
	mid *httplimit.Middleware,
	uploadImpl uploader.Uploader) http.Handler {

	userRepo := postgres.NewUserRepository(db)
	workspaceRepo := postgres.NewWorkspaceRepository(db)
	planRepo := postgres.NewPlanRepository(db)
	contactRepo := postgres.NewContactRepository(db)
	updateRepo := postgres.NewUpdatesRepository(db)

	router := chi.NewRouter()

	referenceGenerator := malak.NewReferenceGenerator()

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
		cfg:                cfg,
		contactRepo:        contactRepo,
		referenceGenerator: referenceGenerator,
	}

	updateHandler := &updatesHandler{
		referenceGenerator: referenceGenerator,
		updateRepo:         updateRepo,
		uploaderImpl:       uploadImpl,
	}

	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(otelchi.Middleware("malak.server", otelchi.WithChiRoutes(router)))
	router.Use(jsonResponse)
	router.Use(mid.Handle)

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/connect/{provider}",
				WrapMalakHTTPHandler(logger, auth.Login, cfg, "Auth.Login"))
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Get("/",
				WrapMalakHTTPHandler(logger, auth.fetchCurrentUser, cfg, "Auth.fetchCurrentUser"))
		})

		r.Route("/workspaces", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Post("/",
				WrapMalakHTTPHandler(logger, workspaceHandler.createWorkspace, cfg, "workspaces.new"))

			r.Post("/{reference}",
				WrapMalakHTTPHandler(logger, workspaceHandler.switchCurrentWorkspaceForUser, cfg, "workspaces.switch"))

			r.Route("/updates", func(r chi.Router) {
				r.Post("/",
					WrapMalakHTTPHandler(logger, updateHandler.create, cfg, "updates.new"))
				r.Get("/",
					WrapMalakHTTPHandler(logger, updateHandler.list, cfg, "updates.list"))
				r.Post("/image",
					WrapMalakHTTPHandler(logger, updateHandler.uploadImage, cfg, "updates.image_upload"))
			})
		})

		r.Route("/contacts", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Post("/",
				WrapMalakHTTPHandler(logger, contactHandler.Create, cfg, "contacts.create"))
		})
	})

	return cors.AllowAll().
		Handler(router)
}
