package postgres

import (
	"context"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type preferenceRepo struct {
	inner *bun.DB
}

func NewPreferenceRepository(inner *bun.DB) malak.PreferenceRepository {
	return &preferenceRepo{
		inner: inner,
	}
}

func (w *preferenceRepo) Update(ctx context.Context,
	preferences *malak.Preference,
) error {
	_, err := w.inner.NewUpdate().
		Model(preferences).
		Where("id = ?", preferences.ID).
		Exec(ctx)
	return err
}

func (w *preferenceRepo) Get(ctx context.Context,
	workspace *malak.Workspace,
) (*malak.Preference, error) {
	preferences := &malak.Preference{}

	err := w.inner.NewSelect().
		Where("workspace_id = ?", workspace.ID).
		Model(preferences).Scan(ctx)

	return preferences, err
}
