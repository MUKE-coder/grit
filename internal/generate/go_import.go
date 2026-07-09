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
// The import runs in the BACKGROUND: the handler reads the upload, creates an
// ImportJob row, then processes the CSV in a goroutine and returns 202 with the
// job id immediately. Clients poll GET /imports/:id (shared handler) for a live
// progress bar and the final counts — so a large file never blocks the request
// and the app can leave the screen while it runs.
//
// The CSV header row uses json field names. Recognised columns map to typed
// fields; unknown columns and the auto-managed id/slug are ignored. Everything
// is optional except fields the model marks required:
//   - file / files columns are skipped (upload images in the app instead).
//   - a belongs_to is given by NAME (column "category", not "category_id"):
//     the related record is looked up by name and created if it doesn't exist.
//   - rows that violate a unique constraint are skipped (ON CONFLICT DO NOTHING),
//     so re-importing the same file is safe. Other failures are reported per-row.
// belongsToLookup describes how a CSV column resolves a belongs_to relation.
// When ByName is true the related record is matched (and created if missing)
// on the NaturalKeyJSON/NaturalKeyGo string column; otherwise the CSV cell is
// treated as the related record's ID (no phantom-create). ByName is only
// chosen when the related model actually HAS a usable string column — this is
// what stops `belongs_to:User` (no Name field) from emitting an uncompilable
// `models.User{Name: v}`.
type belongsToLookup struct {
	ByName        bool
	NaturalKeyJSON string
	NaturalKeyGo   string
}

