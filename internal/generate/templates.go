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
%s	CreatedAt time.Time      `+"`"+`json:"created_at"`+"`"+`
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
func (s *{{Pascal}}Service) GetByID(id uint) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.DB.First(&item, id).Error; err != nil {
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

// Update modifies an existing {{lower}}.
func (s *{{Pascal}}Service) Update(id uint, updates map[string]interface{}) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.DB.First(&item, id).Error; err != nil {
		return nil, fmt.Errorf("{{lower}} not found: %w", err)
	}

	if err := s.DB.Model(&item).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("updating {{lower}}: %w", err)
	}

	s.DB.First(&item, id)
	return &item, nil
}

// Delete soft-deletes a {{lower}}.
func (s *{{Pascal}}Service) Delete(id uint) error {
	var item models.{{Pascal}}
	if err := s.DB.First(&item, id).Error; err != nil {
		return fmt.Errorf("{{lower}} not found: %w", err)
	}
	if err := s.DB.Delete(&item).Error; err != nil {
		return fmt.Errorf("deleting {{lower}}: %w", err)
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

	// Collect preloads for relationships
	var preloads []string

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}

		// belongs_to: add FK field (uint) to create/update
		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkGoName := toPascalCase(baseName) + "ID"
			fkJson := toSnakeCase(baseName) + "_id"
			assocName := toPascalCase(baseName)
			preloads = append(preloads, assocName)

			createFields += fmt.Sprintf("\t\t%s uint `json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
			createAssignments += fmt.Sprintf("\t\t%s: req.%s,\n", fkGoName, fkGoName)

			updateFields += fmt.Sprintf("\t\t%s *uint `json:\"%s\"`\n", fkGoName, fkJson)
			updateMap += fmt.Sprintf("	if req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", fkGoName, fkJson, fkGoName)
			continue
		}

		// many_to_many: add []uint for create, *[]uint for update
		if f.IsManyToMany() {
			relModel := f.RelatedModelName()
			assocName := toPascalCase(f.Name)
			idsName := toPascalCase(f.Name) + "IDs"
			idsJson := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			preloads = append(preloads, assocName)

			createFields += fmt.Sprintf("\t\t%s []uint `json:\"%s\"`\n", idsName, idsJson)
			updateFields += fmt.Sprintf("\t\t%s *[]uint `json:\"%s\"`\n", idsName, idsJson)

			varName := toSnakeCase(f.Name)
			m2mCreateCode += fmt.Sprintf("\n\tif len(req.%s) > 0 {\n\t\tvar %s []models.%s\n\t\th.DB.Find(&%s, req.%s)\n\t\th.DB.Model(&item).Association(\"%s\").Replace(%s)\n\t}\n", idsName, varName, relModel, varName, idsName, assocName, varName)
			m2mUpdateCode += fmt.Sprintf("\n\tif req.%s != nil {\n\t\tvar %s []models.%s\n\t\tif len(*req.%s) > 0 {\n\t\t\th.DB.Find(&%s, *req.%s)\n\t\t}\n\t\th.DB.Model(&item).Association(\"%s\").Replace(%s)\n\t}\n", idsName, varName, relModel, idsName, varName, idsName, assocName, varName)
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
	reloadLine := "\th.DB.First(&item, item.ID)"
	if preloadChain != "" {
		reloadLine = fmt.Sprintf("\th.DB%s.First(&item, item.ID)", preloadChain)
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

	searchArgs := g.buildHandlerSearchArgs(names)

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
		"{{SEARCH_ARGS}}", searchArgs,
		"{{CREATE_FIELDS}}", createFields,
		"{{CREATE_ASSIGN}}", createAssignments,
		"{{UPDATE_FIELDS}}", updateFields,
		"{{UPDATE_MAP}}", updateMap,
		"{{PRELOADS}}", preloadChain,
		"{{M2M_CREATE}}", m2mCreateCode,
		"{{M2M_UPDATE}}", m2mUpdateCode,
		"{{RELOAD}}", reloadLine,
		"{{TIME_IMPORT}}", timeImport,
		"{{DATATYPES_IMPORT}}", datatypesImport,
	)

	content := r.Replace(`package handlers

import (
	"math"
	"net/http"
	"strconv"{{TIME_IMPORT}}

	"github.com/gin-gonic/gin"{{DATATYPES_IMPORT}}
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// {{Pascal}}Handler handles {{lower}} endpoints.
type {{Pascal}}Handler struct {
	DB *gorm.DB
}

// List returns a paginated list of {{plural}}.
func (h *{{Pascal}}Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	allowedSorts := map[string]bool{{{SORT_COLS}}}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}

	query := h.DB.Model(&models.{{Pascal}}{}){{PRELOADS}}

	if search != "" {
		query = query.Where({{SEARCH_ARGS}})
	}

	var total int64
	query.Count(&total)

	var items []models.{{Pascal}}
	offset := (page - 1) * pageSize
	if err := query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch {{plural}}",
			},
		})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// GetByID returns a single {{lower}} by ID.
func (h *{{Pascal}}Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB{{PRELOADS}}.First(&item, id).Error; err != nil {
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
	if err := h.DB.First(&item, id).Error; err != nil {
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

// Delete soft-deletes a {{lower}}.
func (h *{{Pascal}}Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	var item models.{{Pascal}}
	if err := h.DB.First(&item, id).Error; err != nil {
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

// buildHandlerSearchArgs builds the search Where clause for the handler.
func (g *Generator) buildHandlerSearchArgs(names Names) string {
	var searchFields []string
	for _, f := range g.Definition.Fields {
		if f.GoType() == "string" {
			searchFields = append(searchFields, toSnakeCase(f.Name)+" ILIKE ?")
		}
	}
	if len(searchFields) == 0 {
		return `"id::text ILIKE ?", "%"+search+"%"`
	}

	clause := `"` + strings.Join(searchFields, " OR ") + `"`
	args := ""
	for range searchFields {
		args += `, "%"+search+"%"`
	}

	return clause + args
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
			fkName := toCamelCase(f.FKColumnName())
			createFields += fmt.Sprintf("  %s: %s,\n", fkName, f.ZodType())
			updateFields += fmt.Sprintf("  %s: %s,\n", fkName, f.ZodType()+".optional()")
			continue
		}

		// many_to_many: use <name>_ids (e.g., tag_ids)
		if f.IsManyToMany() {
			idsName := strings.TrimSuffix(toCamelCase(f.Name), "s") + "Ids"
			createFields += fmt.Sprintf("  %s: %s,\n", idsName, f.ZodType())
			updateFields += fmt.Sprintf("  %s: %s,\n", idsName, f.ZodType())
			continue
		}

		camelName := toCamelCase(f.Name)
		zodType := f.ZodType()
		createFields += fmt.Sprintf("  %s: %s,\n", camelName, zodType)

		// Update schema: make all fields optional
		updateZod := f.ZodType()
		if !strings.Contains(updateZod, ".optional()") && !strings.Contains(updateZod, ".nullable()") {
			updateZod += ".optional()"
		}
		updateFields += fmt.Sprintf("  %s: %s,\n", camelName, updateZod)
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
			fields += fmt.Sprintf("  %s: number;\n", fkSnake)
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
			result += fmt.Sprintf("  %s: number;\n", fkSnake)
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

	// Build column definitions
	columns := fmt.Sprintf(`    { key: "id", label: "ID", sortable: true, width: "80px" },`)
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
%s
    ],
    filters: [%s
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    searchable: true,
    pageSize: 20,
  },
  form: {
    fields: [%s
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
