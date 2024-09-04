package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWorkspaceHandler_Create(t *testing.T) {
	for _, v := range generateWorkspaceTestTable() {

		t.Run(v.name, func(t *testing.T) {

			logrus.SetOutput(io.Discard)

			controller := gomock.NewController(t)
			defer controller.Finish()

			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
			planRepo := malak_mocks.NewMockPlanRepository(controller)

			v.mockFn(workspaceRepo, planRepo)

			a := &workspaceHandler{
				cfg:           getConfig(),
				workspaceRepo: workspaceRepo,
				planRepo:      planRepo,
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))

			WrapMalakHTTPHandler(a.createWorkspace, getConfig(), "workspaces.new").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateWorkspaceTestTable() []struct {
	name               string
	mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository, planRepo *malak_mocks.MockPlanRepository)
	expectedStatusCode int
	req                createWorkspaceRequest
} {

	// var reusedID = uuid.MustParse("37f41afb-afff-45cc-bcc0-71249d95df90")

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
	}
}
