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

	return r
}
