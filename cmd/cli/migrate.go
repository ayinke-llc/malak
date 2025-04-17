package cli

import (
	"errors"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

func runMigrations(logger *zap.Logger, cfg *config.Config) error {

	d, err := iofs.New(malak.Migrations, "internal/datastore/postgres/migrations")
	if err != nil {
		logger.Error("could not set up embedded migrations", zap.Error(err))
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cfg.Database.Postgres.DSN)
	if err != nil {
		logger.Error("could not set up migrations", zap.Error(err))
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("no new migration to run")
		return nil
	}

	if err != nil {
		logger.Error("could not run migrations", zap.Error(err))
		return err
	}

	logger.Info("migrations successful")
	return nil
}
