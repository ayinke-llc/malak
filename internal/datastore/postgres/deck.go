package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
func (d *decksRepo) List(ctx context.Context,
	workspace *malak.Workspace) ([]malak.Deck, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	// TODO: pagination? will people have like 20/30 decks?
	// so wait till we get there
	decks := make([]malak.Deck, 0)

	q := d.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", workspace.ID)

	err := q.Model(&decks).
		Scan(ctx)

	return decks, err
}

func (d *decksRepo) Create(ctx context.Context,
	deck *malak.Deck, opts *malak.CreateDeckOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			deck.CreatedAt = time.Now()
			deck.UpdatedAt = time.Now()

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

func (d *decksRepo) Get(ctx context.Context, opts malak.FetchDeckOptions) (
	*malak.Deck, error) {

	ctx, cancel := withContext(ctx)
	defer cancel()

	deck := &malak.Deck{}

	err := d.inner.NewSelect().Model(deck).
		Where("reference = ?", opts.Reference).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrDeckNotFound
	}

	return deck, err
}

func (d *decksRepo) Delete(ctx context.Context, deck *malak.Deck) error {

	ctx, cancel := withContext(ctx)
	defer cancel()

	_, err := d.inner.NewDelete().Model(deck).
		Where("reference = ?", deck.Reference).
		Exec(ctx)
	return err
}