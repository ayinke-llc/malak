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
	"syscall"
	"time"

	"github.com/adelowo/gulter"
	"github.com/adelowo/gulter/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/cache/rediscache"
	"github.com/ayinke-llc/malak/internal/pkg/email/smtp"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	watermillqueue "github.com/ayinke-llc/malak/internal/pkg/queue/watermill"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/util"
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

			emailClient, err := smtp.New(*cfg)
			if err != nil {
				logger.Fatal("could not set up smtp client",
					zap.Error(err))
			}

			queueHandler, err := watermillqueue.New(redisClient, *cfg, logger,
				emailClient, userRepo, workspaceRepo, updateRepo, contactRepo)
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
			})
			if err != nil {
				logger.Fatal("could not set up S3 client",
					zap.Error(err))
			}

			gulterHandler, err := gulter.New(
				gulter.WithMaxFileSize(cfg.Uploader.MaxUploadSize),
				gulter.WithValidationFunc(
					gulter.MimeTypeValidator("image/jpeg", "image/png", "application/pdf")),
				gulter.WithStorage(s3Store),
				gulter.WithIgnoreNonExistentKey(true),
				gulter.WithErrorResponseHandler(func(err error) http.HandlerFunc {
					return func(w http.ResponseWriter, _ *http.Request) {
						logger.Error("could not upload file", zap.Error(err))

						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_ = json.NewEncoder(w).Encode(server.APIStatus{
							Message: fmt.Sprintf("could not upload file...%s", err.Error()),
						})
					}
				}),
				gulter.WithNameFuncGenerator(func(s string) string {
					return uuid.New().String()
				}),
			)
			if err != nil {
				logger.Fatal("could not set up gulter uploader",
					zap.Error(err))
			}

			mid, err := httplimit.NewMiddleware(rateLimiterStore, server.HTTPThrottleKeyFunc)
			if err != nil {
				logger.Fatal("could not rate limiting middleware",
					zap.Error(err))
			}

			srv, cleanupSrv := server.New(logger,
				util.DeRef(cfg), db,
				tokenManager, googleAuthProvider,
				userRepo, workspaceRepo, planRepo, contactRepo, updateRepo,
				contactlistRepo, deckRepo, shareRepo,
				mid, gulterHandler, queueHandler, redisCache)

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
