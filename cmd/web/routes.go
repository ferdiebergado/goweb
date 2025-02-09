package main

import (
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/handler"
)

func mountRoutes(r *goexpress.Router, h *handler.BaseHandler) {
	r.Group("/api", func(gr *goexpress.Router) *goexpress.Router {
		gr.Get("/health", h.HandleHealth)
		return gr
	})
}

func mountBaseHTMLRoutes(r *goexpress.Router, h *handler.BaseHTMLHandler) {
	r.Get("/dashboard", h.HandleDashboard)
}

func mountUserRoutes(r *goexpress.Router, h *handler.UserHandler) {
	r.Post("/api/auth/register", h.HandleUserRegister)
}
