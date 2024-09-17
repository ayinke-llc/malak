package malak

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
	ID          uuid.UUID     `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID     `json:"workspace_id,omitempty"`
	Status      UpdateStatus  `json:"status,omitempty"`
	Reference   Reference     `json:"reference,omitempty"`
	CreatedBy   uuid.UUID     `json:"created_by,omitempty"`
	SentBy      uuid.UUID     `json:"sent_by,omitempty" bun:",nullzero"`
	Content     UpdateContent `json:"content,omitempty"`

	// Not persisted at all
	// Only calculated at runtime
	Title string `json:"title" bun:"-"`

	Metadata UpdateMetadata `json:"metadata,omitempty"`

	SentAt    *time.Time `bun:",nullzero" json:"sent_at,omitempty"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

func (u *Update) IsSent() bool { return u.Status == UpdateStatusSent }

func (u *Update) MarshalJSON() ([]byte, error) {
	type Alias Update

	title, err := getFirstHeader(u.Content)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		*Alias
		Title string `json:"title"`
	}{
		Alias: (*Alias)(u),
		Title: title,
	})
}

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

type UpdateRecipient struct {
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference Reference `json:"reference,omitempty"`
	UpdateID  uuid.UUID `json:"update_id,omitempty"`
	Email     Email     `json:"email,omitempty"`

	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" `
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" `
	bun.BaseModel `json:"-"`
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
	// Update(context.Context, *Update) error
	// Get(context.Context, FetchUpdateOptions) (*Update, error)
	List(context.Context, ListUpdateOptions) ([]Update, error)
}
