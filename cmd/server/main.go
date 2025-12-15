package main

import (
	"log"

	"workpulse/internal/config"
	"workpulse/internal/db"
	"workpulse/internal/handler"
	"workpulse/internal/repo"
	"workpulse/internal/router"
	"workpulse/internal/service"
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

	wiRepo := repo.NewWorkItemRepo(gdb)
	wiSvc := service.NewWorkItemService(wiRepo)
	wiHandler := handler.NewWorkItemHandler(wiSvc)

	r := router.New(router.Deps{
		JWTSecret:       cfg.JWTSecret,
		WorkItemHandler: wiHandler,
	})

	log.Printf("listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
