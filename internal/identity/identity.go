package identity

import (
	"context"
)

type IdentityProviderAdapter interface {
	SignUp(ctx context.Context, email, password string) (string, error)
	SignIn(ctx context.Context, email, password string) (*string, error)
	GetUserId(ctx context.Context, email string) (string, error)
	IsEmailVerified(ctx context.Context, email string) (bool, error)
	ConfirmEmail(ctx context.Context, email, code string) error
}
