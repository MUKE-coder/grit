package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeDesktopModel generates internal/models/<snake>.go with a GORM model struct.
func (g *DesktopGenerator) writeDesktopModel(names Names) error {
	fields := g.Definition.Fields

	// Detect slug field and resolve source
	var slugField *Field
	for i, f := range fields {
		if f.IsSlug() {
			slugField = &fields[i]
			break
		}
	}

	slugSourceGo := ""
	if slugField != nil {
		if slugField.SlugSource != "" {
			slugSourceGo = toPascalCase(slugField.SlugSource)
		} else {
			for _, f := range fields {
				if FieldType(f.Type) == FieldString {
					slugSourceGo = toPascalCase(f.Name)
					break
				}
			}
			if slugSourceGo == "" {
				slugSourceGo = "ID"
			}
		}
	}

	// Check if we need time import
	needsTime := false
	for _, f := range fields {
		if f.NeedsTimeImport() {
			needsTime = true
			break
		}
	}

	hasSlug := slugField != nil

	// Build imports
	stdImports := ""
	if hasSlug && needsTime {
		stdImports = "\t\"fmt\"\n\t\"time\""
	} else if hasSlug {
		stdImports = "\t\"fmt\""
	} else if needsTime {
		stdImports = "\t\"time\""
	}

	extImports := "\t\"gorm.io/gorm\""

	imports := "import (\n"
	if stdImports != "" {
		imports += stdImports + "\n\n"
	}
	imports += extImports + "\n)"

	// Build struct fields
	structFields := ""
	for _, f := range fields {
		goName := toPascalCase(f.Name)
		goType := f.GoType()
		jsonTag := toSnakeCase(f.Name)
		gormTag := f.GORMTag()

		tagParts := ""
		if gormTag != "" {
			tagParts = "gorm:\"" + gormTag + "\" json:\"" + jsonTag + "\""
		} else {
			tagParts = "json:\"" + jsonTag + "\""
		}

		structFields += "\t" + goName + " " + goType + " " + "`" + tagParts + "`" + "\n"
	}

	// Build the full file content
	content := "package models\n\n"
	content += imports + "\n\n"
	content += "// " + names.Pascal + " represents a " + names.Lower + " in the system.\n"
	content += "type " + names.Pascal + " struct {\n"
	content += "\tID        string         " + "`" + "gorm:\"primarykey;size:36\" json:\"id\"" + "`" + "\n"
	content += structFields
	content += "\tCreatedAt time.Time      " + "`" + "json:\"created_at\"" + "`" + "\n"
	content += "\tUpdatedAt time.Time      " + "`" + "json:\"updated_at\"" + "`" + "\n"
	content += "\tDeletedAt gorm.DeletedAt " + "`" + "gorm:\"index\" json:\"-\"" + "`" + "\n"
	content += "}\n"

	// Add BeforeCreate hook for slug
	if hasSlug {
		slugGoName := toPascalCase(slugField.Name)
		content += "\n// BeforeCreate auto-generates the slug before inserting.\n"
		content += "func (m *" + names.Pascal + ") BeforeCreate(tx *gorm.DB) error {\n"
		content += "\tif m." + slugGoName + " == \"\" {\n"
		content += "\t\tm." + slugGoName + " = slugify(fmt.Sprintf(\"%v\", m." + slugSourceGo + "))\n"
		content += "\t}\n"
		content += "\treturn nil\n"
		content += "}\n"
	}

	// Also need time import even without explicit date fields because CreatedAt/UpdatedAt use time.Time
	// Re-check: we always need time because of CreatedAt/UpdatedAt
	if !needsTime && !hasSlug {
		// Need to add time import since CreatedAt/UpdatedAt require it
		imports = "import (\n\t\"time\"\n\n\t\"gorm.io/gorm\"\n)"
		content = "package models\n\n"
		content += imports + "\n\n"
		content += "// " + names.Pascal + " represents a " + names.Lower + " in the system.\n"
		content += "type " + names.Pascal + " struct {\n"
		content += "\tID        string         " + "`" + "gorm:\"primarykey;size:36\" json:\"id\"" + "`" + "\n"
		content += structFields
		content += "\tCreatedAt time.Time      " + "`" + "json:\"created_at\"" + "`" + "\n"
		content += "\tUpdatedAt time.Time      " + "`" + "json:\"updated_at\"" + "`" + "\n"
		content += "\tDeletedAt gorm.DeletedAt " + "`" + "gorm:\"index\" json:\"-\"" + "`" + "\n"
		content += "}\n"
	} else if !needsTime && hasSlug {
		// slug needs fmt but also need time for CreatedAt/UpdatedAt
		imports = "import (\n\t\"fmt\"\n\t\"time\"\n\n\t\"gorm.io/gorm\"\n)"
		tmpContent := "package models\n\n"
		tmpContent += imports + "\n\n"
		tmpContent += "// " + names.Pascal + " represents a " + names.Lower + " in the system.\n"
		tmpContent += "type " + names.Pascal + " struct {\n"
		tmpContent += "\tID        string         " + "`" + "gorm:\"primarykey;size:36\" json:\"id\"" + "`" + "\n"
		tmpContent += structFields
		tmpContent += "\tCreatedAt time.Time      " + "`" + "json:\"created_at\"" + "`" + "\n"
		tmpContent += "\tUpdatedAt time.Time      " + "`" + "json:\"updated_at\"" + "`" + "\n"
		tmpContent += "\tDeletedAt gorm.DeletedAt " + "`" + "gorm:\"index\" json:\"-\"" + "`" + "\n"
		tmpContent += "}\n"

		slugGoName := toPascalCase(slugField.Name)
		tmpContent += "\n// BeforeCreate auto-generates the slug before inserting.\n"
		tmpContent += "func (m *" + names.Pascal + ") BeforeCreate(tx *gorm.DB) error {\n"
		tmpContent += "\tif m." + slugGoName + " == \"\" {\n"
		tmpContent += "\t\tm." + slugGoName + " = slugify(fmt.Sprintf(\"%v\", m." + slugSourceGo + "))\n"
		tmpContent += "\t}\n"
		tmpContent += "\treturn nil\n"
		tmpContent += "}\n"

		content = tmpContent
	}

	content = strings.ReplaceAll(content, "<MODULE>", g.Module)

	path := filepath.Join(g.Root, "internal", "models", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// writeDesktopService generates internal/service/<snake>.go with CRUD methods.
func (g *DesktopGenerator) writeDesktopService(names Names) error {
	// Build search WHERE clause
	searchWhere := g.buildDesktopSearchWhere()

	// Build input-to-model field assignments for Create
	createAssignments := ""
	updateAssignments := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		goName := toPascalCase(f.Name)
		createAssignments += "\t\t" + goName + ": input." + goName + ",\n"
		updateAssignments += "\titem." + goName + " = input." + goName + "\n"
	}

	// Build the input struct string (for reference, inject handles writing to types.go)
	inputFields := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		goName := toPascalCase(f.Name)
		goType := f.GoType()
		jsonTag := toSnakeCase(f.Name)
		inputFields += "\t" + goName + " " + goType + " " + "`" + "json:\"" + jsonTag + "\"" + "`" + "\n"
	}

	content := "package service\n\n"
	content += "import (\n"
	content += "\t\"fmt\"\n"
	content += "\t\"math\"\n"
	content += "\n"
	content += "\t\"<MODULE>/internal/models\"\n"
	content += "\t\"gorm.io/gorm\"\n"
	content += ")\n\n"

	// Service struct
	content += "type " + names.Pascal + "Service struct {\n"
	content += "\tdb *gorm.DB\n"
	content += "}\n\n"

	// Constructor
	content += "func New" + names.Pascal + "Service(db *gorm.DB) *" + names.Pascal + "Service {\n"
	content += "\treturn &" + names.Pascal + "Service{db: db}\n"
	content += "}\n\n"

	// List with pagination + search
	content += "func (s *" + names.Pascal + "Service) List(page, pageSize int, search string) (*models.PaginatedResult, error) {\n"
	content += "\tif page < 1 {\n"
	content += "\t\tpage = 1\n"
	content += "\t}\n"
	content += "\tif pageSize < 1 {\n"
	content += "\t\tpageSize = 10\n"
	content += "\t}\n\n"
	content += "\tvar total int64\n"
	content += "\tquery := s.db.Model(&models." + names.Pascal + "{})\n"
	content += "\tif search != \"\" {\n"
	content += "\t\tquery = query.Where(" + searchWhere + ")\n"
	content += "\t}\n"
	content += "\tquery.Count(&total)\n\n"
	content += "\tvar items []models." + names.Pascal + "\n"
	content += "\toffset := (page - 1) * pageSize\n"
	content += "\tif err := query.Order(\"created_at DESC\").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {\n"
	content += "\t\treturn nil, fmt.Errorf(\"failed to list " + names.Plural + ": %w\", err)\n"
	content += "\t}\n\n"
	content += "\tpages := int(math.Ceil(float64(total) / float64(pageSize)))\n"
	content += "\treturn &models.PaginatedResult{\n"
	content += "\t\tData:     items,\n"
	content += "\t\tTotal:    total,\n"
	content += "\t\tPage:     page,\n"
	content += "\t\tPageSize: pageSize,\n"
	content += "\t\tPages:    pages,\n"
	content += "\t}, nil\n"
	content += "}\n\n"

	// ListAll
	content += "func (s *" + names.Pascal + "Service) ListAll() ([]models." + names.Pascal + ", error) {\n"
	content += "\tvar items []models." + names.Pascal + "\n"
	content += "\tif err := s.db.Order(\"created_at DESC\").Find(&items).Error; err != nil {\n"
	content += "\t\treturn nil, err\n"
	content += "\t}\n"
	content += "\treturn items, nil\n"
	content += "}\n\n"

	// GetByID
	content += "func (s *" + names.Pascal + "Service) GetByID(id uint) (*models." + names.Pascal + ", error) {\n"
	content += "\tvar item models." + names.Pascal + "\n"
	content += "\tif err := s.db.First(&item, id).Error; err != nil {\n"
	content += "\t\treturn nil, fmt.Errorf(\"" + names.Lower + " not found\")\n"
	content += "\t}\n"
	content += "\treturn &item, nil\n"
	content += "}\n\n"

	// Create
	content += "func (s *" + names.Pascal + "Service) Create(input models." + names.Pascal + "Input) (*models." + names.Pascal + ", error) {\n"
	content += "\titem := models." + names.Pascal + "{\n"
	content += createAssignments
	content += "\t}\n"
	content += "\tif err := s.db.Create(&item).Error; err != nil {\n"
	content += "\t\treturn nil, fmt.Errorf(\"failed to create " + names.Lower + ": %w\", err)\n"
	content += "\t}\n"
	content += "\treturn &item, nil\n"
	content += "}\n\n"

	// Update
	content += "func (s *" + names.Pascal + "Service) Update(id uint, input models." + names.Pascal + "Input) (*models." + names.Pascal + ", error) {\n"
	content += "\tvar item models." + names.Pascal + "\n"
	content += "\tif err := s.db.First(&item, id).Error; err != nil {\n"
	content += "\t\treturn nil, fmt.Errorf(\"" + names.Lower + " not found\")\n"
	content += "\t}\n"
	content += updateAssignments
	content += "\tif err := s.db.Save(&item).Error; err != nil {\n"
	content += "\t\treturn nil, fmt.Errorf(\"failed to update " + names.Lower + ": %w\", err)\n"
	content += "\t}\n"
	content += "\treturn &item, nil\n"
	content += "}\n\n"

	// Delete
	content += "func (s *" + names.Pascal + "Service) Delete(id uint) error {\n"
	content += "\tif err := s.db.Delete(&models." + names.Pascal + "{}, id).Error; err != nil {\n"
	content += "\t\treturn fmt.Errorf(\"failed to delete " + names.Lower + ": %w\", err)\n"
	content += "\t}\n"
	content += "\treturn nil\n"
	content += "}\n"

	content = strings.ReplaceAll(content, "<MODULE>", g.Module)

	path := filepath.Join(g.Root, "internal", "service", names.Snake+".go")
	return writeFileWithDirs(path, content)
}

