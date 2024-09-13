package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockReferenceGenerator struct{}

func (m *mockReferenceGenerator) Generate(
	e malak.EntityType) malak.Reference {
	return malak.Reference(fmt.Sprintf("%s_%s", e.String(), "oopsoops"))
}

func TestContactHandler_Create(t *testing.T) {
	for _, v := range generateContactTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			contactRepo := malak_mocks.NewMockContactRepository(controller)

			v.mockFn(contactRepo)

			a := &contactHandler{
				cfg:                getConfig(),
				contactRepo:        contactRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), a.Create, getConfig(), "contacts.new").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateContactTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	req                createContactRequest
} {

	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		req                createContactRequest
	}{
		{
			name: "no email provided",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {

			},
			expectedStatusCode: http.StatusBadRequest,
			req:                createContactRequest{},
		},
		{
			name: "duplicate contact",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(malak.ErrContactExists)
			},
			expectedStatusCode: http.StatusConflict,
			req: createContactRequest{
				Email: malak.Email("oopsoops@gmail.com"),
			},
		},
		{
			name: "unknown error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("oops"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactRequest{
				Email: malak.Email("oopsoops@gmail.com"),
			},
		},
		{
			name: "success",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactRequest{
				Email: malak.Email("oopsoops@gmail.com"),
			},
		},
	}
}
