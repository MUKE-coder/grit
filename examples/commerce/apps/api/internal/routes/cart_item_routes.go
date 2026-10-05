package routes

import (
	"commerce/apps/api/internal/handlers"
	"commerce/apps/api/internal/middleware"
)

// CartItems routes.
//
// This is the whole surface for cart_items: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.CartItemHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Staff.GET("/cart_items", middleware.RequireRole("ADMIN", "perm:cart_items.view"), h.List)
		m.Staff.GET("/cart_items/export", middleware.RequireRole("ADMIN", "perm:cart_items.view"), h.Export)
		m.Staff.POST("/cart_items/import", middleware.RequireRole("ADMIN", "perm:cart_items.create"), h.Import)
		m.Staff.GET("/cart_items/import/template", middleware.RequireRole("ADMIN", "perm:cart_items.view"), h.Template)
		m.Staff.GET("/cart_items/:id", middleware.RequireRole("ADMIN", "perm:cart_items.view"), h.GetByID)
		m.Staff.GET("/cart_items/:id/pdf", middleware.RequireRole("ADMIN", "perm:cart_items.view"), h.PDF)
		m.Staff.POST("/cart_items", middleware.RequireRole("ADMIN", "perm:cart_items.create"), h.Create)
		m.Staff.PUT("/cart_items/:id", middleware.RequireRole("ADMIN", "perm:cart_items.edit"), h.Update)
		m.Staff.PATCH("/cart_items/:id", middleware.RequireRole("ADMIN", "perm:cart_items.edit"), h.Patch)

		m.Staff.DELETE("/cart_items/:id", middleware.RequireRole("ADMIN", "perm:cart_items.delete"), h.Delete)
		m.Staff.POST("/cart_items/bulk", middleware.RequireRole("ADMIN", "perm:cart_items.delete"), h.Bulk)
	})
}