// buildDesktopSearchWhere creates the Where clause arguments for the service search.
func (g *DesktopGenerator) buildDesktopSearchWhere() string {
	var searchFields []string
	for _, f := range g.Definition.Fields {
		if f.IsSearchable() {
			searchFields = append(searchFields, toSnakeCase(f.Name)+" LIKE ?")
		}
	}
	if len(searchFields) == 0 {
		return "\"id = ?\", search"
	}

	clause := strings.Join(searchFields, " OR ")
	args := ""
	for range searchFields {
		args += ", \"%\"+search+\"%\""
	}

	return "\"" + clause + "\"" + args
}

// writeDesktopListRoute generates frontend/src/routes/_layout/<plural>.index.tsx.
func (g *DesktopGenerator) writeDesktopListRoute(names Names) error {
	// Build table header columns (skip slug fields)
	headerCols := ""
	colCount := 0
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		label := toPascalCase(f.Name)
		// Add spaces between pascal words
		labelParts := splitPascal(label)
		label = strings.Join(labelParts, " ")
		headerCols += "              <th className=\"text-left text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3\">" + label + "</th>\n"
		colCount++
	}
	// Add Actions column
	colCount += 1 // for actions

	totalColSpan := fmt.Sprintf("%d", colCount)

	// Build table body cells
	bodyCells := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		jsonName := toSnakeCase(f.Name)
		format := f.ColumnFormat()

		switch format {
		case "boolean":
			bodyCells += "                  <td className=\"px-4 py-3\">\n"
			bodyCells += "                    <span className={\"inline-flex items-center px-2 py-0.5 rounded text-xs font-medium \" + (item." + jsonName + " ? \"bg-success/10 text-success\" : \"bg-text-muted/10 text-text-muted\")}>\n"
			bodyCells += "                      {item." + jsonName + " ? \"Yes\" : \"No\"}\n"
			bodyCells += "                    </span>\n"
			bodyCells += "                  </td>\n"
		case "relative":
			bodyCells += "                  <td className=\"px-4 py-3 text-sm text-text-secondary\">\n"
			bodyCells += "                    {item." + jsonName + " ? new Date(item." + jsonName + ").toLocaleDateString(\"en-US\", { year: \"numeric\", month: \"short\", day: \"numeric\" }) : \"-\"}\n"
			bodyCells += "                  </td>\n"
		default:
			bodyCells += "                  <td className=\"px-4 py-3 text-sm text-text-secondary\">{item." + jsonName + " || \"-\"}</td>\n"
		}
	}

	// Determine the display field for delete confirmation and first display column
	firstStringField := ""
	for _, f := range g.Definition.Fields {
		if f.IsSearchable() && !f.IsSlug() {
			firstStringField = toSnakeCase(f.Name)
			break
		}
	}
	if firstStringField == "" {
		firstStringField = "id"
	}

	// Build Wails import names
	getFunc := "Get" + names.PluralPascal
	deleteFunc := "Delete" + names.Pascal
	exportPDFFunc := "Export" + names.PluralPascal + "PDF"
	exportExcelFunc := "Export" + names.PluralPascal + "Excel"

	content := "import { createFileRoute, useNavigate } from \"@tanstack/react-router\";\n"
	content += "import { useState } from \"react\";\n"
	content += "import { useQuery, useMutation, useQueryClient } from \"@tanstack/react-query\";\n"
	content += "import { Plus, Pencil, Trash2, FileDown, FileSpreadsheet, Search } from \"lucide-react\";\n"
	content += "import { toast } from \"sonner\";\n"
	content += "// @ts-ignore\n"
	content += "import { " + getFunc + ", " + deleteFunc + ", " + exportPDFFunc + ", " + exportExcelFunc + " } from \"../../../wailsjs/go/main/App\";\n"
	content += "\n"
	content += "export const Route = createFileRoute(\"/_layout/" + names.Plural + "/\")({\n"
	content += "  component: " + names.PluralPascal + "ListPage,\n"
	content += "});\n"
	content += "\n"
	content += "function " + names.PluralPascal + "ListPage() {\n"
	content += "  const navigate = useNavigate();\n"
	content += "  const queryClient = useQueryClient();\n"
	content += "  const [page, setPage] = useState(1);\n"
	content += "  const [search, setSearch] = useState(\"\");\n"
	content += "  const pageSize = 10;\n"
	content += "\n"
	content += "  const { data, isLoading } = useQuery({\n"
	content += "    queryKey: [\"" + names.Plural + "\", page, search],\n"
	content += "    queryFn: () => " + getFunc + "(page, pageSize, search),\n"
	content += "  });\n"
	content += "\n"
	content += "  const deleteMutation = useMutation({\n"
	content += "    mutationFn: (id: number) => " + deleteFunc + "(id),\n"
	content += "    onSuccess: () => {\n"
	content += "      toast.success(\"" + names.Pascal + " deleted\");\n"
	content += "      queryClient.invalidateQueries({ queryKey: [\"" + names.Plural + "\"] });\n"
	content += "    },\n"
	content += "    onError: (err: any) => toast.error(err?.message || \"Failed to delete\"),\n"
	content += "  });\n"
	content += "\n"
	content += "  const handleDelete = (id: number, label: string) => {\n"
	content += "    if (window.confirm(\"Delete \\\"\" + label + \"\\\"? This cannot be undone.\")) {\n"
	content += "      deleteMutation.mutate(id);\n"
	content += "    }\n"
	content += "  };\n"
	content += "\n"
	content += "  const handleExportPDF = async () => {\n"
	content += "    try {\n"
	content += "      const path = await " + exportPDFFunc + "();\n"
	content += "      toast.success(\"PDF exported to \" + path);\n"
	content += "    } catch (err: any) {\n"
	content += "      toast.error(err?.message || \"Export failed\");\n"
	content += "    }\n"
	content += "  };\n"
	content += "\n"
	content += "  const handleExportExcel = async () => {\n"
	content += "    try {\n"
	content += "      const path = await " + exportExcelFunc + "();\n"
	content += "      toast.success(\"Excel exported to \" + path);\n"
	content += "    } catch (err: any) {\n"
	content += "      toast.error(err?.message || \"Export failed\");\n"
	content += "    }\n"
	content += "  };\n"
	content += "\n"
	content += "  const items = data?.data || [];\n"
	content += "  const total = data?.total || 0;\n"
	content += "  const totalPages = data?.pages || 1;\n"
	content += "\n"
	content += "  return (\n"
	content += "    <div>\n"
	content += "      <div className=\"flex items-center justify-between mb-6\">\n"
	content += "        <div>\n"
	content += "          <h1 className=\"text-2xl font-bold text-foreground\">" + names.PluralPascal + "</h1>\n"
	content += "          <p className=\"text-sm text-text-secondary\">{total + \" total " + names.Lower + "\" + (total !== 1 ? \"s\" : \"\")}</p>\n"
	content += "        </div>\n"
	content += "        <div className=\"flex items-center gap-2\">\n"
	content += "          <button\n"
	content += "            onClick={handleExportPDF}\n"
	content += "            className=\"flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors\"\n"
	content += "          >\n"
	content += "            <FileDown size={16} />\n"
	content += "            PDF\n"
	content += "          </button>\n"
	content += "          <button\n"
	content += "            onClick={handleExportExcel}\n"
	content += "            className=\"flex items-center gap-2 px-3 py-2 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors\"\n"
	content += "          >\n"
	content += "            <FileSpreadsheet size={16} />\n"
	content += "            Excel\n"
	content += "          </button>\n"
	content += "          <button\n"
	content += "            onClick={() => navigate({ to: \"/" + names.Plural + "/new\" })}\n"
	content += "            className=\"flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors\"\n"
	content += "          >\n"
	content += "            <Plus size={16} />\n"
	content += "            New " + names.Pascal + "\n"
	content += "          </button>\n"
	content += "        </div>\n"
	content += "      </div>\n"
	content += "\n"
	content += "      <div className=\"mb-4\">\n"
	content += "        <div className=\"relative\">\n"
	content += "          <Search size={16} className=\"absolute left-3 top-1/2 -translate-y-1/2 text-text-muted\" />\n"
	content += "          <input\n"
	content += "            type=\"text\"\n"
	content += "            value={search}\n"
	content += "            onChange={(e) => { setSearch(e.target.value); setPage(1); }}\n"
	content += "            placeholder=\"Search " + names.Plural + "...\"\n"
	content += "            className=\"w-full pl-9 pr-3 py-2 bg-bg-elevated border border-border rounded-lg text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
	content += "          />\n"
	content += "        </div>\n"
	content += "      </div>\n"
	content += "\n"
	content += "      <div className=\"bg-bg-elevated border border-border rounded-xl overflow-hidden\">\n"
	content += "        <table className=\"w-full\">\n"
	content += "          <thead>\n"
	content += "            <tr className=\"border-b border-border\">\n"
	content += headerCols
	content += "              <th className=\"text-right text-xs font-medium text-text-muted uppercase tracking-wider px-4 py-3\">Actions</th>\n"
	content += "            </tr>\n"
	content += "          </thead>\n"
	content += "          <tbody>\n"
	content += "            {isLoading ? (\n"
	content += "              <tr>\n"
	content += "                <td colSpan={" + totalColSpan + "} className=\"px-4 py-8 text-center text-text-muted\">Loading...</td>\n"
	content += "              </tr>\n"
	content += "            ) : items.length === 0 ? (\n"
	content += "              <tr>\n"
	content += "                <td colSpan={" + totalColSpan + "} className=\"px-4 py-8 text-center text-text-muted\">No " + names.Plural + " found</td>\n"
	content += "              </tr>\n"
	content += "            ) : (\n"
	content += "              items.map((item: any) => (\n"
	content += "                <tr key={item.id} className=\"border-b border-border last:border-0 hover:bg-bg-hover/50 transition-colors\">\n"
	content += bodyCells
	content += "                  <td className=\"px-4 py-3 text-right\">\n"
	content += "                    <div className=\"flex items-center justify-end gap-1\">\n"
	content += "                      <button\n"
	content += "                        onClick={() => navigate({ to: \"/" + names.Plural + "/$id/edit\", params: { id: String(item.id) } })}\n"
	content += "                        className=\"p-1.5 rounded text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors\"\n"
	content += "                      >\n"
	content += "                        <Pencil size={14} />\n"
	content += "                      </button>\n"
	content += "                      <button\n"
	content += "                        onClick={() => handleDelete(item.id, item." + firstStringField + " || String(item.id))}\n"
	content += "                        className=\"p-1.5 rounded text-text-secondary hover:text-danger hover:bg-danger/10 transition-colors\"\n"
	content += "                      >\n"
	content += "                        <Trash2 size={14} />\n"
	content += "                      </button>\n"
	content += "                    </div>\n"
	content += "                  </td>\n"
	content += "                </tr>\n"
	content += "              ))\n"
	content += "            )}\n"
	content += "          </tbody>\n"
	content += "        </table>\n"
	content += "      </div>\n"
	content += "\n"
	content += "      {totalPages > 1 && (\n"
	content += "        <div className=\"flex items-center justify-between mt-4\">\n"
	content += "          <p className=\"text-sm text-text-muted\">\n"
	content += "            {\"Page \" + page + \" of \" + totalPages}\n"
	content += "          </p>\n"
	content += "          <div className=\"flex items-center gap-2\">\n"
	content += "            <button\n"
	content += "              onClick={() => setPage((p) => Math.max(1, p - 1))}\n"
	content += "              disabled={page <= 1}\n"
	content += "              className=\"px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors\"\n"
	content += "            >\n"
	content += "              Previous\n"
	content += "            </button>\n"
	content += "            <button\n"
	content += "              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}\n"
	content += "              disabled={page >= totalPages}\n"
	content += "              className=\"px-3 py-1.5 text-sm text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors\"\n"
	content += "            >\n"
	content += "              Next\n"
	content += "            </button>\n"
	content += "          </div>\n"
	content += "        </div>\n"
	content += "      )}\n"
	content += "    </div>\n"
	content += "  );\n"
	content += "}\n"

	path := filepath.Join(g.Root, "frontend", "src", "routes", "_layout", names.Plural+".index.tsx")
	return writeFileWithDirs(path, content)
}

