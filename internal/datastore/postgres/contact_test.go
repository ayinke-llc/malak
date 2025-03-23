package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	workspaceID = uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526")
)

func TestContact_Delete(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contact := &malak.Contact{
		Email:       malak.Email("oops@oops.com"),
		WorkspaceID: workspaceID,
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
	}

	// Test successful creation
	err := contactRepo.Create(t.Context(), contact)
	require.NoError(t, err)

	// Verify contact was created
	savedContact, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   contact.Reference,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.Email, savedContact.Email)
	require.Equal(t, contact.WorkspaceID, savedContact.WorkspaceID)
	require.Equal(t, contact.Reference, savedContact.Reference)

	// Delete
	err = contactRepo.Delete(t.Context(), contact)
	require.NoError(t, err)

	//  contact was deleted, it should not be found
	_, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   contact.Reference,
		WorkspaceID: workspaceID,
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrContactNotFound, err)

	// get contact from db
	contact, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "contact_kCoC286IR", contact.Reference.String())

	// Delete again
	err = contactRepo.Delete(t.Context(), contact)
	require.NoError(t, err)
}

func TestContact_Get(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contact, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "contact_kCoC286IR", contact.Reference.String())

	contactByID, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		ID:          contact.ID,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.ID, contactByID.ID)

	contactByEmail, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Email:       contact.Email,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.Email, contactByEmail.Email)

	_, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   "contact_kCo",
		WorkspaceID: workspaceID,
	})
	require.Error(t, err)
	require.Equal(t, err, malak.ErrContactNotFound)

	_, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   malak.Reference(contact.Reference.String()),
		WorkspaceID: uuid.New(),
	})
	require.Error(t, err)
	require.Equal(t, err, malak.ErrContactNotFound)
}

func TestContact_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contact := &malak.Contact{
		Email:       malak.Email("oops@oops.com"),
		WorkspaceID: workspaceID,
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
	}

	// Test successful creation
	err := contactRepo.Create(t.Context(), contact)
	require.NoError(t, err)

	// Verify contact was created
	savedContact, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   contact.Reference,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.Email, savedContact.Email)
	require.Equal(t, contact.WorkspaceID, savedContact.WorkspaceID)
	require.Equal(t, contact.Reference, savedContact.Reference)

	// Test duplicate creation
	err = contactRepo.Create(t.Context(), contact)
	require.Error(t, err)
	require.Equal(t, err, malak.ErrContactExists)
}

func TestContact_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contacts := []*malak.Contact{
		{
			Email:       malak.Email("test1@example.com"),
			WorkspaceID: workspaceID,
			CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		},
		{
			Email:       malak.Email("test2@example.com"),
			WorkspaceID: workspaceID,
			CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
		},
	}

	for _, c := range contacts {
		err := contactRepo.Create(t.Context(), c)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		opts          malak.ListContactOptions
		expectedCount int
	}{
		{
			name: "First page",
			opts: malak.ListContactOptions{
				WorkspaceID: workspaceID,
				Paginator: malak.Paginator{
					Page:    1,
					PerPage: 2,
				},
			},
			expectedCount: 2,
		},
		{
			name: "Second page",
			opts: malak.ListContactOptions{
				WorkspaceID: workspaceID,
				Paginator: malak.Paginator{
					Page:    2,
					PerPage: 1,
				},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, total, err := contactRepo.List(t.Context(), tt.opts)
			require.NoError(t, err)
			require.Greater(t, total, int64(0))
			require.Len(t, result, tt.expectedCount)
		})
	}

	result, total, err := contactRepo.List(t.Context(), malak.ListContactOptions{
		WorkspaceID: uuid.New(),
		Paginator: malak.Paginator{
			Page:    1,
			PerPage: 10,
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
	require.Len(t, result, 0)
}

func TestContact_Update(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contact, err := contactRepo.Get(t.Context(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "contact_kCoC286IR", contact.Reference.String())

	newEmail := faker.Email()

	_, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		ID:          contact.ID,
		WorkspaceID: workspaceID,
		Email:       malak.Email(newEmail),
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrContactNotFound, err)

	contact.Email = malak.Email(newEmail)
	require.NoError(t, contactRepo.Update(t.Context(), contact))

	_, err = contactRepo.Get(t.Context(), malak.FetchContactOptions{
		ID:          contact.ID,
		WorkspaceID: workspaceID,
		Email:       malak.Email(newEmail),
	})
	require.NoError(t, err)
}

func TestContact_Overview(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user from fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace from fixtures
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	overview, err := contactRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Equal(t, int64(0), overview.TotalContacts)

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

	overview, err = contactRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), overview.TotalContacts)
}
