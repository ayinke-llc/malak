package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func tokenFromRequest(r *http.Request) (string, error) {

	ss := strings.Split(r.Header.Get("Authorization"), " ")

	if len(ss) != 2 {
		return "", errors.New("invalid bearer token structure")
	}

	return ss[1], nil
}

type contextKey string

const (
	userCtx      contextKey = "user"
	workspaceCtx contextKey = "workspace"
)

// HTTPThrottleKeyFunc throttles unauthenticated users by their IP.
// It goes through cloudflare, X-Forwarded-For and X-Real-IP to determine the
// correct IP
//
// For authenticated requests, it throttles individually instead of IP wild
func HTTPThrottleKeyFunc(r *http.Request) (string, error) {
	if doesUserExistInContext(r.Context()) {
		return getUserFromContext(r.Context()).ID.String(), nil
	}

	return getIP(r), nil
}

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func getIP(r *http.Request) string {

	cloudflareIP := r.Header.Get("CF-Connecting-IP")
	if !util.IsStringEmpty(cloudflareIP) {
		return cloudflareIP
	}

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")

		if i == -1 {
			i = len(xff)
		}

		ip := xff[:i]
		if !util.IsStringEmpty(ip) {
			return ip
		}
	}

	xIP := r.Header.Get(xRealIP)
	if !util.IsStringEmpty(xIP) {
		return xIP
	}

	return r.RemoteAddr
}

func requireWorkspaceValidSubscription(
	cfg config.Config,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx, span, _ := getTracer(r.Context(), r, "middleware.requireWorkspaceValidSubscription", cfg.Otel.IsEnabled)
			defer span.End()

			if r.URL.Path == "/v1/workspaces/billing" ||
				r.URL.Path == "/v1/workspaces" ||
				r.URL.Path == "/v1/user" {
				next.ServeHTTP(w, r)
				return
			}

			workspace := getWorkspaceFromContext(ctx)
			if !workspace.IsSubscriptionActive {
				_ = render.Render(w, r, newAPIStatus(http.StatusPaymentRequired,
					"workspace is not active. You need an active subscription"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func requireAuthentication(
	logger *zap.Logger,
	jwtManager jwttoken.JWTokenManager,
	cfg config.Config,
	userRepo malak.UserRepository,
	workspaceRepo malak.WorkspaceRepository,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx, span, rid := getTracer(r.Context(), r, "middleware.requireAuthentication", cfg.Otel.IsEnabled)
			defer span.End()

			logger = logger.With(
				zap.String("request_id", rid),
				zap.String("path", r.URL.Path),
			)

			token, err := tokenFromRequest(r)
			if err != nil {
				logger.Error("token not found in cookie", zap.Error(err))
				_ = render.Render(w, r, newAPIStatus(http.StatusUnauthorized, "session expired"))
				return
			}

			data, err := jwtManager.ParseJWToken(token)
			if err != nil {
				logger.Error("could not parse JWT", zap.Error(err))
				_ = render.Render(w, r, newAPIStatus(http.StatusUnauthorized, "could not validate JWT token"))
				return
			}

			user, err := userRepo.Get(ctx, &malak.FindUserOptions{
				ID: data.UserID,
			})
			if err != nil {
				logger.Error("could not fetch user from database", zap.Error(err))
				_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "an error occurred while checking user"))
				return
			}

			r = r.WithContext(writeUserToCtx(ctx, user))

			if user.Metadata.CurrentWorkspace != uuid.Nil {
				workspace, err := workspaceRepo.Get(ctx, &malak.FindWorkspaceOptions{
					ID: user.Metadata.CurrentWorkspace,
				})
				if err != nil {
					logger.Error("could not fetch workspace from database", zap.Error(err))
					_ = render.Render(w, r, newAPIStatus(http.StatusInternalServerError, "an error occurred while fetching workspace from database"))
					return
				}

				r = r.WithContext(writeWorkspaceToCtx(r.Context(), workspace))
			}

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

func writeWorkspaceToCtx(ctx context.Context, workspace *malak.Workspace) context.Context {
	return context.WithValue(ctx, workspaceCtx, workspace)
}

func getWorkspaceFromContext(ctx context.Context) *malak.Workspace {
	return ctx.Value(workspaceCtx).(*malak.Workspace)
}

func doesWorkspaceExistInContext(ctx context.Context) bool {
	_, ok := ctx.Value(workspaceCtx).(*malak.Workspace)
	return ok
}

func doesUserExistInContext(ctx context.Context) bool {
	_, ok := ctx.Value(userCtx).(*malak.User)
	return ok
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
