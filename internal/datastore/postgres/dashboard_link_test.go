package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardLinkRepo_Create(t *testing.T) {
	db, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewDashboardLinkRepo(db)
	workspaceRepo := NewWorkspaceRepository(db)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}
	_, err = db.NewInsert().Model(dashboard).Exec(context.Background())
	require.NoError(t, err)

	tests := []struct {
		name    string
		opts    *malak.CreateDashboardLinkOptions
		wantErr bool
	}{
		{
			name: "create default link without contact",
			opts: &malak.CreateDashboardLinkOptions{
				WorkspaceID: workspace.ID,
				Link: &malak.DashboardLink{
					DashboardID: dashboard.ID,
					Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
					LinkType:    malak.DashboardLinkType("default"),
					ExpiresAt:   &time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "create link with new contact",
			opts: &malak.CreateDashboardLinkOptions{
				WorkspaceID: workspace.ID,
				Email:       malak.Email("test@example.com"),
				Link: &malak.DashboardLink{
					DashboardID: dashboard.ID,
					Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
					LinkType:    malak.DashboardLinkType("contact"),
					ExpiresAt:   &time.Time{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify link was created
			var link malak.DashboardLink
			err = db.NewSelect().
				Model(&link).
				Where("token = ?", tt.opts.Link.Token).
				Scan(context.Background())
			require.NoError(t, err)

			assert.Equal(t, tt.opts.Link.DashboardID, link.DashboardID)
			assert.Equal(t, tt.opts.Link.LinkType, link.LinkType)

			if tt.opts.Email != "" {
				// Verify contact was created and linked
				var contact malak.Contact
				err = db.NewSelect().
					Model(&contact).
					Where("email = ?", tt.opts.Email).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, contact.ID, link.ContactID)
			}
		})
	}
}

func TestDashboardLinkRepo_DefaultLink(t *testing.T) {
	db, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewDashboardLinkRepo(db)
	workspaceRepo := NewWorkspaceRepository(db)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}
	_, err = db.NewInsert().Model(dashboard).Exec(context.Background())
	require.NoError(t, err)

	link := &malak.DashboardLink{
		DashboardID: dashboard.ID,
		Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
		LinkType:    malak.DashboardLinkType("default"),
		ExpiresAt:   &time.Time{},
	}
	_, err = db.NewInsert().Model(link).Exec(context.Background())
	require.NoError(t, err)

	t.Run("get existing default link", func(t *testing.T) {
		got, err := repo.DefaultLink(context.Background(), dashboard)
		require.NoError(t, err)
		assert.Equal(t, link.Token, got.Token)
		assert.Equal(t, link.DashboardID, got.DashboardID)
	})

	t.Run("get non-existent default link", func(t *testing.T) {
		nonExistentDash := &malak.Dashboard{ID: uuid.New()}
		_, err := repo.DefaultLink(context.Background(), nonExistentDash)
		assert.ErrorIs(t, err, malak.ErrDashboardLinkNotFound)
	})
}

func TestDashboardLinkRepo_PublicDetails(t *testing.T) {
	db, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	repo := NewDashboardLinkRepo(db)
	workspaceRepo := NewWorkspaceRepository(db)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}
	_, err = db.NewInsert().Model(dashboard).Exec(context.Background())
	require.NoError(t, err)

	link := &malak.DashboardLink{
		DashboardID: dashboard.ID,
		Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
		LinkType:    malak.DashboardLinkType("default"),
		ExpiresAt:   &time.Time{},
	}
	_, err = db.NewInsert().Model(link).Exec(context.Background())
	require.NoError(t, err)

	t.Run("get dashboard with valid link", func(t *testing.T) {
		got, err := repo.PublicDetails(context.Background(), malak.Reference(link.Token))
		require.NoError(t, err)
		assert.Equal(t, dashboard.ID, got.ID)
		assert.Equal(t, dashboard.Title, got.Title)
	})

	t.Run("get dashboard with invalid link", func(t *testing.T) {
		invalidRef := malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink)
		_, err := repo.PublicDetails(context.Background(), invalidRef)
		assert.ErrorIs(t, err, malak.ErrDashboardLinkNotFound)
	})
}
