package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
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

// @Summary Switch current workspace
// @Tags workspace
// @Accept  json
// @Produce  json
// @id switchworkspace
// @Param reference path string required "Workspace unique reference.. e.g update_"
// @Success 200 {object} fetchWorkspaceResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/{reference} [post]
func (wo *workspaceHandler) switchCurrentWorkspaceForUser(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("switching workspaces")

	workspace, err := wo.workspaceRepo.Get(ctx, &malak.FindWorkspaceOptions{
		Reference: malak.Reference(ref),
	})
	if err != nil {
		logger.Error("could not fetch workspace",
			zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not fetch the workspace by reference to switch to"), StatusFailed
	}

	user.Metadata.CurrentWorkspace = workspace.ID

	if err := wo.userRepo.Update(ctx, user); err != nil {
		logger.Error("could not update user's current workspace",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError,
			"could not update user's current workspace"), StatusFailed
	}

	return fetchWorkspaceResponse{
		Workspace: util.DeRef(workspace),
		APIStatus: newAPIStatus(http.StatusOK, "user's default workspace updated"),
	}, StatusSuccess
}
