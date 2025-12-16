package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"workpulse/internal/service"
)

// AuthHandler exposes authentication and onboarding endpoints.
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler wires dependencies into AuthHandler.
func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

// LoginRequest captures the payload for credential login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest captures signup data for first-time users.
type RegisterRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name"`
}

// CreateTeamRequest defines the fields to create a new team/org unit.
type CreateTeamRequest struct {
	Name         string  `json:"name" binding:"required"`
	ParentTeamID *string `json:"parent_team_id"`
}

// JoinTeamRequest accepts a team ID to join.
type JoinTeamRequest struct {
	TeamID string `json:"team_id" binding:"required"`
}

// Login exchanges email/password for a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Register creates a new account and seeds a default team.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Me returns the hydrated auth context for the bearer token.
func (h *AuthHandler) Me(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	userID := uuid.MustParse(c.GetString("user_id"))
	var orgUnit *uuid.UUID
	if header := c.GetHeader("X-Org-Unit-ID"); header != "" {
		if parsed, err := uuid.Parse(header); err == nil {
			orgUnit = &parsed
		}
	}
	ctx, err := h.svc.GetContext(c.Request.Context(), orgID, userID, orgUnit)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ctx)
}

// CreateTeam adds a new team and enrolls the current user.
func (h *AuthHandler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := uuid.MustParse(c.GetString("org_id"))
	userID := uuid.MustParse(c.GetString("user_id"))

	var parentID *uuid.UUID
	if req.ParentTeamID != nil {
		id := uuid.MustParse(*req.ParentTeamID)
		parentID = &id
	}

	team, err := h.svc.CreateTeam(c.Request.Context(), orgID, userID, req.Name, parentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// JoinTeam subscribes the user to an existing team.
func (h *AuthHandler) JoinTeam(c *gin.Context) {
	var req JoinTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := uuid.MustParse(c.GetString("org_id"))
	userID := uuid.MustParse(c.GetString("user_id"))
	teamID := uuid.MustParse(req.TeamID)

	if err := h.svc.JoinTeam(c.Request.Context(), orgID, userID, teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, _ := h.svc.GetContext(c.Request.Context(), orgID, userID, nil)
	c.JSON(http.StatusOK, ctx)
}

// ListTeams returns the teams the current user belongs to.
func (h *AuthHandler) ListTeams(c *gin.Context) {
	orgID := uuid.MustParse(c.GetString("org_id"))
	userID := uuid.MustParse(c.GetString("user_id"))

	teams, err := h.svc.ListTeams(c.Request.Context(), orgID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}
