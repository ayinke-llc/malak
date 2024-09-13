package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type workspaceHandler struct {
	cfg                     config.Config
	userRepo                malak.UserRepository
	workspaceRepo           malak.WorkspaceRepository
	planRepo                malak.PlanRepository
	referenceGenerationFunc func(e malak.EntityType) string
}

type createWorkspaceRequest struct {
	Name string `json:"name,omitempty"`
	GenericRequest
}

func (c *createWorkspaceRequest) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if util.IsStringEmpty(c.Name) {
		return errors.New("please provide workspace name")
	}

	if len(c.Name) < 5 {
		return errors.New("workspace name must be a minimum of 5 characters")
	}

	return nil
}

// @Summary Create a new workspace
// @Tags workspace
// @Accept  json
// @Produce  json
// @Param message body createWorkspaceRequest true "request body to create a workspace"
// @Success 200 {object} fetchWorkspaceResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces [post]
func (wo *workspaceHandler) createWorkspace(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	logger.Debug("creating workspace")

	req := new(createWorkspaceRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	plan, err := wo.planRepo.Get(ctx, &malak.FetchPlanOptions{
		Reference: wo.cfg.Billing.DefaultPlanReference,
	})
	if err != nil {
		logger.
			Error("could not fetch default plan",
				zap.Error(err),
				zap.String("plan_reference", wo.cfg.Billing.DefaultPlanReference))
		return newAPIStatus(http.StatusInternalServerError,
			"could not fetch default plan details"), StatusFailed
	}

	workspace := malak.NewWorkspace(req.Name, user, plan,
		wo.referenceGenerationFunc(malak.EntityTypeWorkspace))

	err = wo.workspaceRepo.Create(ctx, &malak.CreateWorkspaceOptions{
		User:      user,
		Workspace: workspace,
	})
	if err != nil {
		logger.Error("could not fetch default plan",
			zap.Error(err),
			zap.String("plan_reference", wo.cfg.Billing.DefaultPlanReference))
		return newAPIStatus(http.StatusInternalServerError,
			"could not create workspace"), StatusFailed
	}

	return fetchWorkspaceResponse{
		Workspace: util.DeRef(workspace),
		APIStatus: newAPIStatus(http.StatusCreated, "workspace successfully created"),
	}, StatusSuccess
}
