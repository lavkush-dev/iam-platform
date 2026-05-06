package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"iam-platform/internal/config"
)

func main() {
	cfg := config.Load()

	action := flag.String("action", "up", "migration action: up | down | force")
	flag.Parse()

	m, err := migrate.New(
		"file://db/migrations",
		cfg.DB.URL,
	)
	if err != nil {
		log.Fatal("migration init failed:", err)
	}

	switch *action {

	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migration up failed:", err)
		}
		log.Println("migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("migration down failed:", err)
		}
		log.Println("migrations rolled back")

	case "force":
		version := 1 // change if needed
		if err := m.Force(version); err != nil {
			log.Fatal("migration force failed:", err)
		}
		log.Println("migration force applied")

	default:
		log.Fatal("invalid action")
	}
}
