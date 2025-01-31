package billing

import (
	"context"

	"github.com/ayinke-llc/malak"
)

type CreateCustomerOptions struct {
	Name  string
	Email malak.Email
}

type CreateBillingPortalOptions struct {
	CustomerID string
}

type Client interface {
	CreateCustomer(context.Context, *CreateCustomerOptions) (string, error)
	Portal(context.Context, *CreateBillingPortalOptions) (string, error)
}
