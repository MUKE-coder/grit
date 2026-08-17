package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The --public flag.
//
// Generated CRUD sits behind auth, which is right for an admin resource and
// wrong for anything a customer reads. A storefront has no logged-in user, so
// calling the generated list endpoint returns a 401, and that is the first wall
// anyone building a public-facing app walks into.
//
// --public adds a second, narrower surface: a list and a get-by-slug, mounted
// outside the auth middleware and guarded by an API key. The protected routes
// are untouched, because the admin panel calls those and because gin panics at
// boot on two handlers for one method and path.
//
// What it deliberately does not do is expose the model. A public response is an
// allowlist struct, so a column added next month is private until somebody says
// otherwise. That is the opposite default to the admin surface, and it is the
// right one when the audience is the internet.

// PublicFields returns the fields safe to publish by default.
//
// Conservative on purpose. Names, slugs, prices, descriptions and images are
// what a catalogue page renders. Anything that smells like cost, margin,
// internal notes, stock counts or a foreign key is left out, and the developer
// opts it in by editing one struct.
//
// The alternative, publishing everything and expecting the developer to remove
// what should not be there, gets you a cost_price column on the internet the
// first time somebody adds one and forgets.
func PublicFields(fields []Field) (included, excluded []Field) {
	for _, f := range fields {
		if publishableByDefault(f) {
			included = append(included, f)
		} else {
			excluded = append(excluded, f)
		}
	}
	return included, excluded
}

// publishableByDefault decides whether one field goes out by default.
func publishableByDefault(f Field) bool {
	name := strings.ToLower(toSnakeCase(f.Name))

	// Never, whatever the type. Money you do not publish, notes you do not
	// publish, credentials you very much do not publish.
	for _, needle := range []string{
		"cost", "margin", "profit", "internal", "note", "secret",
		"password", "token", "commission", "supplier", "wholesale",
	} {
		if strings.Contains(name, needle) {
			return false
		}
	}

	// Held back for reasons other than secrecy.
	switch name {
	// A raw count is a business fact competitors enjoy, and a page almost
	// always wants "in stock" rather than "we have four left". Publish a
	// derived boolean instead, or add this back if you want the urgency.
	case "stock", "quantity", "inventory", "on_hand", "available":
		return false
	// The endpoint already filters on these, so the value is the same on every
	// row it will ever return. Publishing it is noise that reads like meaning.
	case "active", "published", "visible", "enabled", "archived", "is_active",
		"is_published", "is_visible":
		return false
	}

	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldSlug, FieldRichtext,
		FieldFloat, FieldInt, FieldBool, FieldSelect,
		FieldFile, FieldFiles, FieldDate, FieldDatetime, FieldStringArray:
		return true
	case FieldBelongsTo, FieldManyToMany:
		// A relation would publish the whole related record, including columns
		// nobody vetted. Expose a name or a slug from it deliberately instead.
		return false
	}
	return false
}

// writePublicResource emits the public handler and mounts its routes.
func (g *Generator) writePublicResource(names Names) error {
	if !g.Definition.Public {
		return nil
	}

	included, excluded := PublicFields(g.Definition.Fields)
	if len(included) == 0 {
		return fmt.Errorf(
			"--public on %s: none of its fields are publishable by default, so the "+
				"endpoint would return only ids. Rename the fields or write the public "+
				"handler by hand", names.Pascal)
	}

	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_public.go")

	// Written only when absent, unlike every other file this generator emits.
	//
	// The allowlist inside it is the developer's: the whole point is that they
	// add the one or two columns they do want published. Overwriting on a
	// regenerate would silently take those back, and the file itself tells the
	// reader it will not, so this is the code that keeps that promise.
	if fileExists(path) {
		fmt.Printf("  • %sinternal/handlers/%s_public.go exists, left alone\n",
			g.apiPrefix(), names.Snake)
		return nil
	}

	body := g.publicHandlerSource(names, included)
	if err := writeFileWithDirs(path, body); err != nil {
		return fmt.Errorf("writing the public %s handler: %w", names.Lower, err)
	}
	fmt.Printf("  ✓ %sinternal/handlers/%s_public.go (%d field(s) published)\n",
		g.apiPrefix(), names.Snake, len(included))

	if len(excluded) > 0 {
		// held, not names: the enclosing parameter is already called names, and
		// shadowing it here is how the message ended up pointing at a struct
		// that does not exist.
		held := make([]string, 0, len(excluded))
		for _, f := range excluded {
			held = append(held, toSnakeCase(f.Name))
		}
		fmt.Printf("    Held back: %s\n", strings.Join(held, ", "))
		fmt.Printf("    Add any of those to the public%s struct in that file to publish them.\n",
			names.Pascal)
	}

	// The routes are NOT mounted here. injectAll decides a resource is already
	// wired by finding its handler in routes.go, so a public route added first
	// makes it skip every CRUD route for that resource. Run calls
	// ensurePublicRoutes after injectAll instead.
	return nil
}

// apiPrefix is the display prefix for a path, matching Run's output.
func (g *Generator) apiPrefix() string {
	if g.Architecture == "single" {
		return ""
	}
	return "apps/api/"
}

