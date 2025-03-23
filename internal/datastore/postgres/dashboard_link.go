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

					link.ContactID = contact.ID
				}
			}

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
				On("CONFLICT (contact_id,dashboard_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}

			if link.LinkType == malak.DashboardLinkTypeDefault {
				return nil
			}

			sharedItem := malak.ContactShare{
				Reference:     opts.Generator.Generate(malak.EntityTypeContactShare),
				SharedBy:      opts.UserID,
				ContactID:     link.ContactID,
				ItemType:      malak.ContactShareItemTypeDashboard,
				ItemID:        link.DashboardID,
				ItemReference: link.Dashboard.Reference,
			}

			_, err = tx.NewInsert().Model(&sharedItem).
				On("CONFLICT (item_reference,contact_id) DO NOTHING").
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

// same as Get without the workspace_id
// Separate api so as to not potentially misuse
func (d *dashboardLinkRepo) PublicDetails(ctx context.Context,
	ref malak.Reference) (malak.Dashboard, error) {

	ctx, cancel := withContext(ctx)
	defer cancel()

	link := new(malak.DashboardLink)
	err := d.inner.NewSelect().
		Model(link).
		Relation("Dashboard").
		Where("token = ?", ref.String()).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return malak.Dashboard{}, malak.ErrDashboardLinkNotFound
	}

	if err != nil {
		return malak.Dashboard{}, err
	}

	if link.Dashboard == nil {
		return malak.Dashboard{}, malak.ErrDashboardNotFound
	}

	return *link.Dashboard, nil
}

func (d *dashboardLinkRepo) List(ctx context.Context,
	opts malak.ListAccessControlOptions) ([]malak.DashboardLink, int64, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var links []malak.DashboardLink
	count := int64(0)

	totalCount, err := d.inner.NewSelect().
		Model(&links).
		Where("dashboard_id = ?", opts.DashboardID).
		Where("deleted_at IS NULL").
		Where("link_type = ?", malak.DashboardLinkTypeContact).
		Count(ctx)

	if err != nil {
		return nil, 0, err
	}

	count = int64(totalCount)

	return links, count, d.inner.NewSelect().
		Model(&links).
		Relation("Contact").
		Where("dashboard_id = ?", opts.DashboardID).
		Where("dashboard_link.deleted_at IS NULL").
		Order("dashboard_link.created_at DESC").
		Where("link_type = ?", malak.DashboardLinkTypeContact).
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx)
}

func (d *dashboardLinkRepo) Delete(ctx context.Context,
	dash malak.Dashboard, ref malak.Reference) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	_, err := d.inner.NewDelete().
		Model(new(malak.DashboardLink)).
		Where("reference = ?", ref).
		Where("dashboard_id = ?", dash.ID).
		Exec(ctx)

	return err
}
