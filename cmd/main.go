package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {

	os.Setenv("TZ", "")

	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "malak",
		Short: `Investors' relationship hub for founders and Indiehackers`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			envFile, err := cmd.Flags().GetString("env.file")
			if err != nil {
				return err
			}

			err = godotenv.Load(envFile)
			if os.IsNotExist(err) {
				return nil
			}

			return err
		},
	}

	rootCmd.PersistentFlags().String("env.file", ".env", "Load .env file")

	addHTTPCommand(rootCmd)

	return rootCmd.Execute()
}
