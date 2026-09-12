package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// writeGoModel creates the GORM model file for the resource.
func (g *Generator) writeGoModel(names Names) error {
	fields := g.Definition.Fields

	// Detect slug field and resolve source
	var slugField *Field
	for i, f := range fields {
		if f.IsSlug() {
			slugField = &fields[i]
			break
		}
	}

	// Resolve slug source field
	slugSourceGo := ""
	if slugField != nil {
		if slugField.SlugSource != "" {
			slugSourceGo = toPascalCase(slugField.SlugSource)
		} else {
			// Default to first string field
			for _, f := range fields {
				if FieldType(f.Type) == FieldString {
					slugSourceGo = toPascalCase(f.Name)
					break
				}
			}
			if slugSourceGo == "" {
				slugSourceGo = "ID" // fallback
			}
		}
	}

	// Check if any field needs datatypes import
	needsDatatypes := false
	needsFiles := false
	for _, f := range fields {
		if f.NeedsDatatypesImport() {
			needsDatatypes = true
		}
		if f.NeedsFilesImport() {
			needsFiles = true
		}
	}

	// Auto-number fields (name:string:auto) get a sequence.Next call in
	// BeforeCreate and pull in the internal/sequence package.
	var autoFields []Field
	for _, f := range fields {
		if f.IsAuto() {
			autoFields = append(autoFields, f)
		}
	}
	hasAuto := len(autoFields) > 0

	// Build imports
	hasSlug := slugField != nil
	isTree := g.Definition.Tree
	var imports string
	// A tree needs fmt for wrapped errors and strings for path work, and a slug
	// already needs fmt. Built as a sorted list rather than by appending to a
	// string, because "fmt" has to come before "strings" before "time" and
	// nothing downstream reorders them.
	// A date field no longer contributes a "time" import: its Go type comes
	// from internal/jsontime. time is still needed unconditionally for the
	// auto-managed CreatedAt / UpdatedAt columns.
	needsJSONTime := false
	for _, f := range fields {
		if f.NeedsJSONTimeImport() {
			needsJSONTime = true
			break
		}
	}
	std := []string{`"time"`}
	if hasSlug || isTree {
		std = append(std, `"fmt"`)
	}
	if isTree {
		std = append(std, `"strings"`)
	}
	sort.Strings(std)
	stdImports := strings.Join(std, "\n\t")
	// uuid is no longer imported here: primary keys come from internal/ids
	// (UUIDv7), so GORM is the only external dependency a model needs.
	extImports := "\"gorm.io/gorm\""
	if needsDatatypes {
		extImports = "\"gorm.io/datatypes\"\n\t\"gorm.io/gorm\""
	}
	// Project imports (same module). internal/ids is unconditional — every model
	// mints its primary key in BeforeCreate. Appended in alphabetical order
	// (files, ids, sequence) so the block is already gofmt-sorted.
	needsMoney := false
	for _, f := range fields {
		if FieldType(f.Type) == FieldMoney {
			needsMoney = true
			break
		}
	}

	projImports := []string{}
	if g.erasesWithOwner() {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/erasure\"", g.Module))
	}
	for _, f := range fields {
		if f.Encrypted {
			projImports = append(projImports, fmt.Sprintf("\"%s/internal/crypto\"", g.Module))
			break
		}
	}
	if needsFiles {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/files\"", g.Module))
	}
	projImports = append(projImports, fmt.Sprintf("\"%s/internal/ids\"", g.Module))
	// jsontime carries date and datetime fields. Sorted after ids and before
	// sequence so the block stays gofmt-clean.
	if needsJSONTime {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/jsontime\"", g.Module))
	}
	if needsMoney {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/money\"", g.Module))
	}
	if hasAuto {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/sequence\"", g.Module))
	}
	// tenant sorts after sequence, so the block stays gofmt-clean.
	if g.Definition.TenantOwned {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/tenant\"", g.Module))
	}
	projImports = g.appendOnlyImports(projImports)
	if len(projImports) > 0 {
		imports = fmt.Sprintf("import (\n\t%s\n\n\t%s\n\n\t%s\n)", stdImports, extImports, strings.Join(projImports, "\n\t"))
	} else {
		imports = fmt.Sprintf("import (\n\t%s\n\n\t%s\n)", stdImports, extImports)
	}

	// Build the auto-number statements for BeforeCreate (shared by both the slug
	// and non-slug hook variants). sequence.Next is called directly — the
	// services wrapper would be a models→services import cycle.
	autoHook := ""
	for _, f := range autoFields {
		goName := toPascalCase(f.Name)
		prefix := f.AutoPrefix
		if prefix == "" {
			prefix = defaultPrefix(names.Pascal)
		}
		autoHook += fmt.Sprintf(`	if m.%s == "" {
		n, err := sequence.Next(tx, sequence.Config{
			Name:   %q,
			Prefix: %q,
			Reset:  sequence.ResetMonthly,
			Width:  4,
		}, time.Now())
		if err != nil {
			return err
		}
		m.%s = n
	}
`, goName, names.Snake+"_"+toSnakeCase(f.Name), prefix, goName)
	}

	// A workflow's initial state, applied in the same hook.
	//
	// The declaration says the record starts in draft; without this it starts
	// in the empty string, which is not one of the states, so the first
	// transition is refused and the record is unusable from birth. A GORM
	// column default would not do it either: the field is present and empty in
	// the INSERT, so the default never applies.
	if wf := g.Definition.WorkflowField(); wf != nil && wf.Workflow.Initial != "" {
		autoHook += fmt.Sprintf(`	if m.%s == "" {
		m.%s = %q
	}
`, toPascalCase(wf.Name), toPascalCase(wf.Name), wf.Workflow.Initial)
	}

	// A tree's path is derived from its parent's, and it ends in this row's own
	// id, so it cannot be computed until the id is assigned. That happens at the
	// top of BeforeCreate, a few lines above where this lands, which is the only
	// reason the tree work is not a hook of its own.
	if g.Definition.Tree {
		autoHook += treeCreateHook()
	}

	structFields := ""
	for _, f := range fields {
		// belongs_to: emit FK column + association struct
		if f.IsBelongsTo() {
			relModel := f.RelatedModelName()
			baseName := strings.TrimSuffix(f.Name, "_id") // strip _id if user included it
			fkGoName := toPascalCase(baseName) + "ID"     // e.g., CategoryID
			fkJson := toSnakeCase(baseName) + "_id"       // e.g., category_id
			assocName := toPascalCase(baseName)           // e.g., Category

			// A relation pointing at its own model is how you spell a tree:
			// Category with parent:belongs_to:Category is Electronics above
			// Cameras. Two things about it differ from an ordinary relation.
			//
			// The association must be a pointer, because Go rejects a struct
			// containing itself by value outright: "invalid recursive type".
			// A generated project with a self-reference did not compile at all
			// before this.
			//
			// And the FK must be a NULLABLE pointer, not a plain string. GORM
			// creates a real foreign key constraint for the association, and a
			// root has no parent, so the column has to hold something the
			// constraint accepts. SQL has exactly one such value and it is NULL.
			//
			// An earlier version used "" for absent, matching what other
			// belongs_to fields do. It passed every test and failed on the
			// first real project, because the tests ran on SQLite, which does
			// not enforce foreign keys unless asked, and Postgres does:
			//
			//   ERROR: insert or update on table "categories" violates foreign
			//   key constraint "fk_categories_children" (SQLSTATE 23503)
			//
			// The generated tests now switch foreign keys on so SQLite cannot
			// hide this again.
			selfRef := relModel == toPascalCase(g.Definition.Name)

			// FK column.
			//
			// one_to_one differs from belongs_to here and nowhere else: a unique
			// index. Without it the name is a comment, because the database would
			// accept a second row pointing at the same parent.
			fkIndex := "index"
			if f.IsOneToOne() {
				fkIndex = "uniqueIndex"
			}
			if selfRef {
				structFields += fmt.Sprintf("\t%s *string `gorm:\"size:36;%s\" json:\"%s\"`\n", fkGoName, fkIndex, fkJson)
			} else {
				structFields += fmt.Sprintf("\t%s string `gorm:\"size:36;%s\" json:\"%s\" binding:\"required\"`\n", fkGoName, fkIndex, fkJson)
			}
			// Association struct.
			//
			// A pointer with omitempty, so a relation that was not preloaded is
			// absent from the JSON rather than present and blank. As a value it
			// marshalled a complete zero-value object every time, which reads to
			// a client as a real record whose every field happens to be empty:
			// sender.first_name "" instead of undefined, sender.active false
			// instead of unknown. The generated TypeScript said the object was
			// always there, and agreed with the lie.
			structFields += fmt.Sprintf("\t%s *%s `gorm:\"foreignKey:%s\" json:\"%s,omitempty\"`\n",
				assocName, relModel, fkGoName, toSnakeCase(assocName))
			continue
		}

		// many_to_many: emit association slice only
		if f.IsManyToMany() {
			relModel := f.RelatedModelName()
			assocName := toPascalCase(f.Name) // e.g., Tags
			junctionTable := names.Snake + "_" + toSnakeCase(f.Name)
			structFields += fmt.Sprintf("\t%s []%s `gorm:\"many2many:%s\" json:\"%s\"`\n",
				assocName, relModel, junctionTable, toSnakeCase(f.Name))
			continue
		}

		goName := toPascalCase(f.Name)
		goType := f.GoType()
		jsonTag := toSnakeCase(f.Name)

		tags := fmt.Sprintf(`json:"%s"`, jsonTag)

		gormTag := f.GORMTag()
		if FieldType(f.Type) == FieldMoney {
			// Embedded with a prefix, so one declared field becomes
			// <name>_amount and <name>_currency. The amount stays an integer
			// column the database can SUM and index.
			gormTag = "embedded;embeddedPrefix:" + toSnakeCase(f.Name) + "_"
		}
		if gormTag != "" {
			tags = fmt.Sprintf(`gorm:"%s" %s`, gormTag, tags)
		}

		if f.Required && (f.GoType() == "string") && !f.IsSlug() {
			tags += ` binding:"required"`
		}

		structFields += fmt.Sprintf("\t%s %s `%s`\n", goName, goType, tags)
	}

	// Inline has-many child (from --items): the parent owns a slice of children
	// keyed on <Parent>ID. GORM creates these in the same transaction as the
	// parent when they're set before Create, giving atomic invoice+items saves.
	if g.Definition.Items != nil {
		childNames := BuildNames(g.Definition.Items)
		structFields += fmt.Sprintf("\tItems []%s `gorm:\"foreignKey:%sID\" json:\"items\"`\n", childNames.Pascal, names.Pascal)
	}

	// --tree: the hierarchy columns. The parent FK itself is an ordinary
	// self-referential belongs_to, emitted above with the other fields.
	if isTree {
		structFields += treeFields(names)
	}

	// --tenant-owned: OrgID and its index, plus the opt-in that makes the
	// scoping callbacks apply to this model at all. Prepended so it sits
	// directly under ID, where a reader looking for "which tenant owns this"
	// will find it.
	if g.Definition.TenantOwned {
		structFields = "\ttenant.Owned\n" + structFields
	}

	content := fmt.Sprintf(`package models

%s

// %s represents a %s in the system.
type %s struct {
	ID        string         `+"`"+`gorm:"primarykey;size:36" json:"id"`+"`"+`
%s	Version   int            `+"`"+`gorm:"not null;default:1" json:"version"`+"`"+`
	CreatedAt time.Time      `+"`"+`json:"created_at"`+"`"+`
	UpdatedAt time.Time      `+"`"+`json:"updated_at"`+"`"+`
	DeletedAt gorm.DeletedAt `+"`"+`gorm:"index" json:"-"`+"`"+`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time  `+"`"+`gorm:"index" json:"archived_at,omitempty"`+"`"+`
}
`, imports, names.Pascal, names.Lower, names.Pascal, structFields)

	// Add BeforeCreate hook (UUID generation + optional slug)
	if hasSlug {
		slugGoName := toPascalCase(slugField.Name)
		content += fmt.Sprintf(`
// BeforeCreate generates a UUID and auto-generates the slug before inserting.
func (m *%s) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.%s == "" {
		m.%s = slugify(fmt.Sprintf("%%v", m.%s))
	}
%s	return nil
}
`, names.Pascal, slugGoName, slugGoName, slugSourceGo, autoHook)

		// Write shared slugify helper if it doesn't exist yet
		helpersPath := filepath.Join(g.APIRoot(), "internal", "models", "helpers.go")
		if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
			helpersContent := `package models

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// slugify generates a URL-friendly slug with a unique suffix.
func slugify(s string) string {
	slug := strings.ToLower(s)
	re := regexp.MustCompile(` + "`" + `[^a-z0-9]+` + "`" + `)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	b := make([]byte, 4)
	rand.Read(b)
	return slug + "-" + hex.EncodeToString(b)
}
`
			if err := writeFileWithDirs(helpersPath, helpersContent); err != nil {
				return fmt.Errorf("writing helpers.go: %w", err)
			}
		}
	} else {
		// No slug — still need UUID generation (+ any auto-number fields)
		content += fmt.Sprintf(`
// BeforeCreate generates a UUID before inserting.
func (m *%s) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
%s	return nil
}
`, names.Pascal, autoHook)
	}

	// BeforeUpdate increments Version on every server-side write so offline
	// clients can detect that a record they edited has moved on. Pair with
	// /api/sync/push for safe write replay.
	content += fmt.Sprintf(`
// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *%s) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
`, names.Pascal)
	content += g.appendOnlyModelInit(names)

	// GetOwnerID makes the model satisfy authz.Ownable, which is what lets the
	// handlers check ownership on a row fetched by id. Without it the helper
	// has nothing to compare against, which is why it sat unused for so long.
	if owner := g.Definition.OwnerField(); owner != nil {
		content += fmt.Sprintf(`
// GetOwnerID identifies the row's owner for authz.MustOwnUnlessAdmin.
//
// This resource was generated with --owned-by %s: a caller may only read or
// write rows where this matches their own user id, ADMIN excepted.
func (m *%s) GetOwnerID() string {
	return m.%s
}
`, owner.Name, names.Pascal, toPascalCase(owner.Name)+"ID")
	}
	content += g.erasureModelInit(names)

	// Path resolution, the cycle refusal, and the two accessors a breadcrumb
	// needs. Reparenting lives in the tree service instead: a move rewrites a
	// whole subtree, and that belongs somewhere it can be a transaction.
	if isTree {
		content += treeModelMethods(names)
	}

	path := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// buildServiceSearchWhere creates the full Where(...) arguments for the service search.
