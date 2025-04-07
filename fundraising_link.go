package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrFundraisingLinkNotFound = MalakError("fundraising link not found")

// ENUM(default,contact)
type FundraisingLinkType string

type FundraisingLink struct {
	ID                    uuid.UUID           `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference             Reference           `json:"reference,omitempty"`
	FundraisingPipelineID uuid.UUID           `json:"fundraising_pipeline_id,omitempty"`
	LinkType              FundraisingLinkType `json:"link_type,omitempty"`
	Token                 string              `json:"token,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty" bun:",default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" bun:",default:current_timestamp"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}
