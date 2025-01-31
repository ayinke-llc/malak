package server

type EventStripeWebhook struct {
	AccountCountry       string `json:"account_country"`
	AccountName          string `json:"account_name"`
	AccountTaxIds        any    `json:"account_tax_ids"`
	AmountDue            int    `json:"amount_due"`
	AmountPaid           int    `json:"amount_paid"`
	AmountRemaining      int    `json:"amount_remaining"`
	AmountShipping       int    `json:"amount_shipping"`
	Application          any    `json:"application"`
	ApplicationFeeAmount any    `json:"application_fee_amount"`
	AttemptCount         int    `json:"attempt_count"`
	Attempted            bool   `json:"attempted"`
	AutoAdvance          bool   `json:"auto_advance"`
	AutomaticTax         struct {
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	} `json:"automatic_tax"`
	BillingReason    string `json:"billing_reason"`
	Charge           string `json:"charge"`
	CollectionMethod string `json:"collection_method"`
	Created          int    `json:"created"`
	Currency         string `json:"currency"`
	CustomFields     any    `json:"custom_fields"`
	Customer         string `json:"customer"`
	CustomerAddress  struct {
		City       any    `json:"city"`
		Country    string `json:"country"`
		Line1      any    `json:"line1"`
		Line2      any    `json:"line2"`
		PostalCode any    `json:"postal_code"`
		State      any    `json:"state"`
	} `json:"customer_address"`
	CustomerEmail         string `json:"customer_email"`
	CustomerName          string `json:"customer_name"`
	CustomerPhone         string `json:"customer_phone"`
	CustomerShipping      any    `json:"customer_shipping"`
	CustomerTaxExempt     string `json:"customer_tax_exempt"`
	CustomerTaxIds        []any  `json:"customer_tax_ids"`
	DefaultPaymentMethod  any    `json:"default_payment_method"`
	DefaultSource         any    `json:"default_source"`
	DefaultTaxRates       []any  `json:"default_tax_rates"`
	Description           any    `json:"description"`
	Discount              any    `json:"discount"`
	Discounts             []any  `json:"discounts"`
	DueDate               any    `json:"due_date"`
	EndingBalance         int    `json:"ending_balance"`
	Footer                any    `json:"footer"`
	FromInvoice           any    `json:"from_invoice"`
	HostedInvoiceURL      string `json:"hosted_invoice_url"`
	ID                    string `json:"id"`
	InvoicePdf            string `json:"invoice_pdf"`
	LastFinalizationError any    `json:"last_finalization_error"`
	LatestRevision        any    `json:"latest_revision"`
	Lines                 struct {
		Data []struct {
			Amount             int      `json:"amount"`
			AmountExcludingTax int      `json:"amount_excluding_tax"`
			Currency           string   `json:"currency"`
			Description        string   `json:"description"`
			DiscountAmounts    []any    `json:"discount_amounts"`
			Discountable       bool     `json:"discountable"`
			Discounts          []any    `json:"discounts"`
			ID                 string   `json:"id"`
			Livemode           bool     `json:"livemode"`
			Metadata           struct{} `json:"metadata"`
			Object             string   `json:"object"`
			Period             struct {
				End   int `json:"end"`
				Start int `json:"start"`
			} `json:"period"`
			Plan struct {
				Active          bool     `json:"active"`
				AggregateUsage  any      `json:"aggregate_usage"`
				Amount          int      `json:"amount"`
				AmountDecimal   string   `json:"amount_decimal"`
				BillingScheme   string   `json:"billing_scheme"`
				Created         int      `json:"created"`
				Currency        string   `json:"currency"`
				ID              string   `json:"id"`
				Interval        string   `json:"interval"`
				IntervalCount   int      `json:"interval_count"`
				Livemode        bool     `json:"livemode"`
				Metadata        struct{} `json:"metadata"`
				Nickname        any      `json:"nickname"`
				Object          string   `json:"object"`
				Product         string   `json:"product"`
				TiersMode       any      `json:"tiers_mode"`
				TransformUsage  any      `json:"transform_usage"`
				TrialPeriodDays any      `json:"trial_period_days"`
				UsageType       string   `json:"usage_type"`
			} `json:"plan"`
			Price struct {
				Active           bool     `json:"active"`
				BillingScheme    string   `json:"billing_scheme"`
				Created          int      `json:"created"`
				Currency         string   `json:"currency"`
				CustomUnitAmount any      `json:"custom_unit_amount"`
				ID               string   `json:"id"`
				Livemode         bool     `json:"livemode"`
				LookupKey        any      `json:"lookup_key"`
				Metadata         struct{} `json:"metadata"`
				Nickname         any      `json:"nickname"`
				Object           string   `json:"object"`
				Product          string   `json:"product"`
				Recurring        struct {
					AggregateUsage  any    `json:"aggregate_usage"`
					Interval        string `json:"interval"`
					IntervalCount   int    `json:"interval_count"`
					TrialPeriodDays any    `json:"trial_period_days"`
					UsageType       string `json:"usage_type"`
				} `json:"recurring"`
				TaxBehavior       string `json:"tax_behavior"`
				TiersMode         any    `json:"tiers_mode"`
				TransformQuantity any    `json:"transform_quantity"`
				Type              string `json:"type"`
				UnitAmount        int    `json:"unit_amount"`
				UnitAmountDecimal string `json:"unit_amount_decimal"`
			} `json:"price"`
			Proration        bool `json:"proration"`
			ProrationDetails struct {
				CreditedItems any `json:"credited_items"`
			} `json:"proration_details"`
			Quantity               int    `json:"quantity"`
			Subscription           string `json:"subscription"`
			SubscriptionItem       string `json:"subscription_item"`
			TaxAmounts             []any  `json:"tax_amounts"`
			TaxRates               []any  `json:"tax_rates"`
			Type                   string `json:"type"`
			UnitAmountExcludingTax string `json:"unit_amount_excluding_tax"`
		} `json:"data"`
		HasMore    bool   `json:"has_more"`
		Object     string `json:"object"`
		TotalCount int    `json:"total_count"`
		URL        string `json:"url"`
	} `json:"lines"`
	Livemode           bool     `json:"livemode"`
	Metadata           struct{} `json:"metadata"`
	NextPaymentAttempt any      `json:"next_payment_attempt"`
	Number             string   `json:"number"`
	Object             string   `json:"object"`
	OnBehalfOf         any      `json:"on_behalf_of"`
	Paid               bool     `json:"paid"`
	PaidOutOfBand      bool     `json:"paid_out_of_band"`
	PaymentIntent      string   `json:"payment_intent"`
	PaymentSettings    struct {
		DefaultMandate       any `json:"default_mandate"`
		PaymentMethodOptions any `json:"payment_method_options"`
		PaymentMethodTypes   any `json:"payment_method_types"`
	} `json:"payment_settings"`
	PeriodEnd                    int    `json:"period_end"`
	PeriodStart                  int    `json:"period_start"`
	PostPaymentCreditNotesAmount int    `json:"post_payment_credit_notes_amount"`
	PrePaymentCreditNotesAmount  int    `json:"pre_payment_credit_notes_amount"`
	Quote                        any    `json:"quote"`
	ReceiptNumber                any    `json:"receipt_number"`
	RenderingOptions             any    `json:"rendering_options"`
	ShippingCost                 any    `json:"shipping_cost"`
	ShippingDetails              any    `json:"shipping_details"`
	StartingBalance              int    `json:"starting_balance"`
	StatementDescriptor          any    `json:"statement_descriptor"`
	Status                       string `json:"status"`
	StatusTransitions            struct {
		FinalizedAt           int `json:"finalized_at"`
		MarkedUncollectibleAt any `json:"marked_uncollectible_at"`
		PaidAt                int `json:"paid_at"`
		VoidedAt              any `json:"voided_at"`
	} `json:"status_transitions"`
	Subscription         string `json:"subscription"`
	Subtotal             int    `json:"subtotal"`
	SubtotalExcludingTax int    `json:"subtotal_excluding_tax"`
	Tax                  any    `json:"tax"`
	TestClock            any    `json:"test_clock"`
	Total                int    `json:"total"`
	TotalDiscountAmounts []any  `json:"total_discount_amounts"`
	TotalExcludingTax    int    `json:"total_excluding_tax"`
	TotalTaxAmounts      []any  `json:"total_tax_amounts"`
	TransferData         any    `json:"transfer_data"`
	WebhooksDeliveredAt  any    `json:"webhooks_delivered_at"`
}

