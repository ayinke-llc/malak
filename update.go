package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(draft,sent)
type UpdateStatus string

// ENUM(scheduled,cancelled,sent,failed)
type UpdateSendSchedule string

type UpdateContent string

type UpdateMetadata struct {
}

type Update struct {
	ID          uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID     `json:"workspace_id,omitempty"`
	Status      UpdateStatus  `json:"status,omitempty"`
	Reference   Reference     `json:"reference,omitempty"`
	CreatedBy   uuid.UUID     `json:"created_by,omitempty"`
	SentBy      uuid.UUID     `json:"sent_by,omitempty"`
	Content     UpdateContent `json:"content,omitempty"`
	PublicLink  *string       `json:"public_link,omitempty"`

	Metadata UpdateMetadata `json:"metadata,omitempty"`

	SentAt    *time.Time `bun:"nullzero" json:"sent_at,omitempty"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type UpdateRecipient struct {
	ID       uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	UpdateID uuid.UUID `json:"update_id,omitempty"`
	Email    Email     `json:"email,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
}

type UpdateSchedule struct {
	ID          uuid.UUID          `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	UpdateID    uuid.UUID          `json:"update_id,omitempty"`
	ScheduledBy uuid.UUID          `json:"scheduled_by,omitempty"`
	Status      UpdateSendSchedule `json:"status,omitempty"`

	// Time to send this update at?
	SendAt uuid.UUID `json:"send_at,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
}

func (u *Update) IsSent() bool { return u.Status == UpdateStatusSent }

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
