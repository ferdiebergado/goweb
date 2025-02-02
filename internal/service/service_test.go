package service_test

import (
	"context"
	"testing"

	"github.com/ferdiebergado/goweb/internal/repository/mock"
	"github.com/ferdiebergado/goweb/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDBVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockDBRepo(ctrl)
	ctx := context.Background()
	mockRepo.EXPECT().DBVersion(ctx).Return("1", nil)
	mockService := service.NewDBService(mockRepo)

	v, err := mockService.DBVersion(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, v)
}
