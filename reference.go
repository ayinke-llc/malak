package malak

import (
	"crypto/rand"
	"encoding/base64"
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
// update,link,room,
// recipient,schedule,list,
// list_email, update_stat,
// recipient_stat,recipient_log,
// deck,deck_preference)
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

func ShortLink() string {
	b := make([]byte, 8) // 6 bytes = 8 characters in base64
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
