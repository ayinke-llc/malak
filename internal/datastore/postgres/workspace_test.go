package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWorkspace_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	planRepo := NewPlanRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	plan, err := planRepo.Get(context.Background(), &malak.FetchPlanOptions{
		Reference: "prod_QmtErtydaJZymT",
	})
	require.NoError(t, err)

	require.NoError(t, repo.Create(context.Background(), &malak.CreateWorkspaceOptions{
		User:      user,
		Workspace: malak.NewWorkspace("oops", user, plan, malak.GenerateReference(malak.EntityTypeWorkspace)),
	}))
}

func TestWorkspace_Update(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewWorkspaceRepository(client)

	// from workspaces.yml migration
	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	newReference := "new-reference"

	workspace.Reference = newReference

	require.NoError(t, repo.Update(context.TODO(), workspace))

	newWorkspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	require.Equal(t, newWorkspace.Reference, newReference)
}

func TestWorkspace_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewWorkspaceRepository(client)

	// from workspaces.yml migration
	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	require.Equal(t, workspace.WorkspaceName, "First workspace")

	_, err = repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("cb5955cc-be42-4fe9-9155-250f4cc0ecc8"),
	})
	require.Error(t, err)
	require.Equal(t, err, malak.ErrWorkspaceNotFound)
}

func TestWorkspace_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewWorkspaceRepository(client)

	workspaces, err := repo.List(context.Background(), &malak.User{
		ID: uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
	})
	require.NoError(t, err)
	require.Len(t, workspaces, 2)
}
