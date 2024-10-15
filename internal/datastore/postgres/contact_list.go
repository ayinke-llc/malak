package postgres

import (
	"context"
	"database/sql"

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

func (c *contactListRepo) Add(ctx context.Context, contacts ...*malak.Contact) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return c.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(contacts).
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
