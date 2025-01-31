package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/ayinke-llc/malak"
)

// ENUM(billing_trial_ending,billing_create_customer,
// invite_team_member)
type QueueTopic string

type Message struct {
	ID       string
	Metadata map[string]string
	Data     []byte
}

type QueueHandler interface {
	io.Closer
	Add(context.Context, QueueTopic, any) error
	Start(context.Context)
}

func ToPayload(m any) []byte {
	var b = new(bytes.Buffer)

	_ = json.NewEncoder(b).Encode(m)

	return b.Bytes()
}

type BillingCreateCustomerOptions struct {
	Workspace *malak.Workspace
}
