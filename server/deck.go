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
	// GenericRequest

}
