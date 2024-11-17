package malak

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// DEPRECATED
func GenerateReference(e EntityType) string {
	return fmt.Sprintf("%s_%s", e.String(), gonanoid.Must())
}

// ENUM(
// workspace,invoice,
// team,invite,contact,
// deck,update,link,room,
// recipient,schedule,list,
// list_email, update_stat,
// recipient_stat)
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
			gonanoid.Must()))
}
