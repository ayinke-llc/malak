package postgres

import (
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	err := integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)
}

func TestIntegration_List(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
}

func TestIntegration_System(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	integrations, err := integrationRepo.System(t.Context())
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.System(t.Context())
	require.NoError(t, err)
	require.Len(t, integrations, 1)
}

func TestIntegration_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	_, err := integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeWorkspaceIntegration),
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrWorkspaceIntegrationNotFound, err)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	_, err = integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: integrations[0].Reference,
	})
	require.NoError(t, err)

	_, err = integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: integrations[0].Reference,
		ID:        integrations[0].ID,
	})
	require.NoError(t, err)
}

func TestIntegration_Disable(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]
	workspaceIntegration.IsEnabled = true

	require.NoError(t, integrationRepo.Update(t.Context(), &workspaceIntegration))

	updatedIntegration, err := integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)
	require.True(t, updatedIntegration.IsEnabled)

	err = integrationRepo.Disable(t.Context(), &workspaceIntegration)
	require.NoError(t, err)

	updatedIntegration, err = integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)
	require.False(t, updatedIntegration.IsEnabled)
}

func TestIntegration_Update(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test integration
	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	// Get the workspace integration
	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]

	// Verify initial state
	require.False(t, workspaceIntegration.IsEnabled)

	// Update the integration
	workspaceIntegration.IsEnabled = true
	err = integrationRepo.Update(t.Context(), &workspaceIntegration)
	require.NoError(t, err)

	// Fetch updated integration
	updatedIntegration, err := integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
		ID:        workspaceIntegration.ID,
	})
	require.NoError(t, err)
	require.True(t, updatedIntegration.IsEnabled)

	// Update again with different value
	updatedIntegration.IsEnabled = false
	err = integrationRepo.Update(t.Context(), updatedIntegration)
	require.NoError(t, err)

	// Verify the second update
	finalIntegration, err := integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
		ID:        workspaceIntegration.ID,
	})
	require.NoError(t, err)
	require.False(t, finalIntegration.IsEnabled)
}

func TestIntegration_CreateCharts(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Revenue Chart",
			InternalName:   "revenue_chart",
			ProviderID:     "stripe_revenue",
			ChartType:      malak.IntegrationChartTypeBar,
		},
		{
			UserFacingName: "Customer Growth",
			InternalName:   "customer_growth",
			ProviderID:     "stripe_customers",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}

	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	_, err = integrationRepo.Get(t.Context(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)

	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)
}

func TestIntegration_AddDataPoint(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Revenue Chart",
			InternalName:   "revenue_chart",
			ProviderID:     "stripe_revenue",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}

	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	dataPoints := []malak.IntegrationDataValues{
		{
			UserFacingName: "Revenue Chart",
			InternalName:   "revenue_chart",
			ProviderID:     "stripe_revenue",
			Data: malak.IntegrationDataPoint{
				PointName:     malak.GetTodayFormatted(),
				PointValue:    10050, // 100.50 * 100 to store as integer cents
				DataPointType: malak.IntegrationDataPointTypeCurrency,
				Metadata:      malak.IntegrationDataPointMetadata{},
			},
		},
	}

	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, dataPoints)
	require.NoError(t, err)

	invalidDataPoints := []malak.IntegrationDataValues{
		{
			UserFacingName: "Non Existent Chart",
			InternalName:   "non_existent_chart",
			ProviderID:     "stripe_revenue",
			Data: malak.IntegrationDataPoint{
				PointName:     malak.GetTodayFormatted(),
				PointValue:    20000, // 200.00 * 100 to store as integer cents
				DataPointType: malak.IntegrationDataPointTypeCurrency,
				Metadata:      malak.IntegrationDataPointMetadata{},
			},
		},
	}

	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, invalidDataPoints)
	require.Error(t, err)
}

func TestIntegration_ListCharts(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Initially there should be no charts
	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Empty(t, charts)

	err = integrationRepo.Create(t.Context(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]

	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Monthly Revenue",
			InternalName:   "monthly_revenue",
			ProviderID:     "stripe_monthly_revenue",
			ChartType:      malak.IntegrationChartTypeBar,
		},
		{
			UserFacingName: "Customer Count",
			InternalName:   "customer_count",
			ProviderID:     "stripe_customer_count",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}

	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	charts, err = integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 2)

	require.Contains(t, []string{charts[0].UserFacingName, charts[1].UserFacingName}, "Monthly Revenue")
	require.Contains(t, []string{charts[0].UserFacingName, charts[1].UserFacingName}, "Customer Count")
	require.Contains(t, []string{string(charts[0].InternalName), string(charts[1].InternalName)}, "monthly_revenue")
	require.Contains(t, []string{string(charts[0].InternalName), string(charts[1].InternalName)}, "customer_count")

	// verify workspace association
	for _, chart := range charts {
		require.Equal(t, workspace.ID, chart.WorkspaceID)
		require.Equal(t, workspaceIntegration.ID, chart.WorkspaceIntegrationID)
		require.NotEmpty(t, chart.Reference)
	}
}

