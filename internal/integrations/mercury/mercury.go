package mercury

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
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

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

type AccountTransaction struct {
	Total int64 `json:"total"`
}

var tracer = otel.Tracer("integrations.mercury")

type mercuryClient struct {
	httpClient *http.Client
}

func New(cfg config.Config) (malak.IntegrationProviderClient, error) {

	transport := otelhttp.NewTransport(http.DefaultTransport)
	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Secrets.ClientTimeout,
	}

	return &mercuryClient{
		httpClient: client,
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
	ctx context.Context,
	token malak.AccessToken) ([]malak.IntegrationChartValues, error) {

	charts := make([]malak.IntegrationChartValues, 0)

	req, span, err := m.buildRequest(ctx, token, "connection.ping", "/accounts")
	if err != nil {
		return charts, err
	}

	res, err := m.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not send request")
		return charts, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("mercury api request failed with status code: %d", res.StatusCode)
		span.SetAttributes(attribute.Int("response_code", res.StatusCode))
		return charts, err
	}

	span.SetStatus(codes.Ok, "connection to mercury was successful")

	var response AccountsResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not decode response")
		return charts, err
	}

	for _, account := range response.Accounts {
		charts = append(charts, malak.IntegrationChartValues{
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccount,
			UserFacingName: account.Name,
			ProviderID:     account.ID,
		})

		charts = append(charts, malak.IntegrationChartValues{
			InternalName:   malak.IntegrationChartInternalNameTypeMercuryAccountTransaction,
			UserFacingName: "Transactions count for " + account.Name,
		})
	}

	return charts, nil
}

func (m *mercuryClient) Close() error {
	m.httpClient.CloseIdleConnections()
	return nil
}

func (m *mercuryClient) Data(ctx context.Context,
	token malak.AccessToken,
	opts *malak.IntegrationFetchDataOptions) ([]malak.IntegrationDataValues, error) {

	var g errgroup.Group
	var dataPoints = make([]malak.IntegrationDataValues, 0, 2)

	req, span, err := m.buildRequest(ctx, token, "accounts.fetch", "/accounts")
	if err != nil {
		return dataPoints, err
	}

	res, err := m.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not send request")
		return dataPoints, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("mercury api request failed with status code: %d", res.StatusCode)
		span.SetAttributes(attribute.Int("response_code", res.StatusCode))
		return dataPoints, err
	}

	var response AccountsResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not decode response")
		return dataPoints, err
	}

	for _, account := range response.Accounts {
		dataPoints = append(dataPoints, malak.IntegrationDataValues{
			InternalName: malak.IntegrationChartInternalNameTypeMercuryAccount,
			ProviderID:   account.ID,
			Data: malak.IntegrationDataPoint{
				DataPointType:          malak.IntegrationDataPointTypeCurrency,
				WorkspaceIntegrationID: opts.IntegrationID,
				WorkspaceID:            opts.WorkspaceID,
				Reference:              opts.ReferenceGenerator.Generate(malak.EntityTypeIntegrationDatapoint),
				PointName:              malak.GetTodayFormatted(),
				PointValue:             int64(math.Floor(account.AvailableBalance * 100)),
				Metadata:               malak.IntegrationDataPointMetadata{},
			},
		})

		g.Go(func() error {

			dateFormatterd := time.Now().Format("2006-01-02")

			req, span, err := m.buildRequest(ctx, token, "account.transactions.fetch",
				fmt.Sprintf("/account/%s/transactions?start=%s&end=%s&status=sent", account.ID, dateFormatterd, dateFormatterd))
			if err != nil {
				return err
			}

			span.SetAttributes(
				attribute.String("workspace_id", opts.WorkspaceID.String()),
				attribute.String("account_id", account.ID))

			res, err := m.httpClient.Do(req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "could not send request")
				return err
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				err = fmt.Errorf("mercury api request failed with status code: %d", res.StatusCode)
				span.SetAttributes(attribute.Int("response_code", res.StatusCode))
				span.SetStatus(codes.Error, "request failed")
				return err
			}

			var txs AccountTransaction

			if err := json.NewDecoder(res.Body).Decode(&txs); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "could not decode tx response")
				return err
			}

			dataPoints = append(dataPoints, malak.IntegrationDataValues{
				InternalName: malak.IntegrationChartInternalNameTypeMercuryAccountTransaction,
				Data: malak.IntegrationDataPoint{
					DataPointType:          malak.IntegrationDataPointTypeOthers,
					WorkspaceIntegrationID: opts.IntegrationID,
					WorkspaceID:            opts.WorkspaceID,
					Reference:              opts.ReferenceGenerator.Generate(malak.EntityTypeIntegrationDatapoint),
					PointName:              malak.GetTodayFormatted(),
					PointValue:             txs.Total,
					Metadata:               malak.IntegrationDataPointMetadata{},
				},
			})

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not fetch account and transactions")
		return dataPoints, err
	}

	return dataPoints, nil
}
