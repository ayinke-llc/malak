package secretsmanager

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func TestCreateAndGetSecret(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:1.4.0")
	defer func() {
		require.NoError(t, testcontainers.TerminateContainer(localstackContainer))
	}()
	require.NoError(t, err)

	secretKey := "test-secret"
	secretValue := "super-secret-value"

	cfg := config.Config{}
	cfg.Integration.SecretsManager.Region = "us-east-1"
	cfg.Integration.SecretsManager.AccessSecret = "value"
	cfg.Integration.SecretsManager.AccessKey = "key"

	client, err := New(cfg)
	require.NoError(t, err)

	_, err = client.Create(ctx, &secret.CreateSecretOptions{
		Value:       secretValue,
		WorkspaceID: uuid.New(),
	})
	require.NoError(t, err)

	val, err := client.Get(ctx, secretKey)
	require.NoError(t, err)
	require.Equal(t, secretValue, val)
}

func TestCreateAndGetNonexistenstSecret(t *testing.T) {
	ctx := context.Background()

	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:1.4.0")
	defer func() {
		require.NoError(t, testcontainers.TerminateContainer(localstackContainer))
	}()
	require.NoError(t, err)

	cfg := config.Config{}
	cfg.Integration.SecretsManager.Region = "us-east-1"
	cfg.Integration.SecretsManager.AccessSecret = "value"
	cfg.Integration.SecretsManager.AccessKey = "key"

	client, err := New(cfg)
	require.NoError(t, err)

	_, err = client.Get(ctx, "unknwon-key")
	require.Error(t, err)
}
