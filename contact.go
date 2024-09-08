package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrContactNotFound = malakError("contact not found")
	ErrContactExists   = malakError("contact with email already exists")
)

// ENUM(mr,mrs,miss,doctor,chief)
type ContactTitle string

type CustomContactMetadata map[string]string

type Contact struct {
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Email     Email     `json:"email,omitempty"`
	Reference Reference `json:"reference,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Company   string    `json:"company,omitempty"`
	City      string    `json:"city,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Notes     string    `json:"notes,omitempty"`

	// User who owns the contact.
	// Does not mean who added the contact but who chases
	// or follows up officially with the contact
	OwnerID uuid.UUID `json:"owner_id,omitempty"`

	Metadata CustomContactMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `bun:"table:users" json:"-"`
}

type FetchContactOption struct {
	ID        uuid.UUID
	Email     Email
	Reference Reference
}

type ContactRepository interface {
	Create(context.Context, *Contact) error
	Get(context.Context)
}
