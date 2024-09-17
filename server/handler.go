package server

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
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
			workspace := getWorkspaceFromContext(r.Context()).ID.String()
			logger = logger.With(zap.String("workspace_id", workspace))
			span.SetAttributes(
				attribute.String("workspace_id", workspace))
		}

		if doesUserExistInContext(r.Context()) {
			userID := getUserFromContext(r.Context()).ID.String()
			logger = logger.With(zap.String("user_id", userID))
			span.SetAttributes(
				attribute.String("user_id", userID))
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

		err := render.Render(w, r, resp)
		if err != nil {
			logger.Error("could not write http response", zap.Error(err))
		}
	}
}
