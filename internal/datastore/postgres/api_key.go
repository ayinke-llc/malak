package postgres

import (
	"context"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/uptrace/bun"
)

type apiKeyImpl struct {
	inner *bun.DB
}

func NewAPIKeyRepository(db *bun.DB) malak.APIKeyRepository {
	return &apiKeyImpl{
		inner: db,
	}
}

func (r *apiKeyImpl) Create(ctx context.Context, apiKey *malak.APIKey) error {
	_, err := r.inner.NewInsert().Model(apiKey).Exec(ctx)
	return err
}

func (r *apiKeyImpl) Revoke(ctx context.Context, opts malak.RevokeAPIKeyOptions) error {
	now := time.Now()
	var expiresAt *time.Time

	switch opts.RevocationType {
	case malak.RevocationTypeImmediate:
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		expiresAt = &endOfDay
	case malak.RevocationTypeDay:
		tomorrow := now.AddDate(0, 0, 1)
		endOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, tomorrow.Location())
		expiresAt = &endOfDay
	case malak.RevocationTypeWeek:
		weekFromNow := now.AddDate(0, 0, 7)
		endOfDay := time.Date(weekFromNow.Year(), weekFromNow.Month(), weekFromNow.Day(), 23, 59, 59, 0, weekFromNow.Location())
		expiresAt = &endOfDay
	}

	_, err := r.inner.NewUpdate().
		Model(opts.APIKey).
		Set("expires_at = ?", expiresAt).
		Set("updated_at = ?", now).
		Where("id = ?", opts.APIKey.ID).
		Exec(ctx)

	return err
}
