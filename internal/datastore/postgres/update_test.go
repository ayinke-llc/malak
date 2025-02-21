package postgres

import (
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
	updateByID, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		ID:          id,
		Reference:   "update_NCox_gRNg",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	require.NoError(t, updatesRepo.Delete(t.Context(), updateByID))

	_, err = updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		ID:          id,
		Reference:   "update_NCox_gRNg",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrUpdateNotFound, err)
}

func TestUpdates_Update(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	update, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		Reference:   "update_O-54dq6IR",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
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

	require.NoError(t, updatesRepo.Update(t.Context(), update))

	updatedItem, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		Reference:   "update_O-54dq6IR",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	require.Equal(t, updatedItem.Content, updatedContent)
}

func TestUpdates_GetByID(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	_, err := updatesRepo.GetByID(t.Context(),
		uuid.MustParse("07b0c648-12fd-44fc-a280-946de2700e65"))
	require.NoError(t, err)

	_, err = updatesRepo.GetByID(t.Context(), uuid.New())
	require.Error(t, err)
	require.Equal(t, err, malak.ErrUpdateNotFound)
}

func TestUpdates_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	_, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		Reference:   "update_O-54dq6IR",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	updateByID, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		ID:          uuid.MustParse("07b0c648-12fd-44fc-a280-946de2700e65"),
		Reference:   "update_O-54dq6IR",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	update, err := updatesRepo.Get(t.Context(), malak.FetchUpdateOptions{
		Status:      malak.UpdateStatusDraft,
		ID:          uuid.MustParse("07b0c648-12fd-44fc-a280-946de2700e65"),
		Reference:   "update_O-54dq6IR",
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	require.Equal(t, update, updateByID)
}

func TestUpdates_StatUpdate(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     make([]malak.Block, 0),
		Reference:   "update_ifjfkjfo",
	}

	err = updatesRepo.Create(t.Context(), update)
	require.NoError(t, err)

	stat, err := updatesRepo.Stat(t.Context(), update)
	require.NoError(t, err)

	require.Equal(t, int64(0), stat.TotalOpens)

	stat.TotalOpens = 10
	require.NoError(t, updatesRepo.UpdateStat(t.Context(), stat, nil))

	newStat, err := updatesRepo.Stat(t.Context(), update)
	require.NoError(t, err)

	require.Equal(t, int64(10), newStat.TotalOpens)
}

func TestUpdates_Stat(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     make([]malak.Block, 0),
		Reference:   "update_ifjfkjfo",
	}

	err = updatesRepo.Create(t.Context(), update)
	require.NoError(t, err)

	_, err = updatesRepo.Stat(t.Context(), update)
	require.NoError(t, err)
}

func TestUpdates_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	err = updatesRepo.Create(t.Context(), &malak.Update{
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
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	updates, total, err := updatesRepo.List(t.Context(), malak.ListUpdateOptions{
		WorkspaceID: workspace.ID,
		Status:      malak.ListUpdateFilterStatusAll,
	})
	require.NoError(t, err)

	require.Len(t, updates, 0)
	require.Equal(t, int64(0), total)

	err = updatesRepo.Create(t.Context(), &malak.Update{
		WorkspaceID: workspace.ID,
		Status:      malak.UpdateStatusDraft,
		CreatedBy:   user.ID,
		Content:     []malak.Block{},
		Reference:   "update_ifjfkjfo",
	})
	require.NoError(t, err)

	updates, total, err = updatesRepo.List(t.Context(), malak.ListUpdateOptions{
		WorkspaceID: workspace.ID,
		Status:      malak.ListUpdateFilterStatusAll,
	})
	require.NoError(t, err)

	require.Len(t, updates, 1)
	require.Equal(t, int64(1), total)
}

func TestUpdates_TogglePinned(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)
	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	refGenerator := malak.NewReferenceGenerator()

	// add 3 pinnned updates
	for i := 0; i <= 3; i++ {
		err = updatesRepo.Create(t.Context(), &malak.Update{
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

	err = updatesRepo.Create(t.Context(), update)
	require.NoError(t, err)

	// cannot add a 4th pinned item
	err = updatesRepo.TogglePinned(t.Context(), update)
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
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
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

	err = updatesRepo.Create(t.Context(), update)
	require.NoError(t, err)

	err = updatesRepo.SendUpdate(t.Context(), &malak.CreateUpdateOptions{
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
		Plan:        workspace.Plan,
	})
	require.NoError(t, err)
}

func TestUpdates_ListPinned(t *testing.T) {

	t.Run("no pinned items", func(t *testing.T) {

		client, teardownFunc := setupDatabase(t)
		defer teardownFunc()

		updatesRepo := NewUpdatesRepository(client)
		userRepo := NewUserRepository(client)
		workspaceRepo := NewWorkspaceRepository(client)

		// user from the fixtures
		user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
			Email: "lanre@test.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)

		// from workspaces.yml migration
		workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
			ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		})
		require.NoError(t, err)

		for range []int{0, 1, 2, 3, 4, 5} {
			err = updatesRepo.Create(t.Context(), &malak.Update{
				WorkspaceID: workspace.ID,
				Status:      malak.UpdateStatusDraft,
				CreatedBy:   user.ID,
				Content:     []malak.Block{},
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeUpdate),
			})
			require.NoError(t, err)
		}

		updates, err := updatesRepo.ListPinned(t.Context(), workspace.ID)
		require.NoError(t, err)
		require.Len(t, updates, 0)
	})

	t.Run("pinned items count", func(t *testing.T) {

		client, teardownFunc := setupDatabase(t)
		defer teardownFunc()

		updatesRepo := NewUpdatesRepository(client)
		userRepo := NewUserRepository(client)
		workspaceRepo := NewWorkspaceRepository(client)

		// user from the fixtures
		user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
			Email: "lanre@test.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)

		// from workspaces.yml migration
		workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
			ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		})
		require.NoError(t, err)

		for range []int{0, 1, 2, 3} {
			update := &malak.Update{
				WorkspaceID: workspace.ID,
				Status:      malak.UpdateStatusDraft,
				CreatedBy:   user.ID,
				Content:     []malak.Block{},
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeUpdate),
			}

			err = updatesRepo.Create(t.Context(), update)
			require.NoError(t, err)

			require.NoError(t, updatesRepo.TogglePinned(t.Context(), update))
		}

		update := &malak.Update{
			WorkspaceID: workspace.ID,
			Status:      malak.UpdateStatusDraft,
			CreatedBy:   user.ID,
			Content:     []malak.Block{},
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeUpdate),
		}

		err = updatesRepo.Create(t.Context(), update)
		require.NoError(t, err)

		// max 5 have been added as pinned items
		require.Error(t, updatesRepo.TogglePinned(t.Context(), update))

		updates, err := updatesRepo.ListPinned(t.Context(), workspace.ID)
		require.NoError(t, err)
		require.Len(t, updates, malak.MaximumNumberOfPinnedUpdates)
	})
}

func TestUpdates_GetStatByEmailID(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	updatesRepo := NewUpdatesRepository(client)

	_, _, err := updatesRepo.GetStatByEmailID(t.Context(),
		"random", malak.UpdateRecipientLogProviderResend)
	require.Error(t, err)
}
