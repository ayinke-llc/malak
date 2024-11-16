package smtp

import (
	"context"
	"errors"
	pkgmail "net/mail"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"gopkg.in/gomail.v2"
)

type smtpClient struct {
	client *gomail.Dialer
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

	if util.IsStringEmpty(cfg.Email.SenderName) {
		return nil, errors.New("please provide your sender name")
	}

	_, err := pkgmail.ParseAddress(string(cfg.Email.Sender))
	if err != nil {
		return nil, errors.Join(err, errors.New("invalid email sender"))
	}

	client := gomail.NewDialer(cfg.Email.SMTP.Host, cfg.Email.SMTP.Port, cfg.Email.SMTP.Username, cfg.Email.SMTP.Password)

	return &smtpClient{
		client: client,
	}, nil
}

func (s *smtpClient) Close() error { return nil }

func (s *smtpClient) Send(ctx context.Context,
	opts email.SendOptions) error {

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", opts.Sender.String(), opts.Sender.String())
	msg.SetHeader("To", opts.Recipient.String())
	msg.SetHeader("Subject", opts.Subject)
	msg.AddAlternative("text/html", opts.HTML)

	return s.client.DialAndSend(msg)
}
