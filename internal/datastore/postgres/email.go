package postgres

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"github.com/ayinke-llc/malak"
)

type emailVerificationRepo struct {
	inner *bun.DB
}

func NewEmailVerificationRepository(db *bun.DB) malak.EmailVerificationRepository {
	return &emailVerificationRepo{
		inner: db,
	}
}

func (e *emailVerificationRepo) Create(ctx context.Context, ev *malak.EmailVerification) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return e.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewDelete().
				Model(&malak.EmailVerification{}).
				Where("user_id = ?", ev.UserID).
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewInsert().Model(ev).Exec(ctx)
			return err
		})
}
