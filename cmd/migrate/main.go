package main

import (
	"log"

	"workpulse/internal/config"
	"workpulse/internal/db"
	"workpulse/internal/migration"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gdb, err := db.Open(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Run(gdb, cfg.MigrationsTable); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
