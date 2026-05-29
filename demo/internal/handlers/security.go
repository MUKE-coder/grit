package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gritdemo/internal/services"
)

type SecurityHandler struct {
	Bridge *services.SecObsBridge
}

// Summary aggregates Sentinel's dashboard summary, security score,
// recent threats, AuthShield state, CSP violations, and performance
// overview into one envelope so the React page does a single round-trip.
func (h *SecurityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "SENTINEL_OFF", "message": "Sentinel is not enabled"},
		})
		return
	}

	ctx := c.Request.Context()
	out := gin.H{}
	errs := gin.H{}

	fetches := []struct{ key, path string }{
		{"summary", "/sentinel/api/dashboard/summary"},
		{"score", "/sentinel/api/security-score"},
		{"threats", "/sentinel/api/threats?limit=10"},
		{"auth_shield", "/sentinel/api/auth-shield/status"},
		{"csp_top", "/sentinel/api/csp-violations/stats"},
		{"performance", "/sentinel/api/performance/overview"},
		{"trends", "/sentinel/api/analytics/trends"},
		{"top_routes", "/sentinel/api/analytics/top-routes"},
	}
	for _, f := range fetches {
		var body interface{}
		if err := h.Bridge.SentinelGet(ctx, f.path, &body); err != nil {
			errs[f.key] = err.Error()
			continue
		}
		out[f.key] = body
	}
	if len(errs) > 0 {
		out["_errors"] = errs
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}
