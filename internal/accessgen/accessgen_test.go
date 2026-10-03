package accessgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const routesGo = `package routes

const APIVersion = "v1"

func Setup() *gin.Engine {
	r := gin.New()
	r.GET("/api/health", healthCheck)

	v1 := r.Group("/api/" + APIVersion)

	protected := v1.Group("")
	protected.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	protected.GET("/auth/me", authHandler.Me)

	staff := v1.Group("")
	staff.Use(middleware.APIKeyOrAuth(db, middleware.Auth(db, authService)))
	staff.Use(middleware.RequireStaff())
	staff.GET("/widgets", middleware.RequireRole("ADMIN", "perm:widgets.view"), h.List)

	mountResources(&Mount{
		Engine:    r,
		V1:        v1,
		Protected: protected,
		Staff:     staff,
	})
	return r
}
`

func project(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "apps", "api", "internal", "routes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "routes.go"), []byte(routesGo), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func registry(t *testing.T, root string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, "apps", "api", "internal", "access", "registry.go"))
	if err != nil {
		t.Fatalf("reading the registry: %v", err)
	}
	return string(body)
}

func TestGenerateWritesWhatTheRouteFilesSay(t *testing.T) {
	root := project(t)

	res, err := Generate(root)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !res.Changed || res.Routes != 3 {
		t.Fatalf("Result = %+v, want 3 routes and Changed", res)
	}

	// Matched with runs of spaces collapsed, because the output is gofmt'd and
	// gofmt decides the field alignment. A test that encoded the padding would
	// fail the next time a longer field name changed the column width, which
	// says nothing about whether the table is right.
	out := collapse(registry(t, root))
	for _, want := range []string{
		Header,
		`"GET /api/health": {`,
		`Summary: "public",`,
		`"GET /api/v1/auth/me": {`,
		`Authenticated: true,`,
		`"GET /api/v1/widgets": {`,
		`Roles: []string{"ADMIN"},`,
		`Permissions: []string{"widgets.view"},`,
		`Summary: "ADMIN or widgets.view",`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the registry is missing %q", want)
		}
	}

	// The whole point: nothing authenticated may be written as public.
	if strings.Contains(out, `"GET /api/v1/auth/me": { Group: "public"`) {
		t.Error("/auth/me was written as public")
	}
}

// collapse squeezes runs of whitespace to one space, so an assertion can name
// the fields without naming gofmt's column widths.
func collapse(src string) string {
	return strings.Join(strings.Fields(src), " ")
}

// The output has to be byte-identical for the same input, or Generate cannot
// skip a write and Stale cannot compare by content.
func TestGenerateIsDeterministic(t *testing.T) {
	root := project(t)

	if _, err := Generate(root); err != nil {
		t.Fatal(err)
	}
	first := registry(t, root)

	res, err := Generate(root)
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Error("the second run rewrote an identical file")
	}
	if registry(t, root) != first {
		t.Error("two runs over the same source produced different files")
	}
}

func TestStaleFollowsTheRouteFiles(t *testing.T) {
	root := project(t)
	routesFile := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")

	if _, err := Generate(root); err != nil {
		t.Fatal(err)
	}
	stale, err := Stale(root)
	if err != nil {
		t.Fatal(err)
	}
	if stale {
		t.Error("stale immediately after generating")
	}

	// A route added by hand, which is the normal way a project grows past its
	// generated resources and the one case no command catches.
	src := strings.Replace(routesGo,
		`	protected.GET("/auth/me", authHandler.Me)`,
		"\tprotected.GET(\"/auth/me\", authHandler.Me)\n\tprotected.POST(\"/billing/portal\", billingHandler.Portal)",
		1)
	if err := os.WriteFile(routesFile, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	stale, err = Stale(root)
	if err != nil {
		t.Fatal(err)
	}
	if !stale {
		t.Error("a hand-added route did not make the table stale")
	}
}

// A project with no routes.go is not an error. `grit new --api` outside a
// project, and the desktop shape, both reach here.
func TestNoRoutesFileIsNotAnError(t *testing.T) {
	res, err := Generate(t.TempDir())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Changed || res.Routes != 0 {
		t.Errorf("Result = %+v, want zero", res)
	}

	stale, err := Stale(t.TempDir())
	if err != nil || stale {
		t.Errorf("Stale = %v, %v; want false, nil", stale, err)
	}
}

// Overwriting somebody's code because it sits at a path we wanted is the one
// thing the upgrade machinery exists not to do.
func TestAHandWrittenFileIsRefusedRatherThanOverwritten(t *testing.T) {
	root := project(t)
	dir := filepath.Join(root, "apps", "api", "internal", "access")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	mine := "package access\n\n// Mine, thanks.\nfunc Of(string, string) {}\n"
	target := filepath.Join(dir, "registry.go")
	if err := os.WriteFile(target, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Generate(root); err == nil {
		t.Fatal("Generate overwrote a hand-written file")
	}
	body, err := os.ReadFile(target)
	if err != nil || string(body) != mine {
		t.Error("the hand-written file was changed")
	}

	// And it does not then report the project as stale, which would send
	// somebody to run the command that just refused.
	stale, err := Stale(root)
	if err != nil {
		t.Fatal(err)
	}
	if stale {
		t.Error("a hand-written file was reported as a stale generated one")
	}
}
