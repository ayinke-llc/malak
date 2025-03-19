package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/secret"
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

	Title string `json:"title,omitempty"`
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
		return newAPIStatus(http.StatusInternalServerError, "could not create api key"),
			StatusFailed
	}

	return createdAPIKeyResponse{
		APIStatus: newAPIStatus(http.StatusOK, "api key created"),
		Value:     value,
	}, StatusSuccess
}
