package malak

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(oauth2,api_key)
type IntegrationType string

// ENUM(stripe,paystack,flutterwave,mercury,brex)
type IntegrationProvider string

// ENUM(mercury_account)
type IntegrationChartInternalNameType string

type IntegrationMetadata struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type Integration struct {
	ID              uuid.UUID       `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
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
	AccessToken string `json:"access_token,omitempty"`
}

type WorkspaceIntegration struct {
	ID            uuid.UUID    `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference     Reference    `json:"reference,omitempty"`
	WorkspaceID   uuid.UUID    `json:"workspace_id,omitempty"`
	IntegrationID uuid.UUID    `json:"integration_id,omitempty"`
	Integration   *Integration `bun:"rel:belongs-to,join:integration_id=id" json:"integration,omitempty"`

	// IsEnabled - this integration is enabled and data can be fetched
	IsEnabled bool `json:"is_enabled,omitempty"`

	// IsActive determines if the connection to the integration has been tested and works
	IsActive bool `json:"is_active,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

// ENUM(currency,others)
type IntegrationDataPointType string

type IntegrationDataPointMetadata struct {
}

type IntegrationDataPoint struct {
	ID                     uuid.UUID                    `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceIntegrationID uuid.UUID                    `json:"workspace_integration_id,omitempty"`
	WorkspaceID            uuid.UUID                    `json:"workspace_id,omitempty"`
	Reference              Reference                    `json:"reference,omitempty"`
	PointName              string                       `json:"point_name,omitempty"`
	PointValue             int64                        `json:"point_value,omitempty"`
	DataPointType          IntegrationDataPointType     `json:"data_point_type,omitempty"`
	Metadata               IntegrationDataPointMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type IntegrationChart struct {
	ID                     uuid.UUID                        `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceIntegrationID uuid.UUID                        `json:"workspace_integration_id,omitempty"`
	WorkspaceID            uuid.UUID                        `json:"workspace_id,omitempty"`
	Reference              Reference                        `json:"reference,omitempty"`
	UserFacingName         string                           `json:"user_facing_name,omitempty"`
	InternalName           IntegrationChartInternalNameType `json:"internal_name,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type IntegrationProviderClient interface {
	Name() IntegrationProvider
	// Ping tests the connection to make sure we have an
	// active connection
	Ping(context.Context, AccessToken) error
	io.Closer
}

type AccessToken string

func (a AccessToken) String() string { return string(a) }

type IntegrationRepository interface {
	Create(context.Context, *Integration) error
	System(context.Context) ([]Integration, error)

	List(context.Context, *Workspace) ([]WorkspaceIntegration, error)
}
