// internal/server/handlers/tracking.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xmonader/team-timetracker/internal/models"
	"gorm.io/gorm"
)

// TrackingHandler handles start and stop tracking
type TrackingHandler struct {
	DB *gorm.DB
}

// StartTrackingRequest represents the request payload for starting tracking
type StartTrackingRequest struct {
	Username    string `json:"username"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// StartTracking handles the initiation of a time entry
func (h *TrackingHandler) StartTracking(w http.ResponseWriter, r *http.Request) {
	var req StartTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.URL) == "" || strings.TrimSpace(req.Description) == "" {
		http.Error(w, "Username, Description and URL are required", http.StatusBadRequest)
		return
	}

	// Optional: Check if there's already an active entry for this Username and URL
	var activeEntry models.TimeEntry
	if err := h.DB.Where("username = ? AND url = ? AND end_time IS NULL", req.Username, req.URL).First(&activeEntry).Error; err == nil {
		http.Error(w, "An active time entry for this URL already exists", http.StatusBadRequest)
		return
	}

	entry := models.TimeEntry{
		Username:    req.Username,
		URL:         req.URL,
		Description: req.Description,
		StartTime:   time.Now(),
	}

	if err := h.DB.Create(&entry).Error; err != nil {
		http.Error(w, "Failed to create time entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(entry)
}

// StopTrackingRequest represents the request payload for stopping tracking
type StopTrackingRequest struct {
	Username string `json:"username"`
	URL      string `json:"url"`
}

// StopTracking handles stopping a time entry
func (h *TrackingHandler) StopTracking(w http.ResponseWriter, r *http.Request) {
	var req StopTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.URL) == "" {
		http.Error(w, "Username and URL are required", http.StatusBadRequest)
		return
	}

	var entry models.TimeEntry
	// Find the latest active time entry for this Username and URL
	if err := h.DB.Where("username = ? AND url = ? AND end_time IS NULL", req.Username, req.URL).
		Order("start_time desc").
		First(&entry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, fmt.Sprintf("No active time entry found for the provided URL %s", req.URL), http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving time entry", http.StatusInternalServerError)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(entry.StartTime).Minutes()

	entry.EndTime = &endTime
	entry.Duration = int64(duration)

	if err := h.DB.Save(&entry).Error; err != nil {
		http.Error(w, "Failed to update time entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entry)
}

// GetStatus handles the /api/status endpoint
func (h *TrackingHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var activeEntry models.TimeEntry
	if err := h.DB.Where("username = ? AND end_time IS NULL", username).First(&activeEntry).Error; err != nil {
		// No active entry
		status := struct {
			Status string `json:"status"`
		}{
			Status: "idle",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	// Active entry exists
	status := struct {
		Status string `json:"status"`
		ID     uint   `json:"id"`
		URL    string `json:"url"`
	}{
		Status: "active",
		ID:     activeEntry.ID,
		URL:    activeEntry.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
