package routes

import (
	"commerce/apps/api/internal/handlers"
	"commerce/apps/api/internal/middleware"
)

// Carts routes.
//
// This is the whole surface for carts: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.CartHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Staff.GET("/carts", middleware.RequireRole("ADMIN", "perm:carts.view"), h.List)
		m.Staff.GET("/carts/export", middleware.RequireRole("ADMIN", "perm:carts.view"), h.Export)
		m.Staff.POST("/carts/import", middleware.RequireRole("ADMIN", "perm:carts.create"), h.Import)
		m.Staff.GET("/carts/import/template", middleware.RequireRole("ADMIN", "perm:carts.view"), h.Template)
		m.Staff.GET("/carts/:id", middleware.RequireRole("ADMIN", "perm:carts.view"), h.GetByID)
		m.Staff.GET("/carts/:id/pdf", middleware.RequireRole("ADMIN", "perm:carts.view"), h.PDF)
		m.Staff.POST("/carts", middleware.RequireRole("ADMIN", "perm:carts.create"), h.Create)
		m.Staff.PUT("/carts/:id", middleware.RequireRole("ADMIN", "perm:carts.edit"), h.Update)
		m.Staff.PATCH("/carts/:id", middleware.RequireRole("ADMIN", "perm:carts.edit"), h.Patch)

		m.Staff.DELETE("/carts/:id", middleware.RequireRole("ADMIN", "perm:carts.delete"), h.Delete)
		m.Staff.POST("/carts/bulk", middleware.RequireRole("ADMIN", "perm:carts.delete"), h.Bulk)
	})
}
