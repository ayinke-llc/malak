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
