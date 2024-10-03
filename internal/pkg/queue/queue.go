package queue

import (
	"context"
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

// ENUM(update_preview)
type QueueEventSubscriptionMessage string
