package handler

import (
	"net/http"

	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/go-playground/validator/v10"
)

type validationErrors map[string][]string

type APIResponse[T any] struct {
	Message string           `json:"message"`
	Errors  validationErrors `json:"errors,omitempty"`
	Data    T                `json:"data,omitempty"`
}

func validationError(w http.ResponseWriter, r *http.Request, err error) {
	errs := make(map[string][]string, 0)

	for _, e := range err.(validator.ValidationErrors) {
		errs[e.Field()] = append(errs[e.Field()], e.Tag())
	}

	res := APIResponse[any]{
		Message: "Invalid input.",
		Errors:  errs,
	}

	response.JSON(w, r, http.StatusBadRequest, res)
}
