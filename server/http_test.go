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
	"go.uber.org/zap"
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
			integrations.NewManager(),
			malak_mocks.NewMockSecretClient(controller))

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
			integrations.NewManager(),
			malak_mocks.NewMockSecretClient(controller))

		closeFn()

		go func() {
			_ = srv.ListenAndServe()
		}()

		time.Sleep(time.Second * 2)

		require.NoError(t, srv.Close())
	})
}

func TestNew(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userRepo := malak_mocks.NewMockUserRepository(controller)
	workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
	planRepo := malak_mocks.NewMockPlanRepository(controller)
	contactRepo := malak_mocks.NewMockContactRepository(controller)
	updateRepo := malak_mocks.NewMockUpdateRepository(controller)
	contactListRepo := malak_mocks.NewMockContactListRepository(controller)
	deckRepo := malak_mocks.NewMockDeckRepository(controller)
	contactShareRepo := malak_mocks.NewMockContactShareRepository(controller)
	preferenceRepo := malak_mocks.NewMockPreferenceRepository(controller)
	integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
	queueRepo := malak_mocks.NewMockQueueHandler(controller)
	cacheRepo := malak_mocks.NewMockCache(controller)
	billingClient := malak_mocks.NewMockClient(controller)
	secretsClient := malak_mocks.NewMockSecretClient(controller)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getConfig()

	db := &bun.DB{}

	srv, closeFn := New(logger, cfg, db, jwttoken.New(cfg), socialauth.NewGoogle(cfg),
		userRepo, workspaceRepo, planRepo, contactRepo, updateRepo,
		contactListRepo, deckRepo, contactShareRepo, preferenceRepo,
		integrationRepo, &httplimit.Middleware{}, &gulter.Gulter{}, queueRepo, cacheRepo,
		billingClient, integrations.NewManager(), secretsClient)

	require.NotNil(t, srv)
	require.NotNil(t, closeFn)
}

func TestNewWithInvalidConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userRepo := malak_mocks.NewMockUserRepository(controller)
	workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
	planRepo := malak_mocks.NewMockPlanRepository(controller)
	contactRepo := malak_mocks.NewMockContactRepository(controller)
	updateRepo := malak_mocks.NewMockUpdateRepository(controller)
	contactListRepo := malak_mocks.NewMockContactListRepository(controller)
	deckRepo := malak_mocks.NewMockDeckRepository(controller)
	contactShareRepo := malak_mocks.NewMockContactShareRepository(controller)
	preferenceRepo := malak_mocks.NewMockPreferenceRepository(controller)
	integrationRepo := malak_mocks.NewMockIntegrationRepository(controller)
	queueRepo := malak_mocks.NewMockQueueHandler(controller)
	cacheRepo := malak_mocks.NewMockCache(controller)
	billingClient := malak_mocks.NewMockClient(controller)
	secretsClient := malak_mocks.NewMockSecretClient(controller)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	cfg := getConfig()
	cfg.HTTP.Port = 0

	db := &bun.DB{}

	srv, closeFn := New(logger, cfg, db, jwttoken.New(cfg), socialauth.NewGoogle(cfg),
		userRepo, workspaceRepo, planRepo, contactRepo, updateRepo,
		contactListRepo, deckRepo, contactShareRepo, preferenceRepo,
		integrationRepo, &httplimit.Middleware{}, &gulter.Gulter{}, queueRepo, cacheRepo,
		billingClient, integrations.NewManager(), secretsClient)

	require.Nil(t, srv)
	require.Nil(t, closeFn)
}
