package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/server"
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

			userRepo := postgres.NewUserRepository(db)
			workspaceRepo := postgres.NewWorkspaceRepository(db)

			googleAuthProvider := socialauth.NewGoogle(*cfg)

			srv, cleanupSrv := server.New(logger, *cfg, userRepo, workspaceRepo, googleAuthProvider)

			go func() {
				if err := srv.ListenAndServe(); err != nil {
					logger.WithError(err).Error("error with http server")
				}
			}()

			// opts, err := redis.ParseURL(cfg.Database.Redis.DSN)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			//
			// redisClient := redis.NewClient(opts)
			//
			// if err := redisotel.InstrumentTracing(redisClient); err != nil {
			// 	log.Fatal(err)
			// }
			//
			// if err := redisotel.InstrumentMetrics(redisClient); err != nil {
			// 	log.Fatal(err)
			// }
			//
			// if err := redisClient.Ping(context.Background()).Err(); err != nil {
			// 	log.Fatal(err)
			// }

			<-sig

			logger.Debug("shutting down Malak's server")
			cleanupSrv()
		},
	}

	c.AddCommand(cmd)
}
