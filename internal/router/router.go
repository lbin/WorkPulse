package router

import (
	"github.com/gin-gonic/gin"
	"workpulse/internal/handler"
	"workpulse/internal/middleware"
)

type Deps struct {
	JWTSecret       string
	AuthHandler     *handler.AuthHandler
	TeamHandler     *handler.TeamHandler
	WorkItemHandler *handler.WorkItemHandler
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", handler.Health)

	// Public auth endpoints for registration/login.
	r.POST("/api/v1/auth/register", d.AuthHandler.Register)
	r.POST("/api/v1/auth/login", d.AuthHandler.Login)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext())

	api.GET("/teams", d.TeamHandler.List)
	api.GET("/teams/my", d.TeamHandler.ListMine)
	api.POST("/teams", d.TeamHandler.Create)
	api.POST("/teams/join", d.TeamHandler.Join)

	api.POST("/work-items", d.WorkItemHandler.Create)
	api.POST("/work-items/:id/okr-links", d.WorkItemHandler.AddOKRLink)

	return r
}
