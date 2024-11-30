package server

import (
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
)

type deckHandler struct {
	cfg                config.Config
	contactRepo        malak.ContactRepository
	contactListRepo    malak.ContactListRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createDeckRequest struct {
	GenericRequest

	Email     malak.Email `json:"email,omitempty" validate:"'required'"`
	FirstName *string     `json:"first_name,omitempty" validate:"'required'"`

	LastName *string `json:"last_name,omitempty" validate:"'required'"`
}
