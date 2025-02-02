package infisical

import (
	"context"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	infisical "github.com/infisical/go-sdk"
)

type infisicalClient struct {
	client infisical.InfisicalClientInterface
}

func New(cfg config.Config) (secret.SecretClient, error) {
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          cfg.Integration.Infisical.SiteURL,
		AutoTokenRefresh: true,
	})

	return &infisicalClient{
		client: client,
	}, nil
}

func (i *infisicalClient) Close() error {
	return nil
}

func (i *infisicalClient) Get(ctx context.Context,
	key string) (string, error) {

}

func (i *infisicalClient) Create(ctx context.Context,
	opts *secret.CreateSecretOptions) (string, error) {

}
