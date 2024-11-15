package main

import (
	"fmt"
	"os"

	"github.com/ayinke-llc/malak/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func addCronCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use: "cron",
	}

	cmd.AddCommand(sendScheduledUpdates(c, cfg))

	c.AddCommand(cmd)
}

func sendScheduledUpdates(c *cobra.Command, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
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

			h, _ := os.Hostname()

			logger = logger.With(zap.String("host", h),
				zap.String("app", "malak"),
				zap.String("component", "cron.updates-send"))

			// Doing this here as I do not see another place where we
			// need to reuse this bit of code
			// if for some reason, we have to send updates from
			// another place.
			// Just encapsulate this so we don't duplicate code

			// ctx, span := tracer.Start(context.Background(), "queue.sendPreviewEmail")
			// defer span.End()
			//
			// span.SetAttributes(
			// 	attribute.Bool("preview", true),
			// 	attribute.String("update_id", p.UpdateID.String()),
			// 	attribute.String("schedule_id", p.ScheduleID.String()),
			// )
			//
			// update, err := t.updateRepo.Get(ctx, malak.FetchUpdateOptions{
			// 	ID: p.UpdateID,
			// })
			// if err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not fetch update from database",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// schedule, err := t.updateRepo.GetSchedule(ctx, p.ScheduleID)
			// if err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not fetch update schedule from database",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// contact, err := t.contactRepo.Get(ctx, malak.FetchContactOptions{
			// 	Email:       p.Email,
			// 	WorkspaceID: update.WorkspaceID,
			// })
			// if err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not fetch contact from database",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// span.SetAttributes(
			// 	attribute.String("triggered_user_id", schedule.ScheduledBy.String()))
			//
			// templatedFile, err := template.New("template").
			// 	Parse(email.UpdateHTMLEmailTemplate)
			// if err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not create html template",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// var b = new(bytes.Buffer)
			// err = templatedFile.Execute(b, map[string]string{
			// 	"Content": update.Content.HTML(),
			// })
			//
			// if err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not parse html template",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// sendOptions := email.SendOptions{
			// 	HTML:      b.String(),
			// 	Sender:    t.cfg.Email.Sender,
			// 	Recipient: contact.Email,
			// 	Subject:   fmt.Sprintf("[TEST] %s", update.Title),
			// 	DKIM: struct {
			// 		Sign       bool
			// 		PrivateKey []byte
			// 	}{
			// 		Sign:       false,
			// 		PrivateKey: []byte(""),
			// 	},
			// }
			//
			// if err := t.emailClient.Send(ctx, sendOptions); err != nil {
			// 	span.RecordError(err)
			// 	logger.Error("could not send preview email",
			// 		zap.Error(err))
			// 	return err
			// }
			//
			// msg.Ack()
			// return nil

			return nil
		},
	}
}
