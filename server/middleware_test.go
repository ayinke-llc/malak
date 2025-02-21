package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	mock_jwttoken "github.com/ayinke-llc/malak/internal/pkg/jwttoken/mocks"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRequireWorkspaceValidSubscription(t *testing.T) {
	t.Run("sub not active", func(t *testing.T) {

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

		requireWorkspaceValidSubscription(getConfig())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode("{}"))
		})).ServeHTTP(rr, req)

		require.Equal(t, http.StatusPaymentRequired, rr.Code)
		verifyMatch(t, rr)
	})

	t.Run("sub active", func(t *testing.T) {

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Add("Content-Type", "application/json")
		req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
			IsSubscriptionActive: true,
		}))

		requireWorkspaceValidSubscription(getConfig())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode("{}"))
		})).ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		verifyMatch(t, rr)
	})
}

func TestTokenFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "valid bearer token",
			authHeader:    "Bearer abc123",
			expectedToken: "abc123",
			expectError:   false,
		},
		{
			name:        "missing bearer prefix",
			authHeader:  "abc123",
			expectError: true,
		},
		{
			name:        "empty auth header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "malformed bearer token",
			authHeader:  "Bearer",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := tokenFromRequest(req)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestGetIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expectedIP string
	}{
		{
			name: "cloudflare ip",
			headers: map[string]string{
				"CF-Connecting-IP": "1.2.3.4",
			},
			expectedIP: "1.2.3.4",
		},
		{
			name: "x-forwarded-for single ip",
			headers: map[string]string{
				"X-Forwarded-For": "5.6.7.8",
			},
			expectedIP: "5.6.7.8",
		},
		{
			name: "x-forwarded-for multiple ips",
			headers: map[string]string{
				"X-Forwarded-For": "9.10.11.12, 13.14.15.16",
			},
			expectedIP: "9.10.11.12",
		},
		{
			name: "x-real-ip",
			headers: map[string]string{
				"X-Real-IP": "17.18.19.20",
			},
			expectedIP: "17.18.19.20",
		},
		{
			name:       "remote addr fallback",
			remoteAddr: "21.22.23.24",
			expectedIP: "21.22.23.24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			ip := getIP(req)
			require.Equal(t, tt.expectedIP, ip)
		})
	}
}

func TestHTTPThrottleKeyFunc(t *testing.T) {
	t.Run("authenticated user", func(t *testing.T) {
		userID := uuid.New()
		user := &malak.User{ID: userID}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := writeUserToCtx(req.Context(), user)
		req = req.WithContext(ctx)

		key, err := HTTPThrottleKeyFunc(req)
		require.NoError(t, err)
		require.Equal(t, userID.String(), key)
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("CF-Connecting-IP", "1.2.3.4")

		key, err := HTTPThrottleKeyFunc(req)
		require.NoError(t, err)
		require.Equal(t, "1.2.3.4", key)
	})
}

func TestContextHelpers(t *testing.T) {
	t.Run("user context", func(t *testing.T) {
		ctx := t.Context()
		user := &malak.User{ID: uuid.New()}

		// Test writing and reading user
		ctx = writeUserToCtx(ctx, user)
		require.True(t, doesUserExistInContext(ctx))
		require.Equal(t, user, getUserFromContext(ctx))
	})

	t.Run("workspace context", func(t *testing.T) {
		ctx := t.Context()
		workspace := &malak.Workspace{ID: uuid.New()}

		// Test writing and reading workspace
		ctx = writeWorkspaceToCtx(ctx, workspace)
		require.True(t, doesWorkspaceExistInContext(ctx))
		require.Equal(t, workspace, getWorkspaceFromContext(ctx))
	})
}

func TestJsonResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler := jsonResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"test": true}`))
		if err != nil {
			t.Fatal(err)
		}
	}))

	handler.ServeHTTP(rr, req)

	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestRequireAuthentication(t *testing.T) {
	logger := zap.NewNop()
	cfg := getConfig()

	t.Run("missing token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		jwtManager.EXPECT().
			ParseJWToken("invalid-token").
			Return(jwttoken.JWTokenData{}, errors.New("invalid token"))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		jwtManager.EXPECT().
			ParseJWToken("valid-token").
			Return(jwttoken.JWTokenData{UserID: userID}, nil)

		userRepo.EXPECT().
			Get(gomock.Any(), &malak.FindUserOptions{ID: userID}).
			Return(nil, errors.New("user not found"))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("user without workspace accessing protected route", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		jwtManager.EXPECT().
			ParseJWToken("valid-token").
			Return(jwttoken.JWTokenData{UserID: userID}, nil)

		userRepo.EXPECT().
			Get(gomock.Any(), &malak.FindUserOptions{ID: userID}).
			Return(&malak.User{
				ID: userID,
				Metadata: &malak.UserMetadata{
					CurrentWorkspace: uuid.Nil, // No workspace
				},
			}, nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)

		// Verify error message
		var response map[string]string
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "You must be a member of a workspace", response["message"])
	})

	t.Run("user without workspace accessing auth/connect route", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		jwtManager.EXPECT().
			ParseJWToken("valid-token").
			Return(jwttoken.JWTokenData{UserID: userID}, nil)

		userRepo.EXPECT().
			Get(gomock.Any(), &malak.FindUserOptions{ID: userID}).
			Return(&malak.User{
				ID: userID,
				Metadata: &malak.UserMetadata{
					CurrentWorkspace: uuid.Nil, // No workspace
				},
			}, nil)

		// For auth/connect route, workspace repository should never be called
		// since we check the path before attempting to fetch workspace

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/auth/connect", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For auth/connect route, we should still have user in context
			user := getUserFromContext(r.Context())
			require.Equal(t, userID, user.ID)
			require.Equal(t, uuid.Nil, user.Metadata.CurrentWorkspace)

			// Workspace should not be in context since we never fetched it
			require.False(t, doesWorkspaceExistInContext(r.Context()))

			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("user with workspace accessing protected route", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		workspaceID := uuid.New()
		jwtManager := mock_jwttoken.NewMockJWTokenManager(ctrl)
		userRepo := malak_mocks.NewMockUserRepository(ctrl)
		workspaceRepo := malak_mocks.NewMockWorkspaceRepository(ctrl)

		jwtManager.EXPECT().
			ParseJWToken("valid-token").
			Return(jwttoken.JWTokenData{UserID: userID}, nil)

		userRepo.EXPECT().
			Get(gomock.Any(), &malak.FindUserOptions{ID: userID}).
			Return(&malak.User{
				ID: userID,
				Metadata: &malak.UserMetadata{
					CurrentWorkspace: workspaceID,
				},
			}, nil)

		workspaceRepo.EXPECT().
			Get(gomock.Any(), &malak.FindWorkspaceOptions{ID: workspaceID}).
			Return(&malak.Workspace{ID: workspaceID}, nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		handler := requireAuthentication(
			logger,
			jwtManager,
			cfg,
			userRepo,
			workspaceRepo,
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify user and workspace are in context
			user := getUserFromContext(r.Context())
			require.Equal(t, userID, user.ID)

			workspace := getWorkspaceFromContext(r.Context())
			require.Equal(t, workspaceID, workspace.ID)

			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})
}
