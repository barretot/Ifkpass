package repo

import (
	"context"

	"github.com/barretot/ifkpass/internal/store/dynamostore/models"
)

type ProfileRepository interface {
	Save(ctx context.Context, user models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