// writeDesktopNewRoute generates frontend/src/routes/_layout/<plural>.new.tsx.
func (g *DesktopGenerator) writeDesktopNewRoute(names Names) error {
	stateDecls, inputFields, formFields := g.buildFormHelpers(names)

	createFunc := "Create" + names.Pascal

	content := "import { createFileRoute, useNavigate } from \"@tanstack/react-router\";\n"
	content += "import { useState } from \"react\";\n"
	content += "import { toast } from \"sonner\";\n"
	content += "// @ts-ignore\n"
	content += "import { " + createFunc + " } from \"../../../wailsjs/go/main/App\";\n"
	content += "\n"
	content += "export const Route = createFileRoute(\"/_layout/" + names.Plural + "/new\")({\n"
	content += "  component: " + names.Pascal + "NewPage,\n"
	content += "});\n"
	content += "\n"
	content += "function " + names.Pascal + "NewPage() {\n"
	content += "  const navigate = useNavigate();\n"
	content += "\n"
	content += stateDecls
	content += "  const [loading, setLoading] = useState(false);\n"
	content += "\n"
	content += "  const handleSubmit = async (e: React.FormEvent) => {\n"
	content += "    e.preventDefault();\n"
	content += "    setLoading(true);\n"
	content += "    try {\n"
	content += "      await " + createFunc + "({\n"
	content += inputFields
	content += "      });\n"
	content += "      toast.success(\"" + names.Pascal + " created\");\n"
	content += "      navigate({ to: \"/" + names.Plural + "\" });\n"
	content += "    } catch (err: any) {\n"
	content += "      toast.error(err?.message || \"Failed to create " + names.Lower + "\");\n"
	content += "    } finally {\n"
	content += "      setLoading(false);\n"
	content += "    }\n"
	content += "  };\n"
	content += "\n"
	content += "  return (\n"
	content += "    <div className=\"max-w-2xl\">\n"
	content += "      <h1 className=\"text-2xl font-bold text-foreground mb-1\">New " + names.Pascal + "</h1>\n"
	content += "      <p className=\"text-sm text-text-secondary mb-6\">Create a new " + names.Lower + "</p>\n"
	content += "\n"
	content += "      <form onSubmit={handleSubmit} className=\"space-y-5\">\n"
	content += "        <div className=\"bg-bg-elevated border border-border rounded-xl p-6 space-y-4\">\n"
	content += formFields
	content += "        </div>\n"
	content += "\n"
	content += "        <div className=\"flex items-center gap-3\">\n"
	content += "          <button\n"
	content += "            type=\"submit\"\n"
	content += "            disabled={loading}\n"
	content += "            className=\"px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50\"\n"
	content += "          >\n"
	content += "            {loading ? \"Creating...\" : \"Create " + names.Pascal + "\"}\n"
	content += "          </button>\n"
	content += "          <button\n"
	content += "            type=\"button\"\n"
	content += "            onClick={() => navigate({ to: \"/" + names.Plural + "\" })}\n"
	content += "            className=\"px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors\"\n"
	content += "          >\n"
	content += "            Cancel\n"
	content += "          </button>\n"
	content += "        </div>\n"
	content += "      </form>\n"
	content += "    </div>\n"
	content += "  );\n"
	content += "}\n"

	path := filepath.Join(g.Root, "frontend", "src", "routes", "_layout", names.Plural+".new.tsx")
	return writeFileWithDirs(path, content)
}

