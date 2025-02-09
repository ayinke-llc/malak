package malak

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrWorkspaceNotFound = MalakError("workspace not found")
)

type WorkspaceMetadata struct {
}

type Workspace struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	WorkspaceName string `json:"workspace_name,omitempty"`
	Reference     string `json:"reference,omitempty"`
	Timezone      string `json:"-"`
	Website       string `json:"website,omitempty"`
	LogoURL       string `json:"logo_url,omitempty"`

	// Not required
	// Dummy values work really if not using stripe
	StripeCustomerID     string `json:"stripe_customer_id,omitempty"`
	SubscriptionID       string `json:"subscription_id,omitempty"`
	IsSubscriptionActive bool   `json:"is_subscription_active,omitempty"`

	PlanID uuid.UUID `json:"plan_id,omitempty"`
	Plan   *Plan     `json:"plan,omitempty" bun:"rel:belongs-to,join:plan_id=id"`

	Metadata WorkspaceMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

func NewWorkspace(name string, u *User,
	plan *Plan, reference string) *Workspace {
	return &Workspace{
		WorkspaceName: name,
		Reference:     reference,
		Metadata:      WorkspaceMetadata{},
		PlanID:        plan.ID,
	}
}

type FindWorkspaceOptions struct {
	StripeCustomerID string
	ID               uuid.UUID
	Reference        Reference
}

type CreateWorkspaceOptions struct {
	User      *User
	Workspace *Workspace
}

type WorkspaceRepository interface {
	Create(context.Context, *CreateWorkspaceOptions) error
	Get(context.Context, *FindWorkspaceOptions) (*Workspace, error)
	Update(context.Context, *Workspace) error
	List(context.Context, *User) ([]Workspace, error)
	MarkInActive(context.Context, *Workspace) error
	MarkActive(context.Context, *Workspace) error
}

func (w *Workspace) MarshalJSON() ([]byte, error) {
	type Alias Workspace

	timezone := w.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	return json.Marshal(&struct {
		Timezone string `json:"timezone"`
		*Alias
	}{
		Timezone: timezone,
		Alias:    (*Alias)(w),
	})
}
