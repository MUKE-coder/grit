package scaffold

// roleHandlerGo emits internal/handlers/role.go — the roles + permissions API.
//
// Endpoints:
//
//	GET    /api/admin/permissions        the catalog tree, for the UI
//	GET    /api/admin/roles              roles + user counts
//	POST   /api/admin/roles              create
//	GET    /api/admin/roles/:id          one role
//	PUT    /api/admin/roles/:id          update name/description/grants
//	DELETE /api/admin/roles/:id          delete (blocked for system roles)
//	PUT    /api/admin/users/:id/roles    assign roles to a user
//
// Guard rails learned from the reference implementation:
//
//   - System roles are protected SERVER-SIDE, not just greyed out in the UI.
//     Shoppleet only disables the name input in React, so a direct PUT can
//     rename ADMIN and quietly break every route that still names it.
//   - Every write calls authz.Invalidate(), so a revoked permission takes effect
//     on the next request rather than lingering in a cache.
//   - Grants are validated against the catalog. A typo like "product.create"
//     (singular) would otherwise be stored silently and simply never match.
func roleHandlerGo() string { return tmpl("api/handlers/role.go") }
