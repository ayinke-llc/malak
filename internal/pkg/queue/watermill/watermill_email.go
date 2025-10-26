package watermillqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
)

func (t *WatermillClient) sendSubExpiredEmail(msg *message.Message) error {

	ctx, span := tracer.Start(context.Background(),
		"sendSubExpiredEmail")

	defer span.End()

	var opts queue.SubscriptionExpiredOptions

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).
		Decode(&opts); err != nil {
		return err
	}

	logger := t.logger.With(zap.String("method", "sendSubExpiredEmail"),
		zap.String("workspace_id", opts.Workspace.ID.String()))

	logger.Debug("sending sub expired email")

	tmpl, err := template.New("template").Parse(email.BillingEndedTemplate)
	if err != nil {
		logger.Error("could not parse email template", zap.Error(err))
		return err
	}

	var link = t.cfg.Frontend.AppURL + "/settings?tab=billing"

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
		Subject:   "Your Malak subscription has come to an end. Please resubscribe",
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

func (t *WatermillClient) sendEmailVerification(msg *message.Message) error {

	ctx, span := tracer.Start(context.Background(),
		"sendEmailVerification")

	defer span.End()

	var opts queue.EmailVerificationOptions

	if err := json.NewDecoder(bytes.NewBuffer(msg.Payload)).
		Decode(&opts); err != nil {
		return err
	}

	logger := t.logger.With(zap.String("method", "sendEmailVerification"))

	logger.Debug("sending email to user")

	tmpl, err := template.New("template").Parse(email.EmailVerificationTemplate)
	if err != nil {
		logger.Error("could not parse email template", zap.Error(err))
		return err
	}

	user, err := t.userRepo.Get(ctx, &malak.FindUserOptions{
		ID: opts.UserID,
	})
	if err != nil {
		logger.Error("could not fetch user from database", zap.Error(err))
		return err
	}

	var link = t.cfg.Frontend.AppURL + "/email-verify?token=" + opts.Token

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"FullName": user.FullName,
		"Link":     link,
	}); err != nil {
		logger.Error("could not embed content in template", zap.Error(err))
		return err
	}

	emailOpts := email.SendOptions{
		HTML:      buf.String(),
		Sender:    t.cfg.Email.Sender,
		Recipient: user.Email,
		Subject:   "Verify your account to get started with Malak",
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
