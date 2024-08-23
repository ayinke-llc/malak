package config

import (
	"errors"
	"time"

	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/google/uuid"
)

// Only Postgres for now. Later on we can add support for sqlite3
// ENUM(postgres)
type DatabaseType string

// ENUM(json,text)
type LogFormat string

type Config struct {
	Logging struct {
		Level  string    `yaml:"level" mapstructure:"level"`
		Format LogFormat `yaml:"format" mapstructure:"format"`
	} `yaml:"logging" mapstructure:"logging"`

	Database struct {
		DatabaseType DatabaseType `yaml:"database_type" mapstructure:"database_type"`
		Postgres     struct {
			DSN        string `yaml:"dsn" mapstructure:"dsn"`
			LogQueries bool   `yaml:"log_queries" mapstructure:"log_queries"`
			// How much timeout should be used to run queries agains the db or the
			// context dies
			QueryTimeout time.Duration `yaml:"query_timeout" mapstructure:"query_timeout"`
		} `yaml:"postgres" mapstructure:"postgres"`

		Redis struct {
			DSN string `yaml:"dsn" mapstructure:"dsn"`
		} `yaml:"redis" mapstructure:"redis"`
	} `yaml:"database" mapstructure:"database"`

	Otel struct {
		Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
		UseTLS    bool   `yaml:"use_tls" mapstructure:"use_tls"`
		Headers   string `yaml:"headers" mapstructure:"headers"`
		IsEnabled bool   `yaml:"is_enabled" mapstructure:"is_enabled"`
	} `yaml:"otel" mapstructure:"otel"`

	HTTP struct {
		Port int `yaml:"port" mapstructure:"port"`
	} `yaml:"http" mapstructure:"http"`

	Billing struct {
		Stripe struct {
			APIKey    string `yaml:"api_key" mapstructure:"api_key"`
			APISecret string `yaml:"api_secret" mapstructure:"api_secret"`

			// If stripe is not enabled, then fake ids can be used in the
			// plans table really
			// Ideally self hosted users will want to disable this
			IsEnabled bool `yaml:"is_enabled" mapstructure:"is_enabled"`
		} `yaml:"stripe" mapstructure:"stripe"`

		// Newly created workspaces will have this plan automatically
		// applied upon creation
		DefaultPlan uuid.UUID `yaml:"default_plan" mapstructure:"default_plan"`
	} `yaml:"billing" mapstructure:"billing"`

	Email struct {
	}

	Auth struct {
		Google struct {
			ClientID     string   `yaml:"client_id" mapstructure:"client_id"`
			ClientSecret string   `yaml:"client_secret" mapstructure:"client_secret"`
			RedirectURI  string   `yaml:"redirect_uri" mapstructure:"redirect_uri"`
			Scopes       []string `yaml:"scopes" mapstructure:"scopes"`
			IsEnabled    bool     `yaml:"is_enabled" mapstructure:"is_enabled"`
		} `yaml:"google" mapstructure:"google"`

		JWT struct {
			Key string `yaml:"key" mapstructure:"key"`
		} `yaml:"jwt" mapstructure:"jwt"`
	} `yaml:"auth" mapstructure:"auth"`
}

func (c *Config) Validate() error {

	if !c.Database.DatabaseType.IsValid() {
		return errors.New("please use a valid database provider")
	}

	if !c.Auth.Google.IsEnabled {
		return errors.New("at least one oauth authentication provider has to be turned on")
	}

	if c.Auth.Google.IsEnabled {

		if util.IsStringEmpty(c.Auth.Google.ClientID) {
			return errors.New("please provide Google oauth key")
		}

		if util.IsStringEmpty(c.Auth.Google.ClientSecret) {
			return errors.New("please provide Google oauth secret")
		}
	}

	if util.IsStringEmpty(c.Auth.JWT.Key) {
		return errors.New("please provide your JWT key")
	}

	return nil
}
