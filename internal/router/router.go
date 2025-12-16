package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"workpulse/internal/handler"
	"workpulse/internal/middleware"
)

type Deps struct {
	ServiceName     string
	JWTSecret       string
	APIVersion      string
	MetricsPath     string
	MetricsHandler  http.Handler
	TelemetryMW     gin.HandlerFunc
	WorkItemHandler *handler.WorkItemHandler
	ProjectHandler  *handler.ProjectHandler
	OKRHandler      *handler.OKRHandler
	ReportHandler   *handler.ReportHandler
	MeetingHandler  *handler.MeetingHandler
	AuthHandler     *handler.AuthHandler
	ConfigHandler   *handler.ConfigHandler
	PermissionMW    *middleware.PermissionsLoader
}

func New(d Deps) *gin.Engine {
	r := gin.Default()
	if d.TelemetryMW != nil {
		r.Use(d.TelemetryMW)
	}

	if d.MetricsHandler != nil {
		path := d.MetricsPath
		if path == "" {
			path = "/metrics"
		}
		r.GET(path, gin.WrapH(d.MetricsHandler))
	}
	r.GET("/healthz", handler.Health)

	versionedBase := fmt.Sprintf("/api/%s", d.APIVersion)

	public := r.Group(versionedBase)
	public.POST("/auth/login", d.AuthHandler.Login)
	public.POST("/auth/register", d.AuthHandler.Register)

	api := r.Group(versionedBase)
	api.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext(), d.PermissionMW.Handler())

	specs := []RouteSpec{
		{Method: http.MethodGet, Path: "/auth/me", Handler: d.AuthHandler.Me, Summary: "Current authenticated user", Tag: "auth"},
		{Method: http.MethodGet, Path: "/auth/teams", Handler: d.AuthHandler.ListTeams, Summary: "List my teams", Tag: "auth"},
		{Method: http.MethodPost, Path: "/auth/teams", Handler: d.AuthHandler.CreateTeam, Summary: "Create team", Tag: "auth"},
		{Method: http.MethodPost, Path: "/auth/teams/join", Handler: d.AuthHandler.JoinTeam, Summary: "Join team", Tag: "auth"},

		{Method: http.MethodPost, Path: "/work-items", Permission: "workitems.manage", Handler: d.WorkItemHandler.Create, Summary: "Create work item", Tag: "work_items"},
		{Method: http.MethodPost, Path: "/work-items/:id/okr-links", Permission: "workitems.manage", Handler: d.WorkItemHandler.AddOKRLink, Summary: "Link work item to OKR", Tag: "work_items"},

		{Method: http.MethodGet, Path: "/projects", Permission: "projects.view", Handler: d.ProjectHandler.List, Summary: "List projects", Tag: "projects"},
		{Method: http.MethodPost, Path: "/projects", Permission: "projects.manage", Handler: d.ProjectHandler.Create, Summary: "Create project", Tag: "projects"},
		{Method: http.MethodPut, Path: "/projects/:id", Permission: "projects.manage", Handler: d.ProjectHandler.Update, Summary: "Update project", Tag: "projects"},
		{Method: http.MethodGet, Path: "/projects/:id/milestones", Permission: "projects.view", Handler: d.ProjectHandler.ListMilestones, Summary: "List milestones", Tag: "projects"},
		{Method: http.MethodPost, Path: "/projects/:id/milestones", Permission: "projects.manage", Handler: d.ProjectHandler.CreateMilestone, Summary: "Create milestone", Tag: "projects"},
		{Method: http.MethodPut, Path: "/milestones/:id", Permission: "projects.manage", Handler: d.ProjectHandler.UpdateMilestone, Summary: "Update milestone", Tag: "projects"},
		{Method: http.MethodGet, Path: "/projects/:id/tasks", Permission: "projects.view", Handler: d.ProjectHandler.ListTasks, Summary: "List project tasks", Tag: "projects"},
		{Method: http.MethodPost, Path: "/projects/:id/tasks", Permission: "projects.manage", Handler: d.ProjectHandler.CreateTask, Summary: "Create task", Tag: "projects"},
		{Method: http.MethodPut, Path: "/tasks/:id", Permission: "projects.manage", Handler: d.ProjectHandler.UpdateTask, Summary: "Update task", Tag: "projects"},
		{Method: http.MethodPost, Path: "/projects/:id/links", Permission: "projects.manage", Handler: d.ProjectHandler.AddProjectLink, Summary: "Add project link", Tag: "projects"},
		{Method: http.MethodDelete, Path: "/projects/:projectId/links/:linkId", Permission: "projects.manage", Handler: d.ProjectHandler.RemoveProjectLink, Summary: "Remove project link", Tag: "projects"},
		{Method: http.MethodPost, Path: "/tasks/:id/links", Permission: "projects.manage", Handler: d.ProjectHandler.AddTaskLink, Summary: "Add task link", Tag: "projects"},
		{Method: http.MethodDelete, Path: "/tasks/:taskId/links/:linkId", Permission: "projects.manage", Handler: d.ProjectHandler.RemoveTaskLink, Summary: "Remove task link", Tag: "projects"},
		{Method: http.MethodGet, Path: "/tasks/:id/links", Permission: "projects.view", Handler: d.ProjectHandler.ListTaskLinks, Summary: "List task links", Tag: "projects"},
		{Method: http.MethodGet, Path: "/projects/:id/okr-progress", Permission: "projects.view", Handler: d.ProjectHandler.OKRProgress, Summary: "Project OKR progress", Tag: "projects"},

		{Method: http.MethodGet, Path: "/okrs/cycles", Permission: "okr.view", Handler: d.OKRHandler.ListCycles, Summary: "List cycles", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/cycles", Permission: "okr.manage", Handler: d.OKRHandler.CreateCycle, Summary: "Create cycle", Tag: "okr"},
		{Method: http.MethodPut, Path: "/okrs/cycles/:id", Permission: "okr.manage", Handler: d.OKRHandler.UpdateCycle, Summary: "Update cycle", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/cycles/:id/archive", Permission: "okr.manage", Handler: d.OKRHandler.ArchiveCycle, Summary: "Archive cycle", Tag: "okr"},

		{Method: http.MethodGet, Path: "/okrs/objectives", Permission: "okr.view", Handler: d.OKRHandler.ListObjectives, Summary: "List objectives", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/objectives", Permission: "okr.manage", Handler: d.OKRHandler.CreateObjective, Summary: "Create objective", Tag: "okr"},
		{Method: http.MethodPut, Path: "/okrs/objectives/:id", Permission: "okr.manage", Handler: d.OKRHandler.UpdateObjective, Summary: "Update objective", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/objectives/:id/archive", Permission: "okr.manage", Handler: d.OKRHandler.ArchiveObjective, Summary: "Archive objective", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/key-results", Permission: "okr.manage", Handler: d.OKRHandler.CreateKeyResult, Summary: "Create key result", Tag: "okr"},
		{Method: http.MethodPut, Path: "/okrs/key-results/:id", Permission: "okr.manage", Handler: d.OKRHandler.UpdateKeyResult, Summary: "Update key result", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/key-results/:id/archive", Permission: "okr.manage", Handler: d.OKRHandler.ArchiveKeyResult, Summary: "Archive key result", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/key-results/:id/progress", Permission: "okr.manage", Handler: d.OKRHandler.UpdateProgress, Summary: "Update key result progress", Tag: "okr"},
		{Method: http.MethodPost, Path: "/okrs/links", Permission: "okr.manage", Handler: d.OKRHandler.AddLink, Summary: "Create OKR link", Tag: "okr"},
		{Method: http.MethodDelete, Path: "/okrs/links/:id", Permission: "okr.manage", Handler: d.OKRHandler.RemoveLink, Summary: "Delete OKR link", Tag: "okr"},
		{Method: http.MethodGet, Path: "/okrs/metrics", Permission: "okr.view", Handler: d.OKRHandler.Metrics, Summary: "OKR metrics", Tag: "okr"},

		{Method: http.MethodGet, Path: "/reports", Permission: "reports.view", Handler: d.ReportHandler.List, Summary: "List reports", Tag: "reports"},
		{Method: http.MethodGet, Path: "/reports/export", Permission: "reports.view", Handler: d.ReportHandler.Export, Summary: "Export reports", Tag: "reports"},
		{Method: http.MethodGet, Path: "/reports/:id", Permission: "reports.view", Handler: d.ReportHandler.Get, Summary: "Get report", Tag: "reports"},
		{Method: http.MethodPost, Path: "/reports", Permission: "reports.manage", Handler: d.ReportHandler.Create, Summary: "Create report", Tag: "reports"},
		{Method: http.MethodPut, Path: "/reports/:id", Permission: "reports.manage", Handler: d.ReportHandler.Update, Summary: "Update report", Tag: "reports"},
		{Method: http.MethodPost, Path: "/reports/:id/submit", Permission: "reports.manage", Handler: d.ReportHandler.Submit, Summary: "Submit report", Tag: "reports"},
		{Method: http.MethodPost, Path: "/reports/:id/approve", Permission: "reports.review", Handler: d.ReportHandler.Approve, Summary: "Approve report", Tag: "reports"},
		{Method: http.MethodPost, Path: "/reports/:id/reject", Permission: "reports.review", Handler: d.ReportHandler.Reject, Summary: "Reject report", Tag: "reports"},

		{Method: http.MethodGet, Path: "/meetings", Permission: "meetings.view", Handler: d.MeetingHandler.List, Summary: "List meetings", Tag: "meetings"},
		{Method: http.MethodPost, Path: "/meetings", Permission: "meetings.manage", Handler: d.MeetingHandler.Create, Summary: "Create meeting", Tag: "meetings"},
		{Method: http.MethodGet, Path: "/meetings/:id", Permission: "meetings.view", Handler: d.MeetingHandler.Get, Summary: "Get meeting", Tag: "meetings"},
		{Method: http.MethodPut, Path: "/meetings/:id", Permission: "meetings.manage", Handler: d.MeetingHandler.Update, Summary: "Update meeting", Tag: "meetings"},
		{Method: http.MethodPost, Path: "/meetings/:id/actions", Permission: "meetings.manage", Handler: d.MeetingHandler.AddAction, Summary: "Add meeting action", Tag: "meetings"},
		{Method: http.MethodPost, Path: "/meetings/actions/:actionId/assign", Permission: "meetings.manage", Handler: d.MeetingHandler.AssignAction, Summary: "Assign action", Tag: "meetings"},
		{Method: http.MethodPost, Path: "/meetings/actions/:actionId/convert-task", Permission: "meetings.manage", Handler: d.MeetingHandler.ConvertActionToTask, Summary: "Convert action to task", Tag: "meetings"},
		{Method: http.MethodPost, Path: "/meetings/:id/links", Permission: "meetings.manage", Handler: d.MeetingHandler.AddLink, Summary: "Link meeting", Tag: "meetings"},

		{Method: http.MethodGet, Path: "/config", Handler: d.ConfigHandler.Get, Summary: "Client config and feature flags", Tag: "config"},
	}

	registerRoutes(api, specs)

	docs := r.Group(fmt.Sprintf("%s/docs", versionedBase))
	docs.Use(middleware.Auth(d.JWTSecret), middleware.OrgContext())
	docs.GET("/openapi.json", serveOpenAPI(versionedBase, specs))
	docs.GET("/graphql.sdl", serveGraphQLSDL(specs))

	return r
}
