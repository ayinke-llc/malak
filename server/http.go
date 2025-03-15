package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/integrations"
	"github.com/ayinke-llc/malak/internal/pkg/billing"
	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/ayinke-llc/malak/internal/pkg/geolocation"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/secret"
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
	dashboardRepo malak.DashboardRepository,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	contactRepo malak.ContactRepository,
	updateRepo malak.UpdateRepository,
	contactListRepo malak.ContactListRepository,
	deckRepo malak.DeckRepository,
	shareRepo malak.ContactShareRepository,
	preferenceRepo malak.PreferenceRepository,
	integrationRepo malak.IntegrationRepository,
	templatesRepo malak.TemplateRepository,
	dashboardLinkRepo malak.DashboardLinkRepository,
	mid *httplimit.Middleware,
	gulterHandler *gulter.Gulter,
	queueHandler queue.QueueHandler,
	redisCache cache.Cache,
	billingClient billing.Client,
	integrationManager *integrations.IntegrationsManager,
	secretsClient secret.SecretClient,
	geolocationService geolocation.GeolocationService) (*http.Server, func()) {

	if err := cfg.Validate(); err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return nil, nil
	}

	srv := &http.Server{
		Handler: buildRoutes(logger, db, cfg, jwtTokenManager,
			dashboardRepo,
			userRepo, workspaceRepo, planRepo,
			contactRepo, updateRepo, contactListRepo,
			deckRepo, shareRepo, preferenceRepo, integrationRepo, templatesRepo,
			dashboardLinkRepo,
			googleAuthProvider, mid, gulterHandler,
			queueHandler, redisCache, billingClient, integrationManager, secretsClient,
			geolocationService),
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
	_ *bun.DB,
	cfg config.Config,
	jwtTokenManager jwttoken.JWTokenManager,
	dashboardRepo malak.DashboardRepository,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
	planRepo malak.PlanRepository,
	contactRepo malak.ContactRepository,
	updateRepo malak.UpdateRepository,
	contactListRepo malak.ContactListRepository,
	deckRepo malak.DeckRepository,
	shareRepo malak.ContactShareRepository,
	preferenceRepo malak.PreferenceRepository,
	integrationRepo malak.IntegrationRepository,
	templatesRepo malak.TemplateRepository,
	dashboardLinkRepo malak.DashboardLinkRepository,
	googleAuthProvider socialauth.SocialAuthProvider,
	ratelimiterMiddleware *httplimit.Middleware,
	gulterHandler *gulter.Gulter,
	queueHandler queue.QueueHandler,
	redisCache cache.Cache,
	billingClient billing.Client,
	integrationManager *integrations.IntegrationsManager,
	secretsClient secret.SecretClient,
	geolocationService geolocation.GeolocationService) http.Handler {

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
		preferenceRepo:          preferenceRepo,
		integrationRepo:         integrationRepo,
		referenceGenerationFunc: malak.GenerateReference,
		queueClient:             queueHandler,
		billingClient:           billingClient,
		integrationManager:      integrationManager,
		secretsClient:           secretsClient,
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
		templateRepo:       templatesRepo,
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

	stripeHan := &stripeHandler{
		user:            userRepo,
		planRepo:        planRepo,
		logger:          logger,
		billingClient:   billingClient,
		workRepo:        workspaceRepo,
		preferencesRepo: preferenceRepo,
		taskQueue:       queueHandler,
		cfg:             cfg,
	}

	deckHandler := &deckHandler{
		gulterStore:        gulterHandler.Storage(),
		referenceGenerator: referenceGenerator,
		cache:              redisCache,
		deckRepo:           deckRepo,
		cfg:                cfg,
		geolocationService: geolocationService,
		contactRepo:        contactRepo,
	}

	dashHandler := &dashboardHandler{
		cfg:               cfg,
		dashboardRepo:     dashboardRepo,
		generator:         referenceGenerator,
		integrationRepo:   integrationRepo,
		dashboardLinkRepo: dashboardLinkRepo,
		contactRepo:       contactRepo,
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
		r.Post("/stripe", stripeHan.handleWebhook)
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
			r.Use(requireWorkspaceValidSubscription(cfg))
			r.Get("/",
				WrapMalakHTTPHandler(logger, auth.fetchCurrentUser, cfg, "Auth.fetchCurrentUser"))
		})

		r.Route("/workspaces", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Use(requireWorkspaceValidSubscription(cfg))

			r.Post("/",
				WrapMalakHTTPHandler(logger, workspaceHandler.createWorkspace, cfg, "workspaces.new"))

			r.Patch("/",
				WrapMalakHTTPHandler(logger, workspaceHandler.updateWorkspace, cfg, "workspaces.update"))

			r.Get("/preferences",
				WrapMalakHTTPHandler(logger, workspaceHandler.getPreferences, cfg, "workspaces.preferences.get"))

			r.Put("/preferences",
				WrapMalakHTTPHandler(logger, workspaceHandler.updatePreferences, cfg, "workspaces.preferences.update"))

			r.Post("/billing",
				WrapMalakHTTPHandler(logger, workspaceHandler.getBillingPortal, cfg, "workspaces.billing.portal"))

			r.Post("/switch/{reference}",
				WrapMalakHTTPHandler(logger, workspaceHandler.switchCurrentWorkspaceForUser, cfg, "workspaces.switch"))

			r.Route("/integrations", func(r chi.Router) {
				r.Get("/",
					WrapMalakHTTPHandler(logger,
						workspaceHandler.getIntegrations, cfg, "workspaces.integrations.list"))

				r.Post("/{reference}",
					WrapMalakHTTPHandler(logger,
						workspaceHandler.enableIntegration, cfg, "workspaces.integrations.store"))

				r.Put("/{reference}",
					WrapMalakHTTPHandler(logger,
						workspaceHandler.updateAPIKeyForIntegration, cfg, "workspaces.integrations.update_token"))

				r.Delete("/{reference}",
					WrapMalakHTTPHandler(logger,
						workspaceHandler.disableIntegration, cfg, "workspaces.integrations.disable"))

				r.Post("/{reference}/ping",
					WrapMalakHTTPHandler(logger,
						workspaceHandler.pingIntegration, cfg, "workspaces.integrations.ping"))
			})

			r.Route("/updates", func(r chi.Router) {

				r.Post("/",
					WrapMalakHTTPHandler(logger, updateHandler.create, cfg, "updates.new"))
				r.Get("/",
					WrapMalakHTTPHandler(logger, updateHandler.list, cfg, "updates.list"))

				r.Get("/pins",
					WrapMalakHTTPHandler(logger, updateHandler.listPinnedUpdates, cfg, "updates.list.pins"))

				r.Get("/templates",
					WrapMalakHTTPHandler(logger, updateHandler.templates, cfg, "updates.list.templates"))

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
			r.Use(requireWorkspaceValidSubscription(cfg))
			r.Post("/",
				WrapMalakHTTPHandler(logger, contactHandler.Create, cfg, "contacts.create"))

			r.Get("/",
				WrapMalakHTTPHandler(logger, contactHandler.list, cfg, "contacts.list"))

			r.Get("/{reference}",
				WrapMalakHTTPHandler(logger, contactHandler.fetchContact, cfg, "contacts.fetch"))

			r.Delete("/{reference}",
				WrapMalakHTTPHandler(logger, contactHandler.deleteContact, cfg, "contacts.delete"))

			r.Put("/{reference}",
				WrapMalakHTTPHandler(logger, contactHandler.editContact, cfg, "contacts.edit"))

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
			r.Use(requireWorkspaceValidSubscription(cfg))

			r.Post("/",
				WrapMalakHTTPHandler(logger, deckHandler.Create, cfg, "decks.add"))

			r.Get("/",
				WrapMalakHTTPHandler(logger, deckHandler.List, cfg, "decks.list"))

			r.Delete("/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.Delete, cfg, "decks.delete"))

			r.Get("/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.fetch, cfg, "decks.retrieve"))

			r.Get("/{reference}/sessions",
				WrapMalakHTTPHandler(logger, deckHandler.fetchDeckSessions, cfg, "decks.sessions"))

			r.Get("/{reference}/analytics",
				WrapMalakHTTPHandler(logger, deckHandler.fetchEngagements, cfg, "decks.engagements"))

			r.Put("/{reference}/preferences",
				WrapMalakHTTPHandler(logger, deckHandler.updatePreferences, cfg, "decks.preferences.update"))

			r.Post("/{reference}/archive",
				WrapMalakHTTPHandler(logger, deckHandler.toggleArchive, cfg, "decks.archive"))

			r.Post("/{reference}/pin",
				WrapMalakHTTPHandler(logger, deckHandler.togglePinned, cfg, "decks.togglePinned"))

		})

		r.Route("/dashboards", func(r chi.Router) {
			r.Use(requireAuthentication(logger, jwtTokenManager, cfg, userRepo, workspaceRepo))
			r.Use(requireWorkspaceValidSubscription(cfg))

			r.Post("/",
				WrapMalakHTTPHandler(logger, dashHandler.create, cfg, "dashboards.create"))

			r.Get("/",
				WrapMalakHTTPHandler(logger, dashHandler.list, cfg, "dashboards.list"))

			r.Get("/charts",
				WrapMalakHTTPHandler(logger, dashHandler.listAllCharts, cfg, "dashboards.list.charts"))

			r.Get("/charts/{reference}",
				WrapMalakHTTPHandler(logger, dashHandler.fetchChartingData, cfg, "dashboards.charts.datapoints"))

			r.Get("/{reference}",
				WrapMalakHTTPHandler(logger, dashHandler.fetchDashboard, cfg, "dashboards.fetch"))

			r.Post("/{reference}/positions",
				WrapMalakHTTPHandler(logger, dashHandler.updateDashboardPositions, cfg, "dashboards.positions.update"))

			r.Put("/{reference}/charts",
				WrapMalakHTTPHandler(logger, dashHandler.addChart, cfg, "dashboards.charts.add"))

			r.Delete("/{reference}/charts",
				WrapMalakHTTPHandler(logger, dashHandler.removeChart, cfg, "dashboards.charts.remove"))

			r.Post("/{reference}/access-control/link",
				WrapMalakHTTPHandler(logger, dashHandler.generateLink, cfg, "dashboards.access-control.link.generate"))
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

		r.Route("/public", func(r chi.Router) {
			r.Post("/decks/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.publicDeckDetails, cfg, "public.decks.fetch"))
			r.Put("/decks/{reference}",
				WrapMalakHTTPHandler(logger, deckHandler.updateDeckViewerSession, cfg, "public.decks.update"))

			r.Get("/dashboards/{reference}",
				WrapMalakHTTPHandler(logger, dashHandler.publicDashboardDetails, cfg, "public.dashboards.fetch"))

			r.Get("/dashboards/{reference}/charts/{chart_reference}",
				WrapMalakHTTPHandler(logger, dashHandler.publicChartingDataFetch, cfg, "public.charts.datapoints"))
		})
	})

	return cors.AllowAll().
		Handler(router)
}
