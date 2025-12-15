package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/service"
)

type MeetingHandler struct {
	svc *service.MeetingService
}

func NewMeetingHandler(svc *service.MeetingService) *MeetingHandler {
	return &MeetingHandler{svc: svc}
}

func (h *MeetingHandler) List(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var teamID *uuid.UUID
	if v := c.Query("team_id"); v != "" {
		id := uuid.MustParse(v)
		teamID = &id
	}
	var start, end *time.Time
	if v := c.Query("start"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			start = &t
		}
	}
	if v := c.Query("end"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			end = &t
		}
	}
	meetings, err := h.svc.List(c.Request.Context(), orgID, teamID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": meetings})
}

func (h *MeetingHandler) Get(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	meeting, actions, links, err := h.svc.Get(c.Request.Context(), orgID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if meeting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "meeting not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"meeting": meeting, "actions": actions, "links": links}})
}

func (h *MeetingHandler) Create(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req struct {
		Title             string    `json:"title" binding:"required"`
		Agenda            string    `json:"agenda"`
		ScheduledAt       time.Time `json:"scheduled_at" binding:"required"`
		DurationMinutes   int       `json:"duration_minutes"`
		FacilitatorUserID *string   `json:"facilitator_user_id"`
		AttendeeIDs       []string  `json:"attendee_ids"`
		Status            string    `json:"status"`
		TeamID            *string   `json:"team_id"`
		Notes             string    `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	meeting := &models.Meeting{
		ID:                uuid.New(),
		OrgID:             orgID,
		TeamID:            parseUUID(req.TeamID),
		Title:             req.Title,
		Agenda:            req.Agenda,
		ScheduledAt:       req.ScheduledAt,
		DurationMinutes:   req.DurationMinutes,
		FacilitatorUserID: parseUUID(req.FacilitatorUserID),
		Notes:             req.Notes,
		AttendeeIDs:       mustJSON(req.AttendeeIDs),
		Status:            req.Status,
		SchemaVersion:     1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if meeting.DurationMinutes == 0 {
		meeting.DurationMinutes = 60
	}
	if meeting.Status == "" {
		meeting.Status = "scheduled"
	}
	if err := h.svc.Create(c.Request.Context(), meeting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": meeting})
}

func (h *MeetingHandler) Update(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	id := uuid.MustParse(c.Param("id"))
	var req struct {
		Title             *string    `json:"title"`
		Agenda            *string    `json:"agenda"`
		ScheduledAt       *time.Time `json:"scheduled_at"`
		DurationMinutes   *int       `json:"duration_minutes"`
		FacilitatorUserID *string    `json:"facilitator_user_id"`
		AttendeeIDs       *[]string  `json:"attendee_ids"`
		Status            *string    `json:"status"`
		TeamID            *string    `json:"team_id"`
		Notes             *string    `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meeting, err := h.svc.Update(c.Request.Context(), orgID, id, func(m *models.Meeting) {
		if req.Title != nil {
			m.Title = *req.Title
		}
		if req.Agenda != nil {
			m.Agenda = *req.Agenda
		}
		if req.ScheduledAt != nil {
			m.ScheduledAt = *req.ScheduledAt
		}
		if req.DurationMinutes != nil {
			m.DurationMinutes = *req.DurationMinutes
		}
		if req.FacilitatorUserID != nil {
			m.FacilitatorUserID = parseUUID(req.FacilitatorUserID)
		}
		if req.AttendeeIDs != nil {
			m.AttendeeIDs = mustJSON(*req.AttendeeIDs)
		}
		if req.Status != nil {
			m.Status = *req.Status
		}
		if req.TeamID != nil {
			m.TeamID = parseUUID(req.TeamID)
		}
		if req.Notes != nil {
			m.Notes = *req.Notes
		}
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if meeting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "meeting not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": meeting})
}

func (h *MeetingHandler) AddAction(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	meetingID := uuid.MustParse(c.Param("id"))
	var req struct {
		Title         string  `json:"title" binding:"required"`
		OwnerUserID   *string `json:"owner_user_id"`
		DueDate       *string `json:"due_date"`
		Status        string  `json:"status"`
		RelatedTaskID *string `json:"related_task_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var due *time.Time
	if req.DueDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DueDate); err == nil {
			due = &t
		}
	}
	now := time.Now()
	action := &models.MeetingAction{
		ID:            uuid.New(),
		OrgID:         orgID,
		MeetingID:     meetingID,
		Title:         req.Title,
		OwnerUserID:   parseUUID(req.OwnerUserID),
		DueDate:       due,
		Status:        req.Status,
		RelatedTaskID: parseUUID(req.RelatedTaskID),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if action.Status == "" {
		action.Status = "todo"
	}
	if err := h.svc.AddAction(c.Request.Context(), action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": action})
}

func (h *MeetingHandler) AssignAction(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	actionID := uuid.MustParse(c.Param("actionId"))
	var req struct {
		OwnerUserID   *string `json:"owner_user_id"`
		Status        *string `json:"status"`
		RelatedTaskID *string `json:"related_task_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	action, err := h.svc.UpdateAction(c.Request.Context(), orgID, actionID, func(a *models.MeetingAction) {
		if req.OwnerUserID != nil {
			a.OwnerUserID = parseUUID(req.OwnerUserID)
		}
		if req.Status != nil {
			a.Status = *req.Status
		}
		if req.RelatedTaskID != nil {
			a.RelatedTaskID = parseUUID(req.RelatedTaskID)
		}
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if action == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "action not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": action})
}

func (h *MeetingHandler) ConvertActionToTask(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	actionID := uuid.MustParse(c.Param("actionId"))
	newTaskID := uuid.New()
	action, err := h.svc.UpdateAction(c.Request.Context(), orgID, actionID, func(a *models.MeetingAction) {
		a.RelatedTaskID = &newTaskID
		a.Status = "doing"
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if action == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "action not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": action, "task_id": newTaskID})
}

func (h *MeetingHandler) AddLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	meetingID := uuid.MustParse(c.Param("id"))
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   string `json:"target_id" binding:"required"`
		Relation   string `json:"relation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link := &models.MeetingLink{
		ID:         uuid.New(),
		OrgID:      orgID,
		MeetingID:  meetingID,
		TargetType: req.TargetType,
		TargetID:   uuid.MustParse(req.TargetID),
		Relation:   req.Relation,
		CreatedAt:  time.Now(),
	}
	if link.Relation == "" {
		link.Relation = "related"
	}
	if err := h.svc.AddLink(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": link})
}

// helpers reuse workitem handler parse logic
type jsonable interface{}

func mustJSON(v jsonable) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

func parseUUID(id *string) *uuid.UUID {
	if id == nil || *id == "" {
		return nil
	}
	u := uuid.MustParse(*id)
	return &u
}
