package postgres

import (
	"context"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type templateRepo struct {
	inner *bun.DB
}

func NewTemplateRepository(db *bun.DB) malak.TemplateRepository {
	return &templateRepo{
		inner: db,
	}
}

func (t *templateRepo) System(ctx context.Context,
	filter malak.SystemTemplateFilter) ([]malak.SystemTemplate, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	templates := make([]malak.SystemTemplate, 0)

	q := t.inner.NewSelect().
		Model(&templates)

	switch filter {
	case malak.SystemTemplateFilterMostUsed:
		q = q.Order("number_of_uses DESC")
	case malak.SystemTemplateFilterRecentlyCreated:
		q = q.Order("created_at DESC")
	}

	return templates, q.Scan(ctx)
}
