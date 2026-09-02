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
		// The other half of that promise, which is easy to be bitten by: the
		// file keeps the allowlist you edited AND the generator features it was
		// written with. A handler from before filters existed does not gain
		// them, and the symptom is a query parameter that quietly does nothing.
		fmt.Println("    Your allowlist is kept, and so is everything else in there:")
		fmt.Println("    newer generator features do not reach an existing file.")
		fmt.Println("    Delete it and regenerate to pick them up, then re-add any")
		fmt.Println("    columns you had published by hand.")
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

	// A tree publishes its shape, because a storefront needs it to render a menu
	// and to ask for "this category and everything under it".
	if g.Definition.Tree {
		viewFields = append(viewFields,
			"\tParentID string `json:\"parent_id\"`",
			"\tDepth int `json:\"depth\"`",
			// Filled by GetPublic rather than by the mapper, because it costs a
			// query and a list of a hundred categories should not run a hundred
			// of them.
			"\tDescendantIDs []string `json:\"descendant_ids,omitempty\"`")
		assignments = append(assignments,
			"\t\tParentID: deref"+names.Pascal+"ID(m.ParentID),",
			"\t\tDepth: m.Depth,")
	}

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

	// Filters, built from the published fields and nothing else. That is the
	// whole rule, and it is what keeps the two lists honest as fields come and
	// go: a column safe to show is safe to filter on, and a column held back
	// from the response must not be reachable through a query string either.
	//
	// Text and richtext are left out because equality on a description is
	// never what a caller means, and search already covers them.
	var filterable, rangeFilterable []string
	for _, f := range included {
		col := strconvQuote(toSnakeCase(f.Name)) + ": true"
		switch FieldType(f.Type) {
		case FieldString, FieldSlug, FieldBool, FieldToggle, FieldSelect, FieldRadio:
			filterable = append(filterable, col)
		case FieldInt, FieldFloat:
			// Both: ?year=2024 is an equality question, ?price_min=50 a window.
			filterable = append(filterable, col)
			rangeFilterable = append(rangeFilterable, col)
		}
	}

	// Foreign keys, which the response holds back on purpose (publishing a
	// relation publishes a whole related record nobody vetted) but which a
	// storefront cannot do without: a category page is ?category_id=<id>.
	//
	// Filtering by an id is not the same as publishing the relation. The id
	// identifies a row the endpoint was already willing to return, and it
	// reveals nothing about the parent beyond the fact that it exists.
	var inFilterable []string
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			// IN rather than equality, because a category page filtered on a
			// tree wants "Electronics and everything under it", which is a list
			// of ids. IN with one value is the same query equality would have
			// been, so nothing is lost for the flat case.
			inFilterable = append(inFilterable, strconvQuote(f.FKColumnName())+": true")
		}
	}

	slugField := "id"
	for _, f := range included {
		if FieldType(f.Type) == FieldSlug {
			slugField = toSnakeCase(f.Name)
			break
		}
	}

	std, third, local := publicImports(g.Module, included)
	// The related endpoint clamps ?limit, which needs strconv. Only emitted
	// when that endpoint is, because an unused import is a build failure.
	if hasParent(g.Definition) {
		std += "\t\"strconv\"\n"
	}

	return `package handlers

import (
	"net/http"
` + std + `
	"github.com/gin-gonic/gin"
` + third + `
	"` + g.Module + `/internal/models"
	"` + g.Module + `/internal/paginate"
` + local + `)

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
//
// Supports ?search=, ?sort_by= and ?sort_order=, ?page= and ?page_size=,
// equality filters on the published columns, and a ?x_min= / ?x_max= window on
// the numeric ones.
func (h *` + names.Pascal + `Handler) ListPublic(c *gin.Context) {
	// Scoped before paginate sees it, and archived_at is deliberately absent
	// from the filter lists below: a column listed there is settable from the
	// query string, so ?archived=true would hand back rows somebody took down
	// on purpose.
	query := h.DB.Model(&models.` + names.Pascal + `{}).Where("archived_at IS NULL")

	res, err := paginate.List[models.` + names.Pascal + `](
		query,
		paginate.Bind(c),
		paginate.Config{
			Searchable: []string{` + strings.Join(searchable, ", ") + `},
			Sortable:   map[string]bool{` + strings.Join(sortable, ", ") + `},
			// Published columns only, plus foreign keys. Anything not listed
			// is ignored rather than rejected, so an unknown query param is
			// never an error.
			Filterable:      map[string]bool{` + strings.Join(filterable, ", ") + `},
			RangeFilterable: map[string]bool{` + strings.Join(rangeFilterable, ", ") + `},
			// Foreign keys accept a comma-separated list, so one request can
			// ask for a whole category subtree.
			InFilterable: map[string]bool{` + strings.Join(inFilterable, ", ") + `},
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
` + treeGetPublicPrelude(names, g.Definition) + `	c.JSON(http.StatusOK, gin.H{"data": ` + detailPayload(names, g.Definition) + `})
}
` + treePublic(names, g.Definition) + relatedPublic(names, g.Definition, slugField)
}

