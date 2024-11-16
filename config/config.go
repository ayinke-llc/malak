package config

import (
	"errors"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
)

// Only Postgres for now. Later on we can add support for sqlite3
// ENUM(postgres)
type DatabaseType string

// ENUM(prod,dev)
type LogMode string

// ENUM(memory)
// TODO(adelowo): add Redis support?
type RateLimiterType string

// ENUM(s3)
type UploadDriver string

// ENUM(smtp)
type EmailProvider string

type Config struct {
	Logging struct {
		Mode LogMode `yaml:"mode" mapstructure:"mode"`
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
		Port      int `yaml:"port" mapstructure:"port"`
		RateLimit struct {
			// If redis, you have to configure the redis struct in the database field
			Type              RateLimiterType `yaml:"type" mapstructure:"type"`
			IsEnabled         bool            `yaml:"is_enabled" mapstructure:"is_enabled"`
			RequestsPerMinute uint64          `yaml:"requests_per_minute" mapstructure:"requests_per_minute"`
			BurstInterval     time.Duration   `yaml:"burst_interval" mapstructure:"burst_interval"`
		} `yaml:"rate_limit" mapstructure:"rate_limit"`
		Swagger struct {
			Port      int  `mapstructure:"port" yaml:"port"`
			UIEnabled bool `mapstructure:"ui_enabled" yaml:"ui_enabled"`
		} `mapstructure:"swagger" yaml:"swagger"`
	} `yaml:"http" mapstructure:"http"`

	Billing struct {
		Stripe struct {
			APIKey    string `yaml:"api_key" mapstructure:"api_key"`
			APISecret string `yaml:"api_secret" mapstructure:"api_secret"`
		} `yaml:"stripe" mapstructure:"stripe"`

		// If stripe is not enabled, then fake ids can be used in the
		// plans table really
		// Ideally self hosted users will want to disable this
		IsEnabled bool `yaml:"is_enabled" mapstructure:"is_enabled"`

		// Newly created workspaces will have this plan automatically
		// applied upon creation
		DefaultPlanReference string `yaml:"default_plan_reference" mapstructure:"default_plan_reference"`
	} `yaml:"billing" mapstructure:"billing"`

	Uploader struct {
		Driver        UploadDriver `yaml:"driver" mapstructure:"driver"`
		MaxUploadSize int64        `yaml:"max_upload_size" mapstructure:"max_upload_size"`

		S3 struct {
			AccessKey     string `yaml:"access_key" mapstructure:"access_key"`
			AccessSecret  string `yaml:"access_secret" mapstructure:"access_secret"`
			Region        string `yaml:"region" mapstructure:"region"`
			Endpoint      string `yaml:"endpoint" mapstructure:"endpoint"`
			LogOperations bool   `yaml:"log_operations" mapstructure:"log_operations"`
			Bucket        string `yaml:"bucket" mapstructure:"bucket"`
			// Enabled by default but you can disable this if running
			// your own internal Minio or something
			UseTLS bool `yaml:"use_tls" mapstructure:"use_tls"`
		} `yaml:"s3" mapstructure:"s3"`
	} `yaml:"uploader" mapstructure:"uploader"`

	Email struct {
		Provider   EmailProvider `mapstructure:"provider" yaml:"provider"`
		Sender     malak.Email   `mapstructure:"sender" yaml:"sender"`
		SenderName string        `mapstructure:"sender_name" yaml:"sender_name"`
		SMTP       struct {
			Host     string `mapstructure:"host" yaml:"host"`
			Port     int    `mapstructure:"port" yaml:"port"`
			Username string `mapstructure:"username" yaml:"username"`
			Password string `mapstructure:"password" yaml:"password"`
			UseTLS   bool   `yaml:"use_tls" mapstructure:"use_tls"`
		} `mapstructure:"smtp" yaml:"smtp"`
	} `mapstructure:"email" yaml:"email"`

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

	if !c.Logging.Mode.IsValid() {
		return errors.New("please provide a valid logging mode")
	}

	if !c.Database.DatabaseType.IsValid() {
		return errors.New("please use a valid database provider")
	}

	if !c.Uploader.Driver.IsValid() {
		return errors.New("please provide a valid upload driver like s3")
	}

	if util.IsStringEmpty(c.Uploader.S3.AccessKey) {
		return errors.New("please provide your s3 access key")
	}

	if util.IsStringEmpty(c.Uploader.S3.AccessSecret) {
		return errors.New("please provide your s3 access secret key")
	}

	if util.IsStringEmpty(c.Uploader.S3.Bucket) {
		c.Uploader.S3.Bucket = "malak"
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

	if !c.Email.Provider.IsValid() {
		return errors.New("email provider is not currently supported")
	}

	return nil
}
