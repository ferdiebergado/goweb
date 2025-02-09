package handler

import (
	"net/http"
	"time"

	"github.com/ferdiebergado/gopherkit/http/request"
	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/ferdiebergado/goweb/internal/service"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{service: userService}
}

type RegisterUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *userHandler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[service.RegisterUserParams](r)
	if err != nil {
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}

	// TODO: validate params
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
