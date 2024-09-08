package malak

import (
	"fmt"

	"github.com/teris-io/shortid"
)

// ENUM(workspace,invoice,team,invite,contact,deck,update)
type EntityType string

type Reference string

func NewReference(e EntityType) Reference {
	return Reference(fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate()))
}

func GenerateReference(e EntityType) string {
	return fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate())
}
