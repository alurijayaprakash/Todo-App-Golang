package http

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

type HealthProbe struct {
	Status string `json:"status"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Ping returns a simple health check response.
// @Summary      Health check
// @Description  Returns service health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthProbe
// @Router       /health [get]
func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthProbe{Status: "ok"})
}
