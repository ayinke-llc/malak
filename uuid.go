package malak

import "github.com/google/uuid"

type UuidGenerator interface {
	Create() uuid.UUID
}

type googleuuid struct{}

func (g *googleuuid) Create() uuid.UUID { return uuid.New() }

func NewGoogleUUID() UuidGenerator {
	return &googleuuid{}
}
