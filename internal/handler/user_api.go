package handler

import (
	"net/http"
	"time"

	"github.com/ferdiebergado/gopherkit/http/request"
	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/ferdiebergado/goweb/internal/service"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   service.UserService
	validater *validator.Validate
}

func NewUserHandler(userService service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   userService,
		validater: v,
	}
}

type RegisterUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *UserHandler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[service.RegisterUserParams](r)
	if err != nil {
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}

	if err = h.validater.Struct(params); err != nil {
		validationError(w, r, err)
		return
	}

	user, err := h.service.RegisterUser(r.Context(), params)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	res := APIResponse[*RegisterUserResponse]{
		Message: "User registered.",
		Data: &RegisterUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	response.JSON(w, r, http.StatusCreated, res)
}
