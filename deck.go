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

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeckPreference *DeckPreference `bun:"rel:has-one,join:id=deck_id" json:"preferences,omitempty"`
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

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

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

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt     *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
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

type DeckRepository interface {
	Create(context.Context, *Deck, *CreateDeckOptions) error
	List(context.Context, *Workspace) ([]Deck, error)
	Get(context.Context, FetchDeckOptions) (*Deck, error)
	PublicDetails(context.Context, Reference) (*Deck, error)
	Delete(context.Context, *Deck) error
	UpdatePreferences(context.Context, *Deck) error
	ToggleArchive(context.Context, *Deck) error
	TogglePinned(context.Context, *Deck) error
}
