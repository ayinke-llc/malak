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

func TestWorkspaceHandler_CreateChart(t *testing.T) {
	for _, v := range generateCreateChartTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)

			v.mockFn(integrationRepo)

			a := &workspaceHandler{
				cfg:             getConfig(),
				queueClient:     queueRepo,
				integrationRepo: integrationRepo,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "chart_tt7-YieIgz"
				},
				referenceGenerator: &mockReferenceGenerator{},
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "integration_ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.createChart, getConfig(), "workspaces.charts.create").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_AddDataPoint(t *testing.T) {
	for _, v := range generateAddDataPointTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)

			v.mockFn(integrationRepo)

			a := &workspaceHandler{
				cfg:             getConfig(),
				queueClient:     queueRepo,
				integrationRepo: integrationRepo,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "datapoint_tt7-YieIgz"
				},
				referenceGenerator: &mockReferenceGenerator{},
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "integration_ref")
			ctx.URLParams.Add("chart_reference", "chart_ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.addDataPoint, getConfig(), "workspaces.charts.datapoint.add").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateCreateChartTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
	req                createChartRequest
} {
	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
		req                createChartRequest
	}{
		{
			name: "invalid title - empty",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "invalid title - too short",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "abc",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "invalid title - too long",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "this is a very long title that exceeds the maximum allowed length of one hundred characters which is not allowed in our system",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "invalid chart type",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: "invalid_type",
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "integration fetch error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "integration not enabled",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "non-system integration type",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "create chart error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
		{
			name: "successful chart creation",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createChartRequest{
				Title:     "Valid Title",
				ChartType: malak.IntegrationChartTypeBar,
				Datapoint: malak.IntegrationDataPointTypeCurrency,
			},
		},
	}
}

func generateAddDataPointTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
	req                addDataPointRequest
} {
	chartID := uuid.MustParse("8ce0f580-4d6d-429e-9d0e-a78eb99f62c2")
	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
		req                addDataPointRequest
	}{
		{
			name: "invalid value - negative",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: addDataPointRequest{
				Value: -1,
			},
		},
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "integration fetch error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "integration not enabled",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "non-system integration type",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "chart not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "chart fetch error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "add data point error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:             chartID,
						InternalName:   "test_chart",
						UserFacingName: "Test Chart",
					}, nil)

				integrationRepo.EXPECT().
					AddDataPoint(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addDataPointRequest{
				Value: 100,
			},
		},
		{
			name: "successful data point addition",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						IsActive:  true,
						Integration: &malak.Integration{
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)

				integrationRepo.EXPECT().
					GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:             chartID,
						InternalName:   "test_chart",
						UserFacingName: "Test Chart",
					}, nil)

				integrationRepo.EXPECT().
					AddDataPoint(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: addDataPointRequest{
				Value: 100,
			},
		},
	}
}
