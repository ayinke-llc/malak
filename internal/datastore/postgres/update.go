package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
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

func (u *updatesRepo) Update(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := u.inner.NewUpdate().
		Where("id = ?", update.ID).
		Model(update).
		Exec(ctx)
	return err
}

func (u *updatesRepo) Get(ctx context.Context,
	opts malak.FetchUpdateOptions) (*malak.Update, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	update := &malak.Update{}

	sel := u.inner.NewSelect().Model(update)

	if !util.IsStringEmpty(opts.Reference.String()) {
		sel = sel.Where("reference = ?", opts.Reference)
	}

	if opts.ID != uuid.Nil {
		sel = sel.Where("id = ?", opts.ID)
	}

	if opts.Status.IsValid() {
		sel = sel.Where("status = ?", opts.Status)
	}

	err := sel.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrUpdateNotFound
	}

	return update, err
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

func (u *updatesRepo) Delete(ctx context.Context,
	update *malak.Update) error {

	_, err := u.inner.NewDelete().Model(update).
		Where("id = ?", update.ID).
		Exec(ctx)

	return err
}