// hasParent reports whether the resource has a belongs_to for the related
// endpoint to key on.
func hasParent(def *ResourceDefinition) bool {
	for _, f := range def.Fields {
		if f.IsBelongsTo() {
			return true
		}
	}
	return false
}

// relatedPublic emits the "similar items" endpoint, or nothing.
//
// It needs a parent to define similarity, so a resource with no belongs_to
// gets no endpoint rather than one that returns an arbitrary set. Which
// relation defines "similar" is the generator's choice and not the caller's:
// that keeps it a single bounded query, and stops the endpoint becoming a way
// to filter on a column nobody published.
func relatedPublic(names Names, def *ResourceDefinition, slugField string) string {
	var parent *Field
	for i := range def.Fields {
		if def.Fields[i].IsBelongsTo() {
			parent = &def.Fields[i]
			break
		}
	}
	if parent == nil {
		return ""
	}

	// Derived the same way the model does it, from the same base name, rather
	// than by PascalCasing the column: toPascalCase("category_id") gives
	// CategoryId, and the model field is CategoryID.
	base := strings.TrimSuffix(parent.Name, "_id")
	fk := toSnakeCase(base) + "_id"
	fkField := toPascalCase(base) + "ID"

	// A self-referential parent is a nullable *string, every other one is a
	// plain string, so the emitted comparison differs. A tree resource with a
	// related endpoint would otherwise not compile:
	//   invalid operation: item.ParentID != "" (mismatched types *string and untyped string)
	selfRef := parent.RelatedModelName() == toPascalCase(def.Name)
	parentSet := "item." + fkField + ` != ""`
	parentValue := "item." + fkField
	if selfRef {
		parentSet = "item." + fkField + " != nil && *item." + fkField + ` != ""`
		parentValue = "*item." + fkField
	}

	return `
// RelatedPublic handles GET /api/v1/public/` + names.Plural + `/:key/related.
//
// The "similar items" strip on a detail page: others sharing this one's
// ` + parent.Name + `, newest first, this one excluded.
//
// Capped at 24 however large ?limit= asks, because this endpoint exists to
// fill a row of cards and an uncapped limit on a public route is a free way
// to make the database do work.
func (h *` + names.Pascal + `Handler) RelatedPublic(c *gin.Context) {
	var item models.` + names.Pascal + `
	err := h.DB.Where("` + slugField + ` = ? AND archived_at IS NULL", c.Param("key")).
		First(&item).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "` + names.Pascal + ` not found",
		}})
		return
	}

	limit := 8
	if raw := c.Query("limit"); raw != "" {
		if n, convErr := strconv.Atoi(raw); convErr == nil && n > 0 {
			limit = n
		}
	}
	if limit > 24 {
		limit = 24
	}

	query := h.DB.Model(&models.` + names.Pascal + `{}).
		Where("id <> ? AND archived_at IS NULL", item.ID)
	// A row with no ` + parent.Name + ` has no siblings, and an empty strip on
	// a detail page looks broken. Falling back to the newest rows is a worse
	// recommendation than a real match and a better one than nothing.
	if ` + parentSet + ` {
		query = query.Where("` + fk + ` = ?", ` + parentValue + `)
	}

	var rows []models.` + names.Pascal + `
	if err := query.Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to fetch related ` + names.Plural + `",
		}})
		return
	}

	out := make([]public` + names.Pascal + `, 0, len(rows))
	for _, m := range rows {
		out = append(out, toPublic` + names.Pascal + `(m))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
`
}

