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

	tests := []struct {
		name        string
		dashboard   *malak.Dashboard
		expectError bool
	}{
		{
			name: "valid dashboard",
			dashboard: &malak.Dashboard{
				WorkspaceID: workspace.ID,
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
				Title:       "Test Dashboard",
				Description: "Test Dashboard Description",
			},
			expectError: false,
		},
		{
			name: "dashboard with invalid workspace ID",
			dashboard: &malak.Dashboard{
				WorkspaceID: uuid.New(),
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
				Title:       "Test Dashboard",
				Description: "Test Dashboard Description",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dashboardRepo.Create(t.Context(), tt.dashboard)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tt.dashboard.ID)
				require.Equal(t, int64(0), tt.dashboard.ChartCount)
				require.False(t, tt.dashboard.CreatedAt.IsZero())
				require.False(t, tt.dashboard.UpdatedAt.IsZero())
			}
		})
	}
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

	tests := []struct {
		name        string
		chart       *malak.DashboardChart
		expectError bool
	}{
		{
			name: "valid chart",
			chart: &malak.DashboardChart{
				WorkspaceIntegrationID: workspaceIntegration.ID,
				ChartID:                createdChart.ID,
				Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
				WorkspaceID:            workspace.ID,
				DashboardID:            dashboard.ID,
			},
			expectError: false,
		},
		{
			name: "invalid workspace integration ID",
			chart: &malak.DashboardChart{
				WorkspaceIntegrationID: uuid.New(),
				ChartID:                createdChart.ID,
				Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
				WorkspaceID:            workspace.ID,
				DashboardID:            dashboard.ID,
			},
			expectError: true,
		},
		{
			name: "invalid chart ID",
			chart: &malak.DashboardChart{
				WorkspaceIntegrationID: workspaceIntegration.ID,
				ChartID:                uuid.New(),
				Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
				WorkspaceID:            workspace.ID,
				DashboardID:            dashboard.ID,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dashboardRepo.AddChart(t.Context(), tt.chart)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tt.chart.ID)

				// Verify chart count was incremented
				updatedDashboard, err := dashboardRepo.Get(t.Context(), malak.FetchDashboardOption{
					WorkspaceID: workspace.ID,
					Reference:   dashboard.Reference,
				})
				require.NoError(t, err)
				require.Equal(t, int64(1), updatedDashboard.ChartCount)
			}
		})
	}
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
		name        string
		opts        malak.FetchDashboardOption
		expectError bool
		errorType   error
	}{
		{
			name: "existing dashboard",
			opts: malak.FetchDashboardOption{
				WorkspaceID: workspace.ID,
				Reference:   dashboard.Reference,
			},
			expectError: false,
		},
		{
			name: "non-existent dashboard",
			opts: malak.FetchDashboardOption{
				WorkspaceID: workspace.ID,
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
			},
			expectError: true,
			errorType:   malak.ErrDashboardNotFound,
		},
		{
			name: "wrong workspace",
			opts: malak.FetchDashboardOption{
				WorkspaceID: uuid.New(),
				Reference:   dashboard.Reference,
			},
			expectError: true,
			errorType:   malak.ErrDashboardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dashboardRepo.Get(t.Context(), tt.opts)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					require.ErrorIs(t, err, tt.errorType)
				}
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

	// Add charts to dashboard
	for _, chart := range createdCharts {
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

	tests := []struct {
		name          string
		opts          malak.FetchDashboardChartsOption
		expectedCount int
	}{
		{
			name: "get existing charts",
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
			charts, err := dashboardRepo.GetCharts(t.Context(), tt.opts)
			require.NoError(t, err)
			require.Len(t, charts, tt.expectedCount)

			if tt.expectedCount > 0 {
				for _, chart := range charts {
					require.NotNil(t, chart.IntegrationChart)
					require.NotEmpty(t, chart.IntegrationChart.UserFacingName)
					require.NotEmpty(t, chart.IntegrationChart.InternalName)
					require.NotEmpty(t, chart.IntegrationChart.ChartType)
					require.Equal(t, workspaceIntegration.ID, chart.IntegrationChart.WorkspaceIntegrationID)
					require.Equal(t, workspace.ID, chart.IntegrationChart.WorkspaceID)
				}

				// Verify specific chart data
				foundAccountBalance := false
				foundTransactionHistory := false

				for _, chart := range charts {
					switch chart.IntegrationChart.InternalName {
					case malak.IntegrationChartInternalNameTypeMercuryAccount:
						foundAccountBalance = true
						require.Equal(t, "Account Balance", chart.IntegrationChart.UserFacingName)
						require.Equal(t, malak.IntegrationChartTypeBar, chart.IntegrationChart.ChartType)
					case malak.IntegrationChartInternalNameTypeMercuryAccountTransaction:
						foundTransactionHistory = true
						require.Equal(t, "Transaction History", chart.IntegrationChart.UserFacingName)
						require.Equal(t, malak.IntegrationChartTypeBar, chart.IntegrationChart.ChartType)
					}
				}

				require.True(t, foundAccountBalance, "Account Balance chart not found")
				require.True(t, foundTransactionHistory, "Transaction History chart not found")
			}
		})
	}
}