type StripeWebhook struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	APIVersion string `json:"api_version"`
	Created    int    `json:"created"`
	Data       struct {
		Object struct {
			Subscription     string `json:"subscription"`
			ID               string `json:"id"`
			Customer         string `json:"customer"`
			CurrentPeriodEnd int64  `json:"current_period_end"`
			Lines            struct {
				Data []struct {
					Plan struct {
						ID       string `json:"id"`
						Object   string `json:"object"`
						Active   bool   `json:"active"`
						Livemode bool   `json:"livemode"`
						Product  string `json:"product"`
					} `json:"plan"`
				}
			} `json:"lines"`
			Quantity  int    `json:"quantity"`
			Schedule  any    `json:"schedule"`
			StartDate int    `json:"start_date"`
			Status    string `json:"status"`
		} `json:"object"`
	} `json:"data"`
	Livemode bool `json:"livemode"`
	Request  struct {
		ID             string `json:"id"`
		IdempotencyKey string `json:"idempotency_key"`
	} `json:"request"`
	Type string `json:"type"`
}

type SubscriptionDeletedRequest struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Customer string `json:"customer"`
}

type TrialWillEnd struct {
	ID   string `json:"id"`
	Plan struct {
		ID              string   `json:"id"`
		Object          string   `json:"object"`
		Active          bool     `json:"active"`
		AggregateUsage  any      `json:"aggregate_usage"`
		Amount          int      `json:"amount"`
		AmountDecimal   string   `json:"amount_decimal"`
		BillingScheme   string   `json:"billing_scheme"`
		Created         int      `json:"created"`
		Currency        string   `json:"currency"`
		Interval        string   `json:"interval"`
		IntervalCount   int      `json:"interval_count"`
		Livemode        bool     `json:"livemode"`
		Metadata        struct{} `json:"metadata"`
		Nickname        any      `json:"nickname"`
		Product         string   `json:"product"`
		TiersMode       any      `json:"tiers_mode"`
		TransformUsage  any      `json:"transform_usage"`
		TrialPeriodDays any      `json:"trial_period_days"`
		UsageType       string   `json:"usage_type"`
	} `json:"plan"`
	Status     string `json:"status"`
	TrialEnd   int64  `json:"trial_end"`
	TrialStart int    `json:"trial_start"`
	Customer   string `json:"customer"`
}

