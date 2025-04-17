package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/email/resend"
	"github.com/ayinke-llc/malak/internal/pkg/email/smtp"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Version describes the version of the current build.
	Version = "dev"

	// Commit describes the commit of the current build.
	Commit = "none"

	// Date describes the date of the current build.
	Date = time.Now().UTC()
)

const (
	defaultConfigFilePath = "config"
	envPrefix             = "MALAK_"
)

func Execute() error {

	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:   "malak",
		Short: `Investors' relationship hub for founders and Indiehackers`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if cmd.Use == "version" {
				return nil
			}

			confFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			if err := initializeConfig(cfg, confFile); err != nil {
				return err
			}

			return cfg.Validate()
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", Version, Commit, Date.Format(time.RFC3339))
		},
	}

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Config file. This is in YAML")

	addHTTPCommand(rootCmd, cfg)
	addCronCommand(rootCmd, cfg)
	addPlanCommand(rootCmd, cfg)
	addIntegrationCommand(rootCmd, cfg)

	return rootCmd.Execute()
}

func initializeConfig(cfg *config.Config, pathToFile string) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(homePath, ".config", defaultConfigFilePath))
	viper.AddConfigPath(pathToFile)
	viper.AddConfigPath(".")

	viper.SetConfigName(defaultConfigFilePath)
	viper.SetConfigType("yml")

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	setDefaults()

	bindEnvs(viper.GetViper(), "", cfg)

	return viper.Unmarshal(cfg)
}

func bindEnvs(v *viper.Viper, prefix string, iface interface{}) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := t.Name
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}

		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		switch fieldv.Kind() {
		case reflect.Struct:
			bindEnvs(v, path, fieldv.Addr().Interface())
		default:
			envKey := strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
			if err := v.BindEnv(path, envPrefix+envKey); err != nil {
				panic(err)
			}
		}
	}
}

func setDefaults() {

	viper.SetDefault("logging.mode", config.LogModeDev)

	viper.SetDefault("frontend.app_url", "https://app.malak.vc")

	viper.SetDefault("database.redis.dsn", "redis://localhost:9379")
	viper.SetDefault("database.postgres.database_type", config.DatabaseTypePostgres)
	viper.SetDefault("database.postgres.log_queries", true)
	viper.SetDefault("database.database_type", config.DatabaseTypePostgres)
	viper.SetDefault("database.postgres.dsn", "postgres://malak:malak@localhost:9432/malak?sslmode=disable")
	viper.SetDefault("database.postgres.query_timeout", time.Second*5)

	viper.SetDefault("uploader.driver", config.UploadDriverS3)
	viper.SetDefault("uploader.max_upload_size", 10<<20) // 10MB
	viper.SetDefault("uploader.s3.use_tls", true)

	viper.SetDefault("otel.is_enabled", true)
	viper.SetDefault("otel.use_tls", false)
	viper.SetDefault("otel.service_name", "makal")
	viper.SetDefault("otel.endpoint", "localhost:9317")

	viper.SetDefault("http.port", 5300)
	viper.SetDefault("http.rate_limit.is_enabled", true)
	viper.SetDefault("http.rate_limit.type", config.RateLimiterTypeMemory)
	viper.SetDefault("http.rate_limit.requests_per_minute", 300)
	viper.SetDefault("http.rate_limit.burst_interval", time.Minute)

	viper.SetDefault("api.provider", secret.SecretProviderAesGcm)

	viper.SetDefault("biling.is_enabled", false)
	viper.SetDefault("billing.default_plan", uuid.Nil)
	viper.SetDefault("billing.trial_days", 30)

	viper.SetDefault("auth.google.scopes", []string{"profile", "email"})

	viper.SetDefault("secrets.infisical.site_url", "https://app.infisical.com")
	viper.SetDefault("secrets.infisical.environment", "prod")

	viper.SetDefault("analytics.max_mind_country_db", "internal/pkg/geolocation/maxmind/testdata/country.mmdb")
	viper.SetDefault("analytics.max_mind_city_db", "internal/pkg/geolocation/maxmind/testdata/city.mmdb")
}

func getEmailProvider(cfg config.Config) (email.Client, error) {
	switch cfg.Email.Provider {
	case config.EmailProviderResend:
		return resend.New(cfg)

	case config.EmailProviderSmtp:
		return smtp.New(cfg)

	default:
		return nil, errors.New("unsupported email provider")
	}
}

func getLogger(cfg config.Config) (*zap.Logger, error) {
	switch cfg.Logging.Mode {
	case config.LogModeProd:
		return zap.NewProduction()
	case config.LogModeDev:
		return zap.NewDevelopment()
	default:
		return zap.NewDevelopment()
	}
}
