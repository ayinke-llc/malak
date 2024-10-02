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
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
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
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
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
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
					},
				},
			},
		},
		{
			name: "could not fetch update from db",
			req: contentUpdateRequest{
				Title: "Valid title",
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
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
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
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
				Update: []malak.BlockContent{
					{
						ID:   "here is an id",
						Type: "heading",
						Content: []malak.BlockNoteItem{
							malak.BlockNoteItem{
								Type: "paragraph",
								Text: "omo",
							},
						},
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
					Return(nil, errors.New("could not list update"))

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
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
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