func (g *Generator) buildServiceSearchWhere() string {
	var searchFields []string
	for _, f := range g.Definition.Fields {
		// Relationships are excluded, and the same rule the handler uses picks
		// the rest, so the two agree about what "search" means.
		//
		// GoType() is "string" for a belongs_to as well, because the foreign
		// key is a UUID string. Taking that at face value put LOWER(parent) in
		// the WHERE clause of every tree resource: parent is the relation, the
		// column is parent_id, and the query fails at the database. It went
		// unnoticed because nothing calls the generated service until somebody
		// writes the first line of business logic in it.
		if f.IsRelationship() {
			continue
		}
		if f.IsSearchable() {
			searchFields = append(searchFields, "LOWER("+toSnakeCase(f.Name)+") LIKE LOWER(?)")
		}
	}
	if len(searchFields) == 0 {
		// CAST(... AS TEXT) rather than ::text — the latter is Postgres-only.
		searchFields = []string{"LOWER(CAST(id AS TEXT)) LIKE LOWER(?)"}
	}

	clause := strings.Join(searchFields, " OR ")
	args := ""
	for range searchFields {
		args += `, "%"+params.Search+"%"`
	}

	return `"` + clause + `"` + args
}

// buildSortableSet returns a Go map literal of the columns the generated
// service will permit in ORDER BY. It's the whitelist that stops a
// client-supplied sort_by from being interpolated as raw SQL. Always includes
// the system timestamp/id columns, plus every scalar (non-relation, non-file,
// non-array) field's DB column, plus each belongs_to's foreign-key column.
func (g *Generator) buildSortableSet() string {
	cols := []string{"id", "created_at", "updated_at"}
	seen := map[string]bool{"id": true, "created_at": true, "updated_at": true}
	add := func(c string) {
		if c != "" && !seen[c] {
			seen[c] = true
			cols = append(cols, c)
		}
	}
	for _, f := range g.Definition.Fields {
		switch {
		case f.Encrypted, f.IsFile(), f.IsFiles(), f.IsManyToMany(), f.IsStringArray():
			// not a sortable scalar column (ciphertext sorts as noise)
			continue
		case f.IsBelongsTo():
			base := strings.TrimSuffix(toSnakeCase(f.Name), "_id")
			add(base + "_id")
		default:
			add(toSnakeCase(f.Name))
		}
	}
	var b strings.Builder
	b.WriteString("map[string]bool{")
	for i, c := range cols {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%q: true", c))
	}
	b.WriteString("}")
	return b.String()
}

