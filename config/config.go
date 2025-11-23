package config

import (
	"errors"
	"net/mail"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/secret"
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

// ENUM(smtp,resend,sendgrid)
type EmailProvider string

type HTTPConfig struct {
	Port      int `yaml:"port" mapstructure:"port" json:"port"`
	RateLimit struct {
		// If redis, you have to configure the redis struct in the database field
		Type              RateLimiterType `yaml:"type" mapstructure:"type" json:"type"`
		IsEnabled         bool            `yaml:"is_enabled" mapstructure:"is_enabled" json:"is_enabled"`
		RequestsPerMinute uint64          `yaml:"requests_per_minute" mapstructure:"requests_per_minute" json:"requests_per_minute"`
		BurstInterval     time.Duration   `yaml:"burst_interval" mapstructure:"burst_interval" json:"burst_interval"`
	} `yaml:"rate_limit" mapstructure:"rate_limit" json:"rate_limit"`
	Swagger struct {
		Port      int  `mapstructure:"port" yaml:"port" json:"port"`
		UIEnabled bool `mapstructure:"ui_enabled" yaml:"ui_enabled" json:"ui_enabled"`
	} `mapstructure:"swagger" yaml:"swagger" json:"swagger"`
	Metrics struct {
		Enabled  bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
		Username string `mapstructure:"username" yaml:"username" json:"username"`
		Password string `mapstructure:"password" yaml:"password" json:"password"`
	} `mapstructure:"metrics" yaml:"metrics" json:"metrics"`
}

type Config struct {
	Logging struct {
		Mode LogMode `yaml:"mode" mapstructure:"mode" json:"mode"`
	} `yaml:"logging" mapstructure:"logging" json:"logging"`

	Frontend struct {
		AppURL string `mapstructure:"app_url" yaml:"app_url" json:"app_url"`
	} `mapstructure:"frontend" yaml:"frontend" json:"frontend"`

	Database struct {
		DatabaseType DatabaseType `yaml:"database_type" mapstructure:"database_type" json:"database_type"`
		Postgres     struct {
			DSN        string `yaml:"dsn" mapstructure:"dsn" json:"dsn"`
			LogQueries bool   `yaml:"log_queries" mapstructure:"log_queries" json:"log_queries"`
			// How much timeout should be used to run queries agains the db or the
			// context dies
			QueryTimeout time.Duration `yaml:"query_timeout" mapstructure:"query_timeout" json:"query_timeout"`
		} `yaml:"postgres" mapstructure:"postgres" json:"postgres"`

		Redis struct {
			DSN string `yaml:"dsn" mapstructure:"dsn" json:"dsn"`
		} `yaml:"redis" mapstructure:"redis" json:"redis"`
	} `yaml:"database" mapstructure:"database" json:"database"`

	Otel struct {
		Endpoint  string `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint"`
		UseTLS    bool   `yaml:"use_tls" mapstructure:"use_tls" json:"use_tls"`
		Headers   string `yaml:"headers" mapstructure:"headers" json:"headers"`
		IsEnabled bool   `yaml:"is_enabled" mapstructure:"is_enabled" json:"is_enabled"`
	} `yaml:"otel" mapstructure:"otel" json:"otel"`

	HTTP HTTPConfig `yaml:"http" mapstructure:"http" json:"http"`

	Billing struct {
		Stripe struct {
			APIKey        string `yaml:"api_key" mapstructure:"api_key" json:"api_key"`
			WebhookSecret string `yaml:"webhook_secret" mapstructure:"webhook_secret" json:"webhook_secret"`
		} `yaml:"stripe" mapstructure:"stripe" json:"stripe"`

		// If stripe is not enabled, then fake ids can be used in the
		// plans table really
		// Ideally self hosted users will want to disable this
		IsEnabled bool `yaml:"is_enabled" mapstructure:"is_enabled" json:"is_enabled"`

		TrialDays int64 `yaml:"trial_days" mapstructure:"trial_days" json:"trial_days"`

		// Newly created workspaces will have this plan automatically
		// applied upon creation
		DefaultPlanReference string `yaml:"default_plan_reference" mapstructure:"default_plan_reference" json:"default_plan_reference"`
	} `yaml:"billing" mapstructure:"billing" json:"billing"`

	Secrets struct {
		ClientTimeout time.Duration `yaml:"client_timeout" mapstructure:"client_timeout" json:"client_timeout"`

		Provider secret.SecretProvider `yaml:"provider" mapstructure:"provider" json:"provider"`

		Vault struct {
			Address string `yaml:"address" mapstructure:"address" json:"address"`
			Token   string `yaml:"token" mapstructure:"token" json:"token"`
			Path    string `yaml:"path" mapstructure:"path" json:"path"`
		} `yaml:"vault" mapstructure:"vault" json:"vault"`

		Infisical struct {
			ClientID     string `yaml:"client_id" mapstructure:"client_id" json:"client_id"`
			ClientSecret string `yaml:"client_secret" mapstructure:"client_secret" json:"client_secret"`
			SiteURL      string `yaml:"site_url" mapstructure:"site_url" json:"site_url"`
			ProjectID    string `yaml:"project_id" mapstructure:"project_id" json:"project_id"`
			Environment  string `yaml:"environment" mapstructure:"environment" json:"environment"`
		} `yaml:"infisical" mapstructure:"infisical" json:"infisical"`

		AES struct {
			Key string `yaml:"key" mapstructure:"key" json:"key"`
		} `yaml:"aes" mapstructure:"aes" json:"aes"`

		// Merge secrets/key for s3 and this?
		SecretsManager struct {
			Region       string `yaml:"region" mapstructure:"region" json:"region"`
			AccessSecret string `yaml:"access_secret" mapstructure:"access_secret" json:"access_secret"`
			AccessKey    string `yaml:"access_key" mapstructure:"access_key" json:"access_key"`
			Endpoint     string `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint"`
		} `yaml:"secrets_manager" mapstructure:"secrets_manager" json:"secrets_manager"`
	} `yaml:"secrets" mapstructure:"secrets" json:"secrets"`

	APIKey struct {
		HashSecret string `mapstructure:"hash_secret" yaml:"hash_secret" json:"hash_secret"`
	} `mapstructure:"api_key" yaml:"api_key" json:"api_key"`

	Uploader struct {
		Driver        UploadDriver `yaml:"driver" mapstructure:"driver" json:"driver"`
		MaxUploadSize int64        `yaml:"max_upload_size" mapstructure:"max_upload_size" json:"max_upload_size"`

		S3 struct {
			AccessKey     string `yaml:"access_key" mapstructure:"access_key" json:"access_key"`
			AccessSecret  string `yaml:"access_secret" mapstructure:"access_secret" json:"access_secret"`
			Region        string `yaml:"region" mapstructure:"region" json:"region"`
			Endpoint      string `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint"`
			LogOperations bool   `yaml:"log_operations" mapstructure:"log_operations" json:"log_operations"`
			Bucket        string `yaml:"bucket" mapstructure:"bucket" json:"bucket"`
			DeckBucket    string `yaml:"deck_bucket" mapstructure:"deck_bucket" json:"deck_bucket"`
			// Enabled by default but you can disable this if running
			// your own internal Minio or something
			UseTLS bool `yaml:"use_tls" mapstructure:"use_tls" json:"use_tls"`

			CloudflareBucketDomain     string `yaml:"cloudflare_bucket_domain" mapstructure:"cloudflare_bucket_domain" json:"cloudflare_bucket_domain"`
			CloudflareDeckBucketDomain string `yaml:"cloudflare_deck_bucket_domain" mapstructure:"cloudflare_deck_bucket_domain" json:"cloudflare_deck_bucket_domain"`
		} `yaml:"s3" mapstructure:"s3" json:"s_3"`
	} `yaml:"uploader" mapstructure:"uploader" json:"uploader"`

	Email struct {
		Provider   EmailProvider `mapstructure:"provider" yaml:"provider" json:"provider"`
		Sender     malak.Email   `mapstructure:"sender" yaml:"sender" json:"sender"`
		SenderName string        `mapstructure:"sender_name" yaml:"sender_name" json:"sender_name"`
		SMTP       struct {
			Host     string `mapstructure:"host" yaml:"host" json:"host"`
			Port     int    `mapstructure:"port" yaml:"port" json:"port"`
			Username string `mapstructure:"username" yaml:"username" json:"username"`
			Password string `mapstructure:"password" yaml:"password" json:"password"`
			UseTLS   bool   `yaml:"use_tls" mapstructure:"use_tls" json:"use_tls"`
		} `mapstructure:"smtp" yaml:"smtp" json:"smtp"`
		Resend struct {
			APIKey        string `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
			WebhookSecret string `mapstructure:"webhook_secret" yaml:"webhook_secret" json:"webhook_secret"`
		} `mapstructure:"resend" yaml:"resend" json:"resend"`
	} `mapstructure:"email" yaml:"email" json:"email"`

	Auth struct {
		Google struct {
			ClientID     string   `yaml:"client_id" mapstructure:"client_id" json:"client_id"`
			ClientSecret string   `yaml:"client_secret" mapstructure:"client_secret" json:"client_secret"`
			RedirectURI  string   `yaml:"redirect_uri" mapstructure:"redirect_uri" json:"redirect_uri"`
			Scopes       []string `yaml:"scopes" mapstructure:"scopes" json:"scopes"`
			IsEnabled    bool     `yaml:"is_enabled" mapstructure:"is_enabled" json:"is_enabled"`
		} `yaml:"google" mapstructure:"google" json:"google"`

		JWT struct {
			Key string `yaml:"key" mapstructure:"key" json:"key"`
		} `yaml:"jwt" mapstructure:"jwt" json:"jwt"`
	} `yaml:"auth" mapstructure:"auth" json:"auth"`

	Analytics struct {
		// We do not want to embed this into the binary.
		// 1.) Can bloat the binary by 65MB.
		// 2.) Can lead to discussions about supply chain issues
		//
		// Better to stay clear and let people load in their data themselves
		// We will add the files to the docker image and periodically update it just for simplicity sake
		// but you do not have to use them
		//
		// people can always mount volumes themselves and use their own files
		MaxMindCountryDB string `json:"max_mind_country_db,omitempty" yaml:"max_mind_country_db" mapstructure:"max_mind_country_db"`
		MaxMindCityDB    string `json:"max_mind_city_db,omitempty" yaml:"max_mind_city_db" mapstructure:"max_mind_city_db"`
	} `json:"analytics,omitempty" yaml:"analytics" mapstructure:"analytics"`
}

func (c *Config) Validate() error {

	if hermes.IsStringEmpty(c.Email.SenderName) {
		return errors.New("please provide updates sender name")
	}

	if hermes.IsStringEmpty(c.Email.Sender.String()) {
		return errors.New("please provide updates sender email")
	}

	if _, err := mail.ParseAddress(c.Email.Sender.String()); err != nil {
		return errors.New("email sender is invalid")
	}

	if !c.Logging.Mode.IsValid() {
		return errors.New("please provide a valid logging mode")
	}

	if !c.Database.DatabaseType.IsValid() {
		return errors.New("please use a valid database provider")
	}

	if !c.Uploader.Driver.IsValid() {
		return errors.New("please provide a valid upload driver like s3")
	}

	if c.HTTP.Port < 0 {
		return errors.New("please provide a valid HTTP port number greater than 0")
	}

	if c.HTTP.Port == 0 {
		c.HTTP.Port = 5300
	}

	if c.HTTP.Metrics.Enabled {
		if hermes.IsStringEmpty(c.HTTP.Metrics.Password) {
			return errors.New("metrics password must be provided if metrics is enabled")
		}

		if hermes.IsStringEmpty(c.HTTP.Metrics.Username) {
			return errors.New("metrics username must be provided if metrics is enabled")
		}
	}

	if hermes.IsStringEmpty(c.APIKey.HashSecret) {
		return errors.New("you must provide a hash secret for your api keys")
	}

	if hermes.IsStringEmpty(c.Uploader.S3.AccessKey) {
		return errors.New("please provide your s3 access key")
	}

	if hermes.IsStringEmpty(c.Uploader.S3.AccessSecret) {
		return errors.New("please provide your s3 access secret key")
	}

	if hermes.IsStringEmpty(c.Uploader.S3.Bucket) {
		c.Uploader.S3.Bucket = "malak"
	}

	if hermes.IsStringEmpty(c.Uploader.S3.DeckBucket) {
		c.Uploader.S3.DeckBucket = "deck"
	}

	if !c.Auth.Google.IsEnabled {
		return errors.New("at least one oauth authentication provider has to be turned on")
	}

	if c.Auth.Google.IsEnabled {

		if hermes.IsStringEmpty(c.Auth.Google.ClientID) {
			return errors.New("please provide Google oauth key")
		}

		if hermes.IsStringEmpty(c.Auth.Google.ClientSecret) {
			return errors.New("please provide Google oauth secret")
		}
	}

	if hermes.IsStringEmpty(c.Auth.JWT.Key) {
		return errors.New("please provide your JWT key")
	}

	if !c.Email.Provider.IsValid() {
		return errors.New("email provider is not currently supported")
	}

	// trial days is required. We do not want to
	// somehow enforce users to provide a payment method/card
	// before they can get into the app
	if c.Billing.TrialDays < 0 {
		return errors.New("trial days must be 0 or greater than 0")
	}

	return nil
}
