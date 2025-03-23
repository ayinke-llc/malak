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

	link := malak.NewReferenceGenerator().ShortLink()

	// Create a new deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	err = deck.Create(t.Context(), &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Public Deck Test",
		ShortLink:   link,
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
	deckFromDB, err := deck.PublicDetails(t.Context(), malak.Reference(link))
	require.NoError(t, err)
	require.NotNil(t, deckFromDB)
	require.Equal(t, ref.String(), deckFromDB.Reference.String())
	require.Equal(t, "Public Deck Test", deckFromDB.Title)
	require.NotNil(t, deckFromDB.DeckPreference)
	require.True(t, deckFromDB.DeckPreference.RequireEmail)
	require.True(t, deckFromDB.DeckPreference.EnableDownloading)
	require.True(t, deckFromDB.DeckPreference.Password.Enabled)
}

func TestDeck_CreateDeckSession(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test deck first
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	testDeck := &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Session Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), testDeck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Create a deck session
	session := &malak.DeckViewerSession{
		DeckID:    testDeck.ID,
		SessionID: malak.NewReferenceGenerator().Generate(malak.EntityTypeSession),
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckViewerSession),
	}

	err = deck.CreateDeckSession(t.Context(), session)
	require.NoError(t, err)
	require.NotEmpty(t, session.ID)
}

func TestDeck_UpdateDeckSession(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	testDeck := &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Session Update Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), testDeck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Create initial session
	session := &malak.DeckViewerSession{
		DeckID:     testDeck.ID,
		SessionID:  malak.NewReferenceGenerator().Generate(malak.EntityTypeSession),
		DeviceInfo: "Test Device",
		OS:         "macOS",
		Browser:    "Chrome",
		IPAddress:  "127.0.0.1",
		Country:    "NG",
		City:       "Lagos",
	}

	err = deck.CreateDeckSession(t.Context(), session)
	require.NoError(t, err)

	// Create a contact and update session
	contact := &malak.Contact{
		Email:       malak.Email("test@example.com"),
		FirstName:   "Test",
		LastName:    "User",
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
	}

	updateOpts := &malak.UpdateDeckSessionOptions{
		Session:       session,
		Contact:       contact,
		CreateContact: true,
	}

	err = deck.UpdateDeckSession(t.Context(), updateOpts)
	require.NoError(t, err)
	require.NotEmpty(t, session.ContactID)
}

func TestDeck_FindDeckSession(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	testDeck := &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Find Session Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), testDeck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Test finding non-existent session
	_, err = deck.FindDeckSession(t.Context(), "non-existent-session")
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrDeckNotFound)

	// Create a session
	sessionID := malak.NewReferenceGenerator().Generate(malak.EntityTypeSession)
	session := &malak.DeckViewerSession{
		DeckID:     testDeck.ID,
		Reference:  malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckViewerSession),
		SessionID:  sessionID,
		DeviceInfo: "Test Device",
		OS:         "macOS",
		Browser:    "Chrome",
		IPAddress:  "127.0.0.1",
		Country:    "NG",
		City:       "Lagos",
	}

	err = deck.CreateDeckSession(t.Context(), session)
	require.NoError(t, err)

	// Find the session
	foundSession, err := deck.FindDeckSession(t.Context(), sessionID.String())
	require.NoError(t, err)
	require.NotNil(t, foundSession)
	require.Equal(t, sessionID.String(), foundSession.SessionID.String())
	require.Equal(t, testDeck.ID, foundSession.DeckID)
}

func TestDeck_SessionAnalytics(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	testDeck := &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Analytics Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), testDeck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Create multiple sessions
	for i := 0; i < 5; i++ {
		session := &malak.DeckViewerSession{
			DeckID:     testDeck.ID,
			Reference:  malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckViewerSession),
			SessionID:  malak.NewReferenceGenerator().Generate(malak.EntityTypeSession),
			DeviceInfo: "Test Device",
			OS:         "macOS",
			Browser:    "Chrome",
			IPAddress:  "127.0.0.1",
			Country:    "NG",
			City:       "Lagos",
		}
		err = deck.CreateDeckSession(t.Context(), session)
		require.NoError(t, err)
	}

	// Test analytics retrieval
	opts := &malak.ListSessionAnalyticsOptions{
		DeckID: testDeck.ID,
		Days:   7,
		Paginator: malak.Paginator{
			Page:    1,
			PerPage: 10,
		},
	}

	sessions, total, err := deck.SessionAnalytics(t.Context(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, sessions)
	require.Equal(t, int64(5), total)
}

