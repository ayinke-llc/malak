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

	ErrUpdateScheduleNotFound = MalakError("update schedule not found")

	MaximumNumberOfPinnedUpdates = 4
)

// ENUM(draft,sent)
type UpdateStatus string

// ENUM(draft,sent,all)
type ListUpdateFilterStatus string

// ENUM(scheduled,cancelled,sent,failed,processing)
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
	Content     BlockContents `json:"content,omitempty"`
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

// ENUM(list,email)
// List is a flattened group that can contain infinite amount of emails
type RecipientType string

// ENUM(pending,sent,failed)
type RecipientStatus string

type UpdateRecipient struct {
	ID         uuid.UUID       `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference  Reference       `json:"reference,omitempty"`
	UpdateID   uuid.UUID       `json:"update_id,omitempty"`
	ContactID  uuid.UUID       `json:"contact_id,omitempty"`
	ScheduleID uuid.UUID       `json:"schedule_id,omitempty"`
	Status     RecipientStatus `json:"status,omitempty"`

	UpdateRecipientStat *UpdateRecipientStat `json:"update_recipient_stat,omitempty" bun:"rel:has-one,join:id=recipient_id"`
	Contact             *Contact             `json:"contact,omitempty" bun:"rel:has-one,join:contact_id=id"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`

	bun.BaseModel `json:"-"`
}

// ENUM(resend,sendgrid,smtp)
type UpdateRecipientLogProvider string

type UpdateRecipientLog struct {
	ID          uuid.UUID                  `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference                  `json:"reference,omitempty"`
	RecipientID uuid.UUID                  `json:"recipient_id,omitempty"`
	ProviderID  string                     `json:"provider_id,omitempty"`
	Provider    UpdateRecipientLogProvider `json:"provider,omitempty"`

	Recipient *UpdateRecipient `json:"recipient" bun:"rel:has-one,join:recipient_id=id"`

	CreatedAt     time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt     time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt     *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
	bun.BaseModel `json:"-"`
}

type UpdateRecipientStat struct {
	ID          uuid.UUID        `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference        `json:"reference,omitempty"`
	RecipientID uuid.UUID        `json:"recipient_id,omitempty"`
	Recipient   *UpdateRecipient `json:"recipient" bun:"rel:has-one,join:recipient_id=id"`

	LastOpenedAt *time.Time `bun:",soft_delete,nullzero" json:"last_opened_at,omitempty"`
	HasReaction  bool       `json:"has_reaction,omitempty"`
	IsDelivered  bool       `json:"is_delivered,omitempty"`
	IsBounced    bool       `json:"is_bounced,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type UpdateStat struct {
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference Reference `json:"reference,omitempty"`
	UpdateID  uuid.UUID `json:"update_id,omitempty"`

	TotalOpens     int64 `json:"total_opens,omitempty"`
	TotalReactions int64 `json:"total_reactions,omitempty"`
	TotalClicks    int64 `json:"total_clicks,omitempty"`
	TotalSent      int64 `json:"total_sent,omitempty"`
	UniqueOpens    int64 `json:"unique_opens,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

// ENUM(preview,live)
type UpdateType string

type UpdateSchedule struct {
	ID          uuid.UUID          `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference          `json:"reference,omitempty"`
	UpdateID    uuid.UUID          `json:"update_id,omitempty"`
	ScheduledBy uuid.UUID          `json:"scheduled_by,omitempty"`
	Status      UpdateSendSchedule `json:"status,omitempty"`
	UpdateType  UpdateType         `json:"update_type,omitempty"`

	// Time to send this update at?
	SendAt    time.Time `json:"send_at,omitempty"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`

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

type CreateUpdateOptions struct {
	Reference   func(EntityType) string
	Generator   ReferenceGeneratorOperation
	Email       Email
	WorkspaceID uuid.UUID
	Emails      []Email
	Schedule    *UpdateSchedule
	UserID      uuid.UUID
}

type UpdateRepository interface {
	Create(context.Context, *Update) error
	Update(context.Context, *Update) error
	Get(context.Context, FetchUpdateOptions) (*Update, error)
	GetByID(context.Context, uuid.UUID) (*Update, error)
	List(context.Context, ListUpdateOptions) ([]Update, int64, error)
	ListPinned(context.Context, uuid.UUID) ([]Update, error)
	Delete(context.Context, *Update) error
	TogglePinned(context.Context, *Update) error
	GetSchedule(context.Context, uuid.UUID) (*UpdateSchedule, error)
	SendUpdate(context.Context, *CreateUpdateOptions) error
	GetStatByEmailID(context.Context, string,
		UpdateRecipientLogProvider) (*UpdateRecipientLog, *UpdateRecipientStat, error)
	Stat(context.Context, *Update) (*UpdateStat, error)
	UpdateStat(context.Context, *UpdateStat, *UpdateRecipientStat) error
	RecipientStat(context.Context, *Update) ([]UpdateRecipient, error)
}