// resolveBelongsToLookup inspects the related model's generated Go file to pick
// the column a CSV import should resolve the relation by. It prefers a
// human-friendly natural key (name/title/slug/label/email/username), then any
// other string field. If the model file can't be read or has no string field
// (e.g. a pure join/lookup model), it falls back to ID-based resolution so the
// generated handler always compiles.
func (g *Generator) resolveBelongsToLookup(relModel string) belongsToLookup {
	path := filepath.Join(g.APIRoot(), "internal", "models", toSnakeCase(relModel)+".go")
	structs, err := parseGoStructs(path)
	if err != nil {
		return belongsToLookup{ByName: false}
	}

	var fields []GoField
	for _, s := range structs {
		if s.Name == relModel {
			fields = s.Fields
			break
		}
	}
	if len(fields) == 0 {
		return belongsToLookup{ByName: false}
	}

	isString := func(t string) bool { return t == "string" }
	preferred := []string{"name", "title", "slug", "label", "username", "email"}
	for _, want := range preferred {
		for _, f := range fields {
			if f.JSONName == want && isString(f.GoType) {
				return belongsToLookup{ByName: true, NaturalKeyJSON: want, NaturalKeyGo: f.Name}
			}
		}
	}
	// Any other string field (skip the id/system columns).
	for _, f := range fields {
		if isString(f.GoType) && f.JSONName != "id" && f.Name != "ID" {
			return belongsToLookup{ByName: true, NaturalKeyJSON: f.JSONName, NaturalKeyGo: f.Name}
		}
	}
	return belongsToLookup{ByName: false}
}

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
			lookup := g.resolveBelongsToLookup(relModel)

			if lookup.ByName {
				// The related model has a usable string column — resolve by
				// that natural key and create the record if it's missing.
				// Column is the relation name (e.g. "category").
				headers = append(headers, base)
				assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
					"\t\t\tvar rel models.%s\n"+
					"\t\t\tif err := h.DB.Where(%q, v).First(&rel).Error; err != nil {\n"+
					"\t\t\t\trel = models.%s{%s: v}\n"+
					"\t\t\t\th.DB.Create(&rel)\n"+
					"\t\t\t}\n"+
					"\t\t\titem.%s = rel.ID\n"+
					"\t\t}\n", base, relModel, lookup.NaturalKeyJSON+" = ?", relModel, lookup.NaturalKeyGo, fkGo))
			} else {
				// No natural-key string column (e.g. belongs_to:User) —
				// resolve by the related record's ID. Column is "<base>_id".
				// The related record is NOT auto-created; an unknown or empty
				// id simply leaves the foreign key unset.
				idCol := base + "_id"
				headers = append(headers, idCol)
				assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
					"\t\t\tvar rel models.%s\n"+
					"\t\t\tif err := h.DB.Where(\"id = ?\", v).First(&rel).Error; err == nil {\n"+
					"\t\t\t\titem.%s = rel.ID\n"+
					"\t\t\t}\n"+
					"\t\t}\n", idCol, relModel, fkGo))
			}
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

	rep := strings.NewReplacer(
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{Plural}}", names.Plural,
		"{{PluralKebab}}", names.PluralKebab,
		"{{STRCONV}}", strconvImport,
		"{{ASSIGN}}", assign.String(),
		"{{HEADERS}}", templateHeaders,
	)

	content := rep.Replace(`package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"{{STRCONV}}
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"{{MODULE}}/internal/models"
)

// Import kicks off a BACKGROUND CSV import of {{Plural}}. It streams the upload
// to a temp file (so a large file never sits in memory), creates an ImportJob,
// then processes rows in a goroutine and returns 202 immediately. Poll
// GET /imports/:id for progress and the result.
func (h *{{Pascal}}Handler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE", "message": "No CSV file provided"},
		})
		return
	}
	defer file.Close()

	// Stream the upload to a temp file — never ReadAll a large CSV into memory.
	tmp, err := os.CreateTemp("", "grit-import-*.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TEMP_ERROR", "message": "Could not buffer the upload"},
		})
		return
	}
	tmpPath := tmp.Name()
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_CSV", "message": "Could not read the upload"},
		})
		return
	}
	tmp.Close()

	// Count data rows up front (streaming) so the client's progress bar has a
	// denominator without holding the file in memory.
	total, err := countCSVRows{{Pascal}}(tmpPath)
	if err != nil || total < 0 {
		os.Remove(tmpPath)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_CSV", "message": "Could not read the CSV file"},
		})
		return
	}

	job := models.ImportJob{Resource: "{{Plural}}", Status: "processing", Total: total}
	if err := h.DB.Create(&job).Error; err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "JOB_ERROR", "message": "Could not start import"},
		})
		return
	}

	// Process in the background so a large file never blocks the request.
	go h.runImport{{Pascal}}(job.ID, tmpPath)

	c.JSON(http.StatusAccepted, gin.H{
		"data":    gin.H{"job_id": job.ID, "total": total},
		"message": "Import started",
	})
}

// countCSVRows{{Pascal}} counts data rows (excluding the header) without holding
// the file in memory.
func countCSVRows{{Pascal}}(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true
	n := 0
	for i := 0; ; i++ {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		if i == 0 {
			continue // header row
		}
		n++
	}
	return n, nil
}

// runImport{{Pascal}} streams the temp CSV, creating {{Plural}} in batches and
// updating the ImportJob as it goes. belongs_to columns are resolved by their
// natural key (or id); unique-conflict rows are skipped; per-row failures are
// recorded. The temp file is removed when done.
func (h *{{Pascal}}Handler) runImport{{Pascal}}(jobID, tmpPath string) {
	defer os.Remove(tmpPath)
	// This runs in a bare goroutine, so gin.Recovery() does NOT cover it — an
	// unrecovered panic here would crash the whole server. Recover, and mark
	// the job failed so the client's poll terminates instead of hanging.
	defer func() {
		if r := recover(); r != nil {
			h.DB.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
				"status":  "failed",
				"message": fmt.Sprintf("import crashed: %v", r),
			})
		}
	}()

	f, err := os.Open(tmpPath)
	if err != nil {
		h.DB.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
			"status": "failed", "message": "could not reopen upload",
		})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		h.DB.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
			"status": "failed", "message": "empty or invalid CSV",
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
	rowErrors := []map[string]interface{}{}

	// checkpoint writes current progress so the client's poll sees movement.
	checkpoint := func(status, message string) {
		errsJSON, _ := json.Marshal(rowErrors)
		h.DB.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
			"status":    status,
			"processed": created + skipped + failed,
			"created":   created,
			"skipped":   skipped,
			"failed":    failed,
			"errors":    string(errsJSON),
			"message":   message,
		})
	}

	const batchSize = 200
	type pendingRow struct {
		item   models.{{Pascal}}
		rowNum int
	}
	batch := make([]pendingRow, 0, batchSize)

	// flush inserts the accumulated batch. CreateInBatches (with OnConflict
	// DoNothing) amortises the per-row fsync that makes large SQLite imports
	// crawl. If the whole batch errors (a bad row can poison it), we fall back
	// to per-row inserts so created/skipped/failed stay accurate and only the
	// offending row is dropped.
	flush := func() {
		if len(batch) == 0 {
			return
		}
		items := make([]models.{{Pascal}}, len(batch))
		for i := range batch {
			items[i] = batch[i].item
		}
		res := h.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(items, len(items))
		if res.Error == nil {
			created += int(res.RowsAffected)
			skipped += len(items) - int(res.RowsAffected)
		} else {
			for i := range batch {
				one := batch[i].item
				r := h.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&one)
				switch {
				case r.Error != nil:
					failed++
					if len(rowErrors) < 50 {
						rowErrors = append(rowErrors, map[string]interface{}{"row": batch[i].rowNum, "message": r.Error.Error()})
					}
				case r.RowsAffected == 0:
					skipped++
				default:
					created++
				}
			}
		}
		batch = batch[:0]
	}

	rowNum := 1 // header was row 1
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			failed++
			if len(rowErrors) < 50 {
				rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": err.Error()})
			}
			continue
		}

		item := models.{{Pascal}}{}
{{ASSIGN}}		batch = append(batch, pendingRow{item: item, rowNum: rowNum})

		if len(batch) >= batchSize {
			flush()
			checkpoint("processing", "")
		}
	}
	flush()

	checkpoint("completed", fmt.Sprintf("Imported %d, skipped %d, failed %d", created, skipped, failed))
}

// Template returns a ready-to-fill CSV template (header row) for importing {{Plural}}.
// belongs_to columns use the related record's natural key (e.g. "category"), or
// its id column ("<relation>_id") when the related model has no natural key.
func (h *{{Pascal}}Handler) Template(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `+"`"+`attachment; filename="{{PluralKebab}}-template.csv"`+"`"+`)
	c.String(http.StatusOK, "{{HEADERS}}\n")
}
`)

	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_import.go")
	return writeFileWithDirs(path, content)
}