type CustomerCreatedEvent struct {
	Data struct {
		Object struct {
			ID string `json:"id"`
		} `json:"object"`
	} `json:"data"`
	ID         string `json:"id"`
	Object     string `json:"object"`
	APIVersion string `json:"api_version"`
	Created    int    `json:"created"`
	Type       string `json:"type"`
}

type CustomerUpdatedSubscription struct {
	ID       string `json:"id"`
	Customer string `json:"customer"`
	Items    struct {
		Object string `json:"object"`
		Data   []struct {
			Plan struct {
				ID            string `json:"id"`
				Active        bool   `json:"active"`
				Amount        int    `json:"amount"`
				BillingScheme string `json:"billing_scheme"`
				Interval      string `json:"interval"`
				IntervalCount int    `json:"interval_count"`
				Product       string `json:"product"`
				UsageType     string `json:"usage_type"`
			} `json:"plan"`
			Price struct {
				ID            string `json:"id"`
				Active        bool   `json:"active"`
				BillingScheme string `json:"billing_scheme"`
				Product       string `json:"product"`
			} `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

type Invoice struct {
	ID               string `json:"id"`
	AmountDue        int    `json:"amount_due"`
	Created          int    `json:"created"`
	Customer         string `json:"customer"`
	HostedInvoiceURL string `json:"hosted_invoice_url"`
	InvoicePdf       string `json:"invoice_pdf"`
	Lines            struct {
		Object string `json:"object"`
		Data   []struct {
			Plan struct {
				Product string `json:"product"`
			} `json:"plan"`
		} `json:"data"`
	} `json:"lines"`
	Paid         bool   `json:"paid"`
	Status       string `json:"status"`
	Subscription string `json:"subscription"`
}
