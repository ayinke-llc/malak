package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ayinke-llc/malak"
)

func TestEmailVerification(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewEmailVerificationRepository(client)

	userRepo := NewUserRepository(client)
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	ev1, err := malak.NewEmailVerification(user)
	require.NoError(t, err)

	err = repo.Create(t.Context(), ev1)
	require.NoError(t, err)
	require.NotEmpty(t, ev1.ID)
	require.NotEmpty(t, ev1.CreatedAt)

	ev2, err := malak.NewEmailVerification(user)
	require.NoError(t, err)

	err = repo.Create(t.Context(), ev2)
	require.NoError(t, err)
	require.NotEmpty(t, ev2.ID)
	require.NotEmpty(t, ev2.CreatedAt)

	var verifications []malak.EmailVerification
	count, err := client.NewSelect().
		Table("email_verifications").
		Where("user_id = ?", user.ID).
		ScanAndCount(t.Context(), &verifications)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	require.Len(t, verifications, 1)

	var existing malak.EmailVerification
	err = client.NewSelect().
		Table("email_verifications").
		Where("user_id = ?", user.ID).
		Scan(t.Context(), &existing)
	require.NoError(t, err)
	require.Equal(t, ev2.Token, existing.Token)
	require.Equal(t, ev2.ID, existing.ID)
}
