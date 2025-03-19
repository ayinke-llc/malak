package main

import (
	"context"
	"errors"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/datastore/postgres"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func addIntegrationCommand(c *cobra.Command, cfg *config.Config) {

	cmd := &cobra.Command{
		Use:   "integrations",
		Short: "Manage your system wide integrations",
	}

	cmd.AddCommand(createIntegration(c, cfg))

	c.AddCommand(cmd)
}

func createIntegration(_ *cobra.Command, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: `Create a new integration`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				name                 string
				logoURL              string
				description          string
				isEnabledIntegration bool

				integrationType malak.IntegrationType
			)

			logger, err := getLogger(hermes.DeRef(cfg))
			if err != nil {
				return err
			}

			db, err := postgres.New(cfg, logger)
			if err != nil {
				logger.Error("could not connect to postgres database",
					zap.Error(err))
				return err
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("What’s the name of this integration?").
						Value(&name).
						Validate(func(str string) error {
							if hermes.IsStringEmpty(str) {
								return errors.New("please provide the name of your integration")
							}

							if len(str) > 30 {
								return errors.New("integration name cannot be more than 30")
							}

							return nil
						}),

					huh.NewInput().
						Title("Logo url?").
						Value(&logoURL).
						Validate(func(str string) error {
							if hermes.IsStringEmpty(str) {
								return errors.New("please provide the url to your logo")
							}

							return nil
						}),

					huh.NewInput().
						Title("Write a description for this integration?").
						CharLimit(300).
						Value(&description),

					huh.NewSelect[malak.IntegrationType]().
						Title("Pick an integration type").
						Options(
							huh.NewOption("oauth2", malak.IntegrationTypeOauth2),
							huh.NewOption("api key", malak.IntegrationTypeApiKey),
						).
						Value(&integrationType),

					huh.NewConfirm().
						Title("Is enabled?").
						Description("should this integration be enabled?").
						Value(&isEnabledIntegration),
				),
			)

			defer db.Close()

			if err := form.Run(); err != nil {
				logger.Error("could not take in data from user",
					zap.Error(err))
				return err
			}

			logger.Debug("creating a new integration")

			valid, err := malak.IsImageFromURL(logoURL)
			if err != nil {
				return errors.New("integration name cannot be more than 30")
			}

			if !valid {
				return errors.New("please provide a valid logo url")
			}

			integration := &malak.Integration{
				Reference:       malak.NewReferenceGenerator().Generate(malak.EntityTypeIntegration),
				Description:     description,
				IsEnabled:       isEnabledIntegration,
				Metadata:        malak.IntegrationMetadata{},
				IntegrationType: integrationType,
				IntegrationName: name,
				LogoURL:         logoURL,
			}

			integrationRepo := postgres.NewIntegrationRepo(db)

			if err := integrationRepo.Create(context.Background(), integration); err != nil {
				logger.Error("could not create integration", zap.Error(err))
				return err
			}

			logger.Debug("created integration")

			return nil
		},
	}

	return cmd
}
