package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	for _, f := range fields {
		if f.NeedsDatatypesImport() {
			needsDatatypes = true
			break
		}
	}

	// Build imports
	hasSlug := slugField != nil
	var imports string
	stdImports := `"time"`
	if hasSlug {
		stdImports = "\"fmt\"\n\t\"time\""
	}
	extImports := "\"github.com/google/uuid\"\n\t\"gorm.io/gorm\""
	if needsDatatypes {
		extImports = "\"github.com/google/uuid\"\n\t\"gorm.io/datatypes\"\n\t\"gorm.io/gorm\""
	}
	imports = fmt.Sprintf("import (\n\t%s\n\n\t%s\n)", stdImports, extImports)

	structFields := ""
	for _, f := range fields {
		// belongs_to: emit FK column + association struct
		if f.IsBelongsTo() {
			relModel := f.RelatedModelName()
			baseName := strings.TrimSuffix(f.Name, "_id") // strip _id if user included it
			fkGoName := toPascalCase(baseName) + "ID"     // e.g., CategoryID
			fkJson := toSnakeCase(baseName) + "_id"       // e.g., category_id
			assocName := toPascalCase(baseName)            // e.g., Category
			// FK column
			structFields += fmt.Sprintf("\t%s string `gorm:\"size:36;index\" json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
			// Association struct
			structFields += fmt.Sprintf("\t%s %s `gorm:\"foreignKey:%s\" json:\"%s\"`\n",
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
		if gormTag != "" {
			tags = fmt.Sprintf(`gorm:"%s" %s`, gormTag, tags)
		}

		if f.Required && (f.GoType() == "string") && !f.IsSlug() {
			tags += ` binding:"required"`
		}

		structFields += fmt.Sprintf("\t%s %s `%s`\n", goName, goType, tags)
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
}
`, imports, names.Pascal, names.Lower, names.Pascal, structFields)

	// Add BeforeCreate hook (UUID generation + optional slug)
	if hasSlug {
		slugGoName := toPascalCase(slugField.Name)
		content += fmt.Sprintf(`
// BeforeCreate generates a UUID and auto-generates the slug before inserting.
func (m *%s) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.%s == "" {
		m.%s = slugify(fmt.Sprintf("%%v", m.%s))
	}
	return nil
}
`, names.Pascal, slugGoName, slugGoName, slugSourceGo)

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
		// No slug — still need UUID generation
		content += fmt.Sprintf(`
// BeforeCreate generates a UUID before inserting.
func (m *%s) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
`, names.Pascal)
	}

	// BeforeUpdate increments Version on every server-side write so offline
	// clients can detect that a record they edited has moved on. Pair with
	// /api/sync/push for safe write replay.
	content += fmt.Sprintf(`
// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *%s) BeforeUpdate(tx *gorm.DB) error {
	m.Version++
	return nil
}
`, names.Pascal)

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
	if params.SortBy == "" {
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
			searchFields = append(searchFields, toSnakeCase(f.Name)+" ILIKE ?")
		}
	}
	if len(searchFields) == 0 {
		searchFields = []string{"id::text ILIKE ?"}
	}

	clause := strings.Join(searchFields, " OR ")
	args := ""
	for range searchFields {
		args += `, "%"+params.Search+"%"`
	}

	return `"` + clause + `"` + args
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

			createFields += fmt.Sprintf("\t\t%s string `json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
			createAssignments += fmt.Sprintf("\t\t%s: req.%s,\n", fkGoName, fkGoName)

			updateFields += fmt.Sprintf("\t\t%s *string `json:\"%s\"`\n", fkGoName, fkJson)
			updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", fkGoName, fkJson, fkGoName)
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

	// Build preload chain
	preloadChain := ""
	for _, p := range preloads {
		preloadChain += fmt.Sprintf(".Preload(\"%s\")", p)
	}

	// Build reload-with-preloads line (used after Create/Update)
	reloadLine := "\th.DB.First(&item, \"id = ?\", item.ID)"
	if preloadChain != "" {
		reloadLine = fmt.Sprintf("\th.DB%s.First(&item, \"id = ?\", item.ID)", preloadChain)
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
	for _, f := range g.Definition.Fields {
		if f.GoType() == "*time.Time" {
			needsTimeImport = true
		}
		if f.NeedsDatatypesImport() {
			needsHandlerDatatypes = true
		}
	}
	timeImport := ""
	if needsTimeImport {
		timeImport = "\n\t\"time\""
	}
	datatypesImport := ""
	if needsHandlerDatatypes {
		datatypesImport = "\n\t\"gorm.io/datatypes\""
	}

	r := strings.NewReplacer(
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{lower}}", names.Lower,
		"{{plural}}", names.Plural,
		"{{Plural}}", names.PluralPascal,
		"{{SORT_COLS}}", sortCols,
		"{{SEARCH_COLS}}", searchCols,
		"{{EXPORT_COLS}}", exportCols,
		"{{CREATE_FIELDS}}", createFields,
		"{{CREATE_ASSIGN}}", createAssignments,
		"{{UPDATE_FIELDS}}", updateFields,
		"{{UPDATE_MAP}}", updateMap,
		"{{PATCH_ALLOWED}}", patchAllowed,
		"{{PRELOADS}}", preloadChain,
		"{{M2M_CREATE}}", m2mCreateCode,
		"{{M2M_UPDATE}}", m2mUpdateCode,
		"{{RELOAD}}", reloadLine,
		"{{TIME_IMPORT}}", timeImport,
		"{{DATATYPES_IMPORT}}", datatypesImport,
	)

	content := r.Replace(`package handlers

import (
	"net/http"{{TIME_IMPORT}}

	"github.com/gin-gonic/gin"{{DATATYPES_IMPORT}}
	"gorm.io/gorm"

	"{{MODULE}}/internal/export"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
)

// {{Pascal}}Handler handles {{lower}} endpoints.
type {{Pascal}}Handler struct {
	DB *gorm.DB
}

// List returns a paginated list of {{plural}}.
func (h *{{Pascal}}Handler) List(c *gin.Context) {
	query := h.DB.Model(&models.{{Pascal}}{}){{PRELOADS}}

	res, err := paginate.List[models.{{Pascal}}](
		query,
		paginate.Bind(c),
		paginate.Config{
			Searchable: []string{{{SEARCH_COLS}}},
			Sortable:   map[string]bool{{{SORT_COLS}}},
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
			clause += col + " ILIKE ?"
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

// Create adds a new {{lower}}.
func (h *{{Pascal}}Handler) Create(c *gin.Context) {
	var req struct {
{{CREATE_FIELDS}}	}

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

	if err := h.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create {{lower}}",
			},
		})
		return
	}
{{M2M_CREATE}}
{{RELOAD}}

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

	var req struct {
{{UPDATE_FIELDS}}	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	updates := map[string]interface{}{}
{{UPDATE_MAP}}
	if err := h.DB.Model(&item).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update {{lower}}",
			},
		})
		return
	}
{{M2M_UPDATE}}
{{RELOAD}}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "{{Pascal}} deleted successfully",
	})
}
`)

	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// buildHandlerSearchCols returns the comma-separated quoted column names for
// the paginate.Config.Searchable slice literal. Only text-like field types are
// included — FK UUID columns (which happen to be Go string) are skipped so
// search doesn't ILIKE against opaque identifiers.
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
			createFields += fmt.Sprintf("  %s: %s,\n", fkName, f.ZodType())
			updateFields += fmt.Sprintf("  %s: %s,\n", fkName, f.ZodType()+".optional()")
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

	content := fmt.Sprintf(`import { z } from "zod";

export const Create%sSchema = z.object({
%s});

export const Update%sSchema = z.object({
%s});

export type Create%sInput = z.infer<typeof Create%sSchema>;
export type Update%sInput = z.infer<typeof Update%sSchema>;
`, names.Pascal, createFields, names.Pascal, updateFields,
		names.Pascal, names.Pascal, names.Pascal, names.Pascal)

	path := filepath.Join(g.Root, "packages", "shared", "schemas", names.Kebab+".ts")
	return writeFileWithDirs(path, content)
}

