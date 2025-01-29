package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(oauth2,api_key)
type IntegrationType string

type IntegrationMetadata struct {
}

type Integration struct {
	ID              uuid.UUID       `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`
	IntegrationName string          `json:"integration_name,omitempty"`
	Reference       Reference       `json:"reference,omitempty"`
	Description     string          `json:"description,omitempty"`
	IsEnabled       bool            `json:"is_enabled,omitempty"`
	IntegrationType IntegrationType `json:"integration_type,omitempty"`
	LogoURL         string          `json:"logo_url,omitempty"`

	Metadata IntegrationMetadata `json:"metadata,omitempty" bson:"metadata"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type WorkspaceIntegrationMetadata struct {
}

type WorkspaceIntegration struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`
	Reference   Reference `json:"reference,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`

	IntegrationID uuid.UUID    `json:"integration_id,omitempty"`
	Integration   *Integration `json:"Integration,omitempty" bun:"rel:has-one,join:integration_id=id"`

	IsEnabled bool `json:"is_enabled,omitempty"`

	Metadata WorkspaceIntegrationMetadata `json:"metadata,omitempty" bson:"metadata"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type IntegrationRepository interface {
	Create(context.Context, *Integration) error

	List(context.Context, *Workspace) ([]WorkspaceIntegration, error)
	// Disable(context.Context, *Workspace, *Integration) error
}
