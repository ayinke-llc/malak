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

	if hermes.IsStringEmpty(c.Title) {
		return errors.New("please provide the title of the dashboard")
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
// @Param message body createDeckRequest true "deck request body"
// @Success 200 {object} uploadImageResponse
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

	logger.Debug("creating a new workspace")

	req := new(createDashboardRequest)

	workspace := getWorkspaceFromContext(ctx)

	user := getUserFromContext(ctx)

	_, _ = workspace, user

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "deck was uploaded"), StatusSuccess
}
