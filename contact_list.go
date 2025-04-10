package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrContactListNotFound = MalakError("contact list not found")
)

type ContactList struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	Reference   Reference `json:"reference,omitempty"`
	// the number of items in the list. This is postgresql triggered
	// not manually updated
	// NumberInList int `json:"number_in_list,omitempty"`

	CreatedBy uuid.UUID `json:"created_by,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type ContactListMapping struct {
	ID        uuid.UUID    `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	ContactID uuid.UUID    `json:"contact_id,omitempty"`
	ListID    uuid.UUID    `json:"list_id,omitempty"`
	List      *ContactList `json:"list,omitempty" bun:"rel:has-one,join:list_id=id"`
	Reference Reference    `json:"reference,omitempty"`
	CreatedBy uuid.UUID    `json:"created_by,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type FetchContactListOptions struct {
	Reference   Reference
	WorkspaceID uuid.UUID
}

type ContactListOptions struct {
	WorkspaceID   uuid.UUID
	IncludeEmails bool
}

type ContactListRepository interface {
	Create(context.Context, *ContactList) error
	Get(context.Context, FetchContactListOptions) (*ContactList, error)
	Delete(context.Context, *ContactList) error
	Update(context.Context, *ContactList) error
	Add(context.Context, *ContactListMapping) error
	List(context.Context, *ContactListOptions) ([]ContactList, []ContactListMappingWithContact, error)
}

type ContactListMappingWithContact struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ListID    uuid.UUID `json:"list_id,omitempty"`
	ContactID uuid.UUID `json:"contact_id,omitempty"`
	Reference string    `json:"reference,omitempty"`

	// Contact fields
	Email string `json:"email,omitempty"`
}
