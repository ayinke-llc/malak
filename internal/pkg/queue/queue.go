package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

func ToPayload(m any) []byte {
	var b = new(bytes.Buffer)

	_ = json.NewEncoder(b).Encode(m)

	return b.Bytes()
}
