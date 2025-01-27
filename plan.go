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
		Size    Counter `json:"size,omitempty"`
		Enabled bool    `json:"enabled,omitempty"`
	} `json:"team,omitempty"`
	Deck struct {
		Count Counter `json:"count,omitempty"`
	} `json:"deck,omitempty"`
}

type Plan struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`

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
