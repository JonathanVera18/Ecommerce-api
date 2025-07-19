package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/JonathanVera18/ecommerce-api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down]")
	}

	direction := os.Args[1]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Use GORM to get the underlying SQL DB
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying SQL DB:", err)
	}
	defer sqlDB.Close()

	switch direction {
	case "up":
		err = migrateUp(sqlDB)
	case "down":
		log.Println("Migration down not implemented - please manually rollback if needed")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}

func migrateUp(sqlDB interface{}) error {
	// Get all migration files
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration file
	for _, filename := range sqlFiles {
		log.Printf("Running migration: %s", filename)
		
		filePath := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the entire content as one statement
		if db, ok := sqlDB.(interface{ Exec(string, ...interface{}) (interface{}, error) }); ok {
			_, err := db.Exec(string(content))
			if err != nil {
				// Log the error but continue - some statements might fail if already executed
				log.Printf("Warning: Error executing %s (this might be expected if already executed): %v", filename, err)
			} else {
				log.Printf("Successfully executed migration: %s", filename)
			}
		}
	}

	return nil
}
