package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrPlanNotFound       = MalakError("plan does not exists")
	ErrCounterExhausted   = MalakError("no more units left")
	ErrOnlyOneDefaultPlan = MalakError("there can only be one default plan")
)

type Counter int64

func (c Counter) Add() { c++ }

func (c Counter) Take() error {
	if c <= 0 {
		return ErrCounterExhausted
	}

	c--
	return nil
}

func (c Counter) TakeN(n int64) error {
	if c <= 0 {
		return c.Take()
	}

	c -= Counter(n)
	return nil
}

type PlanMetadata struct {
	Team struct {
		Size Counter `json:"size,omitempty"`
	} `json:"team,omitempty"`

	Deck struct {
		AutoTerminateLink bool `json:"auto_terminate_link,omitempty"`
		CustomDomain      bool `json:"custom_domain,omitempty"`
	} `json:"deck,omitempty"`

	Updates struct {
		MaxRecipients Counter `json:"max_recipients,omitempty"`
		CustomDomain  bool    `json:"custom_domain,omitempty"`
	} `json:"updates,omitempty"`

	Integrations struct {
		AvailableForUse Counter `json:"available_for_use,omitempty"`
	} `json:"integrations,omitempty"`

	Dashboard struct {
		ShareDashboardViaLink bool `json:"share_dashboard_via_link,omitempty"`
		EmbedDashboard        bool `json:"embed_dashboard,omitempty"`
	} `json:"dashboard,omitempty"`

	DataRoom struct {
		Size         Counter `json:"size,omitempty"`
		ShareViaLink bool    `json:"share_via_link,omitempty"`
	} `json:"data_room,omitempty"`
}

type Plan struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	PlanName string `json:"plan_name,omitempty"`

	// Can use a fake id really
	// As this only matters if you turn on Stripe
	Reference string `json:"reference,omitempty"`

	Metadata PlanMetadata `json:"metadata,omitempty" bson:"metadata"`

	// Stripe default price id. Again not needed if not using Stripe
	DefaultPriceID string `json:"default_price_id,omitempty"`

	// Defaults to zero
	Amount int64 `json:"amount,omitempty"`

	// IsDefault if this is the default plan for the user to get signed up to
	// on sign up
	//
	// Better to keep this here than to use config
	IsDefault bool `json:"is_default,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type FetchPlanOptions struct {
	Reference string
	ID        uuid.UUID
}

type PlanRepository interface {
	Get(context.Context, *FetchPlanOptions) (*Plan, error)
	List(context.Context) ([]*Plan, error)
	SetDefault(context.Context, *Plan) error
	Create(context.Context, *Plan) error
}
