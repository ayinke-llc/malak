package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
)

// TODO: this is horrible for sure
var timeout time.Duration

func withContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

func New(cfg *config.Config, logger *logrus.Entry) (*bun.DB, error) {

	pgdb := sql.OpenDB(
		pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.Postgres.DSN)))

	db := bun.NewDB(pgdb, pgdialect.New())

	if cfg.Database.Postgres.LogQueries {
		db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
			Logger:     logger,
			QueryLevel: logger.Level,
			SlowLevel:  logger.Level,
			ErrorLevel: logger.Level,
		}))
	}

	if cfg.Otel.IsEnabled {
		db.AddQueryHook(
			bunotel.NewQueryHook(bunotel.WithDBName("malak.database")))
	}

	timeout = cfg.Database.Postgres.QueryTimeout
	return db, db.Ping()
}
