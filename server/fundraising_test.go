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
				ID: uuid.New(),
			}
			user := &malak.User{
				ID: uuid.New(),
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
