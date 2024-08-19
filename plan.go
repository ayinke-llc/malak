package malak

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPlanNotFound     = errors.New("plan does not exists")
	ErrCounterExhausted = errors.New("no more units left")
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
}

type Plan struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`

	PlanName string `json:"plan_name,omitempty"`

	// Can use a fake id really
	StripeReference string       `json:"stripe_reference,omitempty"`
	Metadata        PlanMetadata `json:"metadata,omitempty" bson:"metadata"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`
}

type FetchPlanOptions struct {
	Reference string
	ID        uuid.UUID
}

type PlanRepository interface {
	Get(context.Context, *FetchPlanOptions) (*Plan, error)
	List(context.Context) ([]*Plan, error)
}
