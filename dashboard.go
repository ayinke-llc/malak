package malak

import "github.com/google/uuid"

type Dashboard struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Reference Reference `json:"reference,omitempty"`
}
