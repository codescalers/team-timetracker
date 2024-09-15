// cmd/server/main.go
package main

import (
	"flag"
	"log"
	"os"

	"github.com/xmonader/team-timetracker/internal/server"
	"github.com/xmonader/team-timetracker/internal/server/db"
)

func main() {
	// Define a flag for the config file path
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Check if config file exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file not found: %s\n", *configPath)
	}

	// Load configuration
	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %w", err)
	}

	database, err := db.InitializeDatabase(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %w", err)
	}

	srv := server.NewServer(database)

	if err := srv.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("Error starting server: %w", err)
	}
}