// writeTSTypes creates the TypeScript type file for the resource.
func (g *Generator) writeTSTypes(names Names) error {
	// Collect relationship imports
	imports := ""
	fields := ""
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			relModel := f.RelatedModelName()
			relKebab := strings.ReplaceAll(toSnakeCase(relModel), "_", "-")
			imports += fmt.Sprintf("import type { %s } from \"./%s\";\n", relModel, relKebab)
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkSnake := toSnakeCase(baseName) + "_id"
			// FK matches the referenced model's UUID string PK.
			fields += fmt.Sprintf("  %s: string;\n", fkSnake)
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
		tsName := toSnakeCase(f.Name)
		tsType := f.TSType()
		fields += fmt.Sprintf("  %s: %s;\n", tsName, tsType)
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
func (g *Generator) writeReactQueryHooks(names Names, app string) error {
	content := fmt.Sprintf(`import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

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
func (g *Generator) writeResourceDefinition(names Names) error {
	icon := guessLucideIcon(names.Pascal)

	// Build column definitions. ID is intentionally NOT listed by default —
	// UUIDs are noisy and rarely something an operator scans by eye.
	// Users who want it can add { key: "id", label: "ID", width: "80px" }
	// to the columns array by hand.
	columns := ""
	for _, f := range g.Definition.Fields {
		// belongs_to: show related model's name via dot notation
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			assocSnake := toSnakeCase(baseName)
			colLabel := strings.Join(splitPascal(toPascalCase(baseName)), " ")
			columns += fmt.Sprintf("\n    { key: \"%s.name\", label: \"%s\" },", assocSnake, colLabel)
			continue
		}
		// many_to_many: skip from table columns (arrays are noisy)
		if f.IsManyToMany() {
			continue
		}

		colName := toSnakeCase(f.Name)
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

		columns += "\n    { " + strings.Join(parts, ", ") + " },"
	}
	columns += fmt.Sprintf(`
    { key: "created_at", label: "Created", sortable: true, format: "relative" },`)

	// Build form field definitions (skip slug — auto-generated, not editable)
	formFields := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}

		// belongs_to: relationship-select with endpoint
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkKey := toSnakeCase(baseName) + "_id"
			fieldLabel := strings.Join(splitPascal(toPascalCase(baseName)), " ")
			relSnake := toSnakeCase(f.RelatedModelName())
			relPlural := Pluralize(relSnake)
			formFields += fmt.Sprintf("\n    { key: \"%s\", label: \"%s\", type: \"relationship-select\", required: true, relatedEndpoint: \"/api/%s\", displayField: \"name\" },",
				fkKey, fieldLabel, relPlural)
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

	content := fmt.Sprintf(`import { defineResource } from "@/lib/resource";

export const %sResource = defineResource({
  name: "%s",
  slug: "%s",
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
});
`,
		names.Camel,
		names.Pascal,
		names.PluralKebab,
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

	path := filepath.Join(g.Root, "apps", "admin", "resources", names.PluralKebab+".ts")
	return writeFileWithDirs(path, content)
}

// writeResourcePage creates a thin admin page wrapper for the resource.
func (g *Generator) writeResourcePage(names Names) error {
	content := fmt.Sprintf(`"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { %sResource } from "@/resources/%s";

export default function %sPage() {
  return <ResourcePage resource={%sResource} />;
}
`,
		names.Camel, names.PluralKebab,
		names.PluralPascal,
		names.Camel,
	)

	path := filepath.Join(g.Root, "apps", "admin", "app", "(dashboard)", "resources", names.PluralKebab, "page.tsx")
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
func splitPascal(s string) []string {
	var words []string
	start := 0
	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			words = append(words, s[start:i])
			start = i
		}
	}
	words = append(words, s[start:])
	return words
}
