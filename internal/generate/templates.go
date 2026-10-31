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
// relOptional reports whether a relationship may be empty.
//
// A self-reference always may: the root of a tree has no parent. Every other
// belongs_to follows its own required flag, exactly like every other field
// type does, which it did not until grit#100. The flag was ignored and every
// foreign key came out `binding:"required"`, so a deal with no contact yet
// could not be saved at all, and the workaround (store "") hits
// SQLSTATE 23503 the moment the foreign key constraint exists.
//
// One predicate rather than four copies of the condition, because the model,
// the handler, the TypeScript type and the admin form each decided this
// separately before, and three of them said required while the fourth said
// nullable.
func (g *Generator) relOptional(f Field) bool {
	if f.RelatedModelName() == toPascalCase(g.Definition.Name) {
		return true
	}
	// The owner of an --owned-by resource is never absent. It is stamped from
	// the signed-in caller rather than accepted from the body, every query is
	// scoped by it, and the generated OwnerID() returns it as a string. A
	// nullable owner is a row nobody can see and a model that does not
	// compile, which is what CI caught the first time this change shipped.
	if owner := g.Definition.OwnerField(); owner != nil && owner.Name == f.Name {
		return false
	}
	return !f.Required
}

func (g *Generator) writeGoModel(names Names) error {
	fields := g.Definition.Fields
	owner := g.Definition.OwnerField()

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
	// A tenant-owned auto-number refuses the write when there is no
	// organization on the context, and that refusal is an fmt.Errorf.
	if hasSlug || isTree || (hasAuto && g.Definition.TenantOwned) {
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
	// A tel column is checked by internal/phone, which registers itself with
	// internal/fieldtypes when it is imported. Without the import the format
	// tag names a check nothing provides, and fieldtypes refuses the write.
	for _, f := range fields {
		if f.IsTel() {
			projImports = append(projImports, fmt.Sprintf("_ \"%s/internal/phone\"", g.Module))
			break
		}
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
		counter := names.Snake + "_" + toSnakeCase(f.Name)
		if g.Definition.TenantOwned {
			// One counter per organization, not one for the platform.
			//
			// sequence.Counter is keyed by (name, bucket), so a fixed name
			// means every organization draws from the same counter: A gets
			// 0001, B gets 0002, and A's next is 0003. A sees a gap it cannot
			// explain, and each organization can work out how many numbers the
			// others issue. Gap-free numbering per business is a legal
			// requirement in many countries, which is what auto is for.
			//
			// It fails closed. With no organization on the context the write is
			// refused rather than drawing from a shared counter, because a
			// number from the wrong counter is worse than a write that did not
			// happen. See grit#131.
			autoHook += fmt.Sprintf(`	if m.%s == "" {
		orgID, ok := tenant.FromContext(tx.Statement.Context)
		if !ok {
			return fmt.Errorf("%s: no organization on the context, so there is no counter to number from")
		}
		n, err := sequence.Next(tx, sequence.Config{
			Name:   %q + ":" + orgID,
			Prefix: %q,
			Reset:  sequence.ResetMonthly,
			Width:  4,
		}, time.Now())
		if err != nil {
			return err
		}
		m.%s = n
	}
`, goName, counter, counter, prefix, goName)
			continue
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
`, goName, counter, prefix, goName)
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

	// --unique: one composite unique index over the named columns, in the
	// order they were named. It is what gives the importer's
	// OnConflict{DoNothing} something to conflict with, so a re-imported CSV
	// skips rows instead of duplicating them (grit#140).
	uniquePriority := map[string]int{}
	for i, col := range g.Definition.Unique {
		uniquePriority[toSnakeCase(strings.TrimSuffix(col, "_id"))] = i + 1
	}
	uniqueIndexName := "idx_" + names.Snake + "_unique"
	// uniqueTag returns the fragment to append to a column's gorm tag, by the
	// column's base name. Empty for a column the index does not cover.
	uniqueTag := func(base string) string {
		p, ok := uniquePriority[base]
		if !ok {
			return ""
		}
		return fmt.Sprintf(";uniqueIndex:%s,priority:%d", uniqueIndexName, p)
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

			// FK column.
			//
			// one_to_one differs from belongs_to here and nowhere else: a unique
			// index. Without it the name is a comment, because the database would
			// accept a second row pointing at the same parent.
			fkIndex := "index"
			if f.IsOneToOne() {
				fkIndex = "uniqueIndex"
			}
			// A pointer with default:null, not an empty string. Postgres
			// enforces the foreign key, and "" is not a key: it answers
			// SQLSTATE 23503 rather than storing "no relation".
			// The owner of an --owned-by row is the signed-in user: the
			// service stamps it and ignores whatever the body said. Marked
			// here so a reader of the model knows not to ask for it, which
			// is how `grit sync` came to write a create schema demanding an
			// owner id no form has (grit#108).
			ownerTag := ""
			if owner != nil && f.Name == owner.Name {
				ownerTag = ` owner:"true"`
			}
			// A foreign key can be part of the unique index too, which is how
			// "unique within a company" is spelled: --unique "name,company".
			fkUnique := uniqueTag(strings.TrimSuffix(fkJson, "_id"))
			if g.relOptional(f) {
				structFields += fmt.Sprintf("\t%s *string `gorm:\"size:36;%s%s;default:null\" json:\"%s\"%s`\n", fkGoName, fkIndex, fkUnique, fkJson, ownerTag)
			} else {
				structFields += fmt.Sprintf("\t%s string `gorm:\"size:36;%s%s\" json:\"%s\" binding:\"required\"%s`\n", fkGoName, fkIndex, fkUnique, fkJson, ownerTag)
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
		if u := uniqueTag(jsonTag); u != "" {
			// The fragment carries its own leading semicolon, so a column with
			// no other gorm tag must not start with one.
			gormTag += u
			gormTag = strings.TrimPrefix(gormTag, ";")
		}
		if gormTag != "" {
			tags = fmt.Sprintf(`gorm:"%s" %s`, gormTag, tags)
		}

		// A conditional field is never unconditionally required: binding runs
		// before anything has looked at the row, so binding:"required" here
		// would refuse every record the condition does not hold for. The
		// requirement travels as requiredif instead, which internal/fieldtypes
		// evaluates against the row on every write that has a model.
		if f.ShowWhen != nil {
			// Only when it is required while shown. A field that is merely
			// conditional, and optional when it appears, has nothing to enforce.
			if f.Required {
				tags += fmt.Sprintf(` requiredif:"%s=%s"`, toSnakeCase(f.ShowWhen.Field), f.ShowWhen.Value)
			}
		} else if f.Required && (f.GoType() == "string") && !f.IsSlug() {
			tags += ` binding:"required"`
		} else if f.Required && !f.IsSlug() {
			// Required, and not a string, so gin cannot carry the rule:
			// binding:"required" rejects 0 and false, which are values these
			// columns legitimately hold, so the generator has always left it off.
			// That left the model saying nothing at all about a required int or a
			// required money column, and `grit sync` read the silence as optional
			// (grit#108). Declared here instead: a statement for whoever reads the
			// model, with none of binding's runtime behaviour.
			tags += ` required:"true"`
		}
		if FieldType(f.Type) == FieldRichtext && !f.Encrypted {
			// Rendered as HTML wherever it is shown, so internal/sanitize cleans
			// it on every write.
			tags += ` sanitize:"html"`
		}
		// The choices of a choice column, so anything reading the model back
		// knows the closed set: a seeder generated later, the API reference,
		// and a person opening the file.
		if enum := f.EnumTag(); enum != "" {
			tags += fmt.Sprintf(` enum:"%s"`, enum)
		}
		// internal/fieldtypes checks and normalises the column on every write
		// that has a model, and the OpenAPI reference reads the docs tag.
		if format := f.FormatTag(); format != "" {
			tags += fmt.Sprintf(` format:"%s"`, format)
			if docs := f.DocsTag(); docs != "" {
				tags += fmt.Sprintf(` docs:"%s"`, docs)
			}
		}
		if f.IsSearchable() && !f.IsRelationship() {
			// The same columns List searches. grit migrate gives them a trigram
			// index on Postgres, which LIKE '%term%' can use.
			tags += ` search:"trigram"`
		}
		// A slug is derived from another column by the BeforeCreate hook
		// below, so nothing should ask a caller for one. Named here because
		// the column is otherwise an ordinary unique string: `grit sync`
		// builds the create schema from this file, and without the tag it
		// demanded a slug the hook was about to overwrite (grit#108). The
		// value is the column it comes from, which is worth saying anyway.
		// Filled by the server, so nothing should ask a caller for it. Named
		// in the model because `grit sync` builds the schemas from this file
		// and cannot tell a computed column from an ordinary one.
		if f.ServerSet {
			tags += ` server:"true"`
		}
		if f.IsSlug() {
			source := toSnakeCase(f.SlugSource)
			if source == "" {
				source = toSnakeCase(slugSourceGo)
			}
			if source == "" {
				source = "auto"
			}
			tags += fmt.Sprintf(` slug:"%s"`, source)
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
	CreatedAt time.Time      `+"`"+`gorm:"index" json:"created_at"`+"`"+`
	UpdatedAt time.Time      `+"`"+`gorm:"index" json:"updated_at"`+"`"+`
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

// writeZodSchema writes the resource's Zod schemas, from the model.
//
// The same move as writeTSTypes, for the same reason. buildZodSchema already
// decided nothing of its own: it reconstructs the field each column was
// generated from and calls Field.ZodType(), the emitter this function used to
// call directly. So the two agreed on every field and still produced different
// files, because this one wrote only the create and update pair while sync
// also writes the entity schema. grit sync is the step the generator prints as
// the next thing to run, so it rewrote this file in every project, every time.
//
// What the entity schema is for: create and update describe what goes in, and
// this describes what comes back, id and timestamps included, which is what
// you validate a response against.
func (g *Generator) writeZodSchema(names Names) error {
	modelPath := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		return fmt.Errorf("reading the model at %s: %w", modelPath, err)
	}
	for _, s := range structs {
		if s.Name != names.Pascal {
			continue
		}
		path := filepath.Join(g.SharedRoot(), "schemas", names.Kebab+".ts")
		return writeFileWithDirs(path, buildZodSchema(s))
	}
	return fmt.Errorf("no %s struct in %s", names.Pascal, modelPath)
}

// zodSchemaImports is the import block a schema file needs for the shared
// schemas its fields reference.
//
// Shared with sync, which emits the same file from the Go model: a field that
// renders as MoneySchema has to import MoneySchema whichever writer rendered
// it, and sync used to avoid the question by emitting z.unknown() instead
// (grit#108).
//
// Built as a list. Appending to a slice stays correct the day a third shared
// schema appears; rewriting an anchor line fails silently the day the anchor
// moves.
func zodSchemaImports(fields []Field) string {
	needsFileRef, needsMoney := false, false
	for _, f := range fields {
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
	// The formatted types, one import line per module, names sorted so a
	// regenerate does not reorder them.
	formatNames := map[string]map[string]bool{}
	for _, f := range fields {
		module, name := f.SharedSchema()
		if module == "" {
			continue
		}
		if formatNames[module] == nil {
			formatNames[module] = map[string]bool{}
		}
		formatNames[module][name] = true
	}
	for _, module := range []string{"./field-formats", "./phone"} {
		if len(formatNames[module]) == 0 {
			continue
		}
		var list []string
		for name := range formatNames[module] {
			list = append(list, name)
		}
		sort.Strings(list)
		imports = append(imports, fmt.Sprintf(`import { %s } from "%s";`, strings.Join(list, ", "), module))
	}
	return strings.Join(imports, "\n")
}

// conditionalRequired is a field that is required only while its condition
// holds, which Zod expresses as a refinement rather than as a type.
type conditionalRequired struct {
	Field string // the field this applies to, in snake_case
	On    string // the sibling it depends on
	Value string // the stored value, not the label
}

// refinementsFor drops the rules whose condition reads a field this schema
// does not have, and names what it dropped.
//
// A refinement is the one thing in the schema that reads a sibling. Every
// other validator only sees its own value, so excluding a field is harmless to
// them and fatal to this: the generated superRefine said `v.status === "lost"`
// for a status that `server_set` had just taken out, which is a TypeScript
// error, which stops the shared package compiling, which fails the build of
// every app that imports it (grit#115).
//
// Dropped rather than kept, because a client that is not sent the field cannot
// evaluate the condition anyway. The rule is not lost: it is on the model as a
// requiredif tag, and internal/fieldtypes enforces it on every write that has
// a model, which is the API, a PATCH, a CSV import and a job. What is lost is
// the inline message on the form, which is worth saying out loud rather than
// dropping in silence.
func refinementsFor(rules []conditionalRequired, keys map[string]bool) (kept, dropped []conditionalRequired) {
	for _, r := range rules {
		if keys[r.On] {
			kept = append(kept, r)
			continue
		}
		dropped = append(dropped, r)
	}
	return kept, dropped
}

// warnDroppedRefinements says which ones went and why.
func warnDroppedRefinements(resource string, dropped []conditionalRequired) {
	for _, r := range dropped {
		fmt.Printf("  \u26a0 %s: %s is required when %s is %q, but %s is not in the "+
			"create schema, so the browser cannot check it. The server still does, "+
			"from the model's requiredif tag.\n",
			resource, r.Field, r.On, r.Value, r.On)
	}
}

// zodRefinements renders the superRefine that enforces them, or "" when there
// are none.
//
// One refine for all of them rather than one each: a chain of superRefine
// calls reads worse and reports the same thing.
func zodRefinements(rules []conditionalRequired) string {
	if len(rules) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(".superRefine((v, ctx) => {\n")
	for _, r := range rules {
		fmt.Fprintf(&b, "  if (v.%s === %q && !v.%s) {\n", r.On, r.Value, r.Field)
		b.WriteString("    ctx.addIssue({\n")
		b.WriteString("      code: z.ZodIssueCode.custom,\n")
		fmt.Fprintf(&b, "      path: [%q],\n", r.Field)
		b.WriteString("      message: \"Required\",\n")
		b.WriteString("    });\n")
		b.WriteString("  }\n")
	}
	b.WriteString("})")
	return b.String()
}

// writeTSTypes writes the resource's TypeScript interface, from the model.
//
// It used to build the interface from the definition, which is the list of
// fields somebody asked for rather than the list the API serialises. Every
// model also carries version and archived_at, and a --tree resource carries
// path, depth, position and children, so six fields a response returns were
// fields the TypeScript said did not exist.
//
// And one writer: grit sync rewrites this same file from the model, and is the
// step the generator itself tells people to run. Two writers meant the second
// one won, which is how Money became unknown and an enum's union became
// string. buildTSType is now the only thing that writes it, and it reads what
// the API will actually send.
func (g *Generator) writeTSTypes(names Names) error {
	modelPath := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		return fmt.Errorf("reading the model at %s: %w", modelPath, err)
	}
	for _, s := range structs {
		if s.Name != names.Pascal {
			continue
		}
		path := filepath.Join(g.SharedRoot(), "types", names.Kebab+".ts")
		return writeFileWithDirs(path, buildTSType(s))
	}
	return fmt.Errorf("no %s struct in %s", names.Pascal, modelPath)
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

		// belongs_to: show what the related model is called, by the same
		// rule the form's dropdown uses. Hardcoding ".name" here left an
		// empty column for every related model that has no name column: a
		// User has first_name, an Invoice has a number, a Collection has a
		// title.
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			assocSnake := toSnakeCase(baseName)
			colLabel := strings.Join(splitPascal(toPascalCase(baseName)), " ")
			columns += fmt.Sprintf("\n      { key: \"%s.%s\", label: \"%s\" },",
				assocSnake, g.displayFieldFor(f.RelatedModelName()), colLabel)
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
		if FieldType(f.Type) == FieldRating {
			parts = append(parts, fmt.Sprintf(`ratingMax: %d`, f.RatingMax()))
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
	owner := g.Definition.OwnerField()
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		// The owner of an --owned-by row is whoever is signed in. The service
		// stamps it on create and nothing writes it afterwards, so a picker
		// here is a required field whose answer is thrown away: it refused to
		// submit until an operator chose a user, and then stored a different
		// one. The column stays in the table, which is a read.
		if owner != nil && f.Name == owner.Name {
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
			// Optional in the form when the relationship is optional. A
			// required select gives the person no way to say "none", which is
			// wrong for the top of a tree and wrong for a deal with no contact
			// yet.
			required := "required: true, "
			if g.relOptional(f) {
				required = ""
			}
			formFields += fmt.Sprintf("\n    { key: \"%s\", label: \"%s\", type: \"relationship-select\", %srelatedEndpoint: \"/api/%s\", displayField: \"%s\" },",
				fkKey, fieldLabel, required, relPlural, g.displayFieldFor(f.RelatedModelName()))
			continue
		}

		// many_to_many: multi-relationship-select with endpoint
		if f.IsManyToMany() {
			idsKey := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			fieldLabel := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
			relSnake := toSnakeCase(f.RelatedModelName())
			relPlural := Pluralize(relSnake)
			assocKey := toSnakeCase(f.Name)
			formFields += fmt.Sprintf("\n    { key: \"%s\", label: \"%s\", type: \"multi-relationship-select\", relatedEndpoint: \"/api/%s\", displayField: \"%s\", relationshipKey: \"%s\" },",
				idsKey, fieldLabel, relPlural, g.displayFieldFor(f.RelatedModelName()), assocKey)
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
		// A conditional field's requirement is a predicate rather than a flag:
		// required: true here would hold even while the field is hidden, which
		// is a form that will not submit and does not say why.
		if f.Required && f.ShowWhen == nil {
			parts = append(parts, `required: true`)
		}
		if f.ShowWhen != nil {
			pred := fmt.Sprintf(`(v) => v.%s === %q`, toSnakeCase(f.ShowWhen.Field), f.ShowWhen.Value)
			parts = append(parts, "visibleWhen: "+pred)
			if f.Required {
				parts = append(parts, "requiredWhen: "+pred)
			}
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
		case FieldTel:
			// The country the picker starts on and local numbers are read in.
			if f.DefaultCountry != "" {
				parts = append(parts, fmt.Sprintf(`defaultCountry: %q`, f.DefaultCountry))
			}
		case FieldCountry:
			if f.DefaultCountry != "" {
				parts = append(parts, fmt.Sprintf(`defaultValue: %q`, f.DefaultCountry))
			}
		case FieldRating:
			parts = append(parts, fmt.Sprintf(`max: %d`, f.RatingMax()))
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
				itemCols += fmt.Sprintf("\n        { key: %q, label: %q, type: \"relationship-select\", relatedEndpoint: \"/api/%s\", displayField: %q },", cf.FKColumnName(), label, relKebab, g.displayFieldFor(cf.RelatedModelName()))
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

	statsBlock, insightsBlock := g.insightsBlocks()

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
  },%s%s
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
		statsBlock,
		insightsBlock,
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
