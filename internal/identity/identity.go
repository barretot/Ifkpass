package identity

import (
	"context"

	"github.com/barretot/ifkpass/internal/config"
)

type IdentityProviderAdapter interface {
	SignUp(ctx context.Context, cfg config.AppConfig, email, password string) (string, error)
	SignIn(ctx context.Context, cfg config.AppConfig, email, password string) (*string, error)
	GetUserId(ctx context.Context, cfg config.AppConfig, email string) (string, error)
	IsEmailVerified(ctx context.Context, cfg config.AppConfig, email string) (bool, error)
	ConfirmEmail(ctx context.Context, cfg config.AppConfig, email, code string) error
}
