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
	"golang.org/x/sync/errgroup"
)

func processDeckAnalytics(c *cobra.Command, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "decks-analytics",
		Short: `Process daily deck engagements and countries segments`,
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

			ctx := cmd.Context()

			g, gctx := errgroup.WithContext(ctx)

			// Process daily engagements concurrently
			g.Go(func() error {
				logger.Info("processing daily deck engagements")

				engagementQuery := `
					WITH dates_to_process AS (
						SELECT DISTINCT
							dvs.deck_id,
							d.workspace_id,
							DATE(dvs.viewed_at) as engagement_date
						FROM deck_viewer_sessions dvs
						INNER JOIN decks d ON d.id = dvs.deck_id
						WHERE dvs.deleted_at IS NULL
						AND d.deleted_at IS NULL
						AND DATE(dvs.viewed_at) = CURRENT_DATE
						AND NOT EXISTS (
							SELECT 1 FROM deck_daily_engagements dde
							WHERE dde.deck_id = dvs.deck_id
							AND dde.workspace_id = d.workspace_id
							AND dde.engagement_date = CURRENT_DATE
						)
					),
					daily_stats AS (
						SELECT 
							dvs.deck_id,
							d.workspace_id,
							CURRENT_DATE as engagement_date,
							COUNT(DISTINCT dvs.id) as engagement_count
						FROM deck_viewer_sessions dvs
						INNER JOIN decks d ON d.id = dvs.deck_id
						WHERE dvs.deleted_at IS NULL
						AND d.deleted_at IS NULL
						AND DATE(dvs.viewed_at) = CURRENT_DATE
						GROUP BY dvs.deck_id, d.workspace_id
					)
					INSERT INTO deck_daily_engagements (
						reference,
						deck_id,
						workspace_id,
						engagement_count,
						engagement_date
					)
					SELECT 
						'deck_daily_engagement_' || LOWER(REPLACE(uuid_generate_v4()::text, '-', '')),
						deck_id,
						workspace_id,
						engagement_count,
						engagement_date
					FROM daily_stats
					ON CONFLICT (deck_id, workspace_id, engagement_date)
					DO UPDATE SET
						engagement_count = EXCLUDED.engagement_count,
						updated_at = CURRENT_TIMESTAMP
				`

				_, err := db.ExecContext(gctx, engagementQuery)
				if err != nil {
					logger.Error("failed to process daily engagements",
						zap.Error(err))
					return err
				}

				logger.Info("successfully processed daily engagements")
				return nil
			})

			// Process geographic stats concurrently
			g.Go(func() error {
				logger.Info("processing geographic stats")

				geoQuery := `
					WITH geo_stats AS (
						SELECT 
							deck_id,
							COALESCE(NULLIF(TRIM(country), ''), 'Unknown') as country,
							COUNT(DISTINCT id) as view_count,
							CURRENT_DATE as stat_date
						FROM deck_viewer_sessions
						WHERE deleted_at IS NULL
						GROUP BY deck_id, COALESCE(NULLIF(TRIM(country), ''), 'Unknown')
					)
					INSERT INTO deck_geographic_stats (
						reference,
						deck_id,
						country,
						view_count,
						stat_date
					)
					SELECT 
						'deck_geographic_stat_' || LOWER(REPLACE(uuid_generate_v4()::text, '-', '')),
						deck_id,
						country,
						view_count,
						stat_date
					FROM geo_stats
					ON CONFLICT (deck_id, country, stat_date)
					DO UPDATE SET
						view_count = EXCLUDED.view_count,
						updated_at = CURRENT_TIMESTAMP
				`

				_, err := db.ExecContext(gctx, geoQuery)
				if err != nil {
					logger.Error("failed to process geographic stats",
						zap.Error(err))
					return err
				}

				logger.Info("successfully processed geographic stats")
				return nil
			})

			// Wait for all goroutines to complete and check for errors
			if err := g.Wait(); err != nil {
				logger.Error("failed to process deck analytics", zap.Error(err))
				return err
			}

			logger.Info("successfully processed all deck analytics")
			return nil
		}}
}
