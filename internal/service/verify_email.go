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
)

type VerifyEmailService struct {
	repo repo.ProfileRepository
	idp  identity.IdentityProviderAdapter
}

func NewVerifyEmailService(
	r repo.ProfileRepository,
	idp identity.IdentityProviderAdapter,
) *VerifyEmailService {
	return &VerifyEmailService{repo: r, idp: idp}
}

func (s *VerifyEmailService) VerifyEmail(ctx context.Context, cfg config.AppConfig, email, password, code string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if _, err := s.repo.FindByEmail(ctx, email); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", apperrors.ErrUserNotFound
		}
		return "", fmt.Errorf("repo: find user by email: %w", err)
	}

	isVerified, err := s.idp.IsEmailVerified(ctx, cfg, email)
	if err != nil {
		return "", fmt.Errorf("idp: check email verified: %w", err)
	}
	if isVerified {
		return "", apperrors.ErrEmailAlreadyVerified
	}

	if err := s.idp.ConfirmEmail(ctx, cfg, email, code); err != nil {
		return "", fmt.Errorf("idp: confirm email: %w", err)
	}

	token, err := s.idp.SignIn(ctx, cfg, email, password)
	if err != nil {
		return "", fmt.Errorf("idp: signin: %w", err)
	}

	if token == nil || *token == "" {
		return "", fmt.Errorf("idp: signin: empty token")
	}

	return *token, nil
}
