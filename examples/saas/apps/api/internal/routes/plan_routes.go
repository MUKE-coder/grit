package routes

import (
	"saas/apps/api/internal/handlers"
	"saas/apps/api/internal/middleware"
)

// Plans routes.
//
// This is the whole surface for plans: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.PlanHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Staff.GET("/plans", middleware.RequireRole("ADMIN", "perm:plans.view"), h.List)
		m.Staff.GET("/plans/export", middleware.RequireRole("ADMIN", "perm:plans.view"), h.Export)
		m.Staff.POST("/plans/import", middleware.RequireRole("ADMIN", "perm:plans.create"), h.Import)
		m.Staff.GET("/plans/import/template", middleware.RequireRole("ADMIN", "perm:plans.view"), h.Template)
		m.Staff.GET("/plans/:id", middleware.RequireRole("ADMIN", "perm:plans.view"), h.GetByID)
		m.Staff.GET("/plans/:id/pdf", middleware.RequireRole("ADMIN", "perm:plans.view"), h.PDF)
		m.Staff.POST("/plans", middleware.RequireRole("ADMIN", "perm:plans.create"), h.Create)
		m.Staff.POST("/plans/bulk-create", middleware.RequireRole("ADMIN", "perm:plans.create"), h.BulkCreate)
		m.Staff.PUT("/plans/:id", middleware.RequireRole("ADMIN", "perm:plans.edit"), h.Update)
		m.Staff.PATCH("/plans/:id", middleware.RequireRole("ADMIN", "perm:plans.edit"), h.Patch)
		m.Staff.POST("/plans/bulk-edit", middleware.RequireRole("ADMIN", "perm:plans.edit"), h.BulkEdit)

		m.Staff.DELETE("/plans/:id", middleware.RequireRole("ADMIN", "perm:plans.delete"), h.Delete)
		m.Staff.POST("/plans/bulk", middleware.RequireRole("ADMIN", "perm:plans.delete"), h.Bulk)
	})
}
