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
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	regUrl         = "/api/auth/register"
	ct             = "application/json"
	testEmail      = "abc@example.com"
	testPass       = "test"
	testPassHashed = "hashed"
)

var validate *validator.Validate

func TestMain(t *testing.M) {
	validate = validator.New()
	t.Run()
}

func TestUserHandler_HandleUserRegister_Success(t *testing.T) {
	const msg = "User registered."

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

	mockService.EXPECT().RegisterUser(handler.NewParamsContext(context.Background(), params), params).Return(user, nil)
	userHandler := handler.NewUserHandler(mockService)
	r := goexpress.New()
	r.Post(regUrl, userHandler.HandleUserRegister,
		handler.DecodeJSON[service.RegisterUserParams](), handler.ValidateInput[service.RegisterUserParams](validate))

	paramsJSON, err := json.Marshal(params)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", regUrl, bytes.NewBuffer(paramsJSON))
	req.Header.Set("Content-Type", "application/json")
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

func TestUserHandler_HandleUserRegister_InvalidInput(t *testing.T) {
	const msg = "Invalid input."

	ctrl := gomock.NewController(t)
	mockService := mock.NewMockUserService(ctrl)
	userHandler := handler.NewUserHandler(mockService)
	r := goexpress.New()
	r.Post(regUrl, userHandler.HandleUserRegister,
		handler.DecodeJSON[service.RegisterUserParams](), handler.ValidateInput[service.RegisterUserParams](validate))

	var tests = []struct {
		name   string
		params service.RegisterUserParams
	}{
		{"Empty email", service.RegisterUserParams{Password: testPass, PasswordConfirm: testPass}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0)
			paramsJSON, err := json.Marshal(tt.params)

			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest("POST", regUrl, bytes.NewBuffer(paramsJSON))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, ct, res.Header["Content-Type"][0])
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			var apiRes handler.APIResponse[handler.RegisterUserResponse]
			if err := json.Unmarshal(rr.Body.Bytes(), &apiRes); err != nil {
				t.Fatal("failed to decode json", err)
			}

			assert.Equal(t, msg, apiRes.Message)
		})
	}
}