// workflowHandlerMethod returns the Transition handler, or "" for a resource
// with no workflow.
//
// It is thin on purpose: parse, delegate, translate the error. The guard lives
// in the service, because a handler is only one of the ways in.
func (g *Generator) workflowHandlerMethod(names Names) string {
	if g.Definition.WorkflowField() == nil {
		return ""
	}
	return `

// Transition moves a ` + names.Lower + ` through its workflow.
//
// POST /api/` + names.Plural + `/:id/transitions/:action
func (h *` + names.Pascal + `Handler) Transition(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")

	// The same grants RequireRole reads, checked per transition rather than
	// per route: which permission applies depends on which move is being
	// made, and a route can only know one.
	can := func(perm string) bool {
		grants, ok := c.Get("user_grants")
		if !ok {
			return false
		}
		list, ok := grants.([]string)
		if !ok {
			return false
		}
		return authz.Granted(list, perm)
	}

	item, err := services.Transition` + names.Pascal + `(h.DB.WithContext(c.Request.Context()), c, id, action, can)
	if err != nil {
		// An illegal move or one a hook refused is a 422 saying why, a missing
		// permission a 403. Anything else is a 500 whose detail stays in the
		// log: it used to be answered as a 403, so a database failure read as
		// "you are not allowed".
		status, code := workflow.Classify(err)
		message := err.Error()
		switch status {
		case http.StatusNotFound:
			message = "` + names.Pascal + ` not found"
		case http.StatusInternalServerError:
			log.Printf("transition %s on ` + names.Lower + ` %s: %v", action, id, err)
			message = "The transition could not be completed"
		}
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "` + names.Pascal + ` " + action,
	})
}

// Workflow returns the state machine, so a client can render badges and the
// moves legal from where a record currently is.
//
// GET /api/` + names.Plural + `/workflow
func (h *` + names.Pascal + `Handler) Workflow(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": workflow.` + names.Pascal + `Workflow})
}
`
}

