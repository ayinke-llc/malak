package malak

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(most_used,recently_created,all)
type SystemTemplateFilter string

type SystemTemplate struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Reference    Reference     `json:"reference,omitempty"`
	Content      BlockContents `json:"content,omitempty"`
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	NumberOfUses int           `json:"number_of_uses,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}

type TemplateRepository interface {
	System(context.Context, SystemTemplateFilter) ([]SystemTemplate, error)
}
