package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// userID + isAdmin from context; falls back to 0 / false if not present.
func ctxUser(c *gin.Context) (uint, bool) {
	uid, _ := c.Get("user_id")
	id, _ := uid.(uint)
	role, _ := c.Get("user_role")
	return id, role == "ADMIN" || role == "OWNER" || role == "MANAGER"
}

// List returns unread + recent notifications for the bell dropdown.
func (h *NotificationHandler) List(c *gin.Context) {
	userID, isAdmin := ctxUser(c)

	q := h.DB.Order("created_at DESC").Limit(50)
	if isAdmin {
		q = q.Where("user_id = 0 OR user_id = ?", userID)
	} else {
		q = q.Where("user_id = ?", userID)
	}

	var items []models.Notification
	q.Find(&items)

	var unread int64
	cq := h.DB.Model(&models.Notification{}).Where("read_at IS NULL")
	if isAdmin {
		cq = cq.Where("user_id = 0 OR user_id = ?", userID)
	} else {
		cq = cq.Where("user_id = ?", userID)
	}
	cq.Count(&unread)

	c.JSON(http.StatusOK, gin.H{"data": items, "unread": unread})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	now := time.Now()
	h.DB.Model(&models.Notification{}).Where("id = ?", id).Update("read_at", now)
	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, isAdmin := ctxUser(c)
	now := time.Now()
	q := h.DB.Model(&models.Notification{}).Where("read_at IS NULL")
	if isAdmin {
		q = q.Where("user_id = 0 OR user_id = ?", userID)
	} else {
		q = q.Where("user_id = ?", userID)
	}
	q.Update("read_at", now)
	c.JSON(http.StatusOK, gin.H{"message": "all marked read"})
}
