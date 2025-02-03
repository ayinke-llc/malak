package vault

import (
	"context"
	"errors"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	vault "github.com/hashicorp/vault/api"
)

type hashicorpVault struct {
	client *vault.Client
	path   string
}

func New(cfg config.Config) (secret.SecretClient, error) {

	if hermes.IsStringEmpty(cfg.Integration.Vault.Path) {
		return nil, errors.New("please provide your vault mount secret path")
	}

	if hermes.IsStringEmpty(cfg.Integration.Vault.Token) {
		return nil, errors.New("please provide your vault token")
	}

	if hermes.IsStringEmpty(cfg.Integration.Vault.Address) {
		return nil, errors.New("please provide your vault address")
	}

	c := vault.DefaultConfig()

	c.Address = cfg.Integration.Vault.Address

	client, err := vault.NewClient(c)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Integration.Vault.Token)

	return &hashicorpVault{
		client: client,
	}, nil
}

func (h *hashicorpVault) Close() error {
	return nil
}

func (h *hashicorpVault) Get(ctx context.Context,
	key string) (string, error) {

	resp, err := h.client.KVv2(h.path).
		Get(ctx,
			key)
	if err != nil {
		return "", err
	}

	data := resp.Data["data"]
	val, ok := data.(string)
	if !ok {
		return "", errors.New("data does not exists")
	}

	return val, nil
}

func (h *hashicorpVault) Create(ctx context.Context,
	opts *secret.CreateSecretOptions) (string, error) {

	_, err := h.client.KVv2(h.path).
		Put(ctx,
			opts.Key(), map[string]interface{}{
				"data": opts.Value,
			})

	return opts.Key(), err
}
