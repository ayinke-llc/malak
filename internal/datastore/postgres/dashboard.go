package postgres

import (
	"context"
	"database/sql"
	"time"

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

			dashboard.ChartCount = 0

			_, err := tx.NewInsert().
				Model(dashboard).
				Exec(ctx)
			return err
		})
}

func (d *dashboardRepo) AddChart(ctx context.Context,
	dashboardChart *malak.DashboardChart) error {

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			_, err := tx.NewInsert().
				Model(dashboardChart).
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewUpdate().
				Model(new(malak.Dashboard)).
				Where("id = ?", dashboardChart.DashboardID).
				Set("updated_at = ?", time.Now()).
				Set("chart_count = chart_count + 1").
				Exec(ctx)
			return err
		})
}

func (d *dashboardRepo) List(ctx context.Context,
	opts malak.ListDashboardOptions) ([]malak.Dashboard, int64, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	dashboards := make([]malak.Dashboard, 0, opts.Paginator.PerPage)

	q := d.inner.NewSelect().
		Order("created_at DESC").
		Where("workspace_id = ?", opts.WorkspaceID)

	// Get total count with same filters
	total, err := q.
		Model(&dashboards).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = q.Model(&dashboards).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)

	return dashboards, int64(total), err
}
