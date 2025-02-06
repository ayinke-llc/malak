package postgres

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Create(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	err := integrationRepo.Create(context.Background(), &malak.Integration{
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

	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(context.Background(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
}

func TestIntegration_System(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)

	integrations, err := integrationRepo.System(context.Background())
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(context.Background(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.System(context.Background())
	require.NoError(t, err)
	require.Len(t, integrations, 1)
}

func TestIntegration_Get(t *testing.T) {

	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	_, err := integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeWorkspaceIntegration),
	})
	require.Error(t, err)
	require.Equal(t, malak.ErrWorkspaceIntegrationNotFound, err)

	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	integrations, err := integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 0)

	err = integrationRepo.Create(context.Background(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	integrations, err = integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	_, err = integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: integrations[0].Reference,
	})
	require.NoError(t, err)

	_, err = integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: integrations[0].Reference,
		ID:        integrations[0].ID,
	})
	require.NoError(t, err)
}

func TestIntegration_ToggleEnabled(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test integration
	err = integrationRepo.Create(context.Background(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	// Get the workspace integration
	integrations, err := integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]
	initialEnabledState := workspaceIntegration.IsEnabled

	// Toggle enabled state
	err = integrationRepo.ToggleEnabled(context.Background(), &workspaceIntegration)
	require.NoError(t, err)

	// Fetch updated integration
	updatedIntegration, err := integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)
	require.NotEqual(t, initialEnabledState, updatedIntegration.IsEnabled)

	// Toggle again to verify it switches back
	err = integrationRepo.ToggleEnabled(context.Background(), updatedIntegration)
	require.NoError(t, err)

	finalIntegration, err := integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)
	require.Equal(t, initialEnabledState, finalIntegration.IsEnabled)
}

func TestIntegration_Update(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	integrationRepo := NewIntegrationRepo(client)
	repo := NewWorkspaceRepository(client)

	workspace, err := repo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: uuid.MustParse("a4ae79a2-9b76-40d7-b5a1-661e60a02cb0"),
	})
	require.NoError(t, err)

	// Create a test integration
	err = integrationRepo.Create(context.Background(), &malak.Integration{
		IntegrationName: "Stripe",
		Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
		Description:     "Stripe stripe stripe",
		IsEnabled:       true,
		IntegrationType: malak.IntegrationTypeOauth2,
		LogoURL:         "https://google.com",
	})
	require.NoError(t, err)

	// Get the workspace integration
	integrations, err := integrationRepo.List(context.Background(), workspace)
	require.NoError(t, err)
	require.Len(t, integrations, 1)

	workspaceIntegration := integrations[0]
	initialUpdateTime := workspaceIntegration.UpdatedAt

	// Update the integration
	workspaceIntegration.IsEnabled = true
	err = integrationRepo.Update(context.Background(), &workspaceIntegration)
	require.NoError(t, err)

	// Fetch updated integration
	updatedIntegration, err := integrationRepo.Get(context.Background(), malak.FindWorkspaceIntegrationOptions{
		Reference: workspaceIntegration.Reference,
	})
	require.NoError(t, err)
	require.True(t, updatedIntegration.IsEnabled)
	require.True(t, updatedIntegration.UpdatedAt.After(initialUpdateTime))
}
