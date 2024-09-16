package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type updatesRepo struct {
	inner *bun.DB
}

func NewUpdatesRepository(db *bun.DB) malak.UpdateRepository {
	return &updatesRepo{
		inner: db,
	}
}

func (u *updatesRepo) Create(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(update).
			Exec(ctx)
		return err
	})
}

func (u *updatesRepo) List(ctx context.Context,
	opts malak.ListUpdateOptions) ([]malak.Update, error) {

	updates := make([]malak.Update, 0, opts.Paginator.PerPage)

	q := u.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", opts.WorkspaceID)

	if opts.Status != malak.ListUpdateFilterStatusAll {
		q = q.Where("status = ?", opts.Status)
	}

	err := q.Model(&updates).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)

	return updates, err
}
