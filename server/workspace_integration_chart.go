package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type createChartRequest struct {
	Title     string                         `json:"title,omitempty" validate:"required"`
	ChartType malak.IntegrationChartType     `json:"chart_type,omitempty" validate:"required"`
	Datapoint malak.IntegrationDataPointType `json:"datapoint,omitempty" validate:"required"`

	GenericRequest
}

func (t *createChartRequest) Validate() error {

	if hermes.IsStringEmpty(t.Title) {
		return errors.New("please provide chart title")
	}

	if len(t.Title) > 100 {
		return errors.New("title can not be more than 100 characters")
	}

	if len(t.Title) < 5 {
		return errors.New("title can not be less than 5 characters")
	}

	if t.ChartType != malak.IntegrationChartTypeBar {
		return errors.New("unsupported chart type")
	}

	if !t.Datapoint.IsValid() {
		return errors.New("datapoint type not supported")
	}

	return nil
}

// @Description create chart
// @Tags integrations
// @Accept  json
// @Produce  json
// @Param message body createChartRequest true "request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations/{reference}/charts [post]
func (wo *workspaceHandler) createChart(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("creating chart")

	req := new(createChartRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	logger = logger.With(zap.String("integration_reference", ref))

	integration, err := wo.integrationRepo.Get(ctx, malak.FindWorkspaceIntegrationOptions{
		Reference: malak.Reference(ref),
	})
	if err != nil {
		var msg string = "could not fetch integration"
		var status = http.StatusInternalServerError

		if errors.Is(err, malak.ErrWorkspaceIntegrationNotFound) {
			msg = err.Error()
			status = http.StatusNotFound
		}

		logger.
			Error(msg,
				zap.Error(err))
		return newAPIStatus(status, msg), StatusFailed
	}

	if !integration.Integration.IsEnabled {
		return newAPIStatus(http.StatusBadRequest, "integration not enabled yet and coming soon"), StatusFailed
	}

	if !integration.IsEnabled || !integration.IsActive {
		return newAPIStatus(http.StatusBadRequest, "integration not enabled"), StatusFailed
	}

	if integration.Integration.IntegrationType != malak.IntegrationTypeSystem {
		return newAPIStatus(http.StatusBadRequest,
				"you can only add a chart for a system integration. Other integrations have their chart auto added"),
			StatusFailed
	}

	logger = logger.With(zap.String("integration_name", integration.Integration.IntegrationName))

	chart := malak.IntegrationChartValues{
		InternalName:   malak.IntegrationChartInternalNameType(wo.referenceGenerator.ShortLink()),
		UserFacingName: req.Title,
		ChartType:      req.ChartType,
		DataPointType:  req.Datapoint,
	}

	if err := wo.integrationRepo.CreateCharts(ctx, integration, []malak.IntegrationChartValues{chart}); err != nil {
		logger.Error("could not create chart",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "could not create chart"), StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "created chart for integration"),
		StatusSuccess
}

type addDataPointRequest struct {
	Value int64 `validate:"required" json:"value,omitempty"`
	GenericRequest
}

func (t *addDataPointRequest) Validate() error {

	if t.Value < 0 {
		return errors.New("provide a valid value")
	}

	return nil
}

// @Description add data point values to chart
// @Tags integrations
// @Accept  json
// @Produce  json
// @Param message body addDataPointRequest true "request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations/{reference}/charts/{chart_reference}/points [post]
func (wo *workspaceHandler) addDataPoint(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	chartRef := chi.URLParam(r, "chart_reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("add data point values to chart")

	req := new(addDataPointRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	logger = logger.With(zap.String("integration_reference", ref))

	integration, err := wo.integrationRepo.Get(ctx, malak.FindWorkspaceIntegrationOptions{
		Reference: malak.Reference(ref),
	})
	if err != nil {
		var msg string = "could not fetch integration"
		var status = http.StatusInternalServerError

		if errors.Is(err, malak.ErrWorkspaceIntegrationNotFound) {
			msg = err.Error()
			status = http.StatusNotFound
		}

		logger.
			Error(msg,
				zap.Error(err))
		return newAPIStatus(status, msg), StatusFailed
	}

	if !integration.Integration.IsEnabled {
		return newAPIStatus(http.StatusBadRequest, "integration not enabled yet and coming soon"), StatusFailed
	}

	if !integration.IsEnabled || !integration.IsActive {
		return newAPIStatus(http.StatusBadRequest, "integration not enabled"), StatusFailed
	}

	if integration.Integration.IntegrationType != malak.IntegrationTypeSystem {
		return newAPIStatus(http.StatusBadRequest,
				"you can only manually add a datapoint for a system integration"),
			StatusFailed
	}

	logger = logger.With(zap.String("integration_name", integration.Integration.IntegrationName))

	chart, err := wo.integrationRepo.GetChart(ctx, malak.FetchChartOptions{
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
		Reference:   malak.Reference(chartRef),
	})
	if err != nil {
		logger.Error("could not fetch chart", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching chart"

		if errors.Is(err, malak.ErrChartNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	value := malak.IntegrationDataValues{
		InternalName:   chart.InternalName,
		UserFacingName: chart.UserFacingName,
		Data: malak.IntegrationDataPoint{
			WorkspaceID:            getWorkspaceFromContext(ctx).ID,
			WorkspaceIntegrationID: integration.ID,
			IntegrationChartID:     chart.ID,
			Reference:              wo.referenceGenerator.Generate(malak.EntityTypeIntegrationDatapoint),
			PointName:              malak.GetTodayFormatted(),
			PointValue:             req.Value,
			Metadata:               malak.IntegrationDataPointMetadata{},
		},
	}

	if err := wo.integrationRepo.AddDataPoint(ctx, integration, []malak.IntegrationDataValues{value}); err != nil {
		logger.Error("could not insert data points",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "could not create data points"), StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "added datapoints"),
		StatusSuccess
}
