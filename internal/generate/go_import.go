package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeGoImportHandler generates internal/handlers/<name>_import.go — a bulk CSV
// import endpoint plus a template endpoint. It lives in its own file so its
// extra imports don't touch the main handler. Generated for every architecture;
// POST /<plural>/import and GET /<plural>/import/template are wired by injectAll.
//
// The CSV header row uses json field names. Recognised columns map to typed
// fields; unknown columns and the auto-managed id/slug are ignored. Everything
// is optional except fields the model marks required:
//   - file / files columns are skipped (upload images in the app instead).
//   - a belongs_to is given by NAME (column "category", not "category_id"):
//     the related record is looked up by name and created if it doesn't exist.
//   - rows that violate a unique constraint are skipped (ON CONFLICT DO NOTHING),
//     so re-importing the same file is safe. Other failures are reported per-row.
func (g *Generator) writeGoImportHandler(names Names) error {
	var assign strings.Builder
	var headers []string
	needStrconv := false

	for _, f := range g.Definition.Fields {
		t := FieldType(f.Type)
		if t == FieldSlug || t == FieldFile || t == FieldFiles ||
			t == FieldManyToMany || t == FieldStringArray ||
			t == FieldDatetime || t == FieldDate {
			continue
		}

		if f.IsBelongsTo() {
			relModel := f.RelatedModelName()
			base := strings.TrimSuffix(toSnakeCase(f.Name), "_id")
			fkGo := toPascalCase(base) + "ID"
			headers = append(headers, base)
			// Resolve the relationship by name: look it up, create it if missing.
			// Optional — an empty cell leaves the foreign key unset.
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
				"\t\t\tvar rel models.%s\n"+
				"\t\t\tif err := h.DB.Where(\"name = ?\", v).First(&rel).Error; err != nil {\n"+
				"\t\t\t\trel = models.%s{Name: v}\n"+
				"\t\t\t\th.DB.Create(&rel)\n"+
				"\t\t\t}\n"+
				"\t\t\titem.%s = rel.ID\n"+
				"\t\t}\n", base, relModel, relModel, fkGo))
			continue
		}

		goName := toPascalCase(f.Name)
		jsonName := toSnakeCase(f.Name)
		headers = append(headers, jsonName)
		switch t {
		case FieldInt:
			needStrconv = true
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.Atoi(v)\n\t\t\titem.%s = n\n\t\t}\n", jsonName, goName))
		case FieldUint:
			needStrconv = true
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.Atoi(v)\n\t\t\titem.%s = uint(n)\n\t\t}\n", jsonName, goName))
		case FieldFloat:
			needStrconv = true
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.ParseFloat(v, 64)\n\t\t\titem.%s = n\n\t\t}\n", jsonName, goName))
		case FieldBool:
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v == \"true\" || v == \"1\" || v == \"yes\"\n\t\t}\n", jsonName, goName))
		default: // string, text, richtext
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v\n\t\t}\n", jsonName, goName))
		}
	}

	strconvImport := ""
	if needStrconv {
		strconvImport = "\n\t\"strconv\""
	}
	templateHeaders := strings.Join(headers, ",")

	content := fmt.Sprintf(`package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"%s
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"%s/internal/models"
)

// Import bulk-creates %s from an uploaded CSV file (multipart field "file").
func (h *%sHandler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE", "message": "No CSV file provided"},
		})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_CSV", "message": "Could not read the CSV header row"},
		})
		return
	}

	idx := map[string]int{}
	for i, name := range headers {
		idx[strings.TrimSpace(strings.ToLower(name))] = i
	}
	get := func(rec []string, key string) (string, bool) {
		if i, ok := idx[key]; ok && i < len(rec) {
			return strings.TrimSpace(rec[i]), true
		}
		return "", false
	}

	created, skipped, failed := 0, 0, 0
	rowErrors := []gin.H{}
	rowNum := 1 // header is row 1
	for {
		rowNum++
		rec, rerr := reader.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			failed++
			if len(rowErrors) < 50 {
				rowErrors = append(rowErrors, gin.H{"row": rowNum, "message": "malformed row"})
			}
			continue
		}

		item := models.%s{}
%s		// Skip rows that collide with a unique constraint instead of erroring.
		res := h.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&item)
		if res.Error != nil {
			failed++
			if len(rowErrors) < 50 {
				rowErrors = append(rowErrors, gin.H{"row": rowNum, "message": res.Error.Error()})
			}
			continue
		}
		if res.RowsAffected == 0 {
			skipped++
			continue
		}
		created++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"created": created,
			"skipped": skipped,
			"failed":  failed,
			"errors":  rowErrors,
		},
		"message": fmt.Sprintf("Imported %%d, skipped %%d, failed %%d", created, skipped, failed),
	})
}

// Template returns a ready-to-fill CSV template (header row) for importing %s.
// belongs_to columns use the related record's NAME (e.g. "category").
func (h *%sHandler) Template(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `+"`"+`attachment; filename="%s-template.csv"`+"`"+`)
	c.String(http.StatusOK, "%s\n")
}
`, strconvImport, g.Module, names.Plural, names.Pascal, names.Pascal, assign.String(),
		names.Plural, names.Pascal, names.PluralKebab, templateHeaders)

	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_import.go")
	return writeFileWithDirs(path, content)
}
