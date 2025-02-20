package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDashboard_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}

	err = dashboardRepo.Create(context.Background(), dashboard)
	require.NoError(t, err)
	require.NotEmpty(t, dashboard.ID)
	require.Equal(t, int64(0), dashboard.ChartCount)
}

// func TestDashboard_AddChart(t *testing.T) {
// 	client, teardownFunc := setupDatabase(t)
// 	defer teardownFunc()
//
// 	dashboardRepo := NewDashboardRepo(client)
// 	workspaceRepo := NewWorkspaceRepository(client)
//
// 	workspace, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
// 		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
// 	})
// 	require.NoError(t, err)
//
// 	dashboard := &malak.Dashboard{
// 		WorkspaceID: workspace.ID,
// 		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
// 		Title:       "Test Dashboard",
// 		Description: "Test Dashboard Description",
// 	}
//
// 	err = dashboardRepo.Create(context.Background(), dashboard)
// 	require.NoError(t, err)
// 	require.Equal(t, int64(0), dashboard.ChartCount)
//
// 	workspaceIntegrationID := uuid.New()
//
// 	chart := &malak.DashboardChart{
// 		WorkspaceIntegrationID: workspaceIntegrationID,
// 		Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
// 		WorkspaceID:            workspace.ID,
// 		DashboardID:            dashboard.ID,
// 		DashboardType:          malak.DashboardChartTypeBarchart,
// 	}
//
// 	err = dashboardRepo.AddChart(context.Background(), chart)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, chart.ID)
// }

func TestDashboard_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// Clean up any existing dashboards first
	_, err := client.NewDelete().
		Table("dashboards").
		Where("1=1").
		Exec(context.Background())
	require.NoError(t, err)

	workspace1, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	workspace2, err := workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	dashboards1 := []*malak.Dashboard{
		{
			WorkspaceID: workspace1.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			Title:       "First Dashboard",
			Description: "First Dashboard Description",
		},
		{
			WorkspaceID: workspace1.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			Title:       "Second Dashboard",
			Description: "Second Dashboard Description",
		},
		{
			WorkspaceID: workspace1.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			Title:       "Third Dashboard",
			Description: "Third Dashboard Description",
		},
	}

	dashboards2 := []*malak.Dashboard{
		{
			WorkspaceID: workspace2.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			Title:       "First Dashboard Workspace 2",
			Description: "First Dashboard Description Workspace 2",
		},
		{
			WorkspaceID: workspace2.ID,
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			Title:       "Second Dashboard Workspace 2",
			Description: "Second Dashboard Description Workspace 2",
		},
	}

	for _, d := range dashboards1 {
		err = dashboardRepo.Create(context.Background(), d)
		require.NoError(t, err)
		require.NotEmpty(t, d.ID)
	}

	for _, d := range dashboards2 {
		err = dashboardRepo.Create(context.Background(), d)
		require.NoError(t, err)
		require.NotEmpty(t, d.ID)
	}

	tests := []struct {
		name          string
		opts          malak.ListDashboardOptions
		expectedCount int
		totalCount    int64
	}{
		{
			name: "first page workspace 1",
			opts: malak.ListDashboardOptions{
				WorkspaceID: workspace1.ID,
				Paginator: malak.Paginator{
					Page:    1,
					PerPage: 2,
				},
			},
			expectedCount: 2,
			totalCount:    3,
		},
		{
			name: "second page workspace 1",
			opts: malak.ListDashboardOptions{
				WorkspaceID: workspace1.ID,
				Paginator: malak.Paginator{
					Page:    2,
					PerPage: 2,
				},
			},
			expectedCount: 1,
			totalCount:    3,
		},
		{
			name: "all items workspace 1",
			opts: malak.ListDashboardOptions{
				WorkspaceID: workspace1.ID,
				Paginator: malak.Paginator{
					Page:    1,
					PerPage: 10,
				},
			},
			expectedCount: 3,
			totalCount:    3,
		},
		{
			name: "all items workspace 2",
			opts: malak.ListDashboardOptions{
				WorkspaceID: workspace2.ID,
				Paginator: malak.Paginator{
					Page:    1,
					PerPage: 10,
				},
			},
			expectedCount: 2,
			totalCount:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := dashboardRepo.List(context.Background(), tt.opts)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, len(results))
			require.Equal(t, tt.totalCount, total)

			for _, result := range results {
				require.Equal(t, tt.opts.WorkspaceID, result.WorkspaceID)
			}
		})
	}

	nonExistentResults, total, err := dashboardRepo.List(context.Background(), malak.ListDashboardOptions{
		WorkspaceID: uuid.New(),
		Paginator: malak.Paginator{
			Page:    1,
			PerPage: 10,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(nonExistentResults))
	require.Equal(t, int64(0), total)
}
