package server

import (
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func tokenFromRequest(r *http.Request) (string, error) {

	c, err := r.Cookie(CookieNameUser.String())
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

type contextKey string

const (
	userCtx contextKey = "user"
	orgCtx  contextKey = "org"
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

		})
	}
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", r.Context().Value(middleware.RequestIDKey).(string))
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }

func jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
