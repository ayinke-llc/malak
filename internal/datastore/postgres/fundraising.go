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

type fundingRepo struct {
	inner *bun.DB
}

func NewFundingRepo(db *bun.DB) malak.FundraisingPipelineRepository {
	return &fundingRepo{
		inner: db,
	}
}

func (d *fundingRepo) Create(ctx context.Context, pipeline *malak.FundraisingPipeline, additionalColumns ...malak.FundraisingPipelineColumn) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().
				Model(pipeline).
				Exec(ctx)
			if err != nil {
				return err
			}

			for i := range additionalColumns {
				additionalColumns[i].FundraisingPipelineID = pipeline.ID
			}

			if len(additionalColumns) > 0 {
				_, err = tx.NewInsert().
					Model(&additionalColumns).
					Exec(ctx)
			}

			return err
		})
}

func (d *fundingRepo) List(ctx context.Context,
	opts malak.ListPipelineOptions) ([]malak.FundraisingPipeline, int64, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var pipelines []malak.FundraisingPipeline

	q := d.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", opts.WorkspaceID)

	if opts.ActiveOnly {
		q = q.Where("is_closed = ?", false)
	}

	total, err := q.
		Model(&pipelines).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = q.Model(&pipelines).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return pipelines, int64(total), nil
}

func (d *fundingRepo) Board(ctx context.Context, pipeline *malak.FundraisingPipeline) ([]malak.FundraisingPipelineColumn, []malak.FundraiseContact, []malak.FundraiseContactPosition, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var columns []malak.FundraisingPipelineColumn
	var contacts []malak.FundraiseContact
	var positions []malak.FundraiseContactPosition

	err := d.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().
			Model(&columns).
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Order("created_at ASC").
			Scan(ctx)
		if err != nil {
			return err
		}

		err = tx.NewSelect().
			Model(&contacts).
			Relation("DealDetails").
			Relation("Contact").
			Where("fundraising_pipeline_id = ?", pipeline.ID).
			Scan(ctx)
		if err != nil {
			return err
		}

		if len(contacts) > 0 {
			var contactIDs []uuid.UUID
			for _, contact := range contacts {
				contactIDs = append(contactIDs, contact.ID)
			}

			err = tx.NewSelect().
				Model(&positions).
				Where("fundraising_pipeline_column_contact_id IN (?)", bun.In(contactIDs)).
				Order("order_index ASC").
				Scan(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return columns, contacts, positions, nil
}

func (d *fundingRepo) Get(ctx context.Context, opts malak.FetchPipelineOptions) (*malak.FundraisingPipeline, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	pipeline := new(malak.FundraisingPipeline)
	err := d.inner.NewSelect().
		Model(pipeline).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("reference = ?", opts.Reference).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = malak.ErrPipelineNotFound
		}

		return nil, err
	}

	return pipeline, nil
}

func (d *fundingRepo) CloseBoard(ctx context.Context, pipeline *malak.FundraisingPipeline) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := d.inner.NewUpdate().
		Model(pipeline).
		Set("is_closed = ?", true).
		Where("id = ?", pipeline.ID).
		Exec(ctx)

	return err
}

func (d *fundingRepo) AddContactToBoard(ctx context.Context, opts *malak.AddContactToBoardOptions) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		fundraiseContact := &malak.FundraiseContact{
			Reference:                   opts.ReferenceGenerator.Generate(malak.EntityTypeFundraisingPipelineColumnContact),
			ContactID:                   opts.Contact.ID,
			FundraisingPipelineID:       opts.Column.FundraisingPipelineID,
			FundraisingPipelineColumnID: opts.Column.ID,
		}

		_, err := tx.NewInsert().
			Model(fundraiseContact).
			Exec(ctx)
		if err != nil {
			return err
		}

		position := &malak.FundraiseContactPosition{
			ID:                                 uuid.New(),
			Reference:                          opts.ReferenceGenerator.Generate(malak.EntityTypeFundraisingPipelineColumnContactPosition),
			FundraisingPipelineColumnContactID: fundraiseContact.ID,
			// use timestamp as order index for now, it just means the contact will be at the bottom of the column
			// hack but saves us another query
			OrderIndex: time.Now().Unix(),
		}

		_, err = tx.NewInsert().
			Model(position).
			Exec(ctx)
		if err != nil {
			return err
		}

		dealDetails := &malak.FundraiseContactDealDetails{
			Reference:                          opts.ReferenceGenerator.Generate(malak.EntityTypeFundraisingPipelineColumnContactDeal),
			FundraisingPipelineColumnContactID: fundraiseContact.ID,
			Rating:                             int64(opts.Rating),
			CanLeadRound:                       opts.CanLeadRound,
			InitialContact:                     opts.InitialContact,
			CheckSize:                          opts.CheckSize,
		}

		_, err = tx.NewInsert().
			Model(dealDetails).
			Exec(ctx)
		return err
	})
}

func (d *fundingRepo) DefaultColumn(ctx context.Context, pipeline *malak.FundraisingPipeline) (malak.FundraisingPipelineColumn, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var column malak.FundraisingPipelineColumn
	err := d.inner.NewSelect().
		Model(&column).
		Where("fundraising_pipeline_id = ?", pipeline.ID).
		Where("column_type = ?", malak.FundraisePipelineColumnTypeNormal).
		Where("title = ?", "Backlog").
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)

	return column, err
}

func (d *fundingRepo) UpdateContactDeal(ctx context.Context,
	pipeline *malak.FundraisingPipeline, opts malak.UpdateContactDealOptions) error {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := d.inner.NewUpdate().
		Model((*malak.FundraiseContactDealDetails)(nil)).
		Set("rating = ?", opts.Rating).
		Set("can_lead_round = ?", opts.CanLeadRound).
		Set("check_size = ?", opts.CheckSize).
		Set("updated_at = ?", time.Now()).
		Where("fundraising_pipeline_column_contact_id = ?", opts.ContactID).
		Exec(ctx)

	return err
}

func (d *fundingRepo) GetContact(ctx context.Context, pipelineID, contactID uuid.UUID) (*malak.FundraiseContact, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var contact malak.FundraiseContact
	err := d.inner.NewSelect().
		Model(&contact).
		Relation("DealDetails").
		Relation("Contact").
		Where("fundraise_contact.id = ?", contactID).
		Where("fundraise_contact.fundraising_pipeline_id = ?", pipelineID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, malak.ErrContactNotFoundOnBoard
		}
		return nil, err
	}

	return &contact, nil
}

func (d *fundingRepo) GetColumn(ctx context.Context,
	opts malak.GetBoardOptions) (*malak.FundraisingPipelineColumn, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	column := new(malak.FundraisingPipelineColumn)

	err := d.inner.NewSelect().
		Model(column).
		Where("fundraising_pipeline_id = ?", opts.PipelineID).
		Where("id = ?", opts.ColumnID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, malak.ErrPipelineColumnNotFound
		}

		return nil, err
	}

	return column, nil
}

func (d *fundingRepo) MoveContactColumn(ctx context.Context, contact *malak.FundraiseContact,
	column *malak.FundraisingPipelineColumn) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		contact.UpdatedAt = time.Now()
		contact.FundraisingPipelineColumnID = column.ID

		_, err := tx.NewUpdate().
			Where("id = ?", contact.ID).
			Model(contact).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model(new(malak.FundraiseContactPosition)).
			Where("fundraising_pipeline_column_contact_id = ?", contact.ID).
			Set("order_index = ?", time.Now().Unix()).
			Exec(ctx)
		return err
	})
}
