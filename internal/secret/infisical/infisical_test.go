package infisical

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInfisicalClientInitialization(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		expectError bool
	}{
		{
			name:        "empty config",
			config:      config.Config{},
			expectError: true,
		},
		{
			name: "missing environment",
			config: func() config.Config {
				cfg := config.Config{}
				cfg.Integration.Provider = secret.SecretProvider("infisical")
				cfg.Integration.Infisical.ClientID = "test-id"
				cfg.Integration.Infisical.ClientSecret = "test-secret"
				cfg.Integration.Infisical.SiteURL = "http://localhost:8080"
				cfg.Integration.Infisical.ProjectID = "test-project"
				return cfg
			}(),
			expectError: true,
		},
		{
			name: "missing project ID",
			config: func() config.Config {
				cfg := config.Config{}
				cfg.Integration.Provider = secret.SecretProvider("infisical")
				cfg.Integration.Infisical.ClientID = "test-id"
				cfg.Integration.Infisical.ClientSecret = "test-secret"
				cfg.Integration.Infisical.SiteURL = "http://localhost:8080"
				cfg.Integration.Infisical.Environment = "dev"
				return cfg
			}(),
			expectError: true,
		},
		{
			name: "missing client ID",
			config: func() config.Config {
				cfg := config.Config{}
				cfg.Integration.Provider = secret.SecretProvider("infisical")
				cfg.Integration.Infisical.ClientSecret = "test-secret"
				cfg.Integration.Infisical.SiteURL = "http://localhost:8080"
				cfg.Integration.Infisical.ProjectID = "test-project"
				cfg.Integration.Infisical.Environment = "dev"
				return cfg
			}(),
			expectError: true,
		},
		{
			name: "missing client secret",
			config: func() config.Config {
				cfg := config.Config{}
				cfg.Integration.Provider = secret.SecretProvider("infisical")
				cfg.Integration.Infisical.ClientID = "test-id"
				cfg.Integration.Infisical.SiteURL = "http://localhost:8080"
				cfg.Integration.Infisical.ProjectID = "test-project"
				cfg.Integration.Infisical.Environment = "dev"
				return cfg
			}(),
			expectError: true,
		},
		{
			name: "missing site URL",
			config: func() config.Config {
				cfg := config.Config{}
				cfg.Integration.Provider = secret.SecretProvider("infisical")
				cfg.Integration.Infisical.ClientID = "test-id"
				cfg.Integration.Infisical.ClientSecret = "test-secret"
				cfg.Integration.Infisical.ProjectID = "test-project"
				cfg.Integration.Infisical.Environment = "dev"
				return cfg
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestInfisicalClient(t *testing.T) {
	ctx := context.Background()

	// Set up the container request
	req := testcontainers.ContainerRequest{
		Image:        "infisical/infisical:7090eea",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Server listening on port"),
			wait.ForListeningPort("8080/tcp"),
		).WithStartupTimeout(30 * time.Second),
		Env: map[string]string{
			"ENCRYPTION_KEY": "test-encryption-key",
		},
	}

	// Start the container
	infisicalContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		if err := infisicalContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get the container's host and port
	host, err := infisicalContainer.Host(ctx)
	require.NoError(t, err)
	port, err := infisicalContainer.MappedPort(ctx, "8080")
	require.NoError(t, err)

	// Create test configuration
	cfg := config.Config{}
	cfg.Integration.Provider = secret.SecretProvider("infisical")
	cfg.Integration.Infisical.SiteURL = fmt.Sprintf("http://%s:%s", host, port.Port())
	cfg.Integration.Infisical.Environment = "dev"
	cfg.Integration.Infisical.ProjectID = "test-project"
	cfg.Integration.Infisical.ClientID = "test-client-id"
	cfg.Integration.Infisical.ClientSecret = "test-client-secret"

	// Initialize client
	client, err := New(cfg)
	require.NoError(t, err)
	defer client.Close()

	t.Run("create and get secret", func(t *testing.T) {
		workspaceID := uuid.New()
		value := "test-value"

		createdKey, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           value,
			WorkspaceID:     workspaceID,
			IntegrationName: malak.IntegrationProvider("infisical"),
		})
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s/infisical", workspaceID.String()), createdKey)

		// Test getting the secret
		retrievedValue, err := client.Get(ctx, createdKey)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("get non-existent secret", func(t *testing.T) {
		_, err := client.Get(ctx, "non-existent-key")
		assert.Error(t, err)
	})

	t.Run("create secret with empty value", func(t *testing.T) {
		workspaceID := uuid.New()
		_, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           "",
			WorkspaceID:     workspaceID,
			IntegrationName: malak.IntegrationProvider("infisical"),
		})
		assert.NoError(t, err)
	})

	t.Run("create secret with special characters", func(t *testing.T) {
		workspaceID := uuid.New()
		value := "test@value#with$special&chars*"

		createdKey, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           value,
			WorkspaceID:     workspaceID,
			IntegrationName: malak.IntegrationProvider("infisical"),
		})
		assert.NoError(t, err)

		retrievedValue, err := client.Get(ctx, createdKey)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("create secret with very long value", func(t *testing.T) {
		workspaceID := uuid.New()
		value := string(make([]byte, 4096)) // 4KB of data

		createdKey, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           value,
			WorkspaceID:     workspaceID,
			IntegrationName: malak.IntegrationProvider("infisical"),
		})
		assert.NoError(t, err)

		retrievedValue, err := client.Get(ctx, createdKey)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("create secret with nil workspace ID", func(t *testing.T) {
		_, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           "test",
			IntegrationName: malak.IntegrationProvider("infisical"),
		})
		assert.Error(t, err)
	})

	t.Run("create secret with invalid integration name", func(t *testing.T) {
		workspaceID := uuid.New()
		_, err := client.Create(ctx, &secret.CreateSecretOptions{
			Value:           "test",
			WorkspaceID:     workspaceID,
			IntegrationName: malak.IntegrationProvider("invalid"),
		})
		assert.Error(t, err)
	})
}
