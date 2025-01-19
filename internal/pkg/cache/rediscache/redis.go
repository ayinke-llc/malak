package rediscache

import (
	"context"
	"time"

	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/redis/go-redis/v9"
)

func makeKey(s string) string {
	return "malak-cache-" + s
}

type redisCache struct {
	inner *redis.Client
}

func New(client *redis.Client) (cache.Cache, error) {
	return &redisCache{
		inner: client,
	}, client.Ping(context.Background()).Err()
}

func (r *redisCache) Get(ctx context.Context,
	key string) ([]byte, error) {

	cmd := r.inner.Get(ctx, makeKey(key))
	return cmd.Bytes()
}

func (r *redisCache) Add(ctx context.Context,
	key string, payload []byte, ttl time.Duration) error {
	return r.inner.Set(
		ctx, makeKey(key), payload, ttl).
		Err()
}

func (r *redisCache) Exists(ctx context.Context,
	key string) (bool, error) {

	key = makeKey(key)

	resp := r.inner.Exists(ctx, key)
	if err := resp.Err(); err != nil {
		return false, cache.ErrCacheMiss
	}

	val, err := resp.Result()
	if err != nil {
		return false, err
	}

	if val > 0 {
		return true, nil
	}

	return false, cache.ErrCacheMiss
}
