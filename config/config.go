package config

import (
	"github.com/caarlos0/env/v9"
)

// ENUM(postgres)
type DatabaseType string

type Config struct {
	Logging struct {
		Level string `yaml:"level" mapstructure:"level"`
	} `yaml:"logging" mapstructure:"logging"`

	Database struct {
		Postgres struct {
			DSN        string `yaml:"dsn" mapstructure:"dsn"`
			LogQueries bool   `yaml:"log_queries" mapstructure:"log_queries"`
		} `yaml:"postgres" mapstructure:"postgres"`

		Redis struct {
			DSN string `yaml:"dsn" mapstructure:"dsn"`
		} `yaml:"redis" mapstructure:"redis"`
	} `yaml:"database" mapstructure:"database"`

	Otel struct {
		Endpoint    string            `yaml:"endpoint" mapstructure:"endpoint"`
		UseTLS      bool              `yaml:"use_tls" mapstructure:"use_tls"`
		Headers     map[string]string `yaml:"headers" mapstructure:"headers"`
		IsEnabled   bool              `yaml:"is_enabled" mapstructure:"is_enabled"`
		ServiceName string            `yaml:"service_name" mapstructure:"service_name"`
	} `yaml:"otel" mapstructure:"otel"`

	HTTP struct {
		Port int `yaml:"port" mapstructure:"port"`
	} `yaml:"http" mapstructure:"http"`
}

func Load() (Config, error) {
	var cfg Config

	return cfg, env.Parse(&cfg)
}
