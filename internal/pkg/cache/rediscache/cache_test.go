package rediscache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedis(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	port, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	return redisClient, func() {
		redisClient.Close()
		require.NoError(t, redisContainer.Terminate(ctx))
	}
}

func TestNew(t *testing.T) {
	redisClient, cleanup := setupRedis(t)
	defer cleanup()

	c, err := New(redisClient)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestAdd(t *testing.T) {
	redisClient, cleanup := setupRedis(t)
	defer cleanup()

	c, err := New(redisClient)
	require.NoError(t, err)

	ctx := context.Background()
	key := "testKey"
	payload := []byte("testPayload")
	ttl := time.Hour

	err = c.Add(ctx, key, payload, ttl)
	require.NoError(t, err)

	// Verify the key was actually set in Redis
	val, err := redisClient.Get(ctx, makeKey(key)).Bytes()
	require.NoError(t, err)
	require.Equal(t, payload, val)

	// Verify TTL
	duration, err := redisClient.TTL(ctx, makeKey(key)).Result()
	require.NoError(t, err)
	require.True(t, duration > 0 && duration <= ttl)
}

func TestExists(t *testing.T) {
	redisClient, cleanup := setupRedis(t)
	defer cleanup()

	c, err := New(redisClient)
	require.NoError(t, err)

	ctx := context.Background()
	key := "testKey"
	payload := []byte("testPayload")
	ttl := time.Hour

	// Test when key doesn't exist
	exists, err := c.Exists(ctx, key)
	require.Error(t, err)
	require.False(t, exists)
	require.ErrorIs(t, err, cache.ErrCacheMiss)

	// Add key
	err = c.Add(ctx, key, payload, ttl)
	require.NoError(t, err)

	// Test when key exists
	exists, err = c.Exists(ctx, key)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestExistsError(t *testing.T) {
	redisClient, cleanup := setupRedis(t)
	defer cleanup()

	c, err := New(redisClient)
	require.NoError(t, err)

	ctx := context.Background()
	key := "testKey"

	// Close Redis connection to simulate error
	redisClient.Close()

	_, err = c.Exists(ctx, key)
	require.Error(t, err)
}

func TestMakeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "malak-cache-test"},
		{"", "malak-cache-"},
		{"long-key-name", "malak-cache-long-key-name"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, makeKey(tt.input))
	}
}
