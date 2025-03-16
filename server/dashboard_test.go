package server

import (
	"bytes"
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

func generateDashboardCreateRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
	expectedStatusCode int
	req                createDashboardRequest
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
		expectedStatusCode int
		req                createDashboardRequest
	}{
		{
			name:               "no title provided",
			mockFn:             func(dashboard *malak_mocks.MockDashboardRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDashboardRequest{
				Description: "Test description",
			},
		},
		{
			name:               "no description provided",
			mockFn:             func(dashboard *malak_mocks.MockDashboardRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDashboardRequest{
				Title: "Test Dashboard",
			},
		},
		{
			name:               "title too long",
			mockFn:             func(dashboard *malak_mocks.MockDashboardRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDashboardRequest{
				Title:       string(make([]byte, 101)), // 101 characters
				Description: "Test description",
			},
		},
		{
			name:               "description too long",
			mockFn:             func(dashboard *malak_mocks.MockDashboardRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			req: createDashboardRequest{
				Title:       "Test Dashboard",
				Description: string(make([]byte, 501)), // 501 characters
			},
		},
		{
			name: "error creating dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not create dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createDashboardRequest{
				Title:       "Test Dashboard",
				Description: "Test description",
			},
		},
		{
			name: "successfully created dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: createDashboardRequest{
				Title:       "Test Dashboard",
				Description: "Test description",
			},
		},
	}
}

func TestDashboardHandler_Create(t *testing.T) {
	for _, v := range generateDashboardCreateRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			v.mockFn(dashboardRepo)

			h := &dashboardHandler{
				dashboardRepo: dashboardRepo,
				generator:     &mockReferenceGenerator{},
				cfg:           getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t),
				h.create,
				getConfig(), "dashboards.create").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateDashboardListRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
		expectedStatusCode int
	}{
		{
			name: "error listing dashboards",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("could not list dashboards"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully listed dashboards",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Dashboard{
						{
							ID:          workspaceID,
							Title:       "Test Dashboard",
							Description: "Test description",
							Reference:   "DASH_123",
							ChartCount:  0,
							WorkspaceID: workspaceID,
						},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "empty dashboards list",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Dashboard{}, int64(0), nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDashboardHandler_List(t *testing.T) {
	for _, v := range generateDashboardListRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			v.mockFn(dashboardRepo)

			h := &dashboardHandler{
				dashboardRepo: dashboardRepo,
				generator:     &mockReferenceGenerator{},
				cfg:           getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			WrapMalakHTTPHandler(getLogger(t),
				h.list,
				getConfig(), "dashboards.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateListAllChartsRequest() []struct {
	name               string
	mockFn             func(integration *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(integration *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
	}{
		{
			name: "error listing charts",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().ListCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not list charts"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully listed charts",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().ListCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChart{
						{
							ID:                     workspaceID,
							WorkspaceIntegrationID: workspaceID,
							WorkspaceID:            workspaceID,
							Reference:              "CHART_123",
							UserFacingName:         "Test Chart",
							InternalName:           malak.IntegrationChartInternalNameTypeMercuryAccount,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "empty charts list",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().ListCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChart{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDashboardHandler_ListAllCharts(t *testing.T) {
	for _, v := range generateListAllChartsRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			v.mockFn(integrationRepo)

			h := &dashboardHandler{
				integrationRepo: integrationRepo,
				cfg:             getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			WrapMalakHTTPHandler(getLogger(t),
				h.listAllCharts,
				getConfig(), "dashboards.listAllCharts").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateAddChartRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
	req                addChartToDashboardRequest
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
		req                addChartToDashboardRequest
	}{
		{
			name: "no chart reference provided",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                addChartToDashboardRequest{},
		},
		{
			name: "dashboard not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "chart not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error fetching chart",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, errors.New("error fetching chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error adding chart to dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				dashboard.EXPECT().AddChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error adding chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "successfully added chart to dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				dashboard.EXPECT().AddChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
	}
}

func TestDashboardHandler_AddChart(t *testing.T) {
	for _, v := range generateAddChartRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			v.mockFn(dashboardRepo, integrationRepo)

			h := &dashboardHandler{
				dashboardRepo:   dashboardRepo,
				integrationRepo: integrationRepo,
				generator:       &mockReferenceGenerator{},
				cfg:             getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/dashboards/DASH_123/charts", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			router := chi.NewRouter()
			router.Put("/dashboards/{reference}/charts", WrapMalakHTTPHandler(getLogger(t),
				h.addChart,
				getConfig(), "dashboards.add_chart").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateFetchDashboardRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository)
		expectedStatusCode int
	}{
		{
			name: "dashboard not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching dashboard positions",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				dashboard.EXPECT().GetCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.DashboardChart{}, nil)

				dashboard.EXPECT().GetDashboardPositions(gomock.Any(), workspaceID).
					Times(1).
					Return(nil, errors.New("error fetching dashboard positions"))

				dashboardLink.EXPECT().DefaultLink(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.DashboardLink{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching dashboard charts",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				dashboard.EXPECT().GetCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("error fetching charts"))

				dashboard.EXPECT().GetDashboardPositions(gomock.Any(), workspaceID).
					Times(1).
					Return([]malak.DashboardChartPosition{}, nil)

				dashboardLink.EXPECT().DefaultLink(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.DashboardLink{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully fetched dashboard and charts",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  1,
						WorkspaceID: workspaceID,
					}, nil)

				dashboard.EXPECT().GetCharts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.DashboardChart{
						{
							ID:                     workspaceID,
							Reference:              "DASHCHART_123",
							WorkspaceIntegrationID: workspaceID,
							WorkspaceID:            workspaceID,
							DashboardID:            workspaceID,
							ChartID:                workspaceID,
						},
					}, nil)

				dashboard.EXPECT().GetDashboardPositions(gomock.Any(), workspaceID).
					Times(1).
					Return([]malak.DashboardChartPosition{
						{
							ID:          workspaceID,
							DashboardID: workspaceID,
							ChartID:     workspaceID,
							OrderIndex:  1,
						},
					}, nil)

				dashboardLink.EXPECT().DefaultLink(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.DashboardLink{
						ID:          workspaceID,
						DashboardID: workspaceID,
						Reference:   "DASHLINK_123",
						Token:       "test_token",
						LinkType:    malak.DashboardLinkTypeDefault,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDashboardHandler_FetchDashboard(t *testing.T) {
	for _, v := range generateFetchDashboardRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			dashboardLinkRepo := malak_mocks.NewMockDashboardLinkRepository(controller)
			v.mockFn(dashboardRepo, dashboardLinkRepo)

			h := &dashboardHandler{
				dashboardRepo:     dashboardRepo,
				dashboardLinkRepo: dashboardLinkRepo,
				generator:         &mockReferenceGenerator{},
				cfg:               getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/dashboards/DASH_123", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			router := chi.NewRouter()
			router.Get("/dashboards/{reference}", WrapMalakHTTPHandler(getLogger(t),
				h.fetchDashboard,
				getConfig(), "dashboards.fetch").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateFetchChartingDataRequest() []struct {
	name               string
	mockFn             func(integration *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(integration *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
	}{
		{
			name: "chart not found",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching chart",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, errors.New("error fetching chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching data points",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				integration.EXPECT().GetDataPoints(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("error fetching data points"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully fetched chart data",
			mockFn: func(integration *malak_mocks.MockIntegrationRepository) {
				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				integration.EXPECT().GetDataPoints(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationDataPoint{
						{
							ID:                     workspaceID,
							WorkspaceIntegrationID: workspaceID,
							WorkspaceID:            workspaceID,
							IntegrationChartID:     workspaceID,
							Reference:              "datapoint_123",
							PointName:              "Test Point",
							PointValue:             100,
							DataPointType:          malak.IntegrationDataPointTypeCurrency,
							Metadata:               malak.IntegrationDataPointMetadata{},
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDashboardHandler_FetchChartingData(t *testing.T) {
	for _, v := range generateFetchChartingDataRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			v.mockFn(integrationRepo)

			h := &dashboardHandler{
				integrationRepo: integrationRepo,
				cfg:             getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/dashboards/charts/CHART_123", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			router := chi.NewRouter()
			router.Get("/dashboards/charts/{reference}", WrapMalakHTTPHandler(getLogger(t),
				h.fetchChartingData,
				getConfig(), "dashboards.fetch_charting_data").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateUpdateDashboardPositionsRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
	expectedStatusCode int
	req                updateDashboardPositionsRequest
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository)
		expectedStatusCode int
		req                updateDashboardPositionsRequest
	}{
		{
			name: "dashboard not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: updateDashboardPositionsRequest{
				Positions: []struct {
					ChartID uuid.UUID `json:"chart_id,omitempty" validate:"required"`
					Index   int64     `json:"index,omitempty" validate:"required"`
				}{
					{
						ChartID: workspaceID,
						Index:   1,
					},
				},
			},
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDashboardPositionsRequest{
				Positions: []struct {
					ChartID uuid.UUID `json:"chart_id,omitempty" validate:"required"`
					Index   int64     `json:"index,omitempty" validate:"required"`
				}{
					{
						ChartID: workspaceID,
						Index:   1,
					},
				},
			},
		},
		{
			name: "error updating positions",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  1,
						WorkspaceID: workspaceID,
					}, nil)

				dashboard.EXPECT().UpdateDashboardPositions(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error updating positions"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateDashboardPositionsRequest{
				Positions: []struct {
					ChartID uuid.UUID `json:"chart_id,omitempty" validate:"required"`
					Index   int64     `json:"index,omitempty" validate:"required"`
				}{
					{
						ChartID: workspaceID,
						Index:   1,
					},
				},
			},
		},
		{
			name: "successfully updated positions",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  1,
						WorkspaceID: workspaceID,
					}, nil)

				dashboard.EXPECT().UpdateDashboardPositions(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updateDashboardPositionsRequest{
				Positions: []struct {
					ChartID uuid.UUID `json:"chart_id,omitempty" validate:"required"`
					Index   int64     `json:"index,omitempty" validate:"required"`
				}{
					{
						ChartID: workspaceID,
						Index:   1,
					},
				},
			},
		},
	}
}

func TestDashboardHandler_UpdateDashboardPositions(t *testing.T) {
	for _, v := range generateUpdateDashboardPositionsRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			v.mockFn(dashboardRepo)

			h := &dashboardHandler{
				dashboardRepo: dashboardRepo,
				cfg:           getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/dashboards/DASH_123/positions", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			router := chi.NewRouter()
			router.Post("/dashboards/{reference}/positions", WrapMalakHTTPHandler(getLogger(t),
				h.updateDashboardPositions,
				getConfig(), "dashboards.update_positions").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateRemoveChartRequest() []struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
	req                addChartToDashboardRequest
} {
	return []struct {
		name               string
		mockFn             func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
		req                addChartToDashboardRequest
	}{
		{
			name: "no chart reference provided",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                addChartToDashboardRequest{},
		},
		{
			name: "dashboard not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "chart not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error fetching chart",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, errors.New("error fetching chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "error removing chart from dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				dashboard.EXPECT().RemoveChart(gomock.Any(), workspaceID, workspaceID).
					Times(1).
					Return(errors.New("error removing chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
		{
			name: "successfully removed chart from dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						Title:       "Test Dashboard",
						Description: "Test description",
						Reference:   "DASH_123",
						ChartCount:  0,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{
						ID:                     workspaceID,
						WorkspaceIntegrationID: workspaceID,
						Reference:              "CHART_123",
						WorkspaceID:            workspaceID,
					}, nil)

				dashboard.EXPECT().RemoveChart(gomock.Any(), workspaceID, workspaceID).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: addChartToDashboardRequest{
				ChartReference: "CHART_123",
			},
		},
	}
}

func TestDashboardHandler_RemoveChart(t *testing.T) {
	for _, v := range generateRemoveChartRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			v.mockFn(dashboardRepo, integrationRepo)

			h := &dashboardHandler{
				dashboardRepo:   dashboardRepo,
				integrationRepo: integrationRepo,
				generator:       &mockReferenceGenerator{},
				cfg:             getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/dashboards/DASH_123/charts", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: workspaceID}))

			router := chi.NewRouter()
			router.Delete("/dashboards/{reference}/charts", WrapMalakHTTPHandler(getLogger(t),
				h.removeChart,
				getConfig(), "dashboards.remove_chart").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