// publicHandlerSource builds the handler.
func (g *Generator) publicHandlerSource(names Names, included []Field) string {
	var viewFields, assignments, searchable, sortable []string

	viewFields = append(viewFields, "\tID string `json:\"id\"`")
	assignments = append(assignments, "\t\tID: m.ID,")

	for _, f := range included {
		goName := toPascalCase(f.Name)
		jsonName := toSnakeCase(f.Name)
		viewFields = append(viewFields,
			fmt.Sprintf("\t%s %s `json:%q`", goName, publicGoType(f), jsonName))
		assignments = append(assignments,
			fmt.Sprintf("\t\t%s: m.%s,", goName, goName))

		switch FieldType(f.Type) {
		case FieldString, FieldText, FieldSlug:
			searchable = append(searchable, strconvQuote(jsonName))
			sortable = append(sortable, strconvQuote(jsonName)+": true")
		case FieldFloat, FieldInt, FieldDate, FieldDatetime:
			sortable = append(sortable, strconvQuote(jsonName)+": true")
		}
	}
	sortable = append(sortable, strconvQuote("created_at")+": true")

	slugField := "id"
	for _, f := range included {
		if FieldType(f.Type) == FieldSlug {
			slugField = toSnakeCase(f.Name)
			break
		}
	}

	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"` + g.Module + `/internal/models"
	"` + g.Module + `/internal/paginate"
)

// The public ` + names.Plural + ` surface.
//
// Generated by grit with --public. Two read endpoints, mounted outside auth and
// guarded by an API key, for clients with no logged-in user: a storefront, a
// mobile app, a public directory.
//
// Regenerating the resource does NOT overwrite this file, so the allowlist
// below is yours to edit.

// public` + names.Pascal + ` is what an anonymous caller is allowed to see.
//
// An allowlist rather than the model. A column added next month is private
// until somebody adds it here, which is the correct default when the audience
// is the internet, and the opposite of the admin surface's default.
type public` + names.Pascal + ` struct {
` + strings.Join(viewFields, "\n") + `
}

func toPublic` + names.Pascal + `(m models.` + names.Pascal + `) public` + names.Pascal + ` {
	return public` + names.Pascal + `{
` + strings.Join(assignments, "\n") + `
	}
}

// ListPublic handles GET /api/v1/public/` + names.Plural + `.
func (h *` + names.Pascal + `Handler) ListPublic(c *gin.Context) {
	// Scoped before paginate sees it, and deliberately not in the Filterable
	// list. A column listed there is settable from the query string, so
	// ?archived=true would hand back rows somebody took down on purpose.
	query := h.DB.Model(&models.` + names.Pascal + `{}).Where("archived_at IS NULL")

	res, err := paginate.List[models.` + names.Pascal + `](
		query,
		paginate.Bind(c),
		paginate.Config{
			Searchable: []string{` + strings.Join(searchable, ", ") + `},
			Sortable:   map[string]bool{` + strings.Join(sortable, ", ") + `},
			// No Filterable. Every filter a public caller gets is one you chose
			// to offer, rather than every column being queryable by default.
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to fetch ` + names.Plural + `",
		}})
		return
	}

	out := make([]public` + names.Pascal + `, 0, len(res.Data))
	for _, m := range res.Data {
		out = append(out, toPublic` + names.Pascal + `(m))
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "meta": res.Meta})
}

// GetPublic handles GET /api/v1/public/` + names.Plural + `/:key.
//
// Looks up by ` + slugField + `, so a public URL reads as something a person
// could type rather than a UUID.
func (h *` + names.Pascal + `Handler) GetPublic(c *gin.Context) {
	var item models.` + names.Pascal + `
	err := h.DB.Where("` + slugField + ` = ? AND archived_at IS NULL", c.Param("key")).
		First(&item).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "` + names.Pascal + ` not found",
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toPublic` + names.Pascal + `(item)})
}
`
}

// publicGoType maps a field to the Go type the view struct uses.
func publicGoType(f Field) string {
	switch FieldType(f.Type) {
	case FieldInt:
		return "int"
	case FieldFloat:
		return "float64"
	case FieldBool:
		return "bool"
	case FieldDate, FieldDatetime:
		return "time.Time"
	case FieldFile:
		return "models.FileRef"
	case FieldFiles:
		return "models.FileRefs"
	case FieldStringArray:
		return "datatypes.JSONSlice[string]"
	}
	return "string"
}

func strconvQuote(s string) string { return "\"" + s + "\"" }

// ensurePublicRoutes mounts the two endpoints under a public group.
//
// Run after injectAll for the same reason the workflow routes are: injectAll
// decides a resource is already wired by finding its handler in routes.go, and
// a public route added first makes it skip every CRUD route.
func (g *Generator) ensurePublicRoutes(names Names) error {
	if !g.Definition.Public {
		return nil
	}
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, names.Camel+"Handler.ListPublic") {
		return nil
	}

	route := fmt.Sprintf(
		"\t\tpublicAPI.GET(\"/%s\", %sHandler.ListPublic)\n"+
			"\t\tpublicAPI.GET(\"/%s/:key\", %sHandler.GetPublic)",
		names.Plural, names.Camel, names.Plural, names.Camel)

	if err := injectBefore(path, "// grit:routes:public", route); err != nil {
		fmt.Printf("  Could not mount the public routes: %v\n", err)
		fmt.Println("    This project predates --public. Add a public group to routes.go:")
		fmt.Println("      publicAPI := v1.Group(\"/public\")")
		fmt.Println("      publicAPI.Use(middleware.RequireAPIKey(db))")
		fmt.Println("      { /* grit:routes:public */ }")
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ GET /api/v1/public/%s and /api/v1/public/%s/:key (API key required)\n",
		names.Plural, names.Plural)
	return nil
}
