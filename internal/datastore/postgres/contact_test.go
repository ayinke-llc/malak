package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	workspaceID = uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526")
)

func TestContact_Get(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	contactRepo := NewContactRepository(client)

	contact, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)
	require.Equal(t, "contact_kCoC286IR", contact.Reference.String())

	contactByID, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		ID:          contact.ID,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.ID, contactByID.ID)

	contactByEmail, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Email:       contact.Email,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.Email, contactByEmail.Email)

	_, err = contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   "contact_kCo",
		WorkspaceID: workspaceID,
	})
	require.Error(t, err)
	require.Equal(t, err, malak.ErrContactNotFound)

	_, err = contactRepo.Get(context.Background(), malak.FetchContactOptions{
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
	err := contactRepo.Create(context.Background(), contact)
	require.NoError(t, err)

	// Verify contact was created
	savedContact, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   contact.Reference,
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.Equal(t, contact.Email, savedContact.Email)
	require.Equal(t, contact.WorkspaceID, savedContact.WorkspaceID)
	require.Equal(t, contact.Reference, savedContact.Reference)

	// Test duplicate creation
	err = contactRepo.Create(context.Background(), contact)
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
		err := contactRepo.Create(context.Background(), c)
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
			result, total, err := contactRepo.List(context.Background(), tt.opts)
			require.NoError(t, err)
			require.Greater(t, total, int64(0))
			require.Len(t, result, tt.expectedCount)
		})
	}

	result, total, err := contactRepo.List(context.Background(), malak.ListContactOptions{
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
