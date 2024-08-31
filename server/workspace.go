package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type workspaceHandler struct {
	logger        *logrus.Entry
	cfg           config.Config
	userRepo      malak.UserRepository
	workspaceRepo malak.WorkspaceRepository
}

type createWorkspaceRequest struct {
}

// @Summary Create a new workspace
// @Tags workspace
// @Accept  json
// @Produce  json
// @Param message body createWorkspaceRequest true "request body to create a workspace"
// @Success 200 {object} createdUserResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces [post]
func (wo *workspaceHandler) createWorkspace(
	ctx context.Context,
	span trace.Span,
	logger *logrus.Entry,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	fmt.Println(getUserFromContext(r.Context()))
	return nil, StatusFailed
}
