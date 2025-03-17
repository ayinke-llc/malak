package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/adelowo/gulter"
	"github.com/adelowo/gulter/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/integrations"
	"github.com/ayinke-llc/malak/internal/integrations/brex"
	"github.com/ayinke-llc/malak/internal/integrations/mercury"
	"github.com/ayinke-llc/malak/internal/pkg/billing/stripe"
	"github.com/ayinke-llc/malak/internal/pkg/cache/rediscache"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/email/resend"
	"github.com/ayinke-llc/malak/internal/pkg/email/smtp"
	"github.com/ayinke-llc/malak/internal/pkg/geolocation/maxmind"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	watermillqueue "github.com/ayinke-llc/malak/internal/pkg/queue/watermill"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/ayinke-llc/malak/internal/secret/aes"
	"github.com/ayinke-llc/malak/internal/secret/infisical"
	"github.com/ayinke-llc/malak/internal/secret/secretsmanager"
	"github.com/ayinke-llc/malak/internal/secret/vault"
	"github.com/ayinke-llc/malak/server"
	"github.com/google/uuid"
	redisotel "github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
	"github.com/spf13/cobra"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

const (
	maxExportBatchSize = sdktrace.DefaultMaxExportBatchSize
)

type APIStatus struct {
	Message string `json:"message"`
}

func parseHTTPPortFromEnv() int {
	s := os.Getenv(("ENV_HTTP_PORT"))
	if s == "" {
		return 5300
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return 5300
	}
	return n
}

func addHTTPCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use: "http",
		Run: func(cmd *cobra.Command, args []string) {

			sig := make(chan os.Signal, 1)

			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			h, _ := os.Hostname()

			var logger *zap.Logger
			var err error

			switch cfg.Logging.Mode {
			case config.LogModeProd:

				logger, err = zap.NewProduction()
				if err != nil {
					fmt.Printf(`{"error":%s}`, err)
					os.Exit(1)
				}

			case config.LogModeDev:

				logger, err = zap.NewDevelopment()
				if err != nil {
					fmt.Printf(`{"error":%s}`, err)
					os.Exit(1)
				}
			}

			logger = logger.With(zap.String("host", h),
				zap.String("app", "malak"))

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Fatal("could not set up database connection",
					zap.Error(err))
			}

			userRepo := postgres.NewUserRepository(db)
			workspaceRepo := postgres.NewWorkspaceRepository(db)
			planRepo := postgres.NewPlanRepository(db)
			contactRepo := postgres.NewContactRepository(db)
			updateRepo := postgres.NewUpdatesRepository(db)
			contactlistRepo := postgres.NewContactListRepository(db)
			deckRepo := postgres.NewDeckRepository(db)
			shareRepo := postgres.NewShareRepository(db)
			preferenceRepo := postgres.NewPreferenceRepository(db)
			integrationRepo := postgres.NewIntegrationRepo(db)
			dashRepo := postgres.NewDashboardRepo(db)
			templatesRepo := postgres.NewTemplateRepository(db)
			dashboardLinkRepo := postgres.NewDashboardLinkRepo(db)

			googleAuthProvider := socialauth.NewGoogle(*cfg)

			tokenManager := jwttoken.New(*cfg)

			opts, err := redis.ParseURL(cfg.Database.Redis.DSN)
			if err != nil {
				logger.Fatal("could not parse redis dsn",
					zap.Error(err))
			}

			redisClient := redis.NewClient(opts)

			if cfg.Otel.IsEnabled {
				if err := redisotel.InstrumentTracing(redisClient); err != nil {
					logger.Fatal("could not instrument tracing of redis client",
						zap.Error(err))
				}

				if err := redisotel.InstrumentMetrics(redisClient); err != nil {
					logger.Fatal("could not instrument metrics collection of redis client",
						zap.Error(err))
				}
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*30)
			defer cancelFn()

			if err := redisClient.Ping(ctx).Err(); err != nil {
				logger.Fatal("could not ping redis",
					zap.Error(err))
			}

			var emailClient email.Client

			switch cfg.Email.Provider {
			case config.EmailProviderSmtp:
				var err error

				emailClient, err = smtp.New(*cfg)
				if err != nil {
					logger.Fatal("could not set up smtp client",
						zap.Error(err))
				}

			case config.EmailProviderResend:
				var err error

				emailClient, err = resend.New(*cfg)
				if err != nil {
					logger.Fatal("could not set up smtp client",
						zap.Error(err))
				}

			default:
				logger.Fatal("unsupported email provider", zap.String("provider", cfg.Email.Provider.String()))
			}

			billingClient, err := stripe.New(hermes.DeRef(cfg))
			if err != nil {
				logger.Fatal("could not set up stripe client",
					zap.Error(err))
			}

			queueHandler, err := watermillqueue.New(
				redisClient, hermes.DeRef(cfg),
				logger, emailClient, userRepo, workspaceRepo,
				updateRepo, contactRepo, billingClient)
			if err != nil {
				logger.Fatal("could not set up watermill queue", zap.Error(err))
			}

			go func() {
				queueHandler.Start(context.Background())
			}()

			redisCache, err := rediscache.New(redisClient)
			if err != nil {
				logger.Fatal("could not set up redis cache", zap.Error(err))
			}

			rateLimiterStore, err := getRatelimiter(*cfg)
			if err != nil {
				logger.Fatal("could not create rate limiter",
					zap.Error(err))
			}

			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: !cfg.Uploader.S3.UseTLS,
					},
				},
			}

			geoService, err := maxmind.New(hermes.DeRef(cfg))
			if err != nil {
				logger.Fatal("could not set up maxmind db", zap.Error(err))
			}

			s3Config, err := awsConfig.LoadDefaultConfig(
				context.Background(),
				awsConfig.WithRegion(cfg.Uploader.S3.Region),
				awsConfig.WithHTTPClient(httpClient),
				awsConfig.WithCredentialsProvider(
					awsCreds.NewStaticCredentialsProvider(
						cfg.Uploader.S3.AccessKey,
						cfg.Uploader.S3.AccessSecret,
						"")),
				//nolint:staticcheck
				awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					//nolint:staticcheck
					return aws.Endpoint{
						URL:               cfg.Uploader.S3.Endpoint,
						SigningRegion:     cfg.Uploader.S3.Region,
						HostnameImmutable: true,
					}, nil
				})),
			)
			if err != nil {
				logger.Fatal("could not set up S3 config",
					zap.Error(err))
			}

			s3Store, err := storage.NewS3FromConfig(s3Config, storage.S3Options{
				DebugMode:    cfg.Uploader.S3.LogOperations,
				UsePathStyle: true,
				Bucket:       cfg.Uploader.S3.Bucket,
			})
			if err != nil {
				logger.Fatal("could not set up S3 client",
					zap.Error(err))
			}

			imageUploadGulterHandler, err := gulter.New(
				gulter.WithMaxFileSize(cfg.Uploader.MaxUploadSize),
				gulter.WithValidationFunc(
					gulter.MimeTypeValidator("image/jpeg", "image/png")),
				gulter.WithStorage(s3Store),
				gulter.WithIgnoreNonExistentKey(true),
				gulter.WithErrorResponseHandler(func(err error) http.HandlerFunc {
					return func(w http.ResponseWriter, _ *http.Request) {
						logger.Error("could not upload file", zap.Error(err))

						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_ = json.NewEncoder(w).Encode(APIStatus{
							Message: fmt.Sprintf("could not upload file...%s", err.Error()),
						})
					}
				}),
				gulter.WithNameFuncGenerator(func(s string) string {
					return uuid.New().String() + strings.Replace(s, " ", "", -1)
				}),
			)
			if err != nil {
				logger.Fatal("could not set up gulter uploader",
					zap.Error(err))
			}

			decks3Store, err := storage.NewS3FromConfig(s3Config, storage.S3Options{
				DebugMode:    cfg.Uploader.S3.LogOperations,
				UsePathStyle: true,
				Bucket:       cfg.Uploader.S3.DeckBucket,
			})
			if err != nil {
				logger.Fatal("could not set up S3 client",
					zap.Error(err))
			}

			deckUploadGulterHandler, err := gulter.New(
				gulter.WithMaxFileSize(cfg.Uploader.MaxUploadSize),
				gulter.WithValidationFunc(
					gulter.MimeTypeValidator("application/pdf")),
				gulter.WithStorage(decks3Store),
				gulter.WithIgnoreNonExistentKey(true),
				gulter.WithErrorResponseHandler(func(err error) http.HandlerFunc {
					return func(w http.ResponseWriter, _ *http.Request) {
						logger.Error("could not upload file", zap.Error(err))

						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_ = json.NewEncoder(w).Encode(APIStatus{
							Message: fmt.Sprintf("could not upload file...%s", err.Error()),
						})
					}
				}),
				gulter.WithNameFuncGenerator(func(s string) string {
					return uuid.New().String() + strings.Replace(s, " ", "", -1)
				}),
			)
			if err != nil {
				logger.Fatal("could not set up gulter deck uploader",
					zap.Error(err))
			}

			mid, err := httplimit.NewMiddleware(rateLimiterStore, server.HTTPThrottleKeyFunc)
			if err != nil {
				logger.Fatal("could not rate limiting middleware",
					zap.Error(err))
			}

			integrationManager, err := buildIntegrationManager(integrationRepo, *cfg, logger)
			if err != nil {
				logger.Fatal("could not build integration manager", zap.Error(err))
			}

			secretsProvider, err := buildSecretsProvider(*cfg)
			if err != nil {
				logger.Fatal("could not build secrets provider", zap.Error(err))
			}

			srv, cleanupSrv := server.New(logger,
				util.DeRef(cfg), db,
				tokenManager, googleAuthProvider,
				dashRepo,
				userRepo, workspaceRepo, planRepo, contactRepo,
				updateRepo, contactlistRepo, deckRepo, shareRepo,
				preferenceRepo, integrationRepo,
				templatesRepo, dashboardLinkRepo,
				mid,
				queueHandler, redisCache, billingClient,
				integrationManager, secretsProvider, geoService,
				imageUploadGulterHandler, deckUploadGulterHandler)

			go func() {
				if err := srv.ListenAndServe(); err != nil {
					logger.Error("error with http server",
						zap.Error(err))
				}
			}()

			<-sig

			cleanupSrv()

			logger.Debug("shutting down Malak's server")
			if err := db.Close(); err != nil {
				logger.Error("could not close db",
					zap.Error(err))
			}

			if err := queueHandler.Close(); err != nil {
				logger.Error("could not close the queue handler", zap.Error(err))
			}

			_ = logger.Sync()
		},
	}

	c.AddCommand(cmd)
}

