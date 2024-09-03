package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func tokenFromRequest(r *http.Request) (string, error) {

	ss := strings.Split(r.Header.Get("Authorization"), "Bearer")
	if len(ss) != 2 {
		return "", errors.New("bearer token not found")
	}

	return ss[1], nil
}

type contextKey string

const (
	userCtx      contextKey = "user"
	workspaceCtx            = "workspace"
)

func requireAuthentication(
	logger *logrus.Entry,
	jwtManager jwttoken.JWTokenManager,
	cfg config.Config,
	userRepo malak.UserRepository,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx, span, rid := getTracer(r.Context(), r, "middleware.requireAuthentication", cfg.Otel.IsEnabled)
			defer span.End()

			logger = logger.WithField("request_id", rid)

			token, err := tokenFromRequest(r)
			if err != nil {
				logger.WithError(err).Error("token not found in cookie")
				_ = render.Render(w, r, newAPIStatus(http.StatusUnauthorized, "session expired"))
				return
			}

			data, err := jwtManager.ParseJWToken(token)
			if err != nil {
				logger.WithError(err).Error("could not parse JWT token")
				_ = render.Render(w, r, newAPIStatus(http.StatusUnauthorized, "could not validate JWT token"))
				return
			}

			user, err := userRepo.Get(ctx, &malak.FindUserOptions{
				ID: data.UserID,
			})
			if errors.Is(err, malak.ErrUserNotFound) {
				_ = render.Render(w, r, newAPIStatus(http.StatusForbidden, "user does not exists. JWT is invalid"))
				return
			}

			if err != nil {
				logger.WithError(err).Error("could not fetch user from database")
				_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "an error occurred while checking user"))
				return
			}

			r = r.WithContext(writeUserToCtx(ctx, user))

			next.ServeHTTP(w, r)
		})
	}
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", retrieveRequestID(r))
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }

func writeUserToCtx(ctx context.Context, user *malak.User) context.Context {
	return context.WithValue(ctx, userCtx, user)
}

func getUserFromContext(ctx context.Context) *malak.User {
	return ctx.Value(userCtx).(*malak.User)
}

func jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
