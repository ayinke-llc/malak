package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateUpdateContentTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	req                contentUpdateRequest
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		req                contentUpdateRequest
		expectedStatusCode int
	}{
		{
			name: "empty content",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {

			},
			req:                contentUpdateRequest{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "empty title",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {

			},
			req: contentUpdateRequest{
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "title not up to 5 chars",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {

			},
			req: contentUpdateRequest{
				Title: "tt",
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "update not exists",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrUpdateNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: contentUpdateRequest{
				Title: "Valid title",
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
		},
		{
			name: "could not fetch update from db",
			req: contentUpdateRequest{
				Title: "Valid title",
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("error while fetching"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "updating content failed",
			req: contentUpdateRequest{
				Title: "Valid title",
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "updating content succeeds",
			req: contentUpdateRequest{
				Title: "Valid title",
				Update: []malak.Block{
					{
						ID:   "here is an id",
						Type: "heading",
					},
				},
			},
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_UpdateContent(t *testing.T) {
	for _, v := range generateUpdateContentTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)

			v.mockFn(updateRepo)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_jfnkfjkf")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), u.update, getConfig(), "updates.new").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestUpdatesHandler_Create(t *testing.T) {
	for _, v := range generateUpdateCreateTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)

			v.mockFn(updateRepo)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
				uuidGenerator:      &mockUUIDGenerator{},
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), u.create, getConfig(), "updates.new").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateCreateTestTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	expectedStatusCode int
	req                createUpdateContent
} {

	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		expectedStatusCode int
		req                createUpdateContent
	}{
		{
			name: "title not provided",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "title not up to 5 chars",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: createUpdateContent{
				Title: "abc",
			},
		},
		{
			name: "culd not create update",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not create update"))

			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createUpdateContent{
				Title: "omo let it go",
			},
		},
		{
			name: "created update",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			req: createUpdateContent{
				Title: "omo let it go",
			},
			expectedStatusCode: http.StatusCreated,
		},
	}
}

func generateUpdateListTestTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		expectedStatusCode int
	}{
		{
			name: "culd not list update",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("could not list update"))

			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "listed updates",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Update{
						{
							Reference: malak.Reference("update_12345"),
						},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func generateUpdateListPinsTestTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		expectedStatusCode int
	}{
		{
			name: "culd not list update",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					ListPinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not list update"))

			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "listed pinned updates",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					ListPinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Update{
						{
							Reference: malak.Reference("update_12345"),
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_ListPins(t *testing.T) {
	for _, v := range generateUpdateListPinsTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)

			v.mockFn(updateRepo)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
			}

			var b = bytes.NewBuffer(nil)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), u.listPinnedUpdates, getConfig(), "updates.list.pinned").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestUpdatesHandler_List(t *testing.T) {
	for _, v := range generateUpdateListTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)

			v.mockFn(updateRepo)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
			}

			var b = bytes.NewBuffer(nil)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), u.list, getConfig(), "updates.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestUpdatesHandler_PreviewUpdate(t *testing.T) {
	for _, v := range generatePreviewUpdateTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)
			cache := malak_mocks.NewMockCache(controller)
			queueHandler := malak_mocks.NewMockQueueHandler(controller)

			v.mockFn(updateRepo, cache, queueHandler)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
				cache:              cache,
				queueHandler:       queueHandler,
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/workspaces/updates/update_123/preview", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_123")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), u.previewUpdate, getConfig(), "updates.preview").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generatePreviewUpdateTestTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler)
	req                previewUpdateRequest
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler)
		req                previewUpdateRequest
		expectedStatusCode int
	}{
		{
			name: "email not provided",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
			},
			req:                previewUpdateRequest{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
			},
			req: previewUpdateRequest{
				Email: "invalid-email",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "preview throttled",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			req: previewUpdateRequest{
				Email: "valid@example.com",
			},
			expectedStatusCode: http.StatusTooManyRequests,
		},
		{
			name: "update not found",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, malak.ErrUpdateNotFound)
			},
			req: previewUpdateRequest{
				Email: "valid@example.com",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching update",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			req: previewUpdateRequest{
				Email: "valid@example.com",
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error creating preview",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&malak.Update{}, nil)
				update.EXPECT().SendUpdate(gomock.Any(), gomock.Any()).Return(errors.New("unknown error"))
			},
			req: previewUpdateRequest{
				Email: "valid@example.com",
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successful preview creation",
			mockFn: func(update *malak_mocks.MockUpdateRepository, cache *malak_mocks.MockCache, queueHandler *malak_mocks.MockQueueHandler) {
				cache.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("not found"))
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&malak.Update{}, nil)
				update.EXPECT().SendUpdate(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: previewUpdateRequest{
				Email: "valid@example.com",
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_SendUpdate(t *testing.T) {
	for _, v := range generateUpdateTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			updateRepo := malak_mocks.NewMockUpdateRepository(controller)

			v.mockFn(updateRepo)

			u := &updatesHandler{
				referenceGenerator: &mockReferenceGenerator{},
				updateRepo:         updateRepo,
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/workspaces/updates/update_123/preview", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_123")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), u.sendUpdate, getConfig(), "updates.send").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateTestTable() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	req                sendUpdateRequest
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		req                sendUpdateRequest
		expectedStatusCode int
	}{
		{
			name: "email not provided",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
			},
			req:                sendUpdateRequest{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid email",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "update not found",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, malak.ErrUpdateNotFound)
			},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@oops.com",
				},
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching update",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@oops.com",
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error during update sending",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&malak.Update{}, nil)
				update.EXPECT().SendUpdate(gomock.Any(), gomock.Any()).Return(errors.New("oops"))
			},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@oops.com",
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update sending reached max recipients",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&malak.Update{}, nil)
				update.EXPECT().SendUpdate(gomock.Any(), gomock.Any()).
					Return(malak.ErrCounterExhausted)
			},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@oops.com",
				},
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name: "successful update sending",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&malak.Update{}, nil)
				update.EXPECT().SendUpdate(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: sendUpdateRequest{
				Emails: []malak.Email{
					"oops@oops.com",
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

type mockUUIDGenerator struct{}

func (m *mockUUIDGenerator) Create() uuid.UUID {
	return uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
}
