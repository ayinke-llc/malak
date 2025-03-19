package malak

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrAPIKeyNotFound = errors.New("api key not found")
	ErrAPIKeyMaxLimit = errors.New("you can only have a maximum of 15 active api keys")
)

type APIKey struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by,omitempty"`
	Reference   Reference `json:"reference,omitempty"`
	Value       string    `json:"value,omitempty"`
	KeyName     string    `json:"key_name,omitempty"`

	ExpiresAt *time.Time `bun:",nullzero" json:"expires_at,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

// ENUM(immediate,day,week)
type RevocationType string

type RevokeAPIKeyOptions struct {
	APIKey         *APIKey
	RevocationType RevocationType
}

type APIKeyRepository interface {
	Create(context.Context, *APIKey) error
	Revoke(context.Context, RevokeAPIKeyOptions) error
}
