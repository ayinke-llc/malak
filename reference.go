package malak

import (
	"fmt"

	"github.com/teris-io/shortid"
)

// ENUM(workspace,invoice,team,invite)
type EntityType string

func GenerateReference(e EntityType) string {
	return fmt.Sprintf("%s_%s", e.String(), shortid.MustGenerate())
}
