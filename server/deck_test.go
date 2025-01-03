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
