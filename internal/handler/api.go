package handler

import (
	"encoding/json"
	"net/http"
)

func HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := struct {
		Message string
	}{
		Message: "Hello world!",
	}
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, "an error occured.", http.StatusInternalServerError)
	}
}
