package handler

import (
	"net/http"
)

type Handler struct {
	Base BaseHandler
	User UserHandler
}

func NewHandler(tmpl *Template) *Handler {
	return &Handler{
		Base: *NewBaseHandler(tmpl),
		User: *NewUserHandler(tmpl),
	}
}

type BaseHandler struct {
	template *Template
}

func NewBaseHandler(t *Template) *BaseHandler {
	return &BaseHandler{
		template: t,
	}
}

func (h *BaseHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	h.template.Render(w, r, "dashboard", nil)
}

type UserHandler struct {
	template *Template
}

func NewUserHandler(t *Template) *UserHandler {
	return &UserHandler{
		template: t,
	}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	h.template.Render(w, r, "register", nil)
}