func TestDashboard_List(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)

	workspace1, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	workspace2, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("c12da796-9362-4c70-b2cb-fc8a1eba2526"),
	})
	require.NoError(t, err)

	// Clean up existing dashboards
	_, err = client.NewDelete().Model((*malak.Dashboard)(nil)).Where("1=1").Exec(t.Context())
	require.NoError(t, err)

	// Create test dashboards
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
	}

	for _, d := range dashboards2 {
		err = dashboardRepo.Create(t.Context(), d)
		require.NoError(t, err)
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
		{
			name: "non-existent workspace",
			opts: malak.ListDashboardOptions{
				WorkspaceID: uuid.New(),
				Paginator: malak.Paginator{
					Page:    1,
					PerPage: 10,
				},
			},
			expectedCount: 0,
			totalCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, total, err := dashboardRepo.List(t.Context(), tt.opts)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, len(results))
			require.Equal(t, tt.totalCount, total)

			if tt.expectedCount > 0 {
				for _, result := range results {
					require.Equal(t, tt.opts.WorkspaceID, result.WorkspaceID)
					require.NotEmpty(t, result.ID)
					require.NotEmpty(t, result.Reference)
					require.NotEmpty(t, result.Title)
					require.NotEmpty(t, result.Description)
					require.False(t, result.CreatedAt.IsZero())
					require.False(t, result.UpdatedAt.IsZero())
				}
			}
		})
	}
}

