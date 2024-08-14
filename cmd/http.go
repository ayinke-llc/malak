package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ayinke-llc/malak/config"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
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

			logrus.SetOutput(os.Stdout)

			var formatter logrus.Formatter = &logrus.JSONFormatter{}

			if cfg.Logging.Format == config.LogFormatText {
				formatter = &logrus.TextFormatter{}
			}

			logrus.SetFormatter(formatter)

			lvl, err := logrus.ParseLevel(cfg.Logging.Level)
			if err != nil {
				lvl = logrus.DebugLevel
			}

			logrus.SetLevel(lvl)

			h, _ := os.Hostname()

			logger := logrus.WithField("host", h).
				WithField("app", "malak")

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

			logger.Debug("shutting down server")
			// cleanup()
		},
	}

	c.AddCommand(cmd)
}
