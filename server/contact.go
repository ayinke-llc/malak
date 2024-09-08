package server

import (
	"errors"
	"strings"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/util"
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
