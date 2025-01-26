package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/stretchr/testify/require"
)

func TestShare_All(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	shareRepo := NewShareRepository(client)
	contactRepo := NewContactRepository(client)

	contact, err := contactRepo.Get(context.Background(), malak.FetchContactOptions{
		Reference:   "contact_kCoC286IR", // contacts.yml
		WorkspaceID: workspaceID,
	})
	require.NoError(t, err)
	require.NotNil(t, contact)

	sharedContacts, err := shareRepo.All(context.Background(), contact)
	require.NoError(t, err)
	require.Len(t, sharedContacts, 0)
}
