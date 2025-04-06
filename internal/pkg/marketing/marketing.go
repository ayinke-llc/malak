package marketing

import (
	"context"

	"github.com/ayinke-llc/malak"
)

type CreateContactOptions struct {
	FirstName string
	LastName  string
	Email     malak.Email
}

type MarketingClient interface {
	CreateContact(context.Context, *CreateContactOptions) (string, error)
}
