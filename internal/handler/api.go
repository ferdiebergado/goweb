package handler

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message,omitempty"`
}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := APIResponse{
		Message: "Hello world!",
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, "an error occured.", http.StatusInternalServerError)
	}
}
