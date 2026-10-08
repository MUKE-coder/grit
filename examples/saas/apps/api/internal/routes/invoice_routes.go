package routes

import (
	"saas/apps/api/internal/handlers"
	"saas/apps/api/internal/middleware"
)

// Invoices routes.
//
// This is the whole surface for invoices: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.InvoiceHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Protected.GET("/invoices", h.List)
		m.Protected.GET("/invoices/export", h.Export)
		m.Protected.POST("/invoices/import", h.Import)
		m.Protected.GET("/invoices/import/template", h.Template)
		m.Protected.GET("/invoices/:id", h.GetByID)
		m.Protected.GET("/invoices/:id/pdf", h.PDF)
		m.Protected.POST("/invoices", h.Create)
		m.Protected.POST("/invoices/bulk-create", h.BulkCreate)
		m.Protected.PUT("/invoices/:id", h.Update)
		m.Protected.PATCH("/invoices/:id", h.Patch)
		m.Protected.POST("/invoices/bulk-edit", h.BulkEdit)

		m.Staff.DELETE("/invoices/:id", middleware.RequireRole("ADMIN", "perm:invoices.delete"), h.Delete)
		m.Staff.POST("/invoices/bulk", middleware.RequireRole("ADMIN", "perm:invoices.delete"), h.Bulk)
	})
}
