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
	GenericRequest
	Title     string                     `json:"title,omitempty" validate:"required"`
	ChartType malak.IntegrationChartType `json:"chart_type,omitempty" validate:"required"`
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

	if !t.ChartType.IsValid() {
		return errors.New("unsupported chart type")
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
	}

	if err := wo.integrationRepo.CreateCharts(ctx, integration, []malak.IntegrationChartValues{chart}); err != nil {
		logger.Error("could not create chart",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "could not create chart"), StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "created chart for integration"),
		StatusSuccess
}
