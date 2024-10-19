package malak

import (
	"context"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrUserNotFound = MalakError("user not found")
	ErrUserExists   = MalakError("User with email already exists")
)

// ENUM(admin,member,billing,investor,guest)
type Role string

type UserRole struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id,omitempty"`

	Role Role `json:"role,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	WorkspaceID uuid.UUID `json:"workspace_id,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`

	bun.BaseModel `bun:"table:roles" json:"-"`
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
	ID    uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id"`
	Email Email     `json:"email"`

	FullName string        `json:"full_name"`
	Metadata *UserMetadata `json:"metadata" `

	Roles UserRoles `json:"roles" bun:"rel:has-many,join:id=user_id"`

	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at" `
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at" `
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"-,omitempty" `

	bun.BaseModel `bun:"table:users" json:"-"`
}

func (u *User) HasWorkspace() bool { return u.Metadata.CurrentWorkspace != uuid.Nil }

type FindUserOptions struct {
	Email Email `json:"email,omitempty"`
	ID    uuid.UUID
}

type UserRepository interface {
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	Get(context.Context, *FindUserOptions) (*User, error)
}