// v3.31.39: pickIdentifierExpr returns the Go expression to use as
// the human-readable identifier in CUD activity log lines. Picks the
// first match from Name / Title / Slug / SKU / Subject / Label /
// Email on the model. Falls back to item.ID when none of those exist
// so the log line is never blank ({verb} {entityType} {identifier}
// is the format convention; identifier is never empty).
func pickIdentifierExpr(fields []Field) string {
	// Number / Reference / Code come after the human-readable names but before
	// the ID fallback, so an auto-numbered record (an invoice, an order) is
	// identified by its number rather than an opaque UUID.
	candidates := []string{"Name", "Title", "Slug", "Sku", "SKU", "Subject", "Label", "Email", "Number", "Reference", "Code"}
	available := map[string]bool{}
	for _, f := range fields {
		if f.IsBelongsTo() || f.IsManyToMany() {
			continue
		}
		available[toPascalCase(f.Name)] = true
	}
	for _, c := range candidates {
		if available[c] {
			return "item." + c
		}
	}
	return "item.ID"
}

// buildHandlerSearchCols returns the comma-separated quoted column names for
// the paginate.Config.Searchable slice literal. Only text-like field types are
// included — FK UUID columns (which happen to be Go string) are skipped so
// search doesn't match against opaque identifiers.
func (g *Generator) buildHandlerSearchCols() string {
	var cols []string
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() {
			continue
		}
		if f.IsSearchable() {
			cols = append(cols, `"`+toSnakeCase(f.Name)+`"`)
		}
	}
	return strings.Join(cols, ", ")
}

// writeZodSchema creates the Zod schema file for the resource.
func (g *Generator) writeZodSchema(names Names) error {
	createFields := ""
	updateFields := ""

	for _, f := range g.Definition.Fields {
		// Slug fields are auto-generated — exclude from create/update schemas
		if f.IsSlug() {
			continue
		}

		// belongs_to: use FK column name (e.g., category_id)
		if f.IsBelongsTo() {
			// Zod schemas use snake_case to match the Go handlers' JSON tags.
			// The Go API decodes with ShouldBindJSON using `json:"foo_id"` tags,
			// so the Zod payload must send snake_case keys too.
			fkName := f.FKColumnName()
			zod := f.ZodType()
			// A self-reference has to accept the empty string, because that is
			// how "this is a root" is spelled. Without it the admin form refuses
			// to save the first node: .uuid() rejects "" with "Invalid ID" on a
			// field the person correctly left blank.
			if f.RelatedModelName() == toPascalCase(g.Definition.Name) {
				zod += `.or(z.literal(""))`
			}
			createFields += fmt.Sprintf("  %s: %s,\n", fkName, zod)
			updateFields += fmt.Sprintf("  %s: %s,\n", fkName, zod+".optional()")
			continue
		}

		// many_to_many: use <name>_ids (e.g., tag_ids)
		if f.IsManyToMany() {
			idsName := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			createFields += fmt.Sprintf("  %s: %s,\n", idsName, f.ZodType())
			updateFields += fmt.Sprintf("  %s: %s,\n", idsName, f.ZodType())
			continue
		}

		snakeName := toSnakeCase(f.Name)
		zodType := f.ZodType()
		createFields += fmt.Sprintf("  %s: %s,\n", snakeName, zodType)

		// Update schema: make all fields optional
		updateZod := f.ZodType()
		if !strings.Contains(updateZod, ".optional()") && !strings.Contains(updateZod, ".nullable()") {
			updateZod += ".optional()"
		}
		updateFields += fmt.Sprintf("  %s: %s,\n", snakeName, updateZod)
	}

	// Shared schemas the generated one references directly rather than
	// inlining, so each shape has one definition.
	//
	// Built as a list. Appending to a slice stays correct the day a third
	// shared schema appears; rewriting an anchor line fails silently the
	// day the anchor moves.
	needsFileRef, needsMoney := false, false
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			needsFileRef = true
		}
		if FieldType(f.Type) == FieldMoney {
			needsMoney = true
		}
	}

	imports := []string{`import { z } from "zod";`}
	if needsFileRef {
		imports = append(imports, `import { FileRefSchema } from "./file-ref";`)
	}
	if needsMoney {
		imports = append(imports, `import { MoneySchema } from "./money";`)
	}
	importLines := strings.Join(imports, "\n")

	content := fmt.Sprintf(`%s

export const Create%sSchema = z.object({
%s});

export const Update%sSchema = z.object({
%s});

export type Create%sInput = z.infer<typeof Create%sSchema>;
export type Update%sInput = z.infer<typeof Update%sSchema>;
`, importLines, names.Pascal, createFields, names.Pascal, updateFields,
		names.Pascal, names.Pascal, names.Pascal, names.Pascal)

	path := filepath.Join(g.Root, "packages", "shared", "schemas", names.Kebab+".ts")
	return writeFileWithDirs(path, content)
}

