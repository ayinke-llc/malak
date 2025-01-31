package stripe

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/billing"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var tracer = otel.Tracer("billing.stripe")

type stripeClient struct {
	stripeClient *client.API
}

func New(cfg config.Config) (billing.Client, error) {

	stripeLib := &client.API{}

	config := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(0), // Zero retries
		EnableTelemetry:   hermes.Ref(false),
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}

	stripeLib.Init(cfg.Billing.Stripe.APIKey, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, config),
	})

	return &stripeClient{
		stripeClient: stripeLib,
	}, nil
}

func (s *stripeClient) CreateCustomer(ctx context.Context,
	opts *billing.CreateCustomerOptions) (string, error) {

	ctx, span := tracer.Start(ctx, "customer.create")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(http.Header{}))

	customer, err := s.stripeClient.Customers.New(&stripe.CustomerParams{
		Email: hermes.Ref(opts.Email.String()),
		Name:  hermes.Ref(opts.Name),
	})
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func (s *stripeClient) Portal(ctx context.Context,
	opts *billing.CreateBillingPortalOptions) (string, error) {

	ctx, span := tracer.Start(ctx, "billing.portal")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(http.Header{}))

	portal, err := s.stripeClient.BillingPortalSessions.New(&stripe.BillingPortalSessionParams{
		Customer: hermes.Ref(opts.CustomerID),
	})
	if err != nil {
		return "", err
	}

	return portal.URL, nil
}
