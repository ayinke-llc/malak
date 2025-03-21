package server

import (
	"context"
	"net/http"

	"github.com/adelowo/gulter"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Description Upload an image
// @Tags images
// @id uploadImage
// @Accept  json
// @Produce  json
// @Param image_body formData file true "image body"
// @Success 200 {object} uploadImageResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /uploads/images [post]
func (u *updatesHandler) uploadImage(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("file uploaded using Gulter")

	files, err := gulter.FilesFromContextWithKey(r, "image_body")
	if err != nil {
		logger.Error("could not fetch gulter uploaded files", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal failure while fetching file from storage"), StatusFailed
	}

	// only one file we are expecting at a time
	file := files[0]

	uploadedURL, err := u.gulter.Path(ctx, gulter.PathOptions{
		Key: file.StorageKey,
	})

	if err != nil {
		logger.Error("could not fetch gulter uploaded file path", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"internal failure while fetching path from storage"), StatusFailed
	}

	return uploadImageResponse{
		URL:       uploadedURL,
		APIStatus: newAPIStatus(http.StatusOK, "image was uploaded"),
	}, StatusSuccess
}
