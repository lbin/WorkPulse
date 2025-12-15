package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"workpulse/internal/config"
	"workpulse/internal/db"
	"workpulse/internal/handler"
	"workpulse/internal/middleware"
	"workpulse/internal/migration"
	"workpulse/internal/observability"
	"workpulse/internal/repo"
	"workpulse/internal/router"
	"workpulse/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	obs, err := observability.New(ctx, observability.Config{
		ServiceName:  cfg.ServiceName,
		Environment:  cfg.Env,
		OTLPEndpoint: cfg.OTLPEndpoint,
		OTLPHeaders:  cfg.OTLPHeaders,
		MetricsPath:  cfg.MetricsPath,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = obs.Shutdown(context.Background())
	}()

	gdb, err := db.Open(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Run(gdb, cfg.MigrationsTable); err != nil {
		log.Fatal(err)
	}

	wiRepo := repo.NewWorkItemRepo(gdb)
	wiSvc := service.NewWorkItemService(wiRepo)
	wiHandler := handler.NewWorkItemHandler(wiSvc)

	projectSvc := service.NewProjectService()
	projectHandler := handler.NewProjectHandler(projectSvc)

	auditRepo := repo.NewAuditRepo(gdb)
	rbacRepo := repo.NewRBACRepo(gdb)
	reportRepo := repo.NewReportRepo(gdb)
	reportSvc := service.NewReportService(reportRepo, auditRepo)
	reportHandler := handler.NewReportHandler(reportSvc)
	authSvc := service.NewAuthService(rbacRepo, auditRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	configSvc := service.NewConfigService(cfg)
	configHandler := handler.NewConfigHandler(configSvc, fmt.Sprintf("/api/%s", cfg.APIVersion))

	okrSvc := service.NewOKRService()
	okrHandler := handler.NewOKRHandler(okrSvc)

	meetingSvc := service.NewMeetingService()
	meetingHandler := handler.NewMeetingHandler(meetingSvc)

	permMW := middleware.NewPermissionsLoader(authSvc)

	r := router.New(router.Deps{
		ServiceName:     cfg.ServiceName,
		JWTSecret:       cfg.JWTSecret,
		APIVersion:      cfg.APIVersion,
		MetricsPath:     cfg.MetricsPath,
		MetricsHandler:  obs.MetricsHandler,
		TelemetryMW:     obs.GinMiddleware(),
		WorkItemHandler: wiHandler,
		ProjectHandler:  projectHandler,
		OKRHandler:      okrHandler,
		ReportHandler:   reportHandler,
		MeetingHandler:  meetingHandler,
		AuthHandler:     authHandler,
		ConfigHandler:   configHandler,
		PermissionMW:    permMW,
	})

	log.Printf("listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