// writeDesktopEditRoute generates frontend/src/routes/_layout/<plural>.$id.edit.tsx.
func (g *DesktopGenerator) writeDesktopEditRoute(names Names) error {
	stateDecls, inputFields, formFields := g.buildFormHelpers(names)
	fetchAssignments := g.buildFetchAssignments()

	getFunc := "Get" + names.Pascal
	updateFunc := "Update" + names.Pascal

	content := "import { createFileRoute, useNavigate } from \"@tanstack/react-router\";\n"
	content += "import { useState, useEffect } from \"react\";\n"
	content += "import { toast } from \"sonner\";\n"
	content += "// @ts-ignore\n"
	content += "import { " + getFunc + ", " + updateFunc + " } from \"../../../wailsjs/go/main/App\";\n"
	content += "\n"
	content += "export const Route = createFileRoute(\"/_layout/" + names.Plural + "/$id/edit\")({\n"
	content += "  component: " + names.Pascal + "EditPage,\n"
	content += "});\n"
	content += "\n"
	content += "function " + names.Pascal + "EditPage() {\n"
	content += "  const { id } = Route.useParams();\n"
	content += "  const navigate = useNavigate();\n"
	content += "\n"
	content += stateDecls
	content += "  const [loading, setLoading] = useState(false);\n"
	content += "  const [fetching, setFetching] = useState(true);\n"
	content += "\n"
	content += "  useEffect(() => {\n"
	content += "    " + getFunc + "(Number(id))\n"
	content += "      .then((data: any) => {\n"
	content += fetchAssignments
	content += "      })\n"
	content += "      .catch((err: any) => {\n"
	content += "        toast.error(err?.message || \"Failed to load " + names.Lower + "\");\n"
	content += "        navigate({ to: \"/" + names.Plural + "\" });\n"
	content += "      })\n"
	content += "      .finally(() => setFetching(false));\n"
	content += "  }, [id, navigate]);\n"
	content += "\n"
	content += "  const handleSubmit = async (e: React.FormEvent) => {\n"
	content += "    e.preventDefault();\n"
	content += "    setLoading(true);\n"
	content += "    try {\n"
	content += "      await " + updateFunc + "(Number(id), {\n"
	content += inputFields
	content += "      });\n"
	content += "      toast.success(\"" + names.Pascal + " updated\");\n"
	content += "      navigate({ to: \"/" + names.Plural + "\" });\n"
	content += "    } catch (err: any) {\n"
	content += "      toast.error(err?.message || \"Failed to update " + names.Lower + "\");\n"
	content += "    } finally {\n"
	content += "      setLoading(false);\n"
	content += "    }\n"
	content += "  };\n"
	content += "\n"
	content += "  if (fetching) {\n"
	content += "    return (\n"
	content += "      <div className=\"flex items-center justify-center py-20\">\n"
	content += "        <div className=\"w-6 h-6 border-2 border-accent border-t-transparent rounded-full animate-spin\" />\n"
	content += "      </div>\n"
	content += "    );\n"
	content += "  }\n"
	content += "\n"
	content += "  return (\n"
	content += "    <div className=\"max-w-2xl\">\n"
	content += "      <h1 className=\"text-2xl font-bold text-foreground mb-1\">Edit " + names.Pascal + "</h1>\n"
	content += "      <p className=\"text-sm text-text-secondary mb-6\">Update the " + names.Lower + " details below</p>\n"
	content += "\n"
	content += "      <form onSubmit={handleSubmit} className=\"space-y-5\">\n"
	content += "        <div className=\"bg-bg-elevated border border-border rounded-xl p-6 space-y-4\">\n"
	content += formFields
	content += "        </div>\n"
	content += "\n"
	content += "        <div className=\"flex items-center gap-3\">\n"
	content += "          <button\n"
	content += "            type=\"submit\"\n"
	content += "            disabled={loading}\n"
	content += "            className=\"px-6 py-2.5 text-sm font-medium text-white bg-accent hover:bg-accent/90 rounded-lg transition-colors disabled:opacity-50\"\n"
	content += "          >\n"
	content += "            {loading ? \"Saving...\" : \"Update " + names.Pascal + "\"}\n"
	content += "          </button>\n"
	content += "          <button\n"
	content += "            type=\"button\"\n"
	content += "            onClick={() => navigate({ to: \"/" + names.Plural + "\" })}\n"
	content += "            className=\"px-6 py-2.5 text-sm font-medium text-text-secondary bg-bg-elevated border border-border rounded-lg hover:bg-bg-hover transition-colors\"\n"
	content += "          >\n"
	content += "            Cancel\n"
	content += "          </button>\n"
	content += "        </div>\n"
	content += "      </form>\n"
	content += "    </div>\n"
	content += "  );\n"
	content += "}\n"

	path := filepath.Join(g.Root, "frontend", "src", "routes", "_layout", names.Plural+".$id.edit.tsx")
	return writeFileWithDirs(path, content)
}

