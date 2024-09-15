// internal/server/db/db.go
package db

import (
	"fmt"

	"github.com/xmonader/team-timetracker/internal/models"
	"github.com/xmonader/team-timetracker/internal/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDatabase sets up the database connection and performs migrations
func InitializeDatabase(cfg *server.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.DataSource)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Perform migrations
	if err := db.AutoMigrate(&models.TimeEntry{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
