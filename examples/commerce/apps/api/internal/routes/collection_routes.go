package routes

import (
	"commerce/apps/api/internal/handlers"
	"commerce/apps/api/internal/middleware"
)

// Collections routes.
//
// This is the whole surface for collections: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.CollectionHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
			Storage: m.Svc.Storage,
		}

		// Read-only, outside auth, API key required.
		m.Public.GET("/collections", h.ListPublic)
		m.Public.GET("/collections/:key", h.GetPublic)

		m.Staff.GET("/collections", middleware.RequireRole("ADMIN", "perm:collections.view"), h.List)
		m.Staff.GET("/collections/export", middleware.RequireRole("ADMIN", "perm:collections.view"), h.Export)
		m.Staff.POST("/collections/import", middleware.RequireRole("ADMIN", "perm:collections.create"), h.Import)
		m.Staff.GET("/collections/import/template", middleware.RequireRole("ADMIN", "perm:collections.view"), h.Template)
		m.Staff.GET("/collections/:id", middleware.RequireRole("ADMIN", "perm:collections.view"), h.GetByID)
		m.Staff.GET("/collections/:id/pdf", middleware.RequireRole("ADMIN", "perm:collections.view"), h.PDF)
		m.Staff.POST("/collections", middleware.RequireRole("ADMIN", "perm:collections.create"), h.Create)
		m.Staff.PUT("/collections/:id", middleware.RequireRole("ADMIN", "perm:collections.edit"), h.Update)
		m.Staff.PATCH("/collections/:id", middleware.RequireRole("ADMIN", "perm:collections.edit"), h.Patch)

		m.Staff.DELETE("/collections/:id", middleware.RequireRole("ADMIN", "perm:collections.delete"), h.Delete)
		m.Staff.POST("/collections/bulk", middleware.RequireRole("ADMIN", "perm:collections.delete"), h.Bulk)
	})
}
