package database

import (
	"database/sql"
	"embed"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// /go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate() error {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Read migration files
	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Sort files to ensure they run in order
	var upFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}
	sort.Strings(upFiles)

	// Execute each migration
	for _, filename := range upFiles {
		content, err := migrationFiles.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			log.Printf("Error reading migration file %s: %v", filename, err)
			continue
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Printf("Error executing migration %s: %v", filename, err)
			continue
		}

		log.Printf("Successfully executed migration: %s", filename)
	}

	return nil
}