// publicGoType maps a field to the Go type the view struct uses.
//
// These have to be the types the model actually declares, not plausible
// guesses at them. Getting FileRefs wrong produced a public handler that did
// not compile for any resource with an image on it, which is most of the
// resources anyone would want public.
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
		return "*files.FileRef"
	case FieldFiles:
		return "files.FileRefs"
	case FieldStringArray:
		return "datatypes.JSONSlice[string]"
	}
	return "string"
}

// publicImports returns the extra imports the view struct needs, one string
// per group: stdlib, third party, this module.
//
// Computed from the fields rather than always emitted, because an unused
// import is a build failure and a resource of plain strings needs none of
// these. Returned per group because gofmt sorts inside a group but never moves
// a line between them, so a local import emitted next to "net/http" stays
// there, wrong, forever.
func publicImports(module string, included []Field) (string, string, string) {
	var needTime, needFiles, needDatatypes bool
	for _, f := range included {
		switch FieldType(f.Type) {
		case FieldDate, FieldDatetime:
			needTime = true
		case FieldFile, FieldFiles:
			needFiles = true
		case FieldStringArray:
			needDatatypes = true
		}
	}

	var std, third, local []string
	if needTime {
		std = append(std, "	\"time\"")
	}
	if needDatatypes {
		third = append(third, "	\"gorm.io/datatypes\"")
	}
	if needFiles {
		local = append(local, "	\""+module+"/internal/files\"")
	}

	join := func(lines []string) string {
		out := ""
		for _, line := range lines {
			out += line + "\n"
		}
		return out
	}
	return join(std), join(third), join(local)
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
	// A project with a route registry already has these in the resource's own
	// file, written in one pass with the rest of its routes. Injecting them
	// here as well would mount each path twice, and Gin panics on a duplicate.
	if fileExists(filepath.Join(g.APIRoot(), "internal", "routes", names.Snake+"_routes.go")) {
		return nil
	}
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	treeRoute := ""
	if g.Definition.Tree {
		treeRoute = fmt.Sprintf("		publicAPI.GET(\"/%s/tree\", %sHandler.TreePublic)",
			names.Plural, names.Camel)
	}

	relatedRoute := ""
	if hasParent(g.Definition) {
		relatedRoute = fmt.Sprintf("\t\tpublicAPI.GET(\"/%s/:key/related\", %sHandler.RelatedPublic)",
			names.Plural, names.Camel)
	}

	// Already mounted. Not a no-op though: a project generated before the
	// related endpoint existed has the first two routes and not the third, and
	// returning here would leave a handler nothing ever calls.
	if strings.Contains(content, names.Camel+"Handler.ListPublic") {
		// A project generated before the tree endpoint existed has the first two
		// routes and not this one, so returning early here would leave a handler
		// nothing ever calls.
		if treeRoute != "" && !strings.Contains(content, names.Camel+"Handler.TreePublic") {
			if err := injectBefore(path, "// grit:routes:public", treeRoute); err == nil {
				manifest.Refresh(path)
				fmt.Printf("  ✓ GET /api/v1/public/%s/tree (API key required)\n", names.Plural)
			}
		}
		if relatedRoute == "" || strings.Contains(content, names.Camel+"Handler.RelatedPublic") {
			return nil
		}
		if err := injectBefore(path, "// grit:routes:public", relatedRoute); err != nil {
			fmt.Printf("  Could not mount the related route: %v\n", err)
			return nil
		}
		manifest.Refresh(path)
		fmt.Printf("  ✓ GET /api/v1/public/%s/:key/related (API key required)\n", names.Plural)
		return nil
	}

	route := fmt.Sprintf(
		"\t\tpublicAPI.GET(\"/%s\", %sHandler.ListPublic)\n"+
			"\t\tpublicAPI.GET(\"/%s/:key\", %sHandler.GetPublic)",
		names.Plural, names.Camel, names.Plural, names.Camel)
	if relatedRoute != "" {
		route += "\n" + relatedRoute
	}
	// First in the block. Gin routes a static segment ahead of a parameter
	// whatever the order, but reading /tree before /:key makes that obvious to
	// the next person, who will otherwise wonder why it works.
	if treeRoute != "" {
		route = treeRoute + "\n" + route
	}

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
	if relatedRoute != "" {
		fmt.Printf("  ✓ GET /api/v1/public/%s/:key/related (similar items)\n", names.Plural)
	}
	return nil
}

