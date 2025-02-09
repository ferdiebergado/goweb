package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/model"
	"github.com/ferdiebergado/goweb/internal/service"
	"github.com/ferdiebergado/goweb/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type RegisterUserRequest struct {
}

func TestUserHandler_HandleUserRegister(t *testing.T) {
	const (
		url            = "/api/auth/register"
		msg            = "User registered."
		ct             = "application/json"
		testEmail      = "abc@example.com"
		testPass       = "test"
		testPassHashed = "hashed"
	)

	ctrl := gomock.NewController(t)
	mockService := mock.NewMockUserService(ctrl)
	params := service.RegisterUserParams{
		Email:           testEmail,
		Password:        testPass,
		PasswordConfirm: testPass,
	}

	user := &model.User{
		Model: model.Model{ID: "1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Email: testEmail,
	}

	mockService.EXPECT().RegisterUser(context.Background(), params).Return(user, nil)
	userHandler := handler.NewUserHandler(mockService)
	r := goexpress.New()
	r.Post(url, userHandler.HandleUserRegister)

	paramsJSON, err := json.Marshal(params)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", url, bytes.NewBuffer(paramsJSON))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, ct, res.Header["Content-Type"][0])

	var apiRes handler.APIResponse[handler.RegisterUserResponse]
	if err := json.Unmarshal(rr.Body.Bytes(), &apiRes); err != nil {
		t.Fatal("failed to decode json", err)
	}

	assert.Equal(t, msg, apiRes.Message)
	assert.Equal(t, user.ID, apiRes.Data.ID)
	assert.Equal(t, user.Email, apiRes.Data.Email)
	assert.NotZero(t, apiRes.Data.CreatedAt)
	assert.NotZero(t, apiRes.Data.UpdatedAt)
}
