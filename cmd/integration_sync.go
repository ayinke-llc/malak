package main

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type SyncCheckpoint struct {
	ID                 uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID        uuid.UUID `bun:"workspace_id,notnull" json:"workspace_id,omitempty"`
	IntegrationID      uuid.UUID `bun:"integration_id,notnull" json:"integration_id,omitempty"`
	LastSyncAttempt    time.Time `bun:"last_sync_attempt" json:"last_sync_attempt,omitempty"`
	LastSuccessfulSync time.Time `bun:"last_successful_sync" json:"last_successful_sync,omitempty"`
	Status             string    `bun:"status,notnull" json:"status,omitempty"`
	ErrorMessage       string    `bun:"error_message" json:"error_message,omitempty"`
	CreatedAt          time.Time `bun:"created_at,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt          time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at,omitempty"`
}

func shouldProcessIntegration(ctx context.Context, db *bun.DB, workspace *malak.Workspace, integration *malak.WorkspaceIntegration, resumeFailed bool) (bool, error) {
	return true, nil
	var checkpoint SyncCheckpoint
	err := db.NewSelect().
		Model(&checkpoint).
		Where("workspace_id = ? AND integration_id = ?",
			workspace.ID, integration.ID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	// if resuming failed syncs, only process failed or pending ones
	if resumeFailed {
		return checkpoint.Status == "failed" || checkpoint.Status == "pending", nil
	}

	// No checkpoint exists
	if err == sql.ErrNoRows {
		return true, nil
	}

	// Last sync was too long ago (e.g. within last 10 hours)
	if time.Since(checkpoint.LastSuccessfulSync) > 24*time.Hour {
		return true, nil
	}

	// previous sync failed
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
) {
	defer wg.Done()

	for job := range jobs {
		logger.Info("worker processing integration",
			zap.Int("worker_id", id),
			zap.String("workspace_id", job.workspace.ID.String()),
			zap.String("integration_id", job.integration.ID.String()))

		dataPoints, err := job.client.Data(ctx,
			job.accessToken,
			&malak.IntegrationFetchDataOptions{
				IntegrationID:      job.integration.ID,
				WorkspaceID:        job.workspace.ID,
				ReferenceGenerator: malak.NewReferenceGenerator(),
				LastFetchedAt:      job.integration.Metadata.LastFetchedAt,
			})
		if err != nil {
			logger.Error("could not fetch data points from integration",
				zap.Int("worker_id", id),
				zap.String("workspace_id", job.workspace.ID.String()),
				zap.String("integration_id", job.integration.ID.String()),
				zap.Error(err))
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
		}
	}
}

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

			secretsProvider, err := buildSecretsProvider(hermes.DeRef(cfg))
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
				go worker(cmd.Context(), w, jobs, &wg, logger, integrationRepo)
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
