package server

import (
	"context"
	"errors"
	"net/http"
	"net/mail"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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

	s, err := hermes.Random(20)
	if err != nil {
		logger.Error("could not generate random token")
		return newAPIStatus(http.StatusInternalServerError, "could not generate link token"), StatusFailed
	}

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

	return regenerateLinkResponse{
		APIStatus: newAPIStatus(http.StatusOK, "Default link regenerated"),
		Link:      hermes.DeRef(link),
	}, StatusSuccess
}
