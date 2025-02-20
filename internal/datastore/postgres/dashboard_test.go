package postgres

import (
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

	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)
	require.NotEmpty(t, dashboard.ID)
	require.Equal(t, int64(0), dashboard.ChartCount)
}

func TestDashboard_AddChart(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)
	integrationRepo := NewIntegrationRepo(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integration := &malak.Integration{
		IntegrationName: "Mercury",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Mercury Banking Integration",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://mercury.com/logo.png",
	}
	err = integrationRepo.Create(t.Context(), integration)
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}
	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 1)
	createdChart := charts[0]

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}

	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)
	require.Equal(t, int64(0), dashboard.ChartCount)

	chart := &malak.DashboardChart{
		WorkspaceIntegrationID: workspaceIntegration.ID,
		ChartID:                createdChart.ID,
		Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
		WorkspaceID:            workspace.ID,
		DashboardID:            dashboard.ID,
	}

	err = dashboardRepo.AddChart(t.Context(), chart)
	require.NoError(t, err)
	require.NotEmpty(t, chart.ID)

	// verify chart count was incremented
	updatedDashboard, err := dashboardRepo.Get(t.Context(), malak.FetchDashboardOption{
		WorkspaceID: workspace.ID,
		Reference:   dashboard.Reference,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), updatedDashboard.ChartCount)
}

func TestDashboard_Get(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}

	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)

	tests := []struct {
		name          string
		opts          malak.FetchDashboardOption
		expectedError error
	}{
		{
			name: "existing dashboard",
			opts: malak.FetchDashboardOption{
				WorkspaceID: workspace.ID,
				Reference:   dashboard.Reference,
			},
			expectedError: nil,
		},
		{
			name: "non-existent dashboard",
			opts: malak.FetchDashboardOption{
				WorkspaceID: workspace.ID,
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			},
			expectedError: malak.ErrDashboardNotFound,
		},
		{
			name: "wrong workspace",
			opts: malak.FetchDashboardOption{
				WorkspaceID: uuid.New(),
				Reference:   dashboard.Reference,
			},
			expectedError: malak.ErrDashboardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dashboardRepo.Get(t.Context(), tt.opts)
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.opts.WorkspaceID, result.WorkspaceID)
				require.Equal(t, tt.opts.Reference, result.Reference)
				require.Equal(t, dashboard.Title, result.Title)
				require.Equal(t, dashboard.Description, result.Description)
			}
		})
	}
}

func TestDashboard_GetCharts(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)
	integrationRepo := NewIntegrationRepo(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integration := &malak.Integration{
		IntegrationName: "Mercury",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Mercury Banking Integration",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://mercury.com/logo.png",
	}
	err = integrationRepo.Create(t.Context(), integration)
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
		{
			UserFacingName: "Transaction History",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccountTransaction,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}
	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	createdCharts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, createdCharts, 2)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}

	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)

	// Add multiple charts
	charts := []*malak.DashboardChart{
		{
			WorkspaceIntegrationID: workspaceIntegration.ID,
			ChartID:                createdCharts[0].ID,
			Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
			WorkspaceID:            workspace.ID,
			DashboardID:            dashboard.ID,
		},
		{
			WorkspaceIntegrationID: workspaceIntegration.ID,
			ChartID:                createdCharts[1].ID,
			Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
			WorkspaceID:            workspace.ID,
			DashboardID:            dashboard.ID,
		},
	}

	for _, chart := range charts {
		err = dashboardRepo.AddChart(t.Context(), chart)
		require.NoError(t, err)
	}

	// Test getting charts
	tests := []struct {
		name          string
		opts          malak.FetchDashboardChartsOption
		expectedCount int
	}{
		{
			name: "existing dashboard charts",
			opts: malak.FetchDashboardChartsOption{
				WorkspaceID: workspace.ID,
				DashboardID: dashboard.ID,
			},
			expectedCount: 2,
		},
		{
			name: "non-existent dashboard",
			opts: malak.FetchDashboardChartsOption{
				WorkspaceID: workspace.ID,
				DashboardID: uuid.New(),
			},
			expectedCount: 0,
		},
		{
			name: "wrong workspace",
			opts: malak.FetchDashboardChartsOption{
				WorkspaceID: uuid.New(),
				DashboardID: dashboard.ID,
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := dashboardRepo.GetCharts(t.Context(), tt.opts)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, len(results))

			if tt.expectedCount > 0 {
				for _, result := range results {
					require.Equal(t, tt.opts.WorkspaceID, result.WorkspaceID)
					require.Equal(t, tt.opts.DashboardID, result.DashboardID)
					require.NotEmpty(t, result.ID)
					require.NotEmpty(t, result.Reference)
				}
			}
		})
	}
}

