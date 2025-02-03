package main

import (
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/handler"
)

func mountRoutes(r *goexpress.Router, handler *handler.BaseHandler) {
	r.Get("/api/health", handler.HandleHealth)
}
