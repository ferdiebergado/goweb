package handler

const jsonCT = "application/json"

type APIResponse[T any] struct {
	Message string           `json:"message"`
	Errors  validationErrors `json:"errors,omitempty"`
	Data    T                `json:"data,omitempty"`
}
