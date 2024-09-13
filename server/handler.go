package server

import (
	"context"
	"net/http"
	"os"

	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ENUM(success,failed)
type Status uint8

// MalakHTTPHandler is a wrapper for HTTP handlers
// that helps to centralizes error handling, otel tracing amongst others
type MalakHTTPHandler func(
	context.Context,
	trace.Span,
	*logrus.Entry,
	http.ResponseWriter,
	*http.Request) (render.Renderer, Status)

// WrapMalakHTTPHandler is a middleware that wraps our handlers and manages errors
func WrapMalakHTTPHandler(
	handler MalakHTTPHandler,
	cfg config.Config,
	spanName string) http.HandlerFunc {

	h, _ := os.Hostname()

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span, rid := getTracer(r.Context(), r, spanName, cfg.Otel.IsEnabled)
		defer span.End()

		logger := logrus.WithField("host", h).
			WithField("app", "malak.http").
			WithField("method", spanName).
			WithField("request_id", rid).
			WithContext(ctx)

		if doesWorkspaceExistInContext(r.Context()) {
			logger = logger.WithField("workspace_id", getWorkspaceFromContext(r.Context()).ID)
		}

		if doesUserExistInContext(r.Context()) {
			logger = logger.WithField("user_id", getUserFromContext(r.Context()).ID)
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
