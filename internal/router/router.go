package router

import (
	"github.com/gin-gonic/gin"
	"workpulse/internal/handler"
	"workpulse/internal/middleware"
)

type Deps struct {
	JWTSecret       string
	WorkItemHandler *handler.WorkItemHandler
	ProjectHandler  *handler.ProjectHandler
	OKRHandler      *handler.OKRHandler
	ReportHandler   *handler.ReportHandler
	MeetingHandler  *handler.MeetingHandler
	AuthHandler     *handler.AuthHandler
	PermissionMW    *middleware.PermissionsLoader
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", handler.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext(), d.PermissionMW.Handler())

	api.GET("/auth/me", d.AuthHandler.Me)

	api.POST("/work-items", middleware.RequirePermission("workitems.manage"), d.WorkItemHandler.Create)
	api.POST("/work-items/:id/okr-links", middleware.RequirePermission("workitems.manage"), d.WorkItemHandler.AddOKRLink)

	api.GET("/projects", middleware.RequirePermission("projects.view"), d.ProjectHandler.List)
	api.POST("/projects", middleware.RequirePermission("projects.manage"), d.ProjectHandler.Create)
	api.PUT("/projects/:id", middleware.RequirePermission("projects.manage"), d.ProjectHandler.Update)
	api.GET("/projects/:id/milestones", middleware.RequirePermission("projects.view"), d.ProjectHandler.ListMilestones)
	api.POST("/projects/:id/milestones", middleware.RequirePermission("projects.manage"), d.ProjectHandler.CreateMilestone)
	api.PUT("/milestones/:id", middleware.RequirePermission("projects.manage"), d.ProjectHandler.UpdateMilestone)
	api.GET("/projects/:id/tasks", middleware.RequirePermission("projects.view"), d.ProjectHandler.ListTasks)
	api.POST("/projects/:id/tasks", middleware.RequirePermission("projects.manage"), d.ProjectHandler.CreateTask)
	api.PUT("/tasks/:id", middleware.RequirePermission("projects.manage"), d.ProjectHandler.UpdateTask)
	api.POST("/projects/:id/links", middleware.RequirePermission("projects.manage"), d.ProjectHandler.AddProjectLink)
	api.DELETE("/projects/:projectId/links/:linkId", middleware.RequirePermission("projects.manage"), d.ProjectHandler.RemoveProjectLink)
	api.POST("/tasks/:id/links", middleware.RequirePermission("projects.manage"), d.ProjectHandler.AddTaskLink)
	api.DELETE("/tasks/:taskId/links/:linkId", middleware.RequirePermission("projects.manage"), d.ProjectHandler.RemoveTaskLink)
	api.GET("/tasks/:id/links", middleware.RequirePermission("projects.view"), d.ProjectHandler.ListTaskLinks)
	api.GET("/projects/:id/okr-progress", middleware.RequirePermission("projects.view"), d.ProjectHandler.OKRProgress)

	api.GET("/okrs/objectives", middleware.RequirePermission("okr.view"), d.OKRHandler.ListObjectives)
	api.POST("/okrs/objectives", middleware.RequirePermission("okr.manage"), d.OKRHandler.CreateObjective)
	api.PUT("/okrs/objectives/:id", middleware.RequirePermission("okr.manage"), d.OKRHandler.UpdateObjective)
	api.POST("/okrs/objectives/:id/archive", middleware.RequirePermission("okr.manage"), d.OKRHandler.ArchiveObjective)

	api.POST("/okrs/key-results", middleware.RequirePermission("okr.manage"), d.OKRHandler.CreateKeyResult)
	api.PUT("/okrs/key-results/:id", middleware.RequirePermission("okr.manage"), d.OKRHandler.UpdateKeyResult)
	api.POST("/okrs/key-results/:id/archive", middleware.RequirePermission("okr.manage"), d.OKRHandler.ArchiveKeyResult)
	api.POST("/okrs/key-results/:id/progress", middleware.RequirePermission("okr.manage"), d.OKRHandler.UpdateProgress)

	api.POST("/okrs/links", middleware.RequirePermission("okr.manage"), d.OKRHandler.AddLink)
	api.DELETE("/okrs/links/:id", middleware.RequirePermission("okr.manage"), d.OKRHandler.RemoveLink)

	api.GET("/okrs/metrics", middleware.RequirePermission("okr.view"), d.OKRHandler.Metrics)

	api.GET("/reports", middleware.RequirePermission("reports.view"), d.ReportHandler.List)
	api.GET("/reports/export", middleware.RequirePermission("reports.view"), d.ReportHandler.Export)
	api.GET("/reports/:id", middleware.RequirePermission("reports.view"), d.ReportHandler.Get)
	api.POST("/reports", middleware.RequirePermission("reports.manage"), d.ReportHandler.Create)
	api.PUT("/reports/:id", middleware.RequirePermission("reports.manage"), d.ReportHandler.Update)
	api.POST("/reports/:id/submit", middleware.RequirePermission("reports.manage"), d.ReportHandler.Submit)
	api.POST("/reports/:id/approve", middleware.RequirePermission("reports.review"), d.ReportHandler.Approve)
	api.POST("/reports/:id/reject", middleware.RequirePermission("reports.review"), d.ReportHandler.Reject)

	api.GET("/meetings", middleware.RequirePermission("meetings.view"), d.MeetingHandler.List)
	api.POST("/meetings", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.Create)
	api.GET("/meetings/:id", middleware.RequirePermission("meetings.view"), d.MeetingHandler.Get)
	api.PUT("/meetings/:id", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.Update)
	api.POST("/meetings/:id/actions", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.AddAction)
	api.POST("/meetings/actions/:actionId/assign", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.AssignAction)
	api.POST("/meetings/actions/:actionId/convert-task", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.ConvertActionToTask)
	api.POST("/meetings/:id/links", middleware.RequirePermission("meetings.manage"), d.MeetingHandler.AddLink)

	return r
}
