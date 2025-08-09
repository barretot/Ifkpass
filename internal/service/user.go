package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/identity"
	"github.com/barretot/ifkpass/internal/repo"
	"github.com/barretot/ifkpass/internal/store/dynamostore/models"
)

type UserService struct {
	repo             repo.ProfileRepository
	identityprovider identity.IdentityProviderAdapter
}

func NewUserService(
	r repo.ProfileRepository,
	identityprovider identity.IdentityProviderAdapter,
) *UserService {
	return &UserService{repo: r, identityprovider: identityprovider}
}

func (s *UserService) CreateUser(ctx context.Context, cfg config.AppConfig, name, lastname, email, password string) error {
	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		if !errors.Is(err, apperrors.ErrUserNotFound) {
			return fmt.Errorf("find user by email: %w", err)
		}
	} else if user != nil {
		return apperrors.ErrUserAlreadyExists
	}

	if err := s.identityprovider.SignUp(ctx, cfg, email, password); err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return apperrors.ErrUserAlreadyExists
		}
		return fmt.Errorf("cognito signup: %w", err)
	}

	userID, err := s.identityprovider.GetUserId(ctx, cfg, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return fmt.Errorf("get user id from identity provider: %w", err)
	}

	if err := s.repo.Save(ctx, models.User{
		UserId:   userID,
		Name:     name,
		LastName: lastname,
		Email:    email,
	}); err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	return nil
}
