package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrDeckNotFound = MalakError("deck not found")
)

type PublicDeck struct {
	Reference   Reference `json:"reference,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	ShortLink   string    `json:"short_link,omitempty"`
	DeckSize    int64     `json:"deck_size,omitempty"`

	IsArchived bool `json:"is_archived,omitempty"`

	ObjectLink string `json:"object_link,omitempty"`

	CreatedAt time.Time         `json:"created_at,omitempty"`
	UpdatedAt time.Time         `json:"updated_at,omitempty"`
	Session   DeckViewerSession `json:"session,omitempty"`

	DeckPreference *PublicDeckPreference `bun:"rel:has-one,join:id=deck_id" json:"preferences,omitempty"`
}

type PublicDeckPreference struct {
	EnableDownloading bool `json:"enable_downloading,omitempty"`
	RequireEmail      bool `json:"require_email,omitempty"`
	HasPassword       bool `json:"has_password,omitempty"`
}

type Deck struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference   Reference `json:"reference,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by,omitempty"`
	Title       string    `json:"title,omitempty"`
	ShortLink   string    `json:"short_link,omitempty"`
	DeckSize    int64     `json:"deck_size,omitempty"`

	IsArchived bool `json:"is_archived,omitempty"`

	IsPinned bool `json:"is_pinned,omitempty"`

	ObjectKey string `json:"object_key,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	DeckPreference *DeckPreference `bun:"rel:has-one,join:id=deck_id" json:"preferences,omitempty"`

	bun.BaseModel `bun:"table:decks" json:"-"`
}

type DeckPreference struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference Reference `json:"reference,omitempty"`

	WorkspaceID       uuid.UUID               `json:"workspace_id,omitempty"`
	DeckID            uuid.UUID               `json:"deck_id,omitempty"`
	EnableDownloading bool                    `json:"enable_downloading,omitempty"`
	RequireEmail      bool                    `json:"require_email,omitempty"`
	Password          PasswordDeckPreferences `json:"password,omitempty"`
	ExpiresAt         *time.Time              `bun:",soft_delete,nullzero" json:"expires_at,omitempty"`

	CreatedBy uuid.UUID `json:"created_by,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type PasswordDeckPreferences struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Password Password `json:"password,omitempty"`
}

type CreateDeckOptions struct {
	RequireEmail      bool `json:"require_email,omitempty"`
	EnableDownloading bool `json:"enable_downloading,omitempty"`
	Password          struct {
		Enabled  bool     `json:"enabled,omitempty" validate:"required"`
		Password Password `json:"password,omitempty" validate:"required"`
	} `json:"password,omitempty" validate:"required"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Reference Reference  `json:"reference,omitempty"`
}

type FetchDeckOptions struct {
	Reference   string
	WorkspaceID uuid.UUID
}

type DeckViewerSession struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference Reference `json:"reference,omitempty"`
	DeckID    uuid.UUID `json:"deck_id,omitempty"`
	ContactID uuid.UUID `json:"contact_id,omitempty" bun:",nullzero"`

	SessionID Reference `json:"session_id,omitempty"`

	DeviceInfo string    `json:"device_info,omitempty"`
	OS         string    `json:"os,omitempty"`
	Browser    string    `json:"browser,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	Country    string    `json:"country,omitempty"`
	City       string    `json:"city,omitempty"`
	ViewedAt   time.Time `json:"viewed_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type UpdateDeckSessionOptions struct {
	CreateContact bool
	Contact       *Contact
	Session       *DeckViewerSession
}

type DeckRepository interface {
	Create(context.Context, *Deck, *CreateDeckOptions) error
	List(context.Context, *Workspace) ([]Deck, error)
	Get(context.Context, FetchDeckOptions) (*Deck, error)
	PublicDetails(context.Context, Reference) (*Deck, error)
	Delete(context.Context, *Deck) error
	UpdatePreferences(context.Context, *Deck) error
	ToggleArchive(context.Context, *Deck) error
	TogglePinned(context.Context, *Deck) error

	CreateDeckSession(context.Context, *DeckViewerSession) error
	UpdateDeckSession(context.Context, *UpdateDeckSessionOptions) error
	FindDeckSession(context.Context, string) (*DeckViewerSession, error)
}
