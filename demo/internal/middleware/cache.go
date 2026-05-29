package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gritdemo/internal/cache"
)

// CacheResponse caches GET request responses in Redis for the given duration.
// Only caches 200 OK responses. Skips caching if no cache service is available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache key from URL + query params
		key := fmt.Sprintf("http:%x", sha256.Sum256([]byte(c.Request.URL.String())))

		// Try to serve from cache
		var cached cachedResponse
		found, err := cacheService.Get(c.Request.Context(), key, &cached)
		if err == nil && found {
			c.Header("X-Cache", "HIT")
			c.Data(cached.Status, cached.ContentType, cached.Body)
			c.Abort()
			return
		}

		// Capture the response
		writer := &responseCapture{ResponseWriter: c.Writer, body: make([]byte, 0)}
		c.Writer = writer
		c.Header("X-Cache", "MISS")

		c.Next()

		// Cache successful responses
		if writer.status == http.StatusOK && len(writer.body) > 0 {
			resp := cachedResponse{
				Status:      writer.status,
				ContentType: writer.Header().Get("Content-Type"),
				Body:        writer.body,
			}
			_ = cacheService.Set(c.Request.Context(), key, resp, ttl)
		}
	}
}

type cachedResponse struct {
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	Body        []byte `json:"body"`
}

type responseCapture struct {
	gin.ResponseWriter
	body   []byte
	status int
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *responseCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
