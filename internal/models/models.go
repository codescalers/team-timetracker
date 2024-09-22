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

// GetCurrentTime returns the current time in RFC3339 format
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// CalculateDuration calculates the duration between start and end times in minutes
func CalculateDuration(start, end string) int64 {
	layout := time.RFC3339
	startTime, err1 := time.Parse(layout, start)
	endTime, err2 := time.Parse(layout, end)
	if err1 != nil || err2 != nil {
		return 0
	}
	return int64(endTime.Sub(startTime).Minutes())
}
