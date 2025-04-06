package postgres

import (
	"context"
	"database/sql"
	"errors"

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
