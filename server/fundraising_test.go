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

func generateNewPipelineTestTable() []struct {
	name               string
	mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository)
	req                createNewPipelineRequest
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository)
		req                createNewPipelineRequest
		expectedStatusCode int
	}{
		{
			name: "valid request - seed stage",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, pipeline *malak.FundraisingPipeline, columns ...malak.FundraisingPipelineColumn) error {
						require.Len(t, columns, len(malak.DefaultFundraisingColumns))
						for i, col := range columns {
							require.Equal(t, malak.DefaultFundraisingColumns[i].Title, col.Title)
							require.Equal(t, malak.DefaultFundraisingColumns[i].ColumnType, col.ColumnType)
							require.Equal(t, malak.DefaultFundraisingColumns[i].Description, col.Description)
						}
						return nil
					})
			},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "valid request - series A stage",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, pipeline *malak.FundraisingPipeline, columns ...malak.FundraisingPipelineColumn) error {
						require.Len(t, columns, len(malak.DefaultFundraisingColumns))
						for i, col := range columns {
							require.Equal(t, malak.DefaultFundraisingColumns[i].Title, col.Title)
							require.Equal(t, malak.DefaultFundraisingColumns[i].ColumnType, col.ColumnType)
							require.Equal(t, malak.DefaultFundraisingColumns[i].Description, col.Description)
						}
						return nil
					})
			},
			req: createNewPipelineRequest{
				Title:             "Series A Fundraising",
				Stage:             malak.FundraisePipelineStageSeriesA,
				Amount:            5000000,
				Description:       "Series A round for expansion.",
				ExpectedCloseDate: time.Now().Add(90 * 24 * time.Hour).Unix(), // 90 days in future
				StartDate:         time.Now().Add(24 * time.Hour).Unix(),      // tomorrow
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "invalid request - empty title",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - title too short",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "abc",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - description too long",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       string(make([]byte, 201)), // 201 characters
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - start date in past",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Add(-24 * time.Hour).Unix(), // yesterday
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - expected close date today",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - expected close date in past",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(-24 * time.Hour).Unix(), // yesterday
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "invalid request - invalid stage",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             "invalid_stage",
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "repository error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("repository error"))
			},
			req: createNewPipelineRequest{
				Title:             "Valid Title",
				Stage:             malak.FundraisePipelineStageSeed,
				Amount:            1000,
				Description:       "A valid description.",
				ExpectedCloseDate: time.Now().Add(24 * time.Hour).Unix(),
				StartDate:         time.Now().Unix(),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
}

func TestFundraisingHandler_NewPipeline(t *testing.T) {
	for _, v := range generateNewPipelineTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)
			v.mockFn(t, fundingRepo)

			handler := &fundraisingHandler{
				cfg:                getConfig(),
				fundingRepo:        fundingRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/pipelines", b)
			req.Header.Add("Content-Type", "application/json")

			workspace := &malak.Workspace{
				ID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
			}
			user := &malak.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440008"),
			}

			// Set up context in the correct order
			ctx := req.Context()
			ctx = writeUserToCtx(ctx, user)
			ctx = writeWorkspaceToCtx(ctx, workspace)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, chi.NewRouteContext())
			req = req.WithContext(ctx)

			WrapMalakHTTPHandler(getLogger(t),
				handler.newPipeline,
				getConfig(),
				"fundraising.new_pipeline").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

type listPipelinesTestCase struct {
	name               string
	mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository)
	queryParams        map[string]string
	expectedStatusCode int
}

func generateListPipelinesTestTable() []listPipelinesTestCase {
	return []listPipelinesTestCase{
		{
			name: "successful listing - default params",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, opts malak.ListPipelineOptions) ([]*malak.FundraisingPipeline, int64, error) {
						require.Equal(t, int64(1), opts.Paginator.Page)
						require.Equal(t, int64(8), opts.Paginator.PerPage)
						require.False(t, opts.ActiveOnly)
						return []*malak.FundraisingPipeline{}, int64(0), nil
					})
			},
			queryParams:        map[string]string{},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "successful listing - with pagination",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, opts malak.ListPipelineOptions) ([]*malak.FundraisingPipeline, int64, error) {
						require.Equal(t, int64(2), opts.Paginator.Page)
						require.Equal(t, int64(20), opts.Paginator.PerPage)
						require.False(t, opts.ActiveOnly)
						return []*malak.FundraisingPipeline{}, int64(0), nil
					})
			},
			queryParams: map[string]string{
				"page":     "2",
				"per_page": "20",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "successful listing - active only",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, opts malak.ListPipelineOptions) ([]*malak.FundraisingPipeline, int64, error) {
						require.True(t, opts.ActiveOnly)
						return []*malak.FundraisingPipeline{}, int64(0), nil
					})
			},
			queryParams: map[string]string{
				"active_only": "true",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "repository error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("repository error"))
			},
			queryParams:        map[string]string{},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
}

