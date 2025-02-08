package handler

import (
	"net/http"

	"github.com/ferdiebergado/goweb/internal/pkg/template"
)

type BaseHTMLHandler struct {
	template *template.Template
}

func NewBaseHTMLHandler(t *template.Template) *BaseHTMLHandler {
	return &BaseHTMLHandler{
		template: t,
	}
}

func (h *BaseHTMLHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	h.template.Render(w, r, "dashboard", nil)
}
