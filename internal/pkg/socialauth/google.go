package socialauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var tracer = otel.Tracer("malak.socialauth")
var noopTracer = noop.NewTracerProvider().Tracer("malak.socialauth")

func getTracer(ctx context.Context,
	operationName string, isTracingEnabled bool) (context.Context, trace.Span) {

	if !isTracingEnabled {
		return noopTracer.Start(ctx, operationName)
	}

	return tracer.Start(ctx, operationName)
}

type googleAuthenticator struct {
	cfg    *oauth2.Config
	client *http.Client
	config config.Config
}

func NewGoogle(cfg config.Config) SocialAuthProvider {
	return &googleAuthenticator{
		config: cfg,
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

	ctx, span := getTracer(ctx, "google.Validate", g.config.Otel.IsEnabled)
	defer span.End()

	return g.cfg.Exchange(ctx, opts.Code)
}

func (g *googleAuthenticator) User(ctx context.Context, token *oauth2.Token) (User, error) {

	ctx, span := getTracer(ctx, "google.User", g.config.Otel.IsEnabled)
	defer span.End()

	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"

	urlEndpoint := fmt.Sprintf("%s?access_token=%s", userInfoEndpoint, token.AccessToken)
	req, err := http.NewRequest(http.MethodGet, urlEndpoint, strings.NewReader(""))
	if err != nil {
		return User{}, err
	}

	resp, err := g.client.Do(req.WithContext(ctx))
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
