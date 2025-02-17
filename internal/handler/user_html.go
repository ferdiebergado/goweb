package handler

import (
	"net/http"
)

type UserHTMLHandler struct {
	template *Template
}

func NewUserHTMLHandler(t *Template) *UserHTMLHandler {
	return &UserHTMLHandler{
		template: t,
	}
}

func (h *UserHTMLHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	h.template.Render(w, r, "register", nil)
}
