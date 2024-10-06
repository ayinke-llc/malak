package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type Message struct {
	ID       string
	Metadata map[string]string
	Data     []byte
}

type QueueHandler interface {
	io.Closer
	Add(context.Context, string, *Message) error
	Start(context.Context)
}

// ENUM(update_preview)
type QueueEventSubscriptionMessage string

type PreviewUpdateMessage struct {
	UpdateID   uuid.UUID
	ScheduleID uuid.UUID
}

func ToPayload(m any) []byte {
	var b = new(bytes.Buffer)

	_ = json.NewEncoder(b).Encode(m)

	return b.Bytes()
}
