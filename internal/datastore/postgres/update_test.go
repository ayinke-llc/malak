package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdates_Delete(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	id := uuid.MustParse("0902ef67-903e-47b8-8f9d-111a9e0ca0c7")

	// workspaces.yml testdata
	updateByID, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		ID:        id,
		Reference: "update_NCox_gRNg",
	})
	require.NoError(t, err)

	require.NoError(t, updatesRepo.Delete(context.Background(), updateByID))

	_, err = updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		ID:        id,
		Reference: "update_NCox_gRNg",
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrUpdateNotFound, err)
}

func TestUpdates_Update(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	update, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		Reference: "update_O-54dq6IR",
	})
	require.NoError(t, err)

	updatedContent := malak.BlockContents{
		{
			ID:   "oops",
			Type: "heading",
		},
	}

	require.NotEqual(t, update.Content, updatedContent)

	update.Content = updatedContent

	require.NoError(t, updatesRepo.Update(context.Background(), update))

	updatedItem, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		Reference: "update_O-54dq6IR",
	})
	require.NoError(t, err)

	require.Equal(t, updatedItem.Content, updatedContent)
}

func TestUpdates_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	_, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		Reference: "update_O-54dq6IR",
	})
	require.NoError(t, err)

	updateByID, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		ID:        uuid.MustParse("07b0c648-12fd-44fc-a280-946de2700e65"),
		Reference: "update_O-54dq6IR",
	})
	require.NoError(t, err)

	update, err := updatesRepo.Get(context.Background(), malak.FetchUpdateOptions{
		Status:    malak.UpdateStatusDraft,
		ID:        uuid.MustParse("07b0c648-12fd-44fc-a280-946de2700e65"),
		Reference: "update_O-54dq6IR",
	})
	require.NoError(t, err)

	require.Equal(t, update, updateByID)
}

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
		Content:     make([]malak.Block, 0),
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
		Content:     []malak.Block{},
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

func TestUpdates_TogglePinned(t *testing.T) {

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

	refGenerator := malak.NewReferenceGenerator()

	// add 3 pinnned updates
	for i := 0; i < 3; i++ {
		err = updatesRepo.Create(context.Background(), &malak.Update{
			WorkspaceID: workspace.ID,
			Status:      malak.UpdateStatusDraft,
			CreatedBy:   user.ID,
			Content:     make([]malak.Block, 0),
			Reference:   refGenerator.Generate(malak.EntityTypeUpdate),
			IsPinned:    true,
		})
		require.NoError(t, err)
	}

	ref := refGenerator.Generate(malak.EntityTypeUpdate)

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     []malak.Block{},
		Reference:   ref,
	}

	err = updatesRepo.Create(context.Background(), update)
	require.NoError(t, err)

	// cannot add a 4th pinned item
	err = updatesRepo.TogglePinned(context.Background(), update)
	require.Error(t, err)
	require.Equal(t, malak.ErrPinnedUpdateCapacityExceeded, err)
}

func TestUpdates_SendUpdate(t *testing.T) {

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

	refGenerator := malak.NewReferenceGenerator()

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     make([]malak.Block, 0),
		Reference:   refGenerator.Generate(malak.EntityTypeUpdate),
		IsPinned:    true,
	}

	err = updatesRepo.Create(context.Background(), update)
	require.NoError(t, err)

	err = updatesRepo.SendUpdate(context.Background(), &malak.CreateUpdateOptions{
		Reference: func(et malak.EntityType) string {
			return string(refGenerator.Generate(et))
		},
		Generator: refGenerator,
		Emails:    []malak.Email{malak.Email("oops@oops.com"), malak.Email("oops@gmail.comf")}, // add an existing email
		UserID:    user.ID,
		Schedule: &malak.UpdateSchedule{
			Reference:   refGenerator.Generate(malak.EntityTypeSchedule),
			SendAt:      time.Now(),
			UpdateType:  malak.UpdateTypeLive,
			ScheduledBy: user.ID,
			Status:      malak.UpdateSendScheduleScheduled,
			UpdateID:    update.ID,
		},
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)
}
