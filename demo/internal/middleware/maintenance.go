package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Maintenance returns a middleware that checks for a .maintenance file.
// When the file exists, all requests receive a 503 Service Unavailable response.
// Toggle with: grit down (enable) / grit up (disable)
func Maintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := os.Stat(".maintenance"); err == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"code":    "MAINTENANCE",
					"message": "Application is in maintenance mode. Please try again later.",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
