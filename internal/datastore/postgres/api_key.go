package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (a *apiKeyImpl) Fetch(ctx context.Context, opts malak.FetchAPIKeyOptions) (
	*malak.APIKey, error) {

	var apiKey = new(malak.APIKey)

	err := a.inner.NewSelect().
		Model(apiKey).
		Where("workspace_id = ?", opts.WorkspaceID).
		Where("reference = ?", opts.Reference).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		err = malak.ErrAPIKeyNotFound
	}

	return apiKey, err
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

	q := r.inner.NewUpdate().
		Model(opts.APIKey).
		Set("updated_at = ?", now)

	switch opts.RevocationType {
	case malak.RevocationTypeImmediate:
		q = q.Set("expires_at = CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 hour'")
		q = q.Set("deleted_at = NOW()")
	case malak.RevocationTypeDay:
		q = q.Set("expires_at = CURRENT_DATE + INTERVAL '2 days' - INTERVAL '1 hour'")
	case malak.RevocationTypeWeek:
		q = q.Set("expires_at = CURRENT_DATE + INTERVAL '8 days' - INTERVAL '1 hour'")
	}

	q = q.Where("id = ?", opts.APIKey.ID)

	_, err := q.Exec(ctx)
	return err
}
