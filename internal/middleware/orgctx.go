package middleware

import "github.com/gin-gonic/gin"

// placeholder for team/role/org scoping
func OrgContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
