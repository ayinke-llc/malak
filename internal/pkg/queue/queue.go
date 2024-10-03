package queue

import "io"

type QueueHandler interface {
	io.Closer
}

// ENUM(update_preview)
type QueueEventSubscriptionMessage string
