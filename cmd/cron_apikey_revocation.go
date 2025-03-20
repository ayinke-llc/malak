package main

import (
	"fmt"
	"os"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func revokeAPIKeys(c *cobra.Command, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-keys",
		Short: `delete api keys that has gone past their due time`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var logger *zap.Logger
			var err error

			switch cfg.Logging.Mode {
			case config.LogModeProd:
				logger, err = zap.NewProduction()
				if err != nil {
					fmt.Printf(`{"error":%s}`, err)
					os.Exit(1)
				}

			case config.LogModeDev:
				logger, err = zap.NewDevelopment()
				if err != nil {
					fmt.Printf(`{"error":%s}`, err)
					os.Exit(1)
				}
			}

			// ignoring on purpose
			h, _ := os.Hostname()

			logger = logger.With(zap.String("host", h),
				zap.String("app", "malak"),
				zap.String("component", "delete-keys"))

			cleanupOtelResources := server.InitOTELCapabilities(hermes.DeRef(cfg), logger)
			defer cleanupOtelResources()

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			defer db.Close()

			ctx := cmd.Context()

			_, err = db.ExecContext(ctx, `
				UPDATE api_keys 
				SET deleted_at = NOW()
				WHERE expires_at < NOW() 
				AND deleted_at IS NULL
			`)
			if err != nil {
				logger.Error("could not mark expired api keys as deleted",
					zap.Error(err))
				return err
			}

			logger.Info("successfully marked expired keys as deleted")
			return nil
		}}
}
