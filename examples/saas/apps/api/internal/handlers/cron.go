package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"saas/apps/api/internal/cron"
)

// CronHandler handles admin cron task endpoints.
type CronHandler struct{}

// ListTasks returns all registered cron tasks.
func (h *CronHandler) ListTasks(c *gin.Context) {
	tasks := cron.RegisteredTasks
	if tasks == nil {
		tasks = []cron.Task{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}
