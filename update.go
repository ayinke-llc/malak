package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrUpdateNotFound = MalakError("update not exists")

	ErrPinnedUpdateNotExists        = MalakError("update not pinned")
	ErrPinnedUpdateCapacityExceeded = MalakError(
		`you have exceeded the maximum number of pinned updates. Please unpin an update and pin this again`)

	MaximumNumberOfPinnedUpdates = 3
)

// ENUM(draft,sent)
type UpdateStatus string

// ENUM(draft,sent,all)
type ListUpdateFilterStatus string

// ENUM(scheduled,cancelled,sent,failed)
type UpdateSendSchedule string

type UpdateContent string

type UpdateMetadata struct {
}

type Update struct {
	ID          uuid.UUID      `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID      `json:"workspace_id,omitempty"`
	Status      UpdateStatus   `json:"status,omitempty"`
	Reference   Reference      `json:"reference,omitempty"`
	CreatedBy   uuid.UUID      `json:"created_by,omitempty"`
	SentBy      uuid.UUID      `json:"sent_by,omitempty" bun:",nullzero"`
	Content     []BlockContent `json:"content,omitempty"`
	// If this update is pinned
	IsPinned bool   `json:"is_pinned,omitempty"`
	Title    string `json:"title,omitempty"`

	Metadata UpdateMetadata `json:"metadata,omitempty"`

	SentAt    *time.Time `bun:",nullzero" json:"sent_at,omitempty"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

func (u *Update) IsSent() bool { return u.Status == UpdateStatusSent }

type UpdateLink struct {
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference Reference `json:"reference,omitempty"`
	UpdateID  uuid.UUID `json:"update_id,omitempty"`

	// Sometimes, you want to share a link containing a specific update
	// for a few minutes or seconds :)
	ExpiresAt *time.Time `json:"expires_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`

	CreatedAt     time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" `
	UpdatedAt     time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" `
	DeletedAt     *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
	bun.BaseModel `json:"-"`
}

// ENUM(preview,live)
type RecipientType string

type UpdateRecipient struct {
	ID            uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference     Reference     `json:"reference,omitempty"`
	UpdateID      uuid.UUID     `json:"update_id,omitempty"`
	Email         Email         `json:"email,omitempty"`
	RecipientType RecipientType `json:"recipient_type,omitempty"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	bun.BaseModel
}

type UpdateSchedule struct {
	ID          uuid.UUID          `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference          `json:"reference,omitempty"`
	UpdateID    uuid.UUID          `json:"update_id,omitempty"`
	ScheduledBy uuid.UUID          `json:"scheduled_by,omitempty"`
	Status      UpdateSendSchedule `json:"status,omitempty"`

	// Time to send this update at?
	SendAt        uuid.UUID `json:"send_at,omitempty"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" `
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" `
	bun.BaseModel `json:"-"`
}

type FetchUpdateOptions struct {
	Status    UpdateStatus
	Reference Reference
	ID        uuid.UUID
}

type ListUpdateOptions struct {
	Paginator   Paginator
	WorkspaceID uuid.UUID
	Status      ListUpdateFilterStatus
}

type UpdateRepository interface {
	Create(context.Context, *Update) error
	Update(context.Context, *Update) error
	Get(context.Context, FetchUpdateOptions) (*Update, error)
	List(context.Context, ListUpdateOptions) ([]Update, error)
	Delete(context.Context, *Update) error
	TogglePinned(context.Context, *Update) error
}
