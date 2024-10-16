package smtp

import (
	"context"
	"errors"
	pkgmail "net/mail"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/wneessen/go-mail"
)

type smtpClient struct {
	client *mail.Client
}

func New(cfg config.Config) (email.Client, error) {

	if util.IsStringEmpty(cfg.Email.SMTP.Host) {
		return nil, errors.New("please provide your smtp host")
	}

	if util.IsStringEmpty(cfg.Email.SMTP.Username) {
		return nil, errors.New("please provide your smtp username")
	}

	if util.IsStringEmpty(cfg.Email.SMTP.Password) {
		return nil, errors.New("please provide your smtp password")
	}

	if util.IsStringEmpty(cfg.Email.Sender.String()) {
		return nil, errors.New("please provide a valid sender")
	}

	_, err := pkgmail.ParseAddress(string(cfg.Email.Sender))
	if err != nil {
		return nil, errors.Join(err, errors.New("invalid email sender"))
	}

	client, err := mail.NewClient(cfg.Email.SMTP.Host,
		mail.WithPort(cfg.Email.SMTP.Port),
		mail.WithUsername(cfg.Email.SMTP.Username),
		mail.WithPassword(cfg.Email.SMTP.Password))
	if err != nil {
		return nil, err
	}

	client.SetTLSPolicy(mail.NoTLS)

	if cfg.Email.SMTP.UseTLS {
		client.SetTLSPolicy(mail.TLSMandatory)
	}

	if err := client.DialWithContext(context.Background()); err != nil {
		return nil, err
	}

	return &smtpClient{
		client: client,
	}, nil
}

func (s *smtpClient) Close() error { return s.client.Close() }

func (s *smtpClient) Send(ctx context.Context,
	opts email.SendOptions) error {

	msg := mail.NewMsg()

	if err := msg.From(opts.Sender.String()); err != nil {
		return err
	}

	if err := msg.To(opts.Recipient.String()); err != nil {
		return err
	}

	msg.Subject(opts.Subject)

	msg.SetBodyString(mail.TypeTextHTML, opts.HTML)

	return s.client.Send(msg)
}
