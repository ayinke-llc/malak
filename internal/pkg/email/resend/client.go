package resend

import (
	"context"
	"fmt"

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
	c := resendclient.NewClient("re_123456789")

	return &client{
		inner:       c,
		senderName:  cfg.Email.SenderName,
		senderEmail: cfg.Email.Sender.String(),
	}, nil
}

func (s *client) Close() error { return nil }

func (s *client) Send(ctx context.Context,
	opts email.SendOptions) error {

	// msg := gomail.NewMessage()
	// msg.SetAddressHeader("From", opts.Sender.String(), opts.Sender.String())
	// msg.SetHeader("To", opts.Recipient.String())
	// msg.SetHeader("Subject", opts.Subject)
	// msg.AddAlternative("text/html", opts.HTML)
	//
	// return s.client.DialAndSend(msg)
	return nil
}

func (s *client) SendBatch(ctx context.Context,
	opts []email.SendOptions) error {

	var batchEmails = make([]*resendclient.SendEmailRequest, 0, len(opts))

	for _, v := range opts {
		batchEmails = append(batchEmails, &resendclient.SendEmailRequest{
			From:    fmt.Sprintf("%s %s", s.senderName, s.senderEmail),
			To:      []string{v.Recipient.String()},
			Subject: v.Subject,
			Html:    v.HTML,
		})
	}

	_, err := s.inner.Batch.SendWithContext(ctx, batchEmails)
	return err
}
