package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"net/http"
	"net/mail"

	"github.com/ayinke-llc/hermes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/theopenlane/utils/passwd"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/ayinke-llc/malak/internal/pkg/jwttoken"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/socialauth"
	"github.com/ayinke-llc/malak/internal/pkg/util"
)

// ENUM(user)
type CookieName string

type authHandler struct {
	googleCfg         socialauth.SocialAuthProvider
	cfg               config.Config
	userRepo          malak.UserRepository
	workspaceRepo     malak.WorkspaceRepository
	tokenManager      jwttoken.JWTokenManager
	queue             queue.QueueHandler
	emailVerification malak.EmailVerificationRepository
	cache             cache.Cache
}

type signupRequest struct {
	GenericRequest

	FullName string      `json:"full_name"`
	Email    malak.Email `json:"email"`
	Password string      `json:"password"`
}

func (s *signupRequest) Validate() error {

	if hermes.IsStringEmpty(s.FullName) {
		return errors.New("please provide your full name")
	}

	if hermes.IsStringEmpty(s.Email.String()) {
		return errors.New("please provide your email address")
	}

	_, err := mail.ParseAddress(s.Email.String())
	if err != nil {
		return errors.New("please provide a valid email address")
	}

	if hermes.IsStringEmpty(s.Password) {
		return errors.New("please provide your password")
	}

	if passwd.Strength(s.Password) < passwd.Moderate {
		return errors.New("your password is too week")
	}

	s.Password, err = malak.HashPassword(s.Password)
	if err != nil {
		return err
	}

	return nil
}

type authenticateUserRequest struct {
	GenericRequest

	Code string `json:"code,omitempty" validate:"required"`
}

func (a *authenticateUserRequest) Validate() error {
	if hermes.IsStringEmpty(a.Code) {
		return errors.New("please provide a valid oauth2 code")
	}

	return nil
}

// @Description Sign in with a social login provider
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

		return a.generateUserToken(user, logger)
	}

	if err != nil {
		logger.Error("an error occurred while creating user", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while creating user"), StatusFailed
	}

	return a.generateUserToken(user, logger)
}

// @Description Fetch current user. This api should also double as a token validation api
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

// @Description Sign up with your email address and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param message body signupRequest true "auth exchange data"
// @Success 200 {object} createdUserResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /auth/register [post]
func (a *authHandler) emailSignup(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating user ( email + password )")

	req := new(signupRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	user := &malak.User{
		Email:    req.Email,
		FullName: req.FullName,
		Metadata: &malak.UserMetadata{},
		Roles:    malak.UserRoles{},
		Password: hermes.Ref(req.Password),
	}

	err := a.userRepo.Create(ctx, user)
	if errors.Is(err, malak.ErrUserExists) {
		return newAPIStatus(http.StatusConflict, "Account already exists. Please use a new email"), StatusFailed
	}

	if err != nil {
		logger.Error("an error occurred while creating user account", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create an account at this time. an error occurred"), StatusFailed
	}

	// ignored on purpose
	_ = a.sendVerificationEmail(user, logger)

	return a.generateUserToken(user, logger)
}

func (a *authHandler) sendVerificationEmail(user *malak.User, logger *zap.Logger) error {

	if user.EmailVerifiedAt != nil {
		return nil
	}

	token, err := malak.NewEmailVerification(user)
	if err != nil {
		// do not return an error here.
		// Let the user request for the email again via the ui
		logger.Error("could not generate email verification token", zap.Error(err))
		return nil
	}

	if err := a.emailVerification.Create(context.Background(), token); err != nil {
		logger.Error("could not store email verification token", zap.Error(err))
		return errors.New("could not store verification token")
	}

	return a.queue.Add(context.Background(), queue.QueueTopicVerifyEmail, queue.EmailVerificationOptions{
		UserID: user.ID,
		Token:  token.Token,
	})
}

func (a *authHandler) generateUserToken(user *malak.User, logger *zap.Logger) (render.Renderer, Status) {

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

// @Description Resend email verification email
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 429 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /user/resend-verification [post]
func (a *authHandler) resendVerificationEmail(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("resending verification email")

	user := getUserFromContext(ctx)
	email := user.Email.String()

	logger = logger.With(zap.String("email", email))

	if user.EmailVerifiedAt != nil {
		return newAPIStatus(http.StatusBadRequest, "email is already verified"), StatusFailed
	}

	cacheKey := fmt.Sprintf("email_verification_attempt:%s", user.ID.String())

	exists, err := a.cache.Exists(ctx, cacheKey)
	if err != nil {
		logger.Error("could not check cache for verification attempt", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while processing your request"), StatusFailed
	}

	if exists {
		return newAPIStatus(http.StatusTooManyRequests, "please wait before requesting another verification email"), StatusFailed
	}

	if err := a.cache.Add(ctx, cacheKey, []byte("1"), 5*time.Minute); err != nil {
		logger.Error("could not add verification attempt to cache", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while processing your request"), StatusFailed
	}

	token, err := malak.NewEmailVerification(user)
	if err != nil {
		logger.Error("could not generate email verification token", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not generate verification token"), StatusFailed
	}

	if err := a.emailVerification.Create(ctx, token); err != nil {
		logger.Error("could not store email verification token", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not store verification token"), StatusFailed
	}

	if err := a.queue.Add(ctx, queue.QueueTopicVerifyEmail, queue.EmailVerificationOptions{
		UserID: user.ID,
		Token:  token.Token,
	}); err != nil {
		logger.Error("could not add verification email to queue", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not send verification email"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "verification email sent successfully"), StatusSuccess
}
