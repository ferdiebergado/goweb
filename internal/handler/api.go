package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/go-playground/validator/v10"
)

const jsonCT = "application/json"

type validationErrors map[string]string

type APIResponse[T any] struct {
	Message string           `json:"message"`
	Errors  validationErrors `json:"errors,omitempty"`
	Data    T                `json:"data,omitempty"`
}

func validationError(w http.ResponseWriter, r *http.Request, err error) {
	errs := make(map[string]string, 0)

	for _, e := range err.(validator.ValidationErrors) {
		errs[e.Field()] = e.Tag()
	}

	res := APIResponse[any]{
		Message: "Invalid input.",
		Errors:  errs,
	}

	response.JSON(w, r, http.StatusBadRequest, res)
}

func badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	const status = http.StatusBadRequest
	slog.Error("request error", "reason", err, "request", fmt.Sprint(r))

	if r.Header.Get("Content-Type") == jsonCT {
		res := APIResponse[any]{
			Message: "Invalid input",
		}
		response.JSON(w, r, http.StatusBadRequest, res)
		return
	}
	http.Error(w, http.StatusText(status), status)
}
