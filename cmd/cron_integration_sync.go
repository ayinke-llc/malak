package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func syncDataPointForIntegration(_ *cobra.Command, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integrations",
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

			secretsProvider, err := buildSecretsProvider(cfg.Secrets.Provider, hermes.DeRef(cfg))
			if err != nil {
				logger.Fatal("could not build secrets provider", zap.Error(err))
			}

			workspaces := make([]*malak.Workspace, 0)

			err = db.NewSelect().
				Model(&workspaces).
				Scan(context.Background())
			if err != nil {
				logger.Error("could not fetch workspaces", zap.Error(err))
				return err
			}

			numWorkers := 10
			jobs := make(chan syncJob, len(workspaces)*3)

			var wg sync.WaitGroup

			for w := 1; w <= numWorkers; w++ {
				wg.Add(1)
				go worker(cmd.Context(), w, jobs, &wg, logger, integrationRepo, db)
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
					if !integration.Integration.IsEnabled || !integration.IsEnabled || !integration.IsActive {
						logger.Info("skipping integration",
							zap.String("workspace_integration_id", integration.ID.String()),
							zap.String("integration", integration.Integration.IntegrationName))
						continue
					}

					shouldProcess, err := shouldProcessIntegration(cmd.Context(), db, workspace, hermes.Ref(integration), true)
					if err != nil {
						logger.Error("error checking integration status",
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
						continue
					}

					if !shouldProcess {
						logger.Debug("skipping integration - already processed",
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()))
						continue
					}

					client, err := integrationManager.Get(
						malak.IntegrationProvider(
							strings.ToLower(integration.Integration.IntegrationName)))
					if err != nil {
						logger.Error("could not get integration client",
							zap.String("integration_name", integration.Integration.IntegrationName),
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
						continue
					}

					value, err := secretsProvider.Get(context.Background(), string(integration.Metadata.AccessToken))
					if err != nil {
						logger.Error("could not fetch value from secret vault",
							zap.String("secret_provider", cfg.Secrets.Provider.String()),
							zap.String("workspace_id", workspace.ID.String()),
							zap.String("integration_id", integration.ID.String()),
							zap.Error(err))
						continue
					}

					jobs <- syncJob{
						workspace:   workspace,
						integration: hermes.Ref(integration),
						client:      client,
						accessToken: malak.AccessToken(value),
					}
				}
			}

			close(jobs)
			wg.Wait()

			return nil
		},
	}

	return cmd
}

func shouldProcessIntegration(ctx context.Context,
	db *bun.DB,
	workspace *malak.Workspace,
	integration *malak.WorkspaceIntegration,
	resumeFailed bool) (bool, error) {

	now := time.Now()
	today := now.Format("2006-01-02")

	var checkpoint malak.IntegrationSyncCheckpoint
	err := db.NewSelect().
		Model(&checkpoint).
		Where("workspace_id = ? AND workspace_integration_id = ? AND DATE(created_at) = ?",
			workspace.ID, integration.ID, today).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)

	// No checkpoint exists for today - should process
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	// Only check these conditions if we have a checkpoint for today
	if resumeFailed {
		return checkpoint.Status == "failed" || checkpoint.Status == "pending", nil
	}

	if checkpoint.Status == "failed" {
		return true, nil
	}

	return false, nil
}

type syncJob struct {
	workspace   *malak.Workspace
	integration *malak.WorkspaceIntegration
	client      malak.IntegrationProviderClient
	accessToken malak.AccessToken
}

func worker(
	ctx context.Context,
	id int,
	jobs <-chan syncJob,
	wg *sync.WaitGroup,
	logger *zap.Logger,
	integrationRepo malak.IntegrationRepository,
	db *bun.DB,
) {
	defer wg.Done()

	for job := range jobs {
		logger.Info("worker processing integration",
			zap.Int("worker_id", id),
			zap.String("workspace_id", job.workspace.ID.String()),
			zap.String("integration_id", job.integration.ID.String()))

		generator := malak.NewReferenceGenerator()

		// create new checkpoint for today
		checkpoint := &malak.IntegrationSyncCheckpoint{
			WorkspaceID:            job.workspace.ID,
			WorkspaceIntegrationID: job.integration.ID,
			Status:                 "pending",
			LastSyncAttempt:        time.Now(),
			Reference:              generator.Generate(malak.EntityTypeIntegrationSyncCheckpoint),
		}

		_, err := db.NewInsert().
			Model(checkpoint).
			Exec(ctx)
		if err != nil {
			logger.Error("failed to create checkpoint",
				zap.Error(err))
			continue
		}

		dataPoints, err := job.client.Data(ctx,
			job.accessToken,
			&malak.IntegrationFetchDataOptions{
				IntegrationID:      job.integration.ID,
				WorkspaceID:        job.workspace.ID,
				ReferenceGenerator: generator,
				LastFetchedAt:      job.integration.Metadata.LastFetchedAt,
			})
		if err != nil {
			logger.Error("could not fetch data points from integration",
				zap.Int("worker_id", id),
				zap.String("workspace_id", job.workspace.ID.String()),
				zap.String("integration_id", job.integration.ID.String()),
				zap.Error(err))

			// update current checkpoint with error
			_, updateErr := db.NewUpdate().
				Model(checkpoint).
				Set("status = ?", "failed").
				Set("error_message = ?", err.Error()).
				Set("updated_at = NOW()").
				Where("id = ?", checkpoint.ID).
				Exec(ctx)
			if updateErr != nil {
				logger.Error("failed to update checkpoint", zap.Error(updateErr))
			}
			continue
		}

		logger.Info("fetched data points from integration",
			zap.Int("worker_id", id),
			zap.String("workspace_id", job.workspace.ID.String()),
			zap.String("integration_id", job.integration.ID.String()),
			zap.Int("data_points_count", len(dataPoints)))

		job.integration.Metadata.LastFetchedAt = time.Now()

		if err := integrationRepo.AddDataPoint(ctx, job.integration, dataPoints); err != nil {
			logger.Error("could not save data points from integration",
				zap.Int("worker_id", id),
				zap.String("workspace_id", job.workspace.ID.String()),
				zap.String("integration_id", job.integration.ID.String()),
				zap.Error(err))

			// update current checkpoint with error
			_, updateErr := db.NewUpdate().
				Model(checkpoint).
				Set("status = ?", "failed").
				Set("error_message = ?", err.Error()).
				Set("updated_at = NOW()").
				Where("id = ?", checkpoint.ID).
				Exec(ctx)
			if updateErr != nil {
				logger.Error("failed to update checkpoint", zap.Error(updateErr))
			}
			continue
		}

		// update current checkpoint as successful
		_, err = db.NewUpdate().
			Model(checkpoint).
			Set("status = ?", "success").
			Set("last_successful_sync = NOW()").
			Set("error_message = NULL").
			Set("updated_at = NOW()").
			Where("id = ?", checkpoint.ID).
			Exec(ctx)
		if err != nil {
			logger.Error("failed to update checkpoint", zap.Error(err))
		}
	}
}
