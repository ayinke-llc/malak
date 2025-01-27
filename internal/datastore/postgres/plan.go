package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
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

	if !hermes.IsStringEmpty(opts.Reference) {
		sel = sel.Where("reference = ?", opts.Reference)
	}

	if opts.ID != uuid.Nil {
		sel = sel.Where("id = ?", opts.ID)
	}

	err := sel.Model(plan).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, malak.ErrPlanNotFound
	}

	return plan, err
}

func (o *planRepo) List(ctx context.Context) ([]*malak.Plan, error) {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	apps := make([]*malak.Plan, 3)

	return apps, o.inner.NewSelect().
		Model(&apps).
		Scan(ctx)
}

func (p *planRepo) SetDefault(ctx context.Context,
	plan *malak.Plan) error {

	ctx, cancelFn := withContext(ctx)
	defer cancelFn()

	return p.inner.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {

			// set all plans to not be the default
			// then set this specific plan to be default
			_, err := tx.NewUpdate().
				Model(new(malak.Plan)).
				Where("is_default = ?", true).
				Set("is_default = ?", false).
				Exec(ctx)
			if err != nil {
				return err
			}

			_, err = tx.NewUpdate().
				Model(plan).
				Where("id = ?", plan.ID).
				Set("is_default = ?", true).
				Exec(ctx)
			return err
		})
}
