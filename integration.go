package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(malak,oauth2,api_key)
//
// malak is  the native integration type that allows you submit data via API
// api_key is when an integration does not support oauth2 and you have to share a
// read only token
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

	Metadata IntegrationMetadata `json:"metadata,omitempty" bson:"metadata"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type WorkspaceIntegrationMetadata struct {
}

type WorkspaceIntegration struct {
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()" json:"id,omitempty"`
	Reference     Reference `json:"reference,omitempty"`
	WorkspaceID   uuid.UUID `json:"workspace_id,omitempty"`
	IntegrationID uuid.UUID `json:"integration_id,omitempty"`
	IsEnabled     bool      `json:"is_enabled,omitempty"`

	Metadata WorkspaceIntegrationMetadata `json:"metadata,omitempty" bson:"metadata"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type IntegrationRepository interface {
}
