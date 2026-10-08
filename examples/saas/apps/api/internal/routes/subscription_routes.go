package routes

import (
	"saas/apps/api/internal/handlers"
	"saas/apps/api/internal/middleware"
)

// Subscriptions routes.
//
// This is the whole surface for subscriptions: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.SubscriptionHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Protected.GET("/subscriptions", h.List)
		m.Protected.GET("/subscriptions/export", h.Export)
		m.Protected.POST("/subscriptions/import", h.Import)
		m.Protected.GET("/subscriptions/import/template", h.Template)
		m.Protected.GET("/subscriptions/:id", h.GetByID)
		m.Protected.GET("/subscriptions/:id/pdf", h.PDF)
		m.Protected.POST("/subscriptions", h.Create)
		m.Protected.POST("/subscriptions/bulk-create", h.BulkCreate)
		m.Protected.PUT("/subscriptions/:id", h.Update)
		m.Protected.PATCH("/subscriptions/:id", h.Patch)
		m.Protected.POST("/subscriptions/bulk-edit", h.BulkEdit)

		m.Staff.DELETE("/subscriptions/:id", middleware.RequireRole("ADMIN", "perm:subscriptions.delete"), h.Delete)
		m.Staff.POST("/subscriptions/bulk", middleware.RequireRole("ADMIN", "perm:subscriptions.delete"), h.Bulk)
	})
}
