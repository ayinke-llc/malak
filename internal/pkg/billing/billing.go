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

type AddPlanToCustomerOptions struct {
	Workspace *malak.Workspace
}

type Client interface {
	CreateCustomer(context.Context, *CreateCustomerOptions) (string, error)
	Portal(context.Context, *CreateBillingPortalOptions) (string, error)
	AddPlanToCustomer(context.Context, *AddPlanToCustomerOptions) (string, error)
}
