package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/config"
	"github.com/barretot/ifkpass/internal/identity"
	"github.com/barretot/ifkpass/internal/repo"
	"github.com/barretot/ifkpass/internal/store/dynamostore/models"
)

type UserService struct {
	repo repo.ProfileRepository
	idp  identity.IdentityProviderAdapter
}

func NewUserService(
	r repo.ProfileRepository,
	idp identity.IdentityProviderAdapter,
) *UserService {
	return &UserService{repo: r, idp: idp}
}

func (s *UserService) CreateUser(ctx context.Context, cfg config.AppConfig, name, lastname, email, password string) error {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if _, err := s.repo.FindByEmail(ctx, email); err != nil {
		if !errors.Is(err, apperrors.ErrUserNotFound) {
			return fmt.Errorf("repo: find user by email: %w", err)
		}
	} else {
		return apperrors.ErrUserAlreadyExists
	}

	userID, err := s.idp.SignUp(ctx, cfg, email, password)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return apperrors.ErrUserAlreadyExists
		}
		return fmt.Errorf("idp: signup: %w", err)
	}

	if err := s.repo.Save(ctx, models.User{
		UserId:   userID,
		Name:     name,
		LastName: lastname,
		Email:    email,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return apperrors.ErrUserAlreadyExists
		}
		return fmt.Errorf("repo: save user: %w", err)
	}

	return nil
}
