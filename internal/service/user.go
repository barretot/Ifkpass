package service

import (
	"context"

	"github.com/barretot/ifkpass/internal/apperrors"
	"github.com/barretot/ifkpass/internal/repo"
	"github.com/barretot/ifkpass/internal/store/dynamostore/models"
	"github.com/barretot/ifkpass/internal/util"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, name, lastname, email string) error {
	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return err
	}

	if user != nil {
		return apperrors.ErrorUserAlreadyExists
	}

	return s.repo.Save(ctx, models.User{
		UserId:   util.GenerateUUID(),
		Name:     name,
		LastName: lastname,
		Email:    email,
	})
}
