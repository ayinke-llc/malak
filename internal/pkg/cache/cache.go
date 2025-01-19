package cache

import (
	"context"
	"time"

	"github.com/ayinke-llc/malak"
)

const (
	ErrCacheMiss = malak.MalakError("cache miss")
)

type Cache interface {
	Add(context.Context, string, []byte, time.Duration) error
	Exists(context.Context, string) (bool, error)
	Get(context.Context, string) ([]byte, error)
}
