package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(draft,sent)
type UpdateStatus string

type Update struct {
	ID          uuid.UUID    `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID    `json:"workspace_id,omitempty"`
	Status      UpdateStatus `json:"status,omitempty"`
	Reference   Reference    `json:"reference,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type FetchUpdateOptions struct {
	Status    UpdateStatus
	Reference Reference
	ID        uuid.UUID
}

type UpdateRepository interface {
	// Create(context.Context, *Update) error
	// Update(context.Context, *Update) error
	// Get(context.Context, FetchUpdateOptions) (*Update, error)
}
