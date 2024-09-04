package postgres

import (
	"context"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type planRepo struct {
	inner *bun.DB
}

func NewPlanRepository(db *bun.DB) *planRepo {
	return &planRepo{
		inner: db,
	}
}

func (p *planRepo) Get(ctx context.Context, opts *malak.FetchPlanOptions) (*malak.Plan, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	plan := &malak.Plan{}

	sel := p.inner.NewSelect()

	if !util.IsStringEmpty(opts.Reference) {
		sel = sel.Where("reference = ?", opts.Reference)
	}

	if opts.ID != uuid.Nil {
		sel = sel.Where("id = ?", opts.ID)
	}

	return plan, sel.Model(plan).Scan(ctx)
}

func (o *planRepo) List(ctx context.Context) ([]*malak.Plan, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	apps := make([]*malak.Plan, 3)

	return apps, o.inner.NewSelect().
		Model(&apps).
		Scan(ctx)
}
