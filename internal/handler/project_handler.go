package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/service"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler { return &ProjectHandler{svc: svc} }

func (h *ProjectHandler) List(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	items, err := h.svc.ListProjects(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		StartDate   *string `json:"start_date"`
		EndDate     *string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	p := &models.Project{
		ID:            uuid.New(),
		OrgID:         orgID,
		OwnerUserID:   uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		Status:        req.Status,
		StartDate:     parseDatePtr(req.StartDate),
		EndDate:       parseDatePtr(req.EndDate),
		SchemaVersion: 1,
		Payload:       datatypes.JSON([]byte(`{}`)),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if p.Status == "" {
		p.Status = "draft"
	}
	if err := h.svc.CreateProject(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		StartDate   *string `json:"start_date"`
		EndDate     *string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.svc.UpdateProject(c.Request.Context(), orgID, id, req.Name, req.Description, req.Status, parseDatePtr(req.StartDate), parseDatePtr(req.EndDate))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": project})
}

func (h *ProjectHandler) ListMilestones(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	items, err := h.svc.ListMilestones(c.Request.Context(), orgID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ProjectHandler) CreateMilestone(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := &models.Milestone{
		ID:          uuid.New(),
		OrgID:       orgID,
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     parseDatePtr(req.DueDate),
		Progress:    0,
	}
	if m.Status == "" {
		m.Status = "active"
	}
	if err := h.svc.CreateMilestone(c.Request.Context(), m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ProjectHandler) UpdateMilestone(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Status      string  `json:"status"`
		DueDate     *string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.UpdateMilestone(c.Request.Context(), orgID, id, req.Title, req.Description, req.Status, parseDatePtr(req.DueDate))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if m == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "milestone not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ProjectHandler) ListTasks(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	tasks, err := h.svc.ListTasks(c.Request.Context(), orgID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (h *ProjectHandler) CreateTask(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	var req struct {
		Title        string   `json:"title" binding:"required"`
		MilestoneID  *string  `json:"milestone_id"`
		ParentTaskID *string  `json:"parent_task_id"`
		Status       string   `json:"status"`
		Tags         []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	task := &models.Task{
		ID:           uuid.New(),
		OrgID:        orgID,
		ProjectID:    projectID,
		MilestoneID:  parseUUID(req.MilestoneID),
		ParentTaskID: parseUUID(req.ParentTaskID),
		Title:        req.Title,
		Status:       req.Status,
		Tags:         mustJSON(req.Tags),
		Order:        0,
		CreatedAt:    now,
		UpdatedAt:    now,
		Metadata:     datatypes.JSON([]byte(`{}`)),
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	if err := h.svc.CreateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (h *ProjectHandler) UpdateTask(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Title       string   `json:"title"`
		Status      string   `json:"status"`
		MilestoneID *string  `json:"milestone_id"`
		ParentTask  *string  `json:"parent_task_id"`
		Tags        []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.svc.UpdateTask(c.Request.Context(), orgID, id, req.Title, req.Status, mustJSON(req.Tags), parseUUID(req.MilestoneID), parseUUID(req.ParentTask))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (h *ProjectHandler) AddProjectLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
		EntityID   string `json:"entity_id" binding:"required"`
		Relation   string `json:"relation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link := &models.ProjectLink{
		ID:         uuid.New(),
		OrgID:      orgID,
		ProjectID:  projectID,
		EntityType: req.EntityType,
		EntityID:   uuid.MustParse(req.EntityID),
		Relation:   req.Relation,
		CreatedAt:  time.Now(),
	}
	if link.Relation == "" {
		link.Relation = "related"
	}
	if err := h.svc.AddProjectLink(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": link})
}

func (h *ProjectHandler) RemoveProjectLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	linkID := uuid.MustParse(c.Param("linkId"))
	h.svc.RemoveProjectLink(c.Request.Context(), orgID, linkID)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (h *ProjectHandler) AddTaskLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	taskID := uuid.MustParse(c.Param("id"))
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
		EntityID   string `json:"entity_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link := &models.TaskLink{
		ID:         uuid.New(),
		OrgID:      orgID,
		TaskID:     taskID,
		EntityType: req.EntityType,
		EntityID:   uuid.MustParse(req.EntityID),
		CreatedAt:  time.Now(),
	}
	if err := h.svc.AddTaskLink(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": link})
}

func (h *ProjectHandler) RemoveTaskLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	linkID := uuid.MustParse(c.Param("linkId"))
	h.svc.RemoveTaskLink(c.Request.Context(), orgID, linkID)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (h *ProjectHandler) ListTaskLinks(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	taskID := uuid.MustParse(c.Param("id"))
	links, err := h.svc.ListTaskLinks(c.Request.Context(), orgID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": links})
}

func (h *ProjectHandler) OKRProgress(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	projectID := uuid.MustParse(c.Param("id"))
	progress, err := h.svc.ComputeOKRProgress(c.Request.Context(), orgID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": progress})
}

func parseDatePtr(input *string) *time.Time {
	if input == nil || *input == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *input)
	if err != nil {
		return nil
	}
	return &t
}
