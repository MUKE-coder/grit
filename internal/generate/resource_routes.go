package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeResourceRoutes writes internal/routes/<resource>_routes.go: everything
// this resource mounts, in one file, next to nothing else.
//
// This replaces four separate injections into routes.go (the handler
// construction, the protected block, the admin block and sometimes the public
// one). Those four edits landed hundreds of lines apart in a file that passed
// a thousand lines with a handful of resources, so adding a route by hand
// meant finding which of four blocks it belonged in, and the generator had
// four chances to inject into the wrong one.
//
// Now: one file per resource, registered from its own init(). Creating the
// file mounts the resource, deleting it unmounts the resource, and routes.go
// does not change either way.
//
// Returns false when the project has no route registry, which means a project
// generated before this existed and not yet upgraded. The caller falls back to
// the marker injections so those projects keep working.
func (g *Generator) writeResourceRoutes(names Names) (bool, error) {
	routesDir := filepath.Join(g.APIRoot(), "internal", "routes")
	if !fileExists(filepath.Join(routesDir, "resources.go")) {
		return false, nil
	}

	path := filepath.Join(routesDir, names.Snake+"_routes.go")

	// A file already here is either one we wrote, in which case rewriting it
	// with the current templates is the point, or one somebody has edited, in
	// which case the manifest guard holds it back and says so. Either way the
	// decision belongs to the guard, not to a check here.
	body, err := g.resourceRoutesSource(names)
	if err != nil {
		return false, err
	}
	if err := writeFileWithDirs(path, body); err != nil {
		return false, fmt.Errorf("writing %s: %w", filepath.Base(path), err)
	}
	return true, nil
}

// resourceRoutesSource builds the file.
func (g *Generator) resourceRoutesSource(names Names) (string, error) {
	var b strings.Builder

	hasFileFields := false
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			hasFileFields = true
			break
		}
	}

	// Imports. Built as a list rather than a template with holes: a resource
	// with no roles must not import middleware, and an unused import is a
	// compile error rather than a warning in Go.
	imports := []string{fmt.Sprintf("%q", g.Module+"/internal/handlers")}
	if len(g.Roles) > 0 {
		imports = append(imports, fmt.Sprintf("%q", g.Module+"/internal/middleware"))
	}

	fmt.Fprintf(&b, "package routes\n\nimport (\n")
	for _, imp := range imports {
		fmt.Fprintf(&b, "\t%s\n", imp)
	}
	fmt.Fprintf(&b, ")\n\n")

	fmt.Fprintf(&b, "// %s routes.\n", names.PluralPascal)
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "// This is the whole surface for %s: the handler, and every path that\n", names.Plural)
	fmt.Fprintf(&b, "// reaches it. Add a route by adding a line here. Remove the resource by\n")
	fmt.Fprintf(&b, "// deleting this file; nothing else refers to it.\n")
	fmt.Fprintf(&b, "//\n")
	fmt.Fprintf(&b, "// m.Public is outside the auth middleware and behind an API key, m.Protected\n")
	fmt.Fprintf(&b, "// takes a JWT or an API key, and m.Admin also requires the ADMIN role.\n")
	fmt.Fprintf(&b, "func init() {\n")
	fmt.Fprintf(&b, "\tRegisterRoutes(func(m *Mount) {\n")

	// Handler construction.
	fmt.Fprintf(&b, "\t\th := &handlers.%sHandler{\n\t\t\tDB: m.DB,\n", names.Pascal)
	if hasFileFields {
		// Without Storage the Create and Update flows skip the S3 cleanup on
		// replace and never mark uploads claimed, both silently.
		fmt.Fprintf(&b, "\t\t\tStorage: m.Svc.Storage,\n")
	}
	fmt.Fprintf(&b, "\t\t}\n\n")

	// Public routes, when the resource asked for them.
	if g.Definition.Public {
		fmt.Fprintf(&b, "\t\t// Read-only, outside auth, API key required.\n")
		if g.Definition.Tree {
			// Before /:key so the reader can see that a static segment wins.
			// Gin would route it correctly in either order.
			fmt.Fprintf(&b, "\t\tm.Public.GET(\"/%s/tree\", h.TreePublic)\n", names.Plural)
		}
		fmt.Fprintf(&b, "\t\tm.Public.GET(\"/%s\", h.ListPublic)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Public.GET(\"/%s/:key\", h.GetPublic)\n", names.Plural)
		if hasParent(g.Definition) {
			fmt.Fprintf(&b, "\t\tm.Public.GET(\"/%s/:key/related\", h.RelatedPublic)\n", names.Plural)
		}
		fmt.Fprintf(&b, "\n")
	}

	if len(g.Roles) > 0 {
		roleArgs := make([]string, len(g.Roles))
		for i, r := range g.Roles {
			roleArgs[i] = fmt.Sprintf("%q", r)
		}
		fmt.Fprintf(&b, "\t\t// Restricted to %s.\n", strings.Join(g.Roles, ", "))
		fmt.Fprintf(&b, "\t\tg := m.Protected.Group(\"/%s\")\n", names.Plural)
		fmt.Fprintf(&b, "\t\tg.Use(middleware.RequireRole(%s))\n", strings.Join(roleArgs, ", "))
		// Every verb on the group, so the role guard cannot be true of some
		// routes on this resource and not others.
		for _, r := range []string{
			`g.GET("", h.List)`,
			`g.GET("/export", h.Export)`,
			`g.GET("/:id", h.GetByID)`,
			`g.GET("/:id/pdf", h.PDF)`,
			`g.POST("", h.Create)`,
			`g.PUT("/:id", h.Update)`,
			`g.PATCH("/:id", h.Patch)`,
			`g.DELETE("/:id", h.Delete)`,
			`g.POST("/bulk", h.Bulk)`,
		} {
			fmt.Fprintf(&b, "\t\t%s\n", r)
		}
	} else {
		fmt.Fprintf(&b, "\t\tm.Protected.GET(\"/%s\", h.List)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.GET(\"/%s/export\", h.Export)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.POST(\"/%s/import\", h.Import)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.GET(\"/%s/import/template\", h.Template)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.GET(\"/%s/:id\", h.GetByID)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.GET(\"/%s/:id/pdf\", h.PDF)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.POST(\"/%s\", h.Create)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.PUT(\"/%s/:id\", h.Update)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Protected.PATCH(\"/%s/:id\", h.Patch)\n", names.Plural)
		fmt.Fprintf(&b, "\n")
		// Bulk sits with DELETE rather than with PATCH: it can delete, and a
		// route is only as protected as its most destructive branch.
		fmt.Fprintf(&b, "\t\tm.Admin.DELETE(\"/%s/:id\", h.Delete)\n", names.Plural)
		fmt.Fprintf(&b, "\t\tm.Admin.POST(\"/%s/bulk\", h.Bulk)\n", names.Plural)
	}

	fmt.Fprintf(&b, "\t})\n}\n")
	return b.String(), nil
}

// removeResourceRoutes deletes a resource's route file.
//
// Returns true when it removed one, so the caller can skip the line-by-line
// unpicking of routes.go that a pre-split project still needs.
func removeResourceRoutes(apiRoot, snake string) (bool, error) {
	path := filepath.Join(apiRoot, "internal", "routes", snake+"_routes.go")
	if !fileExists(path) {
		return false, nil
	}
	if err := os.Remove(path); err != nil {
		return false, fmt.Errorf("removing %s: %w", filepath.Base(path), err)
	}
	return true, nil
}
