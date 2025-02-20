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
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{ID: uuid.New()}))

			WrapMalakHTTPHandler(getLogger(t),
				h.list,
				getConfig(), "dashboards.list").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}
