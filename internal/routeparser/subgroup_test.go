package routeparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// gatewayInstanceRoutes is `grit generate resource GatewayInstance --roles ADMIN`
// output, copied from the generator rather than written by hand: --roles is
// what makes internal/generate/resource_routes.go emit a sub-group with the
// role guard on the group instead of per-route middleware.
const gatewayInstanceRoutes = `package routes

func init() {
	RegisterRoutes(func(m *Mount) {
		h := handlers.NewGatewayInstanceHandler(m.DB)

		// Restricted to ADMIN.
		g := m.Protected.Group("/gateway_instances")
		g.Use(middleware.RequireRole("ADMIN"))
		g.GET("", h.List)
		g.GET("/export", h.Export)
		g.GET("/:id", h.GetByID)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	})
}
`

// A versioned public API mounted on the engine root rather than under
// /api/v1, guarded by a hand-written gate. The local variable is called v1 on
// purpose: routes.go calls its /api/v1 group v1 too, and that collision is
// half of grit#308.
const projectKeyAPIRoutes = `package routes

func init() {
	RegisterRoutes(func(m *Mount) {
		v1 := m.Engine.Group("/v1")
		v1.Use(middleware.RequireProjectKey(m.DB))
		v1.GET("/whoami", handlers.WhoAmI)
	})
}
`

func writeProjectWithSubgroups(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "internal", "routes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"routes.go":                  realisticRoutes,
		"widget_routes.go":           widgetRoutes,
		"gateway_instance_routes.go": gatewayInstanceRoutes,
		"project_key_api_routes.go":  projectKeyAPIRoutes,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(dir, "routes.go")
}

// Routes on a sub-group declared from a Mount field belong to that sub-group.
//
// grit#308: `g := m.Protected.Group("/gateway_instances")` matched the group
// pattern not at all, because it allowed one optional receiver segment and
// this has two. The declaration was skipped, `g` was created by its own
// `.Use` line with no parent and no prefix, and the eleven routes registered
// on it resolved to "/" and "/:id". So an auditor reading `grit routes` for
// what is exposed saw none of the ADMIN-only paths, and because the registry
// keeps the first entry per method-and-path, two such resources collided.
func TestSubgroupOnAMountFieldKeepsItsPrefix(t *testing.T) {
	routes, err := ParseProject(writeProjectWithSubgroups(t))
	if err != nil {
		t.Fatalf("ParseProject: %v", err)
	}

	for _, want := range []struct{ method, path string }{
		{"GET", "/api/v1/gateway_instances"},
		{"GET", "/api/v1/gateway_instances/export"},
		{"GET", "/api/v1/gateway_instances/:id"},
		{"POST", "/api/v1/gateway_instances"},
		{"PUT", "/api/v1/gateway_instances/:id"},
		{"DELETE", "/api/v1/gateway_instances/:id"},
	} {
		r := findRoute(routes, want.method, want.path)
		if r == nil {
			t.Errorf("%s %s is missing from the table; got %v", want.method, want.path, paths(routes))
			continue
		}
		// The guard was already read correctly, from the group's .Use line.
		// It is asserted anyway: the fix moves where the group comes from, and
		// losing the guard while gaining the path would be a worse table.
		if got := r.Access.Summary(); got != "ADMIN" {
			t.Errorf("%s %s access = %q, want ADMIN", want.method, want.path, got)
		}
		if r.Group != "admin" {
			t.Errorf("%s %s group = %q, want admin", want.method, want.path, r.Group)
		}
	}

	// Nothing landed on the engine root, which is where the lost prefix used
	// to put them.
	for _, r := range routes {
		if r.Path == "/" || r.Path == "/:id" {
			t.Errorf("a route resolved to %q, so a group prefix was lost: %s %s", r.Path, r.Method, r.Handler)
		}
	}
}

// A group mounted on m.Engine is under what it was assigned, and a local name
// that collides with one of routes.go's does not inherit it.
func TestEngineMountedGroupKeepsItsOwnPath(t *testing.T) {
	routes, err := ParseProject(writeProjectWithSubgroups(t))
	if err != nil {
		t.Fatalf("ParseProject: %v", err)
	}

	if r := findRoute(routes, "GET", "/api/v1/whoami"); r != nil {
		t.Error("the engine-mounted route was read as /api/v1/whoami, so the local v1 inherited routes.go's group")
	}
	r := findRoute(routes, "GET", "/v1/whoami")
	if r == nil {
		t.Fatalf("GET /v1/whoami is missing; got %v", paths(routes))
	}

	// It is behind a gate this parser does not know. "public" is the one
	// answer that must never be inferred, and it was the answer given.
	if got := r.Access.Summary(); got == "public" {
		t.Error("a route behind RequireProjectKey was summarised as public")
	}
	if got := r.Access.Summary(); !strings.Contains(got, "unknown") || !strings.Contains(got, "RequireProjectKey") {
		t.Errorf("access = %q, want an unknown summary naming the middleware", got)
	}
	if r.Group == "public" {
		t.Error("a route behind an unrecognised gate is in the public bucket")
	}
}

// The seed carries pointers, so a resource file must not be able to change a
// group that routes.go declared for every file parsed after it.
func TestOneFilesGroupsDoNotLeakIntoAnother(t *testing.T) {
	routes, err := ParseProject(writeProjectWithSubgroups(t))
	if err != nil {
		t.Fatalf("ParseProject: %v", err)
	}

	// widget_routes.go is parsed with the same seed as the file that declares
	// its own v1 and its own g. Its routes must be unchanged by that.
	w := findRoute(routes, "GET", "/api/v1/widgets")
	if w == nil {
		t.Fatalf("the widget route disappeared; got %v", paths(routes))
	}
	if got := w.Access.Summary(); got != "ADMIN or widgets.view" {
		t.Errorf("widget access = %q, want ADMIN or widgets.view", got)
	}

	// And routes.go's own /api/v1 group still means /api/v1.
	if findRoute(routes, "GET", "/api/v1/auth/me") == nil {
		t.Errorf("routes.go's own protected route moved; got %v", paths(routes))
	}
}

// A middleware that does not decide access still reads as nothing, or every
// route behind a logger would report as unknown and the column would be noise.
func TestOrdinaryMiddlewareIsStillNotAGuard(t *testing.T) {
	routes := ParseSource(`package routes

func Setup() {
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Gzip())
	open := r.Group("/open")
	open.Use(middleware.CacheResponse(time.Minute))
	open.GET("/ping", pingHandler)
}
`)
	r := findRoute(routes, "GET", "/open/ping")
	if r == nil {
		t.Fatalf("route missing; got %v", paths(routes))
	}
	if got := r.Access.Summary(); got != "public" {
		t.Errorf("access = %q, want public: a logger, a gzip and a cache are not gates", got)
	}
}
