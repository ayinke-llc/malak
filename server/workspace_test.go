package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func getLogger(t *testing.T) *zap.Logger {

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	return logger
}

func TestWorkspaceHandler_SwitchWorkspace(t *testing.T) {
	for _, v := range generateWorkspaceSwitchTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
			userRepo := malak_mocks.NewMockUserRepository(controller)

			v.mockFn(workspaceRepo, userRepo)

			a := &workspaceHandler{
				cfg:           getConfig(),
				workspaceRepo: workspaceRepo,
				userRepo:      userRepo,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{
				Metadata: &malak.UserMetadata{},
			}))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("provider", "reference")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			WrapMalakHTTPHandler(getLogger(t), a.switchCurrentWorkspaceForUser, getConfig(), "workspaces.switch").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_Create(t *testing.T) {
	for _, v := range generateWorkspaceTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
			planRepo := malak_mocks.NewMockPlanRepository(controller)

			v.mockFn(workspaceRepo, planRepo)

			a := &workspaceHandler{
				cfg:           getConfig(),
				workspaceRepo: workspaceRepo,
				planRepo:      planRepo,
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

			WrapMalakHTTPHandler(getLogger(t), a.createWorkspace, getConfig(), "workspaces.new").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_Update(t *testing.T) {
	for _, v := range generateWorkspaceUpdateTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
			planRepo := malak_mocks.NewMockPlanRepository(controller)

			v.mockFn(workspaceRepo, planRepo)

			a := &workspaceHandler{
				cfg:           getConfig(),
				workspaceRepo: workspaceRepo,
				planRepo:      planRepo,
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

			WrapMalakHTTPHandler(getLogger(t), a.updateWorkspace, getConfig(), "workspaces.update").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_GetPreferences(t *testing.T) {
	for _, v := range generateWorkspacePreferencesTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			prefRepo := malak_mocks.NewMockPreferenceRepository(controller)

			v.mockFn(prefRepo)

			a := &workspaceHandler{
				cfg:            getConfig(),
				preferenceRepo: prefRepo,
				referenceGenerationFunc: func(e malak.EntityType) string {
					return "workspace_tt7-YieIgz"
				},
			}

			var b = bytes.NewBuffer(nil)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t),
				a.getPreferences, getConfig(), "workspaces.update").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_UpdatePreferences(t *testing.T) {
	for _, v := range generateWorkspacePreferencesUpdateTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			prefRepo := malak_mocks.NewMockPreferenceRepository(controller)

			v.mockFn(prefRepo)

			a := &workspaceHandler{
				cfg:            getConfig(),
				preferenceRepo: prefRepo,
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

			WrapMalakHTTPHandler(getLogger(t),
				a.updatePreferences, getConfig(), "workspaces.update").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestWorkspaceHandler_getIntegrations(t *testing.T) {
	for _, v := range generateWorkspaceGetIntegrationsTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)

			v.mockFn(integrationRepo)

			a := &workspaceHandler{
				cfg:             getConfig(),
				integrationRepo: integrationRepo,
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

			WrapMalakHTTPHandler(getLogger(t),
				a.getIntegrations, getConfig(), "workspaces.list.integrations").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspaceGetIntegrationsTestTable() []struct {
	name               string
	mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
	expectedStatusCode int
	req                updatePreferencesRequest
} {

	return []struct {
		name               string
		mockFn             func(integrationRepo *malak_mocks.MockIntegrationRepository)
		expectedStatusCode int
		req                updatePreferencesRequest
	}{
		{
			name: "could not fetch workspace integrations",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch integrations"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "listed integrations successfully",
			mockFn: func(integrationRepo *malak_mocks.MockIntegrationRepository) {
				integrationRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.WorkspaceIntegration{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func generateWorkspacePreferencesUpdateTestTable() []struct {
	name               string
	mockFn             func(preferencesRepo *malak_mocks.MockPreferenceRepository)
	expectedStatusCode int
	req                updatePreferencesRequest
} {

	return []struct {
		name               string
		mockFn             func(preferencesRepo *malak_mocks.MockPreferenceRepository)
		expectedStatusCode int
		req                updatePreferencesRequest
	}{
		{
			name: "could not fetch workspace preferences",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch preferences"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update fails",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Preference{}, nil)

				preferencesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("updating prefernces failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "update succeeds",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Preference{}, nil)

				preferencesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "update succeeds with request data",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Preference{}, nil)

				preferencesRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updatePreferencesRequest{
				Preferences: struct {
					Billing    malak.BillingPreferences       "json:\"billing,omitempty\" validate:\"required\""
					Newsletter malak.CommunicationPreferences "json:\"newsletter,omitempty\" validate:\"required\""
				}{
					Newsletter: malak.CommunicationPreferences{
						EnableMarketing:      true,
						EnableProductUpdates: true,
					},
				},
			},
		},
	}
}

func generateWorkspacePreferencesTestTable() []struct {
	name               string
	mockFn             func(preferencesRepo *malak_mocks.MockPreferenceRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(preferencesRepo *malak_mocks.MockPreferenceRepository)
		expectedStatusCode int
	}{
		{
			name: "could not fetch workspace preferences",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch preferences"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fetched workspace preferences",
			mockFn: func(preferencesRepo *malak_mocks.MockPreferenceRepository) {
				preferencesRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Preference{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func generateWorkspaceUpdateTestTable() []struct {
	name               string
	mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository)
	expectedStatusCode int
	req                updateWorkspaceRequest
} {

	return []struct {
		name               string
		mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository)
		expectedStatusCode int
		req                updateWorkspaceRequest
	}{
		{
			name: "invalid timezone provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: updateWorkspaceRequest{
				Timezone: hermes.Ref("oops/oops"),
			},
		},
		{
			name: "invalid image provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: updateWorkspaceRequest{
				Timezone: hermes.Ref("Africa/Algiers"),
				Logo:     hermes.Ref("https://google.com"),
			},
		},
		{
			name: "invalid workspace name provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: updateWorkspaceRequest{
				Timezone:      hermes.Ref("Africa/Algiers"),
				Logo:          hermes.Ref("https://images.unsplash.com/photo-1737467023078-a694673d7cb3"),
				WorkspaceName: hermes.Ref("1234"),
			},
		},
		{
			name: "update fails",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("updating workspace failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: updateWorkspaceRequest{
				Timezone:      hermes.Ref("Africa/Algiers"),
				Logo:          hermes.Ref("https://images.unsplash.com/photo-1737467023078-a694673d7cb3"),
				WorkspaceName: hermes.Ref("12345"),
			},
		},
		{
			name: "update succeeds",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updateWorkspaceRequest{
				Timezone:      hermes.Ref("Africa/Algiers"),
				Logo:          hermes.Ref("https://images.unsplash.com/photo-1737467023078-a694673d7cb3"),
				WorkspaceName: hermes.Ref("12345"),
			},
		},
		{
			name: "update succeeds even if partial fields provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: updateWorkspaceRequest{
				Timezone: hermes.Ref("Africa/Algiers"),
			},
		},
		{
			name: "update succeeds if no fields provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req:                updateWorkspaceRequest{},
		},
	}
}

func generateWorkspaceTestTable() []struct {
	name               string
	mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository)
	expectedStatusCode int
	req                createWorkspaceRequest
} {

	return []struct {
		name               string
		mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository)
		expectedStatusCode int
		req                createWorkspaceRequest
	}{
		{
			name: "no name provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: createWorkspaceRequest{
				Name: "",
			},
		},
		{
			name: "invalid name provided",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req: createWorkspaceRequest{
				Name: "iii",
			},
		},
		{
			name: "could not fetch plan",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				planRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch plan"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createWorkspaceRequest{
				Name: "workspance name",
			},
		},
		{
			name: "could not create workspace",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				planRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Plan{}, nil)

				workspaceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("oops"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createWorkspaceRequest{
				Name: "workspance name",
			},
		},
		{
			name: "created workspace",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository) {
				planRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Plan{}, nil)

				workspaceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createWorkspaceRequest{
				Name: "workspance name",
			},
		},
	}
}

func generateWorkspaceSwitchTable() []struct {
	name               string
	mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, userRepo *malak_mocks.MockUserRepository)
	expectedStatusCode int
} {

	return []struct {
		name               string
		mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, userRepo *malak_mocks.MockUserRepository)
		expectedStatusCode int
	}{
		{
			name: "could not find reference",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, usr *malak_mocks.MockUserRepository) {
				workspaceRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not find workspace"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "could not update user repo",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, userRepo *malak_mocks.MockUserRepository) {
				workspaceRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&malak.Workspace{}, nil)

				userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("could not update"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "updated current workspace",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository, userRepo *malak_mocks.MockUserRepository) {
				workspaceRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&malak.Workspace{}, nil)

				userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}
