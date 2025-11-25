package cli

import (
	"fmt"
	"os"

	"github.com/ayinke-llc/malak/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hermesconfig "github.com/ayinke-llc/hermes/config"
)

func addConfigGenerateCommand(rootCmd *cobra.Command) {
	configGenCmd := &cobra.Command{
		Use:   "export-config",
		Short: "export configuration to yml, json, and .env.example",
		Long:  "Exports the default configuration to config/config.example.yml, config/config.example.json, and config/.env.example files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg := &config.Config{}
			viper.SetEnvPrefix(envPrefix)
			setDefaults()

			if err := viper.Unmarshal(cfg); err != nil {
				return fmt.Errorf("failed to unmarshal defaults: %w", err)
			}

		ymlFile, err := os.Create("config/config.example.yml")
		if err != nil {
			return fmt.Errorf("failed to create config/config.example.yml: %w", err)
		}
		defer ymlFile.Close()

		if err := hermesconfig.Export(ymlFile, cfg, hermesconfig.ExportTypeYml, envPrefix); err != nil {
			return fmt.Errorf("failed to export yml: %w", err)
		}
		fmt.Println("Exported config to config/config.example.yml")

		jsonFile, err := os.Create("config/config.example.json")
		if err != nil {
			return fmt.Errorf("failed to create config/config.example.json: %w", err)
		}
		defer jsonFile.Close()

		if err := hermesconfig.Export(jsonFile, cfg, hermesconfig.ExportTypeJson, envPrefix); err != nil {
			return fmt.Errorf("failed to export json: %w", err)
		}
		fmt.Println("Exported config to config/config.example.json")

			envFile, err := os.Create("config/.env.example")
			if err != nil {
				return fmt.Errorf("failed to create config/.env.example: %w", err)
			}
			defer envFile.Close()

			if err := hermesconfig.Export(envFile, cfg, hermesconfig.ExportTypeEnv, envPrefix); err != nil {
				return fmt.Errorf("failed to export env: %w", err)
			}
			fmt.Println("Exported config to config/.env.example")

			return nil
		},
	}

	rootCmd.AddCommand(configGenCmd)
}
