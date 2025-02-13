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
		errs[e.Field()] = getValidationMessage(e)
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
		response.JSON(w, r, status, res)
		return
	}
	http.Error(w, http.StatusText(status), status)
}

func unprocessableError(w http.ResponseWriter, r *http.Request, err error) {
	const status = http.StatusUnprocessableEntity
	slog.Error("request error", "reason", err, "request", fmt.Sprint(r))

	if r.Header.Get("Content-Type") == jsonCT {
		res := APIResponse[any]{
			Message: err.Error(),
		}
		response.JSON(w, r, status, res)
		return
	}
	http.Error(w, http.StatusText(status), status)
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", e.Field())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", e.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
