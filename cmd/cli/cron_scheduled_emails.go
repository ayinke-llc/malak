package cli

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/adelowo/gulter"
	"github.com/adelowo/gulter/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/chart"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/server"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

const (
	defaultBatchSize         = 100
	defaultWorkerCount       = 5
	defaultMaxRetries        = 3
	defaultProcessingTimeout = 30 * time.Minute
	defaultEmailTimeout      = 10 * time.Second
)

type ProcessingState string

const (
	StateInit       ProcessingState = "init"
	StateProcessing ProcessingState = "processing"
	StateFailed     ProcessingState = "failed"
	StateCompleted  ProcessingState = "completed"
)

type ProcessMetrics struct {
	TotalEmails     int64
	SentEmails      int64
	FailedEmails    int64
	LastProcessedID string
	StartTime       time.Time
	EndTime         time.Time
	mu              sync.Mutex
}

func (m *ProcessMetrics) IncrementSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentEmails++
}

func (m *ProcessMetrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FailedEmails++
}

type EmailJob struct {
	Recipient   recipient
	Title       string
	Content     string
	RetryCount  int
	LastError   error
	LastAttempt time.Time
}

type EmailProcessor struct {
	db            *bun.DB
	emailClient   email.Client
	logger        *zap.Logger
	tracer        trace.Tracer
	cfg           *config.Config
	metrics       *ProcessMetrics
	rateLimiter   *rate.Limiter
	workspaceRepo malak.WorkspaceRepository
	chartRenderer malak.ChartRenderer
	storage       gulter.Storage
}

type ProcessorOptions struct {
	BatchSize         int
	WorkerCount       int
	MaxRetries        int
	ProcessingTimeout time.Duration
	EmailTimeout      time.Duration
	RateLimit         int
}

func DefaultProcessorOptions() ProcessorOptions {
	return ProcessorOptions{
		BatchSize:         defaultBatchSize,
		WorkerCount:       defaultWorkerCount,
		MaxRetries:        defaultMaxRetries,
		ProcessingTimeout: defaultProcessingTimeout,
		EmailTimeout:      defaultEmailTimeout,
		RateLimit:         10,
	}
}

type recipient struct {
	malak.UpdateRecipient
	Contact       *malak.Contact `json:"contact" bun:"rel:has-one,join:contact_id=id"`
	bun.BaseModel `bun:"table:update_recipients"`
}

func NewEmailProcessor(db *bun.DB, emailClient email.Client, logger *zap.Logger, tracer trace.Tracer, cfg *config.Config, storage gulter.Storage, opts ProcessorOptions) *EmailProcessor {
	return &EmailProcessor{
		db:            db,
		emailClient:   emailClient,
		logger:        logger,
		tracer:        tracer,
		cfg:           cfg,
		metrics:       &ProcessMetrics{StartTime: time.Now()},
		rateLimiter:   rate.NewLimiter(rate.Limit(opts.RateLimit), opts.RateLimit),
		workspaceRepo: postgres.NewWorkspaceRepository(db),
		chartRenderer: chart.NewEChartsRenderer(storage, hermes.DeRef(cfg), db, postgres.NewIntegrationRepo(db)),
		storage:       storage,
	}
}

func (p *EmailProcessor) acquireProcessingLock(ctx context.Context, updateID uuid.UUID) (bool, error) {
	var locked bool
	err := p.db.NewSelect().
		TableExpr("update_schedules").
		ColumnExpr("TRUE").
		Where("id = ? AND status = ?", updateID, malak.UpdateSendScheduleScheduled).
		For("UPDATE SKIP LOCKED").
		Scan(ctx, &locked)
	return locked, err
}

func fetchPendingUpdates(ctx context.Context, db *bun.DB, lastProcessedID string, limit int) ([]*malak.UpdateSchedule, error) {
	var scheduledUpdates []*malak.UpdateSchedule
	query := db.NewSelect().
		Model(&scheduledUpdates).
		Where("status = ?", malak.UpdateSendScheduleScheduled).
		Limit(limit)

	if lastProcessedID != "" {
		query = query.Where("id > ?", lastProcessedID)
	}

	err := query.Order("id ASC").Scan(ctx)
	return scheduledUpdates, err
}

