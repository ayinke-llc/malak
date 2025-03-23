package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateCreateAPIKeyTestTable() []struct {
	name               string
	mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient)
	req                createAPIKeyRequest
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient)
		req                createAPIKeyRequest
		expectedStatusCode int
	}{
		{
			name: "empty title",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {
			},
			req:                createAPIKeyRequest{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "title too short",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {
			},
			req: createAPIKeyRequest{
				Title: "ab",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "title too long",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {
			},
			req: createAPIKeyRequest{
				Title: "this is a very long title that exceeds twenty characters",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "api key creation error",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {

				apiRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("creation error"))
			},
			req: createAPIKeyRequest{
				Title: "Valid Title",
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "max limit reached",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {

				apiRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.ErrAPIKeyMaxLimit)
			},
			req: createAPIKeyRequest{
				Title: "Valid Title",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "successful creation",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository, secretsClient *malak_mocks.MockSecretClient) {

				apiRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			req: createAPIKeyRequest{
				Title: "Valid Title",
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestAPIKeyHandler_Create(t *testing.T) {
	for _, v := range generateCreateAPIKeyTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			apiRepo := malak_mocks.NewMockAPIKeyRepository(controller)
			secretsClient := malak_mocks.NewMockSecretClient(controller)

			v.mockFn(apiRepo, secretsClient)

			h := &apiKeyHandler{
				apiRepo:   apiRepo,
				generator: &mockReferenceGenerator{},
				cfg:       getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			userID := uuid.New()
			workspaceID := uuid.New()
			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{ID: userID}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			WrapMalakHTTPHandler(getLogger(t), h.create, getConfig(), "developers.keys.create").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateListAPIKeysTestTable() []struct {
	name               string
	mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository)
		expectedStatusCode int
	}{
		{
			name: "list error",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				apiRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("list error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successful list",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				apiRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.APIKey{
						{
							Reference: "api_key_123",
							KeyName:   "Test Key",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestAPIKeyHandler_List(t *testing.T) {
	for _, v := range generateListAPIKeysTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			apiRepo := malak_mocks.NewMockAPIKeyRepository(controller)

			v.mockFn(apiRepo)

			h := &apiKeyHandler{
				apiRepo:   apiRepo,
				generator: &mockReferenceGenerator{},
				cfg:       getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			userID := uuid.New()
			workspaceID := uuid.New()
			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{ID: userID}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			WrapMalakHTTPHandler(getLogger(t), h.list, getConfig(), "developers.keys.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateRevokeAPIKeyTestTable() []struct {
	name               string
	mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository)
	req                revokeAPIKeyRequest
	reference          string
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(apiRepo *malak_mocks.MockAPIKeyRepository)
		req                revokeAPIKeyRequest
		reference          string
		expectedStatusCode int
	}{
		{
			name: "empty reference",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
			},
			req:                revokeAPIKeyRequest{},
			reference:          "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid strategy",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
			},
			req: revokeAPIKeyRequest{
				Strategy: "invalid",
			},
			reference:          "api_key_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "api key not found",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				apiRepo.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrAPIKeyNotFound)
			},
			req: revokeAPIKeyRequest{
				Strategy: malak.RevocationTypeImmediate,
			},
			reference:          "api_key_123",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "api key already revoked",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				now := time.Now()
				apiRepo.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.APIKey{
						Reference: "api_key_123",
						ExpiresAt: &now,
					}, nil)
			},
			req: revokeAPIKeyRequest{
				Strategy: malak.RevocationTypeImmediate,
			},
			reference:          "api_key_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "revocation error",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				apiRepo.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.APIKey{
						Reference: "api_key_123",
					}, nil)

				apiRepo.EXPECT().
					Revoke(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("revocation error"))
			},
			req: revokeAPIKeyRequest{
				Strategy: malak.RevocationTypeImmediate,
			},
			reference:          "api_key_123",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successful revocation",
			mockFn: func(apiRepo *malak_mocks.MockAPIKeyRepository) {
				apiRepo.EXPECT().
					Fetch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.APIKey{
						Reference: "api_key_123",
					}, nil)

				apiRepo.EXPECT().
					Revoke(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			req: revokeAPIKeyRequest{
				Strategy: malak.RevocationTypeImmediate,
			},
			reference:          "api_key_123",
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestAPIKeyHandler_Revoke(t *testing.T) {
	for _, v := range generateRevokeAPIKeyTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			apiRepo := malak_mocks.NewMockAPIKeyRepository(controller)

			v.mockFn(apiRepo)

			h := &apiKeyHandler{
				apiRepo:   apiRepo,
				generator: &mockReferenceGenerator{},
				cfg:       getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/"+v.reference, b)
			req.Header.Add("Content-Type", "application/json")

			userID := uuid.New()
			workspaceID := uuid.New()
			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{ID: userID}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			// Add the reference parameter to the request context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("reference", v.reference)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			WrapMalakHTTPHandler(getLogger(t), h.revoke, getConfig(), "developers.keys.revoke").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