func TestFundraisingHandler_List(t *testing.T) {
	for _, v := range generateListPipelinesTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)
			v.mockFn(t, fundingRepo)

			handler := &fundraisingHandler{
				cfg:                getConfig(),
				fundingRepo:        fundingRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/pipelines", nil)

			q := req.URL.Query()
			for key, value := range v.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			workspace := &malak.Workspace{
				ID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
			}
			user := &malak.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"),
			}

			ctx := req.Context()
			ctx = writeUserToCtx(ctx, user)
			ctx = writeWorkspaceToCtx(ctx, workspace)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, chi.NewRouteContext())
			req = req.WithContext(ctx)

			WrapMalakHTTPHandler(getLogger(t),
				handler.list,
				getConfig(),
				"fundraising.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

type boardTestCase struct {
	name               string
	mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository)
	expectedStatusCode int
}

func generateBoardTestTable() []boardTestCase {
	return []boardTestCase{
		{
			name: "successful fetch",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:          uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb9a"),
					Reference:   "pipeline_123",
					Title:       "Test Pipeline",
					WorkspaceID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
				}
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				repo.EXPECT().
					Board(gomock.Any(), pipeline).
					Times(1).
					Return([]malak.FundraisingPipelineColumn{
						{
							Title:      "Column 1",
							ColumnType: malak.FundraisePipelineColumnTypeNormal,
						},
					}, []malak.FundraiseContact{
						{
							Reference: "contact_1",
						},
					}, []malak.FundraiseContactPosition{
						{
							Reference:  "position_1",
							OrderIndex: 1,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "missing reference parameter",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				// No mock expectations since it should fail before repository calls
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "pipeline not found",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrPipelineNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "get pipeline error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "board fetch error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:          uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb9c"),
					Reference:   "pipeline_123",
					Title:       "Test Pipeline",
					WorkspaceID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
				}
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				repo.EXPECT().
					Board(gomock.Any(), pipeline).
					Times(1).
					Return(nil, nil, nil, errors.New("repository error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
}

func TestFundraisingHandler_Board(t *testing.T) {
	for _, v := range generateBoardTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)
			v.mockFn(t, fundingRepo)

			handler := &fundraisingHandler{
				cfg:                getConfig(),
				fundingRepo:        fundingRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			rr := httptest.NewRecorder()
			var reference string
			switch v.name {
			case "missing reference parameter":
				reference = ""
			case "pipeline not found":
				reference = "non_existent_pipeline"
			default:
				reference = "pipeline_123"
			}

			path := "/pipelines/" + reference + "/board"
			if reference == "" {
				path = "/pipelines//board" // Test missing reference case
			}
			req := httptest.NewRequest(http.MethodGet, path, nil)

			workspace := &malak.Workspace{
				ID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
			}
			user := &malak.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			}

			// Set up context in the correct order
			ctx := req.Context()
			ctx = writeUserToCtx(ctx, user)
			ctx = writeWorkspaceToCtx(ctx, workspace)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, chi.NewRouteContext())
			req = req.WithContext(ctx)

			// Set up route params
			rctx := chi.NewRouteContext()
			if reference != "" {
				rctx.URLParams.Add("reference", reference)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			WrapMalakHTTPHandler(getLogger(t),
				handler.board,
				getConfig(),
				"fundraising.board").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

type closeBoardTestCase struct {
	name               string
	mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository)
	expectedStatusCode int
}

func generateCloseBoardTestTable() []closeBoardTestCase {
	return []closeBoardTestCase{
		{
			name: "successful close",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:          uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb9b"),
					Reference:   "pipeline_123",
					Title:       "Test Pipeline",
					WorkspaceID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
					IsClosed:    false,
				}
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				repo.EXPECT().
					CloseBoard(gomock.Any(), pipeline).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "missing reference parameter",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				// No mock expectations since it should fail before repository calls
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "pipeline not found",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrPipelineNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "pipeline already closed",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:          uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb9d"),
					Reference:   "pipeline_123",
					Title:       "Test Pipeline",
					WorkspaceID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
					IsClosed:    true,
				}
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "get pipeline error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "close board error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:          uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb9f"),
					Reference:   "pipeline_123",
					Title:       "Test Pipeline",
					WorkspaceID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
					IsClosed:    false,
				}
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				repo.EXPECT().
					CloseBoard(gomock.Any(), pipeline).
					Times(1).
					Return(errors.New("repository error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
}

func TestFundraisingHandler_CloseBoard(t *testing.T) {
	for _, v := range generateCloseBoardTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)
			v.mockFn(t, fundingRepo)

			handler := &fundraisingHandler{
				cfg:                getConfig(),
				fundingRepo:        fundingRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			rr := httptest.NewRecorder()
			var reference string
			switch v.name {
			case "missing reference parameter":
				reference = ""
			case "pipeline not found":
				reference = "non_existent_pipeline"
			default:
				reference = "pipeline_123"
			}

			path := "/pipelines/" + reference + "/close"
			if reference == "" {
				path = "/pipelines//close" // Test missing reference case
			}
			req := httptest.NewRequest(http.MethodPost, path, nil)

			workspace := &malak.Workspace{
				ID: uuid.MustParse("1e66cedd-0e53-493a-adfd-7f81221c2248"),
			}
			user := &malak.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
			}

			// Set up context in the correct order
			ctx := req.Context()
			ctx = writeUserToCtx(ctx, user)
			ctx = writeWorkspaceToCtx(ctx, workspace)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, chi.NewRouteContext())
			req = req.WithContext(ctx)

			// Set up route params
			rctx := chi.NewRouteContext()
			if reference != "" {
				rctx.URLParams.Add("reference", reference)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			WrapMalakHTTPHandler(getLogger(t),
				handler.closeBoard,
				getConfig(),
				"fundraising.close_board").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

type addContactTestCase struct {
	name               string
	mockFn             func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository)
	req                addContactRequest
	pipelineReference  string
	expectedStatusCode int
}

func generateAddContactTestTable() []addContactTestCase {
	now := time.Now().UTC()
	return []addContactTestCase{
		{
			name: "successful add contact to board",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Title:    "Test Pipeline",
					Stage:    malak.FundraisePipelineStageSeed,
					IsClosed: false,
				}

				contact := &malak.Contact{
					ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
					FirstName: "Test",
					LastName:  "Contact",
				}

				defaultColumn := malak.FundraisingPipelineColumn{
					Title:       "Initial Contact",
					ColumnType:  malak.FundraisePipelineColumnTypeNormal,
					Description: "Initial contact column",
				}

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(contact, nil)

				repo.EXPECT().
					DefaultColumn(gomock.Any(), pipeline).
					Times(1).
					Return(defaultColumn, nil)

				repo.EXPECT().
					AddContactToBoard(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100, // $1M in cents
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "invalid request - empty contact reference",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid request - empty pipeline reference",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid request - rating too high",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           6,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid request - rating too low",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           -1,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid request - check size too small",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        500 * 100, // $500 in cents
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid request - initial contact date in future",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "pipeline not found",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrPipelineNotFound)
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "contact not found",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:       uuid.MustParse("34bc303d-6ca6-4881-a31f-55503b98eb90"),
					Title:    "Test Pipeline",
					Stage:    malak.FundraisePipelineStageSeed,
					IsClosed: false,
				}

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "default column error",
			mockFn: func(t *testing.T, repo *malak_mocks.MockFundraisingPipelineRepository, contactRepo *malak_mocks.MockContactRepository) {
				pipeline := &malak.FundraisingPipeline{
					ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
					Title:    "Test Pipeline",
					Stage:    malak.FundraisePipelineStageSeed,
					IsClosed: false,
				}

				contact := &malak.Contact{
					ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
					FirstName: "Test",
					LastName:  "Contact",
				}

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(pipeline, nil)

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(contact, nil)

				repo.EXPECT().
					DefaultColumn(gomock.Any(), pipeline).
					Times(1).
					Return(malak.FundraisingPipelineColumn{}, errors.New("default column error"))
			},
			req: addContactRequest{
				ContactReference: "contact_123",
				Rating:           4,
				CanLeadRound:     true,
				InitialContact:   now.Add(-24 * time.Hour).Unix(),
				CheckSize:        1000000 * 100,
			},
			pipelineReference:  "pipeline_123",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
}

func TestFundraisingHandler_AddContact(t *testing.T) {
	workspaceID := uuid.MustParse("56670b6d-48d4-4b17-bc8f-d101b7d0b53c")
	for _, v := range generateAddContactTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			fundingRepo := malak_mocks.NewMockFundraisingPipelineRepository(controller)
			contactRepo := malak_mocks.NewMockContactRepository(controller)
			v.mockFn(t, fundingRepo, contactRepo)

			handler := &fundraisingHandler{
				cfg:                getConfig(),
				fundingRepo:        fundingRepo,
				contactRepo:        contactRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/pipelines/"+v.pipelineReference+"/contacts", b)
			req.Header.Add("Content-Type", "application/json")

			workspace := &malak.Workspace{
				ID: workspaceID,
			}
			user := &malak.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			}

			// Set up context in the correct order
			ctx := req.Context()
			ctx = writeUserToCtx(ctx, user)
			ctx = writeWorkspaceToCtx(ctx, workspace)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("reference", v.pipelineReference)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
			req = req.WithContext(ctx)

			WrapMalakHTTPHandler(getLogger(t),
				handler.addContact,
				getConfig(),
				"fundraising.add_contact").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
