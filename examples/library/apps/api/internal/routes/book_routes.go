package routes

import (
	"library/apps/api/internal/handlers"
	"library/apps/api/internal/middleware"
)

// Books routes.
//
// This is the whole surface for books: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.BookHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
			Storage: m.Svc.Storage,
		}

		m.Staff.GET("/books", middleware.RequireRole("ADMIN", "perm:books.view"), h.List)
		m.Staff.GET("/books/export", middleware.RequireRole("ADMIN", "perm:books.view"), h.Export)
		m.Staff.POST("/books/import", middleware.RequireRole("ADMIN", "perm:books.create"), h.Import)
		m.Staff.GET("/books/import/template", middleware.RequireRole("ADMIN", "perm:books.view"), h.Template)
		m.Staff.GET("/books/:id", middleware.RequireRole("ADMIN", "perm:books.view"), h.GetByID)
		m.Staff.GET("/books/:id/pdf", middleware.RequireRole("ADMIN", "perm:books.view"), h.PDF)
		m.Staff.POST("/books", middleware.RequireRole("ADMIN", "perm:books.create"), h.Create)
		m.Staff.PUT("/books/:id", middleware.RequireRole("ADMIN", "perm:books.edit"), h.Update)
		m.Staff.PATCH("/books/:id", middleware.RequireRole("ADMIN", "perm:books.edit"), h.Patch)

		m.Staff.DELETE("/books/:id", middleware.RequireRole("ADMIN", "perm:books.delete"), h.Delete)
		m.Staff.POST("/books/bulk", middleware.RequireRole("ADMIN", "perm:books.delete"), h.Bulk)
	})
}