func buildIntegrationManager(integrationRepo malak.IntegrationRepository, cfg config.Config, logger *zap.Logger) (
	*integrations.IntegrationsManager, error) {
	i := integrations.NewManager()

	integrations, err := integrationRepo.System(context.Background())
	if err != nil {
		return nil, err
	}

	for _, v := range integrations {
		provider, err := malak.ParseIntegrationProvider(strings.ToLower(v.IntegrationName))
		if err != nil {
			logger.Warn("invalid integration provider",
				zap.String("integration_name", v.IntegrationName),
				zap.Error(err))
			continue
		}

		switch provider {
		case malak.IntegrationProviderMercury:
			client, err := mercury.New(cfg)
			if err != nil {
				return nil, err
			}

			i.Add(provider, client)

		case malak.IntegrationProviderBrex:
			client, err := brex.New(cfg)
			if err != nil {
				return nil, err
			}

			i.Add(provider, client)

		default:
			logger.Warn("provider not yet implemented",
				zap.String("provider", provider.String()))
		}
	}

	return i, nil
}

func getRatelimiter(cfg config.Config) (limiter.Store, error) {

	if !cfg.HTTP.RateLimit.IsEnabled {
		return noopstore.New()
	}

	switch cfg.HTTP.RateLimit.Type {
	case config.RateLimiterTypeMemory:
		return memorystore.New(&memorystore.Config{
			Interval: cfg.HTTP.RateLimit.BurstInterval,
			Tokens:   cfg.HTTP.RateLimit.RequestsPerMinute,
		})

	default:
		return nil, errors.New("unsupported ratelimter")
	}
}

func buildSecretsProvider(cfg config.Config) (secret.SecretClient, error) {
	switch cfg.Secrets.Provider {
	case secret.SecretProviderVault:
		return vault.New(cfg)
	case secret.SecretProviderInfisical:
		return infisical.New(cfg)
	case secret.SecretProviderAesGcm:
		return aes.New(cfg)
	case secret.SecretProviderSecretsmanager:
		return secretsmanager.New(cfg)
	default:
		return nil, fmt.Errorf("unsupported secrets provider: %s", cfg.Secrets.Provider)
	}
}
