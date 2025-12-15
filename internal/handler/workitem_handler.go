package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"workpulse/internal/models"
	"workpulse/internal/service"
)

type WorkItemHandler struct{ svc *service.WorkItemService }

func NewWorkItemHandler(svc *service.WorkItemService) *WorkItemHandler { return &WorkItemHandler{svc: svc} }

type createWorkItemReq struct {
	TeamID        *string `json:"team_id"`
	OwnerUserID   string  `json:"owner_user_id" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	Priority      int     `json:"priority"`
	EffortMinutes *int    `json:"effort_minutes"`
}

func (h *WorkItemHandler) Create(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req createWorkItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ownerID := uuid.MustParse(req.OwnerUserID)

	var teamID *uuid.UUID
	if req.TeamID != nil && *req.TeamID != "" {
		t := uuid.MustParse(*req.TeamID)
		teamID = &t
	}

	wi := &models.WorkItem{
		OrgID:         orgID,
		TeamID:        teamID,
		OwnerUserID:   ownerID,
		Type:          req.Type,
		Title:         req.Title,
		Description:   req.Description,
		Status:        "todo",
		Priority:      req.Priority,
		EffortMinutes: req.EffortMinutes,
		SchemaVersion: 1,
		Tags:          datatypes.JSON([]byte(`[]`)),
		CustomFields:  datatypes.JSON([]byte(`{}`)),
	}

	if wi.Priority == 0 {
		wi.Priority = 3
	}

	if err := h.svc.Create(c.Request.Context(), wi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wi})
}

type addLinkReq struct {
	ObjectiveID         *string `json:"objective_id"`
	KeyResultID         *string `json:"key_result_id"`
	LinkType            string  `json:"link_type" binding:"required"` // planned/actual/both
	ContributionWeight  float64 `json:"contribution_weight"`
	EvidenceRequired    bool    `json:"evidence_required"`
}

func (h *WorkItemHandler) AddOKRLink(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	workItemID := uuid.MustParse(c.Param("id"))

	var req addLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ContributionWeight == 0 {
		req.ContributionWeight = 1.0
	}

	var objID *uuid.UUID
	if req.ObjectiveID != nil && *req.ObjectiveID != "" {
		t := uuid.MustParse(*req.ObjectiveID)
		objID = &t
	}
	var krID *uuid.UUID
	if req.KeyResultID != nil && *req.KeyResultID != "" {
		t := uuid.MustParse(*req.KeyResultID)
		krID = &t
	}

	if objID == nil && krID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "objective_id or key_result_id is required"})
		return
	}

	if err := h.svc.AddOKRLink(c.Request.Context(), orgID, workItemID, objID, krID, req.LinkType, req.ContributionWeight, req.EvidenceRequired); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
