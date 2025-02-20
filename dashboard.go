package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(barchart,piechart)
type DashboardChartType uint8

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
	ID                     uuid.UUID          `json:"id,omitempty"`
	WorkspaceIntegrationID uuid.UUID          `json:"workspace_integration_id,omitempty"`
	Reference              Reference          `json:"reference,omitempty"`
	WorkspaceID            uuid.UUID          `json:"workspace_id,omitempty"`
	DashboardID            uuid.UUID          `json:"dashboard_id,omitempty"`
	DashboardType          DashboardChartType `json:"dashboard_type,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type ListDashboardOptions struct {
	Paginator   Paginator
	WorkspaceID uuid.UUID
}

type DashboardRepository interface {
	Create(context.Context, *Dashboard) error
	AddChart(context.Context, *DashboardChart) error
	List(context.Context, ListDashboardOptions) ([]Dashboard, int64, error)
}
