package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gulter "github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHashURL(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		expected string
		hasError bool
	}{
		{
			name:     "example.com with query",
			value:    "https://example.com/path?query=123",
			expected: "deck-125bf3f1de0f5189",
			hasError: false,
		},
		{
			name:     "example.com",
			value:    "https://example.com",
			expected: "deck-837b2b5793a240b3",
			hasError: false,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {

			val, err := hashURL(v.value)
			if v.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, v.expected, val)
		})
	}
}

func generateDeckCreateRequest() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache)
	expectedStatusCode int
	req                createDeckRequest
} {

	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache)
		expectedStatusCode int
		req                createDeckRequest
	}{
		{
			name:               "no url provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckRequest{
				DeckURL: "",
			},
		},
		{
			name:               "no title provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckRequest{
				DeckURL: "https://google.com",
			},
		},
		{
			name: "file not exists in cache",
			mockFn: func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return([]byte(``), errors.New("could not fetch file"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckRequest{
				DeckURL: "https://google.com",
				Title:   "oops",
			},
		},
		{
			name: "file size in cache is 0",
			mockFn: func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return([]byte(`{"size: 0}`), nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckRequest{
				DeckURL: "https://google.com",
				Title:   "oops",
			},
		},
		{
			name: "could not create deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {

				cache.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return([]byte(`{"Size": 1000000000}`), nil)

				deck.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckRequest{
				DeckURL: "https://google.com",
				Title:   "oops",
			},
		},
		{
			name: "created deck successfully",
			mockFn: func(deck *malak_mocks.MockDeckRepository, cache *malak_mocks.MockCache) {

				cache.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return([]byte(`{"Size": 1000000000}`), nil)

				deck.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: createDeckRequest{
				DeckURL: "https://google.com",
				Title:   "oops",
			},
		},
	}
}

func TestDeckHandler_Create(t *testing.T) {

	for _, v := range generateDeckCreateRequest() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			cacheRepo := malak_mocks.NewMockCache(controller)

			v.mockFn(deckRepo, cacheRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
				cfg:                getConfig(),
				cache:              cacheRepo,
			}

			var b = bytes.NewBuffer(nil)

			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_jfnkfjkf")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.Create,
				getConfig(), "decks.create").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckListRequest() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name: "could not list decks",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Deck{}, errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "listed deck successfully",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Deck{
						{
							Reference: "oops",
						},
						{
							Reference: "opsdfkf",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDeckHandler_List(t *testing.T) {

	for _, v := range generateDeckListRequest() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)

			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t),
				u.List,
				getConfig(), "decks.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckDeleteRequest() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name: "could not fetch deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "could not fetch deck because deck does not exists",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "could not delete deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "osfifjf",
					}, nil)

				deck.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not delete deck"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "deleted deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "osfifjf",
					}, nil)

				deck.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDeckHandler_Delete(t *testing.T) {

	for _, v := range generateDeckDeleteRequest() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)

			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "deck_djdnd")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.Delete,
				getConfig(), "decks.delete").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckFetchRequest() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name: "could not fetch deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "could not fetch deck because deck does not exists",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "fetched deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "osfifjf",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDeckHandler_Fetch(t *testing.T) {

	for _, v := range generateDeckFetchRequest() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)

			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "deck_djdnd")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.fetch,
				getConfig(), "decks.fetch").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckUpdatePreferencesRequest() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository)
	expectedStatusCode int
	req                updateDeckPreferencesRequest
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
		req                updateDeckPreferencesRequest
	}{
		{
			name:               "no reference provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req:                updateDeckPreferencesRequest{},
		},
		{
			name:               "password enabled but no password provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: updateDeckPreferencesRequest{
				PasswordProtection: struct {
					Enabled bool           `json:"enabled,omitempty"`
					Value   malak.Password `json:"value,omitempty"`
				}{
					Enabled: true,
				},
			},
		},
		{
			name: "could not fetch deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDeckPreferencesRequest{
				EnableDownloading: true,
				RequireEmail:      true,
			},
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: updateDeckPreferencesRequest{
				EnableDownloading: true,
				RequireEmail:      true,
			},
		},
		{
			name: "could not update preferences",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference:      "deck_test",
						DeckPreference: &malak.DeckPreference{},
					}, nil)

				deck.EXPECT().UpdatePreferences(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error occurred"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDeckPreferencesRequest{
				EnableDownloading: true,
				RequireEmail:      true,
			},
		},
		{
			name: "successfully updated preferences",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference:      "deck_test",
						DeckPreference: &malak.DeckPreference{},
					}, nil)

				deck.EXPECT().UpdatePreferences(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updateDeckPreferencesRequest{
				EnableDownloading: true,
				RequireEmail:      true,
				PasswordProtection: struct {
					Enabled bool           `json:"enabled,omitempty"`
					Value   malak.Password `json:"value,omitempty"`
				}{
					Enabled: true,
					Value:   "secure_password",
				},
			},
		},
	}
}

