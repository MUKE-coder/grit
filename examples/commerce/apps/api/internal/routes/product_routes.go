package routes

import (
	"commerce/apps/api/internal/handlers"
	"commerce/apps/api/internal/middleware"
)

// Products routes.
//
// This is the whole surface for products: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.ProductHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
			Storage: m.Svc.Storage,
		}

		// Read-only, outside auth, API key required.
		m.Public.GET("/products", h.ListPublic)
		m.Public.GET("/products/:key", h.GetPublic)
		m.Public.GET("/products/:key/related", h.RelatedPublic)

		m.Staff.GET("/products", middleware.RequireRole("ADMIN", "perm:products.view"), h.List)
		m.Staff.GET("/products/export", middleware.RequireRole("ADMIN", "perm:products.view"), h.Export)
		m.Staff.POST("/products/import", middleware.RequireRole("ADMIN", "perm:products.create"), h.Import)
		m.Staff.GET("/products/import/template", middleware.RequireRole("ADMIN", "perm:products.view"), h.Template)
		m.Staff.GET("/products/:id", middleware.RequireRole("ADMIN", "perm:products.view"), h.GetByID)
		m.Staff.GET("/products/:id/pdf", middleware.RequireRole("ADMIN", "perm:products.view"), h.PDF)
		m.Staff.POST("/products", middleware.RequireRole("ADMIN", "perm:products.create"), h.Create)
		m.Staff.PUT("/products/:id", middleware.RequireRole("ADMIN", "perm:products.edit"), h.Update)
		m.Staff.PATCH("/products/:id", middleware.RequireRole("ADMIN", "perm:products.edit"), h.Patch)

		m.Staff.DELETE("/products/:id", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Delete)
		m.Staff.POST("/products/bulk", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Bulk)
	})
}
