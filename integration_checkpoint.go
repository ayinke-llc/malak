package malak

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ENUM(failed,pending,success)
type IntegrationSyncCheckpointStatus string

type IntegrationSyncCheckpoint struct {
	ID                     uuid.UUID  `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Reference              Reference  `json:"reference,omitempty"`
	WorkspaceID            uuid.UUID  `bun:"workspace_id,notnull" json:"workspace_id,omitempty"`
	WorkspaceIntegrationID uuid.UUID  `json:"workspace_integration_id,omitempty"`
	LastSyncAttempt        time.Time  `bun:"last_sync_attempt" json:"last_sync_attempt,omitempty"`
	LastSuccessfulSync     *time.Time `bun:"last_successful_sync" json:"last_successful_sync,omitempty"`
	Status                 string     `bun:"status,notnull" json:"status,omitempty"`
	ErrorMessage           string     `bun:"error_message" json:"error_message,omitempty"`
	CreatedAt              time.Time  `bun:"created_at,default:current_timestamp" json:"created_at,omitempty"`
	UpdatedAt              time.Time  `bun:"updated_at,default:current_timestamp" json:"updated_at,omitempty"`
	DeletedAt              *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty"`

	bun.BaseModel `json:"-"`
}
