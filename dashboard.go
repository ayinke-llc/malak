package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrDashboardNotFound = MalakError("dashboard not found")
)

type Dashboard struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference `json:"reference,omitempty"`
	Description string    `json:"description,omitempty"`
	Title       string    `json:"title,omitempty"`
	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`

	ChartCount int64 `json:"chart_count,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bun:",default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bun:",default:current_timestamp"`

	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type DashboardChart struct {
	ID                     uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceIntegrationID uuid.UUID `json:"workspace_integration_id,omitempty"`
	Reference              Reference `json:"reference,omitempty"`
	WorkspaceID            uuid.UUID `json:"workspace_id,omitempty"`
	DashboardID            uuid.UUID `json:"dashboard_id,omitempty"`

	ChartID          uuid.UUID         `json:"chart_id,omitempty"`
	IntegrationChart *IntegrationChart `json:"chart,omitempty" bun:"rel:belongs-to,join:chart_id=id"`

	CreatedAt time.Time `json:"created_at,omitempty" bun:",default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bun:",default:current_timestamp"`

	bun.BaseModel `json:"-"`
}

type DashboardChartPosition struct {
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	DashboardID uuid.UUID `json:"dashboard_id,omitempty"`
	ChartID     uuid.UUID `json:"chart_id,omitempty"`
	OrderIndex  int64     `json:"order_index,omitempty"`

	bun.BaseModel `json:"-"`
}

type ListDashboardOptions struct {
	Paginator   Paginator
	WorkspaceID uuid.UUID
}

type FetchDashboardOption struct {
	WorkspaceID uuid.UUID
	Reference   Reference
}

type FetchDashboardChartsOption struct {
	WorkspaceID uuid.UUID
	DashboardID uuid.UUID
}

type DashboardRepository interface {
	Create(context.Context, *Dashboard) error
	Get(context.Context, FetchDashboardOption) (Dashboard, error)

	AddChart(context.Context, *DashboardChart) error
	RemoveChart(context.Context, uuid.UUID, uuid.UUID) error

	List(context.Context, ListDashboardOptions) ([]Dashboard, int64, error)
	GetCharts(context.Context, FetchDashboardChartsOption) ([]DashboardChart, error)

	UpdateDashboardPositions(context.Context, uuid.UUID, []DashboardChartPosition) error
	GetDashboardPositions(context.Context, uuid.UUID) ([]DashboardChartPosition, error)
}
