package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdates_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	err = updatesRepo.Create(context.Background(), &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     "",
		Reference:   "update_ifjfkjfo",
	})
	require.NoError(t, err)
}

func TestUpdates_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	updates, err := updatesRepo.List(context.Background(), malak.ListUpdateOptions{
		WorkspaceID: workspace.ID,
		Status:      malak.ListUpdateFilterStatusAll,
	})
	require.NoError(t, err)

	require.Len(t, updates, 0)

	err = updatesRepo.Create(context.Background(), &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     "",
		Reference:   "update_ifjfkjfo",
	})
	require.NoError(t, err)

	updates, err = updatesRepo.List(context.Background(), malak.ListUpdateOptions{
		WorkspaceID: workspace.ID,
		Status:      malak.ListUpdateFilterStatusAll,
	})
	require.NoError(t, err)

	require.Len(t, updates, 1)
}