func TestDashboard_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	// Clean up any existing dashboards first
	_, err := client.NewDelete().
		Table("dashboards").
		Where("1=1").
		Exec(t.Context())
	require.NoError(t, err)

	workspace1, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	workspace2, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
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
		err = dashboardRepo.Create(t.Context(), d)
		require.NoError(t, err)
		require.NotEmpty(t, d.ID)
	}

	for _, d := range dashboards2 {
		err = dashboardRepo.Create(t.Context(), d)
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
			results, total, err := dashboardRepo.List(t.Context(), tt.opts)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, len(results))
			require.Equal(t, tt.totalCount, total)

			for _, result := range results {
				require.Equal(t, tt.opts.WorkspaceID, result.WorkspaceID)
			}
		})
	}

	nonExistentResults, total, err := dashboardRepo.List(t.Context(), malak.ListDashboardOptions{
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

func TestDashboard_GetChartsWithIntegration(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)
	integrationRepo := NewIntegrationRepo(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integration := &malak.Integration{
		IntegrationName: "Mercury",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Mercury Banking Integration",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://mercury.com/logo.png",
	}
	err = integrationRepo.Create(t.Context(), integration)
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
		{
			UserFacingName: "Transaction History",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccountTransaction,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}
	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 2)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Mercury Dashboard",
		Description: "Mercury Banking Dashboard",
	}
	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)

	for _, chart := range charts {
		dashboardChart := &malak.DashboardChart{
			WorkspaceIntegrationID: workspaceIntegration.ID,
			ChartID:                chart.ID,
			Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
			WorkspaceID:            workspace.ID,
			DashboardID:            dashboard.ID,
		}
		err = dashboardRepo.AddChart(t.Context(), dashboardChart)
		require.NoError(t, err)
	}

	dashboardCharts, err := dashboardRepo.GetCharts(t.Context(), malak.FetchDashboardChartsOption{
		WorkspaceID: workspace.ID,
		DashboardID: dashboard.ID,
	})
	require.NoError(t, err)
	require.Len(t, dashboardCharts, 2)

	// verify integration chart data is loaded
	for _, dashboardChart := range dashboardCharts {
		require.NotNil(t, dashboardChart.IntegrationChart)
		require.NotEmpty(t, dashboardChart.IntegrationChart.UserFacingName)
		require.NotEmpty(t, dashboardChart.IntegrationChart.InternalName)
		require.NotEmpty(t, dashboardChart.IntegrationChart.ChartType)
		require.Equal(t, workspaceIntegration.ID, dashboardChart.IntegrationChart.WorkspaceIntegrationID)
		require.Equal(t, workspace.ID, dashboardChart.IntegrationChart.WorkspaceID)
	}

	// vErify specific chart data
	foundAccountBalance := false
	foundTransactionHistory := false

	for _, dashboardChart := range dashboardCharts {
		switch dashboardChart.IntegrationChart.InternalName {
		case malak.IntegrationChartInternalNameTypeMercuryAccount:
			foundAccountBalance = true
			require.Equal(t, "Account Balance", dashboardChart.IntegrationChart.UserFacingName)
			require.Equal(t, malak.IntegrationChartTypeBar, dashboardChart.IntegrationChart.ChartType)
		case malak.IntegrationChartInternalNameTypeMercuryAccountTransaction:
			foundTransactionHistory = true
			require.Equal(t, "Transaction History", dashboardChart.IntegrationChart.UserFacingName)
			require.Equal(t, malak.IntegrationChartTypeBar, dashboardChart.IntegrationChart.ChartType)
		}
	}

	require.True(t, foundAccountBalance, "Account Balance chart not found")
	require.True(t, foundTransactionHistory, "Transaction History chart not found")
}
