package vault

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/vault"
)

func TestCreateAndGetSecret(t *testing.T) {
	ctx := context.Background()

	vaultContainer, err := vault.Run(ctx,
		"hashicorp/vault:1.18.4",
		testcontainers.WithEnv(map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID": "dev-token",
		}))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, vaultContainer.Terminate(ctx))
	}()

	vaultAddr, err := vaultContainer.HttpHostAddress(ctx)
	require.NoError(t, err)

	cfg := config.Config{}
	cfg.Secrets.Vault.Address = vaultAddr
	cfg.Secrets.Vault.Token = "dev-token"
	cfg.Secrets.Vault.Path = "secret"

	client, err := New(cfg)
	require.NoError(t, err)
	defer client.Close()

	secretValue := "super-secret-value"
	workspaceID := uuid.New()

	key, err := client.Create(ctx, &secret.CreateSecretOptions{
		Value:       secretValue,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)

	val, err := client.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, secretValue, val)
}

func TestCreateAndGetNonexistentSecret(t *testing.T) {
	ctx := context.Background()

	vaultContainer, err := vault.Run(ctx,
		"hashicorp/vault:1.18.4",
		testcontainers.WithEnv(map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID": "dev-token",
		}))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, vaultContainer.Terminate(ctx))
	}()

	vaultAddr, err := vaultContainer.HttpHostAddress(ctx)
	require.NoError(t, err)

	cfg := config.Config{}
	cfg.Secrets.Vault.Address = vaultAddr
	cfg.Secrets.Vault.Token = "dev-token"
	cfg.Secrets.Vault.Path = "secret"

	client, err := New(cfg)
	require.NoError(t, err)
	defer client.Close()

	_, err = client.Get(ctx, "unknown-key")
	require.Error(t, err)
}

func TestNewVaultClientValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		{
			name:    "missing path",
			cfg:     config.Config{},
			wantErr: true,
		},
		{
			name:    "missing token",
			cfg:     config.Config{},
			wantErr: true,
		},
		{
			name:    "missing address",
			cfg:     config.Config{},
			wantErr: true,
		},
		{
			name:    "valid config",
			cfg:     config.Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "missing path":
				tt.cfg.Secrets.Vault.Token = "token"
				tt.cfg.Secrets.Vault.Address = "addr"
			case "missing token":
				tt.cfg.Secrets.Vault.Path = "path"
				tt.cfg.Secrets.Vault.Address = "addr"
			case "missing address":
				tt.cfg.Secrets.Vault.Path = "path"
				tt.cfg.Secrets.Vault.Token = "token"
			case "valid config":
				tt.cfg.Secrets.Vault.Path = "path"
				tt.cfg.Secrets.Vault.Token = "token"
				tt.cfg.Secrets.Vault.Address = "addr"
			}

			_, err := New(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
