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
