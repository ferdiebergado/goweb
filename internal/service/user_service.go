//go:generate mockgen -destination=mock/user_service_mock.go -package=mock . UserService
package service

import (
	"context"
	"fmt"

	"github.com/ferdiebergado/goweb/internal/model"
	"github.com/ferdiebergado/goweb/internal/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, params RegisterUserParams) (*model.User, error)
}

type userService struct {
	repo repository.UserRepo
}

var _ UserService = (*userService)(nil)

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

type RegisterUserParams struct {
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"password_confirm,omitempty"`
}

func (s *userService) RegisterUser(ctx context.Context, params RegisterUserParams) (*model.User, error) {
	newUser, err := s.repo.CreateUser(ctx, repository.CreateUserParams{Email: params.Email, PasswordHash: params.Password})

	if err != nil {
		return nil, fmt.Errorf("create user %s: %w", params.Email, err)
	}

	return newUser, nil
}
