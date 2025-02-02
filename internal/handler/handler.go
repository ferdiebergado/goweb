package handler

import "net/http"

func HandleHello(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello world!"))

	if err != nil {
		http.Error(w, "an error occured.", http.StatusInternalServerError)
	}
}
