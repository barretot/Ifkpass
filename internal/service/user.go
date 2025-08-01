package service

import (
	"context"
	"errors"

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

func (s *UserService) CreateUser(ctx context.Context, name, email string) error {
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return errors.New("user already exists")
	}

	user := models.User{
		ID:    util.GenerateUUID(),
		Name:  name,
		Email: email,
	}

	return s.repo.Save(ctx, user)
}
