package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type fundraisingHandler struct {
	cfg                config.Config
	fundingRepo        malak.FundraisingPipelineRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createNewPipelineRequest struct {
	GenericRequest

	Title             string                       `json:"title,omitempty" validate:"required"`
	Stage             malak.FundraisePipelineStage `json:"stage,omitempty" validate:"required"`
	Amount            int64                        `json:"amount,omitempty" validate:"required"`
	Description       string                       `json:"description,omitempty" validate:"required"`
	ExpectedCloseDate int64                        `json:"expected_close_date,omitempty" validate:"required"`
	StartDate         int64                        `json:"start_date,omitempty" validate:"required"`
}

func (c *createNewPipelineRequest) Validate() error {
	if hermes.IsStringEmpty(c.Title) {
		return errors.New("please provide the title of the deck")
	}

	if len(c.Title) < 5 {
		return errors.New("title must be at least 5 characters long")
	}

	if len(c.Description) > 200 {
		return errors.New("description must not exceed 200 characters")
	}

	p := bluemonday.StrictPolicy()

	c.Title = p.Sanitize(c.Title)
	c.Description = p.Sanitize(c.Description)

	if !c.Stage.IsValid() {
		return errors.New("fundraising stage is unsupported at the moment")
	}

	currentTime := time.Now().UTC()
	expectedCloseDate := time.Unix(c.ExpectedCloseDate, 0).UTC()
	startDate := time.Unix(c.StartDate, 0).UTC()

	if startDate.Before(currentTime) && !(startDate.Year() == currentTime.Year() && startDate.YearDay() == currentTime.YearDay()) {
		return errors.New("start date must be today or in the future")
	}

	if !expectedCloseDate.After(currentTime) {
		return errors.New("expected close date must be in the future")
	}

	if expectedCloseDate.Year() == currentTime.Year() && expectedCloseDate.YearDay() == currentTime.YearDay() {
		return errors.New("expected close date cannot be today")
	}

	return nil
}

// @Description Creates a new fundraising pipeline entry
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param message body createNewPipelineRequest true "pipeline request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines [post]
func (d *fundraisingHandler) newPipeline(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("creating fundraising pipeline")

	req := new(createNewPipelineRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	pipeline := &malak.FundraisingPipeline{
		ID:                uuid.New(),
		Reference:         d.referenceGenerator.Generate(malak.EntityTypeFundraisingPipeline),
		WorkspaceID:       getWorkspaceFromContext(ctx).ID,
		Title:             req.Title,
		Stage:             req.Stage,
		TargetAmount:      req.Amount,
		Description:       req.Description,
		StartDate:         time.Unix(req.StartDate, 0).UTC(),
		ExpectedCloseDate: time.Unix(req.ExpectedCloseDate, 0).UTC(),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		IsClosed:          false,
		ClosedAmount:      0,
	}

	defaultColumns := make([]malak.FundraisingPipelineColumn, len(malak.DefaultFundraisingColumns))
	for i, col := range malak.DefaultFundraisingColumns {
		defaultColumns[i] = malak.FundraisingPipelineColumn{
			Reference:      d.referenceGenerator.Generate(malak.EntityTypeFundraisingPipelineColumn),
			Title:          col.Title,
			ColumnType:     col.ColumnType,
			Description:    col.Description,
			InvestorsCount: 0,
		}
	}

	if err := d.fundingRepo.Create(ctx, pipeline, defaultColumns...); err != nil {
		logger.Error("could not create fundraising pipeline",
			zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not create fundraising pipeline"),
			StatusFailed
	}

	return newAPIStatus(http.StatusOK, "pipeline created"), StatusSuccess
}
