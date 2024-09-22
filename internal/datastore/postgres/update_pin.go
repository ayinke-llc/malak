package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type pinnedUpdateRepo struct {
	inner *bun.DB
}

func NewPinUpdateRepository(db *bun.DB) malak.PinnedUpdateRepository {
	return &pinnedUpdateRepo{
		inner: db,
	}
}

func (p *pinnedUpdateRepo) Get(ctx context.Context,
	opts malak.FetchPinnedUpdate) (*malak.PinnedUpdate, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	pinnedUpdate := &malak.PinnedUpdate{}

	sel := p.inner.NewSelect().
		Where("id = ?", opts.UpdateID)

	return pinnedUpdate, sel.Model(pinnedUpdate).Scan(ctx)
}

func (o *pinnedUpdateRepo) List(ctx context.Context,
	opts malak.ListPinnedUpdates) ([]malak.PinnedUpdate, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	pinned := make([]malak.PinnedUpdate, 3)

	return pinned, o.inner.NewSelect().
		Model(&pinned).
		Scan(ctx)
}

func (o *pinnedUpdateRepo) Pin(ctx context.Context,
	update *malak.Update, state malak.PinState,
	user *malak.User) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	err := o.inner.RunInTx(ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			switch state {

			case malak.PinStateUnpin:

				_, err := tx.NewDelete().Model(new(malak.PinnedUpdate)).
					Where("workspace_id = ?", update.WorkspaceID).
					Where("update_id = ?", update.ID).
					Where("user_id = ?", user.ID).
					Exec(ctx)

				return err

			case malak.PinStatePin:

				pinnedUpdate := &malak.PinnedUpdate{
					UpdateID:    update.ID,
					WorkspaceID: update.WorkspaceID,
					PinnedBy:    user.ID,
				}

				_, err := tx.NewInsert().
					Model(pinnedUpdate).
					Exec(ctx)
				if err != nil {
					return err
				}

				count, err := tx.NewSelect().Model(new(malak.PinnedUpdate)).
					Where("workspace_id = ?", update.WorkspaceID).
					Count(ctx)
				if err != nil {
					return err
				}

				if count > malak.MaximumNumberOfPinnedUpdates {
					return malak.ErrPinnedUpdateCapacityExceeded
				}

				return nil
			}

			return nil
		})

	return err
}
