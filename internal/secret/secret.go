package secret

import (
	"context"
	"io"

	"github.com/ayinke-llc/malak"
	"github.com/google/uuid"
)

// ENUM(hashicorp_vault,infisical,aes_gcm)
type SecretProvider string

type CreateSecretOptions struct {
	Value           string    `json:"value,omitempty"`
	WorkspaceID     uuid.UUID `json:"workspace_id,omitempty"`
	IntegrationName malak.IntegrationProvider
}

type SecretClient interface {
	io.Closer
	Create(context.Context, *CreateSecretOptions) (string, error)
	Get(context.Context, string) (string, error)
}
