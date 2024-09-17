// cmd/server/main.go
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/xmonader/team-timetracker/internal/server"
	"github.com/xmonader/team-timetracker/internal/server/db"
)

func main() {
	// Define a flag for the config file path
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if strings.TrimSpace(*configPath) == "" {
		log.Fatal("configuration file path is required")
	}
	// Check if config file exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("configuration file not found: %s\n", *configPath)
	}

	// Load configuration
	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	database, err := db.InitializeDatabase(cfg)
	if err != nil {
		log.Fatalf("error initializing database: %s", err)
	}

	srv := server.NewServer(database)

	if err := srv.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
