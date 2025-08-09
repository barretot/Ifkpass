package identity

import (
	"context"

	"github.com/barretot/ifkpass/internal/config"
)

type IdentityProviderAdapter interface {
	SignUp(ctx context.Context, cfg config.AppConfig, email, password string) error
	GetUserId(ctx context.Context, cfg config.AppConfig, email string) (string, error)
}
