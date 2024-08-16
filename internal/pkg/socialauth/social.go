package socialauth

import (
	"context"

	"golang.org/x/oauth2"
)

type User struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type ValidateOptions struct {
	Code string
}

type SocialAuthProvider interface {
	User(context.Context, string) (User, error)
	Validate(context.Context, ValidateOptions) (*oauth2.Token, error)
}
