package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(update,dashboard,deck)
type ContactShareItemType string

type ContactShare struct {
	ID        uuid.UUID            `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference Reference            `json:"reference,omitempty"`
	SharedBy  uuid.UUID            `json:"shared_by,omitempty"`
	ContactID uuid.UUID            `json:"contact_id,omitempty"`
	ItemType  ContactShareItemType `json:"item_type,omitempty"`
	ItemID    uuid.UUID            `json:"item_id,omitempty"`
	SharedAt  time.Time            `json:"shared_at,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type ContactShareItem struct {
	ContactShare
	Title string `json:"title,omitempty"`
}

type ContactShareRepository interface {
}
