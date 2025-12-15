package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"

        "workpulse/internal/service"
)

type AuthHandler struct {
        svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

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
