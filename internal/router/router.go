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
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", handler.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext())

	api.POST("/work-items", d.WorkItemHandler.Create)
	api.POST("/work-items/:id/okr-links", d.WorkItemHandler.AddOKRLink)

	api.GET("/projects", d.ProjectHandler.List)
	api.POST("/projects", d.ProjectHandler.Create)
	api.PUT("/projects/:id", d.ProjectHandler.Update)
	api.GET("/projects/:id/milestones", d.ProjectHandler.ListMilestones)
	api.POST("/projects/:id/milestones", d.ProjectHandler.CreateMilestone)
	api.PUT("/milestones/:id", d.ProjectHandler.UpdateMilestone)
	api.GET("/projects/:id/tasks", d.ProjectHandler.ListTasks)
	api.POST("/projects/:id/tasks", d.ProjectHandler.CreateTask)
	api.PUT("/tasks/:id", d.ProjectHandler.UpdateTask)
	api.POST("/projects/:id/links", d.ProjectHandler.AddProjectLink)
	api.DELETE("/projects/:projectId/links/:linkId", d.ProjectHandler.RemoveProjectLink)
	api.POST("/tasks/:id/links", d.ProjectHandler.AddTaskLink)
	api.DELETE("/tasks/:taskId/links/:linkId", d.ProjectHandler.RemoveTaskLink)
	api.GET("/tasks/:id/links", d.ProjectHandler.ListTaskLinks)
	api.GET("/projects/:id/okr-progress", d.ProjectHandler.OKRProgress)

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

	api.GET("/meetings", d.MeetingHandler.List)
	api.POST("/meetings", d.MeetingHandler.Create)
	api.GET("/meetings/:id", d.MeetingHandler.Get)
	api.PUT("/meetings/:id", d.MeetingHandler.Update)
	api.POST("/meetings/:id/actions", d.MeetingHandler.AddAction)
	api.POST("/meetings/actions/:actionId/assign", d.MeetingHandler.AssignAction)
	api.POST("/meetings/actions/:actionId/convert-task", d.MeetingHandler.ConvertActionToTask)
	api.POST("/meetings/:id/links", d.MeetingHandler.AddLink)

	return r
}