func TestIntegration_GetChart(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	_, err = integrationRepo.GetChart(t.Context(), malak.FetchChartOptions{
		WorkspaceID: workspace.ID,
		Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegrationChart),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrChartNotFound)

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

	// Test getting chart by reference
	chart, err := integrationRepo.GetChart(t.Context(), malak.FetchChartOptions{
		WorkspaceID: workspace.ID,
		Reference:   charts[0].Reference,
	})
	require.NoError(t, err)
	require.Equal(t, "Account Balance", chart.UserFacingName)
	require.Equal(t, malak.IntegrationChartInternalNameTypeMercuryAccount, chart.InternalName)
	require.Equal(t, malak.IntegrationChartTypeBar, chart.ChartType)
	require.Equal(t, "account_123", chart.Metadata.ProviderID)

	// Test wrong workspace ID
	_, err = integrationRepo.GetChart(t.Context(), malak.FetchChartOptions{
		WorkspaceID: uuid.New(),
		Reference:   charts[0].Reference,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, malak.ErrChartNotFound)
}

func TestIntegration_CreateChartsDuplicate(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
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

	// Create charts with duplicate values
	chartValues := []malak.IntegrationChartValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			ChartType:      malak.IntegrationChartTypeBar,
		},
	}

	// First creation should succeed
	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	// Second creation should not create duplicates
	err = integrationRepo.CreateCharts(t.Context(), &workspaceIntegration, chartValues)
	require.NoError(t, err)

	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 1)
}

func TestIntegration_AddDataPointErrors(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create integration and chart
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

	// Try to add data point for non-existent chart
	dataPoints := []malak.IntegrationDataValues{
		{
			UserFacingName: "Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			Data: malak.IntegrationDataPoint{
				PointName:     "Balance",
				PointValue:    1000,
				DataPointType: malak.IntegrationDataPointTypeCurrency,
			},
		},
	}

	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, dataPoints)
	require.Error(t, err) // Should fail because chart doesn't exist

	// Create chart
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

	// Try to add data point with wrong provider ID
	dataPoints[0].ProviderID = "wrong_account"
	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, dataPoints)
	require.Error(t, err) // Should fail because provider ID doesn't match

	// Try to add data point with wrong workspace integration
	wrongWorkspaceIntegration := workspaceIntegration
	wrongWorkspaceIntegration.ID = uuid.New()
	err = integrationRepo.AddDataPoint(t.Context(), &wrongWorkspaceIntegration, dataPoints)
	require.Error(t, err)
}

func TestIntegration_ListChartsErrors(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	charts, err := integrationRepo.ListCharts(t.Context(), uuid.New())
	require.NoError(t, err)
	require.Empty(t, charts)

	charts, err = integrationRepo.ListCharts(t.Context(), uuid.Nil)
	require.NoError(t, err)
	require.Empty(t, charts)
}

func TestIntegration_GetDataPoints(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(t.Context(), &malak.FindWorkspaceOptions{
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

	integrations, err := integrationRepo.List(t.Context(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	workspaceIntegration := integrations[0]

	// Create chart
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

	// Get the created chart
	charts, err := integrationRepo.ListCharts(t.Context(), workspace.ID)
	require.NoError(t, err)
	require.Len(t, charts, 1)
	chart := charts[0]

	// Initially there should be no data points
	dataPoints, err := integrationRepo.GetDataPoints(t.Context(), chart)
	require.NoError(t, err)
	require.Empty(t, dataPoints)

	// Add data points
	dataPointValues := []malak.IntegrationDataValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			Data: malak.IntegrationDataPoint{
				PointName:     "Day 1",
				PointValue:    10000, // $100.00
				DataPointType: malak.IntegrationDataPointTypeCurrency,
				Metadata:      malak.IntegrationDataPointMetadata{},
			},
		},
	}
	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, dataPointValues)
	require.NoError(t, err)

	// Add another data point
	dataPointValues = []malak.IntegrationDataValues{
		{
			UserFacingName: "Account Balance",
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:     "account_123",
			Data: malak.IntegrationDataPoint{
				PointName:     "Day 2",
				PointValue:    20000, // $200.00
				DataPointType: malak.IntegrationDataPointTypeCurrency,
				Metadata:      malak.IntegrationDataPointMetadata{},
			},
		},
	}
	err = integrationRepo.AddDataPoint(t.Context(), &workspaceIntegration, dataPointValues)
	require.NoError(t, err)

	// Verify data points are returned in order
	dataPoints, err = integrationRepo.GetDataPoints(t.Context(), chart)
	require.NoError(t, err)
	require.Len(t, dataPoints, 2)

	// Verify data points are ordered by creation date
	require.Equal(t, int64(10000), dataPoints[0].PointValue)
	require.Equal(t, "Day 1", dataPoints[0].PointName)
	require.Equal(t, int64(20000), dataPoints[1].PointValue)
	require.Equal(t, "Day 2", dataPoints[1].PointName)

	// Verify data point fields
	for _, dp := range dataPoints {
		require.NotEmpty(t, dp.ID)
		require.Equal(t, workspaceIntegration.ID, dp.WorkspaceIntegrationID)
		require.Equal(t, workspace.ID, dp.WorkspaceID)
		require.Equal(t, chart.ID, dp.IntegrationChartID)
		require.NotEmpty(t, dp.Reference)
		require.Equal(t, malak.IntegrationDataPointTypeCurrency, dp.DataPointType)
		require.NotZero(t, dp.CreatedAt)
		require.NotZero(t, dp.UpdatedAt)
	}

	// Test with non-existent chart ID
	nonExistentChart := chart
	nonExistentChart.ID = uuid.New()
	dataPoints, err = integrationRepo.GetDataPoints(t.Context(), nonExistentChart)
	require.NoError(t, err)
	require.Empty(t, dataPoints)
}
