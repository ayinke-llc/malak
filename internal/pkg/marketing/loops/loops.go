package loops

import (
	"context"
	"errors"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/marketing"
	"github.com/tilebox/loops-go"
)

type loopsClient struct {
	client *loops.Client
}

func New(cfg config.Config) (marketing.MarketingClient, error) {
	if hermes.IsStringEmpty(cfg.Marketing.Loops.APIKey) {
		return nil, errors.New("please provide your loops api key")
	}

	client, err := loops.NewClient(
		loops.WithAPIKey(cfg.Marketing.Loops.APIKey))
	if err != nil {
		return nil, err
	}

	return &loopsClient{
		client: client,
	}, nil
}

func (l *loopsClient) CreateContact(ctx context.Context,
	opts *marketing.CreateContactOptions) (string, error) {

	return l.client.CreateContact(ctx, &loops.Contact{
		Email:      opts.Email.String(),
		FirstName:  hermes.Ref(opts.FirstName),
		LastName:   hermes.Ref(opts.LastName),
		UserGroup:  hermes.Ref("cloud"),
		Subscribed: true,
	})
}
