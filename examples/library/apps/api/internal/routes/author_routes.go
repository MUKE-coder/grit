package routes

import (
	"library/apps/api/internal/handlers"
	"library/apps/api/internal/middleware"
)

// Authors routes.
//
// This is the whole surface for authors: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.AuthorHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Staff.GET("/authors", middleware.RequireRole("ADMIN", "perm:authors.view"), h.List)
		m.Staff.GET("/authors/export", middleware.RequireRole("ADMIN", "perm:authors.view"), h.Export)
		m.Staff.POST("/authors/import", middleware.RequireRole("ADMIN", "perm:authors.create"), h.Import)
		m.Staff.GET("/authors/import/template", middleware.RequireRole("ADMIN", "perm:authors.view"), h.Template)
		m.Staff.GET("/authors/:id", middleware.RequireRole("ADMIN", "perm:authors.view"), h.GetByID)
		m.Staff.GET("/authors/:id/pdf", middleware.RequireRole("ADMIN", "perm:authors.view"), h.PDF)
		m.Staff.POST("/authors", middleware.RequireRole("ADMIN", "perm:authors.create"), h.Create)
		m.Staff.PUT("/authors/:id", middleware.RequireRole("ADMIN", "perm:authors.edit"), h.Update)
		m.Staff.PATCH("/authors/:id", middleware.RequireRole("ADMIN", "perm:authors.edit"), h.Patch)

		m.Staff.DELETE("/authors/:id", middleware.RequireRole("ADMIN", "perm:authors.delete"), h.Delete)
		m.Staff.POST("/authors/bulk", middleware.RequireRole("ADMIN", "perm:authors.delete"), h.Bulk)
	})
}
