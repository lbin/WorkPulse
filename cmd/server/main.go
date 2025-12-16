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

	// repositories
	orgRepo := repo.NewOrgRepo(gdb)
	userRepo := repo.NewUserRepo(gdb)
	teamRepo := repo.NewTeamRepo(gdb)
	membershipRepo := repo.NewMembershipRepo(gdb)
	wiRepo := repo.NewWorkItemRepo(gdb)

	// services
	authSvc := service.NewAuthService(userRepo, orgRepo, cfg.JWTSecret)
	teamSvc := service.NewTeamService(teamRepo, membershipRepo)
	wiSvc := service.NewWorkItemService(wiRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	wiHandler := handler.NewWorkItemHandler(wiSvc)

	r := router.New(router.Deps{
		JWTSecret:       cfg.JWTSecret,
		AuthHandler:     authHandler,
		TeamHandler:     teamHandler,
		WorkItemHandler: wiHandler,
	})

	log.Printf("listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
