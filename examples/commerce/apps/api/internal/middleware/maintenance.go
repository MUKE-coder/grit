package middleware

import (
	"os"

	"github.com/gin-gonic/gin"

	"commerce/apps/api/internal/respond"
)

// Maintenance returns a middleware that checks for a .maintenance file.
// When the file exists, all requests receive a 503 Service Unavailable response.
// Toggle with: grit down (enable) / grit up (disable)
func Maintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := os.Stat(".maintenance"); err == nil {
			respond.Fail(c, respond.CodeMaintenance, "Application is in maintenance mode. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}
