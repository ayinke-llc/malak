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
	opts *malak.ContactListOptions) ([]malak.ContactList, []malak.ContactListMappingWithContact, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	query := `
        WITH list_data AS (
            SELECT 
                id, workspace_id, title, reference, 
                created_by, created_at, updated_at, deleted_at
            FROM contact_lists 
            WHERE workspace_id = ? AND deleted_at IS NULL ORDER BY created_at DESC
        )
        SELECT 
	          COALESCE(json_agg(list_data)::text, '[]') as lists,
            COALESCE(
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
            )::text as mappings
        FROM list_data
        LEFT JOIN contact_list_mappings clm ON clm.list_id = list_data.id 
            AND clm.deleted_at IS NULL
        LEFT JOIN contacts c ON c.id = clm.contact_id 
            AND c.deleted_at IS NULL;
    `

	var listsStr, mappingsStr string
	err := c.inner.QueryRowContext(ctx, query, opts.WorkspaceID).Scan(&listsStr, &mappingsStr)
	if err != nil {
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
