package server

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ENUM(success,failed)
type Status uint8

// MalakHTTPHandler is a wrapper for HTTP handlers
// that helps to centralizes error handling, otel tracing amongst others
type MalakHTTPHandler func(
	context.Context,
	trace.Span,
	*zap.Logger,
	http.ResponseWriter,
	*http.Request) (render.Renderer, Status)

// WrapMalakHTTPHandler is a middleware that wraps our handlers and manages errors
func WrapMalakHTTPHandler(
	logger *zap.Logger,
	handler MalakHTTPHandler,
	cfg config.Config,
	spanName string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span, rid := getTracer(r.Context(), r, spanName, cfg.Otel.IsEnabled)
		defer span.End()

		logger = logger.With(zap.String("request_id", rid))

		if doesWorkspaceExistInContext(r.Context()) {
			logger = logger.With(zap.String("workspace_id", getWorkspaceFromContext(r.Context()).ID.String()))
		}

		if doesUserExistInContext(r.Context()) {
			logger = logger.With(zap.String("user_id", getUserFromContext(r.Context()).ID.String()))
		}

		resp, status := handler(ctx, span, logger, w, r)
		switch status {
		case StatusFailed:

			span.SetStatus(codes.Error, "")

		case StatusSuccess:
			span.SetStatus(codes.Ok, "")

		default:
			_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "unknown error"))
			return
		}

		_ = render.Render(w, r, resp)
	}
}
