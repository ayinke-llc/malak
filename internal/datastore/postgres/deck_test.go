package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeck_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	opts := &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	}

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(context.Background(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	decks := &malak.Deck{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
		Title:       "fojgfnolkgj",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		CreatedBy:   user.ID,
	}

	err = deck.Create(context.Background(), decks, opts)
	require.NoError(t, err)
}