// writeTSTypes creates the TypeScript type file for the resource.
func (g *Generator) writeTSTypes(names Names) error {
	// Collect relationship imports
	imports := ""
	fields := ""
	needsFileRef, needsMoney := false, false
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			relModel := f.RelatedModelName()
			relKebab := strings.ReplaceAll(toSnakeCase(relModel), "_", "-")
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkSnake := toSnakeCase(baseName) + "_id"

			// A self-reference needs no import: the type is declared in this
			// very file. Emitting one produced
			//
			//   import type { Category } from "./category";
			//   export interface Category { ... }
			//
			// in category.ts, which is TS2440, "import declaration conflicts
			// with local declaration", and it failed tsc for the whole
			// workspace rather than only for that file.
			selfRef := relModel == toPascalCase(g.Definition.Name)
			if !selfRef {
				imports += fmt.Sprintf("import type { %s } from \"./%s\";\n", relModel, relKebab)
			}

			// FK matches the referenced model's UUID string PK, and a
			// self-referential one is nullable: a tree root has no parent, and
			// the column holds NULL rather than an empty string.
			if selfRef {
				fields += fmt.Sprintf("  %s: string | null;\n", fkSnake)
			} else {
				fields += fmt.Sprintf("  %s: string;\n", fkSnake)
			}
			fields += fmt.Sprintf("  %s?: %s;\n", toSnakeCase(baseName), relModel)
			continue
		}
		if f.IsManyToMany() {
			relModel := f.RelatedModelName()
			relKebab := strings.ReplaceAll(toSnakeCase(relModel), "_", "-")
			imports += fmt.Sprintf("import type { %s } from \"./%s\";\n", relModel, relKebab)
			fields += fmt.Sprintf("  %s?: %s[];\n", toSnakeCase(f.Name), relModel)
			continue
		}
		// v3.31.37: file/files fields reference the FileRef type from
		// schemas/file-ref.ts. Without an explicit import the generated
		// type fails tsc with "Cannot find name 'FileRef'", which
		// downstream files (the React Query hook, admin resources)
		// inherit when they import from this module.
		if f.IsFileField() {
			needsFileRef = true
		}
		// Same reason as FileRef above: the interface says Money, so the
		// name has to come from somewhere.
		if FieldType(f.Type) == FieldMoney {
			needsMoney = true
		}
		tsName := toSnakeCase(f.Name)
		tsType := f.TSType()
		fields += fmt.Sprintf("  %s: %s;\n", tsName, tsType)
	}
	if needsFileRef {
		imports += "import type { FileRef } from \"../schemas/file-ref\";\n"
	}
	if needsMoney {
		imports += "import type { Money } from \"./money\";\n"
	}

	content := ""
	if imports != "" {
		content = imports + "\n"
	}
	content += fmt.Sprintf(`export interface %s {
  id: string;
%s  created_at: string;
  updated_at: string;
}
`, names.Pascal, fields)

	path := filepath.Join(g.Root, "packages", "shared", "types", names.Kebab+".ts")
	return writeFileWithDirs(path, content)
}

// writeReactQueryHooks creates React Query hooks for the resource.
// v3.31.21: pick the right api-client path per app — apps/admin has
// lib/api-client.ts; apps/web has lib/api.ts that re-exports apiClient.
func (g *Generator) writeReactQueryHooks(names Names, app string) error {
	apiImport := `import { apiClient } from "@/lib/api-client";`
	if app == "web" {
		apiImport = `import { apiClient } from "@/lib/api";`
	}
	// v3.31.42: the hook's inline `interface <Resource>` may reference
	// FileRef when any field is :file: / :files:. Without an explicit
	// import the generated hook fails tsc with "Cannot find name
	// 'FileRef'" -- the same TS2304 the typed shared model used to
	// surface before v3.31.37 patched writeTSTypes. Patch the hook
	// generator the same way.
	// The resource type comes from the shared package rather than being
	// redeclared here.
	//
	// A local copy is what broke "one backend, every client, from shared
	// types": grit sync rewrites packages/shared and cannot reach a duplicate
	// declaration, so a field added to the Go model reached the shared type
	// and nowhere else. Nothing errored, because each file stayed internally
	// consistent; the app simply did not know the field existed. The
	// scaffold's own use-blogs.ts has always imported the type, and this now
	// matches it.
	//
	// FileRef and Money were imported only to satisfy that local copy, so they
	// are no longer needed: an unused import is an error under noUnusedLocals,
	// which the Vite admin sets.
	apiImport += "\nimport type { " + names.Pascal + " } from \"@repo/shared/types\";"
	content := fmt.Sprintf(`import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
%s

interface %sResponse {
  data: %s[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

interface Use%sParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
}

export function use%s({ page = 1, pageSize = 20, search = "", sortBy = "created_at", sortOrder = "desc" }: Use%sParams = {}) {
  return useQuery<%sResponse>({
    queryKey: ["%s", { page, pageSize, search, sortBy, sortOrder }],
    queryFn: async () => {
      const params = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
        sort_by: sortBy,
        sort_order: sortOrder,
      });
      if (search) {
        params.set("search", search);
      }
      const { data } = await apiClient.get(%s);
      return data;
    },
  });
}

export function useGet%s(id: string) {
  return useQuery<%s>({
    queryKey: ["%s", id],
    queryFn: async () => {
      const { data } = await apiClient.get(%s);
      return data.data;
    },
    enabled: !!id,
  });
}

export function useCreate%s() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => {
      const { data } = await apiClient.post("/api/%s", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["%s"] });
    },
  });
}

export function useUpdate%s() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string } & Record<string, unknown>) => {
      const { data } = await apiClient.put(%s, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["%s"] });
    },
  });
}

export function useDelete%s() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(%s);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["%s"] });
    },
  });
}
`,
		apiImport,
		names.PluralPascal, names.Pascal,
		names.PluralPascal,
		names.PluralPascal, names.PluralPascal,
		names.PluralPascal,
		names.Plural,
		"`/api/"+names.Plural+"?${params}`",
		names.Pascal, names.Pascal,
		names.Plural,
		"`/api/"+names.Plural+"/${id}`",
		names.Pascal,
		names.Plural,
		names.Plural,
		names.Pascal,
		"`/api/"+names.Plural+"/${id}`",
		names.Plural,
		names.Pascal,
		"`/api/"+names.Plural+"/${id}`",
		names.Plural,
	)

	path := filepath.Join(g.Root, "apps", app, "hooks", "use-"+names.PluralKebab+".ts")
	return writeFileWithDirs(path, content)
}

