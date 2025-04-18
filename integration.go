package malak

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrWorkspaceIntegrationNotFound = MalakError("integration not found")
	ErrChartNotFound                = MalakError("chart not found")
)

// ENUM(oauth2,api_key,system)
type IntegrationType string

// ENUM(stripe,paystack,flutterwave,mercury,brex)
type IntegrationProvider string

// ENUM(mercury_account,mercury_account_transaction,brex_account,brex_account_transaction)
type IntegrationChartInternalNameType string

// ENUM(bar,pie)
type IntegrationChartType string

// ENUM(daily,monthly)
type IntegrationChartFrequencyType uint8

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
	AccessToken   AccessToken `json:"access_token,omitempty"`
	LastFetchedAt time.Time   `json:"last_fetched_at,omitempty"`
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

	Metadata WorkspaceIntegrationMetadata `json:"metadata,omitempty"`

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
	IntegrationChartID     uuid.UUID                    `json:"integration_chart_id,omitempty"`
	Reference              Reference                    `json:"reference,omitempty"`
	PointName              string                       `json:"point_name,omitempty"`
	PointValue             int64                        `json:"point_value,omitempty"`
	Metadata               IntegrationDataPointMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type IntegrationChartMetadata struct {
	ProviderID string `json:"provider_id,omitempty"`
}

type IntegrationChart struct {
	ID                     uuid.UUID                        `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceIntegrationID uuid.UUID                        `json:"workspace_integration_id,omitempty"`
	WorkspaceID            uuid.UUID                        `json:"workspace_id,omitempty"`
	Reference              Reference                        `json:"reference,omitempty"`
	UserFacingName         string                           `json:"user_facing_name,omitempty"`
	InternalName           IntegrationChartInternalNameType `json:"internal_name,omitempty"`
	Metadata               IntegrationChartMetadata         `json:"metadata,omitempty"`
	ChartType              IntegrationChartType             `json:"chart_type,omitempty"`
	DataPointType          IntegrationDataPointType         `json:"data_point_type,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

// keeping it simple here with a new struct
type IntegrationDataValues struct {
	// here so it is easy to find the chart this data point belongs to
	// without too much voodoo
	// InternalName + ProviderID search in db
	// We cannot use only InternalName becasue some integrations
	// like mercury have the same InternalName twice ( each account has a savings and checkings which we track)
	InternalName   IntegrationChartInternalNameType
	UserFacingName string
	ProviderID     string
	DataPointType  IntegrationDataPointType
	Data           IntegrationDataPoint
}

type IntegrationChartValues struct {
	InternalName   IntegrationChartInternalNameType
	UserFacingName string
	ProviderID     string
	ChartType      IntegrationChartType
	DataPointType  IntegrationDataPointType
}

type IntegrationFetchDataOptions struct {
	IntegrationID      uuid.UUID
	WorkspaceID        uuid.UUID
	ReferenceGenerator ReferenceGeneratorOperation
	LastFetchedAt      time.Time
}

type IntegrationProviderClient interface {
	Name() IntegrationProvider

	// Ping tests the connection to make sure we have an
	// active connection
	Ping(context.Context, AccessToken) ([]IntegrationChartValues, error)

	Data(context.Context, AccessToken, *IntegrationFetchDataOptions) ([]IntegrationDataValues, error)

	io.Closer
}

type AccessToken string

func (a AccessToken) String() string { return string(a) }

type FindWorkspaceIntegrationOptions struct {
	Reference   Reference
	ID          uuid.UUID
	WorkspaceID uuid.UUID
}

type FetchChartOptions struct {
	WorkspaceID uuid.UUID
	Reference   Reference
}

type IntegrationRepository interface {
	Create(context.Context, *Integration) error
	System(context.Context) ([]Integration, error)

	List(context.Context, *Workspace) ([]WorkspaceIntegration, error)
	Get(context.Context, FindWorkspaceIntegrationOptions) (*WorkspaceIntegration, error)
	Disable(context.Context, *WorkspaceIntegration) error
	Update(context.Context, *WorkspaceIntegration) error

	CreateCharts(context.Context, *WorkspaceIntegration, []IntegrationChartValues) error
	AddDataPoint(context.Context, *WorkspaceIntegration, []IntegrationDataValues) error
	ListCharts(context.Context, uuid.UUID) ([]IntegrationChart, error)
	GetChart(context.Context, FetchChartOptions) (IntegrationChart, error)
	GetDataPoints(context.Context, IntegrationChart) ([]IntegrationDataPoint, error)
}