// detailPayload is the expression GetPublic hands back.
//
// A tree resource resolves its subtree ids on the way out. That is the one
// question a storefront cannot answer for itself: it holds a slug, the products
// endpoint filters on ids, and walking the tree client-side would be a request
// per level.
func detailPayload(names Names, def *ResourceDefinition) string {
	if !def.Tree {
		return "toPublic" + names.Pascal + "(item)"
	}
	return "out"
}

// treeGetPublicPrelude is the code GetPublic runs before answering, for a tree.
func treeGetPublicPrelude(names Names, def *ResourceDefinition) string {
	if !def.Tree {
		return ""
	}
	return `
	// This node and everything under it, which is what "products in Electronics"
	// means to a customer. One indexed LIKE on the materialized path.
	out := toPublic` + names.Pascal + `(item)
	var descendantIDs []string
	h.DB.Model(&models.` + names.Pascal + `{}).
		Where("path LIKE ? AND archived_at IS NULL", item.Path+"%").
		Pluck("id", &descendantIDs)
	out.DescendantIDs = descendantIDs
`
}

// treePublic emits GET /api/v1/public/<plural>/tree, or nothing.
//
// A storefront's category menu is a tree, and assembling it from a flat list on
// the client means either a request per level or a client that knows how paths
// work. Neither belongs in a template somebody copies.
func treePublic(names Names, def *ResourceDefinition) string {
	if !def.Tree {
		return ""
	}
	return `
// derefID renders a nullable foreign key for the wire.
//
// The parent column is a *string because a root has no parent and SQL spells
// absent as NULL, which a foreign key constraint accepts and an empty string
// does not. JSON has the same choice to make, and "" reads better to a client
// than null for an id it is only going to compare.
func deref` + names.Pascal + `ID(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}

// publicParentOf resolves a row's parent in the node map, treating NULL and ""
// alike as no parent. Both turn up in real data: NULL from a root written now,
// "" from a row written before the column was nullable.
func public` + names.Pascal + `ParentOf(byID map[string]*public` + names.Pascal + `Node, parentID *string) (*public` + names.Pascal + `Node, bool) {
	if parentID == nil || *parentID == "" {
		return nil, false
	}
	node, ok := byID[*parentID]
	return node, ok
}

// public` + names.Pascal + `Node is a published node with its children attached.
//
// The embedded struct flattens in JSON, so a node reads as its own fields plus
// "children" rather than nesting one level deeper for no reason.
type public` + names.Pascal + `Node struct {
	public` + names.Pascal + `
	Children []*public` + names.Pascal + `Node ` + "`" + `json:"children"` + "`" + `
}

// TreePublic handles GET /api/v1/public/` + names.Plural + `/tree.
//
// The whole published hierarchy in ONE query, assembled in Go. Ordered by depth
// so a parent is always in the map before its children arrive, then by position
// and name so the order matches what the admin arranged.
func (h *` + names.Pascal + `Handler) TreePublic(c *gin.Context) {
	var rows []models.` + names.Pascal + `
	err := h.DB.Where("archived_at IS NULL").
		Order("depth asc, position asc, name asc").
		Find(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": "Failed to load the tree",
		}})
		return
	}

	byID := make(map[string]*public` + names.Pascal + `Node, len(rows))
	roots := make([]*public` + names.Pascal + `Node, 0)
	for i := range rows {
		byID[rows[i].ID] = &public` + names.Pascal + `Node{public` + names.Pascal + `: toPublic` + names.Pascal + `(rows[i])}
	}
	for i := range rows {
		node := byID[rows[i].ID]
		// A node whose parent is archived is promoted to a root rather than
		// dropped: a menu missing a whole branch is worse than one with a branch
		// in the wrong place.
		if parent, ok := public` + names.Pascal + `ParentOf(byID, rows[i].ParentID); ok {
			parent.Children = append(parent.Children, node)
			continue
		}
		roots = append(roots, node)
	}
	c.JSON(http.StatusOK, gin.H{"data": roots})
}
`
}
