package generate

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// --append-only: rows are created and read, never changed or deleted.
//
// Found building a double-entry ledger. The generator produced a JournalEntry
// with PUT, PATCH, DELETE, a bulk endpoint and a CSV import, a soft-delete
// column and an admin with Edit and Delete buttons: the right shape for a CRUD
// resource and the wrong one for a posted accounting record. Hand-rolling the
// fix took a package of GORM callbacks, a main.go edit the framework offered no
// place for, and it still lost: one UPDATE typed into GORM Studio's SQL editor
// unbalanced the books, because raw SQL never passes a GORM callback.
//
// So the flag does all of it: read and create routes only, a GORM guard on the
// connection, a trigger on the table installed by grit migrate, and an admin
// that offers create and view.

// prepareAppendOnly checks the project can honour the flag, and wires it up
// when it is old enough not to be.
func (g *Generator) prepareAppendOnly() error {
	if g.Definition.Tree {
		return fmt.Errorf("--append-only and --tree cannot be combined: moving a node rewrites its path, which is an update")
	}
	api := g.APIRoot()
	// Per-resource route files arrived before this flag did, and the route
	// half of it is written for that layout only.
	if !fileExists(filepath.Join(api, "internal", "routes", "resources.go")) {
		return fmt.Errorf("--append-only needs the per-resource route layout, which this project predates.\n\n" +
			"Run grit upgrade first, then generate again.")
	}
	if err := scaffold.WriteAppendOnlyPackage(api, g.Module, false); err != nil {
		return err
	}
	return scaffold.EnsureAppendOnlyWiring(api, g.Module)
}

// writeAppendOnlyRoutes emits the routes for an append-only resource: read and
// create. Closes the Mount function the caller opened.
//
// Update, Patch, Delete, Bulk and Import stay in the handler, unrouted. If the
// business later decides a correction may be made in place, that is one route
// rather than a regenerate over code someone has edited; the table refuses the
// write regardless until the model stops registering itself.
func (g *Generator) writeAppendOnlyRoutes(w io.Writer, names Names) {
	routes := []string{
		`GET("%s", h.List)`,
		`GET("%s/export", h.Export)`,
		`GET("%s/:id", h.GetByID)`,
		`GET("%s/:id/pdf", h.PDF)`,
		`POST("%s", h.Create)`,
	}

	if len(g.Roles) > 0 {
		roleArgs := make([]string, len(g.Roles))
		for i, r := range g.Roles {
			roleArgs[i] = fmt.Sprintf("%q", r)
		}
		fmt.Fprintf(w, "\t\t// Append-only, restricted to %s.\n", strings.Join(g.Roles, ", "))
		fmt.Fprintf(w, "\t\tg := m.Protected.Group(\"/%s\")\n", names.Plural)
		fmt.Fprintf(w, "\t\tg.Use(middleware.RequireRole(%s))\n", strings.Join(roleArgs, ", "))
		for _, r := range routes {
			fmt.Fprintf(w, "\t\tg."+r+"\n", "")
		}
	} else {
		fmt.Fprintf(w, "\t\t// Append-only: created and read, never changed or deleted.\n")
		for _, r := range routes {
			fmt.Fprintf(w, "\t\tm.Protected."+r+"\n", "/"+names.Plural)
		}
	}
	fmt.Fprintf(w, "\t})\n}\n")
}

// appendOnlyDocs drops the PUT and DELETE entries from the API reference for an
// append-only resource. They are the last two blocks, so everything from the
// first PUT goes. Documenting a route that does not exist is how an API
// reference stops being believed.
func (g *Generator) appendOnlyDocs(docsRoutes string) string {
	if !g.Definition.AppendOnly {
		return docsRoutes
	}
	if i := strings.Index(docsRoutes, "\tdocs.Route(\"PUT "); i >= 0 {
		return strings.TrimRight(docsRoutes[:i], "\n")
	}
	return docsRoutes
}

// appendOnlyImports puts internal/appendonly first in the model's project
// imports, which is where it sorts.
func (g *Generator) appendOnlyImports(imports []string) []string {
	if !g.Definition.AppendOnly {
		return imports
	}
	return append([]string{fmt.Sprintf("\"%s/internal/appendonly\"", g.Module)}, imports...)
}

// appendOnlyModelInit registers the model with the guard.
//
// From the model file, not the routes file, because the migrate command imports
// models and not routes, and the trigger is installed at migrate time. Removing
// the resource deletes this file, and the registration with it.
func (g *Generator) appendOnlyModelInit(names Names) string {
	if !g.Definition.AppendOnly {
		return ""
	}
	return fmt.Sprintf(`
// init marks %[1]s append-only. GORM refuses to update or delete it, and
// grit migrate puts a trigger on the table for everything that is not GORM. A
// correction is a new row. See internal/appendonly.
func init() { appendonly.Register(&%[1]s{}) }
`, names.Pascal)
}

const editableBulkActions = `    // Shown once rows are ticked. Drop "archive" here and the Archived tab
    // goes with it; the model keeps its archived_at either way.
    bulkActions: ["edit", "archive", "restore", "export", "delete"],`

const appendOnlyTableActions = `    // Append-only: rows are created and read, never edited or deleted. The API
    // has no route for either and the table refuses both, so the buttons would
    // only ever produce a refusal.
    actions: ["create", "view"],
    bulkActions: ["export"],`

// adminDefinitionContent is the admin resource definition, with the table's
// actions cut down to create and view for an append-only resource.
func (g *Generator) adminDefinitionContent(names Names) string {
	content := g.resourceDefinitionFileContent(names)
	if g.Definition.AppendOnly {
		content = strings.Replace(content, editableBulkActions, appendOnlyTableActions, 1)
	}
	return content
}
