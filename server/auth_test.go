package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	mock_jwttoken "github.com/ayinke-llc/malak/internal/pkg/jwttoken/mocks"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	socialauth_mocks "github.com/ayinke-llc/malak/internal/pkg/socialauth/mocks"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sebdah/goldie/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func verifyMatch(t *testing.T, v interface{}) {
	g := goldie.New(t, goldie.WithFixtureDir("./testdata"))

	b := new(bytes.Buffer)

	if d, ok := v.(*httptest.ResponseRecorder); ok {
		_, err := io.Copy(b, d.Body)
		require.NoError(t, err)
	} else {
		err := json.NewEncoder(b).Encode(v)
		require.NoError(t, err)
	}

	g.Assert(t, t.Name(), b.Bytes())
}

func getConfig() config.Config {
	return config.Config{
		Otel: struct {
			Endpoint  string "yaml:\"endpoint\" mapstructure:\"endpoint\""
			UseTLS    bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			Headers   string "yaml:\"headers\" mapstructure:\"headers\""
			IsEnabled bool   "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
		}{
			IsEnabled: false,
		},
		Auth: struct {
			Google struct {
				ClientID     string   "yaml:\"client_id\" mapstructure:\"client_id\""
				ClientSecret string   "yaml:\"client_secret\" mapstructure:\"client_secret\""
				RedirectURI  string   "yaml:\"redirect_uri\" mapstructure:\"redirect_uri\""
				Scopes       []string "yaml:\"scopes\" mapstructure:\"scopes\""
				IsEnabled    bool     "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
			} "yaml:\"google\" mapstructure:\"google\""
			JWT struct {
				Key string "yaml:\"key\" mapstructure:\"key\""
			} "yaml:\"jwt\" mapstructure:\"jwt\""
		}{
			JWT: struct {
				Key string "yaml:\"key\" mapstructure:\"key\""
			}{
				Key: "a907e75f80910f5dc5b8c677de1de611ffa80be9d7d9f9dd614c8c7846db1062",
			},
		},
	}
}

func TestAuthHandler_Login(t *testing.T) {
	for _, v := range generateLoginTestTable() {

		t.Run(v.name, func(t *testing.T) {

			logrus.SetOutput(io.Discard)

			logger := logrus.WithField("test", true)

			controller := gomock.NewController(t)
			defer controller.Finish()

			googleCfg := socialauth_mocks.NewMockSocialAuthProvider(controller)
			userRepo := malak_mocks.NewMockUserRepository(controller)

			jwtMock := mock_jwttoken.NewMockJWTokenManager(controller)

			v.mockFn(googleCfg, userRepo)

			a := &authHandler{
				logger:       logger,
				cfg:          getConfig(),
				googleCfg:    googleCfg,
				userRepo:     userRepo,
				tokenManager: jwtMock,
			}

			var b = bytes.NewBuffer(nil)

			require.NoError(t, json.NewEncoder(b).Encode(&v.req))

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("provider", v.provider)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
			req.Header.Add("Content-Type", "application/json")

			if v.expectedStatusCode == http.StatusOK {
				jwtMock.EXPECT().
					GenerateJWToken(gomock.Any()).
					Times(1).
					Return(jwttoken.JWTokenData{
						Token:  "b622268d-4512-4e3c-98da-88097753d4b9",
						UserID: uuid.MustParse("7e6ad0c8-7a96-4add-a270-52615bd808e6"),
					}, nil)
			}

			WrapMalakHTTPHandler(a.Login, getConfig(), "Auth.Login").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateLoginTestTable() []struct {
	name               string
	mockFn             func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository)
	expectedStatusCode int
	req                authenticateUserRequest
	provider           string
} {

	var reusedID = uuid.MustParse("37f41afb-afff-45cc-bcc0-71249d95df90")

	return []struct {
		name               string
		mockFn             func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository)
		expectedStatusCode int
		req                authenticateUserRequest
		provider           string
	}{
		{
			name: "no code to exchange provided",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), gomock.Any()).
					Times(0)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			provider:           "google",
			expectedStatusCode: http.StatusBadRequest,
			req:                authenticateUserRequest{},
		},
		{
			name: "token exchange fails",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "invalid-token",
					}).
					Times(1).
					Return(nil, errors.New("could not valdate token"))

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: authenticateUserRequest{
				Code: "invalid-token",
			},
			provider: "google",
		},
		{
			name: "could not fetch user details",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "token",
					}).
					Times(1).
					Return(&oauth2.Token{
						AccessToken: "access-token",
					}, nil)

				googleMock.EXPECT().
					User(gomock.Any(), gomock.Any()).
					Times(1).
					Return(socialauth.User{}, errors.New("could not fetch user"))

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			req: authenticateUserRequest{
				Code: "token",
			},
			provider: "google",
		},
		{
			name: "duplicate email. user gets logged in inside but could not fetch details from db",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "token",
					}).
					Times(1).
					Return(&oauth2.Token{
						AccessToken: "access-token",
					}, nil)

				user := socialauth.User{
					Email: "test@test.com",
					Name:  "TEST TEST",
				}

				googleMock.EXPECT().
					User(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.ErrUserExists)

				userRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not fetch user"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: authenticateUserRequest{
				Code: "token",
			},
			provider: "google",
		},
		{
			name: "duplicate email. user gets logged in",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "token",
					}).
					Times(1).
					Return(&oauth2.Token{
						AccessToken: "access-token",
					}, nil)

				user := socialauth.User{
					Email: "test@test.com",
					Name:  "TEST TEST",
				}

				googleMock.EXPECT().
					User(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(malak.ErrUserExists)

				userRepo.EXPECT().Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.User{
						ID: reusedID,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			req: authenticateUserRequest{
				Code: "token",
			},
			provider: "google",
		},
		{
			name: "could not create user in datastore",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "token",
					}).
					Times(1).
					Return(&oauth2.Token{
						AccessToken: "access-token",
					}, nil)

				user := socialauth.User{
					Email: "test@test.com",
					Name:  "TEST TEST",
				}

				googleMock.EXPECT().
					User(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("unknown error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			req: authenticateUserRequest{
				Code: "token",
			},
			provider: "google",
		},
		{
			name: "user was succesfully created",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				googleMock.EXPECT().
					Validate(gomock.Any(), socialauth.ValidateOptions{
						Code: "token",
					}).
					Times(1).
					Return(&oauth2.Token{
						AccessToken: "access-token",
					}, nil)

				user := socialauth.User{
					Email: "test@test.com",
					Name:  "TEST TEST",
				}

				googleMock.EXPECT().
					User(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				userRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			req: authenticateUserRequest{
				Code: "token",
			},
			provider: "google",
		},
	}
}
