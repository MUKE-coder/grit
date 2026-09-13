package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	// Delete and bulk go on the staff group, asking for the resource's delete
	// permission, where the project has that group. Role-restricted and
	// append-only resources have their own arrangements.
	staffGroup := g.projectHasStaffGroup() && len(g.Roles) == 0
	needStaff := staffGroup && !g.Definition.AppendOnly
	// An owned append-only resource has no delete and protects everything else,
	// so it is the one shape on a staff project that names no guard.
	staffGate := staffGroup && !(g.Definition.AppendOnly && g.scopedRows())

	imports := []string{fmt.Sprintf("%q", g.Module+"/internal/handlers")}
	if len(g.Roles) > 0 || staffGate {
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
	if staffGate {
		fmt.Fprintf(&b, "// takes a JWT or an API key, and m.Staff also requires the permission its\n")
		fmt.Fprintf(&b, "// route names (an ADMIN holds every one).\n")
	} else {
		fmt.Fprintf(&b, "// takes a JWT or an API key, and m.Admin also requires the ADMIN role.\n")
	}
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

	// Append-only resources get read and create, and nothing else.
	if g.Definition.AppendOnly {
		g.writeAppendOnlyRoutes(&b, names, staffGroup)
		return b.String(), nil
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
		routes := []string{
			`g.GET("", h.List)`,
			`g.GET("/export", h.Export)`,
			`g.GET("/:id", h.GetByID)`,
			`g.GET("/:id/pdf", h.PDF)`,
			`g.POST("", h.Create)`,
			`g.PUT("/:id", h.Update)`,
			`g.PATCH("/:id", h.Patch)`,
			`g.DELETE("/:id", h.Delete)`,
			`g.POST("/bulk", h.Bulk)`,
		}
		if g.Definition.WorkflowField() != nil {
			routes = append(routes, `g.GET("/workflow", h.Workflow)`,
				`g.POST("/:id/transitions/:action", h.Transition)`)
		}
		for _, r := range routes {
			fmt.Fprintf(&b, "\t\t%s\n", r)
		}
	} else {
		routes := append(readCreateRoutes(),
			gatedRoute{"PUT", "/:id", "edit", "Update"},
			gatedRoute{"PATCH", "/:id", "edit", "Patch"})
		if g.Definition.WorkflowField() != nil {
			// The service also checks each transition's own permission.
			routes = append(routes,
				gatedRoute{"GET", "/workflow", "view", "Workflow"},
				gatedRoute{"POST", "/:id/transitions/:action", "edit", "Transition"})
		}
		for _, r := range routes {
			g.writeGatedRoute(&b, names.Plural, r, staffGroup)
		}
		fmt.Fprintf(&b, "\n")
		// Bulk sits with DELETE rather than with PATCH: it can delete, and a
		// route is only as protected as its most destructive branch.
		if needStaff {
			// A role granted <resource>.delete can delete without being made a
			// full admin. Bulk can delete too, so it asks for the same.
			fmt.Fprintf(&b, "\t\tm.Staff.DELETE(\"/%s/:id\", middleware.RequireRole(\"ADMIN\", \"perm:%s.delete\"), h.Delete)\n", names.Plural, names.Plural)
			fmt.Fprintf(&b, "\t\tm.Staff.POST(\"/%s/bulk\", middleware.RequireRole(\"ADMIN\", \"perm:%s.delete\"), h.Bulk)\n", names.Plural, names.Plural)
		} else {
			fmt.Fprintf(&b, "\t\tm.Admin.DELETE(\"/%s/:id\", h.Delete)\n", names.Plural)
			fmt.Fprintf(&b, "\t\tm.Admin.POST(\"/%s/bulk\", h.Bulk)\n", names.Plural)
		}
	}

	fmt.Fprintf(&b, "\t})\n}\n")
	return b.String(), nil
}

// gatedRoute is one generated route and the permission it asks for.
type gatedRoute struct{ method, path, perm, handler string }

// readCreateRoutes are the routes every resource has, append-only ones included.
func readCreateRoutes() []gatedRoute {
	return []gatedRoute{
		{"GET", "", "view", "List"},
		{"GET", "/export", "view", "Export"},
		{"POST", "/import", "create", "Import"},
		{"GET", "/import/template", "view", "Template"},
		{"GET", "/:id", "view", "GetByID"},
		{"GET", "/:id/pdf", "view", "PDF"},
		{"POST", "", "create", "Create"},
	}
}

// writeGatedRoute writes one route of a resource that has no --roles.
//
// An owned resource's service scopes every query to the caller, so a signed-in
// account sees and changes only its own rows, and the route can be protected.
// Anything else is shared data, and "signed in" is not a permission on an app
// with open registration: that is how a scaffolded address book handed every
// contact to anybody who made an account. On the staff group each verb asks for
// the permission the roles UI grants for it; a project without that group keeps
// shared data for ADMIN.
func (g *Generator) writeGatedRoute(w io.Writer, plural string, r gatedRoute, staffGroup bool) {
	switch {
	case g.scopedRows():
		fmt.Fprintf(w, "\t\tm.Protected.%s(\"/%s%s\", h.%s)\n", r.method, plural, r.path, r.handler)
	case staffGroup:
		fmt.Fprintf(w, "\t\tm.Staff.%s(\"/%s%s\", middleware.RequireRole(\"ADMIN\", \"perm:%s.%s\"), h.%s)\n",
			r.method, plural, r.path, plural, r.perm, r.handler)
	default:
		fmt.Fprintf(w, "\t\tm.Admin.%s(\"/%s%s\", h.%s)\n", r.method, plural, r.path, r.handler)
	}
}

// scopedRows reports whether every query on the resource is already narrowed
// to the caller: to their own rows (--owned-by), or to their organization's
// (--tenant-owned, applied by the tenant middleware). Such a resource can sit on
// the protected group, since signing in reaches nothing that is not yours.
func (g *Generator) scopedRows() bool {
	return g.Definition.IsOwned() || g.Definition.TenantOwned
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

// projectHasStaffGroup reports whether the project's route registry has the
// staff group, which a project from before v3.220.0 does not.
func (g *Generator) projectHasStaffGroup() bool {
	data, err := os.ReadFile(filepath.Join(g.APIRoot(), "internal", "routes", "resources.go"))
	return err == nil && staffField.Match(data)
}

var staffField = regexp.MustCompile(`Staff\s+\*gin\.RouterGroup`)
