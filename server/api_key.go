package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type apiKeyHandler struct {
	apiRepo       malak.APIKeyRepository
	secretsClient secret.SecretClient
	generator     malak.ReferenceGeneratorOperation
}

type createAPIKeyRequest struct {
	GenericRequest

	Title string `json:"title,omitempty" validate:"required"`
}

func (c *createAPIKeyRequest) Validate() error {
	if hermes.IsStringEmpty(c.Title) {
		return errors.New("please provide the title of the api key")
	}

	if len(c.Title) < 3 {
		return errors.New("title must be more than 3 characters")
	}

	if len(c.Title) > 20 {
		return errors.New("title must be more than 20 characters")
	}

	p := bluemonday.StrictPolicy()

	c.Title = p.Sanitize(c.Title)

	return nil
}

// @Description Creates a new api key
// @Tags developers
// @Accept  json
// @Produce  json
// @Param message body createAPIKeyRequest true "api key request body"
// @Success 200 {object} createdAPIKeyResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /developers/keys [post]
func (d *apiKeyHandler) create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating api key")

	req := new(createAPIKeyRequest)

	workspace := getWorkspaceFromContext(r.Context())

	user := getUserFromContext(r.Context())

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	value := d.generator.Token()

	encrypted, err := d.secretsClient.Create(ctx, &secret.CreateSecretOptions{
		Value:       value,
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		logger.Error("could not encrypt value", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not encrypt value"),
			StatusFailed
	}

	key := &malak.APIKey{
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Reference:   d.generator.Generate(malak.EntityTypeApiKey),
		Value:       encrypted,
		KeyName:     req.Title,
	}

	if err := d.apiRepo.Create(ctx, key); err != nil {
		logger.Error("could not create api key",
			zap.Error(err))

		var status = http.StatusInternalServerError
		var msg = "could not create api key"

		if errors.Is(err, malak.ErrAPIKeyMaxLimit) {
			status = http.StatusBadRequest
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	return createdAPIKeyResponse{
		APIStatus: newAPIStatus(http.StatusOK, "api key created"),
		Value:     value,
	}, StatusSuccess
}

// @Description list api keys
// @Tags developers
// @Accept  json
// @Produce  json
// @Success 200 {object} listAPIKeysResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /developers/keys [get]
func (d *apiKeyHandler) list(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("listing api keys")

	workspace := getWorkspaceFromContext(r.Context())

	keys, err := d.apiRepo.List(ctx, workspace.ID)
	if err != nil {
		logger.Error("could not list keys",
			zap.Error(err))

		return newAPIStatus(http.StatusInternalServerError, "could not list api keys"), StatusFailed
	}

	return listAPIKeysResponse{
		APIStatus: newAPIStatus(http.StatusOK, "api key created"),
		Keys:      keys,
	}, StatusSuccess
}

type revokeAPIKeyRequest struct {
	GenericRequest

	Strategy malak.RevocationType `json:"strategy,omitempty" validate:"required"`
}

func (c *revokeAPIKeyRequest) Validate() error {
	if !c.Strategy.IsValid() {
		return errors.New("please provide a valid revocation strategy")
	}

	return nil
}

// @Description revoke a specific api key
// @Tags developers
// @Accept  json
// @Produce  json
// @Param reference path string required "api key unique reference.. e.g api_key_"
// @Param message body revokeAPIKeyRequest true "api key request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /developers/keys/{reference} [delete]
func (d *apiKeyHandler) revoke(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("revoking api key")

	req := new(revokeAPIKeyRequest)

	workspace := getWorkspaceFromContext(r.Context())

	ref := chi.URLParam(r, "reference")

	if hermes.IsStringEmpty(ref) {
		return newAPIStatus(http.StatusBadRequest, "reference required"), StatusFailed
	}

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	logger = logger.With(zap.String("reference", ref))

	key, err := d.apiRepo.Fetch(ctx, malak.FetchAPIKeyOptions{
		Reference:   malak.Reference(ref),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		var msg = "could not fetch api key"
		var status = http.StatusInternalServerError

		logger.Error("error fetching api key", zap.Error(err))

		if errors.Is(err, malak.ErrAPIKeyNotFound) {
			msg = err.Error()
			status = http.StatusNotFound
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	if key.IsRevoked() {
		return newAPIStatus(http.StatusBadRequest, "api key already revoked"), StatusFailed
	}

	opts := malak.RevokeAPIKeyOptions{
		APIKey:         key,
		RevocationType: req.Strategy,
	}

	if err := d.apiRepo.Revoke(ctx, opts); err != nil {
		logger.Error("could not create api key",
			zap.Error(err))

		var status = http.StatusInternalServerError
		var msg = "could not create api key"

		if errors.Is(err, malak.ErrAPIKeyMaxLimit) {
			status = http.StatusBadRequest
			msg = err.Error()
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "api key revoked"), StatusSuccess
}
