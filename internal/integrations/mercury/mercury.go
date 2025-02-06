package mercury

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	baseURL = "https://api.mercury.com/api/v1"
)

type Account struct {
	AccountNumber          string    `json:"accountNumber"`
	AvailableBalance       float64   `json:"availableBalance"`
	CreatedAt              time.Time `json:"createdAt"`
	CurrentBalance         float64   `json:"currentBalance"`
	ID                     string    `json:"id"`
	Kind                   string    `json:"kind"`
	Name                   string    `json:"name"`
	RoutingNumber          string    `json:"routingNumber"`
	Status                 string    `json:"status"`
	CanReceiveTransactions bool      `json:"canReceiveTransactions"`
	Type                   string    `json:"type"`
	Nickname               string    `json:"nickname"`
	LegalBusinessName      string    `json:"legalBusinessName"`
}

var tracer = otel.Tracer("integrations.mercury")

type mercuryClient struct {
	httpClient *http.Client
}

func New(cfg config.Config) (malak.IntegrationProviderClient, error) {

	return &mercuryClient{
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   cfg.Integration.ClientTimeout,
		},
	}, nil
}

func (m *mercuryClient) Name() malak.IntegrationProvider {
	return malak.IntegrationProviderMercury
}

func (m *mercuryClient) buildRequest(ctx context.Context,
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

func (m *mercuryClient) Ping(
	ctx context.Context, token malak.AccessToken) error {

	req, span, err := m.buildRequest(ctx, token, "connection.ping", "/accounts")
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

	span.SetStatus(codes.Ok, "connection to mercury was successful")
	return nil
}

func (m *mercuryClient) Close() error {
	m.httpClient.CloseIdleConnections()
	return nil
}
