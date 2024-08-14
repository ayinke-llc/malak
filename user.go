package malak

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("User with details already exists")
)

// ENUM(admin,member,billing)
type Role string

type UserRole struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Role Role `json:"role,omitempty" bson:"role"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`

	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`

	bun.BaseModel `bun:"table:roles"`
}

type UserRoles []*UserRole

func (m UserRoles) IsEmpty() bool { return len(m) == 0 }

type Email string

func (e Email) String() string { return strings.ToLower(string(e)) }

func (e Email) Value() (driver.Value, error) { return driver.Value(e.String()), nil }

type UserMetadata struct {
	// Used to keep track of the last used workspace
	// In the instance of multiple workspaces
	// So when next the user logs in, we remember and take them to the
	// right place rather than always a list of all their workspaces and they
	// have to select one
	CurrentWorkspace uuid.UUID `json:"current_workspace"`
}

type User struct {
	ID    uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`
	Email Email     `json:"email,omitempty" bson:"email"`

	FullName string        `json:"full_name,omitempty" bson:"full_name"`
	Metadata *UserMetadata `json:"metadata,omitempty" bson:"metadata"`

	Roles UserRoles `json:"roles" bun:"rel:has-many,join:id=user_id"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" bson:"deleted_at"`

	bun.BaseModel
}

type FindUserOptions struct {
	Email         Email  `json:"email,omitempty"`
	ClerkRefernce string `json:"clerk_refernce"`
	ID            uuid.UUID
}

type UserRepository interface {
	Create(context.Context, *User, *Workspace) error
	Update(context.Context, *User) error
	Get(context.Context, *FindUserOptions) (*User, error)
}
