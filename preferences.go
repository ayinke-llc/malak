package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CommunicationPreferences struct {
	EnableMarketing      bool `json:"enable_marketing,omitempty"`
	EnableProductUpdates bool `json:"enable_product_updates,omitempty"`
}

type BillingPreferences struct {
	FinanceEmail Email `json:"finance_email,omitempty"`
}

type Preference struct {
	ID            uuid.UUID                `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	WorkspaceID   uuid.UUID                `json:"workspace_id,omitempty"`
	Communication CommunicationPreferences `json:"communication,omitempty"`
	Billing       BillingPreferences       `json:"billing,omitempty"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel `json:"-"`
}

type PreferenceRepository interface {
	Get(context.Context, *Workspace) (*Preference, error)
	Update(context.Context, *Preference) error
}
