package secretsmanager

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type testSuite struct {
	client    secret.SecretClient
	container *localstack.LocalStackContainer
	ctx       context.Context
}

func setupTest(t *testing.T) *testSuite {
	ctx := context.Background()

	container, err := localstack.Run(ctx, "localstack/localstack:4.1.0",
		testcontainers.WithEnv(map[string]string{
			"SERVICES":  "secretsmanager",
			"DEBUG":     "1",
			"EDGE_PORT": "4566",
		}),
	)
	require.NoError(t, err)

	endpoint, err := container.Endpoint(ctx, "secretsmanager")
	require.NoError(t, err)

	// Configure AWS client to use localstack
	cfg := config.Config{}
	cfg.Integration.SecretsManager.Region = "us-east-1"
	cfg.Integration.SecretsManager.AccessKey = "test"
	cfg.Integration.SecretsManager.AccessSecret = "test"
	cfg.Integration.SecretsManager.Endpoint = endpoint

	client, err := New(cfg)
	require.NoError(t, err)

	return &testSuite{
		client:    client,
		container: container,
		ctx:       ctx,
	}
}

func (ts *testSuite) tearDown(t *testing.T) {
	if ts.container != nil {
		err := ts.container.Terminate(ts.ctx)
		require.NoError(t, err)
	}
}

func TestSecretsManager(t *testing.T) {
	t.Run("create and get secret", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.tearDown(t)

		workspaceID := uuid.New()
		secretValue := "super-secret-value"

		opts := &secret.CreateSecretOptions{
			Value:       secretValue,
			WorkspaceID: workspaceID,
		}

		// Create secret
		key, err := ts.client.Create(ts.ctx, opts)
		require.NoError(t, err)
		assert.NotEmpty(t, key)

		// Get secret
		val, err := ts.client.Get(ts.ctx, key)
		require.NoError(t, err)
		assert.Equal(t, secretValue, val)
	})

	t.Run("get nonexistent secret", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.tearDown(t)

		_, err := ts.client.Get(ts.ctx, "unknown-key")
		require.Error(t, err)
		assert.ErrorIs(t, err, &types.ResourceNotFoundException{})
	})

	t.Run("create duplicate secret", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.tearDown(t)

		workspaceID := uuid.New()
		opts := &secret.CreateSecretOptions{
			Value:       "test-value",
			WorkspaceID: workspaceID,
		}

		// Create first secret
		key, err := ts.client.Create(ts.ctx, opts)
		require.NoError(t, err)

		// Try to create another secret with same key
		_, err = ts.client.Create(ts.ctx, opts)
		require.Error(t, err)
		assert.ErrorIs(t, err, &types.ResourceExistsException{})

		// Verify original secret still exists
		val, err := ts.client.Get(ts.ctx, key)
		require.NoError(t, err)
		assert.Equal(t, opts.Value, val)
	})

	t.Run("create secret with empty value", func(t *testing.T) {
		ts := setupTest(t)
		defer ts.tearDown(t)

		opts := &secret.CreateSecretOptions{
			Value:       "",
			WorkspaceID: uuid.New(),
		}

		key, err := ts.client.Create(ts.ctx, opts)
		require.NoError(t, err)

		val, err := ts.client.Get(ts.ctx, key)
		require.NoError(t, err)
		assert.Empty(t, val)
	})
}
