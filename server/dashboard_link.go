package server

import (
	"context"
	"errors"
	"net/http"
	"net/mail"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type generateDashboardLinkRequest struct {
	GenericRequest

	Email malak.Email `json:"email,omitempty" validate:"optional"`
}

func (c *generateDashboardLinkRequest) Validate() error {
	if hermes.IsStringEmpty(c.Email.String()) {
		return nil
	}

	_, err := mail.ParseAddress(c.Email.String())
	if err != nil {
		return errors.New("please provide a valid email address")
	}

	return nil
}

// @Description regenerate the default link for a dashboard
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param message body generateDashboardLinkRequest false "Request body to generate link"
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Success 200 {object} regenerateLinkResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /dashboards/{reference}/access-control/link [post]
func (d *dashboardHandler) generateLink(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("generating a new link")

	req := new(generateDashboardLinkRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

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

	s := d.generator.Token()

	link := &malak.DashboardLink{
		LinkType:    malak.DashboardLinkTypeDefault,
		Reference:   d.generator.Generate(malak.EntityTypeDashboardLink),
		Token:       s,
		DashboardID: dashboard.ID,
	}

	if !hermes.IsStringEmpty(req.Email.String()) {
		link.LinkType = malak.DashboardLinkTypeContact
	}

	opts := &malak.CreateDashboardLinkOptions{
		Link:        link,
		Email:       req.Email,
		WorkspaceID: workspace.ID,
	}

	if err := d.dashboardLinkRepo.Create(ctx, opts); err != nil {
		logger.Error("could not create dashboard link", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could nto create dashboard link"),
			StatusFailed
	}

	if !hermes.IsStringEmpty(req.Email.String()) {
		go func() {
			err := d.queue.Add(context.Background(), queue.QueueTopicShareDashboard, queue.SendEmailOptions{
				Workspace: workspace,
				Recipient: req.Email,
				Token:     s,
			})
			if err != nil {
				logger.Error("could not add item to queue", zap.Error(err))
			}
		}()
	}

	return regenerateLinkResponse{
		APIStatus: newAPIStatus(http.StatusOK, "Link regenerated"),
		Link:      hermes.DeRef(link),
	}, StatusSuccess
}

// @Description fetch public dashboard and charting data points
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Success 200 {object} listDashboardChartsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /public/dashboards/{reference} [get]
func (d *dashboardHandler) publicDashboardDetails(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Fetching public details of dashboard")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	dashboard, err := d.dashboardLinkRepo.PublicDetails(ctx, malak.Reference(ref))
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

	g.Go(func() error {

		var err error

		charts, err = d.dashboardRepo.GetCharts(ctx, malak.FetchDashboardChartsOption{
			WorkspaceID: dashboard.WorkspaceID,
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

	if err := g.Wait(); err != nil {
		return newAPIStatus(http.StatusInternalServerError, err.Error()),
			StatusFailed
	}

	return listDashboardChartsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "dashboards fetched"),
		Dashboard: dashboard,
		Charts:    charts,
		Positions: positions,
	}, StatusSuccess
}

// @Description fetch charting data
// @Tags dashboards
// @Accept  json
// @Produce  json
// @Param reference path string required "dashboard unique reference.. e.g dashboard_"
// @Param chart_reference path string required "chart unique reference.. e.g integration_chart_km31C.e6xV"
// @Success 200 {object} listChartDataPointsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /public/dashboards/{reference}/charts/{chart_reference} [get]
func (d *dashboardHandler) publicChartingDataFetch(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetch charting data")

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	chartRef := chi.URLParam(r, "chart_reference")

	if hermes.IsStringEmpty(chartRef) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	// verify link first
	dashboard, err := d.dashboardLinkRepo.PublicDetails(ctx, malak.Reference(ref))
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
		WorkspaceID: dashboard.WorkspaceID,
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
