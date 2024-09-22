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

func generateUpdateDuplicationContent() []struct {
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
			name: "duplicating content failed",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "duplicating content succeeds",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
	}
}

func TestUpdatesHandler_Delete(t *testing.T) {
	for _, v := range generateUpdateDeletionTable() {

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

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_jfnkfjkf")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), u.delete, getConfig(), "updates.delete").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateDeletionTable() []struct {
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
			name: "deleting content fails",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "deleting content succeeds",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_Duplicate(t *testing.T) {
	for _, v := range generateUpdateDuplicationContent() {

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

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "update_jfnkfjkf")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), u.duplicate, getConfig(), "updates.duplicate").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
