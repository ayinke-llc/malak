package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/server"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func addCronCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use: "cron",
	}

	cmd.AddCommand(sendScheduledUpdates(c, cfg))

	c.AddCommand(cmd)
}

// TODO(adelowo): test at scale before beta mvp release. Email rate scale and errors syncing
func sendScheduledUpdates(c *cobra.Command, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updates-send",
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
				zap.String("component", "cron.updates-send"))

			cleanupOtelResources := server.InitOTELCapabilities(hermes.DeRef(cfg), logger)
			defer cleanupOtelResources()

			emailClient, err := getEmailProvider(*cfg)
			if err != nil {
				logger.Fatal("could not set up email client",
					zap.Error(err))
			}
			defer emailClient.Close()

			// Doing this here as I do not see another place where we
			// need to reuse this bit of code
			// if for some reason, we have to send updates from
			// another place.
			// Just encapsulate this so we don't duplicate code
			//
			// LOGIC
			//
			// 1. Fetch pending scheduled updates
			// 2. Fetch the recipients of each updates
			// 3. Send to each recipient

			var tracer = otel.Tracer("malak.cron")

			ctx, span := tracer.Start(context.Background(), "updates-send")
			defer span.End()

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Fatal("could not set up database connection",
					zap.Error(err))
			}

			defer db.Close()

			var scheduledUpdates = make([]*malak.UpdateSchedule, 0)

			err = db.NewSelect().Model(&scheduledUpdates).
				Limit(30).
				Where("status = ?", malak.UpdateSendScheduleScheduled).
				Scan(ctx)
			if err != nil {
				span.RecordError(err)
				logger.Error("could not fetch pending scheduled updates from db to process",
					zap.Error(err))
				return err
			}

			if len(scheduledUpdates) == 0 {
				logger.Info("no scheduled updates to process")
				os.Exit(0)
			}

			span.SetAttributes(
				attribute.Int("number_of_scheduled_updates", len(scheduledUpdates)))

			updateRepo := postgres.NewUpdatesRepository(db)

			g, ctx := errgroup.WithContext(ctx)
			var wg sync.WaitGroup

			for _, scheduledUpdate := range scheduledUpdates {
				g.Go(func() error {

					ctx, span := tracer.Start(ctx, "sending")
					defer span.End()

					scheduledUpdate.Status = malak.UpdateSendScheduleProcessing
					_, err = db.NewUpdate().Model(scheduledUpdate).Where("id = ? ", scheduledUpdate.ID).
						Exec(ctx)
					if err != nil {
						span.RecordError(err)
						logger.Error("could not update schedule status to processing",
							zap.Error(err))
						return err
					}

					update := &malak.Update{}

					err := db.NewSelect().Model(update).
						Where("id = ?", scheduledUpdate.UpdateID).
						Scan(ctx)
					if err != nil {
						span.RecordError(err)
						logger.Error("could not fetch update from database",
							zap.Error(err))
						return err
					}

					templatedFile, err := template.New("template").
						Parse(email.UpdateHTMLEmailTemplate)
					if err != nil {
						span.RecordError(err)
						logger.Error("could not create html template",
							zap.Error(err))
						return err
					}

					var b = new(bytes.Buffer)
					err = templatedFile.Execute(b, map[string]string{
						"Content": update.Content.HTML(),
					})

					if err != nil {
						span.RecordError(err)
						logger.Error("could not parse html template",
							zap.Error(err))
						return err
					}

					logger = logger.With(zap.String("update_id", update.ID.String()))

					schedule, err := updateRepo.GetSchedule(ctx, scheduledUpdate.ID)
					if err != nil {
						span.RecordError(err)
						logger.Error("could not fetch update schedule from database",
							zap.Error(err))
						return err
					}

					type contacts struct {
						malak.UpdateRecipient
						Contact       *malak.Contact `json:"contact" bun:"rel:has-one,join:contact_id=id"`
						bun.BaseModel `bun:"table:update_recipients"`
					}

					var contactsFromDB = make([]contacts, 0)

					err = db.NewSelect().Model(&contactsFromDB).
						Limit(10).
						Where("update_id = ?", schedule.UpdateID).
						Where("status = ?", malak.RecipientStatusPending).
						Relation("Contact").
						Scan(ctx)
					if err != nil {
						span.RecordError(err)
						logger.Error("could not fetch recipients for update",
							zap.Error(err))
						return err
					}

					if len(contactsFromDB) == 0 {
						// mark as sent
						logger.Info("no recipients for this update")

						scheduledUpdate.Status = malak.UpdateSendScheduleSent
						_, err = db.NewUpdate().Model(scheduledUpdate).Where("id = ? ", scheduledUpdate.ID).
							Exec(ctx)
						if err != nil {
							span.RecordError(err)
							logger.Error("could not update schedule status",
								zap.Error(err))
							return err
						}

						return nil
					}

					title := fmt.Sprintf("[TEST] %s", update.Title)
					if scheduledUpdate.UpdateType == malak.UpdateTypeLive {
						title = update.Title
					}

					// TODO(adelowo): future version should include batching
					// We cannot batch right now because Resend provider does not send tags when you send a batched email
					// Without tags, there is no way to properly update the right recipient stat
					//
					// API calls will most likely be throttled and we will have to deal with lot of failures
					// But we will figure that out
					//
					// This cron needs a lot of clean up anyways especially with synchronizations amongst others

					wg.Add(len(contactsFromDB))

					for _, contact := range contactsFromDB {
						go func() {
							defer wg.Done()

							time.Sleep(time.Second)

							opts := email.SendOptions{
								HTML:      b.String(),
								Sender:    cfg.Email.Sender,
								Recipient: contact.Contact.Email,
								Subject:   title,
								DKIM: struct {
									Sign       bool
									PrivateKey []byte
								}{
									Sign:       false,
									PrivateKey: []byte(""),
								},
							}

							var status = malak.RecipientStatusSent

							emailID, err := emailClient.Send(ctx, opts)
							if err != nil {
								logger.Error("could not send email", zap.Error(err),
									zap.String("recipient_reference", contact.Reference.String()))

								status = malak.RecipientStatusFailed
								return
							}

							err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

								_, err := tx.NewUpdate().
									Model(&malak.UpdateRecipient{}).
									Where("reference = ?", contact.Reference.String()).
									Set("status = ?", status).
									Exec(ctx)
								if err != nil {
									return err
								}

								emailLog := &malak.UpdateRecipientLog{
									ProviderID:  emailID,
									RecipientID: contact.ID,
									Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeRecipientLog),
									Provider:    emailClient.Name(),
									ID:          uuid.New(),
								}

								_, err = tx.NewInsert().
									Model(emailLog).
									Exec(ctx)
								if err != nil {
									return err
								}

								stats := &malak.UpdateRecipientStat{
									Reference:   malak.NewReferenceGenerator().Generate(malak.EntityTypeRecipientStat),
									RecipientID: contact.ID,
									HasReaction: false,
									IsDelivered: status == malak.RecipientStatusSent,
								}

								_, err = tx.NewInsert().
									Model(stats).
									Exec(ctx)
								return err
							})

							if err != nil {
								logger.Error("could not update database", zap.Error(err),
									zap.String("email_id", emailID),
									zap.String("recipient_reference", contact.Reference.String()))

								status = malak.RecipientStatusFailed
								return
							}
						}()
					}

					wg.Wait()

					scheduledUpdate.Status = malak.UpdateSendScheduleSent
					err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
						_, err = tx.NewUpdate().Model(scheduledUpdate).
							Where("id = ? ", scheduledUpdate.ID).
							Exec(ctx)
						if err != nil {
							return err
						}

						if update.Status == malak.UpdateStatusSent {
							return nil
						}

						update.Status = malak.UpdateStatusSent
						update.SentBy = schedule.ScheduledBy

						_, err = tx.NewUpdate().
							Model(update).
							Where("id = ?", update.ID).
							Exec(ctx)
						return err
					})

					if err != nil {
						span.RecordError(err)
						logger.Error("could not update schedule status",
							zap.Error(err))
						return err
					}

					return nil
				})
			}

			if err := g.Wait(); err != nil {
				logger.Error("could not send updates", zap.Error(err))
				return err
			}

			return nil
		},
	}

	cmd.Flags().Int64("n", 10, "number of scheduled updates to process")

	return cmd
}
