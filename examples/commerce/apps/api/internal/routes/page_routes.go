package routes

import (
	"commerce/apps/api/internal/handlers"
	"commerce/apps/api/internal/middleware"
)

// Pages routes.
//
// This is the whole surface for pages: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.PageHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		// Read-only, outside auth, API key required.
		m.Public.GET("/pages", h.ListPublic)
		m.Public.GET("/pages/:key", h.GetPublic)

		m.Staff.GET("/pages", middleware.RequireRole("ADMIN", "perm:pages.view"), h.List)
		m.Staff.GET("/pages/export", middleware.RequireRole("ADMIN", "perm:pages.view"), h.Export)
		m.Staff.POST("/pages/import", middleware.RequireRole("ADMIN", "perm:pages.create"), h.Import)
		m.Staff.GET("/pages/import/template", middleware.RequireRole("ADMIN", "perm:pages.view"), h.Template)
		m.Staff.GET("/pages/:id", middleware.RequireRole("ADMIN", "perm:pages.view"), h.GetByID)
		m.Staff.GET("/pages/:id/pdf", middleware.RequireRole("ADMIN", "perm:pages.view"), h.PDF)
		m.Staff.POST("/pages", middleware.RequireRole("ADMIN", "perm:pages.create"), h.Create)
		m.Staff.PUT("/pages/:id", middleware.RequireRole("ADMIN", "perm:pages.edit"), h.Update)
		m.Staff.PATCH("/pages/:id", middleware.RequireRole("ADMIN", "perm:pages.edit"), h.Patch)

		m.Staff.DELETE("/pages/:id", middleware.RequireRole("ADMIN", "perm:pages.delete"), h.Delete)
		m.Staff.POST("/pages/bulk", middleware.RequireRole("ADMIN", "perm:pages.delete"), h.Bulk)
	})
}
