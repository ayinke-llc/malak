package vault

import (
	"context"
	"io"
)

type Vault interface {
	io.Closer
	Create(context.Context) error
	Delete(context.Context) error
}
