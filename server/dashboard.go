package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type dashboardHandler struct {
	cfg               config.Config
	dashboardRepo     malak.DashboardRepository
	integrationRepo   malak.IntegrationRepository
	generator         malak.ReferenceGeneratorOperation
	dashboardLinkRepo malak.DashboardLinkRepository
	contactRepo       malak.ContactRepository
	queue             queue.QueueHandler
}

type createDashboardRequest struct {
	GenericRequest

	Title       string `json:"title,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
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

// @Description create a new dashboard
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

// @Description List dashboards
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param page query int false "Page to query data from. Defaults to 1"
// @Param per_page query int false "Number to items to return. Defaults to 10 items"
// @Success 200 {object} listDashboardResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards [get]
func (d *dashboardHandler) list(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Listing dashboards")

	workspace := getWorkspaceFromContext(r.Context())

	opts := malak.ListDashboardOptions{
		Paginator:   malak.PaginatorFromRequest(r),
		WorkspaceID: workspace.ID,
	}

	dashboards, total, err := d.dashboardRepo.List(ctx, opts)
	if err != nil {

		logger.Error("could not list dashboards",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not list dashboards"), StatusFailed
	}

	return listDashboardResponse{
		APIStatus:  newAPIStatus(http.StatusOK, "dashboards fetched"),
		Dashboards: dashboards,
		Meta: meta{
			Paging: pagingInfo{
				PerPage: opts.Paginator.PerPage,
				Page:    opts.Paginator.Page,
				Total:   total,
			},
		},
	}, StatusSuccess
}

// @Description List charts
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Success 200 {object} listIntegrationChartsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/charts [get]
func (d *dashboardHandler) listAllCharts(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Listing all charts")

	workspace := getWorkspaceFromContext(r.Context())

	charts, err := d.integrationRepo.ListCharts(ctx, workspace.ID)
	if err != nil {

		logger.Error("could not list charts",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not list charts"), StatusFailed
	}

	return listIntegrationChartsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "dashboards fetched"),
		Charts:    charts,
	}, StatusSuccess
}

type addChartToDashboardRequest struct {
	GenericRequest

	ChartReference malak.Reference `json:"chart_reference,omitempty" validate:"required"`
}

func (c *addChartToDashboardRequest) Validate() error {
	if hermes.IsStringEmpty(c.ChartReference.String()) {
		return errors.New("please provide a valid chart reference")
	}

	return nil
}

// @Description add a chart to a dashboard
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param message body addChartToDashboardRequest true "dashboard request chart data"
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/{reference}/charts [PUT]
func (d *dashboardHandler) addChart(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("adding a chart to the dashboard")

	workspace := getWorkspaceFromContext(r.Context())

	req := new(addChartToDashboardRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	dashboard, err := d.dashboardRepo.Get(ctx, malak.FetchDashboardOption{
		Reference:   malak.Reference(ref),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		logger.Error("could not fetch dashboard", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching dashboard"

		if errors.Is(err, malak.ErrDashboardNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	chart, err := d.integrationRepo.GetChart(ctx, malak.FetchChartOptions{
		WorkspaceID: workspace.ID,
		Reference:   req.ChartReference,
	})
	if err != nil {
		var status = http.StatusInternalServerError
		var message = "an error occurred while fetching chart"

		logger.Error("could not fetch chart from db",
			zap.Error(err))

		if errors.Is(err, malak.ErrChartNotFound) {
			status = http.StatusNotFound
			message = err.Error()
		}

		return newAPIStatus(status, message), StatusFailed
	}

	dashChart := &malak.DashboardChart{
		Reference:              d.generator.Generate(malak.EntityTypeDashboardChart),
		WorkspaceIntegrationID: chart.WorkspaceIntegrationID,
		WorkspaceID:            workspace.ID,
		DashboardID:            dashboard.ID,
		ChartID:                chart.ID,
	}

	if err := d.dashboardRepo.AddChart(ctx, dashChart); err != nil {
		logger.Error("could not add chart to dashboard", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while adding chart to dashboard"),
			StatusFailed
	}

	return newAPIStatus(http.StatusOK, "chart added to dashboard"),
		StatusSuccess
}

// @Description fetch dashboard
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Success 200 {object} listDashboardChartsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/{reference} [GET]
func (d *dashboardHandler) fetchDashboard(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Fetching dashboard")

	workspace := getWorkspaceFromContext(r.Context())

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	dashboard, err := d.dashboardRepo.Get(ctx, malak.FetchDashboardOption{
		Reference:   malak.Reference(ref),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		logger.Error("could not fetch dashboard", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching dashboard"

		if errors.Is(err, malak.ErrDashboardNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	var g errgroup.Group

	var charts []malak.DashboardChart
	var positions []malak.DashboardChartPosition
	var defaultDashLink malak.DashboardLink

	g.Go(func() error {

		var err error

		charts, err = d.dashboardRepo.GetCharts(ctx, malak.FetchDashboardChartsOption{
			WorkspaceID: workspace.ID,
			DashboardID: dashboard.ID,
		})
		if err != nil {

			logger.Error("could not list dashboard charts",
				zap.Error(err))

			return errors.New("could not list dashboard charts")
		}

		return nil
	})

	g.Go(func() error {

		var err error
		positions, err = d.dashboardRepo.GetDashboardPositions(ctx, dashboard.ID)

		if err != nil {

			logger.Error("could not list dashboard positions",
				zap.Error(err))

			return errors.New("could not list dashboard positions")
		}

		return nil
	})

	g.Go(func() error {
		var err error

		defaultDashLink, err = d.dashboardLinkRepo.DefaultLink(ctx, &dashboard)
		if err != nil {
			logger.Error("could not fetch default dashboard link", zap.Error(err))
			return errors.New("could not fetch dashboard link")
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return newAPIStatus(http.StatusInternalServerError, err.Error()),
			StatusFailed
	}

	return listDashboardChartsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "dashboards fetched"),
		Dashboard: dashboard,
		Charts:    charts,
		Positions: positions,
		Link:      defaultDashLink,
	}, StatusSuccess
}

// @Description fetch charting data
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param reference path string required "chart unique reference.. e.g integration_chart_km31C.e6xV"
// @Success 200 {object} listChartDataPointsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/charts/{reference} [GET]
func (d *dashboardHandler) fetchChartingData(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetch charting data")

	workspace := getWorkspaceFromContext(r.Context())

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	chart, err := d.integrationRepo.GetChart(ctx, malak.FetchChartOptions{
		WorkspaceID: workspace.ID,
		Reference:   malak.Reference(ref),
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

	dataPoints, err := d.integrationRepo.GetDataPoints(ctx, chart)
	if err != nil {

		logger.Error("could not charting data",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not fetch charting data"), StatusFailed
	}

	return listChartDataPointsResponse{
		APIStatus:  newAPIStatus(http.StatusOK, "datapoints fetched"),
		DataPoints: dataPoints,
	}, StatusSuccess
}

type updateDashboardPositionsRequest struct {
	Positions []struct {
		ChartID uuid.UUID `json:"chart_id,omitempty" validate:"required"`
		Index   int64     `json:"index,omitempty" validate:"required"`
	} `json:"positions,omitempty" validate:"required"`
	GenericRequest
}

// @Description update dashboard chart positioning
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param message body updateDashboardPositionsRequest true "dashboard chart positions" @Param reference path string required "dashboard unique reference.. e.g dashboard_22" @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/{reference}/positions [POST]
func (d *dashboardHandler) updateDashboardPositions(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("updating dashboard chart positions")

	workspace := getWorkspaceFromContext(r.Context())

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	req := new(updateDashboardPositionsRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	dashboard, err := d.dashboardRepo.Get(ctx, malak.FetchDashboardOption{
		Reference:   malak.Reference(ref),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		logger.Error("could not fetch dashboard", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching dashboard"

		if errors.Is(err, malak.ErrDashboardNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	// No validation here but leaving it to the db layer since we have relationship
	// mapping guarantees
	positions := make([]malak.DashboardChartPosition, 0, len(req.Positions))

	for _, v := range req.Positions {
		positions = append(positions, malak.DashboardChartPosition{
			OrderIndex:  v.Index,
			ChartID:     v.ChartID,
			DashboardID: dashboard.ID,
		})
	}

	if err := d.dashboardRepo.UpdateDashboardPositions(ctx, dashboard.ID, positions); err != nil {

		logger.Error("could not update dashboard positions",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not update dashboard positions"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "datapoints fetched"), StatusSuccess
}

// @Description remove a chart from a dashboard
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param message body addChartToDashboardRequest true "dashboard request chart data"
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/{reference}/charts [DELETE]
func (d *dashboardHandler) removeChart(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("removing a chart from the dashboard")

	workspace := getWorkspaceFromContext(r.Context())

	req := new(addChartToDashboardRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	dashboard, err := d.dashboardRepo.Get(ctx, malak.FetchDashboardOption{
		Reference:   malak.Reference(ref),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		logger.Error("could not fetch dashboard", zap.Error(err))
		status := http.StatusInternalServerError
		msg := "an error occurred while fetching dashboard"

		if errors.Is(err, malak.ErrDashboardNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	chart, err := d.integrationRepo.GetChart(ctx, malak.FetchChartOptions{
		WorkspaceID: workspace.ID,
		Reference:   req.ChartReference,
	})
	if err != nil {
		var status = http.StatusInternalServerError
		var message = "an error occurred while fetching chart"

		logger.Error("could not fetch chart from db",
			zap.Error(err))

		if errors.Is(err, malak.ErrChartNotFound) {
			status = http.StatusNotFound
			message = err.Error()
		}

		return newAPIStatus(status, message), StatusFailed
	}

	if err := d.dashboardRepo.RemoveChart(ctx, dashboard.ID, chart.ID); err != nil {
		logger.Error("could not remove chart to dashboard", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while removing chart to dashboard"),
			StatusFailed
	}

	return newAPIStatus(http.StatusOK, "chart removed from dashboard"),
		StatusSuccess
}
