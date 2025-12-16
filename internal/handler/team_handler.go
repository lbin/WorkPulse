package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"workpulse/internal/service"
)

// TeamHandler manages team creation and membership endpoints.
type TeamHandler struct{ svc *service.TeamService }

// NewTeamHandler constructs the handler.
func NewTeamHandler(svc *service.TeamService) *TeamHandler { return &TeamHandler{svc: svc} }

type createTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

type joinTeamRequest struct {
	TeamID string `json:"team_id" binding:"required"`
}

// Create builds a new team and enrolls the requester as owner.
func (h *TeamHandler) Create(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, err := h.svc.CreateTeam(c.Request.Context(), orgID, userID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": team})
}

// List provides all teams within the organization.
func (h *TeamHandler) List(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	teams, err := h.svc.ListTeams(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// ListMine returns memberships for the current user.
func (h *TeamHandler) ListMine(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	memberships, err := h.svc.ListUserTeams(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

// Join enrolls the current user into an existing team.
func (h *TeamHandler) Join(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	orgID := uuid.MustParse(c.GetString("org_id"))
	var req joinTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teamID := uuid.MustParse(req.TeamID)
	if err := h.svc.JoinTeam(c.Request.Context(), orgID, teamID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
