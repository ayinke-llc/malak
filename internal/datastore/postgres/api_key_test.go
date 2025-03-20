package postgres

import (
	"database/sql"
	"fmt"
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
		testKey := &malak.APIKey{
			ID:          uuid.New(),
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
			CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       "test_value_immediate",
			KeyName:     "immediate_key",
		}
		err := repo.Create(t.Context(), testKey)
		require.NoError(t, err)

		err = repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         testKey,
			RevocationType: malak.RevocationTypeImmediate,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", testKey.ID).Scan(t.Context())
		require.Error(t, err)
		require.Equal(t, err, sql.ErrNoRows)
	})

	t.Run("day revocation", func(t *testing.T) {
		testKey := &malak.APIKey{
			ID:          uuid.New(),
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
			CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       "test_value_day",
			KeyName:     "day_key",
		}
		err := repo.Create(t.Context(), testKey)
		require.NoError(t, err)

		err = repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         testKey,
			RevocationType: malak.RevocationTypeDay,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", testKey.ID).Scan(t.Context())
		require.NoError(t, err)

		now := time.Now()
		expectedExpiry := now.AddDate(0, 0, 1)
		expectedExpiry = time.Date(expectedExpiry.Year(), expectedExpiry.Month(), expectedExpiry.Day(), 23, 0, 0, 0, now.Location())
		require.Equal(t, expectedExpiry.Format("2006-01-02"), revokedAPIKey.ExpiresAt.Format("2006-01-02"))
		require.Nil(t, revokedAPIKey.DeletedAt)
	})

	t.Run("week revocation", func(t *testing.T) {
		testKey := &malak.APIKey{
			ID:          uuid.New(),
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
			CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       "test_value_week",
			KeyName:     "week_key",
		}
		err := repo.Create(t.Context(), testKey)
		require.NoError(t, err)

		err = repo.Revoke(t.Context(), malak.RevokeAPIKeyOptions{
			APIKey:         testKey,
			RevocationType: malak.RevocationTypeWeek,
		})
		require.NoError(t, err)

		var revokedAPIKey malak.APIKey
		err = client.NewSelect().Model(&revokedAPIKey).Where("id = ?", testKey.ID).Scan(t.Context())
		require.NoError(t, err)

		now := time.Now()
		expectedExpiry := now.AddDate(0, 0, 7)
		expectedExpiry = time.Date(expectedExpiry.Year(), expectedExpiry.Month(), expectedExpiry.Day(), 23, 0, 0, 0, now.Location())
		require.Equal(t, expectedExpiry.Format("2006-01-02"), revokedAPIKey.ExpiresAt.Format("2006-01-02"))
		require.Nil(t, revokedAPIKey.DeletedAt)
	})
}

func TestAPIKey_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewAPIKeyRepository(client)
	workspaceID := uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0")
	createdBy := uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6")

	// Create multiple API keys
	apiKeys := []*malak.APIKey{
		{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			CreatedBy:   createdBy,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       "test_value_1",
			KeyName:     "key1",
		},
		{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			CreatedBy:   createdBy,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       "test_value_2",
			KeyName:     "key2",
		},
	}

	for _, key := range apiKeys {
		err := repo.Create(t.Context(), key)
		require.NoError(t, err)
	}

	// Test listing API keys
	listed, err := repo.List(t.Context(), workspaceID)
	require.NoError(t, err)
	require.Len(t, listed, 2)

	// Verify the listed keys match what we created
	for i, key := range apiKeys {
		require.Equal(t, key.ID, listed[i].ID)
		require.Equal(t, key.WorkspaceID, listed[i].WorkspaceID)
		require.Equal(t, key.CreatedBy, listed[i].CreatedBy)
		require.Equal(t, key.Reference, listed[i].Reference)
		require.Equal(t, key.Value, listed[i].Value)
		require.Equal(t, key.KeyName, listed[i].KeyName)
	}

	// Test listing for a different workspace returns no results
	differentWorkspaceID := uuid.New()
	emptyList, err := repo.List(t.Context(), differentWorkspaceID)
	require.NoError(t, err)
	require.Empty(t, emptyList)
}

func TestAPIKey_Fetch(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewAPIKeyRepository(client)
	workspaceID := uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0")
	createdBy := uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6")

	apiKey := &malak.APIKey{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		CreatedBy:   createdBy,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
		Value:       "test_value",
		KeyName:     "test_key",
	}

	err := repo.Create(t.Context(), apiKey)
	require.NoError(t, err)

	t.Run("fetch existing key", func(t *testing.T) {
		fetched, err := repo.Fetch(t.Context(), malak.FetchAPIKeyOptions{
			WorkspaceID: workspaceID,
			Reference:   apiKey.Reference,
		})
		require.NoError(t, err)
		require.Equal(t, apiKey.ID, fetched.ID)
		require.Equal(t, apiKey.WorkspaceID, fetched.WorkspaceID)
		require.Equal(t, apiKey.CreatedBy, fetched.CreatedBy)
		require.Equal(t, apiKey.Reference, fetched.Reference)
		require.Equal(t, apiKey.Value, fetched.Value)
		require.Equal(t, apiKey.KeyName, fetched.KeyName)
	})

	t.Run("fetch non-existent key", func(t *testing.T) {
		_, err := repo.Fetch(t.Context(), malak.FetchAPIKeyOptions{
			WorkspaceID: workspaceID,
			Reference:   "non_existent_reference",
		})
		require.ErrorIs(t, err, malak.ErrAPIKeyNotFound)
	})

	t.Run("fetch from wrong workspace", func(t *testing.T) {
		_, err := repo.Fetch(t.Context(), malak.FetchAPIKeyOptions{
			WorkspaceID: uuid.New(),
			Reference:   apiKey.Reference,
		})
		require.ErrorIs(t, err, malak.ErrAPIKeyNotFound)
	})
}

func TestAPIKey_MaxLimit(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewAPIKeyRepository(client)
	workspaceID := uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0")
	createdBy := uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6")

	// Create 15 API keys (max limit)
	for i := 0; i < 15; i++ {
		apiKey := &malak.APIKey{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			CreatedBy:   createdBy,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
			Value:       fmt.Sprintf("test_value_%d", i),
			KeyName:     fmt.Sprintf("key_%d", i),
		}
		err := repo.Create(t.Context(), apiKey)
		require.NoError(t, err)
	}

	// Try to create one more key
	extraKey := &malak.APIKey{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		CreatedBy:   createdBy,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeApiKey),
		Value:       "extra_key",
		KeyName:     "extra",
	}
	err := repo.Create(t.Context(), extraKey)
	require.ErrorIs(t, err, malak.ErrAPIKeyMaxLimit)

	// Verify we still have only 15 keys
	keys, err := repo.List(t.Context(), workspaceID)
	require.NoError(t, err)
	require.Len(t, keys, 15)
}
