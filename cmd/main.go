package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/ayinke-llc/malak/internal/pkg/email/resend"
	"github.com/ayinke-llc/malak/internal/pkg/email/smtp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	envPrefix             = "MAKAL_"
)

func main() {

	os.Setenv("TZ", "")

	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}

func Execute() error {

	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:   "malak",
		Short: `Investors' relationship hub for founders and Indiehackers`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

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

	rootCmd.SetVersionTemplate(
		fmt.Sprintf("Version: %v\nCommit: %v\nDate: %v\n", Version, Commit, Date))

	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Config file. This is in YAML")

	addHTTPCommand(rootCmd, cfg)
	addCronCommand(rootCmd, cfg)
	addPlanCommand(rootCmd, cfg)

	cmd, _, err := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	// default to http
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"http"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}

func initializeConfig(cfg *config.Config, pathToFile string) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	setDefaults()

	viper.AddConfigPath(filepath.Join(homePath, ".config", defaultConfigFilePath))
	viper.AddConfigPath(pathToFile)
	viper.AddConfigPath(".")

	viper.SetConfigName(defaultConfigFilePath)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix(envPrefix)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return viper.Unmarshal(cfg)
}

func setDefaults() {

	viper.SetDefault("logging.mode", config.LogModeDev)

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
	viper.SetDefault("http.rate_limit.requests_per_minute", 10)
	viper.SetDefault("http.rate_limit.bursst_interval", time.Minute)

	viper.SetDefault("biling.is_enabled", false)
	viper.SetDefault("billing.default_plan", uuid.Nil)

	viper.SetDefault("auth.google.scopes", []string{"profile", "email"})
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
