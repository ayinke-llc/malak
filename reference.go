package malak

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/ayinke-llc/hermes"
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
// deck,deck_preference, contact_share,dashboard,
// plan,price,integration,workspace_integration, integration_datapoint,
// integration_chart, integration_sync_checkpoint,dashboard_chart,system_template,
// deck_daily_engagement, deck_analytic, deck_viewer_session,
// deck_geographic_stat, session,dashboard_link,dashboard_link_access_log)
type EntityType string

type Reference string

func (r Reference) String() string { return string(r) }

type ReferenceGeneratorOperation interface {
	Generate(EntityType) Reference
	ShortLink() string
	Token() string
}

type ReferenceGenerator struct{}

func NewReferenceGenerator() *ReferenceGenerator {
	return &ReferenceGenerator{}
}

func (r *ReferenceGenerator) Generate(e EntityType) Reference {

	nanoID := gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._", 10)
	return Reference(
		fmt.Sprintf(
			"%s_%s",
			e.String(), nanoID))
}

func (r *ReferenceGenerator) Token() string {
	s, err := hermes.Random(20)
	if err != nil {
		panic(err)
	}

	return s
}

func (r *ReferenceGenerator) ShortLink() string {
	b := make([]byte, 8) // 6 bytes = 8 characters in base64
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
