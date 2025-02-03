package infisical

import (
	"context"
	"errors"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	infisical "github.com/infisical/go-sdk"
)

type infisicalClient struct {
	client      infisical.InfisicalClientInterface
	environment string
	projectID   string
}

func New(cfg config.Config) (secret.SecretClient, error) {

	if hermes.IsStringEmpty(cfg.Integration.Infisical.Environment) {
		return nil, errors.New("please add your infisical environment")
	}

	if hermes.IsStringEmpty(cfg.Integration.Infisical.ProjectID) {
		return nil, errors.New("please provide your infisical projectID")
	}

	if hermes.IsStringEmpty(cfg.Integration.Infisical.ClientID) {
		return nil, errors.New("please provide your infisical client id")
	}

	if hermes.IsStringEmpty(cfg.Integration.Infisical.ClientSecret) {
		return nil, errors.New("please provide your infisical client secret")
	}

	if hermes.IsStringEmpty(cfg.Integration.Infisical.SiteURL) {
		return nil, errors.New("please provide your infisical site url")
	}

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          cfg.Integration.Infisical.SiteURL,
		AutoTokenRefresh: true,
	})

	_, err := client.Auth().
		UniversalAuthLogin(cfg.Integration.Infisical.ClientID, cfg.Integration.Infisical.ClientSecret)
	if err != nil {
		return nil, err
	}

	return &infisicalClient{
		client:      client,
		environment: cfg.Integration.Infisical.Environment,
		projectID:   cfg.Integration.Infisical.ProjectID,
	}, nil
}

func (i *infisicalClient) Close() error {
	return nil
}

func (i *infisicalClient) Get(ctx context.Context,
	key string) (string, error) {

	apiKeySecret, err := i.client.Secrets().
		Retrieve(infisical.RetrieveSecretOptions{
			SecretKey:   key,
			Environment: i.environment,
			ProjectID:   i.projectID,
			SecretPath:  "/",
		})
	if err != nil {
		return "", err
	}

	return apiKeySecret.SecretValue, nil
}

func (i *infisicalClient) Create(ctx context.Context,
	opts *secret.CreateSecretOptions) (string, error) {

	_, err := i.client.Secrets().
		Create(infisical.CreateSecretOptions{
			ProjectID:   i.projectID,
			Environment: i.environment,
			SecretKey:   opts.Key(),
			SecretValue: opts.Value,
		})

	return opts.Key(), err
}
