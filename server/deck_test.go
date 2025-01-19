package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
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
	mockFn             func(update *malak_mocks.MockDeckRepository)
	expectedStatusCode int
	req                createDeckRequest
} {

	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository)
		expectedStatusCode int
		req                createDeckRequest
	}{
		{
			name:               "no url provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckRequest{
				DeckURL: "",
			},
		},
		{
			name:               "no title provided",
			mockFn:             func(deck *malak_mocks.MockDeckRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckRequest{
				DeckURL: "https://google.com",
			},
		},
		{
			name: "could not create deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
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
			mockFn: func(deck *malak_mocks.MockDeckRepository) {
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

			v.mockFn(deckRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
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
