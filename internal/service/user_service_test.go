package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ferdiebergado/goweb/internal/model"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/ferdiebergado/goweb/internal/repository/mock"
	"github.com/ferdiebergado/goweb/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_RegisterUser(t *testing.T) {
	const (
		testEmail      = "abc@example.com"
		testPass       = "test"
		testPassHashed = "hashed"
	)
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockUserRepo(ctrl)
	ctx := context.Background()

	regParams := service.RegisterUserParams{
		Email:           testEmail,
		Password:        testPass,
		PasswordConfirm: testPass,
	}

	params := repository.CreateUserParams{
		Email:        regParams.Email,
		PasswordHash: regParams.Password,
	}

	user := &model.User{
		Model: model.Model{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Email: testEmail,
	}

	mockRepo.EXPECT().CreateUser(ctx, params).Return(user, nil)

	userService := service.NewUserService(mockRepo)

	newUser, err := userService.RegisterUser(ctx, regParams)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.NotZero(t, newUser.ID)
	assert.Equal(t, params.Email, newUser.Email, "Emails must match")
	assert.NotZero(t, newUser.CreatedAt)
	assert.NotZero(t, newUser.UpdatedAt)
}
