package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("malak.server")
var noopTracer = noop.NewTracerProvider().Tracer("malak.server")

func getTracer(ctx context.Context, r *http.Request,
	operationName string, isEnabled bool) (context.Context, trace.Span, string) {

	rid := retrieveRequestID(r)
	if !isEnabled {
		ctx, span := noopTracer.Start(ctx, operationName)
		return ctx, span, rid
	}

	ctx, span := tracer.Start(ctx, operationName)

	span.SetAttributes(attribute.String("request_id", rid))
	return ctx, span, rid
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

func InitOTELCapabilities(cfg config.Config, logger *zap.Logger) func() {

	if !cfg.Otel.IsEnabled {
		return func() {}
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{}),
	)

	resources, err := initResources()
	if err != nil {
		logger.Fatal("could not setup OTEL tracing resources",
			zap.Error(err))
	}

	var (
		tracesSuffixEndpoint  = "/v1/traces"
		metricsSuffixEndpoint = "/v1/metrics"
	)

	headers := map[string]string{}
	pairs := strings.Split(cfg.Otel.Headers, ",")
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
	splittedEndpoint := strings.Split(cfg.Otel.Endpoint, "/")

	if len(splittedEndpoint) == 2 {
		// pick out the host
		cfg.Otel.Endpoint = splittedEndpoint[0]

		// make sure to use the remaining path and prepend to the actual
		// standard /v1 paths
		tracesSuffixEndpoint = splittedEndpoint[1] + tracesSuffixEndpoint
		metricsSuffixEndpoint = splittedEndpoint[1] + metricsSuffixEndpoint
	}

	var traceOptions = []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Otel.Endpoint),
		otlptracegrpc.WithHeaders(headers),
	}

	if !cfg.Otel.UseTLS {
		traceOptions = append(traceOptions, otlptracegrpc.WithInsecure())
	}

	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(traceOptions...))

	if err != nil {
		logger.Fatal("could not setup OTEL tracing resources",
			zap.Error(err))
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(traceExporter,
				sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
				sdktrace.WithBatchTimeout(5*time.Second),
			),
			sdktrace.WithResource(resources),
		),
	)

	var metricsOptions = []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Otel.Endpoint),
		otlpmetricgrpc.WithHeaders(headers),
	}

	if !cfg.Otel.UseTLS {
		metricsOptions = append(metricsOptions, otlpmetricgrpc.WithInsecure())
	}

	metricExporter, err := otlpmetricgrpc.New(
		context.Background(), metricsOptions...)
	if err != nil {
		logger.Fatal("could not setup metrics exporter",
			zap.Error(err))
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

func regiterMetrics(logger *zap.Logger) {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		logger.Fatal("could not gather runtime metrics",
			zap.Error(err))
	}
}
