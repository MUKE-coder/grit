package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

// v3.31.47 -- ChartHandler exposes the Preset Chart builder
// endpoint. One endpoint per resource; the preset is a query param
// and the service-side dispatcher decides which aggregation to run.
//
// Mounted at GET /api/admin/dashboard/chart/:resource.
type ChartHandler struct {
	DB *gorm.DB
}

func (h *ChartHandler) Get(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		respond.Fail(c, respond.CodeValidationError, "resource is required")
		return
	}

	preset := c.Query("preset")
	if preset == "" {
		respond.Fail(c, respond.CodeValidationError, "preset is required")
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	params := services.ChartParams{
		Preset: preset,
		Field:  c.Query("field"),
		Limit:  limit,
		Grain:  c.Query("grain"),
	}

	result, err := services.ComputeChart(h.DB.WithContext(c.Request.Context()), resource, params)
	if err != nil {
		respond.ServerError(c, "CHART_FAILED", err, "Could not compute this chart")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
