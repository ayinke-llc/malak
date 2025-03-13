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

func generatePublicDeckDetailsTestTable() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService)
	expectedStatusCode int
	req                createDeckViewerSession
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService)
		expectedStatusCode int
		req                createDeckViewerSession
	}{
		{
			name: "no reference provided",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                createDeckViewerSession{},
		},
		{
			name: "no os provided",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckViewerSession{
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "no device info provided",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createDeckViewerSession{
				OS:      "iOS",
				Browser: "Safari",
			},
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "error getting object path",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey: "test-key",
					}, nil)

				gulter.file = nil // This will cause Path to return an error

				geo.EXPECT().FindByIP(gomock.Any(), gomock.Any()).
					Times(1).
					Return("US", "New York", nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "error getting geolocation",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey: "test-key",
					}, nil)

				gulter.file = &uploadedFileMetadata{
					Size:              1024,
					FolderDestination: "decks",
					Key:               "test-key",
				}

				geo.EXPECT().FindByIP(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", "", errors.New("geolocation error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "error creating session",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey: "test-key",
					}, nil)

				gulter.file = &uploadedFileMetadata{
					Size:              1024,
					FolderDestination: "decks",
					Key:               "test-key",
				}

				geo.EXPECT().FindByIP(gomock.Any(), gomock.Any()).
					Times(1).
					Return("US", "New York", nil)

				deck.EXPECT().CreateDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("session error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
		{
			name: "successfully fetched deck details",
			mockFn: func(deck *malak_mocks.MockDeckRepository, gulter *mockStorage, geo *malak_mocks.MockGeolocationService) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey: "test-key",
						DeckPreference: &malak.DeckPreference{
							EnableDownloading: true,
							RequireEmail:      true,
						},
					}, nil)

				gulter.file = &uploadedFileMetadata{
					Size:              1024,
					FolderDestination: "decks",
					Key:               "test-key",
				}

				geo.EXPECT().FindByIP(gomock.Any(), gomock.Any()).
					Times(1).
					Return("US", "New York", nil)

				deck.EXPECT().CreateDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: createDeckViewerSession{
				OS:         "iOS",
				DeviceInfo: "iPhone",
				Browser:    "Safari",
			},
		},
	}
}

func TestDeckHandler_PublicDeckDetails(t *testing.T) {
	for _, v := range generatePublicDeckDetailsTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			gulterStore := &mockStorage{}
			geoService := malak_mocks.NewMockGeolocationService(controller)

			v.mockFn(deckRepo, gulterStore, geoService)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
				gulterStore:        gulterStore,
				cfg:                getConfig(),
				geolocationService: geoService,
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/public/decks/deck_test", b)
			req.Header.Add("Content-Type", "application/json")

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.publicDeckDetails,
				getConfig(), "decks.public_details").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateDeckViewerSessionTestTable() []struct {
	name               string
	mockFn             func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository)
	expectedStatusCode int
	req                updateDeckViewerSession
} {
	return []struct {
		name               string
		mockFn             func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository)
		expectedStatusCode int
		req                updateDeckViewerSession
	}{
		{
			name: "no reference provided",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                updateDeckViewerSession{},
		},
		{
			name: "invalid email",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: updateDeckViewerSession{
				Email: "invalid-email",
			},
		},
		{
			name: "deck not found",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrDeckNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: updateDeckViewerSession{
				Email:     "test@example.com",
				TimeSpent: 60,
				SessionID: "session_123",
			},
		},
		{
			name: "error fetching deck",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDeckViewerSession{
				Email:     "test@example.com",
				TimeSpent: 60,
				SessionID: "session_123",
			},
		},
		{
			name: "error finding session",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey:   "test-key",
						WorkspaceID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					}, nil)

				contact.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)

				deck.EXPECT().FindDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("session not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDeckViewerSession{
				Email:     "test@example.com",
				TimeSpent: 60,
				SessionID: "session_123",
			},
		},
		{
			name: "error updating session",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey:   "test-key",
						WorkspaceID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					}, nil)

				contact.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)

				deck.EXPECT().FindDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckViewerSession{}, nil)

				deck.EXPECT().UpdateDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("update error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDeckViewerSession{
				Email:     "test@example.com",
				TimeSpent: 60,
				SessionID: "session_123",
			},
		},
		{
			name: "successfully updated session",
			mockFn: func(deck *malak_mocks.MockDeckRepository, contact *malak_mocks.MockContactRepository) {
				deck.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Deck{
						ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						ObjectKey:   "test-key",
						WorkspaceID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					}, nil)

				contact.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)

				deck.EXPECT().FindDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.DeckViewerSession{}, nil)

				deck.EXPECT().UpdateDeckSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updateDeckViewerSession{
				Email:     "test@example.com",
				TimeSpent: 60,
				SessionID: "session_123",
			},
		},
	}
}

func TestDeckHandler_UpdateDeckViewerSession(t *testing.T) {
	for _, v := range generateUpdateDeckViewerSessionTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			deckRepo := malak_mocks.NewMockDeckRepository(controller)
			contactRepo := malak_mocks.NewMockContactRepository(controller)

			v.mockFn(deckRepo, contactRepo)

			u := &deckHandler{
				referenceGenerator: &mockReferenceGenerator{},
				deckRepo:           deckRepo,
				contactRepo:        contactRepo,
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/public/decks/deck_test", b)
			req.Header.Add("Content-Type", "application/json")

			ctx := chi.NewRouteContext()
			if v.name != "no reference provided" {
				ctx.URLParams.Add("reference", "deck_test")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.updateDeckViewerSession,
				getConfig(), "decks.update_session").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
