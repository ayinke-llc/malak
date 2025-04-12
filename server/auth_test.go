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
	"time"

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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

var webhookSecret = "wh_sec"

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
		Logging: struct {
			Mode config.LogMode "yaml:\"mode\" mapstructure:\"mode\""
		}{
			Mode: config.LogModeDev,
		},
		APIKey: struct {
			HashSecret string "mapstructure:\"hash_secret\" yaml:\"hash_secret\""
		}{
			HashSecret: "1234597u8tysdhfjhfk",
		},
		Database: struct {
			DatabaseType config.DatabaseType "yaml:\"database_type\" mapstructure:\"database_type\""
			Postgres     struct {
				DSN          string        "yaml:\"dsn\" mapstructure:\"dsn\""
				LogQueries   bool          "yaml:\"log_queries\" mapstructure:\"log_queries\""
				QueryTimeout time.Duration "yaml:\"query_timeout\" mapstructure:\"query_timeout\""
			} "yaml:\"postgres\" mapstructure:\"postgres\""
			Redis struct {
				DSN string "yaml:\"dsn\" mapstructure:\"dsn\""
			} "yaml:\"redis\" mapstructure:\"redis\""
		}{
			DatabaseType: config.DatabaseTypePostgres,
		},
		HTTP: config.HTTPConfig{
			Port: 8000,
		},
		Uploader: struct {
			Driver        config.UploadDriver "yaml:\"driver\" mapstructure:\"driver\""
			MaxUploadSize int64               "yaml:\"max_upload_size\" mapstructure:\"max_upload_size\""
			S3            struct {
				AccessKey                  string "yaml:\"access_key\" mapstructure:\"access_key\""
				AccessSecret               string "yaml:\"access_secret\" mapstructure:\"access_secret\""
				Region                     string "yaml:\"region\" mapstructure:\"region\""
				Endpoint                   string "yaml:\"endpoint\" mapstructure:\"endpoint\""
				LogOperations              bool   "yaml:\"log_operations\" mapstructure:\"log_operations\""
				Bucket                     string "yaml:\"bucket\" mapstructure:\"bucket\""
				DeckBucket                 string "yaml:\"deck_bucket\" mapstructure:\"deck_bucket\""
				UseTLS                     bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
				CloudflareBucketDomain     string "yaml:\"cloudflare_bucket_domain\" mapstructure:\"cloudflare_bucket_domain\""
				CloudflareDeckBucketDomain string "yaml:\"cloudflare_deck_bucket_domain\" mapstructure:\"cloudflare_deck_bucket_domain\""
			} "yaml:\"s3\" mapstructure:\"s3\""
		}{
			Driver: config.UploadDriverS3,
			S3: struct {
				AccessKey                  string "yaml:\"access_key\" mapstructure:\"access_key\""
				AccessSecret               string "yaml:\"access_secret\" mapstructure:\"access_secret\""
				Region                     string "yaml:\"region\" mapstructure:\"region\""
				Endpoint                   string "yaml:\"endpoint\" mapstructure:\"endpoint\""
				LogOperations              bool   "yaml:\"log_operations\" mapstructure:\"log_operations\""
				Bucket                     string "yaml:\"bucket\" mapstructure:\"bucket\""
				DeckBucket                 string "yaml:\"deck_bucket\" mapstructure:\"deck_bucket\""
				UseTLS                     bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
				CloudflareBucketDomain     string "yaml:\"cloudflare_bucket_domain\" mapstructure:\"cloudflare_bucket_domain\""
				CloudflareDeckBucketDomain string "yaml:\"cloudflare_deck_bucket_domain\" mapstructure:\"cloudflare_deck_bucket_domain\""
			}{
				AccessKey:    "test-key",
				AccessSecret: "test-secret",
				Region:       "us-east-1",
				Bucket:       "test-bucket",
				DeckBucket:   "test-deck-bucket",
			},
		},
		Email: struct {
			Provider   config.EmailProvider "mapstructure:\"provider\" yaml:\"provider\""
			Sender     malak.Email          "mapstructure:\"sender\" yaml:\"sender\""
			SenderName string               "mapstructure:\"sender_name\" yaml:\"sender_name\""
			SMTP       struct {
				Host     string "mapstructure:\"host\" yaml:\"host\""
				Port     int    "mapstructure:\"port\" yaml:\"port\""
				Username string "mapstructure:\"username\" yaml:\"username\""
				Password string "mapstructure:\"password\" yaml:\"password\""
				UseTLS   bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			} "mapstructure:\"smtp\" yaml:\"smtp\""
			Resend struct {
				APIKey        string "mapstructure:\"api_key\" yaml:\"api_key\""
				WebhookSecret string "mapstructure:\"webhook_secret\" yaml:\"webhook_secret\""
			} "mapstructure:\"resend\" yaml:\"resend\""
		}{
			Provider:   config.EmailProviderSmtp,
			Sender:     malak.Email("test@example.com"),
			SenderName: "Test Sender",
			SMTP: struct {
				Host     string "mapstructure:\"host\" yaml:\"host\""
				Port     int    "mapstructure:\"port\" yaml:\"port\""
				Username string "mapstructure:\"username\" yaml:\"username\""
				Password string "mapstructure:\"password\" yaml:\"password\""
				UseTLS   bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			}{
				Host:     "localhost",
				Port:     1025,
				Username: "test",
				Password: "test",
			},
		},
		Otel: struct {
			Endpoint  string "yaml:\"endpoint\" mapstructure:\"endpoint\""
			UseTLS    bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			Headers   string "yaml:\"headers\" mapstructure:\"headers\""
			IsEnabled bool   "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
		}{
			IsEnabled: false,
		},
		Billing: struct {
			Stripe struct {
				APIKey        string "yaml:\"api_key\" mapstructure:\"api_key\""
				WebhookSecret string "yaml:\"webhook_secret\" mapstructure:\"webhook_secret\""
			} "yaml:\"stripe\" mapstructure:\"stripe\""
			IsEnabled            bool   "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
			TrialDays            int64  "yaml:\"trial_days\" mapstructure:\"trial_days\""
			DefaultPlanReference string "yaml:\"default_plan_reference\" mapstructure:\"default_plan_reference\""
		}{
			Stripe: struct {
				APIKey        string "yaml:\"api_key\" mapstructure:\"api_key\""
				WebhookSecret string "yaml:\"webhook_secret\" mapstructure:\"webhook_secret\""
			}{
				WebhookSecret: webhookSecret,
			},
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
			Google: struct {
				ClientID     string   "yaml:\"client_id\" mapstructure:\"client_id\""
				ClientSecret string   "yaml:\"client_secret\" mapstructure:\"client_secret\""
				RedirectURI  string   "yaml:\"redirect_uri\" mapstructure:\"redirect_uri\""
				Scopes       []string "yaml:\"scopes\" mapstructure:\"scopes\""
				IsEnabled    bool     "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
			}{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				IsEnabled:    true,
			},
			JWT: struct {
				Key string "yaml:\"key\" mapstructure:\"key\""
			}{
				Key: "a907e75f80910f5dc5b8c677de1de611ffa80be9d7d9f9dd614c8c7846db1062",
			},
		},
	}
}

