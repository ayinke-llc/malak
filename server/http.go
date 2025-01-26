package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	_ "github.com/ayinke-llc/malak/swagger"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/rs/cors"
	"github.com/sethvargo/go-limiter/httplimit"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func New(logger *zap.Logger,
	cfg config.Config,
	db *bun.DB,
	jwtTokenManager jwttoken.JWTokenManager,
	googleAuthProvider socialauth.SocialAuthProvider,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	contactRepo malak.ContactRepository,
	updateRepo malak.UpdateRepository,
	contactListRepo malak.ContactListRepository,
	deckRepo malak.DeckRepository,
	shareRepo malak.ContactShareRepository,
	mid *httplimit.Middleware,
	gulterHandler *gulter.Gulter,
	queueHandler queue.QueueHandler,
	redisCache cache.Cache) (*http.Server, func()) {

	srv := &http.Server{
		Handler: buildRoutes(logger, db, cfg, jwtTokenManager,
			userRepo, workspaceRepo, planRepo,
			contactRepo, updateRepo, contactListRepo,
			deckRepo, shareRepo, googleAuthProvider, mid,
			gulterHandler, queueHandler, redisCache),
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
	}

	cleanupOtelResources := InitOTELCapabilities(cfg, logger)

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
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	contactRepo malak.ContactRepository,
	updateRepo malak.UpdateRepository,
	contactListRepo malak.ContactListRepository,
	deckRepo malak.DeckRepository,
	shareRepo malak.ContactShareRepository,
	googleAuthProvider socialauth.SocialAuthProvider,
	ratelimiterMiddleware *httplimit.Middleware,
	gulterHandler *gulter.Gulter,
	queueHandler queue.QueueHandler,
	redisCache cache.Cache,
) http.Handler {

	if cfg.HTTP.Swagger.UIEnabled {
		go func() {
			r := chi.NewRouter()

			r.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.HTTP.Swagger.Port)),
			))

			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTP.Swagger.Port), r); err != nil {
				logger.Error("error with swagger server", zap.Error(err))
			}
		}()
	}

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
		contactListRepo:    contactListRepo,
		contactShareRepo:   shareRepo,
	}

	updateHandler := &updatesHandler{
		referenceGenerator: referenceGenerator,
		updateRepo:         updateRepo,
		cfg:                cfg,
		queueHandler:       queueHandler,
		cache:              redisCache,
		uuidGenerator:      malak.NewGoogleUUID(),
	}

	webhookHandler := &webhookHandler{
		workspaceRepo:      workspaceRepo,
		cfg:                cfg,
		userRepo:           userRepo,
		planRepo:           planRepo,
		referenceGenerator: referenceGenerator,
		updateRepo:         updateRepo,
		contactRepo:        contactRepo,
	}

	deckHandler := &deckHandler{
		referenceGenerator: referenceGenerator,
		deckRepo:           deckRepo,
		cfg:                cfg,
		cache:              redisCache,
	}

	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)
	router.Use(
		middleware.AllowContentType("application/json", "multipart/form-data"))
	router.Use(
		otelchi.Middleware("malak.server",
			otelchi.WithChiRoutes(router)))
	router.Use(jsonResponse)
	router.Use(ratelimiterMiddleware.Handle)

	router.Route("/hooks", func(r chi.Router) {
		r.Post("/resend", webhookHandler.handleResend(logger))
	})

	router.Route("/updates", func(r chi.Router) {
		r.Post("/react", updateHandler.handleReaction(logger))
	})

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

				r.Get("/pins",
					WrapMalakHTTPHandler(logger, updateHandler.listPinnedUpdates, cfg, "updates.list.pins"))

				r.Get("/{reference}",
					WrapMalakHTTPHandler(logger, updateHandler.fetchUpdate, cfg, "updates.fetchUpdate"))
				r.Delete("/{reference}",
					WrapMalakHTTPHandler(logger, updateHandler.delete, cfg, "updates.delete"))

				r.Post("/{reference}",
					WrapMalakHTTPHandler(logger, updateHandler.sendUpdate, cfg, "updates.send"))

				r.Put("/{reference}",
					WrapMalakHTTPHandler(logger, updateHandler.update, cfg, "updates.content_update"))

				r.Post("/{reference}/pin",
					WrapMalakHTTPHandler(logger, updateHandler.togglePinned, cfg, "updates.togglePinned"))

				r.Post("/{reference}/duplicate",
					WrapMalakHTTPHandler(logger, updateHandler.duplicate, cfg, "updates.duplicate"))

				r.Post("/{reference}/preview",
					WrapMalakHTTPHandler(logger, updateHandler.previewUpdate, cfg, "updates.preview"))

				r.Get("/{reference}/analytics",
					WrapMalakHTTPHandler(logger, updateHandler.fetchUpdateAnalytics, cfg, "updates.analytics"))
			})
		})

		r.Route("/contacts", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Post("/",
				WrapMalakHTTPHandler(logger, contactHandler.Create, cfg, "contacts.create"))

			r.Get("/",
				WrapMalakHTTPHandler(logger, contactHandler.list, cfg, "contacts.list"))

			r.Get("/{reference}",
				WrapMalakHTTPHandler(logger, contactHandler.fetchContact, cfg, "contacts.fetch"))

			r.Route("/lists", func(r chi.Router) {
				r.Post("/",
					WrapMalakHTTPHandler(logger, contactHandler.createContactList, cfg, "contacts.lists.new"))

				r.Get("/",
					WrapMalakHTTPHandler(logger, contactHandler.fetchContactLists, cfg, "contacts.lists.fetch"))

				r.Delete("/{reference}",
					WrapMalakHTTPHandler(logger, contactHandler.deleteContactList, cfg, "contacts.lists.delete"))

				r.Put("/{reference}",
					WrapMalakHTTPHandler(logger, contactHandler.editContactList, cfg, "contacts.lists.update"))

				r.Post("/{reference}",
					WrapMalakHTTPHandler(logger, contactHandler.addUserToContactList, cfg, "contacts.lists.add"))
			})
		})

		r.Route("/decks", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))

			r.Post("/",
				WrapMalakHTTPHandler(logger, deckHandler.Create, cfg, "decks.add"))

			r.Get("/",
				WrapMalakHTTPHandler(logger, deckHandler.List, cfg, "decks.list"))

			r.Delete("/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.Delete, cfg, "decks.delete"))

			r.Get("/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.fetch, cfg, "decks.retrieve"))

			r.Put("/{reference}/preferences",
				WrapMalakHTTPHandler(logger, deckHandler.updatePreferences, cfg, "decks.preferences.update"))

			r.Post("/{reference}/archive",
				WrapMalakHTTPHandler(logger, deckHandler.toggleArchive, cfg, "decks.archive"))

			r.Post("/{reference}/pin",
				WrapMalakHTTPHandler(logger, deckHandler.togglePinned, cfg, "decks.togglePinned"))

		})

		r.Route("/uploads", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))

			r.Route("/decks", func(r chi.Router) {
				r.Use(gulterHandler.Upload(cfg.Uploader.S3.DeckBucket, "image_body"))

				r.Post("/",
					WrapMalakHTTPHandler(logger, deckHandler.uploadImage, cfg, "decks.upload"))
			})

			r.Route("/images", func(r chi.Router) {
				r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
				r.Use(gulterHandler.Upload(cfg.Uploader.S3.Bucket, "image_body"))

				r.Post("/",
					WrapMalakHTTPHandler(logger, updateHandler.uploadImage, cfg, "updates.image_upload"))
			})
		})
	})

	return cors.AllowAll().
		Handler(router)
}
