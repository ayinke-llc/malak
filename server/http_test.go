package server

import (
	"testing"
	"time"

	"github.com/adelowo/gulter"
	"github.com/ayinke-llc/malak/internal/integrations"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"go.uber.org/mock/gomock"
)

// This does nothing really
// It is only here to verify that the routing works correctly
// e.g middlewares are correctly set, paths are correctly named.
// Will help catch issues with paths like `/updates/{references` that is
// missing an ending brace or wrongly placed middlewares and others
func TestServer_New(t *testing.T) {

	t.Run("without swagger", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()

		cfg := getConfig()

		srv, closeFn := New(getLogger(t), cfg, &bun.DB{},
			jwttoken.New(cfg), socialauth.NewGoogle(cfg),
			malak_mocks.NewMockUserRepository(controller),
			malak_mocks.NewMockWorkspaceRepository(controller),
			malak_mocks.NewMockPlanRepository(controller),
			malak_mocks.NewMockContactRepository(controller),
			malak_mocks.NewMockUpdateRepository(controller),
			malak_mocks.NewMockContactListRepository(controller),
			malak_mocks.NewMockDeckRepository(controller),
			malak_mocks.NewMockContactShareRepository(controller),
			malak_mocks.NewMockPreferenceRepository(controller),
			malak_mocks.NewMockIntegrationRepository(controller),
			&httplimit.Middleware{},
			&gulter.Gulter{},
			malak_mocks.NewMockQueueHandler(controller),
			malak_mocks.NewMockCache(controller),
			malak_mocks.NewMockClient(controller),
			integrations.NewManager())

		closeFn()

		go func() {
			_ = srv.ListenAndServe()
		}()

		time.Sleep(time.Second * 2)

		require.NoError(t, srv.Close())
	})

	t.Run("with swagger enabled", func(t *testing.T) {

		controller := gomock.NewController(t)
		defer controller.Finish()

		cfg := getConfig()

		cfg.HTTP.Swagger.UIEnabled = true
		cfg.HTTP.Swagger.Port = 9990

		srv, closeFn := New(getLogger(t), cfg, &bun.DB{},
			jwttoken.New(cfg), socialauth.NewGoogle(cfg),
			malak_mocks.NewMockUserRepository(controller),
			malak_mocks.NewMockWorkspaceRepository(controller),
			malak_mocks.NewMockPlanRepository(controller),
			malak_mocks.NewMockContactRepository(controller),
			malak_mocks.NewMockUpdateRepository(controller),
			malak_mocks.NewMockContactListRepository(controller),
			malak_mocks.NewMockDeckRepository(controller),
			malak_mocks.NewMockContactShareRepository(controller),
			malak_mocks.NewMockPreferenceRepository(controller),
			malak_mocks.NewMockIntegrationRepository(controller),
			&httplimit.Middleware{},
			&gulter.Gulter{},
			malak_mocks.NewMockQueueHandler(controller),
			malak_mocks.NewMockCache(controller),
			malak_mocks.NewMockClient(controller),
			integrations.NewManager())

		closeFn()

		go func() {
			_ = srv.ListenAndServe()
		}()

		time.Sleep(time.Second * 2)

		require.NoError(t, srv.Close())
	})
}
