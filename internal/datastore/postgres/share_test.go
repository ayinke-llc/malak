package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestShare_All(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	shareRepo := NewShareRepository(client)
	contactRepo := NewContactRepository(client)

	contact, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)

	sharedContacts, err := shareRepo.All(t.Context(), contact)
	require.NoError(t, err)
	require.Len(t, sharedContacts, 0)
}

func TestShare_Overview(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	shareRepo := NewShareRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)
	updateRepo := NewUpdatesRepository(client)
	contactRepo := NewContactRepository(client)

	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// initially there should be no shares
	overview, err := shareRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Empty(t, overview.RecentShares)

	contact := &malak.Contact{
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		WorkspaceID: workspace.ID,
		Email:       "test@example.com",
		FirstName:   "Test",
		LastName:    "User",
		CreatedBy:   user.ID,
		OwnerID:     user.ID,
		Metadata:    make(malak.CustomContactMetadata),
	}
	err = contactRepo.Create(t.Context(), contact)
	require.NoError(t, err)

	// create an update to share
	update := &malak.Update{
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeUpdate),
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Test Update",
		Status:      malak.UpdateStatusSent,
	}
	err = updateRepo.Create(t.Context(), update, &malak.TemplateCreateUpdateOptions{})
	require.NoError(t, err)

	share := &malak.ContactShare{
		Reference:     malak.NewReferenceGenerator().Generate(malak.EntityTypeContactShare),
		SharedBy:      user.ID,
		ContactID:     contact.ID,
		ItemType:      malak.ContactShareItemTypeUpdate,
		ItemID:        update.ID,
		ItemReference: update.Reference,
	}
	_, err = client.NewInsert().Model(share).Exec(t.Context())
	require.NoError(t, err)

	// check overview after creating share
	overview, err = shareRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, overview.RecentShares, 1)
	require.Equal(t, "Test Update", overview.RecentShares[0].Title)
	require.Equal(t, "test@example.com", overview.RecentShares[0].Email)
	require.Equal(t, "Test", overview.RecentShares[0].FirstName)
	require.Equal(t, "User", overview.RecentShares[0].LastName)
}
