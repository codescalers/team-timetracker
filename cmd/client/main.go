// cmd/client/main.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xmonader/team-timetracker/internal/apiclient"
	"github.com/xmonader/team-timetracker/internal/cli"
)

func main() {
	// Define a flag for the config file path
	configPath := flag.String("config", "", "Path to configuration file")

	// Parse global flags first
	flag.Parse()

	// Determine config file path
	var cfgPath string
	if *configPath != "" {
		cfgPath = *configPath
	} else {
		cfgPath = apiclient.DefaultConfigPath()
	}

	// Check if config file exists
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("configuration file not found: %s\n", cfgPath)
	}

	// Load configuration
	cfg, err := apiclient.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	api := apiclient.NewAPIClient(cfg.BackendURL)

	// Retrieve non-flag arguments
	args := flag.Args()

	// Ensure at least one subcommand is provided
	if len(args) < 1 {
		fmt.Println("error: No subcommand provided.")
		printUsage()
		os.Exit(1)
	}

	// Parse subcommand
	subcommand := args[0]
	switch subcommand {
	case "start":
		err = handleStartCommand(cfg, api, args[1:])
	case "stop":
		err = handleStopCommand(cfg, api, args[1:])
	case "entries":
		err = handleEntriesCommand(cfg, api, args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// handleStartCommand parses and executes the 'start' subcommand
func handleStartCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, args []string) error {
	// Define a new FlagSet for the 'start' command
	startFlagSet := flag.NewFlagSet("start", flag.ExitOnError)
	// Parse the flags specific to the 'start' command
	err := startFlagSet.Parse(args)
	if err != nil {
		log.Fatalf("error parsing start command flags: %s\n", err)
	}

	remainingArgs := startFlagSet.Args()
	if len(remainingArgs) < 1 {
		return fmt.Errorf("url and description are required for the start command")
	}
	url := remainingArgs[0]
	description := remainingArgs[1]

	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("url is required for the start command")
	}

	if strings.TrimSpace(description) == "" {
		return fmt.Errorf("description is required for the start command")
	}
	return cli.StartCommand(cfg, api, url, description)
}

// handleStopCommand parses and executes the 'stop' subcommand
func handleStopCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, args []string) error {
	stopFlagSet := flag.NewFlagSet("stop", flag.ExitOnError)

	if err := stopFlagSet.Parse(args); err != nil {
		return err
	}

	// After parsing, the remaining args should include the URL
	remainingArgs := stopFlagSet.Args()
	if len(remainingArgs) < 1 {
		return errors.New("url is required for the stop command")
	}
	url := remainingArgs[0]

	return cli.StopCommand(cfg, api, url)
}

// handleEntriesCommand parses and executes the 'entries' subcommand
func handleEntriesCommand(cfg *apiclient.ClientConfig, api *apiclient.APIClient, args []string) error {
	// Define a new FlagSet for the 'entries' command
	entriesFlagSet := flag.NewFlagSet("entries", flag.ExitOnError)
	format := entriesFlagSet.String("format", "json", "Output format: json or csv")
	username := entriesFlagSet.String("username", "", "Filter entries by username")
	urlFilter := entriesFlagSet.String("url", "", "Filter entries by URL")

	// Parse the flags specific to the 'entries' command
	if err := entriesFlagSet.Parse(args); err != nil {
		return err
	}

	// Execute the entries command
	return cli.EntriesCommand(cfg, api, *format, *username, *urlFilter)
}

// printUsage displays the usage information for the CLI
func printUsage() {
	usage := `Usage:
    timetracker [command] [options]

Commands:
    start [URL] ["Your description"]   Start a new time entry
    stop [URL]                                      Stop an existing time entry by URL
    entries [--username="username"] [--url="URL"] [--format=csv]  List all time entries with optional filtering and format

Examples:
    timetracker start https://github.com/yourrepo/issue/1 "Working on feature X"
    timetracker stop https://github.com/yourrepo/issue/1
    timetracker entries --username="xmonader" --url="https://github.com/yourrepo/issue/1" --format=csv
    timetracker entries --format=csv
`
	fmt.Println(usage)
}
