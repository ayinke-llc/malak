package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
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

	return c.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewDelete().Model(new(malak.ContactListMapping)).
			Where("list_id = ?", list.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().Model(list).
			Where("id = ?", list.ID).
			Exec(ctx)
		return err

	})

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
	opts *malak.ContactListOptions) ([]malak.ContactList, []malak.ContactListMappingWithContact, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	query := `
		SELECT 
			(
				SELECT COALESCE(json_agg(cl)::text, '[]')
				FROM (
					SELECT id, workspace_id, title, reference, created_by, created_at, updated_at, deleted_at
					FROM contact_lists 
					WHERE workspace_id = ? AND deleted_at IS NULL 
					ORDER BY created_at DESC
				) cl
			) as lists,
			(
				SELECT COALESCE(
					json_agg(
						DISTINCT jsonb_build_object(
							'id', clm.id,
							'list_id', clm.list_id,
							'contact_id', clm.contact_id,
							'reference', clm.reference,
							'email', c.email
						)
					) FILTER (WHERE clm.id IS NOT NULL),
					'[]'
				)::text
				FROM contact_lists cl
				LEFT JOIN contact_list_mappings clm ON clm.list_id = cl.id AND clm.deleted_at IS NULL
				LEFT JOIN contacts c ON c.id = clm.contact_id AND c.deleted_at IS NULL
				WHERE cl.workspace_id = ? AND cl.deleted_at IS NULL
			) as mappings`

	var listsStr, mappingsStr string
	if err := c.inner.QueryRowContext(ctx, query, opts.WorkspaceID, opts.WorkspaceID).Scan(&listsStr, &mappingsStr); err != nil {
		return nil, nil, err
	}

	var lists []malak.ContactList
	var mappings []malak.ContactListMappingWithContact

	if err := json.Unmarshal([]byte(listsStr), &lists); err != nil {
		return nil, nil, fmt.Errorf("unmarshal lists: %w", err)
	}

	if err := json.Unmarshal([]byte(mappingsStr), &mappings); err != nil {
		return nil, nil, fmt.Errorf("unmarshal mappings: %w", err)
	}

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

			if malak.IsDuplicateUniqueError(err) {
				return nil
			}

			return err
		})
}
