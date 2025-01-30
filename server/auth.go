package server

import (
	"context"
	"errors"

	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stripe/stripe-go/v81/client"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ENUM(user)
type CookieName string

type authHandler struct {
	googleCfg     socialauth.SocialAuthProvider
	cfg           config.Config
	userRepo      malak.UserRepository
	workspaceRepo malak.WorkspaceRepository
	tokenManager  jwttoken.JWTokenManager
	stripeClient  *client.API
}

type authenticateUserRequest struct {
	GenericRequest

	Code string `json:"code,omitempty" validate:"required"`
}

func (a *authenticateUserRequest) Validate() error {
	if util.IsStringEmpty(a.Code) {
		return errors.New("please provide a valid oauth2 code")
	}

	return nil
}

// @Summary Sign in with a social login provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param message body authenticateUserRequest true "auth exchange data"
// @Success 200 {object} createdUserResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /auth/connect/{provider} [post]
// @Param provider  path string true "oauth2 provider"
func (a *authHandler) Login(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	provider := chi.URLParam(r, "provider")

	logger = logger.With(zap.String("provider", provider))

	span.SetAttributes(attribute.String("auth_provider", provider))

	logger.Debug("Authenticating user")

	if provider != "google" {
		return newAPIStatus(http.StatusBadRequest, "unspported provider"), StatusFailed
	}

	req := new(authenticateUserRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	token, err := a.googleCfg.Validate(ctx, socialauth.ValidateOptions{
		Code: req.Code,
	})
	if err != nil {
		logger.Error("could not exchange token", zap.Error(err))
		return newAPIStatus(http.StatusBadRequest, "could not verify your sign in with Google"), StatusFailed
	}

	u, err := a.googleCfg.User(ctx, token)
	if err != nil {
		logger.Error("could not fetch user details from google", zap.Error(err))
		return newAPIStatus(http.StatusBadRequest, "could not fetch user details from oauth2 provider"), StatusFailed
	}

	user := &malak.User{
		Email:    malak.Email(u.Email),
		FullName: u.Name,
		Metadata: &malak.UserMetadata{},
		Roles:    malak.UserRoles{},
	}

	err = a.userRepo.Create(ctx, user)
	if errors.Is(err, malak.ErrUserExists) {
		// if user exists
		// fetch the user by email,
		// and generate the token
		user, err := a.userRepo.Get(ctx, &malak.FindUserOptions{
			Email: user.Email,
		})
		if err != nil {
			logger.Error("an error occurred while fetching user", zap.Error(err))
			return newAPIStatus(http.StatusInternalServerError, "an error occurred while logging user into app"), StatusFailed
		}

		token, err := a.tokenManager.GenerateJWToken(jwttoken.JWTokenData{
			UserID: user.ID,
		})
		if err != nil {
			logger.Error("an error occurred while generating jwt token", zap.Error(err))
			return newAPIStatus(http.StatusInternalServerError, "an error occurred while generating jwt token"), StatusFailed
		}

		resp := createdUserResponse{
			User:      util.DeRef(user),
			APIStatus: newAPIStatus(http.StatusOK, "Logged in Successfully"),
			Token:     token.Token,
		}
		return resp, StatusSuccess
	}

	if err != nil {
		logger.Error("an error occurred while creating user", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while creating user"), StatusFailed
	}

	authToken, err := a.tokenManager.GenerateJWToken(jwttoken.JWTokenData{
		UserID: user.ID,
	})
	if err != nil {
		logger.Error("an error occurred while generating jwt token", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while generating jwt token"), StatusFailed
	}

	resp := createdUserResponse{
		User:      util.DeRef(user),
		APIStatus: newAPIStatus(http.StatusOK, "user Successfully created"),
		Token:     authToken.Token,
	}
	return resp, StatusSuccess
}

// @Summary Fetch current user. This api should also double as a token validation api
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} createdUserResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /user [get]
func (a *authHandler) fetchCurrentUser(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Fetching user profile")

	user := getUserFromContext(ctx)

	var workspace *malak.Workspace = nil
	if doesWorkspaceExistInContext(ctx) {
		workspace = getWorkspaceFromContext(ctx)
	}

	workspaces, err := a.workspaceRepo.List(ctx, user)
	if err != nil {
		logger.Error("could not list workspaces", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError, "could not list all your current workspaces"), StatusFailed
	}

	return createdUserResponse{
		User:             util.DeRef(user),
		CurrentWorkspace: workspace,
		Workspaces:       workspaces,
		APIStatus:        newAPIStatus(http.StatusOK, "user data successfully retrieved"),
	}, StatusSuccess
}
