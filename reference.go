package malak

import (
	"fmt"

	"github.com/teris-io/shortid"
)

// DEPRECATED
func GenerateReference(e EntityType) string {
	return fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate())
}

// ENUM(workspace,invoice,team,invite,contact,deck,update)
type EntityType string

type Reference string

type ReferenceGeneratorOperation interface {
	Generate(EntityType) Reference
}

type ReferenceGenerator struct{}

func NewReferenceGenerator() *ReferenceGenerator {
	return &ReferenceGenerator{}
}

func (r *ReferenceGenerator) Generate(e EntityType) Reference {
	return Reference(fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate()))
}
