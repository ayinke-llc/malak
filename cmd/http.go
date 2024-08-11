package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
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

func addHTTPCommand(c *cobra.Command) {

	cmd := &cobra.Command{
		Use: "http",
		Run: func(cmd *cobra.Command, args []string) {

			sig := make(chan os.Signal, 1)

			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			logrus.SetOutput(os.Stdout)
			logrus.SetFormatter(&logrus.JSONFormatter{})

			cfg, err := config.Load()
			if err != nil {
				logrus.WithError(err).Fatal("could not load configuration details")
			}

			json.NewEncoder(os.Stdout).Encode(cfg)

			lvl, err := logrus.ParseLevel(cfg.LogLevel)
			if err != nil {
				lvl = logrus.DebugLevel
			}

			logrus.SetLevel(lvl)

			h, _ := os.Hostname()

			logger := logrus.WithField("host", h).
				WithField("app", "malak")

			cleanup := initOTELCapabilities(cfg, logger)

			opts, err := redis.ParseURL("dood")
			if err != nil {
				log.Fatal(err)
			}

			redisClient := redis.NewClient(opts)

			if err := redisotel.InstrumentTracing(redisClient); err != nil {
				log.Fatal(err)
			}

			if err := redisotel.InstrumentMetrics(redisClient); err != nil {
				log.Fatal(err)
			}

			if err := redisClient.Ping(context.Background()).Err(); err != nil {
				log.Fatal(err)
			}

			<-sig

			logger.Debug("shutting down server")
			cleanup()
		},
	}

	c.AddCommand(cmd)
}

func initResources() (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "malak"),
			attribute.String("library.language", "go"),
		),
	)
}

func initOTELCapabilities(cfg config.Config, logger *logrus.Entry) func() {

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{}),
	)

	resources, err := initResources()
	if err != nil {
		logger.WithError(err).Fatal("could not setup OTEL tracing resources")
	}

	var (
		tracesSuffixEndpoint  = "/v1/traces"
		metricsSuffixEndpoint = "/v1/metrics"
	)

	headers := map[string]string{}
	pairs := strings.Split(cfg.OtelHeaders, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	// By default, Otel sends traces and metrics, logs to v1/* paths
	// but some providers like Grafana have their OTEL collector on a subpath
	// so /otlp/v1/*
	// The sdk is pretty stringent as that format does not match the standard
	// so it doesn't accept, this makes sure to split out the url and make it match
	splittedEndpoint := strings.Split(cfg.OtelEndpoint, "/")

	if len(splittedEndpoint) == 2 {
		// pick out the host
		cfg.OtelEndpoint = splittedEndpoint[0]

		// make sure to use the remaining path and prepend to the actual
		// standard /v1 paths
		tracesSuffixEndpoint = splittedEndpoint[1] + tracesSuffixEndpoint
		metricsSuffixEndpoint = splittedEndpoint[1] + metricsSuffixEndpoint
	}

	var traceOptions = []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.OtelEndpoint),
		otlptracehttp.WithURLPath(tracesSuffixEndpoint),
		otlptracehttp.WithHeaders(headers),
	}

	if !cfg.OtelUseTLS {
		traceOptions = append(traceOptions, otlptracehttp.WithInsecure())
	}

	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(traceOptions...))

	if err != nil {
		logger.WithError(err).Fatal("could not setup OTEL tracing")
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(traceExporter,
				sdktrace.WithMaxExportBatchSize(maxExportBatchSize),
				sdktrace.WithBatchTimeout(5*time.Second),
			),
			sdktrace.WithResource(resources),
		),
	)

	var metricsOptions = []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.OtelEndpoint),
		otlpmetrichttp.WithURLPath(metricsSuffixEndpoint),
		otlpmetrichttp.WithHeaders(headers),
	}

	if !cfg.OtelUseTLS {
		metricsOptions = append(metricsOptions, otlpmetrichttp.WithInsecure())
	}

	metricExporter, err := otlpmetrichttp.New(
		context.Background(), metricsOptions...)
	if err != nil {
		logger.WithError(err).Fatal("could not set up Metrics exporter")
	}

	otel.SetMeterProvider(
		metric.NewMeterProvider(
			metric.WithResource(resources),
			metric.WithReader(
				metric.NewPeriodicReader(metricExporter))))

	regiterMetrics(logger)

	return func() {
		_ = traceExporter.Shutdown(context.Background())
		_ = metricExporter.Shutdown(context.Background())
	}
}

func regiterMetrics(logger *logrus.Entry) {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		logger.WithError(err).Fatal("could not gather runtime metrics")
	}
}
