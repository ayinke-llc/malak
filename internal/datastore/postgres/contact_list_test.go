package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContactList_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
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

	err = contactRepo.Create(context.Background(), list)
	require.NoError(t, err)

	newList, err := contactRepo.Get(context.Background(), malak.FetchContactListOptions{
		Reference:   list.Reference,
		WorkspaceID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)
	require.Equal(t, list.Title, newList.Title)

	newList.Title = "Series A"

	err = contactRepo.Update(context.Background(), newList)
	require.NoError(t, err)

	err = contactRepo.Delete(context.Background(), newList)
	require.NoError(t, err)

	// fetch again
	//

	_, err = contactRepo.Get(context.Background(), malak.FetchContactListOptions{
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

	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
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

		err = contactRepo.Create(context.Background(), list)
		require.NoError(t, err)

	}

	lists, err := contactRepo.List(context.Background(), uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"))
	require.NoError(t, err)

	require.Len(t, lists, 2)
}
