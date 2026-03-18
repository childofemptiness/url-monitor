package http

import (
	"net/http"

	"url-monitor/internal/monitor"
)

type Handler struct {
	service *monitor.Service
}

func NewHandler(service *monitor.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *Handler) ListMonitors(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *Handler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	//
}
