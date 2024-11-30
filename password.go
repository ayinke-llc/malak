package malak

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/ayinke-llc/hermes"
	"golang.org/x/crypto/bcrypt"
)

type Password string

func (p Password) IsZero() bool { return hermes.IsStringEmpty(string(p)) }

func (p Password) String() string { return "****" }

func (p *Password) MarshalJSON() ([]byte, error) {
	if p.IsZero() {
		return json.Marshal("")
	}

	return json.Marshal(p.String())
}

func (p Password) Equals(other Password) bool { return strings.EqualFold(string(p), string(other)) }

func (p Password) Value() (driver.Value, error) {
	return HashPassword(string(p))
}

func HashPassword(p string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(s), err
}

func VerifyPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
