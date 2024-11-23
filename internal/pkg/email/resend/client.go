package resend

import (
	"context"
	"fmt"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	resendclient "github.com/resend/resend-go/v2"
)

type client struct {
	inner       *resendclient.Client
	senderName  string
	senderEmail string
}

func New(cfg config.Config) (email.Client, error) {
	c := resendclient.NewClient(cfg.Email.Resend.APIKey)

	return &client{
		inner:       c,
		senderName:  cfg.Email.SenderName,
		senderEmail: cfg.Email.Sender.String(),
	}, nil
}

func (s *client) Close() error { return nil }

func (s *client) Send(ctx context.Context,
	opts email.SendOptions) (string, error) {

	params := &resendclient.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.senderName, s.senderEmail),
		To:      []string{opts.Recipient.String()},
		Subject: opts.Subject,
		Html:    opts.HTML,
	}

	res, err := s.inner.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (s *client) Name() malak.UpdateRecipientLogProvider {
	return malak.UpdateRecipientLogProviderResend
}
