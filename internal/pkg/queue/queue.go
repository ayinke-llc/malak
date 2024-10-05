package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/ayinke-llc/malak"
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

type GenericMessage struct{}

func (g GenericMessage) Payload() ([]byte, error) {
	var b = new(bytes.Buffer)
	return b.Bytes(), json.NewEncoder(b).Encode(g)
}

type PreviewUpdateMessage struct {
	Update   *malak.Update
	Schedule *malak.UpdateSchedule
	GenericMessage
}
