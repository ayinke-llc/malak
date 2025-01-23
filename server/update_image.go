package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adelowo/gulter"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// @Summary Upload an image
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

	uploadedURL := fmt.Sprintf("%s/%s/%s",
		u.cfg.Uploader.S3.Endpoint,
		file.FolderDestination,
		file.UploadedFileName)

	return uploadImageResponse{
		URL:       uploadedURL,
		APIStatus: newAPIStatus(http.StatusOK, "image was uploaded"),
	}, StatusSuccess
}