func (p *EmailProcessor) processUpdate(ctx context.Context, update *malak.UpdateSchedule) error {
	ctx, span := p.tracer.Start(ctx, "processUpdate")
	defer span.End()

	dbCtx, dbCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer dbCancel()

	locked, err := p.acquireProcessingLock(dbCtx, update.ID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		return errors.New("update is being processed by another instance")
	}

	if err := updateScheduleStatus(dbCtx, p.db, update, malak.UpdateSendScheduleProcessing); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	updateDetails, err := fetchUpdateDetails(dbCtx, p.db, update.UpdateID)
	if err != nil {
		return fmt.Errorf("failed to fetch update details: %w", err)
	}

	emailCtx, emailCancel := context.WithTimeout(ctx, defaultProcessingTimeout)
	defer emailCancel()

	g, emailCtx := errgroup.WithContext(emailCtx)

	for {
		recipients, err := p.fetchNextBatch(dbCtx, update.UpdateID)
		if err != nil {
			return fmt.Errorf("failed to fetch recipients: %w", err)
		}
		if len(recipients) == 0 {
			break
		}

		jobs := p.createEmailJobs(recipients, updateDetails)
		p.metrics.TotalEmails += int64(len(jobs))

		results := make(chan error, len(jobs))
		g.Go(func() error {
			return p.processEmailBatch(emailCtx, jobs, results)
		})

		for i := 0; i < len(jobs); i++ {
			if err := <-results; err != nil {
				p.logger.Error("email processing failed",
					zap.Error(err),
					zap.String("update_id", update.ID.String()))
			}
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("batch processing failed: %w", err)
	}

	if p.metrics.FailedEmails > 0 {
		if err := updateScheduleStatus(dbCtx, p.db, update, malak.UpdateSendScheduleFailed); err != nil {
			return fmt.Errorf("failed to update final status: %w", err)
		}
		return fmt.Errorf("some emails failed to send: %d/%d", p.metrics.FailedEmails, p.metrics.TotalEmails)
	}

	return updateScheduleStatus(dbCtx, p.db, update, malak.UpdateSendScheduleSent)
}

func (p *EmailProcessor) fetchNextBatch(ctx context.Context, updateID uuid.UUID) ([]recipient, error) {
	var recipients []recipient
	err := p.db.NewSelect().
		Model(&recipients).
		Where("update_id = ?", updateID).
		Where("status = ?", malak.RecipientStatusPending).
		Relation("Contact").
		Limit(defaultBatchSize).
		Scan(ctx)
	return recipients, err
}

func (p *EmailProcessor) createEmailJobs(recipients []recipient, update *malak.Update) []*EmailJob {
	workspace, err := p.workspaceRepo.Get(context.Background(), &malak.FindWorkspaceOptions{
		ID: update.WorkspaceID,
	})
	if err != nil {
		// the workspace is supposed to exist so this is okay to do
		panic(err.Error())
	}

	content, err := prepareEmailTemplate(update, workspace.WorkspaceName, p.chartRenderer)
	if err != nil {
		// the template is supposed to be fine so this is okay to do
		panic(err.Error())
	}

	jobs := make([]*EmailJob, len(recipients))
	for i, r := range recipients {
		jobs[i] = &EmailJob{
			Recipient:  r,
			Title:      update.Title,
			Content:    content,
			RetryCount: 0,
		}
	}
	return jobs
}

func (p *EmailProcessor) processEmailBatch(ctx context.Context, jobs []*EmailJob, results chan<- error) error {
	workerPool := make(chan struct{}, defaultWorkerCount)
	var wg sync.WaitGroup

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case workerPool <- struct{}{}:
			wg.Add(1)
			go func(j *EmailJob) {
				defer wg.Done()
				defer func() { <-workerPool }()

				err := p.processEmailWithRetry(ctx, j)
				if err != nil {
					p.metrics.IncrementFailed()
					results <- err
					return
				}
				p.metrics.IncrementSent()
				results <- nil
			}(job)
		}
	}

	wg.Wait()
	return nil
}

func (p *EmailProcessor) processEmailWithRetry(ctx context.Context, job *EmailJob) error {
	var err error
	for attempt := 0; attempt <= defaultMaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}

		err = p.rateLimiter.Wait(ctx)
		if err != nil {
			return fmt.Errorf("rate limiter error: %w", err)
		}

		emailCtx, cancel := context.WithTimeout(ctx, defaultEmailTimeout)
		err = p.sendEmail(emailCtx, job)
		cancel()

		if err == nil {
			return nil
		}

		job.RetryCount++
		job.LastError = err
		job.LastAttempt = time.Now()

		p.logger.Warn("email send failed, will retry",
			zap.Error(err),
			zap.Int("attempt", attempt+1),
			zap.String("recipient", job.Recipient.Contact.Email.String()))
	}

	return fmt.Errorf("max retries exceeded: %w", err)
}

func (p *EmailProcessor) sendEmail(ctx context.Context, job *EmailJob) error {
	opts := email.SendOptions{
		HTML:      job.Content,
		Sender:    p.cfg.Email.Sender,
		Recipient: job.Recipient.Contact.Email,
		Subject:   job.Title,
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign:       false,
			PrivateKey: []byte(""),
		},
	}

	emailID, err := p.emailClient.Send(ctx, opts)
	if err != nil {
		return updateRecipientStatus(ctx, p.db, p.emailClient, job.Recipient, malak.RecipientStatusFailed, "")
	}

	return updateRecipientStatus(ctx, p.db, p.emailClient, job.Recipient, malak.RecipientStatusSent, emailID)
}

