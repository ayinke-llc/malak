package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type dashboardHandler struct {
	cfg           config.Config
	dashboardRepo malak.DashboardRepository
	generator     malak.ReferenceGeneratorOperation
}

type createDashboardRequest struct {
	GenericRequest

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *createDashboardRequest) Validate() error {
	if hermes.IsStringEmpty(c.Description) {
		return errors.New("please provide a description")
	}

	if len(c.Description) > 500 {
		return errors.New("description cannot be more than 500 characters")
	}

	if hermes.IsStringEmpty(c.Title) {
		return errors.New("please provide the title of the dashboard")
	}

	if len(c.Title) > 100 {
		return errors.New("title cannot be more than 100 characters")
	}

	p := bluemonday.StrictPolicy()

	c.Title = p.Sanitize(c.Title)
	c.Description = p.Sanitize(c.Description)

	return nil
}

// @Summary create a new dashboard
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param message body createDashboardRequest true "dashboard request body"
// @Success 200 {object} fetchDashboardResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards [post]
func (d *dashboardHandler) create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating a new dashboard")

	req := new(createDashboardRequest)

	workspace := getWorkspaceFromContext(ctx)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	dashboard := &malak.Dashboard{
		Description: req.Description,
		Title:       req.Title,
		Reference:   d.generator.Generate(malak.EntityTypeDashboard),
		ChartCount:  0,
		WorkspaceID: workspace.ID,
	}

	if err := d.dashboardRepo.Create(ctx, dashboard); err != nil {
		logger.Error("could not create dashboard", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create dashboard"),
			StatusFailed
	}

	return fetchDashboardResponse{
		APIStatus: newAPIStatus(http.StatusOK, "dashboard was successfully created"),
		Dashboard: hermes.DeRef(dashboard),
	}, StatusSuccess
}
