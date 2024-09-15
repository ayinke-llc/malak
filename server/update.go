package server

import (
	"context"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type updatesHandler struct {
	referenceGenerator malak.ReferenceGeneratorOperation
	updateRepo         malak.UpdateRepository
}

// @Summary Create a new update
// @Tags updates
// @Accept  json
// @Produce  json
// @Success 200 {object} createdUpdateResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates [post]
func (u *updatesHandler) create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating a new update")

	user := getUserFromContext(r.Context())
	workspace := getWorkspaceFromContext(r.Context())

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeUpdate),
		Status:      malak.UpdateStatusDraft,
		Content:     malak.UpdateContent(""),
		Metadata:    malak.UpdateMetadata{},
	}

	if err := u.updateRepo.Create(ctx, update); err != nil {

		logger.Error("could not create update",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not create a new update"), StatusFailed
	}

	return createdUpdateResponse{
		Update:    util.DeRef(update),
		APIStatus: newAPIStatus(http.StatusCreated, "update successfully created"),
	}, StatusSuccess
}
