package secretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
)

type awsSecretsManagerClient struct {
	svc *secretsmanager.Client
	URL string
}

func New(cfg config.Config) (secret.SecretClient, error) {
	opts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.Secrets.SecretsManager.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.Secrets.SecretsManager.AccessKey,
				cfg.Secrets.SecretsManager.AccessSecret,
				"")),
	}

	if cfg.Secrets.SecretsManager.Endpoint != "" {
		opts = append(opts, awsConfig.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.Secrets.SecretsManager.Endpoint,
					SigningRegion: region,
				}, nil
			}),
		))
	}

	conf, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	svc := secretsmanager.NewFromConfig(conf)

	return &awsSecretsManagerClient{
		svc: svc,
		URL: cfg.Secrets.SecretsManager.Endpoint,
	}, nil
}

func (i *awsSecretsManagerClient) Close() error {
	return nil
}

func (i *awsSecretsManagerClient) Get(ctx context.Context,
	key string) (string, error) {

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     hermes.Ref(key),
		VersionStage: hermes.Ref("AWSCURRENT"),
	}

	result, err := i.svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return hermes.DeRef(result.SecretString), nil
}

func (i *awsSecretsManagerClient) Create(ctx context.Context,
	opts *secret.CreateSecretOptions) (string, error) {

	_, err := i.svc.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         hermes.Ref(opts.Key()),
		SecretString: hermes.Ref(opts.Value),
	})
	return opts.Key(), err
}