func (g *Generator) buildTSInterfaceFields() string {
	result := ""
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkSnake := toSnakeCase(baseName) + "_id"
			// FK matches the referenced model's UUID string PK.
			result += fmt.Sprintf("  %s: string;\n", fkSnake)
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			result += fmt.Sprintf("  %s?: any;\n", toSnakeCase(baseName))
			continue
		}
		if f.IsManyToMany() {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			result += fmt.Sprintf("  %s?: any[];\n", toSnakeCase(f.Name))
			continue
		}
		result += fmt.Sprintf("  %s: %s;\n", toSnakeCase(f.Name), f.TSType())
	}
	return result
}

// writeResourceDefinition creates the resource definition file (resources/<plural>.ts).
// resourceDefinitionFileContent builds the admin resource definition file.
//
// Both admins consume the SAME apps/admin/**/lib/resource.ts defineResource(),
// so both must get byte-identical content — only the destination path differs
// (Next.js: apps/admin/resources, TanStack: apps/admin/src/resources). Keeping
// one builder is deliberate: a TanStack-specific copy previously drifted and
// emitted a flat {plural, apiEndpoint, columns, fields} shape that defineResource
// could not read, blanking the whole admin at import time.
func (g *Generator) resourceDefinitionFileContent(names Names) string {
	icon := guessLucideIcon(names.Pascal)

	// v3.31.19: column-pack heuristic. When a resource has both `name`
	// and `email` (or both `first_name` and `last_name`), pack them
	// into a single stacked column instead of two narrow ones. The
	// helper component lives at apps/admin/components/tables/stacked-cell.tsx
	// and is imported automatically when a pack fires.
	packs, usesStackedCell := detectColumnPacks(g.Definition.Fields)

	// Build column definitions. ID is intentionally NOT listed by default —
	// UUIDs are noisy and rarely something an operator scans by eye.
	// Users who want it can add { key: "id", label: "ID", width: "80px" }
	// to the columns array by hand.
	columns := ""
	// The first plain (non-relationship) column becomes click-to-open — the
	// primary identifier (number / name / title) links to the detail page out
	// of the box. Developers can move it, switch it to "copy", or give any
	// column a custom onClick.
	linkedFirstColumn := false
	for _, f := range g.Definition.Fields {
		colName := toSnakeCase(f.Name)

		// Is this field part of a pack? Either emit the pack (when this
		// field is the pack's primary key) or skip silently.
		if pack, ok := packs[colName]; ok {
			if pack.primary == colName {
				columns += "\n      " + pack.line
			}
			continue
		}

		// belongs_to: show related model's name via dot notation
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			assocSnake := toSnakeCase(baseName)
			colLabel := strings.Join(splitPascal(toPascalCase(baseName)), " ")
			columns += fmt.Sprintf("\n      { key: \"%s.name\", label: \"%s\" },", assocSnake, colLabel)
			continue
		}
		// many_to_many: skip from table columns (arrays are noisy)
		if f.IsManyToMany() {
			continue
		}

		colLabel := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
		sortable := f.IsSortable()
		searchable := f.IsSearchable()
		format := f.ColumnFormat()

		parts := []string{
			fmt.Sprintf(`key: "%s"`, colName),
			fmt.Sprintf(`label: "%s"`, colLabel),
		}
		if sortable {
			parts = append(parts, `sortable: true`)
		}
		if searchable {
			parts = append(parts, `searchable: true`)
		}
		if format != "text" {
			parts = append(parts, fmt.Sprintf(`format: "%s"`, format))
		}
		if !linkedFirstColumn {
			parts = append(parts, `onClick: "link"`)
			linkedFirstColumn = true
		}

		columns += "\n      { " + strings.Join(parts, ", ") + " },"
	}
	columns += fmt.Sprintf(`
      { key: "created_at", label: "Created", sortable: true, format: "relative" },`)

	// Build form field definitions (skip slug — auto-generated, not editable)
	formFields := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		// Auto-number fields are filled by the server in BeforeCreate — keep them
		// out of the create/edit form entirely (they still show in the table/detail).
		if f.IsAuto() {
			continue
		}

		// belongs_to: relationship-select with endpoint
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkKey := toSnakeCase(baseName) + "_id"
			fieldLabel := strings.Join(splitPascal(toPascalCase(baseName)), " ")
			relSnake := toSnakeCase(f.RelatedModelName())
			relPlural := Pluralize(relSnake)
			// A self-reference is optional in the form, because the top of a
			// tree has no parent and a required select gives the person no way
			// to say so.
			required := "required: true, "
			if f.RelatedModelName() == toPascalCase(g.Definition.Name) {
				required = ""
			}
			formFields += fmt.Sprintf("\n    { key: \"%s\", label: \"%s\", type: \"relationship-select\", %srelatedEndpoint: \"/api/%s\", displayField: \"name\" },",
				fkKey, fieldLabel, required, relPlural)
			continue
		}

		// many_to_many: multi-relationship-select with endpoint
		if f.IsManyToMany() {
			idsKey := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			fieldLabel := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
			relSnake := toSnakeCase(f.RelatedModelName())
			relPlural := Pluralize(relSnake)
			assocKey := toSnakeCase(f.Name)
			formFields += fmt.Sprintf("\n    { key: \"%s\", label: \"%s\", type: \"multi-relationship-select\", relatedEndpoint: \"/api/%s\", displayField: \"name\", relationshipKey: \"%s\" },",
				idsKey, fieldLabel, relPlural, assocKey)
			continue
		}

		fieldKey := toSnakeCase(f.Name)
		fieldLabel := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
		fieldType := f.FormFieldType()

		parts := []string{
			fmt.Sprintf(`key: "%s"`, fieldKey),
			fmt.Sprintf(`label: "%s"`, fieldLabel),
			fmt.Sprintf(`type: "%s"`, fieldType),
		}
		if f.Required {
			parts = append(parts, `required: true`)
		}
		// Carried into the admin so bulk edit can skip it. A slug is unique
		// too, but it is excluded by type rather than by this flag.
		if f.Unique {
			parts = append(parts, `unique: true`)
		}
		// v3.31.30: emit file/files accepts + size knobs so the runtime
		// FileField can build the per-field upload URL without round-
		// tripping to the API for field metadata.
		if f.IsFileField() && len(f.FileAccepts) > 0 {
			quoted := make([]string, 0, len(f.FileAccepts))
			for _, a := range f.FileAccepts {
				quoted = append(quoted, fmt.Sprintf("%q", a))
			}
			parts = append(parts, fmt.Sprintf("accepts: [%s]", strings.Join(quoted, ", ")))
			// Default size cap: 300MB for video, 5MB everything else.
			// Resource def can still override by hand.
			size := 5
			for _, a := range f.FileAccepts {
				if a == "video" {
					size = 300
					break
				}
			}
			parts = append(parts, fmt.Sprintf("maxSizeMB: %d", size))
			if f.IsFiles() {
				parts = append(parts, "max: 5")
			}
		}
		// v3.31.38: numberKind hints the comma-formatting NumberField
		// at the Go-side domain. int allows negatives but no decimals,
		// uint disallows both, float allows both.
		switch FieldType(f.Type) {
		case FieldInt:
			parts = append(parts, `numberKind: "int"`)
		case FieldUint:
			parts = append(parts, `numberKind: "uint"`)
		case FieldFloat:
			parts = append(parts, `numberKind: "float"`)
		}
		// select / check carry their value=label choices so the SelectField /
		// CheckboxGroupField can render the dropdown or checkboxes.
		if f.HasOptions() {
			parts = append(parts, "options: "+f.OptionsLiteral())
		}

		formFields += "\n    { " + strings.Join(parts, ", ") + " },"
	}

	// Build filter definitions (auto-detect boolean and select-like fields)
	filters := ""
	for _, f := range g.Definition.Fields {
		if FieldType(f.Type) == FieldBool {
			filterKey := toSnakeCase(f.Name)
			filterLabel := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
			filters += fmt.Sprintf(`
    { key: "%s", label: "%s", type: "boolean" },`, filterKey, filterLabel)
		}
	}

	// Inline line-items field (parent of --items): an editable child table inside
	// the parent form. itemFields are the child's editable columns (its FK back
	// to the parent is set server-side, so it's never a column).
	if g.Definition.Items != nil {
		childNames := BuildNames(g.Definition.Items)
		itemCols := ""
		for _, cf := range g.Definition.Items.Fields {
			if cf.IsSlug() || cf.IsManyToMany() {
				continue
			}
			if cf.IsBelongsTo() && cf.RelatedModelName() == names.Pascal {
				continue // the back-link to the parent
			}
			label := strings.Join(splitPascal(toPascalCase(strings.TrimSuffix(cf.Name, "_id"))), " ")
			if cf.IsBelongsTo() {
				relKebab := replaceAll(Pluralize(toSnakeCase(cf.RelatedModelName())), "_", "-")
				itemCols += fmt.Sprintf("\n        { key: %q, label: %q, type: \"relationship-select\", relatedEndpoint: \"/api/%s\", displayField: \"name\" },", cf.FKColumnName(), label, relKebab)
				continue
			}
			// The same mapping the resource's own form uses.
			//
			// This used to be a second switch covering five types and falling
			// through to "text", so money, toggles, selects, textareas and rich
			// text were all plain text inputs. An invoice line's unit_rate came
			// out as text while the identical column in the child's own
			// resource came out as money, in the same run. A partial copy of a
			// mapping is a copy that falls behind, and this one had.
			typ := cf.FormFieldType()
			if typ == "" {
				continue // slug, and anything else with no form control
			}
			extra := ""
			switch FieldType(cf.Type) {
			case FieldInt:
				extra = `, numberKind: "int"`
			case FieldUint:
				extra = `, numberKind: "uint"`
			case FieldFloat:
				extra = `, numberKind: "float"`
			case FieldSelect, FieldRadio, FieldCheck:
				// Without its options a select renders an empty dropdown, which
				// reads as a loading bug rather than a missing argument.
				extra = ", options: " + cf.OptionsLiteral()
			}
			itemCols += fmt.Sprintf("\n        { key: %q, label: %q, type: %q%s },", toSnakeCase(cf.Name), label, typ, extra)
		}
		// itemEndpoint is the API path (snake plural, matching the child's routes),
		// NOT the kebab frontend slug. The label is the spaced plural.
		itemsLabel := strings.Join(splitPascal(childNames.PluralPascal), " ")
		formFields += fmt.Sprintf("\n    { key: \"items\", label: %q, type: \"line-items\", colSpan: 2, itemEndpoint: \"/api/%s\", foreignKey: %q, itemFields: [%s\n    ] },",
			itemsLabel, childNames.Plural, names.Snake+"_id", itemCols)
	}

	// Hidden resources (inline --items children) are generated fully but kept
	// out of the sidebar — managed via the parent's form + detail page.
	hiddenLine := ""
	if g.Definition.Hidden {
		hiddenLine = "\n  hidden: true,"
	}
	// --tree turns on the Table / Tree toggle on the list page. The flag and the
	// endpoints the tree view calls are generated together, so one can never
	// arrive without the other.
	if g.Definition.Tree {
		hiddenLine += "\n  tree: true,"
	}

	// v3.31.19: conditionally pull in the StackedCell helper. Only
	// emitted when the column-pack heuristic actually fires — keeps
	// resources without a pack from carrying a dead import.
	stackedCellImport := ""
	if usesStackedCell {
		stackedCellImport = "\nimport { StackedCell } from \"@/components/tables/stacked-cell\";"
	}

	content := fmt.Sprintf(`import { defineResource } from "@/lib/resource";%s
import custom from "./%s.custom";

export const %sResource = defineResource({
  name: "%s",
  slug: "%s",%s
  endpoint: "/api/%s",
  icon: "%s",
  label: { singular: "%s", plural: "%s" },
  table: {
    columns: [
      // grit:cols:auto-start%s
      // grit:cols:auto-end
    ],
    filters: [%s
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    searchable: true,
    pageSize: 20,
    // Shown once rows are ticked. Drop "archive" here and the Archived tab
    // goes with it; the model keeps its archived_at either way.
    bulkActions: ["edit", "archive", "restore", "export", "delete"],
  },
  form: {
    fields: [
      // grit:fields:auto-start%s
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total %s",
        endpoint: "/api/%s",
        icon: "%s",
        color: "accent",
      },
    ],
  },
}, custom);
`,
		stackedCellImport,
		names.PluralKebab,
		names.Camel,
		names.Pascal,
		names.PluralKebab,
		hiddenLine,
		names.Plural,
		icon,
		// Spaced, because this is what the sidebar, the page title and every
		// form heading show: "Inventory Items", not "InventoryItems".
		strings.Join(splitPascal(names.Pascal), " "), strings.Join(splitPascal(names.PluralPascal), " "),
		columns,
		filters,
		formFields,
		names.PluralPascal,
		names.Plural,
		icon,
	)

	return content
}

