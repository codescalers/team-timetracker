// internal/server/handlers/health.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// LivenessHandler handles the /livez endpoint
func (h *HealthHandler) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HealthHandler handles the /healthz endpoint
func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	type HealthStatus struct {
		Status   string `json:"status"`
		Database string `json:"database"`
		Time     string `json:"time"`
	}

	status := HealthStatus{
		Status: "ok",
		Time:   time.Now().Format(time.RFC3339),
	}

	// Check database connectivity
	dbHandle, err := h.DB.DB()
	if err != nil {
		status.Status = "error"
		status.Database = "unreachable"
	}

	defer dbHandle.Close() // safe it's connectionpooled
	if err := dbHandle.Ping(); err != nil {
		status.Status = "error"
		status.Database = "unreachable"
	} else {
		status.Database = "ok"
	}

	// Add more checks here if necessary

	if status.Status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// HealthHandler encapsulates dependencies for health checks
type HealthHandler struct {
	DB *gorm.DB
}
