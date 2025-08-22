package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/identity"
	"github.com/barretot/ifkpass/internal/repo"
)

type AuthenticateService struct {
	repo repo.ProfileRepository
	idp  identity.IdentityProviderAdapter
}

func NewAuthenticateService(
	r repo.ProfileRepository,
	idp identity.IdentityProviderAdapter,
) *AuthenticateService {
	return &AuthenticateService{repo: r, idp: idp}
}

func (s *AuthenticateService) Authenticate(ctx context.Context, email, password string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if _, err := s.repo.FindByEmail(ctx, email); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", apperrors.ErrUserNotFound
		}
		return "", fmt.Errorf("repo: find user by email: %w", err)
	}

	token, err := s.idp.SignIn(ctx, email, password)
	if err != nil {
		return "", fmt.Errorf("idp: signin: %w", err)
	}

	if token == nil || *token == "" {
		return "", fmt.Errorf("idp: signin: empty token")
	}

	return *token, nil
}
