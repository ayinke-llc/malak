package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func (i *integrationRepo) ToggleEnabled(ctx context.Context,
	integration *malak.WorkspaceIntegration) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return i.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewUpdate().
				Where("id = ?", integration.ID).
				Set("is_enabled = CASE WHEN is_enabled = true THEN false ELSE true END").
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
				On("CONFLICT (internal_name,workspace_id,workspace_integration_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// update the integration too
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

	// generator := malak.NewReferenceGenerator()
	//
	// return i.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
	//
	// 	for _, value := range dataPoints {
	//
	// 		var chart malak.IntegrationChart
	//
	// 		err := tx.NewSelect().
	// 			Where("workspace_integration_id = ?", workspaceIntegration.ID).
	// 			Where("workspace_id = ?", workspaceIntegration.WorkspaceID).
	// 			Where("internal_name = ?", value.InternalName).
	// 			Scan(ctx, &chart)
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			chart = malak.IntegrationChart{
	// 				WorkspaceIntegrationID: workspaceIntegration.ID,
	// 				WorkspaceID:            workspaceIntegration.WorkspaceID,
	// 				Reference:              generator.Generate(malak.EntityTypeIntegrationChart),
	// 				UserFacingName:         value.InternalName.String(),
	// 			}
	// 			tx.NewInsert().
	// 				Exec(ctx, &chart)
	// 		}
	//
	// 	}
	//
	// 	return nil
	// })
	return nil
}
