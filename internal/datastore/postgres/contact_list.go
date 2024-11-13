package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayinke-llc/malak"
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

func (c *contactListRepo) List(ctx context.Context,
	opts *malak.ContactListOptions) ([]malak.ContactList, []malak.ContactListMapping, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var lists []malak.ContactList
	err := c.inner.NewSelect().Model(&lists).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, nil, err
	}

	if !opts.IncludeEmails {
		return lists, []malak.ContactListMapping{}, nil
	}

	var listIDs []string
	for _, v := range lists {
		listIDs = append(listIDs, v.ID.String())
	}

	var mappings []malak.ContactListMapping
	err = c.inner.NewSelect().Model(&mappings).
		Where("list_id IN (?)", bun.In(listIDs)).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, nil, err
	}

	fmt.Println(mappings)

	return lists, mappings, nil
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
