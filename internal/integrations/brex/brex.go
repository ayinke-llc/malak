package brex

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

type AccountTransaction struct {
	Total int64 `json:"total"`
}

type Account struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	CurrentBalance   Balance `json:"current_balance"`
	AvailableBalance Balance `json:"available_balance"`
	AccountNumber    string  `json:"account_number"`
	RoutingNumber    string  `json:"routing_number"`
	Primary          bool    `json:"primary"`
}

type Balance struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type AccountsResponse struct {
	NextCursor string    `json:"next_cursor"`
	Items      []Account `json:"items"`
}

type TransactionAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Transaction struct {
	ID              string            `json:"id"`
	Description     string            `json:"description"`
	Amount          TransactionAmount `json:"amount"`
	InitiatedAtDate string            `json:"initiated_at_date"`
	PostedAtDate    string            `json:"posted_at_date"`
	Type            string            `json:"type"`
	TransferID      string            `json:"transfer_id"`
}

type TransactionsResponse struct {
	NextCursor string        `json:"next_cursor"`
	Items      []Transaction `json:"items"`
}

var tracer = otel.Tracer("integrations.brex")

type brexClient struct {
	httpClient *http.Client
}

func New(cfg config.Config) (malak.IntegrationProviderClient, error) {

	transport := otelhttp.NewTransport(http.DefaultTransport)
	client := &http.Client{
		Transport: transport,
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

func checkResponse(res *http.Response, span trace.Span) error {
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("brex api request failed with status code: %d", res.StatusCode)
		span.SetAttributes(attribute.Int("response_code", res.StatusCode))
		span.SetStatus(codes.Error, "request failed with non-200 status code")
		return err
	}
	return nil
}

func (m *brexClient) Ping(
	ctx context.Context,
	token malak.AccessToken) ([]malak.IntegrationChartValues, error) {

	charts := make([]malak.IntegrationChartValues, 0)

	req, span, err := m.buildRequest(ctx, token, "connection.ping", "/users/me")
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

	if err := checkResponse(res, span); err != nil {
		return charts, err
	}

	span.SetStatus(codes.Ok, "connection to brex was successful")

	var user CurrentUser

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not decode response")
		return charts, err
	}

	// Now fetch accounts
	req, span, err = m.buildRequest(ctx, token, "accounts.fetch", "/accounts/cash")
	if err != nil {
		return charts, err
	}

	res, err = m.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not send request")
		return charts, err
	}

	defer res.Body.Close()

	if err := checkResponse(res, span); err != nil {
		return charts, err
	}

	var accountsResponse AccountsResponse

	if err := json.NewDecoder(res.Body).Decode(&accountsResponse); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not decode response")
		return charts, err
	}

	for _, account := range accountsResponse.Items {
		charts = append(charts, malak.IntegrationChartValues{
			InternalName:   malak.IntegrationChartInternalNameTypeBrexAccount,
			UserFacingName: account.Name,
			ProviderID:     account.ID,
			ChartType:      malak.IntegrationChartTypeBar,
		})

		charts = append(charts, malak.IntegrationChartValues{
			InternalName:   malak.IntegrationChartInternalNameTypeBrexAccount,
			UserFacingName: "Transactions count for " + account.Name,
			ChartType:      malak.IntegrationChartTypeBar,
		})
	}

	return charts, nil
}

func (m *brexClient) Close() error {
	m.httpClient.CloseIdleConnections()
	return nil
}

func (m *brexClient) fetchTransactions(ctx context.Context,
	token malak.AccessToken,
	accountID string,
	startTime string,
	cursor string) ([]Transaction, error) {

	var allTransactions []Transaction
	currentCursor := cursor

	for {
		endpoint := fmt.Sprintf("/transactions/cash/%s?posted_at_start=%s", accountID, startTime)
		if currentCursor != "" {
			endpoint = fmt.Sprintf("%s&cursor=%s", endpoint, currentCursor)
		}

		req, span, err := m.buildRequest(ctx, token, "account.transactions.fetch", endpoint)
		if err != nil {
			return nil, err
		}

		res, err := m.httpClient.Do(req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "could not send request")
			return nil, err
		}

		defer res.Body.Close()

		if err := checkResponse(res, span); err != nil {
			return nil, err
		}

		var txs TransactionsResponse

		if err := json.NewDecoder(res.Body).Decode(&txs); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "could not decode tx response")
			return nil, err
		}

		allTransactions = append(allTransactions, txs.Items...)

		if txs.NextCursor == "" {
			break
		}

		currentCursor = txs.NextCursor
	}

	return allTransactions, nil
}

func (m *brexClient) Data(ctx context.Context,
	token malak.AccessToken,
	opts *malak.IntegrationFetchDataOptions) ([]malak.IntegrationDataValues, error) {

	var g errgroup.Group
	var dataPoints = make([]malak.IntegrationDataValues, 0, 2)

	req, span, err := m.buildRequest(ctx, token, "accounts.fetch", "/accounts/cash")
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

	if err := checkResponse(res, span); err != nil {
		return dataPoints, err
	}

	var accountsResponse AccountsResponse

	if err := json.NewDecoder(res.Body).Decode(&accountsResponse); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "could not decode response")
		return dataPoints, err
	}

	for _, account := range accountsResponse.Items {
		dataPoints = append(dataPoints, malak.IntegrationDataValues{
			InternalName:   malak.IntegrationChartInternalNameTypeBrexAccount,
			ProviderID:     account.ID,
			UserFacingName: account.Name,
			Data: malak.IntegrationDataPoint{
				DataPointType:          malak.IntegrationDataPointTypeCurrency,
				WorkspaceIntegrationID: opts.IntegrationID,
				WorkspaceID:            opts.WorkspaceID,
				Reference:              opts.ReferenceGenerator.Generate(malak.EntityTypeIntegrationDatapoint),
				PointName:              malak.GetTodayFormatted(),
				PointValue:             int64(math.Floor(account.AvailableBalance.Amount * 100)),
				Metadata:               malak.IntegrationDataPointMetadata{},
			},
		})

		g.Go(func() error {
			startTimeFormatted := opts.LastFetchedAt.Format(time.RFC3339)

			transactions, err := m.fetchTransactions(ctx, token, account.ID, startTimeFormatted, "")
			if err != nil {
				return err
			}

			dataPoints = append(dataPoints, malak.IntegrationDataValues{
				InternalName:   malak.IntegrationChartInternalNameTypeBrexAccountTransaction,
				UserFacingName: "Transactions count for " + account.Name,
				Data: malak.IntegrationDataPoint{
					DataPointType:          malak.IntegrationDataPointTypeOthers,
					WorkspaceIntegrationID: opts.IntegrationID,
					WorkspaceID:            opts.WorkspaceID,
					Reference:              opts.ReferenceGenerator.Generate(malak.EntityTypeIntegrationDatapoint),
					PointName:              malak.GetTodayFormatted(),
					PointValue:             int64(len(transactions)),
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
