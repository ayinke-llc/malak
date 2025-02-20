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
