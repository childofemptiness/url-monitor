package http

import (
	"encoding/json"
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
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) ListMonitors(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *Handler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	//
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
