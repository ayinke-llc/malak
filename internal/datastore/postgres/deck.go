package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type decksRepo struct {
	inner *bun.DB
}

func NewDeckRepository(db *bun.DB) malak.DeckRepository {
	return &decksRepo{
		inner: db,
	}
}

func (d *decksRepo) Create(ctx context.Context,
	deck *malak.Deck, opts *malak.CreateDeckOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewInsert().
				Model(deck).
				Exec(ctx)
			if err != nil {
				return err
			}

			p, err := malak.HashPassword(string(opts.Password.Password))
			if err != nil {
				return err
			}

			deckPreferences := &malak.DeckPreference{
				DeckID: deck.ID,
				Password: malak.PasswordDeckPreferences{
					Enabled:  opts.Password.Enabled,
					Password: malak.Password(p),
				},
				Reference:         opts.Reference,
				ExpiresAt:         opts.ExpiresAt,
				WorkspaceID:       deck.WorkspaceID,
				RequireEmail:      opts.RequireEmail,
				EnableDownloading: opts.EnableDownloading,
				CreatedBy:         deck.CreatedBy,
			}

			_, err = tx.NewInsert().
				Model(deckPreferences).
				Exec(ctx)
			return err
		})
}
