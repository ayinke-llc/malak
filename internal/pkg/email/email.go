package email

import (
	"context"
	_ "embed"
	"errors"
	"io"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
)

// TODO: migrate these to embed.Fs instead
var (
	//go:embed templates/updates/view.html
	UpdateHTMLEmailTemplate string

	//go:embed templates/sharing/dashboard_share.html
	DashboardSharingTemplate string

	//go:embed templates/billing/trial.html
	BillingTrialTemplate string

	//go:embed templates/billing/expired.html
	BillingEndedTemplate string

	//go:embed templates/auth/email_verify.html
	EmailVerificationTemplate string
)

type SendOptionsBatch []SendOptions

func (s SendOptionsBatch) Validate() error {
	if len(s) > 25 {
		return errors.New("maximum of 25 allowed per batch")
	}

	for _, v := range s {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type SendOptions struct {
	HTML      string
	Sender    malak.Email
	Recipient malak.Email
	Subject   string
	DKIM      struct {
		Sign       bool
		PrivateKey []byte
	}
}

func (s SendOptions) Validate() error {

	if util.IsStringEmpty(s.HTML) {
		return errors.New("html copy of email must be provided")
	}

	if util.IsStringEmpty(s.Subject) {
		return errors.New("please provide subject")
	}

	if util.IsStringEmpty(s.Recipient.String()) {
		return errors.New("please provide recipient")
	}

	if util.IsStringEmpty(s.Sender.String()) {
		return errors.New("please provide sender")
	}

	return nil
}

type Client interface {
	io.Closer
	Send(context.Context, SendOptions) (string, error)
	Name() malak.UpdateRecipientLogProvider
}
