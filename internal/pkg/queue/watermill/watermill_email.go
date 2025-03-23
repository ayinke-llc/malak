package watermillqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"go.uber.org/zap"
)

func (t *WatermillClient) sendBillingTrialEmail(msg *message.Message) error {

	ctx, span := tracer.Start(context.Background(),
		"sendBillingTrialEmail")

	defer span.End()

	var opts queue.SendBillingTrialEmailOptions

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).
		Decode(&opts); err != nil {
		return err
	}

	logger := t.logger.With(zap.String("method", "sendBillingTrialEmail"),
		zap.String("workspace_id", opts.Workspace.ID.String()))

	logger.Debug("sending email to user for free trial")

	tmpl, err := template.New("template").Parse(email.BillingTrialTemplate)
	if err != nil {
		logger.Error("could not parse email template", zap.Error(err))
		return err
	}

	var link = t.cfg.Frontend.AppURL + "/settings?tab=billing"

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"WorkspaceName": opts.Workspace.WorkspaceName,
		"Link":          link,
		"Expiration":    opts.Expiration,
	}); err != nil {
		logger.Error("could not embed content in template", zap.Error(err))
		return err
	}

	emailOpts := email.SendOptions{
		HTML:      buf.String(),
		Sender:    t.cfg.Email.Sender,
		Recipient: opts.Recipient,
		Subject:   "Your Malak trial is coming to an end",
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign:       false,
			PrivateKey: []byte(""),
		},
	}

	_, err = t.emailClient.Send(ctx, emailOpts)
	if err != nil {
		logger.Error("could not send email", zap.Error(err))
		return err
	}

	return nil
}

func (t *WatermillClient) sendDashboardSharingEmail(msg *message.Message) error {

	ctx, span := tracer.Start(context.Background(),
		"sendDashboardSharingEmail")

	defer span.End()

	var opts queue.SendEmailOptions

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).
		Decode(&opts); err != nil {
		return err
	}

	logger := t.logger.With(zap.String("method", "sendDashboardSharingEmail"),
		zap.String("workspace_id", opts.Workspace.ID.String()))

	logger.Debug("sending email to user")

	tmpl, err := template.New("template").Parse(email.DashboardSharingTemplate)
	if err != nil {
		logger.Error("could not parse email template", zap.Error(err))
		return err
	}

	var link = t.cfg.Frontend.AppURL + "/shared/dashboards/" + opts.Token

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"WorkspaceName": opts.Workspace.WorkspaceName,
		"Link":          link,
	}); err != nil {
		logger.Error("could not embed content in template", zap.Error(err))
		return err
	}

	emailOpts := email.SendOptions{
		HTML:      buf.String(),
		Sender:    t.cfg.Email.Sender,
		Recipient: opts.Recipient,
		Subject:   "Metrics dashboard shared with you by " + opts.Workspace.WorkspaceName,
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign:       false,
			PrivateKey: []byte(""),
		},
	}

	_, err = t.emailClient.Send(ctx, emailOpts)
	if err != nil {
		logger.Error("could not send email", zap.Error(err))
		return err
	}

	return nil
}
