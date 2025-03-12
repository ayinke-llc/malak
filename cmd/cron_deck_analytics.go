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

func processDeckAnalytics(c *cobra.Command, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "decks-analytics",
		Short: `Send scheduled updates`,
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
				zap.String("component", "decks-analytics"))

			cleanupOtelResources := server.InitOTELCapabilities(hermes.DeRef(cfg), logger)
			defer cleanupOtelResources()

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			defer db.Close()

			logger.Debug("syncing datapoints")

			deckRepo := postgres.NewDeckRepository(db)
			_ = deckRepo

			return nil
		}}
}
