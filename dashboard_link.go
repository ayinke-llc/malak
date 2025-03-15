package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(default,contact)
type DashboardLinkType uint8

type DashboardLink struct {
	ID          uuid.UUID         `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference   Reference         `json:"reference,omitempty"`
	DashboardID uuid.UUID         `json:"dashboard_id,omitempty"`
	LinkType    DashboardLinkType `json:"link_type,omitempty"`
	Token       string            `json:"token,omitempty"`
	ContactID   uuid.UUID         `json:"contact_id,omitempty" bun:",nullzero"`
	Contact     *Contact          `json:"contact,omitempty" bun:"rel:has-one,join:contact_id=id"`

	ExpiresAt *time.Time `bun:",nullzero" json:"expires_at,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty" bun:",default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" bun:",default:current_timestamp"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

// type DashboardLinkAccessLog struct {
// 	ID              uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
// 	Reference       Reference `json:"reference,omitempty"`
// 	DashboardLinkID uuid.UUID `json:"dashboard_link_id,omitempty"`
// 	ContactID       uuid.UUID `json:"contact_id,omitempty" bun:",nullzero"`
// 	Contact         *Contact  `json:"contact,omitempty" bun:"rel:has-one,join:contact_id=id"`
//
// 	IPAddress string `json:"ip_address,omitempty"`
// 	Device    string `json:"device,omitempty"`
//
// 	CreatedAt time.Time  `json:"created_at,omitempty" bun:",default:current_timestamp"`
// 	UpdatedAt time.Time  `json:"updated_at,omitempty" bun:",default:current_timestamp"`
// 	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`
//
// 	bun.BaseModel `json:"-"`
// }

type DashboardLinkRepository interface {
	Create(context.Context, *DashboardLink) error
}
