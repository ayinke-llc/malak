package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func (u *updatesRepo) Create(ctx context.Context,
	update *malak.Update, opts *malak.TemplateCreateUpdateOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().
			Model(update).
			Exec(ctx)
		if err != nil {
			return err
		}

		if opts.IsSystemTemplate {
			res, err := tx.NewUpdate().Model(new(malak.SystemTemplate)).
				Where("reference = ?", opts.Reference).
				Set("number_of_uses = number_of_uses + 1").
				Exec(ctx)
			if err != nil {
				return err
			}

			affected, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if affected == 0 {
				return sql.ErrNoRows
			}
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

	update.UpdatedAt = time.Now()

	_, err := u.inner.NewUpdate().
		Where("id = ?", update.ID).
		Model(update).
		Exec(ctx)
	return err
}

func (u *updatesRepo) UpdateStat(ctx context.Context,
	stat *malak.UpdateStat,
	recipientStat *malak.UpdateRecipientStat) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		if stat != nil {
			_, err := tx.NewUpdate().
				Where("id = ?", stat.ID).
				Model(stat).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		if recipientStat == nil {
			return nil
		}

		_, err := tx.NewUpdate().
			Where("id = ?", recipientStat.ID).
			Model(recipientStat).
			Exec(ctx)
		return err
	})
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

	if errors.Is(err, sql.ErrNoRows) {
		return stat, nil
	}

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
		Where("reference = ?", opts.Reference).
		Where("workspace_id = ?", opts.WorkspaceID)

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

func (u *updatesRepo) List(ctx context.Context,
	opts malak.ListUpdateOptions) ([]malak.Update, int64, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	updates := make([]malak.Update, 0, opts.Paginator.PerPage)

	q := u.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", opts.WorkspaceID)

	if opts.Status != malak.ListUpdateFilterStatusAll {
		q = q.Where("status = ?", opts.Status)
	}

	// Get total count with same filters
	total, err := q.
		Model(&updates).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = q.Model(&updates).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)

	return updates, int64(total), err
}

func (u *updatesRepo) ListPinned(ctx context.Context,
	workspaceID uuid.UUID) ([]malak.Update, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	updates := make([]malak.Update, 0, malak.MaximumNumberOfPinnedUpdates)

	q := u.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", workspaceID).
		Where("is_pinned = ?", true)

	err := q.Model(&updates).
		Limit(malak.MaximumNumberOfPinnedUpdates).
		Scan(ctx)

	return updates, err
}

func (u *updatesRepo) Delete(ctx context.Context,
	update *malak.Update) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

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

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return u.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewUpdate().
				Model(new(malak.Update)).
				Where("reference = ?", opts.UpdateReference).
				Set("status = ?", malak.UpdateStatusSent).
				Exec(ctx)
			if err != nil {
				return err
			}

			contacts := make([]malak.Contact, 0, len(opts.Emails))
			var insertedContactIDs = make([]uuid.UUID, 0, len(opts.Emails))

			for _, email := range opts.Emails {
				contacts = append(contacts, malak.Contact{
					WorkspaceID: opts.WorkspaceID,
					Reference:   malak.Reference(opts.Reference(malak.EntityTypeContact)),
					Email:       email,
					FirstName:   email.String(),
					Metadata:    make(malak.CustomContactMetadata),
					OwnerID:     opts.UserID,
					CreatedBy:   opts.UserID,
				})
			}

			_, err = tx.NewInsert().
				Model(&contacts).
				On("CONFLICT (email,workspace_id) DO NOTHING").
				Returning("id").
				Exec(ctx, &insertedContactIDs)
			if err != nil {
				return err
			}

			// Retrieve IDs of all contacts (both newly inserted and existing ones)
			// Refetching since on CONFLICT skips the existing ids and do not return them
			err = tx.NewSelect().
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

			var sharedItems = make([]malak.ContactShare, 0, len(opts.Emails))

			for _, contact := range insertedContactIDs {
				recipients = append(recipients, malak.UpdateRecipient{
					ContactID:  contact,
					UpdateID:   opts.Schedule.UpdateID,
					ScheduleID: opts.Schedule.ID,
					Reference:  opts.Generator.Generate(malak.EntityTypeRecipient),
					Status:     malak.RecipientStatusPending,
				})

				sharedItems = append(sharedItems, malak.ContactShare{
					Reference:     opts.Generator.Generate(malak.EntityTypeContactShare),
					SharedBy:      opts.UserID,
					ContactID:     contact,
					ItemType:      malak.ContactShareItemTypeUpdate,
					ItemID:        opts.Schedule.UpdateID,
					ItemReference: opts.UpdateReference,
				})
			}

			_, err = tx.NewInsert().Model(&recipients).
				On("CONFLICT (contact_id,update_id) DO NOTHING").
				Returning("id").
				Exec(ctx)
			if err != nil {
				return err
			}

			count, err := tx.NewSelect().Model(new(malak.UpdateRecipient)).
				Where("update_id = ?", opts.Schedule.UpdateID).
				Count(ctx)
			if err != nil {
				return err
			}

			if err := opts.Plan.Metadata.Updates.MaxRecipients.TakeN(int64(count)); err != nil {
				return err
			}

			_, err = tx.NewInsert().Model(&sharedItems).
				Returning("id").
				On("CONFLICT (item_reference,contact_id) DO NOTHING").
				Exec(ctx)
			return err
		})
}

func (u *updatesRepo) GetStatByEmailID(ctx context.Context,
	emailID string,
	provider malak.UpdateRecipientLogProvider) (
	*malak.UpdateRecipientLog, *malak.UpdateRecipientStat, error) {

	// can just JOIN this into one query
	// but meh for now currently in the trenches
	//
	//
	//
	//

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var log *malak.UpdateRecipientLog
	var stat *malak.UpdateRecipientStat

	log = &malak.UpdateRecipientLog{
		Recipient: &malak.UpdateRecipient{},
	}

	err := u.inner.NewSelect().
		Model(log).
		Where("provider_id = ?", emailID).
		Where("provider = ?", provider.String()).
		Relation("Recipient").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	stat = &malak.UpdateRecipientStat{}

	err = u.inner.NewSelect().
		Model(stat).
		Where("recipient_id = ?", log.RecipientID).
		Scan(ctx)

	return log, stat, err
}

func (u *updatesRepo) RecipientStat(ctx context.Context,
	update *malak.Update) ([]malak.UpdateRecipient, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	recipients := make([]malak.UpdateRecipient, 0)

	return recipients, u.inner.NewSelect().
		Model(&recipients).
		Order("created_at DESC").
		Where("update_id = ?", update.ID).
		Relation("UpdateRecipientStat").
		Relation("Contact").
		Scan(ctx)
}

func (u *updatesRepo) Overview(ctx context.Context, workspaceID uuid.UUID) (*malak.UpdateOverview, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	// Get total count of updates for the workspace
	total, err := u.inner.NewSelect().
		Model((*malak.Update)(nil)).
		Where("workspace_id = ?", workspaceID).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get last 10 sent updates
	lastUpdates := make([]malak.Update, 0, 10)
	err = u.inner.NewSelect().
		Model(&lastUpdates).
		Where("workspace_id = ?", workspaceID).
		Where("status = ?", malak.UpdateStatusSent).
		Order("created_at DESC").
		Limit(10).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &malak.UpdateOverview{
		Total:       int64(total),
		LastUpdates: lastUpdates,
	}, nil
}
