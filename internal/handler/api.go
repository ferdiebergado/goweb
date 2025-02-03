package handler

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/ferdiebergado/goweb/internal/service"
)

type BaseHandler struct {
	service service.Service
}

func NewBaseHandler(service service.Service) *BaseHandler {
	return &BaseHandler{service: service}
}

type APIResponse struct {
	Message string `json:"message"`
}

func (h *BaseHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	msg := "healthy"

	if err := h.service.PingDB(r.Context()); err != nil {
		status = http.StatusServiceUnavailable
		msg = "unhealthy"
		slog.Error("failed to connect to the database", "reason", err)
	}

	response.JSON(w, r, status, APIResponse{Message: msg})
}
