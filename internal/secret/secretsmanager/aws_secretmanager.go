package secretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"

	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
)

type awsSecretsManagerClient struct {
	svc *secretsmanager.Client
}

func New(cfg config.Config) (secret.SecretClient, error) {
	opts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.Integration.SecretsManager.Region),
		awsConfig.WithCredentialsProvider(
			awsCreds.NewStaticCredentialsProvider(
				cfg.Integration.SecretsManager.AccessKey,
				cfg.Integration.SecretsManager.AccessSecret,
				"")),
	}

	if cfg.Integration.SecretsManager.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.Integration.SecretsManager.Endpoint,
				SigningRegion: region,
			}, nil
		})
		opts = append(opts, awsConfig.WithEndpointResolverWithOptions(customResolver))
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
