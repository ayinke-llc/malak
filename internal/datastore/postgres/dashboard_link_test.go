package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
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
					Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
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
					Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
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
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify link was created
			var link malak.DashboardLink
			err = db.NewSelect().
				Model(&link).
				Where("token = ?", tt.opts.Link.Token).
				Scan(context.Background())
			require.NoError(t, err)

			require.Equal(t, tt.opts.Link.DashboardID, link.DashboardID)
			require.Equal(t, tt.opts.Link.LinkType, link.LinkType)
			require.Equal(t, tt.opts.Link.Reference, link.Reference)

			if tt.opts.Email != "" {
				// Verify contact was created and linked
				var contact malak.Contact
				err = db.NewSelect().
					Model(&contact).
					Where("email = ?", tt.opts.Email).
					Scan(context.Background())
				require.NoError(t, err)
				require.Equal(t, contact.ID, link.ContactID)
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
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
		Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
		LinkType:    malak.DashboardLinkType("default"),
		ExpiresAt:   &time.Time{},
	}
	_, err = db.NewInsert().Model(link).Exec(context.Background())
	require.NoError(t, err)

	t.Run("get existing default link", func(t *testing.T) {
		got, err := repo.DefaultLink(context.Background(), dashboard)
		require.NoError(t, err)
		require.Equal(t, link.Token, got.Token)
		require.Equal(t, link.DashboardID, got.DashboardID)
		require.Equal(t, link.Reference, got.Reference)
	})

	t.Run("get non-existent default link", func(t *testing.T) {
		nonExistentDash := &malak.Dashboard{ID: uuid.New()}
		_, err := repo.DefaultLink(context.Background(), nonExistentDash)
		require.ErrorIs(t, err, malak.ErrDashboardLinkNotFound)
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
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
		Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
		LinkType:    malak.DashboardLinkType("default"),
		ExpiresAt:   &time.Time{},
	}
	_, err = db.NewInsert().Model(link).Exec(context.Background())
	require.NoError(t, err)

	t.Run("get dashboard with valid link", func(t *testing.T) {
		got, err := repo.PublicDetails(context.Background(), malak.Reference(link.Token))
		require.NoError(t, err)
		require.Equal(t, dashboard.ID, got.ID)
		require.Equal(t, dashboard.Title, got.Title)
	})

	t.Run("get dashboard with invalid link", func(t *testing.T) {
		invalidRef := malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink)
		_, err := repo.PublicDetails(context.Background(), invalidRef)
		require.ErrorIs(t, err, malak.ErrDashboardLinkNotFound)
	})
}

func TestDashboardLinkRepo_List(t *testing.T) {
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

	links := []*malak.DashboardLink{
		{
			DashboardID: dashboard.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
			Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
			LinkType:    malak.DashboardLinkType("default"),
			ExpiresAt:   &time.Time{},
		},
		{
			DashboardID: dashboard.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
			Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
			LinkType:    malak.DashboardLinkType("contact"),
			ExpiresAt:   &time.Time{},
		},
		{
			DashboardID: dashboard.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
			Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
			LinkType:    malak.DashboardLinkType("contact"),
			ExpiresAt:   &time.Time{},
		},
	}

	for _, link := range links {
		_, err = db.NewInsert().Model(link).Exec(context.Background())
		require.NoError(t, err)
	}

	t.Run("list all dashboard links", func(t *testing.T) {
		opts := malak.ListAccessControlOptions{
			DashboardID: dashboard.ID,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		}

		got, count, err := repo.List(context.Background(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(2), count)
		require.Len(t, got, len(links))
	})

	t.Run("list with pagination", func(t *testing.T) {
		opts := malak.ListAccessControlOptions{
			DashboardID: dashboard.ID,
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 2,
			},
		}

		got, count, err := repo.List(context.Background(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(2), count)
		require.Len(t, got, 2)
	})

	t.Run("list for non-existent dashboard", func(t *testing.T) {
		opts := malak.ListAccessControlOptions{
			DashboardID: uuid.New(),
			Paginator: malak.Paginator{
				Page:    1,
				PerPage: 10,
			},
		}

		got, count, err := repo.List(context.Background(), opts)
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
		require.Empty(t, got)
	})
}

func TestDashboardLinkRepo_Delete(t *testing.T) {
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
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink),
		Token:       malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink).String(),
		LinkType:    malak.DashboardLinkType("default"),
		ExpiresAt:   &time.Time{},
	}
	_, err = db.NewInsert().Model(link).Exec(context.Background())
	require.NoError(t, err)

	t.Run("delete existing link", func(t *testing.T) {
		err := repo.Delete(context.Background(), *dashboard, link.Reference)
		require.NoError(t, err)

		var count int
		count, err = db.NewSelect().Model((*malak.DashboardLink)(nil)).Where("reference = ?", link.Reference).Count(context.Background())
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("delete non-existent link", func(t *testing.T) {
		nonExistentRef := malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardLink)
		err := repo.Delete(context.Background(), *dashboard, nonExistentRef)
		require.NoError(t, err)
	})
}
