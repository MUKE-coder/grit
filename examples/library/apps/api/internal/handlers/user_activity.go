package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/paginate"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

type UserActivityHandler struct {
	DB *gorm.DB
}

// activity is the service the reads below go through: which columns a client may
// sort or filter by is a decision about data, because each name ends up in SQL.
// See internal/services/activity_reads.go.
func (h *UserActivityHandler) activity() *services.ActivityService {
	return &services.ActivityService{DB: h.DB}
}

// List returns paginated activity events. Filters: user_id, action,
// severity, resource_type, q (substring against summary). Default sort
// is newest first.
//
//	GET /api/user-activity?severity=critical&page=1&page_size=25
//	GET /api/user-activity?q=login&user_id=...
func (h *UserActivityHandler) List(c *gin.Context) {
	params := paginate.Bind(c).
		With("user_id", c.Query("user_id")).
		With("action", c.Query("action")).
		With("severity", c.Query("severity")).
		With("resource_type", c.Query("resource_type"))

	res, err := h.activity().Events(c.Request.Context(), params, c.Query("q"))
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Stats returns severity bucket counts for the last 24h. Powers the
// header chips on the activity dashboard.
//
//	GET /api/user-activity/stats
//	→ { "data": { "info": 142, "warn": 8, "critical": 1, "total": 151 } }
func (h *UserActivityHandler) Stats(c *gin.Context) {
	// The window is this endpoint's promise, so it is chosen here; the count and
	// the reason it binds a cutoff rather than writing an interval are the
	// service's. Its error is reported now: it used to be dropped, so a failure
	// drew the same zeros as a quiet day.
	out, err := h.activity().SeverityCounts(c.Request.Context(), time.Now().Add(-24*time.Hour))
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
