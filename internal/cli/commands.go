// internal/cli/commands.go
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/xmonader/team-timetracker/internal/apiclient"
)

// StartCommand handles the 'start' command.
// It takes the URL as an argument and a description.
func StartCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, url string, description string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("empty url")
	}
	if strings.TrimSpace(description) == "" {
		return fmt.Errorf("empty description")
	}

	entry, err := api.StartTracking(cfg.Username, url, description)
	if err != nil {
		return fmt.Errorf("starting tracking: %w", err)
	}

	fmt.Printf("Started tracking ID %d for URL '%s'.\n", entry.ID, entry.URL)
	return nil
}

// StopCommand handles the 'stop' command.
// It takes the URL of the time entry as a string.
func StopCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, url string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("empty url")
	}

	entry, err := api.StopTracking(cfg.Username, url)
	if err != nil {
		return fmt.Errorf("error stopping tracking: %w", err)
	}

	fmt.Printf("Stopped tracking for URL '%s'. Duration: %d minutes.\n", entry.URL, entry.Duration)
	return nil
}

// EntriesCommand handles the 'entries' command.
// It takes the desired output format, and optional username and url for filtering.
func EntriesCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, format string, username string, url string) error {
	if strings.TrimSpace(format) == "" {
		format = "json" // Default format
	}

	data, err := api.GetEntries(username, url, format)
	if err != nil {
		return fmt.Errorf("error retrieving entries: %w", err)
	}

	if strings.ToLower(format) == "csv" {
		// Write CSV to stdout
		fmt.Println(string(data))
		return nil
	}

	// Assume JSON format
	var entries []apiclient.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("error parsing JSON response: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No time entries found.")
		return nil
	}

	fmt.Println("Time Entries:")
	fmt.Println("------------------------------")
	for _, entry := range entries {
		fmt.Printf("ID: %d\n", entry.ID)
		fmt.Printf("Username: %s\n", entry.Username)
		fmt.Printf("URL: %s\n", entry.URL)
		fmt.Printf("Description: %s\n", entry.Description)
		fmt.Printf("Start Time: %s\n", formatTime(entry.StartTime))
		if entry.EndTime != nil {
			fmt.Printf("End Time: %s\n", formatTime(*entry.EndTime))
			fmt.Printf("Duration: %d minutes\n", entry.Duration)
		} else {
			fmt.Println("End Time: ---")
			fmt.Println("Duration: ---")
		}
		fmt.Println("------------------------------")
	}
	return nil
}

// formatTime formats time strings from RFC3339 to RFC1123
func formatTime(t string) string {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}
	return parsedTime.Format(time.RFC1123)
}
