package server

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateFetchUpdateAnalaytics() []struct {
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
			name: "update not exists",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrUpdateNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "could not fetch update from db",
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
			name: "could not fetch recipient stats",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().
					RecipientStat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could nto fetch recipient stats"))

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "could not fetch update stats",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().
					RecipientStat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could nto fetch update stats"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fetched successfully",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().
					RecipientStat(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.UpdateRecipient{}, nil)

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.UpdateStat{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdateHandler_FetchUpdateAnalytics(t *testing.T) {

	for _, v := range generateFetchUpdateAnalaytics() {

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

			req := httptest.NewRequest(http.MethodGet, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_jfnkfjkf")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				u.fetchUpdateAnalytics,
				getConfig(), "updates.fetchUpdateAnalytics").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
