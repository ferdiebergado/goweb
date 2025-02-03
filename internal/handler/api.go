package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ferdiebergado/gopherkit/http/response"
)

type APIResponse struct {
	Message string `json:"message"`
}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := APIResponse{
		Message: "Hello world!",
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		response.ServerError(w, r, fmt.Errorf("encode json: %w", err))
	}
}
