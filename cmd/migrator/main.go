package main

import (
	"echoFramework/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config is required: %v", err)
	}

	log.Printf("try to migrate from %s to %s\n", cfg.MigrationsPath, cfg.DBPath)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.DBPath, cfg.MigrationsTable),
	)

	if err != nil {
		log.Fatalf("can't create migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}

	log.Println("migrations successfully applied")
}
