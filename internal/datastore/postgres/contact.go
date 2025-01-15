package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type contactRepo struct {
	inner *bun.DB
}

func NewContactRepository(inner *bun.DB) malak.ContactRepository {
	return &contactRepo{
		inner: inner,
	}
}

// func (o *contactRepo) Update(ctx context.Context,
// 	org *malak.Workspace) error {
//
// 	ctx, cancelFn := withContext(ctx)
// 	defer cancelFn()
//
// 	_, err := o.inner.NewUpdate().
// 		Where("id = ?", org.ID).
// 		Model(org).
// 		Exec(ctx)
// 	return err
// }

func (o *contactRepo) Create(ctx context.Context,
	contact *malak.Contact) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return o.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(contact).
			Exec(ctx)

		if err != nil {

			if strings.Contains(err.Error(), "duplicate key value violates") {
				return malak.ErrContactExists
			}

			return err
		}

		return nil
	})

}

func (o *contactRepo) Get(ctx context.Context,
	opts malak.FetchContactOptions) (*malak.Contact, error) {

	contact := new(malak.Contact)

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	q := o.inner.NewSelect().
		Where("workspace_id = ?", opts.WorkspaceID)

	if !util.IsStringEmpty(opts.Reference.String()) {
		q = q.Where("reference = ?", opts.Reference.String())
	}

	if opts.ID != uuid.Nil {
		q = q.Where("id = ?", opts.ID)
	}

	if !util.IsStringEmpty(opts.Email.String()) {
		q = q.Where("email = ?", opts.Email.String())
	}

	err := q.Model(contact).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrContactNotFound
	}

	return contact, err
}

func (o *contactRepo) List(ctx context.Context,
	opts malak.ListContactOptions) ([]malak.Contact, int64, error) {

	// simple pagination but can always be improved if it proves to be a data hog
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var contacts []malak.Contact
	count := int64(0)

	totalCount, err := o.inner.NewSelect().
		Model(&contacts).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("deleted_at IS NULL").
		Count(ctx)

	if err != nil {
		return nil, 0, err
	}

	count = int64(totalCount)

	return contacts, count, o.inner.NewSelect().
		Model(&contacts).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)
}
