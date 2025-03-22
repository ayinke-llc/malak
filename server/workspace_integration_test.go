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
	"github.com/ayinke-llc/malak/internal/integrations"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWorkspaceHandler_PingIntegration(t *testing.T) {
	for _, v := range generateWorkspacePingIntegrationTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			integrationManager := malak_mocks.NewMockIntegrationProviderClient(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)

			v.mockFn(integrationRepo, integrationManager)

			integrationManager.EXPECT().
				Name().
				AnyTimes().
				Return(malak.IntegrationProviderMercury)

			manager := integrations.NewManager()
			manager.Add(malak.IntegrationProviderMercury, integrationManager)

			a := &workspaceHandler{
				cfg:                getConfig(),
				queueClient:        queueRepo,
				integrationRepo:    integrationRepo,
				integrationManager: manager,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "test-ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.pingIntegration, getConfig(), "workspaces.ping.integration").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspacePingIntegrationTestTable() []struct {
	name   string
	mockFn func(integrationRepo *malak_mocks.MockIntegrationRepository,
		integrationManager *malak_mocks.MockIntegrationProviderClient)
	expectedStatusCode int
	req                testAPIIntegrationRequest
} {
	return []struct {
		name   string
		mockFn func(integrationRepo *malak_mocks.MockIntegrationRepository,
			integrationManager *malak_mocks.MockIntegrationProviderClient)
		expectedStatusCode int
		req                testAPIIntegrationRequest
	}{
		{
			name: "invalid request - no api key",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                testAPIIntegrationRequest{},
		},
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration fetching error from database",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch integration"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "invalid integration provider",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "invalid_provider",
							IsEnabled:       true,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration provider valid but no client",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "paystack",
							IsEnabled:       false,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "root integration not enabled yet",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       false,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "only api integration can be pinged",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeOauth2,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "ping error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, errors.New("ping error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "successful ping",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{
						{
							InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
							UserFacingName: "Test Account",
							ProviderID:     "test-id",
							ChartType:      malak.IntegrationChartTypeBar,
							DataPointType:  malak.IntegrationDataPointTypeCurrency,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "cannot ping system integration",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "system",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
	}
}

func TestWorkspaceHandler_EnableIntegration(t *testing.T) {
	for _, v := range generateWorkspaceEnableIntegrationTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			integrationManager := malak_mocks.NewMockIntegrationProviderClient(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)
			secretsClient := malak_mocks.NewMockSecretClient(controller)

			v.mockFn(integrationRepo, integrationManager, secretsClient)

			// Set up the mock provider's Name() method
			integrationManager.EXPECT().
				Name().
				AnyTimes().
				Return(malak.IntegrationProviderMercury)

			manager := integrations.NewManager()
			manager.Add(malak.IntegrationProviderMercury, integrationManager)

			a := &workspaceHandler{
				cfg:                getConfig(),
				queueClient:        queueRepo,
				integrationRepo:    integrationRepo,
				integrationManager: manager,
				secretsClient:      secretsClient,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "test-ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.enableIntegration, getConfig(), "workspaces.enable.integration").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspaceEnableIntegrationTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
	expectedStatusCode int
	req                testAPIIntegrationRequest
} {
	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
		expectedStatusCode int
		req                testAPIIntegrationRequest
	}{
		{
			name: "invalid request - no api key",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                testAPIIntegrationRequest{},
		},
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration failure from database",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("Integration failure"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "root integration not enabled yet",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       false,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration already enabled",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsActive:  true,
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration cannot update if not api key type",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsActive:  true,
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeOauth2,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration provider not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsActive:  true,
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "dmercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration ping error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, errors.New("ping error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "secrets storage error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", errors.New("secrets storage error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "chart creation error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("test-key", nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("chart creation error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "successful enable",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("test-key", nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "cannot enable system integration",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "system",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
	}
}

func generateWorkspaceUpdateAPIKeyTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
	expectedStatusCode int
	req                testAPIIntegrationRequest
} {
	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
		expectedStatusCode int
		req                testAPIIntegrationRequest
	}{
		{
			name: "invalid request - no api key",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                testAPIIntegrationRequest{},
		},
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration db error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("oops"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "root integration not enabled yet",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       false,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "integration not enabled for workspace",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "non-api key integration type",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeOauth2,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "invalid integration provider",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "invalid_provider",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "ping error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, errors.New("ping error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "secrets storage error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", errors.New("secrets storage error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "chart creation error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("test-key", nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("chart creation error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
		{
			name: "successful update",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationManager.EXPECT().
					Ping(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.IntegrationChartValues{}, nil)

				secretsClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return("test-key", nil)

				integrationRepo.EXPECT().
					CreateCharts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: testAPIIntegrationRequest{
				APIKey: "test-key",
			},
		},
	}
}

func TestWorkspaceHandler_UpdateAPIKey(t *testing.T) {
	for _, v := range generateWorkspaceUpdateAPIKeyTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			integrationManager := malak_mocks.NewMockIntegrationProviderClient(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)
			secretsClient := malak_mocks.NewMockSecretClient(controller)

			v.mockFn(integrationRepo, integrationManager, secretsClient)

			// Set up the mock provider's Name() method
			integrationManager.EXPECT().
				Name().
				AnyTimes().
				Return(malak.IntegrationProviderMercury)

			manager := integrations.NewManager()
			manager.Add(malak.IntegrationProviderMercury, integrationManager)

			a := &workspaceHandler{
				cfg:                getConfig(),
				queueClient:        queueRepo,
				integrationRepo:    integrationRepo,
				integrationManager: manager,
				secretsClient:      secretsClient,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "test-ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.updateAPIKeyForIntegration, getConfig(), "workspaces.update.integration").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_DisableIntegration(t *testing.T) {
	for _, v := range generateWorkspaceDisableIntegrationTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
			integrationManager := malak_mocks.NewMockIntegrationProviderClient(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)
			secretsClient := malak_mocks.NewMockSecretClient(controller)

			v.mockFn(integrationRepo, integrationManager, secretsClient)

			integrationManager.EXPECT().
				Name().
				AnyTimes().
				Return(malak.IntegrationProviderMercury)

			manager := integrations.NewManager()
			manager.Add(malak.IntegrationProviderMercury, integrationManager)

			a := &workspaceHandler{
				cfg:                getConfig(),
				queueClient:        queueRepo,
				integrationRepo:    integrationRepo,
				integrationManager: manager,
				secretsClient:      secretsClient,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("reference", "test-ref")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t),
				a.disableIntegration, getConfig(), "workspaces.disable.integration").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspaceDisableIntegrationTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient)
		expectedStatusCode int
	}{
		{
			name: "integration not found",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrWorkspaceIntegrationNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "integration db error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "root integration not enabled",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       false,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "integration already disabled",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: false,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "disable error",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationRepo.EXPECT().
					Disable(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("disable error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "successful disable",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						IsEnabled: true,
						Integration: &malak.Integration{
							IntegrationName: "mercury",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeApiKey,
						},
					}, nil)

				integrationRepo.EXPECT().
					Disable(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "cannot disable system integration",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository, integrationManager *malak_mocks.MockIntegrationProviderClient, secretsClient *malak_mocks.MockSecretClient) {
				integrationRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.WorkspaceIntegration{
						Integration: &malak.Integration{
							IntegrationName: "system",
							IsEnabled:       true,
							IntegrationType: malak.IntegrationTypeSystem,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
}
