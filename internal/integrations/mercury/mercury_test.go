package mercury

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates client with config", func(t *testing.T) {
		cfg := config.Config{}
		cfg.Secrets.ClientTimeout = 30 * time.Second

		client, err := New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)

		mercuryClient, ok := client.(*mercuryClient)
		require.True(t, ok)
		assert.Equal(t, 30*time.Second, mercuryClient.httpClient.Timeout)
	})
}

func TestMercuryClient_Name(t *testing.T) {
	client := &mercuryClient{}
	assert.Equal(t, malak.IntegrationProviderMercury, client.Name())
}

func TestMercuryClient_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	token := os.Getenv("MERCURY_API_TOKEN")
	if token == "" {
		t.Skip("MERCURY_API_TOKEN not set")
	}

	tests := []struct {
		name       string
		token      malak.AccessToken
		wantErr    bool
		errMessage string
	}{
		{
			name:    "successful ping with valid token",
			token:   malak.AccessToken(token),
			wantErr: false,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantErr:    true,
			errMessage: "invalid api key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.ClientTimeout = 10 * time.Second

			client, err := New(cfg)
			require.NoError(t, err)
			defer client.Close()

			err = client.Ping(context.Background(), tt.token)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMercuryClient_Close(t *testing.T) {
	client := &mercuryClient{
		httpClient: &http.Client{},
	}
	err := client.Close()
	require.NoError(t, err)
}
