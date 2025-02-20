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
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
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
