package socialauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var tracer = otel.Tracer("malak.socialauth")

func getTracer(ctx context.Context, r *http.Request,
	operationName string) (context.Context, trace.Span) {

	ctx, span := tracer.Start(ctx, operationName)

	return ctx, span
}

type googleAuthenticator struct {
	cfg    *oauth2.Config
	client *http.Client
}

func NewGoogle(cfg config.Config) SocialAuthProvider {
	return &googleAuthenticator{
		cfg: &oauth2.Config{
			ClientID:     cfg.Auth.Google.ClientID,
			ClientSecret: cfg.Auth.Google.ClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  cfg.Auth.Google.RedirectURI,
			Scopes:       cfg.Auth.Google.Scopes,
		},
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (g *googleAuthenticator) Validate(
	ctx context.Context, opts ValidateOptions) (*oauth2.Token, error) {
	return g.cfg.Exchange(ctx, opts.Code)
}

func (g *googleAuthenticator) User(ctx context.Context, accessToken string) (User, error) {

	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"

	resp, err := g.client.Get(fmt.Sprintf("%s?access_token=%s", userInfoEndpoint, accessToken))
	if err != nil {
		return User{}, err
	}

	defer resp.Body.Close()

	user := User{}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