// writeResourceDefinition writes the resource definition for the Next.js admin.
func (g *Generator) writeResourceDefinition(names Names) error {
	root := filepath.Join(g.AdminCodeRoot(), "resources")
	// One folder per resource. If this project is still flat, the definition
	// that was there is left where it is and the migration in grit upgrade
	// moves it: silently writing a second copy in a folder would leave two
	// definitions and a registry pointing at the stale one.
	if flat := filepath.Join(root, names.PluralKebab+".ts"); fileExists(flat) {
		if err := g.writeAdminFile(flat, g.adminDefinitionContent(names)); err != nil {
			return err
		}
		return writeResourceCustomStub(root, names)
	}
	dir, path := scaffold.ResourceDefPath(root, names.PluralKebab)
	if err := writeFileWithDirs(path, g.adminDefinitionContent(names)); err != nil {
		return err
	}
	return writeResourceCustomStub(dir, names)
}

// writeResourceCustomStub drops resources/<slug>.custom.tsx next to the
// generated definition, and refuses to touch it if it is already there.
//
// That refusal is the whole point. The .ts half is rewritten on every generate;
// this half is where components live precisely because it never is. Overwriting
// it would delete the custom table someone spent an afternoon on, silently, as
// a side effect of adding a field to an unrelated model.
func writeResourceCustomStub(dir string, names Names) error {
	path := filepath.Join(dir, names.PluralKebab+".custom.tsx")
	if _, err := os.Stat(path); err == nil {
		return nil // already exists — yours, not ours
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking %s: %w", path, err)
	}
	// names.Pascal is also the shared row type: grit emits
	// packages/shared/types/<name>.ts with an interface of the same name.
	return writeFileWithDirs(path, scaffold.AdminResourceCustomStub(names.Pascal, names.Pascal))
}

