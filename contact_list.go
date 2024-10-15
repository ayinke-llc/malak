package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ContactList struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	Reference   Reference `json:"reference,omitempty"`

	CreatedBy uuid.UUID `json:"created_by,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type FetchContactListOptions struct {
	Reference   Reference
	WorkspaceID uuid.UUID
}

type ContactListRepository interface {
	Create(context.Context, *ContactList) error
	// Get(context.Context, FetchContactListOptions) (ContactList, error)
	// Delete(context.Context, *ContactList) error
	Add(context.Context, ...*Contact) error
	List(context.Context, uuid.UUID) ([]ContactList, error)
}
