package malak

import (
	"context"
	"time"

	"github.com/ayinke-llc/hermes"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ErrEmailVerificationNotFound = MalakError("email verification token not found")
)

type EmailVerification struct {
	Token     string    `json:"token"`
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	bun.BaseModel `bun:"table:email_verifications" json:"-"`
}

func NewEmailVerification(u *User) (*EmailVerification, error) {
	val, err := hermes.Random(40)
	if err != nil {
		return nil, err
	}

	return &EmailVerification{
		Token:     val,
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}, nil
}

type EmailVerificationRepository interface {
	Create(context.Context, *EmailVerification) error
	Get(context.Context, string) (*EmailVerification, error)
}
