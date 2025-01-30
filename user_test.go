package malak

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUser_AccessWorkspace(t *testing.T) {

	workspaceID := uuid.New()

	t.Run("has access", func(t *testing.T) {
		user := &User{
			Roles: UserRoles{
				{
					WorkspaceID: workspaceID,
				},
				{
					WorkspaceID: uuid.New(),
				},
			},
		}

		require.True(t, user.CanAccessWorkspace(workspaceID))
	})

	t.Run("no access", func(t *testing.T) {
		user := &User{
			Roles: UserRoles{
				{
					WorkspaceID: workspaceID,
				},
				{
					WorkspaceID: uuid.New(),
				},
			},
		}

		require.False(t, user.CanAccessWorkspace(uuid.New()))
	})
}
