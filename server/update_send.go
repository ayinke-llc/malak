package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
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
// @Param message body previewUpdateRequest true "request body to create a workspace"
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

	if _, err := u.cache.Exists(ctx, key); err == nil {
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

	// get contact from db
	// if not exists, create one
	schedule := &malak.UpdateSchedule{
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeSchedule),
		SendAt:      time.Now(),
		UpdateType:  malak.UpdateTypePreview,
		ScheduledBy: user.ID,
		Status:      malak.UpdateSendScheduleScheduled,
		UpdateID:    update.ID,
		CreatedAt:   time.Now(),
	}

	opts := &malak.CreateUpdateOptions{
		Reference: func(et malak.EntityType) string {
			return u.referenceGenerator.Generate(et).String()
		},
		Emails:      []malak.Email{req.Email},
		WorkspaceID: workspace.ID,
		Schedule:    schedule,
		Generator:   u.referenceGenerator,
		UserID:      user.ID,
	}

	if err := u.updateRepo.SendUpdate(ctx, opts); err != nil {
		logger.Error("could not create schedule update", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not send update"), StatusFailed
	}

	span.SetAttributes(
		attribute.String("schedule.type", "live"),
		attribute.String("schedule.id", schedule.ID.String()))

	span.AddEvent("update.preview")

	return newAPIStatus(http.StatusOK, "Your preview is now scheduled and will be sent out immediately"),
		StatusSuccess
}

type sendUpdateRequest struct {
	Emails []malak.Email `json:"emails,omitempty"`
	SendAt *int64        `json:"send_at,omitempty"`

	GenericRequest
}

func (s *sendUpdateRequest) Validate() error {
	if len(s.Emails) == 0 {
		return errors.New("please provide atleast one email")
	}

	for _, v := range s.Emails {
		_, err := mail.ParseAddress(v.String())
		if err != nil {
			return err
		}
	}

	if s.SendAt == nil {
		return nil
	}

	scheduledTime := hermes.DeRef(s.SendAt)

	if time.Now().Before(time.Unix(scheduledTime, 0)) {
		return errors.New("you can only schedule to the future not past")
	}

	return nil
}

// @Tags updates
// @Summary Send an update to real users
// @id sendUpdate
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Param message body sendUpdateRequest true "request body to send an update"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference} [post]
func (u *updatesHandler) sendUpdate(
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

	logger.Debug("Sending update")

	req := new(sendUpdateRequest)

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

	var sendAt = time.Now()
	if req.SendAt != nil {
		sendAt = time.Unix(hermes.DeRef(req.SendAt), 0)
	}

	// get contact from db
	// if not exists, create one
	schedule := &malak.UpdateSchedule{
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeSchedule),
		SendAt:      sendAt,
		UpdateType:  malak.UpdateTypeLive,
		ScheduledBy: user.ID,
		Status:      malak.UpdateSendScheduleScheduled,
		UpdateID:    update.ID,
		CreatedAt:   time.Now(),
	}

	opts := &malak.CreateUpdateOptions{
		Reference: func(et malak.EntityType) string {
			return u.referenceGenerator.Generate(et).String()
		},
		Emails:      req.Emails,
		WorkspaceID: workspace.ID,
		Schedule:    schedule,
		Generator:   u.referenceGenerator,
		UserID:      user.ID,
	}

	if err := u.updateRepo.SendUpdate(ctx, opts); err != nil {
		logger.Error("could not create schedule update", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not send update"), StatusFailed
	}

	span.SetAttributes(
		attribute.String("schedule.type", "live"),
		attribute.String("schedule.id", schedule.ID.String()))

	span.AddEvent("update.sending.live")

	return newAPIStatus(http.StatusOK, "Your update is now scheduled and will be sent out"),
		StatusSuccess
}
