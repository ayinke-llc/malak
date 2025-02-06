package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Summary fetch workspace preferences
// @Tags integrations
// @Accept  json
// @Produce  json
// @Success 200 {object} listIntegrationResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations [get]
func (wo *workspaceHandler) getIntegrations(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request,
) (render.Renderer, Status) {

	logger.Debug("fetching workspace integrations")

	workspace := getWorkspaceFromContext(ctx)

	integrations, err := wo.integrationRepo.List(ctx, workspace)
	if err != nil {
		logger.Error("could not list integrations", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not list integrations"), StatusFailed
	}

	return &listIntegrationResponse{
		Integrations: integrations,
		APIStatus:    newAPIStatus(http.StatusOK, "workspace integrations retrieved"),
	}, StatusSuccess
}

type testAPIIntegrationRequest struct {
	APIKey malak.AccessToken `json:"api_key,omitempty" validate:"required"`
	GenericRequest
}

func (t *testAPIIntegrationRequest) Validate() error {

	if util.IsStringEmpty(t.APIKey.String()) {
		return errors.New("please provide API key")
	}

	return nil
}

// @Summary test an api key is valid and can reach the integration
// @Tags integrations
// @Accept  json
// @Produce  json
// @Param message body testAPIIntegrationRequest true "request body to test an integration"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations/{reference}/ping [post]
func (wo *workspaceHandler) pingIntegration(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("pinging integration")

	req := new(testAPIIntegrationRequest)

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

	logger = logger.With(zap.String("integration_name", integration.Integration.IntegrationName))

	if integration.IsActive {
		return newAPIStatus(http.StatusBadRequest, "Integration is currently active"), StatusFailed
	}

	provider, err := malak.ParseIntegrationProvider(strings.ToLower(integration.Integration.IntegrationName))
	if err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	integrationImpl, err := wo.integrationManager.Get(provider)
	if err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	if err := integrationImpl.Ping(ctx, req.APIKey); err != nil {
		logger.Error("could not ping Integration",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, err.Error()), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "integration successfully pinged"),
		StatusSuccess
}

// @Summary enable integration
// @Tags integrations
// @Accept  json
// @Produce  json
// @Param message body testAPIIntegrationRequest true "request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/integrations/{reference} [post]
func (wo *workspaceHandler) enableIntegration(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("enabling integration")

	req := new(testAPIIntegrationRequest)

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

	logger = logger.With(zap.String("integration_name", integration.Integration.IntegrationName))

	provider, err := malak.ParseIntegrationProvider(strings.ToLower(integration.Integration.IntegrationName))
	if err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	integrationImpl, err := wo.integrationManager.Get(provider)
	if err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	if err := integrationImpl.Ping(ctx, req.APIKey); err != nil {
		logger.Error("could not ping Integration",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, err.Error()), StatusFailed
	}

	integration.IsEnabled = true
	integration.Metadata.AccessToken = req.APIKey
	integration.IsActive = true

	if err := wo.integrationRepo.Update(ctx, integration); err != nil {
		logger.Error("could not update integration",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "could not update integration"), StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "integration successfully enabled"),
		StatusSuccess
}
