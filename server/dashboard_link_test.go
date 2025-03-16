package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type dashboardLinkTestCase struct {
	name               string
	mockFn             func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler)
	expectedStatusCode int
	req                dashboardLinkRequest
}

type dashboardLinkRequest struct {
	GenericRequest
	Email malak.Email `json:"email,omitempty" validate:"optional"`
}

func generateDashboardLinkTestCases() []dashboardLinkTestCase {
	return []dashboardLinkTestCase{
		{
			name: "invalid email address",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: dashboardLinkRequest{
				Email: malak.Email("invalid-email"),
			},
		},
		{
			name: "dashboard not found",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: dashboardLinkRequest{
				Email: malak.Email("test@example.com"),
			},
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler) {
				dashboard.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: dashboardLinkRequest{
				Email: malak.Email("test@example.com"),
			},
		},
		{
			name: "error creating dashboard link",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler) {
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

				dashboardLink.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("error creating dashboard link"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: dashboardLinkRequest{
				Email: malak.Email("test@example.com"),
			},
		},
		{
			name: "successfully generated link with email",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queueHandler *malak_mocks.MockQueueHandler) {
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

				dashboardLink.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				queueHandler.EXPECT().Add(gomock.Any(), queue.QueueTopicShareDashboard, gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: dashboardLinkRequest{
				Email: malak.Email("test@example.com"),
			},
		},
		{
			name: "successfully generated default link",
			mockFn: func(dashboard *malak_mocks.MockDashboardRepository, dashboardLink *malak_mocks.MockDashboardLinkRepository, queue *malak_mocks.MockQueueHandler) {
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

				dashboardLink.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req:                dashboardLinkRequest{},
		},
	}
}

func TestDashboardHandler_GenerateLink(t *testing.T) {
	for _, v := range generateDashboardLinkTestCases() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			dashboardLinkRepo := malak_mocks.NewMockDashboardLinkRepository(controller)
			queueMock := malak_mocks.NewMockQueueHandler(controller)
			v.mockFn(dashboardRepo, dashboardLinkRepo, queueMock)

			h := &dashboardHandler{
				dashboardRepo:     dashboardRepo,
				dashboardLinkRepo: dashboardLinkRepo,
				generator:         &mockReferenceGenerator{},
				queue:             queueMock,
				cfg:               getConfig(),
			}

			var b = bytes.NewBuffer(nil)
			err := json.NewEncoder(b).Encode(v.req)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/dashboards/DASH_123/access-control/link", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
				ID:            workspaceID,
				WorkspaceName: "Test Workspace",
			}))

			router := chi.NewRouter()
			router.Post("/dashboards/{reference}/access-control/link", WrapMalakHTTPHandler(getLogger(t),
				h.generateLink,
				getConfig(), "dashboards.generate_link").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generatePublicDashboardDetailsRequest() []struct {
	name               string
	mockFn             func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository)
		expectedStatusCode int
	}{
		{
			name: "dashboard not found",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching charts",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
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
					Return(nil, errors.New("error fetching charts"))

				dashboard.EXPECT().GetDashboardPositions(gomock.Any(), workspaceID).
					Times(1).
					Return([]malak.DashboardChartPosition{}, nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching positions",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
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
					Return([]malak.DashboardChart{}, nil)

				dashboard.EXPECT().GetDashboardPositions(gomock.Any(), workspaceID).
					Times(1).
					Return(nil, errors.New("error fetching positions"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully fetched public dashboard details",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, dashboard *malak_mocks.MockDashboardRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
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
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestDashboardHandler_PublicDashboardDetails(t *testing.T) {
	for _, v := range generatePublicDashboardDetailsRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardLinkRepo := malak_mocks.NewMockDashboardLinkRepository(controller)
			dashboardRepo := malak_mocks.NewMockDashboardRepository(controller)
			v.mockFn(dashboardLinkRepo, dashboardRepo)

			h := &dashboardHandler{
				dashboardRepo:     dashboardRepo,
				dashboardLinkRepo: dashboardLinkRepo,
				cfg:               getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/public/dashboards/DASH_123", nil)
			req.Header.Add("Content-Type", "application/json")

			router := chi.NewRouter()
			router.Get("/public/dashboards/{reference}", WrapMalakHTTPHandler(getLogger(t),
				h.publicDashboardDetails,
				getConfig(), "dashboards.public_details").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generatePublicChartingDataRequest() []struct {
	name               string
	mockFn             func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
	}{
		{
			name: "dashboard not found",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, malak.ErrDashboardNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching dashboard",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{}, errors.New("error fetching dashboard"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "chart not found",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, malak.ErrChartNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "error fetching chart",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
						WorkspaceID: workspaceID,
					}, nil)

				integration.EXPECT().GetChart(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.IntegrationChart{}, errors.New("error fetching chart"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error fetching data points",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
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

				integration.EXPECT().GetDataPoints(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("error fetching data points"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successfully fetched public charting data",
			mockFn: func(dashboardLink *malak_mocks.MockDashboardLinkRepository, integration *malak_mocks.MockIntegrationRepository) {
				dashboardLink.EXPECT().PublicDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.Dashboard{
						ID:          workspaceID,
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

func TestDashboardHandler_PublicChartingDataFetch(t *testing.T) {
	for _, v := range generatePublicChartingDataRequest() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dashboardLinkRepo := malak_mocks.NewMockDashboardLinkRepository(controller)
			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			v.mockFn(dashboardLinkRepo, integrationRepo)

			h := &dashboardHandler{
				dashboardLinkRepo: dashboardLinkRepo,
				integrationRepo:   integrationRepo,
				cfg:               getConfig(),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/public/dashboards/DASH_123/charts/CHART_123", nil)
			req.Header.Add("Content-Type", "application/json")

			router := chi.NewRouter()
			router.Get("/public/dashboards/{reference}/charts/{chart_reference}", WrapMalakHTTPHandler(getLogger(t),
				h.publicChartingDataFetch,
				getConfig(), "dashboards.public_charting_data").ServeHTTP)

			router.ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
