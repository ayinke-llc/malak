package email

import (
	"context"
	"errors"
	"io"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
)

type SendOptions struct {
	HTML string

	Sender    malak.Email
	Recipient malak.Email

	DKIM struct {
		Sign       bool
		PrivateKey []byte
	}
}

func (s SendOptions) Validate() error {

	if util.IsStringEmpty(s.HTML) {
		return errors.New("html copy of email must be provided")
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
	Send(context.Context, SendOptions) error
}
