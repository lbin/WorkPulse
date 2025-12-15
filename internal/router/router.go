package router

import (
	"github.com/gin-gonic/gin"
	"workpulse/internal/handler"
	"workpulse/internal/middleware"
)

type Deps struct {
	JWTSecret       string
	WorkItemHandler *handler.WorkItemHandler
	OKRHandler      *handler.OKRHandler
	ReportHandler   *handler.ReportHandler
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", handler.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext())

	api.POST("/work-items", d.WorkItemHandler.Create)
	api.POST("/work-items/:id/okr-links", d.WorkItemHandler.AddOKRLink)

	api.GET("/okrs/objectives", d.OKRHandler.ListObjectives)
	api.POST("/okrs/objectives", d.OKRHandler.CreateObjective)
	api.PUT("/okrs/objectives/:id", d.OKRHandler.UpdateObjective)
	api.POST("/okrs/objectives/:id/archive", d.OKRHandler.ArchiveObjective)

	api.POST("/okrs/key-results", d.OKRHandler.CreateKeyResult)
	api.PUT("/okrs/key-results/:id", d.OKRHandler.UpdateKeyResult)
	api.POST("/okrs/key-results/:id/archive", d.OKRHandler.ArchiveKeyResult)
	api.POST("/okrs/key-results/:id/progress", d.OKRHandler.UpdateProgress)

	api.POST("/okrs/links", d.OKRHandler.AddLink)
	api.DELETE("/okrs/links/:id", d.OKRHandler.RemoveLink)

	api.GET("/reports", d.ReportHandler.List)
	api.GET("/reports/export", d.ReportHandler.Export)
	api.GET("/reports/:id", d.ReportHandler.Get)
	api.POST("/reports", d.ReportHandler.Create)
	api.PUT("/reports/:id", d.ReportHandler.Update)
	api.POST("/reports/:id/submit", d.ReportHandler.Submit)
	api.POST("/reports/:id/approve", d.ReportHandler.Approve)
	api.POST("/reports/:id/reject", d.ReportHandler.Reject)

	return r
}
