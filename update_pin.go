package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrPinnedUpdateNotExists        = malakError("update not pinned")
	ErrPinnedUpdateCapacityExceeded = malakError(
		`you have exceeded the maximum number of pinned updates. Please unpin an update and pin this again`)

	MaximumNumberOfPinnedUpdates = 3
)

// ENUM(pin,unpin)
type PinState uint8

type PinnedUpdate struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	UpdateID    uuid.UUID `json:"update_id,omitempty"`
	PinnedBy    uuid.UUID `json:"pinned_by,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type ListPinnedUpdates struct {
	WorkspaceID uuid.UUID
}

type FetchPinnedUpdate struct {
	UpdateID uuid.UUID
}

type PinnedUpdateRepository interface {
	Pin(context.Context, *Update, PinState, *User) error
	List(context.Context, ListPinnedUpdates) ([]PinnedUpdate, error)
	Get(context.Context, FetchPinnedUpdate) (*PinnedUpdate, error)
}
