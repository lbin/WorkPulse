package middleware

import (
	"github.com/gin-gonic/gin"
)

// OrgContext captures org-level context such as scoped org unit from headers.
func OrgContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if unit := c.GetHeader("X-Org-Unit-ID"); unit != "" {
			c.Set("org_unit_id", unit)
		}
		c.Next()
	}
}