func TestDeck_DeckEngagements(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deck := NewDeckRepository(client)
	workspaceRepo := NewWorkspaceRepository(client)
	userRepo := NewUserRepository(client)

	// Get test user
	user, err := userRepo.Get(t.Context(), &malak.FindUserOptions{
		Email: "lanre@test.com",
	})
	require.NoError(t, err)

	// Get test workspace
	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test deck
	ref := malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck)
	testDeck := &malak.Deck{
		Reference:   ref,
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Engagements Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}

	err = deck.Create(t.Context(), testDeck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	// Create sessions from different countries
	countries := []string{"NG", "US", "GB", "CA", "FR"}
	for _, country := range countries {
		session := &malak.DeckViewerSession{
			DeckID:     testDeck.ID,
			Reference:  malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckViewerSession),
			SessionID:  malak.NewReferenceGenerator().Generate(malak.EntityTypeSession),
			DeviceInfo: "Test Device",
			OS:         "macOS",
			Browser:    "Chrome",
			IPAddress:  "127.0.0.1",
			Country:    country,
			City:       "Test City",
		}
		err = deck.CreateDeckSession(t.Context(), session)
		require.NoError(t, err)
	}

	// Manually run the aggregation queries that would normally be run by the cron job
	// Copied from cmd/cron_analytics
	engagementQuery := `
		WITH dates_to_process AS (
			SELECT DISTINCT
				dvs.deck_id,
				d.workspace_id,
				DATE(dvs.viewed_at) as engagement_date
			FROM deck_viewer_sessions dvs
			INNER JOIN decks d ON d.id = dvs.deck_id
			WHERE dvs.deleted_at IS NULL
			AND d.deleted_at IS NULL
			AND DATE(dvs.viewed_at) = CURRENT_DATE
			AND NOT EXISTS (
				SELECT 1 FROM deck_daily_engagements dde
				WHERE dde.deck_id = dvs.deck_id
				AND dde.workspace_id = d.workspace_id
				AND dde.engagement_date = CURRENT_DATE
			)
		),
		daily_stats AS (
			SELECT 
				dvs.deck_id,
				d.workspace_id,
				CURRENT_DATE as engagement_date,
				COUNT(DISTINCT dvs.id) as engagement_count
			FROM deck_viewer_sessions dvs
			INNER JOIN decks d ON d.id = dvs.deck_id
			WHERE dvs.deleted_at IS NULL
			AND d.deleted_at IS NULL
			AND DATE(dvs.viewed_at) = CURRENT_DATE
			GROUP BY dvs.deck_id, d.workspace_id
		)
		INSERT INTO deck_daily_engagements (
			reference,
			deck_id,
			workspace_id,
			engagement_count,
			engagement_date
		)
		SELECT 
			'deck_daily_engagement_' || LOWER(REPLACE(uuid_generate_v4()::text, '-', '')),
			deck_id,
			workspace_id,
			engagement_count,
			engagement_date
		FROM daily_stats
		ON CONFLICT (deck_id, workspace_id, engagement_date)
		DO UPDATE SET
			engagement_count = EXCLUDED.engagement_count,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = client.ExecContext(t.Context(), engagementQuery)
	require.NoError(t, err)

	geoQuery := `
		WITH geo_stats AS (
			SELECT 
				deck_id,
				COALESCE(NULLIF(TRIM(country), ''), 'Unknown') as country,
				COUNT(DISTINCT id) as view_count,
				CURRENT_DATE as stat_date
			FROM deck_viewer_sessions
			WHERE deleted_at IS NULL
			GROUP BY deck_id, COALESCE(NULLIF(TRIM(country), ''), 'Unknown')
		)
		INSERT INTO deck_geographic_stats (
			reference,
			deck_id,
			country,
			view_count,
			stat_date
		)
		SELECT 
			'deck_geographic_stat_' || LOWER(REPLACE(uuid_generate_v4()::text, '-', '')),
			deck_id,
			country,
			view_count,
			stat_date
		FROM geo_stats
		ON CONFLICT (deck_id, country, stat_date)
		DO UPDATE SET
			view_count = EXCLUDED.view_count,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = client.ExecContext(t.Context(), geoQuery)
	require.NoError(t, err)

	// Test engagements retrieval
	opts := &malak.ListDeckEngagementsOptions{
		DeckID: testDeck.ID,
	}

	engagements, err := deck.DeckEngagements(t.Context(), opts)
	require.NoError(t, err)
	require.NotNil(t, engagements)
	require.NotEmpty(t, engagements.GeographicStats)
	require.NotEmpty(t, engagements.DailyEngagements)
}

func TestDeck_Overview(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	deckRepo := NewDeckRepository(client)
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

	overview, err := deckRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Equal(t, int64(0), overview.TotalDecks)
	require.Equal(t, int64(0), overview.TotalViewerSessions)

	deck := &malak.Deck{
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDeck),
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Title:       "Test Deck",
		ShortLink:   malak.NewReferenceGenerator().ShortLink(),
		ObjectKey:   uuid.NewString(),
	}
	err = deckRepo.Create(t.Context(), deck, &malak.CreateDeckOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckPreference),
	})
	require.NoError(t, err)

	session := &malak.DeckViewerSession{
		DeckID:    deck.ID,
		SessionID: malak.NewReferenceGenerator().Generate(malak.EntityTypeSession),
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeDeckViewerSession),
	}
	err = deckRepo.CreateDeckSession(t.Context(), session)
	require.NoError(t, err)

	overview, err = deckRepo.Overview(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), overview.TotalDecks)
	require.Equal(t, int64(1), overview.TotalViewerSessions)
}
