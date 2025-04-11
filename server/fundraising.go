package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type fundraisingHandler struct {
	cfg                config.Config
	fundingRepo        malak.FundraisingPipelineRepository
	referenceGenerator malak.ReferenceGeneratorOperation
	contactRepo        malak.ContactRepository
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

// @Description list all fundraising pipelines with pagination and filtering
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param page query int false "Page to query data from. Defaults to 1"
// @Param per_page query int false "Number to items to return. Defaults to 10 items"
// @Param active_only query bool false "Whether to return only active pipelines. Defaults to false"
// @Success 200 {object} fetchPipelinesResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines [get]
func (d *fundraisingHandler) list(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("listing fundraising pipelines")

	workspace := getWorkspaceFromContext(r.Context())
	paginator := malak.PaginatorFromRequest(r)

	activeOnly := false
	if activeOnlyStr := r.URL.Query().Get("active_only"); activeOnlyStr != "" {
		if active, err := strconv.ParseBool(activeOnlyStr); err == nil {
			activeOnly = active
		}
	}

	opts := malak.ListPipelineOptions{
		Paginator:   paginator,
		ActiveOnly:  activeOnly,
		WorkspaceID: workspace.ID,
	}

	pipelines, total, err := d.fundingRepo.List(ctx, opts)
	if err != nil {
		logger.Error("could not list fundraising pipelines", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not list fundraising pipelines"),
			StatusFailed
	}

	return fetchPipelinesResponse{
		APIStatus: newAPIStatus(http.StatusOK, "fetched fundraising pipelines"),
		Pipelines: pipelines,
		Meta: meta{
			Paging: pagingInfo{
				Total:   total,
				PerPage: paginator.PerPage,
				Page:    paginator.Page,
			},
		},
	}, StatusSuccess
}

// @Description Fetch a fundraising board with its columns
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param reference path string true "Pipeline reference"
// @Success 200 {object} fetchBoardResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines/{reference} [get]
func (d *fundraisingHandler) board(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("fetching fundraising board")

	reference := chi.URLParam(r, "reference")
	if hermes.IsStringEmpty(reference) {
		return newAPIStatus(http.StatusBadRequest, "please provide the pipeline reference"), StatusFailed
	}

	logger = logger.With(zap.String("reference", reference))
	workspace := getWorkspaceFromContext(ctx)

	pipeline, err := d.fundingRepo.Get(ctx, malak.FetchPipelineOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrPipelineNotFound) {
			return newAPIStatus(http.StatusNotFound, "fundraising pipeline not found"), StatusFailed
		}

		logger.Error("could not fetch fundraising pipeline", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising pipeline"), StatusFailed
	}

	columns, contacts, positions, err := d.fundingRepo.Board(ctx, pipeline)
	if err != nil {
		logger.Error("could not fetch fundraising board", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising board"), StatusFailed
	}

	return fetchBoardResponse{
		Pipeline:  hermes.DeRef(pipeline),
		Columns:   columns,
		Contacts:  contacts,
		Positions: positions,
		APIStatus: newAPIStatus(http.StatusOK, "fetched fundraising board"),
	}, StatusSuccess
}

// @Description Close a fundraising board permanently
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param reference path string true "Pipeline reference"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines/{reference} [delete]
func (d *fundraisingHandler) closeBoard(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("closing fundraising board")

	reference := chi.URLParam(r, "reference")
	if hermes.IsStringEmpty(reference) {
		return newAPIStatus(http.StatusBadRequest, "please provide the pipeline reference"), StatusFailed
	}

	logger = logger.With(zap.String("reference", reference))

	workspace := getWorkspaceFromContext(ctx)

	pipeline, err := d.fundingRepo.Get(ctx, malak.FetchPipelineOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrPipelineNotFound) {
			return newAPIStatus(http.StatusNotFound, "fundraising pipeline not found"), StatusFailed
		}

		logger.Error("could not fetch fundraising pipeline", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising pipeline"), StatusFailed
	}

	if pipeline.IsClosed {
		return newAPIStatus(http.StatusOK, "fundraising pipeline is already closed"), StatusSuccess
	}

	if err := d.fundingRepo.CloseBoard(ctx, pipeline); err != nil {
		logger.Error("could not close fundraising board", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not close fundraising board"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "fundraising board closed successfully"), StatusSuccess
}

type addContactRequest struct {
	GenericRequest
	ContactReference string `json:"contact_reference" validate:"required"`
	Rating           int    `json:"rating" validate:"required,min=0,max=5"`
	CanLeadRound     bool   `json:"can_lead_round"`
	InitialContact   int64  `json:"initial_contact" validate:"required"`
	CheckSize        int64  `json:"check_size" validate:"required"`
}

func (c *addContactRequest) Validate() error {
	if hermes.IsStringEmpty(c.ContactReference) {
		return errors.New("please provide the contact reference")
	}

	if c.Rating < 0 || c.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}

	if c.CheckSize < (1000 * 100) {
		return errors.New("check size must be at least 1000 USD ($1,000)")
	}

	initialContactDate := time.Unix(c.InitialContact, 0).UTC()
	if initialContactDate.After(time.Now().UTC()) {
		return errors.New("initial contact date cannot be in the future")
	}

	return nil
}

// @Description Add a contact to a fundraising board
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param reference path string true "Pipeline reference"
// @Param message body addContactRequest true "contact request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines/{reference}/contacts [post]
func (d *fundraisingHandler) addContact(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("adding contact to fundraising board")

	reference := chi.URLParam(r, "reference")
	if hermes.IsStringEmpty(reference) {
		return newAPIStatus(http.StatusBadRequest, "please provide the pipeline reference"), StatusFailed
	}

	req := new(addContactRequest)
	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	workspace := getWorkspaceFromContext(ctx)

	logger = logger.With(zap.String("reference", reference),
		zap.String("contact_reference", req.ContactReference))

	pipeline, err := d.fundingRepo.Get(ctx, malak.FetchPipelineOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrPipelineNotFound) {
			return newAPIStatus(http.StatusNotFound, "fundraising pipeline not found"), StatusFailed
		}

		logger.Error("could not fetch fundraising pipeline", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising pipeline"), StatusFailed
	}

	if pipeline.IsClosed {
		return newAPIStatus(http.StatusBadRequest, "this pipeline is closed already"), StatusFailed
	}

	contact, err := d.contactRepo.Get(ctx, malak.FetchContactOptions{
		Reference:   malak.Reference(req.ContactReference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrContactNotFound) {
			return newAPIStatus(http.StatusNotFound, "contact not found"), StatusFailed
		}

		logger.Error("could not fetch contact", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch contact"), StatusFailed
	}

	defaultColumn, err := d.fundingRepo.DefaultColumn(ctx, pipeline)
	if err != nil {
		logger.Error("could not get default column", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not get default column"), StatusFailed
	}

	err = d.fundingRepo.AddContactToBoard(ctx, &malak.AddContactToBoardOptions{
		Column:             &defaultColumn,
		Contact:            contact,
		ReferenceGenerator: d.referenceGenerator,
		Rating:             req.Rating,
		CanLeadRound:       req.CanLeadRound,
		InitialContact:     time.Unix(req.InitialContact, 0).UTC(),
		CheckSize:          req.CheckSize,
	})
	if err != nil {
		logger.Error("could not add contact to board", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not add contact to board"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "contact added to board successfully"), StatusSuccess
}

type updateContactDealRequest struct {
	GenericRequest
	Rating       int   `json:"rating,omitempty" validate:"required,min=0,max=5"`
	CanLeadRound bool  `json:"can_lead_round,omitempty" validate:"required"`
	CheckSize    int64 `json:"check_size,omitempty" validate:"required"`
}

func (c *updateContactDealRequest) Validate() error {
	if c.Rating < 0 || c.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}

	if c.CheckSize < (1000 * 100) {
		return errors.New("check size must be at least 1000 USD ($1,000)")
	}

	return nil
}

// @Description Update deal details for a contact on the fundraising board
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param reference path string true "Pipeline reference"
// @Param contact_id path string true "Contact ID"
// @Param message body updateContactDealRequest true "update deal request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines/{reference}/contacts/{contact_id} [patch]
func (d *fundraisingHandler) updateContactDeal(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("updating contact deal details")

	reference := chi.URLParam(r, "reference")
	if hermes.IsStringEmpty(reference) {
		return newAPIStatus(http.StatusBadRequest, "please provide the pipeline reference"), StatusFailed
	}

	contactID := chi.URLParam(r, "contact_id")
	if hermes.IsStringEmpty(contactID) {
		return newAPIStatus(http.StatusBadRequest, "please provide the contact id"), StatusFailed
	}

	contactUUID, err := uuid.Parse(contactID)
	if err != nil {
		return newAPIStatus(http.StatusBadRequest, "you must provide a valid contact uuid"), StatusFailed
	}

	req := new(updateContactDealRequest)
	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	workspace := getWorkspaceFromContext(ctx)

	logger = logger.With(zap.String("reference", reference),
		zap.String("contact_id", contactID))

	pipeline, err := d.fundingRepo.Get(ctx, malak.FetchPipelineOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrPipelineNotFound) {
			return newAPIStatus(http.StatusNotFound, "fundraising pipeline not found"), StatusFailed
		}

		logger.Error("could not fetch fundraising pipeline contact", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising pipeline"), StatusFailed
	}

	if pipeline.IsClosed {
		return newAPIStatus(http.StatusBadRequest, "this pipeline is closed already"), StatusFailed
	}

	contact, err := d.fundingRepo.GetContact(ctx, pipeline.ID, contactUUID)
	if err != nil {
		if errors.Is(err, malak.ErrContactNotFoundOnBoard) {
			return newAPIStatus(http.StatusNotFound, "this contact is not on this board"), StatusFailed
		}

		logger.Error("could not fetch contact", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "an error occurred while fetching a contact"), StatusFailed
	}

	err = d.fundingRepo.UpdateContactDeal(ctx, pipeline, malak.UpdateContactDealOptions{
		Rating:       int64(req.Rating),
		CanLeadRound: req.CanLeadRound,
		CheckSize:    req.CheckSize,
		ContactID:    contact.ID,
	})
	if err != nil {
		logger.Error("could not update contact deal details", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not update contact deal details"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "contact deal details updated successfully"), StatusSuccess
}

type moveContactAcrossBoardRequest struct {
	GenericRequest
	ColumnID  uuid.UUID `json:"column_id,omitempty" validate:"required"`
	ContactID uuid.UUID `json:"contact_id,omitempty" validate:"required"`
}

func (c *moveContactAcrossBoardRequest) Validate() error {

	if c.ColumnID == uuid.Nil {
		return errors.New("please provide a valid column uuid")
	}

	if c.ContactID == uuid.Nil {
		return errors.New("Please provide a valid contact uuid")
	}

	return nil
}

// @Description move contact across board
// @Tags fundraising
// @Accept  json
// @Produce  json
// @Param reference path string true "Pipeline reference"
// @Param contact_id path string true "Contact ID"
// @Param message body moveContactAcrossBoardRequest true "move cotnact across board"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /pipelines/{reference}/contacts/board [post]
func (d *fundraisingHandler) moveContactAcrossBoard(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("moving contact across board")

	reference := chi.URLParam(r, "reference")
	if hermes.IsStringEmpty(reference) {
		return newAPIStatus(http.StatusBadRequest, "please provide the pipeline reference"), StatusFailed
	}

	req := new(moveContactAcrossBoardRequest)
	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	workspace := getWorkspaceFromContext(ctx)

	logger = logger.With(zap.String("reference", reference),
		zap.String("board_column_id", req.ColumnID.String()),
		zap.String("contact_id", req.ContactID.String()))

	pipeline, err := d.fundingRepo.Get(ctx, malak.FetchPipelineOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: workspace.ID,
	})
	if err != nil {
		if errors.Is(err, malak.ErrPipelineNotFound) {
			return newAPIStatus(http.StatusNotFound, "fundraising pipeline not found"), StatusFailed
		}

		logger.Error("could not fetch fundraising pipeline contact", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch fundraising pipeline"), StatusFailed
	}

	if pipeline.IsClosed {
		return newAPIStatus(http.StatusBadRequest, "this pipeline is closed already"), StatusFailed
	}

	var contact *malak.FundraiseContact
	var column *malak.FundraisingPipelineColumn
	var g errgroup.Group

	g.Go(func() error {

		var err error

		contact, err = d.fundingRepo.GetContact(ctx, pipeline.ID, req.ContactID)
		if err != nil {
			logger.Error("could not fetch contact", zap.Error(err))
			return err
		}

		return nil
	})

	g.Go(func() error {
		var err error

		column, err = d.fundingRepo.GetColumn(ctx, malak.GetBoardOptions{
			PipelineID: pipeline.ID,
			ColumnID:   req.ColumnID,
		})
		if err != nil {
			logger.Error("could not fetch board column", zap.Error(err))
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, malak.ErrPipelineColumnNotFound) || errors.Is(err, malak.ErrContactNotFoundOnBoard) {
			return newAPIStatus(http.StatusNotFound, err.Error()), StatusFailed
		}

		return newAPIStatus(http.StatusInternalServerError, "could not fetch your contact or column"), StatusFailed
	}

	err = d.fundingRepo.UpdateContactDeal(ctx, pipeline, nil)
	if err != nil {
		logger.Error("could not update contact deal details", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not update contact deal details"), StatusFailed
	}

	return newAPIStatus(http.StatusOK, "contact deal details updated successfully"), StatusSuccess
}