func sendScheduledUpdates(c *cobra.Command, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "updates",
		Short: `Send scheduled updates`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := setupLogger(cfg)
			if err != nil {
				return fmt.Errorf("failed to setup logger: %w", err)
			}

			cleanupOtel := server.InitOTELCapabilities(hermes.DeRef(cfg), logger)
			defer cleanupOtel()

			emailClient, err := getEmailProvider(hermes.DeRef(cfg))
			if err != nil {
				return fmt.Errorf("failed to setup email client: %w", err)
			}
			defer emailClient.Close()

			db, err := postgres.New(cfg, logger)
			if err != nil {
				return fmt.Errorf("failed to setup database: %w", err)
			}
			defer db.Close()

			tracer := otel.Tracer("malak.cron")
			ctx, span := tracer.Start(context.Background(), "updates-send")
			defer span.End()

			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: !cfg.Uploader.S3.UseTLS,
					},
				},
			}

			s3Config, err := awsConfig.LoadDefaultConfig(
				context.Background(),
				awsConfig.WithRegion(cfg.Uploader.S3.Region),
				awsConfig.WithHTTPClient(httpClient),
				awsConfig.WithCredentialsProvider(
					awsCreds.NewStaticCredentialsProvider(
						cfg.Uploader.S3.AccessKey,
						cfg.Uploader.S3.AccessSecret,
						"")),
				//nolint:staticcheck
				awsConfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					//nolint:staticcheck
					return aws.Endpoint{
						URL:               cfg.Uploader.S3.Endpoint,
						SigningRegion:     cfg.Uploader.S3.Region,
						HostnameImmutable: true,
					}, nil
				})),
			)
			if err != nil {
				logger.Fatal("could not set up S3 config",
					zap.Error(err))
			}

			s3Store, err := storage.NewS3FromConfig(s3Config, storage.S3Options{
				DebugMode:        cfg.Uploader.S3.LogOperations,
				UsePathStyle:     true,
				Bucket:           cfg.Uploader.S3.Bucket,
				CloudflareDomain: cfg.Uploader.S3.CloudflareBucketDomain,
			})
			if err != nil {
				logger.Fatal("could not set up S3 client",
					zap.Error(err))
			}

			opts := DefaultProcessorOptions()
			processor := NewEmailProcessor(db, emailClient, logger, tracer, cfg, s3Store, opts)

			var lastProcessedID string
			for {
				updates, err := fetchPendingUpdates(ctx, db, lastProcessedID, opts.BatchSize)
				if err != nil {
					return fmt.Errorf("failed to fetch updates: %w", err)
				}

				if len(updates) == 0 {
					break
				}

				for _, update := range updates {
					if err := processor.processUpdate(ctx, update); err != nil {
						logger.Error("failed to process update",
							zap.Error(err),
							zap.String("update_id", update.ID.String()))
						continue
					}
					lastProcessedID = update.ID.String()
				}
			}

			logger.Info("email processing completed",
				zap.Int64("total_emails", processor.metrics.TotalEmails),
				zap.Int64("sent_emails", processor.metrics.SentEmails),
				zap.Int64("failed_emails", processor.metrics.FailedEmails),
				zap.Duration("duration", time.Since(processor.metrics.StartTime)))

			return nil
		},
	}
}

func setupLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch cfg.Logging.Mode {
	case config.LogModeProd:
		logger, err = zap.NewProduction()
	case config.LogModeDev:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	hostname, _ := os.Hostname()
	return logger.With(
		zap.String("host", hostname),
		zap.String("app", "malak"),
		zap.String("component", "cron.updates-send"),
	), nil
}

func updateScheduleStatus(ctx context.Context, db *bun.DB, schedule *malak.UpdateSchedule,
	status malak.UpdateSendSchedule) error {
	schedule.Status = status
	_, err := db.NewUpdate().
		Model(schedule).
		Where("id = ?", schedule.ID).
		Exec(ctx)
	return err
}

func fetchUpdateDetails(ctx context.Context, db *bun.DB, updateID uuid.UUID) (*malak.Update, error) {
	update := &malak.Update{}
	err := db.NewSelect().
		Model(update).
		Where("id = ?", updateID).
		Scan(ctx)
	return update, err
}

func prepareEmailTemplate(update *malak.Update, workspaceName string, renderer malak.ChartRenderer) (string, error) {
	tmpl, err := template.New("template").Parse(email.UpdateHTMLEmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"Content": update.Content.HTML(update.WorkspaceID, renderer),
		"Company": workspaceName,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func updateRecipientStatus(ctx context.Context, db *bun.DB, emailClient email.Client, r recipient, status malak.RecipientStatus, emailID string) error {
	return db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().
			Model(&malak.UpdateRecipient{}).
			Where("reference = ?", r.Reference.String()).
			Set("status = ?", status).
			Exec(ctx)
		if err != nil {
			return err
		}

		if emailID != "" {
			emailLog := &malak.UpdateRecipientLog{
				ProviderID:  emailID,
				RecipientID: r.ID,
				Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeRecipientLog),
				Provider:    emailClient.Name(),
				ID:          uuid.New(),
			}

			if _, err := tx.NewInsert().Model(emailLog).Exec(ctx); err != nil {
				return err
			}
		}

		stats := &malak.UpdateRecipientStat{
			Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeRecipientStat),
			RecipientID: r.ID,
			HasReaction: false,
			IsDelivered: status == malak.RecipientStatusSent,
		}

		_, err = tx.NewInsert().Model(stats).Exec(ctx)
		return err
	})
}
