package middleware

import (
	"github.com/gin-gonic/gin"

	"library/apps/api/internal/services"
)

// RequestMeta records the client IP, the user agent and the request id on the
// request's context, so a service can read them without taking a *gin.Context.
//
// It runs after RequestID, which is where the id comes from, and before Auth,
// which is where the actor comes from. The actor is therefore not in here: a
// service that wants it reads it from the context it was handed, which Auth
// refreshes once it knows who is calling.
func RequestMeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(services.ContextOf(c))
		c.Next()
	}
}
