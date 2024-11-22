package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
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

func (u *updatesRepo) TogglePinned(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewUpdate().
				Where("id = ?", update.ID).
				Set("is_pinned = CASE WHEN is_pinned = true THEN false ELSE true END").
				Model(update).
				Exec(ctx)

			if err != nil {
				return err
			}

			count, err := tx.NewSelect().
				Model(new(malak.Update)).
				Where("is_pinned = ?", true).
				Where("workspace_id = ?", update.WorkspaceID).
				Count(ctx)
			if err != nil {
				return err
			}

			if count > malak.MaximumNumberOfPinnedUpdates {
				return malak.ErrPinnedUpdateCapacityExceeded
			}

			return nil
		})
}

func (u *updatesRepo) Update(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := u.inner.NewUpdate().
		Where("id = ?", update.ID).
		Model(update).
		Exec(ctx)
	return err
}

func (u *updatesRepo) Stat(ctx context.Context, update *malak.Update) (
	*malak.UpdateStat, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	stat := new(malak.UpdateStat)

	err := u.inner.NewSelect().
		Model(stat).
		Where("update_id = ?", update.ID).
		Scan(ctx)

	return stat, err
}

func (u *updatesRepo) GetByID(ctx context.Context, id uuid.UUID) (*malak.Update, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	update := &malak.Update{}

	sel := u.inner.NewSelect().Model(update).
		Where("id = ?", id.String())

	err := sel.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrUpdateNotFound
	}

	return update, err
}

func (u *updatesRepo) Get(ctx context.Context,
	opts malak.FetchUpdateOptions) (*malak.Update, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	update := &malak.Update{}

	sel := u.inner.NewSelect().Model(update).
		Where("reference = ?", opts.Reference)

	if opts.ID != uuid.Nil {
		sel = sel.Where("id = ?", opts.ID)
	}

	if opts.Status.IsValid() {
		sel = sel.Where("status = ?", opts.Status)
	}

	err := sel.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrUpdateNotFound
	}

	return update, err
}

func (u *updatesRepo) Create(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(update).
			Exec(ctx)
		if err != nil {
			return err
		}

		updateStats := &malak.UpdateStat{
			UpdateID:  update.ID,
			Reference: malak.NewReferenceGenerator().Generate(malak.EntityTypeUpdateStat),
		}

		_, err = tx.NewInsert().
			Model(updateStats).
			Exec(ctx)
		return err
	})
}

func (u *updatesRepo) List(ctx context.Context,
	opts malak.ListUpdateOptions) ([]malak.Update, error) {

	updates := make([]malak.Update, 0, opts.Paginator.PerPage)

	q := u.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", opts.WorkspaceID)

	if opts.Status != malak.ListUpdateFilterStatusAll {
		q = q.Where("status = ?", opts.Status)
	}

	err := q.Model(&updates).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)

	return updates, err
}

func (u *updatesRepo) ListPinned(ctx context.Context,
	workspaceID uuid.UUID) ([]malak.Update, error) {

	updates := make([]malak.Update, 0, malak.MaximumNumberOfPinnedUpdates)

	return updates, u.inner.NewSelect().
		Model(&updates).
		Order("created_at DESC").
		Where("workspace_id = ?", workspaceID).
		Where("is_pinned = ?", true).
		Limit(malak.MaximumNumberOfPinnedUpdates).
		Scan(ctx)
}

func (u *updatesRepo) Delete(ctx context.Context,
	update *malak.Update) error {

	_, err := u.inner.NewDelete().Model(update).
		Where("id = ?", update.ID).
		Exec(ctx)

	return err
}

func (u *updatesRepo) GetSchedule(ctx context.Context, scheduleID uuid.UUID) (
	*malak.UpdateSchedule, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	schedule := &malak.UpdateSchedule{}

	err := u.inner.NewSelect().
		Where("id = ?", scheduleID).
		Model(schedule).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrUpdateNotFound
	}

	return schedule, err
}

func (u *updatesRepo) SendUpdate(ctx context.Context,
	opts *malak.CreateUpdateOptions) error {

	return u.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			contacts := make([]malak.Contact, 0, len(opts.Emails))
			var insertedContactIDs = make([]uuid.UUID, 0, len(opts.Emails))

			for _, email := range opts.Emails {
				contacts = append(contacts, malak.Contact{
					WorkspaceID: opts.WorkspaceID,
					Reference:   malak.Reference(opts.Reference(malak.EntityTypeContact)),
					Email:       email,
					FirstName:   "Investor",
					Metadata:    make(malak.CustomContactMetadata),
					OwnerID:     opts.UserID,
					CreatedBy:   opts.UserID,
				})
			}

			_, err := u.inner.NewInsert().
				Model(&contacts).
				// if we already have this email in this workspace
				On("CONFLICT (email,workspace_id) DO NOTHING").
				Returning("id").
				Exec(ctx, &insertedContactIDs)
			if err != nil {
				return err
			}

			// Retrieve IDs of all contacts (both newly inserted and existing ones)
			// Refetching since on CONFLICT skips the existing ids and do not return them
			err = u.inner.NewSelect().
				Model(&contacts).
				Column("id").
				Where("email IN (?)", bun.In(opts.Emails)).
				Where("workspace_id = ?", opts.WorkspaceID).
				Scan(ctx, &insertedContactIDs)
			if err != nil {
				return err
			}

			_, err = tx.NewInsert().Model(opts.Schedule).
				Exec(ctx)
			if err != nil {
				return err
			}

			var recipients = make([]malak.UpdateRecipient, 0, len(opts.Emails))

			for _, contact := range insertedContactIDs {
				recipients = append(recipients, malak.UpdateRecipient{
					ContactID:  contact,
					UpdateID:   opts.Schedule.UpdateID,
					ScheduleID: opts.Schedule.ID,
					Reference:  opts.Generator.Generate(malak.EntityTypeRecipient),
				})
			}

			_, err = tx.NewInsert().Model(&recipients).
				Exec(ctx)
			return err
		})
}
