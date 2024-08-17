package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ayinke-llc/malak/config"
	socialauth_mocks "github.com/ayinke-llc/malak/internal/pkg/socialauth/mocks"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

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
	}
}

func TestAuthHandler_Login(t *testing.T) {

	tt := []struct {
		name               string
		mockFn             func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository)
		expectedStatusCode int
	}{
		{
			name: "no code to exchange provided",
			mockFn: func(googleMock *socialauth_mocks.MockSocialAuthProvider, userRepo *malak_mocks.MockUserRepository) {
				// googleMock.EXPECT().Validate()
			},
		},
	}

	for _, v := range tt {

		t.Run(v.name, func(t *testing.T) {

			logrus.SetOutput(io.Discard)

			logger := logrus.WithField("test", true)

			controller := gomock.NewController(t)
			defer controller.Finish()

			googleCfg := socialauth_mocks.NewMockSocialAuthProvider(controller)
			userRepo := malak_mocks.NewMockUserRepository(controller)

			a := &authHandler{
				logger:    logger,
				cfg:       getConfig(),
				googleCfg: googleCfg,
				userRepo:  userRepo,
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			a.Login(rr, req)
		})
	}
}
