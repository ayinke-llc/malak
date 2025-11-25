package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	testfixtures "github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/uptrace/bun"
	"go.uber.org/zap"

	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/ayinke-llc/malak/config"
)

func getConfig(dsn string) *config.Config {
	return &config.Config{
		Logging: struct {
			Mode config.LogMode "yaml:\"mode\" mapstructure:\"mode\" json:\"mode\""
		}{
			Mode: config.LogModeDev,
		},
		Database: struct {
			DatabaseType config.DatabaseType "yaml:\"database_type\" mapstructure:\"database_type\" json:\"database_type\""
			Postgres     struct {
				DSN          string        "yaml:\"dsn\" mapstructure:\"dsn\" json:\"dsn\""
				LogQueries   bool          "yaml:\"log_queries\" mapstructure:\"log_queries\" json:\"log_queries\""
				QueryTimeout time.Duration "yaml:\"query_timeout\" mapstructure:\"query_timeout\" json:\"query_timeout\""
			} "yaml:\"postgres\" mapstructure:\"postgres\" json:\"postgres\""
			Redis struct {
				DSN string "yaml:\"dsn\" mapstructure:\"dsn\" json:\"dsn\""
			} "yaml:\"redis\" mapstructure:\"redis\" json:\"redis\""
		}{
			DatabaseType: config.DatabaseTypePostgres,
			Postgres: struct {
				DSN          string        "yaml:\"dsn\" mapstructure:\"dsn\" json:\"dsn\""
				LogQueries   bool          "yaml:\"log_queries\" mapstructure:\"log_queries\" json:\"log_queries\""
				QueryTimeout time.Duration "yaml:\"query_timeout\" mapstructure:\"query_timeout\" json:\"query_timeout\""
			}{
				DSN:          dsn,
				LogQueries:   true,
				QueryTimeout: time.Second * 10,
			},
		},
	}
}

func prepareTestDatabase(t *testing.T, dsn string) {
	t.Helper()

	var err error

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"), "postgres", driver)
	require.NoError(t, err)

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("testdata/fixtures"),
	)
	require.NoError(t, err)

	err = fixtures.Load()
	require.NoError(t, err)
}

// setupDatabase spins up a new Postgres container and returns a closure
// please always make sure to call the closure as it is the teardown function
func setupDatabase(t *testing.T) (*bun.DB, func()) {
	t.Helper()

	os.Setenv("TZ", "")

	dbName := "malaktest"

	postgresContainer, err := testpostgres.Run(
		t.Context(),
		"postgres:18",
		testpostgres.WithDatabase(dbName),
		testpostgres.WithUsername(dbName),
		testpostgres.WithPassword(dbName),
		testpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	dsn, err := postgresContainer.ConnectionString(t.Context(), "sslmode=disable")
	require.NoError(t, err)

	prepareTestDatabase(t, dsn)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	db, err := New(getConfig(dsn), logger)
	require.NoError(t, err)

	return db, func() {
		require.NoError(t, testcontainers.TerminateContainer(postgresContainer))
	}
}
