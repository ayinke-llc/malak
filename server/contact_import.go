package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type createContactRequestBatch struct {
	GenericRequest

	Contacts []struct {
		Email     malak.Email `json:"email,omitempty" validate:"'required'"`
		FirstName *string     `json:"first_name,omitempty" validate:"'required'"`

		LastName *string `json:"last_name,omitempty" validate:"'required'"`
	} `json:"contacts,omitempty"`
}

func (v *createContactRequestBatch) Validate() error {
	if len(v.Contacts) == 0 {
		return errors.New("please provide at least one contact")
	}

	for _, c := range v.Contacts {

		if hermes.IsStringEmpty(c.Email.String()) {
			return errors.New("please provide the email address of the contact")
		}

		firstName := strings.TrimSpace(hermes.DeRef(c.FirstName))

		if !hermes.IsStringEmpty(firstName) {
			if len(firstName) > 100 {
				return errors.New("contact's last name must be less than 100")
			}
		}

		if hermes.IsStringEmpty(firstName) {
			firstName = c.Email.String()
		}

		c.FirstName = hermes.Ref(firstName)

		lastName := strings.TrimSpace(hermes.DeRef(c.LastName))

		if !hermes.IsStringEmpty(lastName) {
			if len(lastName) > 100 {
				return errors.New("contact's last name must be less than 100")
			}
		}

		c.LastName = hermes.Ref(lastName)
	}

	return nil
}

// @Description batch create a new contact
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param message body createContactRequestBatch true "contact request body"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/batch [post]
func (c *contactHandler) batchCreate(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	logger.Debug("creating contact batch")

	req := new(createContactRequestBatch)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	var contacts = make([]*malak.Contact, 0, len(req.Contacts))

	for _, v := range req.Contacts {
		contact := &malak.Contact{
			Email:       v.Email,
			FirstName:   hermes.DeRef(v.FirstName),
			LastName:    hermes.DeRef(v.LastName),
			Metadata:    make(malak.CustomContactMetadata),
			WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
			Reference:   c.referenceGenerator.Generate(malak.EntityTypeContact),
			OwnerID:     user.ID,
			CreatedBy:   user.ID,
		}

		contacts = append(contacts, contact)
	}

	err := c.contactRepo.Create(ctx, contacts...)
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

	return newAPIStatus(http.StatusCreated, "your contacts were uploaded successfully"), StatusSuccess
}
