package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"workpulse/internal/service"
)

type permissionKey struct{}

type PermissionsLoader struct {
	Auth *service.AuthService
}

func NewPermissionsLoader(auth *service.AuthService) *PermissionsLoader {
	return &PermissionsLoader{Auth: auth}
}

func (m *PermissionsLoader) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIDStr := c.GetString("org_id")
		userIDStr := c.GetString("user_id")
		if orgIDStr == "" || userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}
		var orgUnitID *uuid.UUID
		if header := c.GetHeader("X-Org-Unit-ID"); header != "" {
			if parsed, err := uuid.Parse(header); err == nil {
				orgUnitID = &parsed
			}
		}
		orgID := uuid.MustParse(orgIDStr)
		userID := uuid.MustParse(userIDStr)
		ctx, err := m.Auth.GetContext(c.Request.Context(), orgID, userID, orgUnitID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("permissions", ctx.Permissions)
		if ctx.ActiveOrgUnit != nil {
			c.Set("org_unit_id", ctx.ActiveOrgUnit.String())
		}
		c.Set("auth_ctx", ctx)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsVal, ok := c.Get("permissions")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission context missing"})
			return
		}
		perms, ok := permsVal.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid permissions"})
			return
		}
		for _, p := range perms {
			if p == permission {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	}
}
