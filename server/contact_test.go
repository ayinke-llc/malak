package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockReferenceGenerator struct{}

func (m *mockReferenceGenerator) Token() string {
	return "oops"
}

func (m *mockReferenceGenerator) ShortLink() string {
	return "oops"
}

func (m *mockReferenceGenerator) Generate(e malak.EntityType) malak.Reference {
	return malak.Reference(fmt.Sprintf("%s_%s", e.String(), "test_reference"))
}

func TestContactHandler_Create(t *testing.T) {
	for _, v := range generateCreateContactTestTable() {
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

func TestContactHandler_AddUserToContactList(t *testing.T) {
	for _, v := range generateAddUserToContactListTestTable() {
		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			contactListRepo := malak_mocks.NewMockContactListRepository(controller)
			contactRepo := malak_mocks.NewMockContactRepository(controller)

			v.mockFn(contactListRepo, contactRepo)

			a := &contactHandler{
				cfg:                getConfig(),
				contactListRepo:    contactListRepo,
				referenceGenerator: &mockReferenceGenerator{},
				contactRepo:        contactRepo,
			}

			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(&v.req))
			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/contacts/lists", b)
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), a.addUserToContactList, getConfig(), "contacts.list.add").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)

			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_CreateContactList(t *testing.T) {
	for _, v := range generateCreateContactListTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactListRepo := malak_mocks.NewMockContactListRepository(controller)
			v.mockFn(contactListRepo)
			a := &contactHandler{
				cfg:                getConfig(),
				contactListRepo:    contactListRepo,
				referenceGenerator: &mockReferenceGenerator{},
			}
			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(&v.req))
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/contacts/lists", b)
			req.Header.Add("Content-Type", "application/json")
			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			WrapMalakHTTPHandler(getLogger(t), a.createContactList, getConfig(), "contacts.list.new").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_FetchContactLists(t *testing.T) {
	for _, v := range generateFetchContactListsTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactListRepo := malak_mocks.NewMockContactListRepository(controller)
			v.mockFn(contactListRepo)
			a := &contactHandler{
				cfg:             getConfig(),
				contactListRepo: contactListRepo,
			}
			rr := httptest.NewRecorder()
			url := "/contacts/lists"
			if v.includeEmails {
				url += "?include_emails=true"
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			}))
			WrapMalakHTTPHandler(getLogger(t), a.fetchContactLists, getConfig(), "contacts.list.fetch").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_EditContactList(t *testing.T) {
	for _, v := range generateEditContactListTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactListRepo := malak_mocks.NewMockContactListRepository(controller)
			v.mockFn(contactListRepo)
			a := &contactHandler{
				cfg:             getConfig(),
				contactListRepo: contactListRepo,
			}
			var b = bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(&v.req))
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/contacts/lists/test_reference", b)
			req.Header.Add("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("reference", "test_reference")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			WrapMalakHTTPHandler(getLogger(t), a.editContactList, getConfig(), "contacts.list.edit").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_DeleteContactList(t *testing.T) {
	for _, v := range generateDeleteContactListTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactListRepo := malak_mocks.NewMockContactListRepository(controller)
			v.mockFn(contactListRepo)
			a := &contactHandler{
				cfg:             getConfig(),
				contactListRepo: contactListRepo,
			}
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/contacts/lists/test_reference", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("reference", "test_reference")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			WrapMalakHTTPHandler(getLogger(t), a.deleteContactList, getConfig(), "contacts.list.delete").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_List(t *testing.T) {
	for _, v := range generateListContactsTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactRepo := malak_mocks.NewMockContactRepository(controller)
			v.mockFn(contactRepo)
			a := &contactHandler{
				cfg:         getConfig(),
				contactRepo: contactRepo,
			}
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
			if v.withPagination {
				q := req.URL.Query()
				q.Add("page", "2")
				q.Add("per_page", "20")
				req.URL.RawQuery = q.Encode()
			}
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			WrapMalakHTTPHandler(getLogger(t), a.list, getConfig(), "contacts.list").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_GetSingleContact(t *testing.T) {
	for _, v := range generateContactFetchContact() {
		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			contactRepo := malak_mocks.NewMockContactRepository(controller)
			shareRepo := malak_mocks.NewMockContactShareRepository(controller)

			v.mockFn(contactRepo, shareRepo)

			a := &contactHandler{
				cfg:              getConfig(),
				contactRepo:      contactRepo,
				contactShareRepo: shareRepo,
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), a.fetchContact, getConfig(), "contacts.fetchContact").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_DeleteContact(t *testing.T) {
	for _, v := range generateDeleteContactTestTable() {
		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			contactRepo := malak_mocks.NewMockContactRepository(controller)

			v.mockFn(contactRepo)

			a := &contactHandler{
				cfg:         getConfig(),
				contactRepo: contactRepo,
			}

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, "/contacts", nil)
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))

			WrapMalakHTTPHandler(getLogger(t), a.deleteContact, getConfig(), "contacts.delete").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_Edit(t *testing.T) {
	for _, v := range generateEditContactTestTable() {
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
			WrapMalakHTTPHandler(getLogger(t), a.editContact, getConfig(), "contacts.edit").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_BatchCreate(t *testing.T) {
	for _, v := range generateBatchCreateTestTable() {
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
			req := httptest.NewRequest(http.MethodPost, "/contacts/batch", b)
			req.Header.Add("Content-Type", "application/json")
			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			WrapMalakHTTPHandler(getLogger(t), a.batchCreate, getConfig(), "contacts.batch").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestContactHandler_Search(t *testing.T) {
	for _, v := range generateSearchContactTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			contactRepo := malak_mocks.NewMockContactRepository(controller)
			v.mockFn(contactRepo)
			a := &contactHandler{
				cfg:         getConfig(),
				contactRepo: contactRepo,
			}
			rr := httptest.NewRecorder()
			url := "/contacts/search"
			if v.searchValue != "" {
				url += "?search=" + v.searchValue
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			}))
			WrapMalakHTTPHandler(getLogger(t), a.search, getConfig(), "contacts.search").
				ServeHTTP(rr, req)
			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateEditContactTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	req                editContactRequest
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		req                editContactRequest
	}{
		{
			name: "first name too short",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: "abc",
			},
		},
		{
			name: "first name too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: strings.Repeat(faker.Name(), 50),
			},
		},
		{
			name: "last name too short",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  "abc",
			},
		},
		{
			name: "last name too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  strings.Repeat(faker.Name(), 1000),
			},
		},
		{
			name: "company name too short",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   "abc",
			},
		},
		{
			name: "company name too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   strings.Repeat(faker.Name(), 100),
			},
		},
		{
			name: "address too short",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   "abc",
			},
		},
		{
			name: "address too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   strings.Repeat(faker.Name(), 1000),
			},
		},
		{
			name: "notes too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   faker.Name(),
				Notes:     strings.Repeat(faker.Name(), 1000),
			},
		},
		{
			name: "notes too short",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   faker.Name(),
				Notes:     "abc",
			},
		},
		{
			name: "could not fetch contact due to db error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch from db"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   faker.Name(),
				Notes:     faker.Sentence(),
			},
		},
		{
			name: "could not fetch contact because contact does not exists",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: editContactRequest{
				FirstName: faker.Name(),
				LastName:  faker.Name(),
				Company:   faker.Name(),
				Address:   faker.Name(),
				Notes:     faker.Sentence(),
			},
		},
		{
			name: "updating contact fails at db level",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not update"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: editContactRequest{
				FirstName: "first name",
				LastName:  "last name",
				Company:   "malak inc",
				Address:   "240 Delaware",
				Notes:     "here is my random note",
			},
		},
		{
			name: "updating contact succeeds",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: editContactRequest{
				FirstName: "first name",
				LastName:  "last name",
				Company:   "malak inc",
				Address:   "240 Delaware",
				Notes:     "here is my random note",
			},
		},
		{
			name: "updating contact succeeds with partial data",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: editContactRequest{
				Address: "240 Delaware",
				Notes:   "here is my random note",
			},
		},
		{
			name: "updating contact succeeds even without any data in request",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req:                editContactRequest{},
		},
	}
}

func generateDeleteContactTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
	}{
		{
			name: "contact not exists",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "contact error while fetching",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "deleting contact fails",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("could not create contact"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "deleting contact succeeds",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func generateContactFetchContact() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository)
		expectedStatusCode int
	}{
		{
			name: "get contact fails because it does not exists",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository) {

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)

			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "get contact fails because of db error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository) {

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("failed"))

			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "getting shared items failed",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository) {

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				shareRepo.EXPECT().
					All(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch shared items"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "getting shared items succeeds",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository, shareRepo *malak_mocks.MockContactShareRepository) {

				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				shareRepo.EXPECT().
					All(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.ContactShareItem{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func generateCreateContactTestTable() []struct {
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
				Email: malak.Email("test@example.com"),
			},
		},
		{
			name: "unknown error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactRequest{
				Email: malak.Email("test@example.com"),
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
				Email: malak.Email("test@example.com"),
			},
		},
	}
}

func generateAddUserToContactListTestTable() []struct {
	name   string
	mockFn func(contactListRepo *malak_mocks.MockContactListRepository,
		contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	req                addContactToListRequest
} {
	return []struct {
		name               string
		mockFn             func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		req                addContactToListRequest
	}{
		{
			name: "no reference provided",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                addContactToListRequest{},
		},
		{
			name: "contact not found",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, malak.ErrContactNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
		{
			name: "error fetching contact",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
		{
			name: "list not found",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactListRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, malak.ErrContactListNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
		{
			name: "Error while fetching list",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactListRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, errors.New("oops"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
		{
			name: "could not create contact list mappings",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactListRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&malak.ContactList{}, nil)

				contactListRepo.EXPECT().Add(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unkwown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
		{
			name: "added contact to list",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository, contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Contact{}, nil)

				contactListRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).Return(&malak.ContactList{}, nil)

				contactListRepo.EXPECT().
					Add(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: addContactToListRequest{
				Reference: "oops",
			},
		},
	}
}

func generateCreateContactListTestTable() []struct {
	name               string
	mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
	expectedStatusCode int
	req                createContactListRequest
} {
	return []struct {
		name               string
		mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
		expectedStatusCode int
		req                createContactListRequest
	}{
		{
			name: "no name provided",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                createContactListRequest{},
		},
		{
			name: "name too long",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createContactListRequest{
				Name: "This is a very long name that exceeds the 50 character limit for contact list names",
			},
		},
		{
			name: "unknown error",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactListRequest{
				Name: "Test List",
			},
		},
		{
			name: "success",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactListRequest{
				Name: "Test List",
			},
		},
	}
}

func generateFetchContactListsTestTable() []struct {
	name               string
	mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
	expectedStatusCode int
	includeEmails      bool
} {
	return []struct {
		name               string
		mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
		expectedStatusCode int
		includeEmails      bool
	}{
		{
			name: "unknown error",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(nil, nil, errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "empty lists",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return([]malak.ContactList{}, []malak.ContactListMappingWithContact{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "lists with no contacts",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				lists := []malak.ContactList{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Title:     "Test List",
						Reference: "list_test",
					},
				}
				contactListRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(lists, []malak.ContactListMappingWithContact{}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "lists with contacts",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				lists := []malak.ContactList{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Title:     "Test List",
						Reference: "list_test",
					},
				}
				mappings := []malak.ContactListMappingWithContact{
					{
						ListID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						ContactID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
						Email:     "test@example.com",
					},
				}
				contactListRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(lists, mappings, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "lists with include_emails parameter",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				lists := []malak.ContactList{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Title:     "Test List",
						Reference: "list_test",
					},
				}
				mappings := []malak.ContactListMappingWithContact{
					{
						ListID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						ContactID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
						Email:     "test@example.com",
					},
				}
				contactListRepo.EXPECT().List(gomock.Any(), &malak.ContactListOptions{
					WorkspaceID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					IncludeEmails: true,
				}).Return(lists, mappings, nil)
			},
			expectedStatusCode: http.StatusOK,
			includeEmails:      true,
		},
	}
}

func generateEditContactListTestTable() []struct {
	name               string
	mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
	expectedStatusCode int
	req                createContactListRequest
} {
	return []struct {
		name               string
		mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
		expectedStatusCode int
		req                createContactListRequest
	}{
		{
			name: "no name provided",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                createContactListRequest{},
		},
		{
			name: "contact list not found",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(nil, malak.ErrContactListNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			req: createContactListRequest{
				Name: "Updated List",
			},
		},
		{
			name: "unknown error during fetch",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactListRequest{
				Name: "Updated List",
			},
		},
		{
			name: "unknown error during update",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(&malak.ContactList{Title: "Old Name"}, nil)
				contactListRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactListRequest{
				Name: "Updated List",
			},
		},
		{
			name: "success",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(&malak.ContactList{Title: "Old Name"}, nil)
				contactListRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactListRequest{
				Name: "Updated List",
			},
		},
	}
}

func generateDeleteContactListTestTable() []struct {
	name               string
	mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
	expectedStatusCode int
} {
	return []struct {
		name               string
		mockFn             func(contactListRepo *malak_mocks.MockContactListRepository)
		expectedStatusCode int
	}{
		{
			name: "contact list not found",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(nil, malak.ErrContactListNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "unknown error during fetch",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "unknown error during delete",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(&malak.ContactList{}, nil)
				contactListRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			mockFn: func(contactListRepo *malak_mocks.MockContactListRepository) {
				contactListRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Return(&malak.ContactList{}, nil)
				contactListRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
	}
}

func generateListContactsTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	withPagination     bool
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		withPagination     bool
	}{
		{
			name: "unknown error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "empty contacts",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return([]malak.Contact{}, int64(0), nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "with contacts",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contacts := []malak.Contact{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Email:     "test@example.com",
						Reference: "contact_test",
						FirstName: "Test",
						LastName:  "User",
					},
				}
				contactRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(contacts, int64(1), nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "with pagination",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contacts := []malak.Contact{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Email:     "test@example.com",
						Reference: "contact_test",
						FirstName: "Test",
						LastName:  "User",
					},
				}
				contactRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(contacts, int64(25), nil)
			},
			expectedStatusCode: http.StatusCreated,
			withPagination:     true,
		},
	}
}

func generateBatchCreateTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	req                createContactRequestBatch
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		req                createContactRequestBatch
	}{
		{
			name: "empty contacts array",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req:                createContactRequestBatch{},
		},
		{
			name: "missing email",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref("Doe"),
					},
				},
			},
		},
		{
			name: "first name too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test@example.com",
						FirstName: hermes.Ref(strings.Repeat("a", 101)),
						LastName:  hermes.Ref("Doe"),
					},
				},
			},
		},
		{
			name: "last name too long",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test@example.com",
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref(strings.Repeat("a", 101)),
					},
				},
			},
		},
		{
			name: "duplicate contact",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(malak.ErrContactExists)
			},
			expectedStatusCode: http.StatusConflict,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test@example.com",
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref("Doe"),
					},
				},
			},
		},
		{
			name: "unknown error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test@example.com",
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref("Doe"),
					},
				},
			},
		},
		{
			name: "success with single contact",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test@example.com",
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref("Doe"),
					},
				},
			},
		},
		{
			name: "success with multiple contacts",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:     "test1@example.com",
						FirstName: hermes.Ref("John"),
						LastName:  hermes.Ref("Doe"),
					},
					{
						Email:     "test2@example.com",
						FirstName: hermes.Ref("Jane"),
						LastName:  hermes.Ref("Smith"),
					},
				},
			},
		},
		{
			name: "success with email as first name",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			req: createContactRequestBatch{
				Contacts: []struct {
					Email     malak.Email `json:"email,omitempty" validate:"'required'"`
					FirstName *string     `json:"first_name,omitempty" validate:"'required'"`
					LastName  *string     `json:"last_name,omitempty" validate:"'required'"`
				}{
					{
						Email:    "test@example.com",
						LastName: hermes.Ref("Doe"),
					},
				},
			},
		},
	}
}

func generateSearchContactTestTable() []struct {
	name               string
	mockFn             func(contactRepo *malak_mocks.MockContactRepository)
	expectedStatusCode int
	searchValue        string
} {
	return []struct {
		name               string
		mockFn             func(contactRepo *malak_mocks.MockContactRepository)
		expectedStatusCode int
		searchValue        string
	}{
		{
			name: "search with no query parameter",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				// No expectation since the handler returns 400 before calling Search
			},
			expectedStatusCode: http.StatusBadRequest,
			searchValue:        "",
		},
		{
			name: "search with query parameter",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Search(gomock.Any(), malak.SearchContactOptions{
					WorkspaceID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					SearchValue: "john",
				}).Return([]malak.Contact{
					{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
						Email:     "john@example.com",
						FirstName: "John",
						LastName:  "Doe",
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			searchValue:        "john",
		},
		{
			name: "search with error",
			mockFn: func(contactRepo *malak_mocks.MockContactRepository) {
				contactRepo.EXPECT().Search(gomock.Any(), malak.SearchContactOptions{
					WorkspaceID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					SearchValue: "error",
				}).Return(nil, errors.New("search failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			searchValue:        "error",
		},
	}
}
