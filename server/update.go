package server

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/cache"
	"github.com/ayinke-llc/malak/internal/pkg/queue"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type updatesHandler struct {
	referenceGenerator malak.ReferenceGeneratorOperation
	updateRepo         malak.UpdateRepository
	cfg                config.Config
	cache              cache.Cache
	queueHandler       queue.QueueHandler
}

type createUpdateContent struct {
	Title string `json:"title,omitempty" validate:"required"`
	GenericRequest
}

func (c *createUpdateContent) Validate() error {

	p := bluemonday.StrictPolicy()

	c.Title = p.Sanitize(c.Title)

	if util.IsStringEmpty(c.Title) {
		return errors.New("please provide update title")
	}

	if len(c.Title) < 5 {
		return errors.New("title must be more than 5 characters")
	}

	return nil
}

// @Summary Create a new update
// @Tags updates
// @Accept  json
// @Produce  json
// @Param message body createUpdateContent true "update content body"
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

	req := new(createUpdateContent)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	update := &malak.Update{
		WorkspaceID: workspace.ID,
		CreatedBy:   user.ID,
		Content:     malak.BlockContents{},
		Reference:   u.referenceGenerator.Generate(malak.EntityTypeUpdate),
		Status:      malak.UpdateStatusDraft,
		Metadata:    malak.UpdateMetadata{},
		Title:       req.Title,
	}

	if err := u.updateRepo.Create(ctx, update); err != nil {

		logger.Error("could not create update",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not create a new update"), StatusFailed
	}

	span.AddEvent("workspace.new", trace.WithAttributes(
		attribute.String("id", update.Reference.String())))

	return createdUpdateResponse{
		Update:    util.DeRef(update),
		APIStatus: newAPIStatus(http.StatusCreated, "update successfully created"),
	}, StatusSuccess
}

// @Summary List updates
// @Tags updates
// @Accept  json
// @Produce  json
// @Param page query int false "Page to query data from. Defaults to 1"
// @Param per_page query int false "Number to items to return. Defaults to 10 items"
// @Param status query string false "filter results by the status of the update."
// @Success 200 {object} listUpdateResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates [get]
func (u *updatesHandler) list(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Listing updates")

	workspace := getWorkspaceFromContext(r.Context())

	filterStatus := malak.ListUpdateFilterStatus(r.URL.Query().Get("view"))

	if !filterStatus.IsValid() {
		filterStatus = malak.ListUpdateFilterStatusAll
	}

	opts := malak.ListUpdateOptions{
		Status:      filterStatus,
		Paginator:   malak.PaginatorFromRequest(r),
		WorkspaceID: workspace.ID,
	}

	span.SetAttributes(
		append(opts.Paginator.OTELAttributes(),
			attribute.String("view",
				filterStatus.String()))...)

	updates, err := u.updateRepo.List(ctx, opts)
	if err != nil {

		logger.Error("could not list updates",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not list updates"), StatusFailed
	}

	return listUpdateResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "updates fetched"),
		Updates:   updates,
		Meta: meta{
			Paging: pagingInfo{
				PerPage: opts.Paginator.PerPage,
				Page:    opts.Paginator.Page,
			},
		},
	}, StatusSuccess
}

type contentUpdateRequest struct {
	Update malak.BlockContents `json:"update,omitempty" validate:"required"`
	Title  string              `json:"title,omitempty" validate:"required"`
	GenericRequest
}

var compiledAllowRegexp = regexp.MustCompile(`[a-z; -]*`)

func (c *contentUpdateRequest) Validate() error {

	p := bluemonday.UGCPolicy()

	// p.AllowDataAttributes()

	// Youtube iframe check
	// p.AllowElements("iframe")
	// p.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	// p.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	// p.AllowAttrs("src").OnElements("iframe")
	// p.AllowAttrs("frameborder").Matching(bluemonday.Number).OnElements("iframe")
	// p.AllowAttrs("allow").Matching(compiledAllowRegexp).OnElements("iframe")
	// p.AllowAttrs("allowfullscreen").OnElements("iframe")
	//
	// TWITTER embed
	// p.AllowAttrs("src").OnElements("div")
	// p.AllowStyles("color").OnElements("span")

	sanitized, err := malak.SanitizeBlocks(c.Update)
	if err != nil {
		return err
	}

	if len(sanitized) == 0 {
		return errors.New("please provide the content")
	}

	c.Title = p.Sanitize(c.Title)
	c.Update = sanitized

	if util.IsStringEmpty(c.Title) {
		return errors.New("please provide update title")
	}

	if len(c.Title) < 5 {
		return errors.New("title must be more than 5 characters")
	}

	return nil
}

type BlockContent struct {
	Type    string                 `json:"type"`
	Content []BlockContentItem     `json:"content"`
	Props   map[string]interface{} `json:"props,omitempty"`
}

type BlockContentItem struct {
	Type  string                 `json:"type"`
	Text  string                 `json:"text,omitempty"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// @Summary Update a specific update
// @Tags updates
// @id updateContent
// @Accept  json
// @Produce  json
// @Param reference path string required "update unique reference.. e.g update_"
// @Param message body contentUpdateRequest true "update content body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /workspaces/updates/{reference} [put]
func (u *updatesHandler) update(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	ref := chi.URLParam(r, "reference")

	span.SetAttributes(attribute.String("reference", ref))

	logger = logger.With(zap.String("reference", ref))

	logger.Debug("Updating specific update")

	req := new(contentUpdateRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
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

	update.Content = req.Update
	update.Title = req.Title

	if err := u.updateRepo.Update(ctx, update); err != nil {
		logger.Error("could not update content", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError,
			"could not update content"), StatusFailed
	}

	return newAPIStatus(http.StatusOK,
		"updates stored"), StatusSuccess
}
