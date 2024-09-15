package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ayinke-llc/malak"
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

func (u *updatesRepo) Create(ctx context.Context,
	contact *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

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

// func (o *contactRepo) Get(ctx context.Context,
// 	opts malak.FetchContactOptions) (*malak.Contact, error) {
//
// 	contact := new(malak.Contact)
//
// 	ctx, cancelFn := withContext(ctx)
// 	defer cancelFn()
//
// 	q := o.inner.NewSelect()
//
// 	if !util.IsStringEmpty(opts.Reference.String()) {
// 		q = q.Where("reference = ?", opts.Reference.String())
// 	}
//
// 	if opts.ID != uuid.Nil {
// 		q = q.Where("id = ?", opts.ID)
// 	}
//
// 	if !util.IsStringEmpty(opts.Email.String()) {
// 		q = q.Where("email = ?", opts.Email.String())
// 	}
//
// 	err := q.Model(contact).Scan(ctx)
//
// 	if errors.Is(err, sql.ErrNoRows) {
// 		err = malak.ErrContactNotFound
// 	}
//
// 	return contact, err
// }
