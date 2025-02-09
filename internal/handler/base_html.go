package handler

import (
	"net/http"
)

type BaseHTMLHandler struct {
	template *Template
}

func NewBaseHTMLHandler(t *Template) *BaseHTMLHandler {
	return &BaseHTMLHandler{
		template: t,
	}
}

func (h *BaseHTMLHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	h.template.Render(w, r, "dashboard", nil)
}
