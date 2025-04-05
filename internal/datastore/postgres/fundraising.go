package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
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
