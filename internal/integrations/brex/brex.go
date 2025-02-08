package brex

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	baseURL = "https://platform.brexapis.com/v2"
)

type CurrentUser struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Status          string `json:"status"`
	ManagerID       string `json:"manager_id"`
	DepartmentID    string `json:"department_id"`
	LocationID      string `json:"location_id"`
	TitleID         string `json:"title_id"`
	RemoteDisplayID string `json:"remote_display_id"`
}

var tracer = otel.Tracer("integrations.brex")

type brexClient struct {
	httpClient *http.Client
}

func New(cfg config.Config) (malak.IntegrationProviderClient, error) {

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   cfg.Secrets.ClientTimeout,
	}

	return &brexClient{
		httpClient: client,
	}, nil
}

func (m *brexClient) Name() malak.IntegrationProvider {
	return malak.IntegrationProviderBrex
}

func (m *brexClient) buildRequest(ctx context.Context,
	token malak.AccessToken,
	spanName, endpoint string) (*http.Request, trace.Span, error) {

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not build request")
		return nil, span, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")

	return req.WithContext(ctx), span, err
}

func (m *brexClient) Ping(
	ctx context.Context, token malak.AccessToken) error {

	req, span, err := m.buildRequest(ctx, token, "connection.ping", "/users/me")
	if err != nil {
		return err
	}

	res, err := m.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not send request")
		return err
	}

	defer res.Body.Close()

	// ignored on purpose
	_, _ = io.Copy(io.Discard, res.Body)

	if res.StatusCode != http.StatusOK {
		err = errors.New("invalid api key")
		span.SetAttributes(attribute.Int("response_code", res.StatusCode))
		return err
	}

	span.SetStatus(codes.Ok, "connection to brex was successful")
	return nil
}

func (m *brexClient) Close() error {
	m.httpClient.CloseIdleConnections()
	return nil
}