// writeResourcePage creates a thin admin page wrapper for the resource.
func (g *Generator) writeResourcePage(names Names) error {
	content := fmt.Sprintf(`"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { %sResource } from "@/resources/%s/%s";

export default function %sPage() {
  return <ResourcePage resource={%sResource} />;
}
`,
		names.Camel, names.PluralKebab, names.PluralKebab,
		names.PluralPascal,
		names.Camel,
	)

	path := filepath.Join(g.AdminRoutesRoot(), "(dashboard)", "resources", names.PluralKebab, "page.tsx")
	return g.writeAdminFile(path, content)
}

// writeResourceDetailPage writes the per-resource [id] detail route (Next.js).
// Every "view" action navigates here.
func (g *Generator) writeResourceDetailPage(names Names) error {
	content := fmt.Sprintf(`"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { %sResource } from "@/resources/%s/%s";

export default function %sDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={%sResource} id={id} />;
}
`,
		names.Camel, names.PluralKebab, names.PluralKebab,
		names.PluralPascal,
		names.Camel,
	)

	path := filepath.Join(g.AdminRoutesRoot(), "(dashboard)", "resources", names.PluralKebab, "[id]", "page.tsx")
	return g.writeAdminFile(path, content)
}

// toCamelCase converts snake_case to camelCase.
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	result := parts[0]
	for _, p := range parts[1:] {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

// splitPascal splits PascalCase into words: "AuthorId" -> ["Author", "Id"]
// splitPascal breaks a PascalCase identifier into display words.
//
// Runs of capitals stay together, so PortfolioURL reads "Portfolio URL" rather
// than "Portfolio U R L", and APIKey reads "API Key". A capital that starts a
// new lowercase run still opens a word, which is what keeps the trailing
// letter of a run attached to the word it belongs to (PDFExport → PDF Export,
// not PDFE xport).
func splitPascal(s string) []string {
	var words []string
	start := 0
	for i := 1; i < len(s); i++ {
		isUpper := s[i] >= 'A' && s[i] <= 'Z'
		if !isUpper {
			continue
		}
		prevUpper := s[i-1] >= 'A' && s[i-1] <= 'Z'
		nextLower := i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z'
		// Mid-acronym: only break when this capital begins a new word.
		if prevUpper && !nextLower {
			continue
		}
		words = append(words, s[start:i])
		start = i
	}
	words = append(words, s[start:])
	return words
}

// dedentOneTab strips exactly one leading tab from every line.
//
// The create/update field lists are built for a struct declared *inside* a
// handler (two tabs). Promoting them to a package-level type moves them out one
// level, and nested blocks — the inline `Items []struct{...}` from --items —
// have to move with them, which is why this shifts each line by one tab rather
// than collapsing all leading whitespace.
func dedentOneTab(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, "\t")
	}
	return strings.Join(lines, "\n")
}
