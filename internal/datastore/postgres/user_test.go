package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUser_Update(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.Equal(t, "Lanre Adelowo", user.FullName)

	newName := "Lebron James"

	user.FullName = newName
	require.NoError(t, userRepo.Update(context.TODO(), user))

	// fetch the user again and check the name
	fetchUser, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)
	require.Equal(t, newName, fetchUser.FullName)
}

func TestUser_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	tt := []struct {
		email    string
		name     string
		hasError bool
	}{
		{
			email:    "lanre@test.com",
			name:     "user from fixtures trying to be recreated",
			hasError: true,
		},
		{
			email: "newuser@google.com",
			name:  "new user added",
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			err := userRepo.Create(context.Background(), &malak.User{
				Email:    malak.Email(v.email),
				FullName: "Lanre",
			})
			if v.hasError {
				require.Error(t, err)
				require.Equal(t, malak.ErrUserExists, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUser_GetUserID(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	tt := []struct {
		id       string
		name     string
		hasError bool
	}{
		{
			id:   "1aa6b38e-33d3-499f-bc9d-3090738f29e6",
			name: "User 1 from fixtures",
		},
		{
			id:       "fe76e7a4-9e9b-4cb6-934f-79e528b7c016",
			name:     "user does not exists",
			hasError: true,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
				ID: uuid.MustParse(v.id),
			})
			if v.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, user.FullName)
			require.Equal(t, v.id, user.ID.String())
		})
	}
}

func TestUser_Get(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)

	tt := []struct {
		email    string
		name     string
		hasError bool
	}{
		{
			email: "lanre@test.com",
			name:  "User 1 from fixtures",
		},
		{
			email:    "unknown@google.com",
			name:     "user does not exists",
			hasError: true,
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
				Email: malak.Email(v.email),
			})
			if v.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, user.FullName)
			require.Equal(t, v.email, user.Email.String())
		})
	}
}
