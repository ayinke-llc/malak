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
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type contactHandler struct {
	cfg config.Config
}

type createContactRequest struct {
	GenericRequest

	Email     malak.Email `json:"email,omitempty"`
	FirstName string      `json:"first_name,omitempty"`

	LastName *string `json:"last_name,omitempty"`
}

func (c *createContactRequest) Validate() error {
	if util.IsStringEmpty(c.Email.String()) {
		return errors.New("please provide the email address of the contact")
	}

	c.FirstName = strings.TrimSpace(c.FirstName)

	if util.IsStringEmpty(c.FirstName) {
		return errors.New("please provide the first name of the contact")
	}

	if len(c.FirstName) > 100 {
		return errors.New("contact's first name must be less than 100")
	}

	lastName := strings.TrimSpace(util.DeRef(c.LastName))

	if !util.IsStringEmpty(lastName) {
		if len(lastName) > 100 {
			return errors.New("contact's last name must be less than 100")
		}
	}

	c.LastName = util.Ref(lastName)

	return nil
}

// @Summary Create a new contact
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
	logger *logrus.Entry,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	return newAPIStatus(http.StatusInternalServerError, "o"), StatusFailed
}
