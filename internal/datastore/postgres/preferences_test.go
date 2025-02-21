package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPreferences_Get(t *testing.T) {
	// create and fetch
	t.Run("create and fetch", func(t *testing.T) {

		client, teardownFunc := setupDatabase(t)
		defer teardownFunc()

		prefRepo := NewPreferenceRepository(client)

		repo := NewWorkspaceRepository(client)

		userRepo := NewUserRepository(client)

		planRepo := NewPlanRepository(client)

		// user from the fixtures
		user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
			Email: "lanre@test.com",
		})
		require.NoError(t, err)

		plan, err := planRepo.Get(t.Context(), &malak.FetchPlanOptions{
			Reference: "prod_QmtErtydaJZymT",
		})
		require.NoError(t, err)

		opts := &malak.CreateWorkspaceOptions{
			User:      user,
			Workspace: malak.NewWorkspace("oops", user, plan, malak.GenerateReference(malak.EntityTypeWorkspace)),
		}

		require.NoError(t, repo.Create(t.Context(), opts))

		pref, err := prefRepo.Get(t.Context(), opts.Workspace)
		require.NoError(t, err)

		require.True(t, pref.Communication.EnableMarketing)
		require.True(t, pref.Communication.EnableProductUpdates)
	})

	t.Run("no exists", func(t *testing.T) {

		client, teardownFunc := setupDatabase(t)
		defer teardownFunc()

		prefRepo := NewPreferenceRepository(client)

		_, err := prefRepo.Get(t.Context(), &malak.Workspace{
			ID: uuid.New(),
		})
		require.Error(t, err)
	})
}

func TestPreferences_Update(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	prefRepo := NewPreferenceRepository(client)

	repo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	planRepo := NewPlanRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	plan, err := planRepo.Get(t.Context(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtydaJZymT",
	})
	require.NoError(t, err)

	opts := &malak.CreateWorkspaceOptions{
		User:      user,
		Workspace: malak.NewWorkspace("oops", user, plan, malak.GenerateReference(malak.EntityTypeWorkspace)),
	}

	require.NoError(t, repo.Create(t.Context(), opts))

	pref, err := prefRepo.Get(t.Context(), opts.Workspace)
	require.NoError(t, err)

	require.True(t, pref.Communication.EnableMarketing)
	require.True(t, pref.Communication.EnableProductUpdates)

	pref.Communication.EnableMarketing = false
	require.NoError(t, prefRepo.Update(t.Context(), pref))

	newPref, err := prefRepo.Get(t.Context(), opts.Workspace)
	require.NoError(t, err)

	require.False(t, newPref.Communication.EnableMarketing)
	require.True(t, newPref.Communication.EnableProductUpdates)
}
