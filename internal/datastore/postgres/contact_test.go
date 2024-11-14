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

	_, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)

	_, err = contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   "contact_kCo",
		WorkspaceID: workspaceID,
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
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
	}

	err := contactRepo.Create(context.Background(), contact)
	require.NoError(t, err)

	// create again
	err = contactRepo.Create(context.Background(), contact)
	require.Error(t, err)
	require.Equal(t, err, malak.ErrContactExists)
}
