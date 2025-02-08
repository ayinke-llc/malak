package mercury

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/google/uuid"
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

			_, err = client.Ping(context.Background(), tt.token)
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

func TestMercuryClient_buildRequest(t *testing.T) {
	tests := []struct {
		name       string
		token      malak.AccessToken
		spanName   string
		endpoint   string
		wantErr    bool
		errMessage string
	}{
		{
			name:     "successful request build",
			token:    "test-token",
			spanName: "test.span",
			endpoint: "/accounts",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.ClientTimeout = 10 * time.Second

			client, err := New(cfg)
			require.NoError(t, err)
			mercuryClient := client.(*mercuryClient)

			req, span, err := mercuryClient.buildRequest(context.Background(), tt.token, tt.spanName, tt.endpoint)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req)
			require.NotNil(t, span)

			assert.Equal(t, fmt.Sprintf("Bearer %s", tt.token), req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, baseURL+tt.endpoint, req.URL.String())
		})
	}
}

func TestMercuryClient_Data(t *testing.T) {
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
		opts       *malak.IntegrationFetchDataOptions
		wantErr    bool
		errMessage string
	}{
		{
			name:  "successful data fetch",
			token: malak.AccessToken(token),
			opts: &malak.IntegrationFetchDataOptions{
				IntegrationID:      uuid.New(),
				WorkspaceID:        uuid.New(),
				ReferenceGenerator: malak.NewReferenceGenerator(),
				LastFetchedAt:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:  "invalid token",
			token: "invalid-token",
			opts: &malak.IntegrationFetchDataOptions{
				IntegrationID:      uuid.New(),
				WorkspaceID:        uuid.New(),
				ReferenceGenerator: malak.NewReferenceGenerator(),
				LastFetchedAt:      time.Now(),
			},
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

			dataPoints, err := client.Data(context.Background(), tt.token, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, dataPoints)

			// Verify we have both account balance and transaction data points
			hasAccountData := false
			hasTransactionData := false

			for _, dp := range dataPoints {
				switch dp.InternalName {
				case malak.IntegrationChartInternalNameTypeMercuryAccount:
					hasAccountData = true
					assert.Equal(t, malak.IntegrationDataPointTypeCurrency, dp.Data.DataPointType)
				case malak.IntegrationChartInternalNameTypeMercuryAccountTransaction:
					hasTransactionData = true
					assert.Equal(t, malak.IntegrationDataPointTypeOthers, dp.Data.DataPointType)
				}

				assert.Equal(t, tt.opts.WorkspaceID, dp.Data.WorkspaceID)
				assert.Equal(t, tt.opts.IntegrationID, dp.Data.WorkspaceIntegrationID)
				assert.NotEmpty(t, dp.Data.Reference)
				assert.NotEmpty(t, dp.Data.PointName)
			}

			assert.True(t, hasAccountData, "should have account balance data")
			assert.True(t, hasTransactionData, "should have transaction count data")
		})
	}
}
