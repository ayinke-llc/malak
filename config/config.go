package config

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
}
