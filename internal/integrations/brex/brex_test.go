package brex

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates client with config", func(t *testing.T) {
		cfg := config.Config{}
		cfg.Secrets.ClientTimeout = 30 * time.Second

		client, err := New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)

		brexClient, ok := client.(*brexClient)
		require.True(t, ok)
		require.Equal(t, 30*time.Second, brexClient.httpClient.Timeout)
	})
}

func TestBrexClient_Name(t *testing.T) {
	client := &brexClient{}
	require.Equal(t, malak.IntegrationProviderBrex, client.Name())
}

func TestBrexClient_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	token := os.Getenv("BREX_API_TOKEN")
	if token == "" {
		t.Skip("BREX_API_TOKEN not set")
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
			errMessage: "brex api request failed with status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.ClientTimeout = 10 * time.Second

			client, err := New(cfg)
			require.NoError(t, err)
			defer client.Close()

			charts, err := client.Ping(t.Context(), tt.token)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMessage)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, charts)
			}
		})
	}
}

func TestBrexClient_Close(t *testing.T) {
	client := &brexClient{
		httpClient: &http.Client{},
	}
	err := client.Close()
	require.NoError(t, err)
}

func TestBrexClient_buildRequest(t *testing.T) {
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
			endpoint: "/accounts/cash",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.ClientTimeout = 10 * time.Second

			client, err := New(cfg)
			require.NoError(t, err)
			brexClient := client.(*brexClient)

			req, span, err := brexClient.buildRequest(t.Context(), tt.token, tt.spanName, tt.endpoint)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMessage)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req)
			require.NotNil(t, span)

			require.Equal(t, fmt.Sprintf("Bearer %s", tt.token), req.Header.Get("Authorization"))
			require.Equal(t, "application/json", req.Header.Get("Accept"))
			require.Equal(t, baseURL+tt.endpoint, req.URL.String())
		})
	}
}

func TestBrexClient_Data(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	token := os.Getenv("BREX_API_TOKEN")
	if token == "" {
		t.Skip("BREX_API_TOKEN not set")
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
			errMessage: "brex api request failed with status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.ClientTimeout = 10 * time.Second

			client, err := New(cfg)
			require.NoError(t, err)
			defer client.Close()

			dataPoints, err := client.Data(t.Context(), tt.token, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMessage)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, dataPoints)

			// Verify we have both account balance and transaction data points
			hasAccountData := false
			hasTransactionData := false

			for _, dp := range dataPoints {
				switch dp.InternalName {
				case malak.IntegrationChartInternalNameTypeBrexAccount:
					hasAccountData = true
					require.Equal(t, malak.IntegrationDataPointTypeCurrency, dp.Data.DataPointType)
				case malak.IntegrationChartInternalNameTypeBrexAccountTransaction:
					hasTransactionData = true
					require.Equal(t, malak.IntegrationDataPointTypeOthers, dp.Data.DataPointType)
				}

				require.Equal(t, tt.opts.WorkspaceID, dp.Data.WorkspaceID)
				require.Equal(t, tt.opts.IntegrationID, dp.Data.WorkspaceIntegrationID)
				require.NotEmpty(t, dp.Data.Reference)
				require.NotEmpty(t, dp.Data.PointName)
			}

			require.True(t, hasAccountData, "should have account balance data")
			require.True(t, hasTransactionData, "should have transaction count data")
		})
	}
}
