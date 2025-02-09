package main

import (
	"context"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func syncDataPointForIntegration(_ *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "sync",
		Short: `Sync integration data points`,
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, err := getLogger(hermes.DeRef(cfg))
			if err != nil {
				return err
			}

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			defer db.Close()

			logger.Debug("syncing datapoints")

			integrationRepo := postgres.NewIntegrationRepo(db)

			integrationManager, err := buildIntegrationManager(integrationRepo, hermes.DeRef(cfg), logger)
			if err != nil {
				logger.Fatal("could not build integration manager", zap.Error(err))
			}

			workspaces := make([]*malak.Workspace, 0)

			err = db.NewSelect().
				Model(&workspaces).
				Scan(context.Background())
			if err != nil {
				logger.Error("could not fetch workspaces", zap.Error(err))
				return err
			}

			for _, workspace := range workspaces {
				integrations, err := integrationRepo.List(cmd.Context(), workspace)
				if err != nil {
					logger.Error("could not fetch integrations for workspace",
						zap.String("workspace_id", workspace.ID.String()),
						zap.Error(err))
					continue
				}

				logger.Info("processing integrations for workspace",
					zap.String("workspace_id", workspace.ID.String()),
					zap.Int("integration_count", len(integrations)))

				for _, integration := range integrations {
					if !integration.IsEnabled {
						continue
					}

					client, err := integrationManager.Get(malak.IntegrationProvider(integration.Integration.IntegrationName))
					if err != nil {
						logger.Error("could not get integration client",
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
						continue
					}

					lastFetchedAt := time.Now()
					if integration.Metadata.LastFetchedAt != nil {
						lastFetchedAt = hermes.DeRef(integration.Metadata.LastFetchedAt)
					}

					dataPoints, err := client.Data(cmd.Context(), integration.Metadata.AccessToken, &malak.IntegrationFetchDataOptions{
						IntegrationID:      integration.ID,
						WorkspaceID:        workspace.ID,
						ReferenceGenerator: malak.NewReferenceGenerator(),
						LastFetchedAt:      lastFetchedAt,
					})
					if err != nil {
						logger.Error("could not fetch data points from integration",
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
						continue
					}

					logger.Info("fetched data points from integration",
						zap.String("workspace_id", workspace.ID.String()),
						zap.String("integration_id", integration.ID.String()),
						zap.Int("data_points_count", len(dataPoints)))

					integration.Metadata.LastFetchedAt = hermes.Ref(time.Now())

					if err := integrationRepo.Update(context.Background(), &integration); err != nil {
						logger.Error("could not fetch data points from integration",
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
					}
				}
			}

			return nil
		},
	}

	return cmd
}
