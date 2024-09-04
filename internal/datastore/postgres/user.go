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

type userRepo struct {
	inner *bun.DB
}

func NewUserRepository(db *bun.DB) malak.UserRepository {
	return &userRepo{
		inner: db,
	}
}

func (u *userRepo) Update(ctx context.Context, user *malak.User) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := u.inner.NewUpdate().
		Where("id = ?", user.ID).
		Model(user).
		Exec(ctx)
	return err
}

func (u *userRepo) Get(ctx context.Context, opts *malak.FindUserOptions) (*malak.User, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	user := &malak.User{
		Roles: malak.UserRoles{},
	}

	sel := u.inner.NewSelect().Model(user).Relation("Roles")

	if !util.IsStringEmpty(opts.Email.String()) {
		sel = sel.Where("email = ?", opts.Email.String())
	}

	if opts.ID != uuid.Nil {
		sel = sel.Where("id = ?", opts.ID)
	}

	err := sel.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrUserNotFound
	}

	return user, err
}

func (u *userRepo) Create(ctx context.Context, user *malak.User) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().Model(user).
				Exec(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					return malak.ErrUserExists
				}

				return err
			}

			return nil
		})
}
