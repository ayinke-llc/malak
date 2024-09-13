package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/util"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contactHandler struct {
	cfg                config.Config
	contactRepo        malak.ContactRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createContactRequest struct {
	GenericRequest

	Email     malak.Email `json:"email,omitempty" validate:"'required'"`
	FirstName *string     `json:"first_name,omitempty" validate:"'required'"`

	LastName *string `json:"last_name,omitempty" validate:"'required'"`
}

func (c *createContactRequest) Validate() error {
	if util.IsStringEmpty(c.Email.String()) {
		return errors.New("please provide the email address of the contact")
	}

	firstName := strings.TrimSpace(util.DeRef(c.FirstName))

	if !util.IsStringEmpty(firstName) {
		if len(firstName) > 100 {
			return errors.New("contact's last name must be less than 100")
		}
	}

	c.FirstName = util.Ref(firstName)

	lastName := strings.TrimSpace(util.DeRef(c.LastName))

	if !util.IsStringEmpty(lastName) {
		if len(lastName) > 100 {
			return errors.New("contact's last name must be less than 100")
		}
	}

	c.LastName = util.Ref(lastName)

	return nil
}

// @Summary Creates a new contact
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param message body createContactRequest true "contact request body"
// @Success 200 {object} fetchContactResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts [post]
func (c *contactHandler) Create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	logger.Debug("creating workspace")

	req := new(createContactRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	contact := &malak.Contact{
		Email:       req.Email,
		FirstName:   util.DeRef(req.FirstName),
		LastName:    util.DeRef(req.LastName),
		Metadata:    make(malak.CustomContactMetadata),
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
		Reference:   c.referenceGenerator.Generate(malak.EntityTypeContact),
		// Default to the user who created it
		OwnerID:   user.ID,
		CreatedBy: user.ID,
	}

	err := c.contactRepo.Create(ctx, contact)
	if errors.Is(err, malak.ErrContactExists) {
		return newAPIStatus(http.StatusConflict, err.Error()), StatusFailed
	}

	if err != nil {
		logger.
			Error("an error occurred while storing contact to the database", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while creating contact"), StatusFailed
	}

	return fetchContactResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "contact was successfully created"),
		Contact:   util.DeRef(contact),
	}, StatusSuccess
}
