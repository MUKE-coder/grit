package routeparser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRoutes(t *testing.T, src string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "routes.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatalf("writing routes file: %v", err)
	}
	return path
}

func findRoute(routes []Route, method, path string) *Route {
	for i := range routes {
		if routes[i].Method == method && routes[i].Path == path {
			return &routes[i]
		}
	}
	return nil
}

// The regression this package was fixed for. Grit mounts everything under
// r.Group("/api/" + APIVersion); before the fix that expression did not match
// the literal-only pattern, the prefix evaluated to nothing, and every route
// was reported a segment short — /users instead of /api/v1/users. Plausible
// enough to act on, and wrong.
func TestParse_ResolvesPrefixBuiltFromAConstant(t *testing.T) {
	path := writeRoutes(t, `package routes

const APIVersion = "v1"

func Setup(r *gin.Engine) {
	v1 := r.Group("/api/" + APIVersion)
	v1.GET("/users", userHandler.List)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if findRoute(routes, "GET", "/api/v1/users") == nil {
		t.Errorf("expected /api/v1/users, got %v", routes)
	}
}

func TestParse_ResolvesConstantDeclaredInABlock(t *testing.T) {
	path := writeRoutes(t, `package routes

const (
	APIVersion = "v2"
	Unrelated  = "nope"
)

func Setup(r *gin.Engine) {
	v2 := r.Group("/api/" + APIVersion)
	v2.GET("/posts", postHandler.List)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if findRoute(routes, "GET", "/api/v2/posts") == nil {
		t.Errorf("expected /api/v2/posts, got %v", routes)
	}
}

// Nested groups must inherit from the receiver they were actually created on.
func TestParse_NestedGroupsInheritFromTheirReceiver(t *testing.T) {
	path := writeRoutes(t, `package routes

const APIVersion = "v1"

func Setup(r *gin.Engine) {
	v1 := r.Group("/api/" + APIVersion)
	admin := v1.Group("/admin")
	admin.GET("/blogs", blogHandler.List)
	v1.GET("/health", healthHandler.Check)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if findRoute(routes, "GET", "/api/v1/admin/blogs") == nil {
		t.Errorf("nested group did not inherit parent prefix, got %v", routes)
	}
	if findRoute(routes, "GET", "/api/v1/health") == nil {
		t.Errorf("sibling route lost its prefix, got %v", routes)
	}
}

// An identifier the parser cannot resolve must be visible in the output, not
// silently dropped — an incomplete path signals "parser hit something it does
// not understand", where a short one reads as correct.
func TestParse_UnresolvableIdentifierStaysVisible(t *testing.T) {
	path := writeRoutes(t, `package routes

func Setup(r *gin.Engine) {
	g := r.Group("/api/" + versionFromEnv())
	g.GET("/users", userHandler.List)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d: %v", len(routes), routes)
	}
	if got := routes[0].Path; got == "/api/users" {
		t.Errorf("unresolved identifier was silently dropped: %q", got)
	}
}

func TestParse_PlainLiteralPrefixStillWorks(t *testing.T) {
	path := writeRoutes(t, `package routes

func Setup(r *gin.Engine) {
	auth := r.Group("/api/auth")
	auth.POST("/login", authHandler.Login)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if findRoute(routes, "POST", "/api/auth/login") == nil {
		t.Errorf("expected /api/auth/login, got %v", routes)
	}
}

func TestParse_ClassifiesMiddlewareGroups(t *testing.T) {
	path := writeRoutes(t, `package routes

func Setup(r *gin.Engine) {
	pub := r.Group("/api")
	pub.GET("/health", healthHandler.Check)

	prot := r.Group("/api")
	prot.Use(middleware.Auth(cfg))
	prot.GET("/me", userHandler.Me)

	adm := r.Group("/api")
	adm.Use(middleware.RequireRole("admin"))
	adm.GET("/users", userHandler.List)
}
`)
	routes, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	for _, tc := range []struct{ method, path, group string }{
		{"GET", "/api/health", "public"},
		{"GET", "/api/me", "protected"},
		{"GET", "/api/users", "admin"},
	} {
		r := findRoute(routes, tc.method, tc.path)
		if r == nil {
			t.Errorf("route %s %s not found", tc.method, tc.path)
			continue
		}
		if r.Group != tc.group {
			t.Errorf("%s %s: group = %q, want %q", tc.method, tc.path, r.Group, tc.group)
		}
	}
}

func TestParse_MissingFileReturnsError(t *testing.T) {
	if _, err := Parse(filepath.Join(t.TempDir(), "nope.go")); err == nil {
		t.Error("expected an error for a missing routes file")
	}
}
