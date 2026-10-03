package routeparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The generated routes.go in miniature: the four groups it really builds, the
// two forms of auth middleware it really writes, and a route with its guard
// passed inline, which is the form that used to parse as nothing at all.
const realisticRoutes = `package routes

const APIVersion = "v1"

func mountAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/passkeys/login/begin", passkeyHandler.BeginLogin)
	}
}

func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.GET("/api/health", healthCheck)
	r.GET("/api/"+APIVersion+"/health", healthCheck)

	v1 := r.Group("/api/" + APIVersion)

	publicAPI := v1.Group("/public")
	publicAPI.Use(middleware.RequireAPIKey(db, svc.Cache))
	publicAPI.GET("/products", productHandler.PublicList)

	protected := v1.Group("")
	protected.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	protected.GET("/auth/me", authHandler.Me)
	protected.GET("/auth/passkeys", passkeyHandler.List)

	staff := v1.Group("")
	staff.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	staff.Use(middleware.RequireStaff())
	staff.PUT("/users/:id/roles", middleware.RequireRole("ADMIN", "perm:users.edit"), roleHandler.AssignUserRoles)
	staff.GET("/dashboard/:resource/stats", middleware.RequirePermissionFor("resource", "view"), statsHandler.Show)

	admin := v1.Group("")
	admin.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	admin.Use(middleware.RequireRole("ADMIN"))
	admin.POST("/admin/flags", featureFlagHandler.Create)

	mountResources(&Mount{
		Engine:    r,
		V1:        v1,
		Public:    publicAPI,
		Protected: protected,
		Admin:     admin,
		Staff:     staff,
	})

	return r
}
`

const widgetRoutes = `package routes

func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.WidgetHandler{DB: m.DB}
		m.Staff.GET("/widgets", middleware.RequireRole("ADMIN", "perm:widgets.view"), h.List)
		m.Staff.DELETE("/widgets/:id", middleware.RequireRole("ADMIN", "perm:widgets.delete"), h.Delete)
		m.Public.GET("/widgets/feed", h.Feed)
	})
}
`

func writeProject(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "internal", "routes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"routes.go":        realisticRoutes,
		"widget_routes.go": widgetRoutes,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(dir, "routes.go")
}

// A route with its guard passed inline used to match nothing, which silently
// dropped 63 of the 171 routes a fresh project registers, and specifically the
// 55 that name a permission.
func TestInlineMiddlewareDoesNotHideTheRoute(t *testing.T) {
	routes := ParseSource(realisticRoutes)

	r := findRoute(routes, "PUT", "/api/v1/users/:id/roles")
	if r == nil {
		t.Fatalf("the route with inline middleware was dropped; got %d routes", len(routes))
	}
	if r.Handler != "roleHandler.AssignUserRoles" {
		t.Errorf("handler = %q, want roleHandler.AssignUserRoles", r.Handler)
	}
	if got := r.Access.Summary(); got != "ADMIN or users.edit" {
		t.Errorf("access = %q, want %q", got, "ADMIN or users.edit")
	}
}

// The one defect worth a test of its own: an authenticated route reported as
// public is read as a security statement.
func TestAuthenticatedRoutesAreNeverCalledPublic(t *testing.T) {
	routes := ParseSource(realisticRoutes)

	for _, path := range []string{"/api/v1/auth/me", "/api/v1/auth/passkeys"} {
		r := findRoute(routes, "GET", path)
		if r == nil {
			t.Errorf("%s not found", path)
			continue
		}
		if !r.Access.Authenticated {
			t.Errorf("%s: Authenticated = false, and it is behind APIKeyOrAuth", path)
		}
		if got := r.Access.Summary(); got == "public" {
			t.Errorf("%s: reported as public", path)
		}
	}

	// And the genuinely public ones still say so, or the check above would pass
	// by calling everything protected.
	for _, path := range []string{"/api/v1/auth/login", "/api/v1/auth/passkeys/login/begin"} {
		r := findRoute(routes, "POST", path)
		if r == nil {
			t.Errorf("%s not found", path)
			continue
		}
		if got := r.Access.Summary(); got != "public" {
			t.Errorf("%s: access = %q, want public", path, got)
		}
	}
}

// "/api/" + APIVersion + "/health" was unquoted whole, producing the path
// /api/"+APIVersion+"/health.
func TestPathBuiltFromAConstantResolves(t *testing.T) {
	routes := ParseSource(realisticRoutes)
	if findRoute(routes, "GET", "/api/v1/health") == nil {
		t.Errorf("the concatenated path did not resolve; got %v", paths(routes))
	}
	for _, r := range routes {
		if strings.Contains(r.Path, `"`) || strings.Contains(r.Path, "+") {
			t.Errorf("path %q was not resolved", r.Path)
		}
	}
}

