package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	LogLevel string `json:"level,omitempty" env:"LOG_LEVEL"`

	OtelEndpoint string `json:"otel_endpoint,omitempty" env:"OTEL_ENDPOINT"`
	OtelUseTLS   bool   `env:"OTEL_USE_TLS" json:"otel_use_tls,omitempty"`
	OtelHeaders  string `env:"OTEL_HEADERS"`

	PostgresDSN string `env:"POSTGRES_DSN" json:"postgres_dsn,omitempty"`

	// DISABLE IN PROD
	PostgresLogQueries bool `env:"POSTGRES_LOG_QUERIES" json:"postgres_log_queries,omitempty"`

	RedisDSN string `env:"REDIS_DSN" json:"redis_dsn,omitempty"`
}

func Load() (Config, error) {
	var cfg Config

	return cfg, env.Parse(&cfg)
}