func getFetchCurrentUserData() []struct {
	name               string
	mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository)
	expectedStatusCode int
	addWorkspace       bool
} {

	return []struct {
		name               string
		mockFn             func(workspaceRepo *malak_mocks.MockWorkspaceRepository)
		expectedStatusCode int
		addWorkspace       bool
	}{
		{
			name: "could not list workspaces",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository) {
				workspaceRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("could not list workspaces"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "listed workspaces",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository) {
				workspaceRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Workspace{
						{
							Reference: "workspace_oops",
						},
						{
							Reference: "workspace_oopskfjk",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "listed workspaces with current workspace",
			mockFn: func(workspaceRepo *malak_mocks.MockWorkspaceRepository) {
				workspaceRepo.EXPECT().List(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]malak.Workspace{
						{
							Reference: "workspace_oops",
						},
						{
							Reference: "workspace_oopskfjk",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}
}

func TestAuthHandler_FetchCurrentUser(t *testing.T) {
	for _, v := range getFetchCurrentUserData() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)

			v.mockFn(workspaceRepo)

			a := &authHandler{
				cfg:           getConfig(),
				workspaceRepo: workspaceRepo,
			}

			var b = bytes.NewBuffer(nil)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", b)
			ctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
			req.Header.Add("Content-Type", "application/json")

			req = req.WithContext(writeUserToCtx(req.Context(), &malak.User{}))

			if v.addWorkspace {
				req = req.WithContext(writeWorkspaceToCtx(req.Context(), &malak.Workspace{}))
			}

			WrapMalakHTTPHandler(getLogger(t), a.fetchCurrentUser, getConfig(), "Auth.fetchCurrentUser").
				ServeHTTP(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	for _, v := range generateLoginTestTable() {

		t.Run(v.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			googleCfg := socialauth_mocks.NewMockSocialAuthProvider(controller)
			userRepo := malak_mocks.NewMockUserRepository(controller)

			jwtMock := mock_jwttoken.NewMockJWTokenManager(controller)

			v.mockFn(googleCfg, userRepo)

			a := &authHandler{
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

			WrapMalakHTTPHandler(getLogger(t), a.Login, getConfig(), "Auth.Login").
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
