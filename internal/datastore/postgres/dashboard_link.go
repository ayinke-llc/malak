package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type dashboardLinkRepo struct {
	inner *bun.DB
}

func NewDashboardLinkRepo(inner *bun.DB) malak.DashboardLinkRepository {
	return &dashboardLinkRepo{
		inner: inner,
	}
}

func (d *dashboardLinkRepo) Create(ctx context.Context,
	link *malak.DashboardLink) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().
				Model(link).
				Exec(ctx)
			return err
		})
}

func (d *dashboardLinkRepo) DefaultLink(ctx context.Context,
	dash *malak.Dashboard) (malak.DashboardLink, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	link := malak.DashboardLink{}

	err := d.inner.NewSelect().
		Model(&link).
		Where("dashboard_id = ?", dash.ID).
		Where("link_type = ?", malak.DashboardLinkTypeDefault).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrDashboardLinkNotFound
	}

	return link, err
}
