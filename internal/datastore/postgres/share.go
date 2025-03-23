package postgres

import (
	"context"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type shareRepo struct {
	inner *bun.DB
}

func NewShareRepository(inner *bun.DB) malak.ContactShareRepository {
	return &shareRepo{
		inner: inner,
	}
}

func (o *shareRepo) All(ctx context.Context,
	contact *malak.Contact) ([]malak.ContactShareItem, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var sharedItems = make([]malak.ContactShareItem, 0)

	err := o.inner.NewSelect().
		Table("contact_shares").
		ColumnExpr("contact_shares.shared_at, contact_shares.item_type").
		ColumnExpr("COALESCE(updates.title, dashboards.title) AS title").
		ColumnExpr("COALESCE(updates.updated_at, dashboards.updated_at) AS updated_at").
		Column("item_reference", "shared_by").
		Join("LEFT JOIN updates ON contact_shares.item_type = 'update' AND contact_shares.item_id = updates.id").
		Join("LEFT JOIN dashboards ON contact_shares.item_type = 'dashboard' AND contact_shares.item_id = dashboards.id").
		Where("contact_shares.contact_id = ?", contact.ID).
		Scan(ctx, &sharedItems)

	return sharedItems, err
}

func (o *shareRepo) Overview(ctx context.Context, workspaceID uuid.UUID) (*malak.ShareOverview, error) {
	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	var sharedItems = make([]malak.ContactShareItem, 0)

	err := o.inner.NewSelect().
		Table("contact_shares").
		ColumnExpr("contact_shares.id, contact_shares.shared_at, contact_shares.item_type").
		ColumnExpr("COALESCE(updates.title, dashboards.title) AS title").
		ColumnExpr("COALESCE(updates.updated_at, dashboards.updated_at) AS updated_at").
		Column("item_reference", "shared_by").
		ColumnExpr("contacts.email AS email").
		ColumnExpr("contacts.first_name AS first_name").
		ColumnExpr("contacts.last_name AS last_name").
		Join("LEFT JOIN updates ON contact_shares.item_type = 'update' AND contact_shares.item_id = updates.id").
		Join("LEFT JOIN dashboards ON contact_shares.item_type = 'dashboard' AND contact_shares.item_id = dashboards.id").
		Join("JOIN contacts ON contact_shares.contact_id = contacts.id").
		Where("contacts.workspace_id = ?", workspaceID).
		Order("contact_shares.shared_at DESC").
		Limit(15).
		Scan(ctx, &sharedItems)

	if err != nil {
		return nil, err
	}

	return &malak.ShareOverview{
		RecentShares: sharedItems,
	}, nil
}
