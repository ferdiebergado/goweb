package handler

import (
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/service"
)

func mountRoutes(r *goexpress.Router, h *BaseHandler) {
	r.Group("/api", func(gr *goexpress.Router) *goexpress.Router {
		gr.Get("/health", h.HandleHealth)
		return gr
	})
}

func mountBaseHTMLRoutes(r *goexpress.Router, h *BaseHTMLHandler) {
	r.Get("/dashboard", h.HandleDashboard)
}

func mountUserRoutes(r *goexpress.Router, h *UserHandler) {
	r.Post("/api/auth/register", h.HandleUserRegister,
		DecodeJSON[service.RegisterUserParams](), ValidateInput[service.RegisterUserParams](h.validater))
}
