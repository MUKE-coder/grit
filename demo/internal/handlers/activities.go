package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// ActivityHandler exposes the audit log to admins. Read-only — there is no
// way to mutate or delete activity rows through the API.
type ActivityHandler struct {
	DB *gorm.DB
}

// List returns activities for the active business with optional filters:
// user_id, action, resource, from (date), to (date), q (search description).
// Paginated, default 100/page (activity volumes get high quickly).
//
// Admin only — checked at the handler level so non-admins get a clean 403.
func (h *ActivityHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	q := h.DB.Preload("User").Where("business_id = ?", businessID)

	if uid := c.Query("user_id"); uid != "" {
		q = q.Where("user_id = ?", uid)
	}
	if action := c.Query("action"); action != "" {
		q = q.Where("action = ?", action)
	}
	if resource := c.Query("resource"); resource != "" {
		q = q.Where("resource = ?", resource)
	}
	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			// `to` is inclusive on the day — so push to the next midnight.
			q = q.Where("created_at < ?", t.AddDate(0, 0, 1))
		}
	}
	if needle := c.Query("q"); needle != "" {
		q = q.Where("description ILIKE ?", "%"+needle+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 500 {
		perPage = 100
	}

	var total int64
	q.Model(&models.Activity{}).Count(&total)

	var rows []models.Activity
	q.Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&rows)

	// Flatten user → user_name to keep payload small.
	out := make([]gin.H, 0, len(rows))
	for _, a := range rows {
		out = append(out, gin.H{
			"id":          a.ID,
			"user_id":     a.UserID,
			"user_name":   a.User.Name,
			"user_email":  a.User.Email,
			"action":      a.Action,
			"resource":    a.Resource,
			"resource_id": a.ResourceID,
			"description": a.Description,
			"metadata":    a.Metadata,
			"ip_address":  a.IPAddress,
			"created_at":  a.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     out,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}
