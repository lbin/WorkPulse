package router

import (
	"github.com/gin-gonic/gin"
	"workpulse/internal/handler"
	"workpulse/internal/middleware"
)

type Deps struct {
	JWTSecret       string
	WorkItemHandler *handler.WorkItemHandler
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", handler.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext())

	api.POST("/work-items", d.WorkItemHandler.Create)
	api.POST("/work-items/:id/okr-links", d.WorkItemHandler.AddOKRLink)

	return r
}
