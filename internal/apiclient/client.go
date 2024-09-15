// internal/apiclient/client.go
package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIClient handles API requests
type APIClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Client:  &http.Client{},
	}
}

// StartTracking starts a new time entry
func (api *APIClient) StartTracking(username, url, description string) (*TimeEntry, error) {
	payload := map[string]string{
		"username":    username,
		"url":         url,
		"description": description,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := api.Client.Post(fmt.Sprintf("%s/api/start", api.BaseURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make start request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("start request failed: %s", string(body))
	}

	var entry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode start response: %w", err)
	}

	return &entry, nil
}

// StopTracking stops an existing time entry by URL
func (api *APIClient) StopTracking(username, url string) (*TimeEntry, error) {
	payload := map[string]string{
		"username": username,
		"url":      url,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := api.Client.Post(fmt.Sprintf("%s/api/stop", api.BaseURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make stop request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stop request failed: %s", string(body))
	}

	var entry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode stop response: %w", err)
	}

	return &entry, nil
}

// GetEntries retrieves time entries, optionally filtered by username and URL
// If format is "csv", it returns the data as CSV bytes
func (api *APIClient) GetEntries(username, url, format string) ([]byte, error) {
	query := "?"
	if username != "" {
		query += fmt.Sprintf("username=%s&", username)
	}
	if url != "" {
		query += fmt.Sprintf("url=%s&", url)
	}
	if format != "" {
		query += fmt.Sprintf("format=%s&", format)
	}
	// Remove trailing '&' or '?' if present
	query = strings.TrimRight(query, "&")
	if query == "?" {
		query = ""
	}

	reqURL := fmt.Sprintf("%s/api/entries%s", api.BaseURL, query)
	resp, err := api.Client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make entries request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("entries request failed: %s", string(body))
	}

	if strings.ToLower(format) == "csv" {
		// Read CSV bytes directly
		csvBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV response: %w", err)
		}
		return csvBytes, nil
	}

	// Default to JSON
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON response: %w", err)
	}

	return jsonBytes, nil
}

// TimeEntry represents a time tracking entry
type TimeEntry struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	StartTime   string  `json:"start_time"`
	EndTime     *string `json:"end_time,omitempty"`
	Duration    int64   `json:"duration"` // in minutes
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
