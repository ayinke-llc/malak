package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type workspaceHandler struct {
	cfg                     config.Config
	userRepo                malak.UserRepository
	workspaceRepo           malak.WorkspaceRepository
	planRepo                malak.PlanRepository
	preferenceRepo          malak.PreferenceRepository
	integrationRepo         malak.IntegrationRepository
	referenceGenerationFunc func(e malak.EntityType) string
	stripeClient            *client.API
	queueClient             queue.QueueHandler
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
		logger.Error("an error occurred while adding user to queue to create stripe customer",
			zap.Error(err))
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

	if !user.CanAccessWorkspace(workspace.ID) {
		return newAPIStatus(http.StatusForbidden,
			"You can only join workspaces you are a member of"), StatusFailed
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

type updateWorkspaceRequest struct {
	Timezone      *string `json:"timezone,omitempty"`
	WorkspaceName *string `json:"workspace_name,omitempty"`
	Website       *string `json:"website,omitempty"`
	Logo          *string `json:"logo,omitempty"`

	GenericRequest
}

func (u *updateWorkspaceRequest) Validate() error {

	timezone := strings.TrimSpace(hermes.DeRef(u.Timezone))
	website := strings.TrimSpace(hermes.DeRef(u.Website))
	workspaceName := strings.TrimSpace(hermes.DeRef(u.WorkspaceName))
	logo := strings.TrimSpace(hermes.DeRef(u.Logo))

	if !hermes.IsStringEmpty(timezone) {
		_, err := time.LoadLocation(timezone)
		if err != nil {
			return errors.New("invalid or unsupported timezone")
		}
	}

	if !hermes.IsStringEmpty(logo) {
		isValid, err := malak.IsImageFromURL(logo)
		if err != nil {
			return err
		}

		if !isValid {
			return errors.New("logo is not a valid image url")
		}
	}

	if !hermes.IsStringEmpty(website) {
		_, err := url.Parse(website)
		if err != nil {
			return errors.New("invalid website")
		}
	}

	if !hermes.IsStringEmpty(workspaceName) {

		if len(workspaceName) < 5 {
			return errors.New("workspace name must be a minimum of 5 characters")
		}
	}

	return nil
}

// @Summary update workspace details
// @Tags workspace
// @Accept  json
// @Produce  json
// @Param message body updateWorkspaceRequest true "request body to create a workspace"
// @Success 200 {object} fetchWorkspaceResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces [patch]
func (wo *workspaceHandler) updateWorkspace(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("updating workspace")

	req := new(updateWorkspaceRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	workspace := getWorkspaceFromContext(ctx)

	timezone := strings.TrimSpace(hermes.DeRef(req.Timezone))
	website := strings.TrimSpace(hermes.DeRef(req.Website))
	workspaceName := strings.TrimSpace(hermes.DeRef(req.WorkspaceName))
	logo := strings.TrimSpace(hermes.DeRef(req.Logo))

	if !hermes.IsStringEmpty(timezone) {
		workspace.Timezone = timezone
	}

	if !hermes.IsStringEmpty(website) {
		workspace.Website = website
	}

	if !hermes.IsStringEmpty(workspaceName) {
		workspace.WorkspaceName = workspaceName
	}

	if !hermes.IsStringEmpty(logo) {
		workspace.LogoURL = logo
	}

	if err := wo.workspaceRepo.Update(ctx, workspace); err != nil {
		logger.Error("could not update workspace",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError,
			"could not update workspace"), StatusFailed
	}

	return fetchWorkspaceResponse{
		Workspace: util.DeRef(workspace),
		APIStatus: newAPIStatus(http.StatusOK, "workspace updated"),
	}, StatusSuccess
}

// @Summary fetch workspace preferences
// @Tags workspace
// @Accept  json
// @Produce  json
// @Success 200 {object} preferenceResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/preferences [get]
func (wo *workspaceHandler) getPreferences(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request,
) (render.Renderer, Status) {

	logger.Debug("fetching workspace preferences")

	workspace := getWorkspaceFromContext(ctx)

	preferences, err := wo.preferenceRepo.Get(ctx, workspace)
	if err != nil {
		logger.Error("could not fetch preferences", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not fetch preferences"), StatusFailed
	}

	return &preferenceResponse{
		Preferences: preferences,
		APIStatus:   newAPIStatus(http.StatusOK, "workspace preferences retrieved"),
	}, StatusSuccess
}

type updatePreferencesRequest struct {
	Preferences struct {
		Billing    malak.BillingPreferences       `json:"billing,omitempty" validate:"required"`
		Newsletter malak.CommunicationPreferences `json:"newsletter,omitempty" validate:"required"`
	} `json:"preferences,omitempty" validate:"required"`
	GenericRequest
}

func (u *updatePreferencesRequest) Validate() error { return nil }

func (u *updatePreferencesRequest) Make(current *malak.Preference) *malak.Preference {

	if u.Preferences.Newsletter.EnableMarketing != current.Communication.EnableMarketing {
		current.Communication.EnableMarketing = u.Preferences.Newsletter.EnableMarketing
	}

	if u.Preferences.Newsletter.EnableProductUpdates != current.Communication.EnableProductUpdates {
		current.Communication.EnableProductUpdates = u.Preferences.Newsletter.EnableProductUpdates
	}

	return current
}

// @Summary update workspace preferences
// @Tags workspace
// @Accept  json
// @Produce  json
// @Param message body updatePreferencesRequest true "request body to updare a workspace preference"
// @Success 200 {object} preferenceResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/preferences [put]
func (wo *workspaceHandler) updatePreferences(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Updating workspace preferences")

	req := new(updatePreferencesRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	workspace := getWorkspaceFromContext(ctx)

	preferences, err := wo.preferenceRepo.Get(ctx, workspace)
	if err != nil {
		logger.Error("could not fetch preferences", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch preferences"), StatusFailed
	}

	pref := req.Make(preferences)

	if err := wo.preferenceRepo.Update(ctx, pref); err != nil {
		logger.Error("could not update preferences", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not update preferences"), StatusFailed
	}

	return &preferenceResponse{
		Preferences: preferences,
		APIStatus:   newAPIStatus(http.StatusOK, "workspace preferences updated"),
	}, StatusSuccess
}

// @Summary get billing portal
// @Tags billing
// @Accept  json
// @Produce  json
// @Success 200 {object} fetchBillingPortalResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/billing [post]
func (wo *workspaceHandler) getBillingPortal(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetching billing portal")

	workspace := getWorkspaceFromContext(ctx)

	billingSession, err := wo.stripeClient.BillingPortalSessions.New(&stripe.BillingPortalSessionParams{
		Customer: &workspace.StripeCustomerID,
	})
	if err != nil {
		logger.Error("could not create billing portal", zap.Error(err))
		return newAPIStatus(http.StatusFailedDependency, "could not create billing portal link"),
			StatusFailed
	}

	return &fetchBillingPortalResponse{
		APIStatus: newAPIStatus(http.StatusOK, "Billing portal link created"),
		Link:      billingSession.URL,
	}, StatusSuccess
}
