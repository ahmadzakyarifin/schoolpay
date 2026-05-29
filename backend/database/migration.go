package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *bun.DB) error {
	log.Println("Memulai proses migrasi database...")

	// Set database dialect
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("gagal mengatur dialect goose: %w", err)
	}

	// Tell goose to use our embedded files
	goose.SetBaseFS(embedMigrations)

	// Run migrations
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return fmt.Errorf("gagal menjalankan migrasi up: %w", err)
	}

	log.Println("Migrasi database selesai dengan sukses!")
	return nil
}
