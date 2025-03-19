package postgres

import (
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAPIKey_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewAPIKeyRepository(client)

	apiKey := &malak.APIKey{
		ID:          uuid.New(),
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"), // First workspace
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"), // lanre@test.com
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
		Value:       "test_value",
		KeyName:     "oopsf",
	}

	err := repo.Create(t.Context(), apiKey)
	require.NoError(t, err)

	var createdAPIKey malak.APIKey
	err = client.NewSelect().Model(&createdAPIKey).Where("id = ?", apiKey.ID).Scan(t.Context())
	require.NoError(t, err)
	require.Equal(t, apiKey.ID, createdAPIKey.ID)
	require.Equal(t, apiKey.WorkspaceID, createdAPIKey.WorkspaceID)
	require.Equal(t, apiKey.CreatedBy, createdAPIKey.CreatedBy)
	require.Equal(t, apiKey.Reference, createdAPIKey.Reference)
	require.Equal(t, apiKey.Value, createdAPIKey.Value)
}

func TestAPIKey_Revoke(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewAPIKeyRepository(client)

	apiKey := &malak.APIKey{
		ID:          uuid.New(),
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"), // First workspace
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"), // lanre@test.com
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
		Value:       "test_value",
		KeyName:     "oopsf",
	}

	err := repo.Create(t.Context(), apiKey)
	require.NoError(t, err)

	t.Run("immediate revocation", func(t *testing.T) {
		err := repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         apiKey,
			RevocationType: malak.RevocationTypeImmediate,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", apiKey.ID).Scan(t.Context())
		require.NoError(t, err)

		now := time.Now()
		expectedExpiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		require.Equal(t, expectedExpiry.Unix(), revokedAPIKey.ExpiresAt.Unix())
	})

	t.Run("day revocation", func(t *testing.T) {
		err := repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         apiKey,
			RevocationType: malak.RevocationTypeDay,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", apiKey.ID).Scan(t.Context())
		require.NoError(t, err)

		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)
		expectedExpiry := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, tomorrow.Location())
		require.Equal(t, expectedExpiry.Unix(), revokedAPIKey.ExpiresAt.Unix())
	})

	t.Run("week revocation", func(t *testing.T) {
		err := repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         apiKey,
			RevocationType: malak.RevocationTypeWeek,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", apiKey.ID).Scan(t.Context())
		require.NoError(t, err)

		now := time.Now()
		weekFromNow := now.AddDate(0, 0, 7)
		expectedExpiry := time.Date(weekFromNow.Year(), weekFromNow.Month(), weekFromNow.Day(), 23, 59, 59, 0, weekFromNow.Location())
		require.Equal(t, expectedExpiry.Unix(), revokedAPIKey.ExpiresAt.Unix())
	})
}
