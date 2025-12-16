package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"workpulse/internal/service"
)

// AuthHandler exposes registration and login endpoints.
type AuthHandler struct{ svc *service.AuthService }

// NewAuthHandler wires a new auth handler.
func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

type registerRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
	OrgName     string `json:"org_name"`
}

type loginRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required"`
	OrgID    *string `json:"org_id"`
}

// Register creates a user and organization when none exist.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Register(c.Request.Context(), service.RegisterInput{
		Email:       strings.ToLower(req.Email),
		Password:    req.Password,
		DisplayName: req.DisplayName,
		OrgName:     req.OrgName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token, "user": res.User})
}

// Login validates email/password credentials and responds with JWT.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res *service.AuthResult
	var err error
	if req.OrgID != nil && *req.OrgID != "" {
		orgID := uuid.MustParse(*req.OrgID)
		res, err = h.svc.LoginWithOrg(c.Request.Context(), orgID, req.Email, req.Password)
	} else {
		res, err = h.svc.LoginDefault(c.Request.Context(), req.Email, req.Password)
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token, "user": res.User})
}
