package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrContactNotFound = MalakError("contact not found")
	ErrContactExists   = MalakError("contact with email already exists")
)

// ENUM(mr,mrs,miss,doctor,chief)
type ContactTitle string

type CustomContactMetadata map[string]string

type Contact struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Email       Email     `json:"email,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	Reference   Reference `json:"reference,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Company     string    `json:"company,omitempty"`

	// Legacy lmao. should be address but migrations bit ugh :))
	City  string `json:"city,omitempty"`
	Phone string `json:"phone,omitempty"`

	Notes string               `json:"notes,omitempty"`
	Lists []ContactListMapping `json:"lists" bun:"rel:has-many,join:id=contact_id"`

	// User who owns the contact.
	// Does not mean who added the contact but who chases
	// or follows up officially with the contact
	OwnerID uuid.UUID `json:"owner_id,omitempty" bun:",nullzero"`

	// User who added/created this contact
	CreatedBy uuid.UUID `json:"created_by,omitempty" bun:",nullzero"`

	Metadata CustomContactMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type FetchContactOptions struct {
	ID          uuid.UUID
	Email       Email
	Reference   Reference
	WorkspaceID uuid.UUID
}

type ListContactOptions struct {
	Paginator   Paginator
	WorkspaceID uuid.UUID
	Status      ListUpdateFilterStatus
}

type ContactOverview struct {
	TotalContacts int64 `json:"total_contacts,omitempty"`
}

type SearchContactOptions struct {
	WorkspaceID uuid.UUID
	SearchValue string
}

type ContactRepository interface {
	Create(context.Context, ...*Contact) error
	Get(context.Context, FetchContactOptions) (*Contact, error)
	List(context.Context, ListContactOptions) ([]Contact, int64, error)
	Delete(context.Context, *Contact) error
	Update(context.Context, *Contact) error
	Overview(context.Context, uuid.UUID) (*ContactOverview, error)
	Search(context.Context, SearchContactOptions) ([]Contact, error)

	// This should only be used for updates sending
	// ideally moste people have under 50 contacts so it is fine
	// If we see people have 200-1k contacts, then we can optimise this even better
	// EBut we really may never know. OTEL will tell us if it makes
	// any sense at all or not
	All(context.Context, uuid.UUID) ([]Contact, error)
}
