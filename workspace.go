package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var ErrWorkspaceNotFound = malakError("workspace not found")

type Workspace struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	WorkspaceName string    `json:"workspace_name,omitempty"`
	PlanID        uuid.UUID `json:"plan_id,omitempty"`
	Reference     string    `json:"reference,omitempty"`

	// Not required
	// Dummy values work really
	StripeCustomerID string `json:"stripe_customer_id,omitempty"`
	SubscriptionID   string `json:"subscription_id,omitempty"`

	Metadata PlanMetadata `json:"metadata,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`
}

func NewWorkspace(name string, u *User,
	plan *Plan, reference string) *Workspace {
	return &Workspace{
		WorkspaceName: name,
		Reference:     reference,
		Metadata:      plan.Metadata,
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
	// MarkInActive(context.Context, *Workspace) error
	// MarkActive(context.Context, *Workspace) error
}
