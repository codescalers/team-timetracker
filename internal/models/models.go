package models

import (
    "time"
)

// TimeEntry represents a time tracking entry
type TimeEntry struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    Username    string     `json:"username"`
    URL         string     `json:"url"`
    Description string     `json:"description"`
    StartTime   time.Time  `json:"start_time"`
    EndTime     *time.Time `json:"end_time,omitempty"`
    Duration    int64      `json:"duration"` // in minutes
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
