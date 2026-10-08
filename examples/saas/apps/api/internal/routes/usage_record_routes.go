package routes

import (
	"saas/apps/api/internal/handlers"
	"saas/apps/api/internal/middleware"
)

// UsageRecords routes.
//
// This is the whole surface for usage_records: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.UsageRecordHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Protected.GET("/usage_records", h.List)
		m.Protected.GET("/usage_records/export", h.Export)
		m.Protected.POST("/usage_records/import", h.Import)
		m.Protected.GET("/usage_records/import/template", h.Template)
		m.Protected.GET("/usage_records/:id", h.GetByID)
		m.Protected.GET("/usage_records/:id/pdf", h.PDF)
		m.Protected.POST("/usage_records", h.Create)
		m.Protected.POST("/usage_records/bulk-create", h.BulkCreate)
		m.Protected.PUT("/usage_records/:id", h.Update)
		m.Protected.PATCH("/usage_records/:id", h.Patch)
		m.Protected.POST("/usage_records/bulk-edit", h.BulkEdit)

		m.Staff.DELETE("/usage_records/:id", middleware.RequireRole("ADMIN", "perm:usage_records.delete"), h.Delete)
		m.Staff.POST("/usage_records/bulk", middleware.RequireRole("ADMIN", "perm:usage_records.delete"), h.Bulk)
	})
}
