package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
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
