package routes

import (
	"library/apps/api/internal/handlers"
	"library/apps/api/internal/middleware"
)

// Loans routes.
//
// This is the whole surface for loans: the handler, and every path that
// reaches it. Add a route by adding a line here. Remove the resource by
// deleting this file; nothing else refers to it.
//
// m.Public is outside the auth middleware and behind an API key, m.Protected
// takes a JWT or an API key, and m.Staff also requires the permission its
// route names (an ADMIN holds every one).
func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.LoanHandler{
			DB:      m.DB,
			AppName: m.Cfg.AppName,
		}

		m.Protected.GET("/loans", h.List)
		m.Protected.GET("/loans/export", h.Export)
		m.Protected.POST("/loans/import", h.Import)
		m.Protected.GET("/loans/import/template", h.Template)
		m.Protected.GET("/loans/:id", h.GetByID)
		m.Protected.GET("/loans/:id/pdf", h.PDF)
		m.Protected.POST("/loans", h.Create)
		m.Protected.PUT("/loans/:id", h.Update)
		m.Protected.PATCH("/loans/:id", h.Patch)

		m.Staff.DELETE("/loans/:id", middleware.RequireRole("ADMIN", "perm:loans.delete"), h.Delete)
		m.Staff.POST("/loans/bulk", middleware.RequireRole("ADMIN", "perm:loans.delete"), h.Bulk)
	})
}
