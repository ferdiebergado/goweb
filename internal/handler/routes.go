package handler

import (
	"github.com/ferdiebergado/goexpress"
	"github.com/go-playground/validator/v10"
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

func mountUserRoutes(r *goexpress.Router, h *UserHandler, v *validator.Validate) {
	r.Post("/api/auth/register", h.HandleUserRegister,
		DecodeJSON[RegisterUserRequest](), ValidateInput[RegisterUserRequest](v))
}
