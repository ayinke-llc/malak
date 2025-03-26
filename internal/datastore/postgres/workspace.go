package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type workspaceRepo struct {
	inner *bun.DB
}

func NewWorkspaceRepository(inner *bun.DB) malak.WorkspaceRepository {
	return &workspaceRepo{
		inner: inner,
	}
}

func (o *workspaceRepo) Update(ctx context.Context,
	workspace *malak.Workspace) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	workspace.UpdatedAt = time.Now()

	_, err := o.inner.NewUpdate().
		Where("id = ?", workspace.ID).
		Model(workspace).
		Exec(ctx)
	return err
}

func (o *workspaceRepo) Get(ctx context.Context,
	opts *malak.FindWorkspaceOptions) (*malak.Workspace, error) {
	workspace := new(malak.Workspace)

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	q := o.inner.NewSelect().
		Model(workspace)

	if !util.IsStringEmpty(opts.StripeCustomerID) {
		q = q.Where("stripe_customer_id = ?", opts.StripeCustomerID)
	}

	if !util.IsStringEmpty(opts.Reference.String()) {
		q = q.Where("workspace.reference = ?", opts.Reference)
	}

	if opts.ID != uuid.Nil {
		q = q.Where("workspace.id = ?", opts.ID)
	}

	err := q.Relation("Plan").
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrWorkspaceNotFound
	}

	return workspace, err
}

func (o *workspaceRepo) List(ctx context.Context, user *malak.User) (
	[]malak.Workspace, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	workspaces := make([]malak.Workspace, 0)

	return workspaces, o.inner.NewSelect().
		Model(&workspaces).
		Join(`JOIN roles as role on "role".workspace_id = "workspace".id`).
		Where("role.user_id = ?", user.ID).
		Order("workspace.created_at DESC").
		Relation("Plan").
		Scan(ctx)
}

func (o *workspaceRepo) Create(ctx context.Context,
	opts *malak.CreateWorkspaceOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return o.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(opts.Workspace).Exec(ctx)
		if err != nil {
			return err
		}

		prefs := malak.NewPreference(opts.Workspace)

		prefs.Billing.FinanceEmail = opts.User.Email

		_, err = tx.NewInsert().
			Model(prefs).
			Exec(ctx)
		if err != nil {
			return err
		}

		integrations := make([]*malak.Integration, 0)

		err = tx.NewSelect().
			Model(&integrations).
			Scan(ctx)
		if err != nil {
			return err
		}

		if len(integrations) > 0 {

			var workspaceIntegrations = make([]*malak.WorkspaceIntegration, 0, len(integrations))

			gen := malak.NewReferenceGenerator()

			for _, integration := range integrations {
				workspaceIntegrations = append(workspaceIntegrations, &malak.WorkspaceIntegration{
					WorkspaceID:   opts.Workspace.ID,
					Reference:     gen.Generate(malak.EntityTypeWorkspaceIntegration),
					IntegrationID: integration.ID,
					// enable by default if system level
					IsEnabled: integration.IntegrationType == malak.IntegrationTypeSystem,
					IsActive:  integration.IntegrationType == malak.IntegrationTypeSystem,
				})
			}

			_, err = tx.NewInsert().
				Model(&workspaceIntegrations).
				On("CONFLICT (workspace_id,integration_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Always update the user's current workspace
		opts.User.Metadata.CurrentWorkspace = opts.Workspace.ID
		if _, err := tx.NewUpdate().
			Model(opts.User).
			Where("id = ?", opts.User.ID).
			Exec(ctx); err != nil {
			return err
		}

		var roles malak.UserRoles

		roles = append(roles, &malak.UserRole{
			WorkspaceID: opts.Workspace.ID,
			Role:        malak.RoleAdmin,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      opts.User.ID,
		})

		_, err = tx.NewInsert().Model(&roles).Exec(ctx)
		return err
	})
}

func (w *workspaceRepo) MarkActive(ctx context.Context,
	workspace *malak.Workspace) error {
	return w.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		workspace.IsSubscriptionActive = true
		workspace.UpdatedAt = time.Now()

		_, err := tx.NewUpdate().Model(workspace).
			Where("id = ?", workspace.ID).
			Column("is_subscription_active", "updated_at").
			Exec(ctx)

		return err
	})
}

func (w *workspaceRepo) MarkInActive(ctx context.Context,
	workspace *malak.Workspace) error {
	return w.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		workspace.IsSubscriptionActive = false
		workspace.UpdatedAt = time.Now()

		_, err := tx.NewUpdate().Model(workspace).
			Where("id = ?", workspace.ID).
			Column("is_subscription_active", "updated_at").
			Exec(ctx)
		return err
	})
}
