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
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := config.Config{}
	cfg.Secrets.ClientTimeout = 30 * time.Second

	client, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	mercuryClient, ok := client.(*mercuryClient)
	require.True(t, ok)
	require.Equal(t, 30*time.Second, mercuryClient.httpClient.Timeout)
}

func TestMercuryClient_Name(t *testing.T) {
	client := &mercuryClient{}
	require.Equal(t, malak.IntegrationProviderMercury, client.Name())
}

func TestMercuryClient_Close(t *testing.T) {
	client := &mercuryClient{
		httpClient: &http.Client{},
	}
	err := client.Close()
	require.NoError(t, err)
}

func TestMercuryClient_buildRequest(t *testing.T) {
	cfg := config.Config{}
	cfg.Secrets.ClientTimeout = 10 * time.Second

	client, err := New(cfg)
	require.NoError(t, err)
	mercuryClient := client.(*mercuryClient)

	req, span, err := mercuryClient.buildRequest(context.Background(), "test-token", "test.span", "/accounts")
	require.NoError(t, err)
	require.NotNil(t, req)
	require.NotNil(t, span)

	require.Equal(t, fmt.Sprintf("Bearer %s", "test-token"), req.Header.Get("Authorization"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, baseURL+"/accounts", req.URL.String())
}

func TestMercuryClient_InvalidToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := config.Config{}
	cfg.Secrets.ClientTimeout = 10 * time.Second

	client, err := New(cfg)
	require.NoError(t, err)
	defer client.Close()

	// Test Ping with invalid token
	_, err = client.Ping(context.Background(), "invalid-token")
	require.Error(t, err)

	// Test Data with invalid token
	opts := &malak.IntegrationFetchDataOptions{
		IntegrationID:      uuid.New(),
		WorkspaceID:        uuid.New(),
		ReferenceGenerator: malak.NewReferenceGenerator(),
		LastFetchedAt:      time.Now(),
	}
	_, err = client.Data(context.Background(), "invalid-token", opts)
	require.Error(t, err)
}

func TestMercuryClient_ValidToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	token := os.Getenv("MERCURY_API_TOKEN")
	if token == "" {
		t.Skip("MERCURY_API_TOKEN not set")
	}

	cfg := config.Config{}
	cfg.Secrets.ClientTimeout = 10 * time.Second

	client, err := New(cfg)
	require.NoError(t, err)
	defer client.Close()

	// Test Ping with valid token
	_, err = client.Ping(context.Background(), malak.AccessToken(token))
	require.NoError(t, err)

	// Test Data with valid token
	opts := &malak.IntegrationFetchDataOptions{
		IntegrationID:      uuid.New(),
		WorkspaceID:        uuid.New(),
		ReferenceGenerator: malak.NewReferenceGenerator(),
		LastFetchedAt:      time.Now(),
	}
	dataPoints, err := client.Data(context.Background(), malak.AccessToken(token), opts)
	if err != nil {
		t.Skip("MERCURY_API_TOKEN does not have access to any accounts")
	}

	require.NotEmpty(t, dataPoints)
	for _, dp := range dataPoints {
		require.NotEmpty(t, dp.Data.Reference)
		require.NotEmpty(t, dp.Data.PointName)
		require.Equal(t, opts.WorkspaceID, dp.Data.WorkspaceID)
		require.Equal(t, opts.IntegrationID, dp.Data.WorkspaceIntegrationID)
	}
}
