package scaffold

// The admin's Routes screen: the route table in a browser.
//
// `grit routes` has printed it in a terminal for a long time, and that remains
// the fuller answer because the CLI parses routes.go and can therefore report
// the middleware group and the permission each route wants. What nothing
// provided was a searchable version for somebody who is in the admin already,
// looking for the URL of an endpoint they are about to call.
//
// The screen shows method, path and handler, which is what the router knows.
// There is deliberately no access column: Gin's RouteInfo carries the last
// handler and nothing about the chain, and a guessed one would be wrong for
// /api/v1/auth/me and every passkey endpoint, which are under /auth/ and all
// need a signed-in user. Somebody would read that as a security statement.

// apiRouteExplorerHandlerGo emits internal/handlers/routes_explorer.go.
func apiRouteExplorerHandlerGo() string { return tmpl("api/handlers/routes_explorer.go") }

// routeExplorerHandlerBlock builds the handler. It takes the engine, which no
// other handler does, because the table is the thing being served.
const routeExplorerHandlerBlock = `	// The Routes screen reads the router's own table. See internal/handlers.
	routeExplorerHandler := &handlers.RouteExplorerHandler{Engine: r}
`

// routeExplorerRoutesBlock mounts it. Behind the same guard as the other
// system screens: a route list is not a secret and is also not something to
// hand to anybody with an account.
const routeExplorerRoutesBlock = `		staff.GET("/admin/routes", middleware.RequireRole("ADMIN", "perm:settings.view"), routeExplorerHandler.List)
`
