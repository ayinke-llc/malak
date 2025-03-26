package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
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

	// TODO:(adelowo): pagination? will people have like 20/30 decks?
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

	err := d.inner.NewSelect().
		Model(deck).
		Where("deck.reference = ?", opts.Reference).
		Where("deck.workspace_id = ?", opts.WorkspaceID).
		Relation("DeckPreference").
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrDeckNotFound
	}

	return deck, err
}

func (d *decksRepo) Delete(ctx context.Context, deck *malak.Deck) error {

	ctx, cancel := withContext(ctx)
	defer cancel()

	return d.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewDelete().Model(&malak.DeckPreference{}).
			Where("deck_id = ?", deck.ID).
			Exec(ctx)

		if err != nil {
			return err
		}
		_, err = d.inner.NewDelete().Model(deck).
			Where("reference = ?", deck.Reference).
			Exec(ctx)
		return err
	})
}

func (d *decksRepo) UpdatePreferences(ctx context.Context, deck *malak.Deck) error {

	ctx, cancel := withContext(ctx)
	defer cancel()

	return d.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		deck.DeckPreference.UpdatedAt = time.Now()

		p, err := malak.HashPassword(string(deck.DeckPreference.Password.Password))
		if err != nil {
			return err
		}

		deck.DeckPreference.Password.Password = malak.Password(p)

		_, err = d.inner.NewUpdate().
			Model(deck.DeckPreference).
			Where("id = ?", deck.DeckPreference.ID).
			Exec(ctx)
		return err
	})
}

func (d *decksRepo) ToggleArchive(ctx context.Context,
	deck *malak.Deck) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewUpdate().
				Where("id = ?", deck.ID).
				Set("is_archived = CASE WHEN is_archived = true THEN false ELSE true END").
				Model(deck).
				Exec(ctx)

			return err
		})
}

func (d *decksRepo) TogglePinned(ctx context.Context,
	deck *malak.Deck) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewUpdate().
				Where("id = ?", deck.ID).
				Set("is_pinned = CASE WHEN is_pinned = true THEN false ELSE true END").
				Model(deck).
				Exec(ctx)

			if err != nil {
				return err
			}

			count, err := tx.NewSelect().
				Model(new(malak.Deck)).
				Where("is_pinned = ?", true).
				Where("workspace_id = ?", deck.WorkspaceID).
				Count(ctx)
			if err != nil {
				return err
			}

			if count > malak.MaximumNumberOfPinnedUpdates {
				return malak.ErrPinnedDeckCapacityExceeded
			}

			return nil
		})
}

// same as Get without the workspace_id
// Separate api so as to not potentially misuse
func (d *decksRepo) PublicDetails(ctx context.Context,
	ref malak.Reference) (*malak.Deck, error) {

	ctx, cancel := withContext(ctx)
	defer cancel()

	deck := &malak.Deck{}

	err := d.inner.NewSelect().
		Model(deck).
		Where("deck.short_link = ?", ref).
		Relation("DeckPreference").
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrDeckNotFound
	}

	return deck, err
}

func (d *decksRepo) CreateDeckSession(ctx context.Context,
	session *malak.DeckViewerSession) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().
				Model(session).
				Exec(ctx)
			return err
		})
}

func (d *decksRepo) UpdateDeckSession(ctx context.Context,
	opts *malak.UpdateDeckSessionOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if opts.CreateContact {
				_, err := tx.NewInsert().
					Model(opts.Contact).
					Exec(ctx)
				if err != nil {
					return err
				}
			}

			if opts.Contact != nil {
				opts.Session.ContactID = opts.Contact.ID
			}

			_, err := tx.NewUpdate().
				Model(opts.Session).
				Where("id = ?", opts.Session.ID).
				Exec(ctx)
			return err
		})
}

func (d *decksRepo) FindDeckSession(ctx context.Context,
	sessionID string) (*malak.DeckViewerSession, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	session := new(malak.DeckViewerSession)
	err := d.inner.NewSelect().
		Model(session).
		Where("session_id = ?", sessionID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, malak.ErrDeckNotFound
	}

	return session, err
}

func (d *decksRepo) SessionAnalytics(ctx context.Context,
	opts *malak.ListSessionAnalyticsOptions) ([]*malak.DeckViewerSession, int64, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	sessions := make([]*malak.DeckViewerSession, 0)

	query := d.inner.NewSelect().
		Model(&sessions).
		Where("deck_id = ?", opts.DeckID).
		Where("deck_viewer_session.created_at >= NOW() - INTERVAL '? days'", opts.Days).
		Order("deck_viewer_session.created_at DESC")

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Relation("Contact").
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)

	return sessions, int64(total), err
}

func (d *decksRepo) DeckEngagements(ctx context.Context,
	opts *malak.ListDeckEngagementsOptions) (*malak.DeckEngagementResponse, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var dailyEngagements []malak.DeckDailyEngagement
	var geographicStats []malak.DeckGeographicStat

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return d.inner.NewSelect().
			Model(&dailyEngagements).
			Where("deck_id = ?", opts.DeckID).
			Order("engagement_date DESC").
			Scan(ctx)
	})

	g.Go(func() error {
		return d.inner.NewSelect().
			Model(&geographicStats).
			Where("deck_id = ?", opts.DeckID).
			Order("view_count DESC").
			Scan(ctx)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &malak.DeckEngagementResponse{
		DailyEngagements: dailyEngagements,
		GeographicStats:  geographicStats,
	}, nil
}

func (d *decksRepo) Overview(ctx context.Context, workspaceID uuid.UUID) (*malak.DeckOverview, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	// Use errgroup to run both counts concurrently
	g, ctx := errgroup.WithContext(ctx)

	var totalDecks, totalSessions int

	g.Go(func() error {
		var err error
		totalDecks, err = d.inner.NewSelect().
			Model((*malak.Deck)(nil)).
			Where("workspace_id = ?", workspaceID).
			Count(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		totalSessions, err = d.inner.NewSelect().
			Model((*malak.DeckViewerSession)(nil)).
			Join("JOIN decks ON decks.id = deck_viewer_session.deck_id").
			Where("decks.workspace_id = ?", workspaceID).
			Count(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &malak.DeckOverview{
		TotalDecks:          int64(totalDecks),
		TotalViewerSessions: int64(totalSessions),
	}, nil
}
