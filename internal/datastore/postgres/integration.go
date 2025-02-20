package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type integrationRepo struct {
	inner *bun.DB
}

func NewIntegrationRepo(db *bun.DB) malak.IntegrationRepository {
	return &integrationRepo{
		inner: db,
	}
}

func (i *integrationRepo) Create(ctx context.Context,
	integration *malak.Integration) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return i.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().
				Model(integration).
				Exec(ctx)
			if err != nil {
				return err
			}

			workspaces := make([]*malak.Workspace, 0)

			err = tx.NewSelect().
				Model(&workspaces).
				Scan(ctx)
			if err != nil {
				return err
			}

			var workspaceIntegrations = make([]*malak.WorkspaceIntegration, 0, len(workspaces))

			gen := malak.NewReferenceGenerator()

			for _, workspace := range workspaces {
				workspaceIntegrations = append(workspaceIntegrations, &malak.WorkspaceIntegration{
					WorkspaceID:   workspace.ID,
					Reference:     gen.Generate(malak.EntityTypeWorkspaceIntegration),
					IntegrationID: integration.ID,
					IsEnabled:     false,
				})
			}

			_, err = tx.NewInsert().
				Model(&workspaceIntegrations).
				On("CONFLICT (workspace_id,integration_id) DO NOTHING").
				Exec(ctx)
			return err
		})
}

func (i *integrationRepo) List(ctx context.Context,
	workspace *malak.Workspace) ([]malak.WorkspaceIntegration, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	integrations := make([]malak.WorkspaceIntegration, 0)

	return integrations, i.inner.NewSelect().
		Model(&integrations).
		Where("workspace_id = ?", workspace.ID).
		Order("created_at ASC").
		Relation("Integration").
		Scan(ctx)
}

func (i *integrationRepo) System(ctx context.Context) ([]malak.Integration, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	integrations := make([]malak.Integration, 0)

	return integrations, i.inner.NewSelect().
		Model(&integrations).
		Order("created_at ASC").
		Scan(ctx)
}

func (i *integrationRepo) Get(ctx context.Context,
	opts malak.FindWorkspaceIntegrationOptions) (*malak.WorkspaceIntegration, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	integration := &malak.WorkspaceIntegration{}

	sel := i.inner.NewSelect().Model(integration).
		Where("workspace_integration.reference = ?", opts.Reference)

	if opts.ID != uuid.Nil {
		sel = sel.Where("workspace_integration.id = ?", opts.ID)
	}

	err := sel.Relation("Integration").
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrWorkspaceIntegrationNotFound
	}

	return integration, err
}

func (i *integrationRepo) Disable(ctx context.Context,
	integration *malak.WorkspaceIntegration) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return i.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			integration.UpdatedAt = time.Now()

			_, err := tx.NewUpdate().
				Where("id = ?", integration.ID).
				Set("is_enabled = ?", false).
				Model(integration).
				Exec(ctx)

			return err
		})
}

func (i *integrationRepo) Update(ctx context.Context,
	integration *malak.WorkspaceIntegration) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	integration.UpdatedAt = time.Now()

	_, err := i.inner.NewUpdate().
		Where("id = ?", integration.ID).
		Model(integration).
		Exec(ctx)
	return err
}

func (i *integrationRepo) CreateCharts(ctx context.Context,
	workspaceIntegration *malak.WorkspaceIntegration,
	chartValues []malak.IntegrationChartValues) error {

	generator := malak.NewReferenceGenerator()

	return i.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		for _, value := range chartValues {

			chart := &malak.IntegrationChart{
				WorkspaceIntegrationID: workspaceIntegration.ID,
				WorkspaceID:            workspaceIntegration.WorkspaceID,
				Reference:              generator.Generate(malak.EntityTypeIntegrationChart),
				UserFacingName:         value.UserFacingName,
				InternalName:           value.InternalName,
				Metadata: malak.IntegrationChartMetadata{
					ProviderID: value.ProviderID,
				},
			}

			_, err := tx.NewInsert().Model(chart).
				On("CONFLICT (user_facing_name,internal_name,workspace_id,workspace_integration_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		workspaceIntegration.UpdatedAt = time.Now()

		_, err := tx.NewUpdate().
			Where("id = ?", workspaceIntegration.ID).
			Model(workspaceIntegration).
			Exec(ctx)
		return err
	})
}

func (i *integrationRepo) AddDataPoint(ctx context.Context,
	workspaceIntegration *malak.WorkspaceIntegration,
	dataPoints []malak.IntegrationDataValues) error {

	generator := malak.NewReferenceGenerator()

	return i.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		for _, value := range dataPoints {

			var chart malak.IntegrationChart

			query := tx.NewSelect().
				Model(&chart).
				Where("workspace_integration_id = ?", workspaceIntegration.ID).
				Where("workspace_id = ?", workspaceIntegration.WorkspaceID).
				Where("internal_name = ?", value.InternalName)

			if !hermes.IsStringEmpty(value.ProviderID) {
				query = query.Where("metadata->>'provider_id' = ?", value.ProviderID)
			}

			err := query.Scan(ctx)
			if err != nil {
				return err
			}

			value.Data.Reference = generator.Generate(malak.EntityTypeIntegrationDatapoint)
			value.Data.WorkspaceIntegrationID = workspaceIntegration.ID
			value.Data.WorkspaceID = workspaceIntegration.WorkspaceID
			value.Data.IntegrationChartID = chart.ID

			_, err = tx.NewInsert().
				Model(&value.Data).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *integrationRepo) ListCharts(ctx context.Context,
	workspaceID uuid.UUID) ([]malak.IntegrationChart, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	charts := make([]malak.IntegrationChart, 0)

	return charts, i.inner.NewSelect().
		Model(&charts).
		Where("workspace_id = ?", workspaceID).
		Order("created_at ASC").
		Scan(ctx)
}
