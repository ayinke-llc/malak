package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
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

	contact := &malak.Contact{
		Lists: make([]malak.ContactListMapping, 0),
	}

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	q := o.inner.NewSelect().
		Where("workspace_id = ?", opts.WorkspaceID)

	if !hermes.IsStringEmpty(opts.Reference.String()) {
		q = q.Where("reference = ?", opts.Reference.String())
	}

	if opts.ID != uuid.Nil {
		q = q.Where("id = ?", opts.ID)
	}

	if !hermes.IsStringEmpty(opts.Email.String()) {
		q = q.Where("email = ?", opts.Email.String())
	}

	err := q.Model(contact).
		Relation("Lists").
		Relation("Lists.List").
		Scan(ctx)

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

func (o *contactRepo) Delete(ctx context.Context,
	contact *malak.Contact) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return o.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewDelete().
			Where("contact_id = ?", contact.ID).
			Model(new(malak.ContactListMapping)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Where("contact_id = ?", contact.ID).
			Model(new(malak.ContactShare)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Where("contact_id = ?", contact.ID).
			Model(new(malak.UpdateRecipient)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Where("id = ?", contact.ID).
			Model(contact).
			Exec(ctx)
		return err
	})

}
