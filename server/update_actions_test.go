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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateFetchUpdateTable() []struct {
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
			name: "Fetched item successfully",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdateHandler_FetchUpdate(t *testing.T) {
	for _, v := range generateFetchUpdateTable() {

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

			WrapMalakHTTPHandler(getLogger(t), u.fetchUpdate, getConfig(), "updates.fetchUpdate").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

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

				update.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
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

				update.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
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
			name: "you cannot delete an update that is sent already",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{
						Status: malak.UpdateStatusSent,
					}, nil)

				update.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			expectedStatusCode: http.StatusBadRequest,
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

func generateUpdatePin() []struct {
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
			name: "toggling failed",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{}, nil)

				update.EXPECT().
					TogglePinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "toggling to false succeeds",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{
						IsPinned: true,
					}, nil)

				update.EXPECT().
					TogglePinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "toggling to true succeeds",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Update{
						IsPinned: false,
					}, nil)

				update.EXPECT().TogglePinned(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_Pin(t *testing.T) {
	for _, v := range generateUpdatePin() {

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

			WrapMalakHTTPHandler(getLogger(t), u.togglePinned, getConfig(), "updates.toggle").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateReaction() []struct {
	name               string
	mockFn             func(update *malak_mocks.MockUpdateRepository)
	expectedStatusCode int
} {

	id := uuid.MustParse("ed3daf26-5839-4b91-bf19-64c8a6b3bb7a")

	return []struct {
		name               string
		mockFn             func(update *malak_mocks.MockUpdateRepository)
		expectedStatusCode int
	}{
		{
			name: "recipient stat does not exists",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					GetStatByEmailID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil, errors.New("oops"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update stat failed",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					GetStatByEmailID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &malak.UpdateRecipientStat{
						Recipient: &malak.UpdateRecipient{
							UpdateID: id,
						},
					}, nil)

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update stat updating failed",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					GetStatByEmailID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &malak.UpdateRecipientStat{
						Recipient: &malak.UpdateRecipient{
							UpdateID: id,
						},
					}, nil)

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Return(&malak.UpdateStat{}, nil)

				update.EXPECT().
					UpdateStat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("could not update"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update stat reaction",
			mockFn: func(update *malak_mocks.MockUpdateRepository) {
				update.
					EXPECT().
					GetStatByEmailID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &malak.UpdateRecipientStat{
						Recipient: &malak.UpdateRecipient{
							UpdateID: id,
						},
					}, nil)

				update.EXPECT().
					Stat(gomock.Any(), gomock.Any()).
					Return(&malak.UpdateStat{}, nil)

				update.EXPECT().
					UpdateStat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestUpdatesHandler_HandleReaction(t *testing.T) {
	for _, v := range generateUpdateReaction() {

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

			u.handleReaction(getLogger(t)).ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
