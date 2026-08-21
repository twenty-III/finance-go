package db

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations runs all pending database migrations using goose.
func RunMigrations(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose: failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Goose: migration failed: %v", err)
	}

	log.Println("Goose: migrations completed successfully")
}

// RollbackMigrations rolls back all database migrations using goose.
func RollbackMigrations(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose: failed to set dialect: %v", err)
	}

	if err := goose.DownTo(db, "migrations", 0); err != nil {
		log.Fatalf("Goose: migration rollback failed: %v", err)
	}

	log.Println("Goose: migrations rolled back successfully")
}
