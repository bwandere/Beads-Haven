package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
	"time"

	"beads-haven/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

// RunMigrations applies any pending embedded SQL migrations in order within transactions.
func RunMigrations(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, filename := range files {
		var applied bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&applied)
		if err != nil {
			return fmt.Errorf("failed checking status of migration %s: %w", filename, err)
		}
		if applied {
			continue
		}

		content, err := fs.ReadFile(migrations.FS, filename)
		if err != nil {
			return fmt.Errorf("failed reading migration %s: %w", filename, err)
		}

		log.Printf("Applying migration: %s", filename)
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", filename, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for migration %s: %w", filename, err)
		}
		log.Printf("Applied migration %s successfully", filename)
	}

	return nil
}

// EnsureSchema is retained for compatibility and delegates directly to RunMigrations.
func EnsureSchema(db *sql.DB) error {
	return RunMigrations(db)
}
