package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/events"
	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
	"saas/apps/api/internal/sync"
)

/*
 * The bin, and what comes out of it.
 *
 * Every generated model soft-deletes, so a row a person deleted is still on
 * disk and has been invisible since the day the feature shipped. These are the
 * four things worth doing to it: see it, put it back, remove it now, or let the
 * retention window remove it later.
 *
 * Restoring and purging are both ADMIN, and deliberately. Deleting a row is a
 * permission a resource grants; undoing somebody else's delete, or making one
 * permanent, is an operator's job and belongs with the other operator surfaces.
 */

// TrashHandler serves /admin/trash.
type TrashHandler struct {
	Service *services.TrashService
}

// NewTrashHandler builds the handler and the service behind it.
func NewTrashHandler(db *gorm.DB, registry *sync.Registry) *TrashHandler {
	return &TrashHandler{Service: &services.TrashService{DB: db, Registry: registry}}
}

// Buckets answers GET /admin/trash: what is in the bin, per resource.
func (h *TrashHandler) Buckets(c *gin.Context) {
	buckets, err := h.Service.Buckets(c.Request.Context())
	if err != nil {
		respond.WriteError(c, err, "Failed to read the bin")
		return
	}
	var total int64
	for _, b := range buckets {
		total += b.Count
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"buckets":        buckets,
			"total":          total,
			"retention_days": services.TrashRetentionDays,
		},
	})
}

// List answers GET /admin/trash/:table.
func (h *TrashHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	items, total, err := h.Service.List(c.Request.Context(), c.Param("table"), limit, offset)
	if err != nil {
		respond.NotFound(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{"total": total, "retention_days": services.TrashRetentionDays},
	})
}

// Restore answers POST /admin/trash/:table/:id/restore.
func (h *TrashHandler) Restore(c *gin.Context) {
	table, id := c.Param("table"), c.Param("id")
	if err := h.Service.Restore(c.Request.Context(), table, id); err != nil {
		respond.NotFound(c, err.Error())
		return
	}
	// Worth a line in the activity log: somebody undid a delete, and the row is
	// back in lists it had left.
	events.Emitted(c, table, table, "restore", id, "restored a deleted "+table+" row", "restored from the bin", nil, nil)
	respond.OK(c, gin.H{"table": table, "id": id}, "Restored")
}

// Purge answers DELETE /admin/trash/:table/:id.
func (h *TrashHandler) Purge(c *gin.Context) {
	table, id := c.Param("table"), c.Param("id")
	if err := h.Service.Purge(c.Request.Context(), table, id); err != nil {
		respond.NotFound(c, err.Error())
		return
	}
	events.Emitted(c, table, table, "purge", id, "permanently deleted a "+table+" row", "removed from the bin", nil, nil)
	respond.OK(c, gin.H{"table": table, "id": id}, "Deleted for good")
}

// Empty answers DELETE /admin/trash, and DELETE /admin/trash/:table for one
// resource.
//
// It asks for a confirmation in the query string rather than taking the word of
// the method alone: this is the one call on the page that cannot be undone by
// any other call on the page.
func (h *TrashHandler) Empty(c *gin.Context) {
	if c.Query("confirm") != "true" {
		respond.BadRequest(c, "Emptying the bin cannot be undone: send confirm=true")
		return
	}
	table := c.Param("table")
	removed, err := h.Service.Empty(c.Request.Context(), table)
	if err != nil {
		respond.WriteError(c, err, "Failed to empty the bin")
		return
	}
	what := "the bin"
	if table != "" {
		what = table
	}
	events.Emitted(c, "trash", "Trash", "empty", "", "emptied "+what, "removed "+strconv.FormatInt(removed, 10)+" row(s) for good", nil, nil)
	respond.OK(c, gin.H{"removed": removed}, "Emptied")
}
