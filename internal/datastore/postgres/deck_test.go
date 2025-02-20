package postgres

import (
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
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	decks := &malak.Deck{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
		Title:       "fojgfnolkgj",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		CreatedBy:   user.ID,
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), decks, opts)
	require.NoError(t, err)
}

func TestDeck_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	_ = user

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	decks, err := deck.List(t.Context(), workspace)

	require.NoError(t, err)
	require.Len(t, decks, 0)

	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "oops",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})

	require.NoError(t, err)

	decks, err = deck.List(t.Context(), workspace)

	require.NoError(t, err)
	require.Len(t, decks, 1)
}

func TestDeck_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	_ = user

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	_, err = deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   "oops",
		WorkspaceID: workspace.ID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrDeckNotFound)

	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)

	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "oops",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})

	require.NoError(t, err)

	_, err = deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)
}

func TestDeck_Delete(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	_ = user

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)

	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "oops",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})

	require.NoError(t, err)

	deckFromDB, err := deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)

	require.NoError(t, deck.Delete(t.Context(), deckFromDB))

	_, err = deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrDeckNotFound)
}

func TestDeck_UpdatePreferences(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	_ = user

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)

	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "oops",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
		RequireEmail:      false,
		EnableDownloading: false,
	})

	require.NoError(t, err)

	deckFromDB, err := deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)

	deckFromDB.DeckPreference.RequireEmail = true
	deckFromDB.DeckPreference.EnableDownloading = true

	require.NoError(t, deck.UpdatePreferences(t.Context(), deckFromDB))

	deckFromDatabase, err := deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)

	require.True(t, deckFromDatabase.DeckPreference.RequireEmail)
	require.True(t, deckFromDatabase.DeckPreference.EnableDownloading)
}

func TestDeck_ToggleArchive(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)

	// Create a new deck
	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Test Archive Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Get the deck
	deckFromDB, err := deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)
	require.False(t, deckFromDB.IsArchived)

	// Toggle archive (true)
	err = deck.ToggleArchive(t.Context(), deckFromDB)
	require.NoError(t, err)

	// Verify it's archived
	deckFromDB, err = deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)
	require.True(t, deckFromDB.IsArchived)

	// Toggle archive again (false)
	err = deck.ToggleArchive(t.Context(), deckFromDB)
	require.NoError(t, err)

	// Verify it's unarchived
	deckFromDB, err = deck.Get(t.Context(), malak.FetchDeckOptions{
		Reference:   ref.String(),
		WorkspaceID: workspace.ID,
	})
	require.NoError(t, err)
	require.False(t, deckFromDB.IsArchived)
}

func TestDecks_TogglePinned(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	userRepo := NewUserRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)

	deck := NewDeckRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// add 3 pinnned decks
	for i := 0; i <= 3; i++ {
		opts := &malak.CreateDeckOptions{
			Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
		}

		decks := &malak.Deck{
			WorkspaceID: workspace.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
			Title:       "fojgfnolkgj",
			ShortLink:   malak.NewReferenceGenerator().ShortLink(),
			CreatedBy:   user.ID,
			IsPinned:    true,
			ObjectKey:   uuid.NewString(),
		}

		err = deck.Create(t.Context(), decks, opts)
		require.NoError(t, err)
	}

	opts := &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	}

	decks := &malak.Deck{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
		Title:       "fojgfnolkgj",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		CreatedBy:   user.ID,
		ObjectKey:   uuid.NewString(),
	}

	// create another without pinning
	err = deck.Create(t.Context(), decks, opts)
	require.NoError(t, err)

	// cannot add a 5th pinned item
	err = deck.TogglePinned(t.Context(), decks)
	require.Error(t, err)
	require.Equal(t, malak.ErrPinnedDeckCapacityExceeded, err)
}

func TestDeck_PublicDetails(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// user from the fixtures
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// from workspaces.yml migration
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Test non-existent deck
	nonExistentRef := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	_, err = deck.PublicDetails(t.Context(), nonExistentRef)
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrDeckNotFound)

	// Create a new deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Public Deck Test",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}, &malak.CreateDeckOptions{
		Reference:         malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
		RequireEmail:      true,
		EnableDownloading: true,
		Password: struct {
			Enabled  bool           `json:"enabled,omitempty" validate:"required"`
			Password malak.Password `json:"password,omitempty" validate:"required"`
		}{
			Enabled:  true,
			Password: malak.Password("test-password"),
		},
	})
	require.NoError(t, err)

	// Test fetching the deck's public details
	deckFromDB, err := deck.PublicDetails(t.Context(), ref)
	require.NoError(t, err)
	require.NotNil(t, deckFromDB)
	require.Equal(t, ref.String(), deckFromDB.Reference.String())
	require.Equal(t, "Public Deck Test", deckFromDB.Title)
	require.NotNil(t, deckFromDB.DeckPreference)
	require.True(t, deckFromDB.DeckPreference.RequireEmail)
	require.True(t, deckFromDB.DeckPreference.EnableDownloading)
	require.True(t, deckFromDB.DeckPreference.Password.Enabled)
}
