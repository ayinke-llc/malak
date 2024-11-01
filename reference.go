package malak

import (
	"fmt"

	"github.com/teris-io/shortid"
)

// DEPRECATED
func GenerateReference(e EntityType) string {
	return fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate())
}

// ENUM(
// workspace,invoice,
// team,invite,contact,
// deck,update,link,room,
// recipient,schedule,list,list_email)
type EntityType string

type Reference string

func (r Reference) String() string { return string(r) }

type ReferenceGeneratorOperation interface {
	Generate(EntityType) Reference
}

type ReferenceGenerator struct{}

func NewReferenceGenerator() *ReferenceGenerator {
	return &ReferenceGenerator{}
}

func (r *ReferenceGenerator) Generate(e EntityType) Reference {
	return Reference(
		fmt.Sprintf(
			"%s_%s",
			e.String(),
			shortid.MustGenerate()))
}