func TestDashboard_UpdatePositions(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	dashboardRepo := NewDashboardRepo(client)
	workspaceRepo := NewWorkspaceRepository(client)
	integrationRepo := NewIntegrationRepo(client)

	workspace, err := workspaceRepo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create integration
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

	// Get workspace integration
	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	workspaceIntegration := integrations[0]

	// Create integration charts
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

	// Get created charts
	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 2)

	dashboard := &malak.Dashboard{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboard),
		Title:       "Test Dashboard",
		Description: "Test Dashboard Description",
	}

	err = dashboardRepo.Create(t.Context(), dashboard)
	require.NoError(t, err)

	// Create dashboard charts
	dashboardCharts := make([]*malak.DashboardChart, 0, len(charts))
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
		dashboardCharts = append(dashboardCharts, dashboardChart)
	}

	tests := []struct {
		name          string
		positions     []malak.DashboardChartPosition
		expectedCount int
		expectedOrder []int64
		expectError   bool
		dashboardID   uuid.UUID
	}{
		{
			name: "update with valid positions",
			positions: []malak.DashboardChartPosition{
				{
					DashboardID: dashboard.ID,
					ChartID:     dashboardCharts[0].ID,
					OrderIndex:  1,
				},
				{
					DashboardID: dashboard.ID,
					ChartID:     dashboardCharts[1].ID,
					OrderIndex:  2,
				},
			},
			expectedCount: 2,
			expectedOrder: []int64{1, 2},
			dashboardID:   dashboard.ID,
		},
		{
			name:          "update with empty positions",
			positions:     []malak.DashboardChartPosition{},
			expectedCount: 0,
			expectedOrder: []int64{},
			dashboardID:   dashboard.ID,
		},
		{
			name: "update with reversed order",
			positions: []malak.DashboardChartPosition{
				{
					DashboardID: dashboard.ID,
					ChartID:     dashboardCharts[1].ID,
					OrderIndex:  1,
				},
				{
					DashboardID: dashboard.ID,
					ChartID:     dashboardCharts[0].ID,
					OrderIndex:  2,
				},
			},
			expectedCount: 2,
			expectedOrder: []int64{1, 2},
			dashboardID:   dashboard.ID,
		},
		{
			name: "update with invalid dashboard ID",
			positions: []malak.DashboardChartPosition{
				{
					DashboardID: uuid.New(),
					ChartID:     dashboardCharts[0].ID,
					OrderIndex:  1,
				},
			},
			expectedCount: 0,
			expectedOrder: []int64{},
			expectError:   true,
			dashboardID:   uuid.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update positions
			err := dashboardRepo.UpdateDashboardPositions(t.Context(), tt.dashboardID, tt.positions)
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify positions were updated
			savedPositions, err := dashboardRepo.GetDashboardPositions(t.Context(), tt.dashboardID)
			require.NoError(t, err)
			require.Len(t, savedPositions, tt.expectedCount)

			// Verify position order if there are expected positions
			if tt.expectedCount > 0 {
				for i, expectedOrder := range tt.expectedOrder {
					require.Equal(t, expectedOrder, savedPositions[i].OrderIndex)
				}
			}
		})
	}
}

func TestDashboard_RemoveChart(t *testing.T) {
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

	dashboardChart := &malak.DashboardChart{
		WorkspaceIntegrationID: workspaceIntegration.ID,
		ChartID:                createdChart.ID,
		Reference:              malak.NewReferenceGenerator().Generate(malak.EntityTypeDashboardChart),
		WorkspaceID:            workspace.ID,
		DashboardID:            dashboard.ID,
	}
	err = dashboardRepo.AddChart(t.Context(), dashboardChart)
	require.NoError(t, err)

	positions := []malak.DashboardChartPosition{
		{
			DashboardID: dashboard.ID,
			ChartID:     dashboardChart.ID,
			OrderIndex:  1,
		},
	}
	err = dashboardRepo.UpdateDashboardPositions(t.Context(), dashboard.ID, positions)
	require.NoError(t, err)

	tests := []struct {
		name        string
		dashboardID uuid.UUID
		chartID     uuid.UUID
		expectError bool
	}{
		{
			name:        "successfully remove chart",
			dashboardID: dashboard.ID,
			chartID:     createdChart.ID,
			expectError: false,
		},
		{
			name:        "non-existent chart",
			dashboardID: dashboard.ID,
			chartID:     uuid.New(),
			expectError: true,
		},
		{
			name:        "non-existent dashboard",
			dashboardID: uuid.New(),
			chartID:     createdChart.ID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dashboardRepo.RemoveChart(t.Context(), tt.dashboardID, tt.chartID)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// verify chart was removed
				updatedDashboard, err := dashboardRepo.Get(t.Context(), malak.FetchDashboardOption{
					WorkspaceID: workspace.ID,
					Reference:   dashboard.Reference,
				})
				require.NoError(t, err)
				require.Equal(t, int64(0), updatedDashboard.ChartCount)

				// verify chart positions were removed
				positions, err := dashboardRepo.GetDashboardPositions(t.Context(), dashboard.ID)
				require.NoError(t, err)
				require.Empty(t, positions)

				// verify chart is no longer in dashboard charts
				charts, err := dashboardRepo.GetCharts(t.Context(), malak.FetchDashboardChartsOption{
					WorkspaceID: workspace.ID,
					DashboardID: dashboard.ID,
				})
				require.NoError(t, err)
				require.Empty(t, charts)
			}
		})
	}
}
