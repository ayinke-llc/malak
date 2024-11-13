package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type contactListRepo struct {
	inner *bun.DB
}

func NewContactListRepository(db *bun.DB) malak.ContactListRepository {
	return &contactListRepo{
		inner: db,
	}
}

func (c *contactListRepo) Get(ctx context.Context,
	opts malak.FetchContactListOptions) (*malak.ContactList, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	list := &malak.ContactList{}

	err := c.inner.NewSelect().Model(list).
		Where("reference = ?", opts.Reference).
		Where("workspace_id = ?", opts.WorkspaceID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrContactListNotFound
	}

	return list, err
}

func (c *contactListRepo) Delete(ctx context.Context,
	list *malak.ContactList) error {

	_, err := c.inner.NewDelete().Model(list).
		Where("id = ?", list.ID).
		Exec(ctx)
	return err
}

func (c *contactListRepo) Update(ctx context.Context,
	list *malak.ContactList) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := c.inner.NewUpdate().
		Where("id = ?", list.ID).
		Model(list).
		Exec(ctx)
	return err
}

func (c *contactListRepo) Create(ctx context.Context,
	list *malak.ContactList) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return c.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewInsert().
				Model(list).
				Exec(ctx)

			return err
		})

}

func (c *contactListRepo) List(ctx context.Context, id uuid.UUID) ([]malak.ContactList, error) {
	list := make([]malak.ContactList, 0)

	err := c.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", id).
		Model(&list).
		Scan(ctx)

	return list, err
}

func (c *contactListRepo) Add(ctx context.Context,
	mapping *malak.ContactListMapping) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return c.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewInsert().
				Model(mapping).
				Exec(ctx)

			return err
		})
}
