package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gritdemo/internal/services"
)

type ObservabilityHandler struct {
	Bridge *services.SecObsBridge
}

// Summary aggregates Pulse's overview, SLOs, USE grid, N+1 ranked,
// errors, runtime, health checks, and alerts into one envelope.
func (h *ObservabilityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "PULSE_OFF", "message": "Pulse is not enabled"},
		})
		return
	}

	ctx := c.Request.Context()
	out := gin.H{}
	errs := gin.H{}

	fetches := []struct{ key, path string }{
		{"overview", "/pulse/api/overview?range=1h"},
		{"slos", "/pulse/api/slos"},
		{"use", "/pulse/api/use"},
		{"n1_ranked", "/pulse/api/database/n1/ranked?range=1h&limit=10"},
		{"errors", "/pulse/api/errors?limit=10&resolved=false"},
		{"runtime", "/pulse/api/runtime/current"},
		{"health_checks", "/pulse/api/health/checks"},
		{"alerts", "/pulse/api/alerts?state=firing&limit=20"},
	}
	for _, f := range fetches {
		var body interface{}
		if err := h.Bridge.PulseGet(ctx, f.path, &body); err != nil {
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
