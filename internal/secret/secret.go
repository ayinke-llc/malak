package secret

import (
	"context"
	"fmt"
	"io"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
)

// ENUM(vault,infisical,aes_gcm,secretsmanager)
type SecretProvider string

type CreateSecretOptions struct {
	Value           string    `json:"value,omitempty"`
	WorkspaceID     uuid.UUID `json:"workspace_id,omitempty"`
	IntegrationName malak.IntegrationProvider
}

func (s *CreateSecretOptions) Key() string {
	return fmt.Sprintf("%s/%s",
		s.WorkspaceID.String(), s.IntegrationName.String())
}

type SecretClient interface {
	io.Closer
	Create(context.Context, *CreateSecretOptions) (string, error)
	Get(context.Context, string) (string, error)
}