func TestGroupIsTheMountAndAccessIsTheRequirement(t *testing.T) {
	routes := ParseSource(realisticRoutes)

	// A route on the staff group naming a permission belongs to staff and
	// demands that permission. One column cannot say both, which is why these
	// are two fields.
	r := findRoute(routes, "PUT", "/api/v1/users/:id/roles")
	if r.Group != "staff" {
		t.Errorf("group = %q, want staff", r.Group)
	}
	// RequireStaff is a gate, not a demand: it admits an ADMIN or anybody with
	// any permission, and the route behind it names the one it needs. Reporting
	// "staff or ADMIN or users.edit" would state a requirement that is not one.
	if got := r.Access.Summary(); strings.Contains(got, "staff") {
		t.Errorf("access = %q, and RequireStaff is not a demand", got)
	}

	if adm := findRoute(routes, "POST", "/api/v1/admin/flags"); adm == nil || adm.Group != "admin" {
		t.Errorf("admin group route = %v", adm)
	}
	if key := findRoute(routes, "GET", "/api/v1/public/products"); key == nil {
		t.Error("the api-key route was dropped")
	} else if key.Access.Summary() != "api key" {
		t.Errorf("api key route access = %q", key.Access.Summary())
	}
}

// A permission whose resource is a path parameter is only knowable per request.
// Naming it ":resource.view" says which half is dynamic, instead of pretending
// the route asks for nothing.
func TestDynamicPermissionNamesItsParameter(t *testing.T) {
	routes := ParseSource(realisticRoutes)
	r := findRoute(routes, "GET", "/api/v1/dashboard/:resource/stats")
	if r == nil {
		t.Fatal("the dynamic-permission route was dropped")
	}
	if got := r.Access.Summary(); got != ":resource.view" {
		t.Errorf("access = %q, want :resource.view", got)
	}
}

// Resources own their own route files and register through Mount, so a table
// built from routes.go alone misses most of what a real project serves.
func TestResourceRouteFilesAreIncluded(t *testing.T) {
	routesFile := writeProject(t)

	routes, err := ParseProject(routesFile)
	if err != nil {
		t.Fatalf("ParseProject: %v", err)
	}

	r := findRoute(routes, "GET", "/api/v1/widgets")
	if r == nil {
		t.Fatalf("m.Staff routes were not resolved; got %v", paths(routes))
	}
	if got := r.Access.Summary(); got != "ADMIN or widgets.view" {
		t.Errorf("access = %q", got)
	}
	if r.Group != "staff" {
		t.Errorf("group = %q, want staff", r.Group)
	}

	// m.Public is behind an API key, and a resource route on it inherits that
	// rather than reading as public.
	feed := findRoute(routes, "GET", "/api/v1/public/widgets/feed")
	if feed == nil {
		t.Fatalf("the m.Public route was not resolved; got %v", paths(routes))
	}
	if !feed.Access.Authenticated {
		t.Error("a route on the api-key group reported as needing nothing")
	}
}

func paths(routes []Route) []string {
	out := make([]string, 0, len(routes))
	for _, r := range routes {
		out = append(out, r.Method+" "+r.Path)
	}
	return out
}

// Three routes in a generated project pass their handler as a literal written
// across four lines. Skipping them left three gaps in the middle of the
// project's own API that the routes screen had to show as unknown.
func TestARouteWrittenAcrossLinesIsStillARoute(t *testing.T) {
	routes := ParseSource(`package routes

func Setup() *gin.Engine {
	v1 := r.Group("/api/v1")
	protected := v1.Group("")
	protected.Use(middleware.Auth(db, authService))
	protected.GET("/system/modules", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": cfg.Modules.Map()})
	})
	protected.POST("/auth/logout", authHandler.Logout)
}
`)

	r := findRoute(routes, "GET", "/api/v1/system/modules")
	if r == nil {
		t.Fatalf("the multi-line registration was skipped; got %v", paths(routes))
	}
	if !r.Access.Authenticated {
		t.Error("the joined route lost its group's middleware")
	}
	// And the line after it is still found, which is what breaks if the joiner
	// consumes too much.
	if findRoute(routes, "POST", "/api/v1/auth/logout") == nil {
		t.Error("the route after the multi-line one was swallowed")
	}
}

// Joining only the call being read, never every unbalanced line: a whole
// RegisterRoutes body collapsed into one line leaves only its first route
// findable, which is how the first attempt at this lost four of five.
func TestJoiningDoesNotSwallowTheEnclosingBlock(t *testing.T) {
	routes := ParseSource(`package routes

func init() {
	RegisterRoutes(func(m *Mount) {
		m.Admin.GET("/a", h.A)
		m.Admin.GET("/b", h.B)
		m.Admin.GET("/c", h.C)
	})
}
`)
	if len(routes) != 3 {
		t.Fatalf("found %d routes, want 3: %v", len(routes), paths(routes))
	}
}

// A const block's parens must not be joined away, or every prefix built from a
// constant resolves to {APIVersion}.
func TestAConstBlockSurvives(t *testing.T) {
	routes := ParseSource(`package routes

const (
	APIVersion = "v2"
)

func Setup() *gin.Engine {
	v := r.Group("/api/" + APIVersion)
	v.GET("/posts", postHandler.List)
}
`)
	if findRoute(routes, "GET", "/api/v2/posts") == nil {
		t.Errorf("the constant did not resolve; got %v", paths(routes))
	}
}
