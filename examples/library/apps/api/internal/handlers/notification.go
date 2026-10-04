package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/models"
	"library/apps/api/internal/respond"
)

type NotificationHandler struct {
	DB *gorm.DB
}

// notificationListLimit is how many rows the bell's dropdown shows. The unread
// count beside it is exact; the list is only the most recent.
const notificationListLimit = 50

// scopeNotifications limits a query to the rows one viewer may see: their own,
// plus the broadcast rows (user_id "") when they are an ADMIN.
//
// One definition of the rule, spelled out three times before this and left out
// of MarkRead altogether, where any signed-in account could clear any
// notification by id.
func scopeNotifications(q *gorm.DB, c *gin.Context) *gorm.DB {
	userID, _ := c.Get("user_id")
	if role, _ := c.Get("user_role"); role == models.RoleAdmin {
		return q.Where("user_id = '' OR user_id = ?", userID)
	}
	return q.Where("user_id = ?", userID)
}

// List returns unread + recent notifications for the bell dropdown.
// Visible to any authenticated user; admins see system-wide ones too.
func (h *NotificationHandler) List(c *gin.Context) {
	var items []models.Notification
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()), c).
		Order("created_at DESC").Limit(notificationListLimit)
	if err := q.Find(&items).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}

	// Quick unread count for the bell badge
	var unread int64
	cq := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("read_at IS NULL")
	if err := cq.Count(&unread).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"unread": unread,
	})
}

// MarkRead marks one notification as read.
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	now := time.Now()
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("id = ?", c.Param("id"))
	res := q.Update("read_at", now)
	if res.Error != nil {
		respond.ServerError(c, "DB_ERROR", res.Error, "Internal server error")
		return
	}
	if res.RowsAffected == 0 {
		respond.Fail(c, respond.CodeNotFound, "notification not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

// MarkAllRead clears the bell for the current viewer.
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	now := time.Now()
	q := scopeNotifications(h.DB.WithContext(c.Request.Context()).Model(&models.Notification{}), c).
		Where("read_at IS NULL")
	if err := q.Update("read_at", now).Error; err != nil {
		respond.ServerError(c, "DB_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all marked read"})
}
