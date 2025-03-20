package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
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

func (r *apiKeyImpl) List(ctx context.Context, worskpaceID uuid.UUID) ([]malak.APIKey, error) {

	var apiKeys = make([]malak.APIKey, 0, 15)

	return apiKeys, r.inner.NewSelect().Model(&apiKeys).
		Where("workspace_id = ?", worskpaceID).
		Scan(ctx)
}

func (r *apiKeyImpl) Create(ctx context.Context, apiKey *malak.APIKey) error {
	return r.inner.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().Model(apiKey).Exec(ctx)
		if err != nil {
			return err
		}

		count, err := tx.NewSelect().
			Model(new(malak.APIKey)).
			Where("workspace_id = ?", apiKey.WorkspaceID).
			Count(ctx)
		if err != nil {
			return err
		}

		if count > 15 {
			return malak.ErrAPIKeyMaxLimit
		}

		return nil
	})
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
