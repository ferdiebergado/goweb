//go:generate mockgen -destination=mock/user_service_mock.go -package=mock . UserService
package service

import (
	"context"
	"fmt"

	"github.com/ferdiebergado/goweb/internal/model"
	"github.com/ferdiebergado/goweb/internal/pkg/security"
	"github.com/ferdiebergado/goweb/internal/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, params RegisterUserParams) (*model.User, error)
}

type userService struct {
	repo   repository.UserRepo
	hasher security.Hasher
}

var _ UserService = (*userService)(nil)

func NewUserService(repo repository.UserRepo, hasher security.Hasher) UserService {
	return &userService{
		repo:   repo,
		hasher: hasher,
	}
}

type RegisterUserParams struct {
	Email           string
	Password        string
}

func (s *userService) RegisterUser(ctx context.Context, params RegisterUserParams) (*model.User, error) {
	hash, err := s.hasher.Hash(params.Password)

	if err != nil {
		return nil, fmt.Errorf("hasher hash: %w", err)
	}

	newUser, err := s.repo.CreateUser(ctx, repository.CreateUserParams{Email: params.Email, PasswordHash: hash})

	if err != nil {
		return nil, fmt.Errorf("create user %s: %w", params.Email, err)
	}

	return newUser, nil
}
