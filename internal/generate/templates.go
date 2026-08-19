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
	projImports := []string{}
	if needsFiles {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/files\"", g.Module))
	}
	projImports = append(projImports, fmt.Sprintf("\"%s/internal/ids\"", g.Module))
	if hasAuto {
		projImports = append(projImports, fmt.Sprintf("\"%s/internal/sequence\"", g.Module))
	}
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

			// FK column
			if selfRef {
				structFields += fmt.Sprintf("\t%s *string `gorm:\"size:36;index\" json:\"%s\"`\n", fkGoName, fkJson)
			} else {
				structFields += fmt.Sprintf("\t%s string `gorm:\"size:36;index\" json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
			}
			// Association struct
			assocType := relModel
			if selfRef {
				assocType = "*" + relModel
			}
			structFields += fmt.Sprintf("\t%s %s `gorm:\"foreignKey:%s\" json:\"%s\"`\n",
				assocName, assocType, fkGoName, toSnakeCase(assocName))
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

	// Path resolution, the cycle refusal, and the two accessors a breadcrumb
	// needs. Reparenting lives in the tree service instead: a move rewrites a
	// whole subtree, and that belongs somewhere it can be a transaction.
	if isTree {
		content += treeModelMethods(names)
	}

	path := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// writeGoService creates the service layer for the resource.
func (g *Generator) writeGoService(names Names) error {
	searchWhere := g.buildServiceSearchWhere()

	r := strings.NewReplacer(
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{lower}}", names.Lower,
		"{{plural}}", names.Plural,
		"{{SEARCH_WHERE}}", searchWhere,
		"{{SORTABLE_SET}}", g.buildSortableSet(),
	)

	content := r.Replace(`package services

import (
	"fmt"
	"math"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// {{Pascal}}Service handles business logic for {{plural}}.
type {{Pascal}}Service struct {
	DB *gorm.DB
}

// {{Pascal}}ListParams holds pagination and filter parameters.
type {{Pascal}}ListParams struct {
	Page      int
	PageSize  int
	Search    string
	SortBy    string
	SortOrder string
}

// List returns a paginated list of {{plural}}.
func (s *{{Pascal}}Service) List(params {{Pascal}}ListParams) ([]models.{{Pascal}}, int64, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "desc"
	}
	// SortBy is interpolated into ORDER BY below, so it MUST be whitelisted
	// against real columns — never trust a client-supplied sort column.
	sortable{{Pascal}} := {{SORTABLE_SET}}
	if !sortable{{Pascal}}[params.SortBy] {
		params.SortBy = "created_at"
	}

	query := s.DB.Model(&models.{{Pascal}}{})

	if params.Search != "" {
		query = query.Where({{SEARCH_WHERE}})
	}

	var total int64
	query.Count(&total)

	var items []models.{{Pascal}}
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(params.SortBy + " " + params.SortOrder).Offset(offset).Limit(params.PageSize).Find(&items).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("fetching {{plural}}: %w", err)
	}

	pages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	return items, total, pages, nil
}

// GetByID returns a single {{lower}} by ID.
func (s *{{Pascal}}Service) GetByID(id string) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.DB.First(&item, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("{{lower}} not found: %w", err)
	}
	return &item, nil
}

// Create creates a new {{lower}}.
func (s *{{Pascal}}Service) Create(item *models.{{Pascal}}) error {
	if err := s.DB.Create(item).Error; err != nil {
		return fmt.Errorf("creating {{lower}}: %w", err)
	}
	return nil
}

// Update modifies an existing {{lower}}. Two queries: First() loads
// the row + verifies existence; Updates() persists the diff. The
// loaded struct is mutated by Updates() so we can return it directly
// without a third refetch.
func (s *{{Pascal}}Service) Update(id string, updates map[string]interface{}) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.DB.First(&item, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("{{lower}} not found: %w", err)
	}

	if err := s.DB.Model(&item).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("updating {{lower}}: %w", err)
	}

	return &item, nil
}

// Delete soft-deletes a {{lower}}. One query — we don't need to load
// the row first; GORM's Delete is atomic and rows-affected tells us
// whether it existed.
func (s *{{Pascal}}Service) Delete(id string) error {
	res := s.DB.Where("id = ?", id).Delete(&models.{{Pascal}}{})
	if res.Error != nil {
		return fmt.Errorf("deleting {{lower}}: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("{{lower}} not found")
	}
	return nil
}
`)

	path := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// buildServiceSearchWhere creates the full Where(...) arguments for the service search.
func (g *Generator) buildServiceSearchWhere() string {
	var searchFields []string
	for _, f := range g.Definition.Fields {
		if f.GoType() == "string" {
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
		case f.IsFile(), f.IsFiles(), f.IsManyToMany(), f.IsStringArray():
			// not a sortable scalar column
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

// writeGoHandler creates the Gin handler file for the resource.
func (g *Generator) writeGoHandler(names Names) error {
	// Build create/update request struct fields
	createFields := ""
	updateFields := ""
	createAssignments := ""
	updateMap := ""
	m2mCreateCode := ""
	m2mUpdateCode := ""
	// v3.31.18: writable-field whitelist for the new Patch handler.
	// Built once at generation time so the Patch endpoint can refuse any
	// JSON key that isn't a model field — preventing operators (or
	// malicious clients) from setting id, created_at, deleted_at, etc.
	patchAllowed := ""

	// Collect preloads for relationships
	var preloads []string

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}

		// belongs_to: add FK field (string UUID) to create/update
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkGoName := toPascalCase(baseName) + "ID"
			fkJson := toSnakeCase(baseName) + "_id"
			assocName := toPascalCase(baseName)
			preloads = append(preloads, assocName)

			// A self-reference is optional, because the root of a tree has no
			// parent. Required here and the API refuses to create the first row
			// with a validation error about a field the caller correctly left
			// empty, which is where this was first noticed.
			selfRef := f.RelatedModelName() == toPascalCase(g.Definition.Name)
			if selfRef {
				createFields += fmt.Sprintf("\t\t%s string `json:\"%s\"`\n", fkGoName, fkJson)
			} else {
				createFields += fmt.Sprintf("\t\t%s string `json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
			}

			if selfRef {
				// The column is nullable, so "" from an empty form select has to
				// become nil rather than reaching the database. An empty string
				// is not a row any foreign key constraint can point at, and
				// Postgres says so with SQLSTATE 23503 while SQLite quietly
				// stores a dangling reference.
				createAssignments += fmt.Sprintf("\t\t%s: optional%sID(req.%s),\n", fkGoName, names.Pascal, fkGoName)
			} else {
				createAssignments += fmt.Sprintf("\t\t%s: req.%s,\n", fkGoName, fkGoName)
			}

			updateFields += fmt.Sprintf("\t\t%s *string `json:\"%s\"`\n", fkGoName, fkJson)
			if selfRef {
				updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = optional%sID(*req.%s)\n\t}\n", fkGoName, fkJson, names.Pascal, fkGoName)
			} else {
				updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", fkGoName, fkJson, fkGoName)
			}
			patchAllowed += fmt.Sprintf("\t\t\"%s\": true,\n", fkJson)
			continue
		}

		// many_to_many: add []string for create, *[]string for update
		// (all related models use UUID string PKs).
		if f.IsManyToMany() {
			relModel := f.RelatedModelName()
			assocName := toPascalCase(f.Name)
			idsName := toPascalCase(f.Name) + "IDs"
			idsJson := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			preloads = append(preloads, assocName)

			createFields += fmt.Sprintf("\t\t%s []string `json:\"%s\"`\n", idsName, idsJson)
			updateFields += fmt.Sprintf("\t\t%s *[]string `json:\"%s\"`\n", idsName, idsJson)

			varName := toSnakeCase(f.Name)
			m2mCreateCode += fmt.Sprintf("\n\tif len(req.%s) > 0 {\n\t\tvar %s []models.%s\n\t\th.DB.Find(&%s, req.%s)\n\t\th.DB.Model(&item).Association(\"%s\").Replace(%s)\n\t}\n", idsName, varName, relModel, varName, idsName, assocName, varName)
			m2mUpdateCode += fmt.Sprintf("\n\tif req.%s != nil {\n\t\tvar %s []models.%s\n\t\tif len(*req.%s) > 0 {\n\t\t\th.DB.Find(&%s, *req.%s)\n\t\t}\n\t\th.DB.Model(&item).Association(\"%s\").Replace(%s)\n\t}\n", idsName, varName, relModel, idsName, varName, idsName, assocName, varName)
			// m2m IDs aren't a plain column — the Patch handler skips
			// them. (Set them via the full Update endpoint instead.)
			continue
		}

		goName := toPascalCase(f.Name)
		goType := f.GoType()
		jsonTag := toSnakeCase(f.Name)

		bindingTag := ""
		if f.Required {
			bindingTag = ` binding:"required"`
		}

		createFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"%s`\n", goName, goType, jsonTag, bindingTag)
		createAssignments += fmt.Sprintf("\t\t%s: req.%s,\n", goName, goName)
		patchAllowed += fmt.Sprintf("\t\t\"%s\": true,\n", jsonTag)

		// For update, use pointer types to detect "provided" vs "missing"
		if goType == "bool" {
			updateFields += fmt.Sprintf("\t\t%s *%s `json:\"%s\"`\n", goName, goType, jsonTag)
			updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", goName, jsonTag, goName)
		} else if goType == "string" {
			updateFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", goName, goType, jsonTag)
			updateMap += fmt.Sprintf("	if req.%s != \"\" {\n\t\tupdates[\"%s\"] = req.%s\n\t}\n", goName, jsonTag, goName)
		} else if goType == "*time.Time" {
			updateFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", goName, goType, jsonTag)
			updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = req.%s\n\t}\n", goName, jsonTag, goName)
		} else {
			updateFields += fmt.Sprintf("\t\t%s *%s `json:\"%s\"`\n", goName, goType, jsonTag)
			updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", goName, jsonTag, goName)
		}
	}

	// Inline has-many items (from --items): a request slice, an atomic create
	// (GORM cascades the children in Create's transaction), and a replace-all on
	// update. The parent FK is set by the association (create) or by hand
	// (update) — clients never send it per row.
	itemsReqField := ""
	itemsBuild := ""
	itemsUpdate := ""
	if g.Definition.Items != nil {
		childModel := BuildNames(g.Definition.Items).Pascal
		reqStructFields := ""
		assign := ""
		for _, cf := range g.Definition.Items.Fields {
			if cf.IsBelongsTo() || cf.IsSlug() || cf.IsManyToMany() {
				continue
			}
			gName := toPascalCase(cf.Name)
			gType := cf.GoType()
			jTag := toSnakeCase(cf.Name)
			reqStructFields += fmt.Sprintf("\t\t\t%s %s `json:\"%s\"`\n", gName, gType, jTag)
			assign += fmt.Sprintf("\t\t\t\t%s: it.%s,\n", gName, gName)
		}
		itemsReqField = fmt.Sprintf("\t\tItems []struct {\n%s\t\t} `json:\"items\"`\n", reqStructFields)
		itemsBuild = fmt.Sprintf("\n\tif len(req.Items) > 0 {\n\t\titems := make([]models.%s, 0, len(req.Items))\n\t\tfor _, it := range req.Items {\n\t\t\titems = append(items, models.%s{\n%s\t\t\t})\n\t\t}\n\t\titem.Items = items\n\t}\n", childModel, childModel, assign)
		fkCol := names.Snake + "_id"
		itemsUpdate = fmt.Sprintf("\n\tif req.Items != nil {\n\t\th.DB.Where(\"%s = ?\", item.ID).Delete(&models.%s{})\n\t\tif len(req.Items) > 0 {\n\t\t\tnewItems := make([]models.%s, 0, len(req.Items))\n\t\t\tfor _, it := range req.Items {\n\t\t\t\trow := models.%s{\n%s\t\t\t\t}\n\t\t\t\trow.%sID = item.ID\n\t\t\t\tnewItems = append(newItems, row)\n\t\t\t}\n\t\t\th.DB.Create(&newItems)\n\t\t}\n\t}\n", fkCol, childModel, childModel, childModel, assign, names.Pascal)
	}
	createFields += itemsReqField
	updateFields += itemsReqField

	// Build preload chain
	preloadChain := ""
	for _, p := range preloads {
		preloadChain += fmt.Sprintf(".Preload(\"%s\")", p)
	}
	// Load inline items on reads so the detail form + related table are populated.
	if g.Definition.Items != nil {
		preloadChain += ".Preload(\"Items\")"
	}

	// Build reload-with-preloads line (used after Create/Update)
	reloadLine := "\th.DB.First(&item, \"id = ?\", item.ID)"
	if preloadChain != "" {
		reloadLine = fmt.Sprintf("\th.DB%s.First(&item, \"id = ?\", item.ID)", preloadChain)
	}

	// A write that is exactly one statement needs neither the reload nor the
	// transaction GORM wraps writes in.
	//
	// The reload existed to pick up columns the database fills in — defaults,
	// values assigned in a hook. RETURNING gets those from the INSERT itself,
	// so the follow-up SELECT is a second round trip buying nothing. Where the
	// model has relations it still earns its place, because RETURNING cannot
	// populate a preloaded association.
	//
	// The transaction is the same argument. GORM wraps every write because a
	// Create expands into several INSERTs once a model has children, and a
	// half-written parent is worse than a slow one. When the generator can see
	// there are no children, no join rows and no sequence hook writing
	// alongside, the statement is already atomic in Postgres and BEGIN/COMMIT
	// costs two round trips for a guarantee that is already held.
	//
	// The global DB_SKIP_DEFAULT_TRANSACTION stays off. This is decided per
	// resource, from the definition, where it can be known rather than assumed.
	hasAutoField := false
	for _, f := range g.Definition.Fields {
		if f.Auto {
			hasAutoField = true
			break
		}
	}
	singleStatementWrite := preloadChain == "" &&
		itemsBuild == "" &&
		m2mCreateCode == "" &&
		!hasAutoField

	createCall := "h.DB.Create(&item)"
	updateCall := "h.DB.Model(&item).Updates(updates)"
	createReload := reloadLine
	updateReload := reloadLine
	clauseImport := ""
	databaseImport := ""

	// Only a resource whose status field is a state machine imports these.
	// An unused import is a build failure, so an ordinary resource must not
	// get them.
	workflowImport := ""
	if g.Definition.WorkflowField() != nil {
		// g.Module directly, not a {{MODULE}} placeholder: NewReplacer is
		// single-pass, so a placeholder inside an inserted value is never
		// expanded and lands in the file verbatim.
		workflowImport = "\n\t\"" + g.Module + "/internal/authz\"" +
			"\n\t\"" + g.Module + "/internal/workflow\""
	}
	if singleStatementWrite {
		// database.Write applies RETURNING only where the dialect has it, so
		// the reload has to stay for the dialects that do not. It is guarded
		// rather than dropped: on Postgres and SQLite this costs one boolean,
		// and on MySQL it is the difference between a complete record and a
		// half-empty one.
		createCall = "database.Write(h.DB).Create(&item)"
		updateCall = "database.Write(h.DB).Model(&item).Updates(updates)"
		guard := "\tif !database.SupportsReturning(h.DB) {\n\t"
		createReload = guard + reloadLine + "\n\t}"
		updateReload = guard + reloadLine + "\n\t}"
		databaseImport = "\n\t\"" + g.Module + "/internal/database\""
	}

	// Build allowed sort columns (skip relationship fields)
	sortCols := `"id": true, "created_at": true`
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() {
			continue
		}
		if f.GoType() == "string" || f.GoType() == "int" || f.GoType() == "uint" {
			sortCols += fmt.Sprintf(`, "%s": true`, toSnakeCase(f.Name))
		}
	}

	// Columns the admin may filter on from the query string. Wider than the
	// sortable set: a bool or a foreign key is worth filtering by and pointless
	// to sort by. Whitelisted rather than open, because the column name goes
	// into the WHERE clause.
	filterCols := `"id": true`
	for _, f := range g.Definition.Fields {
		if f.IsManyToMany() {
			continue
		}
		col := toSnakeCase(f.Name)
		if f.IsBelongsTo() {
			col = f.FKColumnName()
		}
		filterCols += fmt.Sprintf(`, "%s": true`, col)
	}

	searchCols := g.buildHandlerSearchCols()

	// Build export columns from the field list. Skips relationships
	// (we don't try to traverse associations in the default export).
	// Time fields get a "date:..." format string; bools get "bool".
	exportCols := "\t\t\t{Header: \"ID\", Field: \"ID\"},\n"
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() {
			continue
		}
		header := strings.ReplaceAll(toPascalCase(f.Name), "ID", " ID")
		fieldGoName := toPascalCase(f.Name)
		format := ""
		switch {
		case f.GoType() == "time.Time" || f.GoType() == "*time.Time":
			format = "date:2006-01-02"
		case f.GoType() == "bool":
			format = "bool"
		}
		if format != "" {
			exportCols += fmt.Sprintf("\t\t\t{Header: %q, Field: %q, Format: %q},\n", header, fieldGoName, format)
		} else {
			exportCols += fmt.Sprintf("\t\t\t{Header: %q, Field: %q},\n", header, fieldGoName)
		}
	}
	exportCols += "\t\t\t{Header: \"Created At\", Field: \"CreatedAt\", Format: \"date:2006-01-02\"},"

	// Check if any field needs "time" import
	needsTimeImport := false
	needsHandlerDatatypes := false
	hasFileFields := false
	for _, f := range g.Definition.Fields {
		if f.GoType() == "*time.Time" {
			needsTimeImport = true
		}
		if f.NeedsDatatypesImport() {
			needsHandlerDatatypes = true
		}
		if f.IsFileField() {
			hasFileFields = true
		}
	}
	// The PDF endpoint always stamps the footer with time.Now() and reads
	// APP_NAME, so "time" and "os" are imported unconditionally now — the
	// needsTimeImport check below is kept only to document why other code
	// paths wanted it.
	_ = needsTimeImport
	// strconv joined the unconditional set with the Bulk handler.
	timeImport := "\n\t\"os\"\n\t\"strconv\"\n\t\"time\""
	datatypesImport := ""
	if needsHandlerDatatypes {
		datatypesImport = "\n\t\"gorm.io/datatypes\""
	}

	// v3.31.33 -- file lifecycle. When the resource has any file/files
	// fields, the handler gets a Storage field and the Create + Update
	// flows pick up cleanup-on-replace + claim-on-save. The
	// {{HANDLER_*}} replacements degrade to empty strings on resources
	// without file fields so non-file handlers stay unchanged.
	filesImport := ""
	handlerStorageField := ""
	createClaim := ""
	updateSnapshot := ""
	updateCleanup := ""
	if hasFileFields {
		filesImport = "\n\t\"" + g.Module + "/internal/files\"\n\t\"" + g.Module + "/internal/storage\""
		handlerStorageField = "\n\tStorage *storage.Storage // v3.31.33"
		createClaim = "\n\tif h.Storage != nil {\n\t\tfiles.ClaimRefs(c.Request.Context(), h.DB, &item)\n\t}"
		updateSnapshot = "\n\toldItem := item // v3.31.33: snapshot for file diff"
		updateCleanup = "\n\tif h.Storage != nil {\n\t\tfiles.CleanupRemoved(c.Request.Context(), h.Storage, &oldItem, &item)\n\t\tfiles.ClaimRefs(c.Request.Context(), h.DB, &item)\n\t}"
	}

	// v3.31.39: identifier expression for the activity log calls --
	// picks the first human-readable field on the model (Name, Title,
	// Slug, ...) and falls back to item.ID so the {verb} {entityType}
	// {identifier} log line is never blank.
	identExpr := pickIdentifierExpr(g.Definition.Fields)

	// belongs_to resources are filterable by their foreign key, so
	// GET /<plural>?<fk>=<id> returns only the children of that parent
	// (e.g. /products?category_id=... for a category's products). .With
	// ignores empty values, so each filter is inert when its param is absent.
	fkFilters := ""
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			fk := f.FKColumnName()
			fkFilters += fmt.Sprintf(".With(%q, c.Query(%q))", fk, fk)
		}
	}

	// ── PDF endpoint ─────────────────────────────────────────────────────────
	// Build the detail grid from the resource's own fields. Binary/collection
	// types (files, images, video) and many-to-many have no sensible text
	// rendering, so they're left out; belongs_to prints the related record's
	// human label via pdf.Display, which reflects for Name/Title/Email/etc.
	pdfFields := ""
	for _, f := range g.Definition.Fields {
		if f.IsManyToMany() || f.IsSlug() {
			continue
		}
		switch f.FormFieldType() {
		case "image", "images", "video", "videos", "file", "files", "richtext":
			continue
		}
		label := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
		if f.IsBelongsTo() {
			assoc := toPascalCase(strings.TrimSuffix(f.Name, "_id"))
			pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Display(item.%s)},\n", label, assoc)
			continue
		}
		pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Value(item.%s)},\n", label, toPascalCase(f.Name))
	}
	pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Value(item.CreatedAt)},\n", "Created")

	// Line items become a table section, with a Total column when the child
	// carries a quantity and a rate (the same name-based pairing the admin's
	// inline line-items editor uses).
	pdfSections := ""
	if g.Definition.Items != nil {
		childNames := BuildNames(g.Definition.Items)
		headers := ""
		cells := ""
		aligns := ""
		for _, cf := range g.Definition.Items.Fields {
			if cf.IsBelongsTo() || cf.IsManyToMany() {
				continue
			}
			switch cf.FormFieldType() {
			case "image", "images", "video", "videos", "file", "files", "richtext":
				continue
			}
			headers += fmt.Sprintf("%q, ", strings.Join(splitPascal(toPascalCase(cf.Name)), " "))
			cells += fmt.Sprintf("pdf.Value(row.%s), ", toPascalCase(cf.Name))
			if cf.FormFieldType() == "number" {
				aligns += `"R", `
			} else {
				aligns += `"L", `
			}
		}
		pdfSections = fmt.Sprintf(`
	itemRows := make([][]string, 0, len(item.Items))
	for _, row := range item.Items {
		itemRows = append(itemRows, []string{%s})
	}
	if len(itemRows) > 0 {
		rec.Sections = append(rec.Sections, pdf.Section{
			Title:   %q,
			Headers: []string{%s},
			Aligns:  []string{%s},
			Rows:    itemRows,
		})
	}
`, strings.TrimSuffix(cells, ", "),
			strings.Join(splitPascal(childNames.PluralPascal), " "),
			strings.TrimSuffix(headers, ", "), strings.TrimSuffix(aligns, ", "))
	}

	// optionalID is emitted only for a resource with a self-reference, because
	// that is the only nullable foreign key the generator produces, and an
	// unused function is a compile error.
	optionalIDHelper := ""
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() && f.RelatedModelName() == toPascalCase(g.Definition.Name) {
			// Expanded here rather than left to the replacer below: a
			// strings.Replacer is single pass, so {{Pascal}} inside text it
			// just inserted is never looked at again.
			optionalIDHelper = strings.ReplaceAll(OPTIONAL_ID_HELPER_SRC, "{{Pascal}}", names.Pascal)
			break
		}
	}

	r := strings.NewReplacer(
		"{{OPTIONAL_ID_HELPER}}", optionalIDHelper,
		"{{FK_FILTERS}}", fkFilters,
		"{{UPPER_LABEL}}", strings.ToUpper(strings.Join(splitPascal(names.Pascal), " ")),
		"{{PDF_SUBTITLE}}", "pdf.Value("+identExpr+")",
		"{{PDF_FIELDS}}", pdfFields,
		"{{PDF_SECTIONS}}", pdfSections,
		"{{kebab}}", names.Kebab,
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{lower}}", names.Lower,
		"{{plural}}", names.Plural,
		"{{Plural}}", names.PluralPascal,
		"{{SORT_COLS}}", sortCols,
		"{{FILTER_COLS}}", filterCols,
		"{{SEARCH_COLS}}", searchCols,
		"{{EXPORT_COLS}}", exportCols,
		"{{CREATE_FIELDS}}", createFields,
		"{{CREATE_FIELDS_TOP}}", dedentOneTab(createFields),
		"{{UPDATE_FIELDS_TOP}}", dedentOneTab(updateFields),
		"{{CREATE_ASSIGN}}", createAssignments,
		"{{UPDATE_FIELDS}}", updateFields,
		"{{UPDATE_MAP}}", updateMap,
		"{{PATCH_ALLOWED}}", patchAllowed,
		"{{PRELOADS}}", preloadChain,
		"{{M2M_CREATE}}", m2mCreateCode,
		"{{M2M_UPDATE}}", m2mUpdateCode,
		"{{ITEMS_BUILD}}", itemsBuild,
		"{{ITEMS_UPDATE}}", itemsUpdate,
		"{{RELOAD}}", reloadLine,
		"{{CREATE_CALL}}", createCall,
		"{{UPDATE_CALL}}", updateCall,
		"{{CREATE_RELOAD}}", createReload,
		"{{UPDATE_RELOAD}}", updateReload,
		"{{CLAUSE_IMPORT}}", clauseImport,
		"{{DATABASE_IMPORT}}", databaseImport,
		"{{WORKFLOW_IMPORT}}", workflowImport,
		"{{TIME_IMPORT}}", timeImport,
		"{{DATATYPES_IMPORT}}", datatypesImport,
		"{{FILES_IMPORT}}", filesImport,
		"{{HANDLER_STORAGE_FIELD}}", handlerStorageField,
		"{{CREATE_CLAIM}}", createClaim,
		"{{UPDATE_SNAPSHOT}}", updateSnapshot,
		"{{UPDATE_CLEANUP}}", updateCleanup,
		"{{IDENT_EXPR}}", identExpr,
	)

	content := r.Replace(`package handlers

import (
	"net/http"{{TIME_IMPORT}}

	"github.com/gin-gonic/gin"{{DATATYPES_IMPORT}}
	"gorm.io/gorm"{{CLAUSE_IMPORT}}

	"{{MODULE}}/internal/events"
	"{{MODULE}}/internal/export"{{FILES_IMPORT}}{{DATABASE_IMPORT}}{{WORKFLOW_IMPORT}}
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
	"{{MODULE}}/internal/pdf"
	"{{MODULE}}/internal/services"
)

// {{Pascal}}Handler handles {{lower}} endpoints.
type {{Pascal}}Handler struct {
	DB *gorm.DB{{HANDLER_STORAGE_FIELD}}
}
{{OPTIONAL_ID_HELPER}}

// List returns a paginated list of {{plural}}.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *{{Pascal}}Handler) List(c *gin.Context) {
	query := h.DB.Model(&models.{{Pascal}}{}){{PRELOADS}}

	// Archived rows are excluded by default. Anything else means an operator
	// archives twelve rows, sees the count go down, and finds them again the
	// next time somebody sorts by a different column.
	switch c.Query("archived") {
	case "true", "1":
		query = query.Where("archived_at IS NOT NULL")
	case "all":
		// no filter
	default:
		query = query.Where("archived_at IS NULL")
	}

	res, err := paginate.List[models.{{Pascal}}](
		query,
		paginate.Bind(c){{FK_FILTERS}},
		paginate.Config{
			Searchable: []string{{{SEARCH_COLS}}},
			Sortable:   map[string]bool{{{SORT_COLS}}},
			Filterable: map[string]bool{{{FILTER_COLS}}},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch {{plural}}",
			},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX.
// Honours the same search/filter query params as List but skips
// pagination — you get every matching row in one file.
//
// Memory-bounded: reads in chunks of exportBatchSize so a million-row
// export doesn't OOM the process. CSV streams directly to the response
// writer; XLSX has to buffer (excelize requires the full sheet in
// memory before Write), so we still chunk the SCAN to avoid loading
// every row at once.
//
//	GET /api/{{plural}}/export?format=csv
//	GET /api/{{plural}}/export?format=xlsx&search=foo
func (h *{{Pascal}}Handler) Export(c *gin.Context) {
	const exportBatchSize = 1000

	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	query := h.DB.Model(&models.{{Pascal}}{}){{PRELOADS}}.Order("created_at desc")
	if search != "" && len([]string{{{SEARCH_COLS}}}) > 0 {
		// Reuse the same searchable columns as List.
		searchable := []string{{{SEARCH_COLS}}}
		clause := ""
		args := []any{}
		wild := "%" + search + "%"
		for i, col := range searchable {
			if i > 0 {
				clause += " OR "
			}
			clause += "LOWER(" + col + ") LIKE LOWER(?)"
			args = append(args, wild)
		}
		query = query.Where(clause, args...)
	}

	opts := export.Options{
		Sheet: "{{Plural}}",
		Columns: []export.Column{
{{EXPORT_COLS}}
		},
	}

	// Stream rows in batches via GORM's FindInBatches. CSV writes each
	// batch straight to the wire; XLSX accumulates into a slice (no
	// streaming API in excelize) but at least we never load the whole
	// table at once.
	if format == "xlsx" {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.xlsx"` + "`" + `)
		var all []models.{{Pascal}}
		if err := query.FindInBatches(&[]models.{{Pascal}}{}, exportBatchSize, func(tx *gorm.DB, batch int) error {
			var rows []models.{{Pascal}}
			if err := tx.Scan(&rows).Error; err != nil {
				return err
			}
			all = append(all, rows...)
			return nil
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "EXPORT_FAILED", "message": err.Error()},
			})
			return
		}
		if err := export.XLSX(c.Writer, all, opts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "EXPORT_FAILED", "message": err.Error()},
			})
		}
		return
	}

	// CSV path — true streaming. Write headers once, then each batch
	// flushes its rows directly to the response writer.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.csv"` + "`" + `)

	headerWritten := false
	if err := query.FindInBatches(&[]models.{{Pascal}}{}, exportBatchSize, func(tx *gorm.DB, batch int) error {
		var rows []models.{{Pascal}}
		if err := tx.Scan(&rows).Error; err != nil {
			return err
		}
		if !headerWritten {
			if err := export.CSV(c.Writer, rows, opts); err != nil {
				return err
			}
			headerWritten = true
		} else {
			// Subsequent batches: write rows only, no header.
			if err := export.CSVRows(c.Writer, rows, opts); err != nil {
				return err
			}
		}
		return nil
	}).Error; err != nil {
		// Headers already sent — best we can do is log + truncate.
		// The client will see a malformed CSV; ops should re-run.
		// (We don't write a JSON error body once streaming has begun.)
		_ = err
	}
}

// GetByID returns a single {{lower}} by ID.
func (h *{{Pascal}}Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB{{PRELOADS}}.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this {{lower}} as a print-ready PDF — a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *{{Pascal}}Handler) PDF(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB{{PRELOADS}}.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
		return
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "{{Pascal}}"
	}

	rec := pdf.Record{
		Title:      "{{UPPER_LABEL}}",
		Subtitle:   {{PDF_SUBTITLE}},
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
{{PDF_FIELDS}}		},
	}
{{PDF_SECTIONS}}
	out, err := pdf.RenderRecord(rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PDF_ERROR",
				"message": "could not render the PDF",
			},
		})
		return
	}

	filename := "{{kebab}}-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// Create{{Pascal}}Request is the JSON body accepted by POST /{{plural}}.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type, and an anonymous
// struct inside a handler gives it nothing to reflect over. routes.go passes
// this type to docs.Route(...).RequestBody().
type Create{{Pascal}}Request struct {
{{CREATE_FIELDS_TOP}}}

// Update{{Pascal}}Request is the JSON body accepted by PUT /{{plural}}/:id.
// Every field is optional — only what the client sends is applied.
type Update{{Pascal}}Request struct {
{{UPDATE_FIELDS_TOP}}}

// Create adds a new {{lower}}.
func (h *{{Pascal}}Handler) Create(c *gin.Context) {
	var req Create{{Pascal}}Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	item := models.{{Pascal}}{
{{CREATE_ASSIGN}}	}
{{ITEMS_BUILD}}
	if err := {{CREATE_CALL}}.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create {{lower}}",
			},
		})
		return
	}
{{M2M_CREATE}}
{{CREATE_RELOAD}}{{CREATE_CLAIM}}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "created", item.ID, {{IDENT_EXPR}}, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "{{Pascal}} created successfully",
	})
}

// Update modifies an existing {{lower}}.
func (h *{{Pascal}}Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
		return
	}

	var req Update{{Pascal}}Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}
{{UPDATE_SNAPSHOT}}
	updates := map[string]interface{}{}
{{UPDATE_MAP}}
	if err := {{UPDATE_CALL}}.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update {{lower}}",
			},
		})
		return
	}
{{M2M_UPDATE}}{{ITEMS_UPDATE}}
{{UPDATE_RELOAD}}{{UPDATE_CLEANUP}}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "updated", item.ID, {{IDENT_EXPR}}, services.DiffSummary(updates), nil, item)

	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "{{Pascal}} updated successfully",
	})
}

// Patch applies a partial update to a {{lower}}. Used by the admin's
// grouped update view — each form group's Save button calls PATCH with
// only the fields it owns, so editing "Address" doesn't rewrite
// "Pricing". Refuses any key that isn't a writable model column.
func (h *{{Pascal}}Handler) Patch(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Whitelist: only writable model columns may be patched. id,
	// created_at, updated_at, deleted_at, version are owned by the
	// framework and silently dropped here.
	allowed := map[string]bool{
{{PATCH_ALLOWED}}	}
	updates := map[string]interface{}{}
	for k, v := range body {
		if allowed[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "No writable fields in request body",
			},
		})
		return
	}

	if err := h.DB.Model(&item).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to patch {{lower}}",
			},
		})
		return
	}
{{RELOAD}}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "updated", item.ID, {{IDENT_EXPR}}, services.DiffSummary(updates), nil, item)

	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "{{Pascal}} updated successfully",
	})
}

// Delete soft-deletes a {{lower}}.
func (h *{{Pascal}}Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
		return
	}

	if err := h.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete {{lower}}",
			},
		})
		return
	}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "deleted", item.ID, {{IDENT_EXPR}}, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "{{Pascal}} deleted successfully",
	})
}

// Bulk{{Pascal}}Request is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails: the operator sees "failed" while ten rows
// are already gone.
type Bulk{{Pascal}}Request struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string ` + "`" + `json:"action" binding:"required,oneof=delete archive restore patch"` + "`" + `
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string ` + "`" + `json:"ids" binding:"required,min=1,max=500"` + "`" + `
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} ` + "`" + `json:"patch"` + "`" + `
}

// Bulk applies one action to many {{plural}} in a single transaction.
func (h *{{Pascal}}Handler) Bulk(c *gin.Context) {
	var req Bulk{{Pascal}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Unarchived scope for archive, archived scope for restore: without it a
	// mixed selection reports "12 archived" having changed three rows.
	var items []models.{{Pascal}}
	scope := h.DB.Where("id IN ?", req.IDs)
	if req.Action == "restore" {
		scope = scope.Where("archived_at IS NOT NULL")
	} else if req.Action == "archive" {
		scope = scope.Where("archived_at IS NULL")
	}
	if err := scope.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to load {{plural}}"},
		})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    gin.H{"affected": 0, "requested": len(req.IDs)},
			"message": "Nothing to do",
		})
		return
	}

	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	var updates map[string]interface{}
	if req.Action == "patch" {
		// Same whitelist as Patch. Framework-owned columns are dropped rather
		// than rejected, so a client sending the whole row is not an error.
		allowed := map[string]bool{
{{PATCH_ALLOWED}}		}
		updates = map[string]interface{}{}
		for k, v := range req.Patch {
			if allowed[k] {
				updates[k] = v
			}
		}
		if len(updates) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "No writable fields in patch",
				},
			})
			return
		}
	}

	// One transaction: all of it lands or none of it does.
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		switch req.Action {
		case "delete":
			return tx.Where("id IN ?", ids).Delete(&models.{{Pascal}}{}).Error
		case "archive":
			now := time.Now()
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", ids).
				Update("archived_at", now).Error
		case "restore":
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", ids).
				Update("archived_at", nil).Error
		default:
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", ids).
				Updates(updates).Error
		}
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to " + req.Action + " {{plural}}",
			},
		})
		return
	}

	// One audit entry naming the action and the count, not N entries that bury
	// everything else somebody did today.
	// A local map, not a package-level helper: every resource gets its own
	// handler file in package handlers, so a shared func would be redeclared
	// once per resource.
	past := map[string]string{
		"delete":  "deleted",
		"archive": "archived",
		"restore": "restored",
		"patch":   "updated",
	}[req.Action]

	noun := "{{plural}}"
	if len(ids) == 1 {
		noun = "{{lower}}"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and joining
	// five hundred UUIDs into it makes the column unusable for the thing it is
	// for. The count lives in the summary, where it can be read.
	events.Emitted(c, "{{plural}}", "{{Pascal}}", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
`)

	// The transition endpoint, only for a resource whose status field is a
	// state machine. Appended rather than templated in, so an ordinary
	// resource's handler is byte-identical to what it was.
	content += g.workflowHandlerMethod(names)

	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+".go")
	return writeFileWithDirs(path, content)
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

	item, err := services.Transition` + names.Pascal + `(h.DB, c, id, action, can)
	if err != nil {
		switch err.(type) {
		case workflow.ErrInvalidTransition, workflow.ErrUnknownAction:
			// 422, not 400 and not 500. The request was well-formed and the
			// server is fine; the process does not allow this move, and the
			// message says which moves it does allow.
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
				"code":    "INVALID_TRANSITION",
				"message": err.Error(),
			}})
		default:
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
					"code": "NOT_FOUND", "message": "` + names.Pascal + ` not found",
				}})
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{
				"code": "FORBIDDEN", "message": err.Error(),
			}})
		}
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

	// If any field is a file/files, pull FileRefSchema from the
	// shared package so the generated schema references it directly
	// instead of inlining the FileRef shape (single source of truth).
	needsFileRef := false
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			needsFileRef = true
			break
		}
	}

	importLines := `import { z } from "zod";`
	if needsFileRef {
		importLines = "import { z } from \"zod\";\nimport { FileRefSchema } from \"./file-ref\";"
	}

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
	needsFileRef := false
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
		tsName := toSnakeCase(f.Name)
		tsType := f.TSType()
		fields += fmt.Sprintf("  %s: %s;\n", tsName, tsType)
	}
	if needsFileRef {
		imports += "import type { FileRef } from \"../schemas/file-ref\";\n"
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
	needsFileRef := false
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			needsFileRef = true
			break
		}
	}
	if needsFileRef {
		apiImport += "\nimport type { FileRef } from \"@repo/shared/schemas\";"
	}
	content := fmt.Sprintf(`import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
%s

interface %s {
  id: string;
%s  created_at: string;
  updated_at: string;
}

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
		names.Pascal,
		g.buildTSInterfaceFields(),
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
			typ := "text"
			extra := ""
			switch FieldType(cf.Type) {
			case FieldInt:
				typ, extra = "number", `, numberKind: "int"`
			case FieldUint:
				typ, extra = "number", `, numberKind: "uint"`
			case FieldFloat:
				typ, extra = "number", `, numberKind: "float"`
			case FieldDate, FieldDatetime:
				typ = "date"
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
		names.Pascal, names.PluralPascal,
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
	root := filepath.Join(g.Root, "apps", "admin", "resources")
	// One folder per resource. If this project is still flat, the definition
	// that was there is left where it is and the migration in grit upgrade
	// moves it: silently writing a second copy in a folder would leave two
	// definitions and a registry pointing at the stale one.
	if flat := filepath.Join(root, names.PluralKebab+".ts"); fileExists(flat) {
		if err := writeFileWithDirs(flat, g.resourceDefinitionFileContent(names)); err != nil {
			return err
		}
		return writeResourceCustomStub(root, names)
	}
	dir, path := scaffold.ResourceDefPath(root, names.PluralKebab)
	if err := writeFileWithDirs(path, g.resourceDefinitionFileContent(names)); err != nil {
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

	path := filepath.Join(g.Root, "apps", "admin", "app", "(dashboard)", "resources", names.PluralKebab, "page.tsx")
	return writeFileWithDirs(path, content)
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

	path := filepath.Join(g.Root, "apps", "admin", "app", "(dashboard)", "resources", names.PluralKebab, "[id]", "page.tsx")
	return writeFileWithDirs(path, content)
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