// buildFormHelpers extracts shared state declarations, input fields, and form fields for new/edit routes.
func (g *DesktopGenerator) buildFormHelpers(names Names) (stateDecls, inputFields, formFields string) {
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		camel := toLowerFirst(toPascalCase(f.Name))
		setter := "set" + toPascalCase(f.Name)

		switch f.FormFieldType() {
		case "number":
			stateDecls += "  const [" + camel + ", " + setter + "] = useState(0);\n"
		case "toggle":
			stateDecls += "  const [" + camel + ", " + setter + "] = useState(false);\n"
		default:
			stateDecls += "  const [" + camel + ", " + setter + "] = useState(\"\");\n"
		}
	}

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		camel := toLowerFirst(toPascalCase(f.Name))
		jsonName := toSnakeCase(f.Name)

		switch f.FormFieldType() {
		case "number":
			inputFields += "        " + jsonName + ": Number(" + camel + "),\n"
		default:
			inputFields += "        " + jsonName + ": " + camel + ",\n"
		}
	}

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		camel := toLowerFirst(toPascalCase(f.Name))
		setter := "set" + toPascalCase(f.Name)
		labelParts := splitPascal(toPascalCase(f.Name))
		label := strings.Join(labelParts, " ")
		formFieldType := f.FormFieldType()

		switch formFieldType {
		case "text":
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <input\n"
			formFields += "              type=\"text\"\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(e.target.value)}\n"
			formFields += "              placeholder=\"Enter " + toLower(label) + "\"\n"
			if f.Required {
				formFields += "              required\n"
			}
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"

		case "textarea", "richtext":
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <textarea\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(e.target.value)}\n"
			formFields += "              placeholder=\"Enter " + toLower(label) + "\"\n"
			formFields += "              rows={6}\n"
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors resize-y\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"

		case "number":
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <input\n"
			formFields += "              type=\"number\"\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(Number(e.target.value))}\n"
			formFields += "              placeholder=\"0\"\n"
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"

		case "toggle":
			formFields += "          <div className=\"flex items-center gap-3\">\n"
			formFields += "            <button\n"
			formFields += "              type=\"button\"\n"
			formFields += "              onClick={() => " + setter + "(!" + camel + ")}\n"
			formFields += "              className={\"relative inline-flex h-6 w-11 items-center rounded-full transition-colors \" + (" + camel + " ? \"bg-accent\" : \"bg-border\")}\n"
			formFields += "            >\n"
			formFields += "              <span className={\"inline-block h-4 w-4 rounded-full bg-white transition-transform \" + (" + camel + " ? \"translate-x-6\" : \"translate-x-1\")} />\n"
			formFields += "            </button>\n"
			formFields += "            <label className=\"text-sm text-foreground\">" + label + "</label>\n"
			formFields += "          </div>\n\n"

		case "date":
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <input\n"
			formFields += "              type=\"date\"\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(e.target.value)}\n"
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"

		case "datetime":
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <input\n"
			formFields += "              type=\"datetime-local\"\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(e.target.value)}\n"
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"

		default:
			formFields += "          <div>\n"
			formFields += "            <label className=\"block text-sm font-medium text-foreground mb-1.5\">" + label + "</label>\n"
			formFields += "            <input\n"
			formFields += "              type=\"text\"\n"
			formFields += "              value={" + camel + "}\n"
			formFields += "              onChange={(e) => " + setter + "(e.target.value)}\n"
			formFields += "              placeholder=\"Enter " + toLower(label) + "\"\n"
			formFields += "              className=\"w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors\"\n"
			formFields += "            />\n"
			formFields += "          </div>\n\n"
		}
	}

	return stateDecls, inputFields, formFields
}

// buildFetchAssignments generates the setter calls for populating form state from fetched data.
func (g *DesktopGenerator) buildFetchAssignments() string {
	fetchAssignments := ""
	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}
		setter := "set" + toPascalCase(f.Name)
		jsonName := toSnakeCase(f.Name)

		switch f.FormFieldType() {
		case "number":
			fetchAssignments += "        " + setter + "(data." + jsonName + " || 0);\n"
		case "toggle":
			fetchAssignments += "        " + setter + "(data." + jsonName + " || false);\n"
		case "date", "datetime":
			fetchAssignments += "        " + setter + "(data." + jsonName + " ? data." + jsonName + ".slice(0, 10) : \"\");\n"
		default:
			fetchAssignments += "        " + setter + "(data." + jsonName + " || \"\");\n"
		}
	}
	return fetchAssignments
}
