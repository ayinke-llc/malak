package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			logrus.SetOutput(os.Stdout)

			var formatter logrus.Formatter = &logrus.JSONFormatter{}

			if cfg.Logging.Format == config.LogFormatText {
				formatter = &logrus.TextFormatter{}
			}

			logrus.SetFormatter(formatter)

			lvl, err := logrus.ParseLevel(cfg.Logging.Level)
			if err != nil {
				lvl = logrus.DebugLevel
			}

			logrus.SetLevel(lvl)

			if err := cfg.Validate(); err != nil {
				logrus.WithError(err).Error("could not validate configuration")
				return err
			}

			return nil
		},
	}

	rootCmd.SetVersionTemplate(
		fmt.Sprintf("Version: %v\nCommit: %v\nDate: %v\n", Version, Commit, Date))

	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Config file. This is in YAML")

	addHTTPCommand(rootCmd, cfg)

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

	viper.SetDefault("logging.level", "debug")
	viper.SetDefault("logging.format", config.LogFormatJson)

	viper.SetDefault("database.redis.dsn", "redis://localhost:3379")
	viper.SetDefault("database.postgres.database_type", config.DatabaseTypePostgres)
	viper.SetDefault("database.postgres.log_queries", true)
	viper.SetDefault("database.postgres.dsn", "postgres://makal:makal@localhost:3432/makal?sslmode=disable")

	viper.SetDefault("otel.is_enabled", true)
	viper.SetDefault("otel.use_tls", false)
	viper.SetDefault("otel.service_name", "makal")
	viper.SetDefault("otel.endpoint", "localhost:3318")

	viper.SetDefault("http.port", 3200)

	viper.SetDefault("biling.stripe.is_enabled", false)
	viper.SetDefault("billing.default_plan", uuid.Nil)
}
