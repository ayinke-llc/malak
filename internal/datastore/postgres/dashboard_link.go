package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/hermes"
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
	opts *malak.CreateDashboardLinkOptions) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return d.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			contact := new(malak.Contact)
			link := opts.Link

			if !hermes.IsStringEmpty(opts.Email.String()) {
				err := tx.NewSelect().
					Model(contact).
					Where("email = ?", opts.Email.String()).
					Where("workspace_id = ?", opts.WorkspaceID).
					Scan(ctx)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return err
				}

				if err == nil {
					link.ContactID = contact.ID
				}

				if errors.Is(err, sql.ErrNoRows) {
					contact := &malak.Contact{
						WorkspaceID: opts.WorkspaceID,
						Email:       opts.Email,
						Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeContact),
					}

					_, err = tx.NewInsert().
						Model(contact).
						Exec(ctx)
					if err != nil {
						return err
					}
				}
			}

			// we can only have one default link type per dashboard at any given time
			if link.LinkType == malak.DashboardLinkTypeDefault {
				_, err := tx.NewDelete().
					Model(new(malak.DashboardLink)).
					Where("dashboard_id = ?", link.DashboardID).
					Where("link_type = ?", malak.DashboardLinkTypeDefault).
					Exec(ctx)
				if err != nil {
					return err
				}
			}

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
