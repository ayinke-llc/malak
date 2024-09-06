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
	org *malak.Workspace) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := o.inner.NewUpdate().
		Where("id = ?", org.ID).
		Model(org).
		Exec(ctx)
	return err
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

		if len(opts.User.Roles) == 0 {
			opts.User.Metadata.CurrentWorkspace = opts.Workspace.ID
			if _, err := tx.NewUpdate().
				Model(opts.User).
				Where("id = ?", opts.User.ID).
				Exec(ctx); err != nil {
				return err
			}
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

func (o *workspaceRepo) Get(ctx context.Context,
	opts *malak.FindWorkspaceOptions) (*malak.Workspace, error) {
	workspace := new(malak.Workspace)

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	q := o.inner.NewSelect()

	if !util.IsStringEmpty(opts.StripeCustomerID) {
		q = q.Where("stripe_customer_id = ?", opts.StripeCustomerID)
	}

	if opts.ID != uuid.Nil {
		q = q.Where("id = ?", opts.ID)
	}

	err := q.Model(workspace).Scan(ctx)

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
		Scan(ctx)
}
