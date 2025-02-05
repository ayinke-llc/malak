package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/render"
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

	opts := &queue.BillingCreateCustomerOptions{
		Workspace: workspace,
		Email:     user.Email,
	}

	if err := wo.queueClient.Add(ctx, queue.QueueTopicBillingCreateCustomer, opts); err != nil {
		// aware of logic here. no error sent to client as
		// 1. this would rarely fail
		// 2. in the event, it fails, it is not a fatal error and can be sorted
		// out by contacting support usually, or we can even make a slack bot/cli that fixes this
		// 3. given 2, it's fine but if it comes up often, then we should fail
		// the request instead
		logger.Error("an error occurred while adding user to queue to create billing customer",
			zap.Error(err))
	}

	return fetchWorkspaceResponse{
		Workspace: util.DeRef(workspace),
		APIStatus: newAPIStatus(http.StatusCreated, "workspace successfully created"),
	}, StatusSuccess
}
