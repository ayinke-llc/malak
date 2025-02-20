package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContactList_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	contactRepo := NewContactListRepository(client)

	list := &malak.ContactList{
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		Title:       "My contact list",
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeList),
		CreatedBy:   user.ID,
	}

	err = contactRepo.Create(t.Context(), list)
	require.NoError(t, err)

	newList, err := contactRepo.Get(t.Context(), malak.FetchContactListOptions{
		Reference:   list.Reference,
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)
	require.Equal(t, list.Title, newList.Title)

	newList.Title = "Series A"

	err = contactRepo.Update(t.Context(), newList)
	require.NoError(t, err)

	err = contactRepo.Delete(t.Context(), newList)
	require.NoError(t, err)

	// fetch again
	//

	_, err = contactRepo.Get(t.Context(), malak.FetchContactListOptions{
		Reference:   list.Reference,
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrContactListNotFound)
}

func TestContactList(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	contactRepo := NewContactListRepository(client)

	for range []string{"", ""} {

		list := &malak.ContactList{
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
			Title:       "My contact list",
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeList),
			CreatedBy:   user.ID,
		}

		err = contactRepo.Create(t.Context(), list)
		require.NoError(t, err)

	}

	lists, _, err := contactRepo.List(t.Context(),
		&malak.ContactListOptions{
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		})
	require.NoError(t, err)

	require.Len(t, lists, 2)
}

func TestContactList_Add(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	contactListRepo := NewContactListRepository(client)

	contactRepo := NewContactRepository(client)

	contact := &malak.Contact{
		Email:       malak.Email("oops@oops.com"),
		WorkspaceID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
		CreatedBy:   uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		OwnerID:     uuid.MustParse("1aa6b38e-33d3-499f-bc9d-3090738f29e6"),
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
	}

	err = contactRepo.Create(t.Context(), contact)
	require.NoError(t, err)

	list := &malak.ContactList{
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		Title:       "My contact list",
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeList),
		CreatedBy:   user.ID,
	}

	err = contactListRepo.Create(t.Context(), list)
	require.NoError(t, err)

	newList, err := contactListRepo.Get(t.Context(), malak.FetchContactListOptions{
		Reference:   list.Reference,
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)
	require.Equal(t, list.Title, newList.Title)

	_, mappings, err := contactListRepo.List(t.Context(),
		&malak.ContactListOptions{
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		})
	require.NoError(t, err)

	require.Len(t, mappings, 0)

	err = contactListRepo.Add(t.Context(), &malak.ContactListMapping{
		ListID:    newList.ID,
		ContactID: contact.ID,
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeListEmail),
		CreatedBy: user.ID,
	})
	require.NoError(t, err)

	_, mappings, err = contactListRepo.List(t.Context(),
		&malak.ContactListOptions{
			WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
		})
	require.NoError(t, err)

	require.Len(t, mappings, 1)
}
