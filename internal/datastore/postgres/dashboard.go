package postgres

import (
	"context"
	"database/sql"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type dashboardRepo struct {
	inner *bun.DB
}

func NewDashboardRepo(inner *bun.DB) malak.DashboardRepository {
	return &dashboardRepo{
		inner: inner,
	}
}

func (d *dashboardRepo) Create(ctx context.Context,
	dashboard *malak.Dashboard) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewInsert().
				Model(dashboard).
				Exec(ctx)
			return err
		})
}
