// internal/server/handlers/entries.go
package handlers

import (
    "encoding/csv"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/xmonader/team-timetracker/internal/models"
    "gorm.io/gorm"
)

// EntriesHandler handles retrieving time entries
type EntriesHandler struct {
    DB *gorm.DB
}

// GetEntries handles retrieving time entries, with optional filtering and format
func (h *EntriesHandler) GetEntries(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    url := r.URL.Query().Get("url")
    format := r.URL.Query().Get("format") // e.g., "csv"

    var entries []models.TimeEntry
    query := h.DB.Model(&models.TimeEntry{})

    if username != "" {
        query = query.Where("username = ?", username)
    }

    if url != "" {
        query = query.Where("url = ?", url)
    }

    if err := query.Order("start_time desc").Find(&entries).Error; err != nil {
        http.Error(w, "Failed to retrieve entries", http.StatusInternalServerError)
        return
    }

    if strings.ToLower(format) == "csv" {
        // Set CSV headers
        w.Header().Set("Content-Type", "text/csv")
        w.Header().Set("Content-Disposition", "attachment;filename=entries.csv")

        writer := csv.NewWriter(w)
        defer writer.Flush()

        // Write CSV headers
        writer.Write([]string{"ID", "Username", "URL", "Description", "Start Time", "End Time", "Duration (minutes)", "Created At", "Updated At"})

        // Write entries
        for _, entry := range entries {
            endTime := "N/A"
            if entry.EndTime != nil {
                endTime = entry.EndTime.Format(time.RFC3339)
            }

            record := []string{
                strconv.Itoa(int(entry.ID)),
                entry.Username,
                entry.URL,
                entry.Description,
                entry.StartTime.Format(time.RFC3339),
                endTime,
                strconv.FormatInt(entry.Duration, 10),
                entry.CreatedAt.Format(time.RFC3339),
                entry.UpdatedAt.Format(time.RFC3339),
            }

            writer.Write(record)
        }
        return
    }

    // Default to JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(entries)
}
