package stripe

import (
	"os"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/billing"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupStripeClient() (billing.Client, error) {
	cfg := config.Config{
		Billing: struct {
			Stripe struct {
				APIKey        string "yaml:\"api_key\" mapstructure:\"api_key\""
				WebhookSecret string "yaml:\"webhook_secret\" mapstructure:\"webhook_secret\""
			} "yaml:\"stripe\" mapstructure:\"stripe\""
			IsEnabled            bool   "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
			TrialDays            int64  "yaml:\"trial_days\" mapstructure:\"trial_days\""
			DefaultPlanReference string "yaml:\"default_plan_reference\" mapstructure:\"default_plan_reference\""
		}{
			Stripe: struct {
				APIKey        string "yaml:\"api_key\" mapstructure:\"api_key\""
				WebhookSecret string "yaml:\"webhook_secret\" mapstructure:\"webhook_secret\""
			}{
				APIKey: os.Getenv("STRIPE_SECRET"),
			},
			TrialDays: 14,
		},
	}

	return New(cfg)
}

func TestCreate_Integration(t *testing.T) {
	stripeClient, err := setupStripeClient()
	require.NoError(t, err)

	opts := &billing.CreateCustomerOptions{
		Email: malak.Email(faker.Email()),
		Name:  faker.Name(),
	}

	customerID, err := stripeClient.CreateCustomer(t.Context(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, customerID)

	billingURL, err := stripeClient.Portal(t.Context(), &billing.CreateBillingPortalOptions{
		CustomerID: customerID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, billingURL)

	subscriptionID, err := stripeClient.AddPlanToCustomer(t.Context(), &billing.AddPlanToCustomerOptions{
		Workspace: &malak.Workspace{
			ID:               uuid.New(),
			StripeCustomerID: customerID,
			Plan: &malak.Plan{
				DefaultPriceID: "price_1PvJSMIuzgc0GUapiNyXBEaH",
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, subscriptionID)
}
