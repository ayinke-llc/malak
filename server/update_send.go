package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"sync"
	"time"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type previewUpdateRequest struct {
	Email malak.Email `json:"email,omitempty" validate:"required"`

	GenericRequest
}

func (p *previewUpdateRequest) Validate() error {
	if util.IsStringEmpty(p.Email.String()) {
		return errors.New("please provide the email to send the preview to")
	}

	_, err := mail.ParseAddress(p.Email.String())
	if err != nil {
		return errors.New("email is invalid")
	}

	return nil
}

// @Tags updates
// @Summary Send preview of an update
// @id previewUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference}/preview [post]
func (u *updatesHandler) previewUpdate(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	workspace := getWorkspaceFromContext(ctx)

	user := getUserFromContext(ctx)

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Sending preview of update")

	// workspaceID -> update_ref
	// This makes sure we can throttle the rate at which
	// preview emails are sent because they are not charged and
	// can thus be abused
	//
	// This blockage is quite simplistic as it does not account for
	// silmutaneous requests. It is as simple as can be for now.
	key := fmt.Sprintf("%s-%s", workspace.ID, ref)

	if _, err := u.cache.Exists(ctx, key); err != nil {
		return newAPIStatus(http.StatusTooManyRequests,
				"please wait a few more minutes before sending another preview of this email"),
			StatusFailed
	}

	req := new(previewUpdateRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"),
			StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	update, err := u.updateRepo.Get(ctx, malak.FetchUpdateOptions{
		Reference: malak.Reference(ref),
	})
	if errors.Is(err, malak.ErrUpdateNotFound) {
		return newAPIStatus(http.StatusNotFound,
			"update does not exists"), StatusFailed
	}

	if err != nil {
		logger.Error("could not fetch update", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"an error occurred while fetching update"), StatusFailed
	}

	schedule := &malak.UpdateSchedule{
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeSchedule),
		SendAt:      time.Now(),
		UpdateType:  malak.RecipientTypePreview,
		ScheduledBy: user.ID,
		Status:      malak.UpdateSendScheduleScheduled,
		UpdateID:    update.ID,
		CreatedAt:   time.Now(),
	}

	if err := u.updateRepo.CreateSchedule(ctx, schedule); err != nil {
		logger.Error("could not create schedule update", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not send preview update"), StatusFailed
	}

	span.SetAttributes(
		attribute.String("schedule.type", "preview"),
		attribute.String("schedule.id", schedule.ID.String()))

	span.AddEvent("update.preview")

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()
		err := u.cache.Add(ctx, key, []byte("ok"), time.Hour)
		if err != nil {
			logger.Error("could not add user throttling to cache",
				zap.Error(err))
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		err := u.queueHandler.Add(ctx, queue.QueueEventSubscriptionMessageUpdatePreview, nil)
		if err != nil {
			logger.Error("could not add schedule to queue to be processed",
				zap.Error(err))
		}
	}()

	wg.Wait()

	return newAPIStatus(http.StatusOK, "Preview email sent"),
		StatusSuccess
}
