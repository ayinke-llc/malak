package malak

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/ayinke-llc/hermes"
	"github.com/theopenlane/utils/passwd"
)

// deprecated
type Password string

func (p Password) IsZero() bool { return hermes.IsStringEmpty(string(p)) }

func (p Password) String() string { return "****" }

func (p *Password) MarshalJSON() ([]byte, error) {
	if p.IsZero() {
		return json.Marshal("")
	}

	return json.Marshal(p.String())
}

func (p Password) Value() (driver.Value, error) {
	return HashPassword(string(p))
}

func HashPassword(p string) (string, error) { return passwd.CreateDerivedKey(p) }

func VerifyPassword(hashed, plain string) bool {
	ok, _ := passwd.VerifyDerivedKey(hashed, plain)
	return ok
}
