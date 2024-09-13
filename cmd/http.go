package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/ayinke-llc/malak/server"
	redisotel "github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

			logger := logrus.WithField("host", h).
				WithField("app", "malak")

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.WithError(err).Fatal("could not set up database connection")
			}

			googleAuthProvider := socialauth.NewGoogle(*cfg)

			tokenManager := jwttoken.New(*cfg)

			opts, err := redis.ParseURL(cfg.Database.Redis.DSN)
			if err != nil {
				logger.WithError(err).Fatal("could not parse redis dsn")
			}

			redisClient := redis.NewClient(opts)

			if err := redisotel.InstrumentTracing(redisClient); err != nil {
				logger.WithError(err).Fatal("could not instrument tracing of redis client")
			}

			if err := redisotel.InstrumentMetrics(redisClient); err != nil {
				logger.WithError(err).Fatal("could not instrument metrics of redis client")
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*30)
			defer cancelFn()

			if err := redisClient.Ping(ctx).Err(); err != nil {
				logger.WithError(err).Fatal("could not ping Redis")
			}

			_ = redisClient

			rateLimiterStore, err := getRatelimiter(*cfg)
			if err != nil {
				logger.WithError(err).Fatal("could not create rate limiter")
			}

			mid, err := httplimit.NewMiddleware(rateLimiterStore, server.HTTPThrottleKeyFunc)
			if err != nil {
				logger.WithError(err).Fatal("could not set up rate limiting middleware")
			}

			srv, cleanupSrv := server.New(logger, util.DeRef(cfg), db,
				tokenManager, googleAuthProvider, mid)

			go func() {
				if err := srv.ListenAndServe(); err != nil {
					logger.WithError(err).Error("error with http server")
				}
			}()

			<-sig

			logger.Debug("shutting down Malak's server")
			if err := db.Close(); err != nil {
				logger.WithError(err).Error("could not close db")
			}

			cleanupSrv()
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
