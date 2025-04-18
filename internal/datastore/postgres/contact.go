package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

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

func (o *contactRepo) Update(ctx context.Context,
	contact *malak.Contact) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	contact.UpdatedAt = time.Now()

	_, err := o.inner.NewUpdate().
		Where("id = ?", contact.ID).
		Model(contact).
		Exec(ctx)
	return err
}

func (o *contactRepo) Create(ctx context.Context,
	contacts ...*malak.Contact) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return o.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(&contacts).
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

func (o *contactRepo) Overview(ctx context.Context, workspaceID uuid.UUID) (*malak.ContactOverview, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	total, err := o.inner.NewSelect().
		Model((*malak.Contact)(nil)).
		Where("workspace_id = ?", workspaceID).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return nil, err
	}

	return &malak.ContactOverview{
		TotalContacts: int64(total),
	}, nil
}

func (o *contactRepo) Search(ctx context.Context, opts malak.SearchContactOptions) ([]malak.Contact, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	if opts.WorkspaceID == uuid.Nil {
		return []malak.Contact{}, nil
	}

	var contacts []malak.Contact

	q := o.inner.NewSelect().
		Model(&contacts).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("deleted_at IS NULL")

	if !hermes.IsStringEmpty(opts.SearchValue) {
		searchValue := strings.ToLower(opts.SearchValue)
		q = q.Where(`(
			LOWER(email) LIKE ? OR 
			LOWER(first_name) LIKE ? OR 
			LOWER(last_name) LIKE ? OR 
			LOWER(company) LIKE ?
		)`,
			"%"+searchValue+"%",
			"%"+searchValue+"%",
			"%"+searchValue+"%",
			"%"+searchValue+"%")
	}

	err := q.Order("created_at DESC").Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []malak.Contact{}, nil
		}
		return nil, err
	}

	return contacts, nil
}

func (o *contactRepo) Delete(ctx context.Context,
	contact *malak.Contact) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return o.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		// _, err := tx.NewDelete().
		// 	Where("fundraising_pipeline_column_contact_id IN (SELECT id FROM fundraising_pipeline_column_contacts WHERE contact_id = ?)", contact.ID).
		// 	Model(new(malak.FundraiseContactActivity)).
		// 	Exec(ctx)
		// if err != nil {
		// 	return err
		// }

		_, err := tx.NewDelete().
			Where("fundraising_pipeline_column_contact_id IN (SELECT id FROM fundraising_pipeline_column_contacts WHERE contact_id = ?)", contact.ID).
			Model(new(malak.FundraiseContactDealDetails)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Where("fundraising_pipeline_column_contact_id IN (SELECT id FROM fundraising_pipeline_column_contacts WHERE contact_id = ?)", contact.ID).
			Model(new(malak.FundraiseContactPosition)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Where("contact_id = ?", contact.ID).
			Model(new(malak.FundraiseContact)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
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

		// finally delete the contact
		_, err = tx.NewDelete().
			Where("id = ?", contact.ID).
			Model(contact).
			Exec(ctx)
		return err
	})
}

func (o *contactRepo) All(ctx context.Context,
	workspaceID uuid.UUID) ([]malak.Contact, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var contacts []malak.Contact

	err := o.inner.NewSelect().
		Model(&contacts).
		Where("workspace_id = ?", workspaceID).
		Order("created_at DESC").
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []malak.Contact{}, nil
		}

		return nil, err
	}

	return contacts, nil
}