func TestDeckHandler_UpdatePreferences(t *testing.T) {
	for _, v := range generateDeckUpdatePreferencesRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.updatePreferences,
				getConfig(), "decks.update_preferences").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckToggleArchiveTestTable() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name:               "no reference provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error toggling archive status",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().ToggleArchive(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not toggle archive status"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully toggled archive status",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().ToggleArchive(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDeckHandler_ToggleArchive(t *testing.T) {
	for _, v := range generateDeckToggleArchiveTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/workspaces/decks/deck_test/archive", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.toggleArchive,
				getConfig(), "decks.toggle_archive").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckPins() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name:               "no reference provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error toggling pinned status",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().TogglePinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not toggle pinned status"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully toggled pinned status",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().TogglePinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDecksHandler_Pin(t *testing.T) {
	for _, v := range generateDeckPins() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/workspaces/decks/deck_test/pin", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}

			req = req.WithContext(
				context.WithValue(req.Context(),
					chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.togglePinned,
				getConfig(), "decks.toggle_archive").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func getExpectedCacheKey(s3Endpoint, folderDestination, fileName string) string {
	url := fmt.Sprintf("%s/%s/%s", s3Endpoint, folderDestination, fileName)
	key, err := hashURL(url)
	if err != nil {
		panic(err)
	}
	return key
}

func TestGetExpectedCacheKey(t *testing.T) {
	tt := []struct {
		name     string
		endpoint string
		folder   string
		file     string
		expected string
	}{
		{
			name:     "s3 example url",
			endpoint: "https://s3.example.com",
			folder:   "decks",
			file:     "test.pdf",
			expected: "deck-3dad11f2911c16ef",
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/%s/%s", v.endpoint, v.folder, v.file)
			t.Logf("URL being hashed: %s", url)
			key := getExpectedCacheKey(v.endpoint, v.folder, v.file)
			require.Equal(t, v.expected, key)
		})
	}
}

type uploadedFileMetadata struct {
	Size              int64
	FolderDestination string
	Key               string
}

type mockStorage struct {
	file *uploadedFileMetadata
}

func generateDeckUploadImageRequest() []struct {
	name               string
	mockFn             func(t *testing.T, cache *malak_mocks.MockCache)
	expectedStatusCode int
	file               *uploadedFileMetadata
} {
	return []struct {
		name               string
		mockFn             func(t *testing.T, cache *malak_mocks.MockCache)
		expectedStatusCode int
		file               *uploadedFileMetadata
	}{
		{
			name:               "no file uploaded",
			mockFn:             func(t *testing.T, cache *malak_mocks.MockCache) {},
			expectedStatusCode: http.StatusInternalServerError,
			file:               nil,
		},
		{
			name: "could not store file details in cache",
			mockFn: func(t *testing.T, cache *malak_mocks.MockCache) {
				cache.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not store in cache"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			file: &uploadedFileMetadata{
				Size:              1024,
				FolderDestination: "decks",
				Key:               "decks/test.pdf",
			},
		},
		{
			name: "successfully uploaded file",
			mockFn: func(t *testing.T, cache *malak_mocks.MockCache) {
				cache.EXPECT().Add(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Times(1).
					Do(func(_ context.Context, key string, value []byte, _ time.Duration) {
						t.Logf("Cache key being used: %s", key)
						var f struct {
							Size int64  `json:"Size"`
							Key  string `json:"Key"`
						}
						if err := json.Unmarshal(value, &f); err != nil {
							panic(err)
						}
						if f.Size != 1024 {
							panic("invalid file size")
						}
						if !strings.HasPrefix(f.Key, "decks/gulter-") || !strings.HasSuffix(f.Key, "test.pdf") {
							panic("invalid file key format")
						}
					}).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			file: &uploadedFileMetadata{
				Size:              1024,
				FolderDestination: "decks",
				Key:               "decks/test.pdf",
			},
		},
	}
}

func mockValidator(f gulter.File) error {
	return nil
}

func (m *mockStorage) Upload(ctx context.Context, r io.Reader, opts *gulter.UploadFileOptions) (*gulter.UploadedFileMetadata, error) {
	if m.file == nil {
		return nil, errors.New("no file")
	}
	return &gulter.UploadedFileMetadata{
		Size:              m.file.Size,
		FolderDestination: "decks",
		Key:               fmt.Sprintf("decks/%s", opts.FileName),
	}, nil
}

func (m *mockStorage) Path(ctx context.Context, opts gulter.PathOptions) (string, error) {
	if m.file == nil {
		return "", errors.New("no file")
	}
	return fmt.Sprintf("https://s3.amazonaws.com/%s", opts.Key), nil
}

func (m *mockStorage) Close() error {
	return nil
}

func TestDeckHandler_UploadImage(t *testing.T) {
	for _, v := range generateDeckUploadImageRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			cacheRepo := malak_mocks.NewMockCache(controller)
			v.mockFn(t, cacheRepo)

			cfg := getConfig()
			cfg.Uploader.S3.Endpoint = "https://s3.example.com"

			// Log the expected URL and key
			if v.file != nil {
				expectedURL := fmt.Sprintf("%s/%s/%s", cfg.Uploader.S3.Endpoint, "decks", "test.pdf")
				expectedKey, err := hashURL(expectedURL)
				require.NoError(t, err)
				t.Logf("Expected URL: %s", expectedURL)
				t.Logf("Expected key: %s", expectedKey)
			}

			u := &deckHandler{
				cache: cacheRepo,
				cfg:   cfg,
			}

			rr := httptest.NewRecorder()

			var req *http.Request
			if v.file != nil {
				// Create a multipart form request
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("image_body", "test.pdf")
				require.NoError(t, err)
				_, err = part.Write([]byte("test content"))
				require.NoError(t, err)
				err = writer.Close()
				require.NoError(t, err)

				req = httptest.NewRequest(http.MethodPost, "/", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				req = httptest.NewRequest(http.MethodPost, "/", nil)
				req.Header.Add("Content-Type", "application/json")
			}

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			if v.file != nil {
				// Set up gulter middleware with mock storage
				mockStorage := &mockStorage{
					file: v.file,
				}

				g, err := gulter.New(
					gulter.WithStorage(mockStorage),
					gulter.WithMaxFileSize(1024*1024*10), // 10MB
					gulter.WithIgnoreNonExistentKey(true),
					gulter.WithValidationFunc(mockValidator),
					gulter.WithNameFuncGenerator(func(s string) string {
						return fmt.Sprintf("gulter-1738628312-%s", s) // Fixed timestamp for consistent test output
					}),
				)
				require.NoError(t, err)

				// Wrap our handler with gulter middleware
				handler := g.Upload("decks", "image_body")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					files, err := gulter.FilesFromContextWithKey(r, "image_body")
					require.NoError(t, err)
					require.Len(t, files, 1)
					file := files[0]
					t.Logf("File after gulter processing: FolderDestination=%s, UploadedFileName=%s, StorageKey=%s",
						file.FolderDestination, file.UploadedFileName, file.StorageKey)

					WrapMalakHTTPHandler(getLogger(t),
						u.uploadImage,
						cfg, "decks.upload").
						ServeHTTP(w, r)
				}))

				handler.ServeHTTP(rr, req)
			} else {
				WrapMalakHTTPHandler(getLogger(t),
					u.uploadImage,
					cfg, "decks.upload").
					ServeHTTP(rr, req)
			}

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckSessionsTestTable() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository)
	expectedStatusCode int
	queryParams        map[string]string
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
		queryParams        map[string]string
	}{
		{
			name:               "no reference provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			queryParams:        map[string]string{},
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			queryParams:        map[string]string{},
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			queryParams:        map[string]string{},
		},
		{
			name: "invalid days parameter",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			queryParams: map[string]string{
				"days": "15", // Only 7, 14, or 30 are valid
			},
		},
		{
			name: "error fetching analytics",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().SessionAnalytics(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("analytics error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			queryParams:        map[string]string{},
		},
		{
			name: "successfully fetched sessions",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().SessionAnalytics(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]*malak.DeckViewerSession{
						{
							ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
			queryParams: map[string]string{
				"days":     "7",
				"page":     "1",
				"per_page": "10",
			},
		},
	}
}

func TestDeckHandler_FetchDeckSessions(t *testing.T) {
	for _, v := range generateDeckSessionsTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()

			// Build URL with query parameters
			url := "/decks/deck_test/sessions"
			if len(v.queryParams) > 0 {
				params := make([]string, 0)
				for key, value := range v.queryParams {
					params = append(params, fmt.Sprintf("%s=%s", key, value))
				}
				url = fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
				Plan: &malak.Plan{
					Metadata: malak.PlanMetadata{
						Deck: struct {
							AutoTerminateLink bool "json:\"auto_terminate_link,omitempty\""
							CustomDomain      bool "json:\"custom_domain,omitempty\""
							Analytics         struct {
								CanViewHistoricalSessions bool "json:\"can_view_historical_sessions,omitempty\""
							} "json:\"analytics,omitempty\""
						}{
							Analytics: struct {
								CanViewHistoricalSessions bool "json:\"can_view_historical_sessions,omitempty\""
							}{
								CanViewHistoricalSessions: true,
							},
						},
					},
				},
			}))

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.fetchDeckSessions,
				getConfig(), "decks.fetch_sessions").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDeckEngagementsTestTable() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
	}{
		{
			name:               "no reference provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching engagements",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().DeckEngagements(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("engagements error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully fetched engagements",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
				deck.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						Reference: "deck_test",
					}, nil)

				deck.EXPECT().DeckEngagements(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckEngagementResponse{
						DailyEngagements: []malak.DeckDailyEngagement{
							{
								EngagementCount: 100,
							},
						},
						GeographicStats: []malak.DeckGeographicStat{
							{
								ViewCount: 100,
							},
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDeckHandler_FetchEngagements(t *testing.T) {
	for _, v := range generateDeckEngagementsTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/decks/deck_test/analytics", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.fetchEngagements,
				getConfig(), "decks.fetch_engagements").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
