package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

// ResourceStatsHandler exposes per-resource dashboard stats. One
// endpoint per resource, mounted at
// GET /api/admin/dashboard/resource-stats/:resource.
//
// The handler is intentionally thin -- date-range parsing and the
// per-resource dispatch live in services.ComputeResourceStats so the
// generator only has to inject one line per new resource (a switch
// case) without touching the HTTP layer.
type ResourceStatsHandler struct {
	DB *gorm.DB
}

// Get returns the dashboard stat bundle for one resource. Accepts the
// same query params the DateFilter component already sends to the
// resource list endpoint, so the dashboard widgets and the per-page
// stats agree on the wire shape.
//
// Query params:
//
//	?created_since=7d                 — relative window
//	?created_from=2026-01-01&created_to=2026-01-31
//	                                  — explicit range
//	?limit=10                         — latest-rows cap (default 10)
func (h *ResourceStatsHandler) Get(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		respond.Fail(c, respond.CodeValidationError, "resource is required")
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	filter := services.ResourceStatsFilter{
		Since:       c.Query("created_since"),
		From:        c.Query("created_from"),
		To:          c.Query("created_to"),
		LatestLimit: limit,
	}

	stats, err := services.ComputeResourceStats(h.DB.WithContext(c.Request.Context()), resource, filter)
	if err != nil {
		// An unregistered resource and a database error both land here. Both are
		// the server's to fix, so the cause is logged and the widget gets a
		// message it can show, not the database's error text.
		respond.ServerError(c, "STATS_FAILED", err, "Could not compute the stats for this resource")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
