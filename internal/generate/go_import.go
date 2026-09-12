package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeGoImportHandler generates the bulk CSV import of a resource, in two files.
//
// internal/handlers/<name>_import.go is the HTTP half: it takes the upload,
// counts its rows, starts the job and returns 202, and it serves the template.
// internal/services/<name>_import.go is the rest: it reads the rows, resolves
// their relations and writes them in batches, so the import runs no query from
// the handler. Their own files, so their extra imports don't touch the main
// handler and service. POST /<plural>/import and GET /<plural>/import/template
// are wired by injectAll.
//
// The import runs in the BACKGROUND: the handler streams the upload to a temp
// file, starts an ImportJob, hands the file to the service in a goroutine and
// returns 202 with the job id at once. Clients poll GET /imports/:id (shared
// handler) for a live progress bar and the final counts, so a large file never
// blocks the request and the app can leave the screen while it runs.
//
// The CSV header row uses json field names. Recognised columns map to typed
// fields; unknown columns and the auto-managed id/slug are ignored. Everything
// is optional except fields the model marks required:
//   - file / files columns are skipped (upload images in the app instead).
//   - a belongs_to is given by NAME (column "category", not "category_id"):
//     the related record is looked up by name and created if it doesn't exist.
//     A user is the exception: it must already exist, or the row fails.
//   - an --owned-by resource's rows belong to whoever imports them. Only an
//     ADMIN may name a different owner in the CSV.
//   - rows that violate a unique constraint are skipped (ON CONFLICT DO NOTHING),
//     so re-importing the same file is safe. Other failures are reported per-row.
//
// belongsToLookup describes how a CSV column resolves a belongs_to relation.
// When ByName is true the related record is matched (and created if missing)
// on the NaturalKeyJSON/NaturalKeyGo string column; otherwise the CSV cell is
// treated as the related record's ID (no phantom-create). ByName is only
// chosen when the related model actually HAS a usable string column: this is
// what stops `belongs_to:User` (no Name field) from emitting an uncompilable
// `models.User{Name: v}`.
type belongsToLookup struct {
	ByName         bool
	NaturalKeyJSON string
	NaturalKeyGo   string
}

// resolveBelongsToLookup inspects the related model's generated Go file to pick
// the column a CSV import should resolve the relation by. It prefers a
// human-friendly natural key (name/title/slug/label/email/username), then any
// other string field. If the model file can't be read or has no string field
// (e.g. a pure join/lookup model), it falls back to ID-based resolution so the
// generated code always compiles.
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
	needMoney := false
	needDatatypes := false

	needCrypto := false
	owned := false
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

			// A self-reference is a nullable column, so the importer assigns a
			// pointer to it. Every other foreign key is a plain string.
			assignExpr := "rel.ID"
			if relModel == toPascalCase(g.Definition.Name) {
				assignExpr = "&rel.ID"
			}

			if relModel == "User" {
				col, where := base+"_id", "id"
				if lookup.ByName {
					col, where = base, lookup.NaturalKeyJSON
				}
				headers = append(headers, col)
				cond := ""
				if owner := g.Definition.OwnerField(); owner != nil && owner.Name == f.Name {
					owned = true
					cond = " && canAssignOwner"
					assign.WriteString(fmt.Sprintf(
						"\t\t// --owned-by %s: a row belongs to whoever imports it. An ADMIN may\n"+
							"\t\t// name another owner in the %s column; for anyone else it is\n"+
							"\t\t// ignored, or a CSV could file records under somebody else's name.\n"+
							"\t\titem.%s = ownerID\n", f.Name, col, fkGo))
				}
				assign.WriteString(importUserLookup(col, where, fkGo, cond))
				continue
			}

			if lookup.ByName {
				// The related model has a usable string column: resolve by
				// that natural key and create the record if it's missing.
				// Column is the relation name (e.g. "category").
				headers = append(headers, base)
				assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
					"\t\t\tvar rel models.%s\n"+
					"\t\t\tif err := db.Where(%q, v).First(&rel).Error; err != nil {\n"+
					"\t\t\t\trel = models.%s{%s: v}\n"+
					"\t\t\t\tdb.Create(&rel)\n"+
					"\t\t\t}\n"+
					"\t\t\titem.%s = "+assignExpr+"\n"+
					"\t\t}\n", base, relModel, lookup.NaturalKeyJSON+" = ?", relModel, lookup.NaturalKeyGo, fkGo))
			} else {
				// No natural-key string column (e.g. belongs_to:User):
				// resolve by the related record's ID. Column is "<base>_id".
				// The related record is NOT auto-created; an unknown or empty
				// id simply leaves the foreign key unset.
				idCol := base + "_id"
				headers = append(headers, idCol)
				assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
					"\t\t\tvar rel models.%s\n"+
					"\t\t\tif err := db.Where(\"id = ?\", v).First(&rel).Error; err == nil {\n"+
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
		case FieldMoney:
			// A CSV cell carries a human figure like 19.99, so it is read as
			// major units and converted once. The currency comes from a
			// sibling <field>_currency column when the file has one; guessing
			// it from the amount is not possible and defaulting silently is
			// better than importing an amount with no currency at all.
			needStrconv = true
			needMoney = true
			assign.WriteString(fmt.Sprintf(
				"\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n"+
					"\t\t\tcur := \"USD\"\n"+
					"\t\t\tif c, ok := get(rec, %q); ok && c != \"\" {\n"+
					"\t\t\t\tcur = c\n"+
					"\t\t\t}\n"+
					"\t\t\tmajor, _ := strconv.ParseFloat(v, 64)\n"+
					"\t\t\titem.%s = money.FromMajor(major, cur)\n"+
					"\t\t}\n",
				jsonName, jsonName+"_currency", goName))
		case FieldBool, FieldToggle:
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v == \"true\" || v == \"1\" || v == \"yes\"\n\t\t}\n", jsonName, goName))
		case FieldCheck:
			// Multi-value cell: pipe-separated values → JSON string slice.
			needDatatypes = true
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\" {\n\t\t\titem.%s = datatypes.JSONSlice[string](strings.Split(v, \"|\"))\n\t\t}\n", jsonName, goName))
		default: // string, text, richtext, select
			if f.Encrypted {
				needCrypto = true
				assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = crypto.EncryptedString(v)\n\t\t}\n", jsonName, goName))
				continue
			}
			assign.WriteString(fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v\n\t\t}\n", jsonName, goName))
		}
	}

	cryptoImport := ""
	if needCrypto {
		cryptoImport = "\n\t\"" + g.Module + "/internal/crypto\""
	}
	moneyImport := ""
	if needMoney {
		moneyImport = "\n\t\"" + g.Module + "/internal/money\""
	}
	strconvImport := ""
	if needStrconv {
		strconvImport = "\n\t\"strconv\""
	}
	datatypesImport := ""
	if needDatatypes {
		datatypesImport = "\n\t\"gorm.io/datatypes\""
	}
	templateHeaders := strings.Join(headers, ",")

	ownerSetup, authzImport := "", ""
	if owned {
		// From the caller on the context rather than from parameters: the
		// handler puts it there, and a job importing for somebody says who.
		ownerSetup = "\t// --owned-by: rows belong to whoever imports them. Only an ADMIN may\n" +
			"\t// name another owner in the CSV.\n" +
			"\tactor, _ := authz.ActorFrom(ctx)\n" +
			"\townerID, canAssignOwner := actor.UserID, actor.Admin\n"
		authzImport = "\"" + g.Module + "/internal/authz\"\n\t"
	}

	rep := strings.NewReplacer(
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{Plural}}", names.Plural,
		"{{PluralKebab}}", names.PluralKebab,
		"{{STRCONV}}", strconvImport,
		"{{DATATYPES}}", datatypesImport+moneyImport+cryptoImport,
		"{{ASSIGN}}", assign.String(),
		"{{HEADERS}}", templateHeaders,
		"{{OWNER_SETUP}}", ownerSetup,
		"{{AUTHZ_IMPORT}}", authzImport,
	)

	handlerPath := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_import.go")
	if err := writeFileWithDirs(handlerPath, rep.Replace(importHandlerTemplate)); err != nil {
		return err
	}
	servicePath := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+"_import.go")
	return writeFileWithDirs(servicePath, rep.Replace(importServiceTemplate))
}

// importHandlerTemplate is the HTTP half of the import: the upload, the job,
// the 202, and the template download. It runs no query.
const importHandlerTemplate = `package handlers

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/respond"
)

// Import kicks off a BACKGROUND CSV import of {{Plural}}. It streams the upload
// to a temp file (so a large file never sits in memory), starts an ImportJob,
// then hands the file to the service in a goroutine and returns 202
// immediately. Poll GET /imports/:id for progress and the result.
func (h *{{Pascal}}Handler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		respond.Fail(c, respond.CodeInvalidFile, "No CSV file provided")
		return
	}
	defer file.Close()

	// Stream the upload to a temp file: never ReadAll a large CSV into memory.
	tmp, err := os.CreateTemp("", "grit-import-*.csv")
	if err != nil {
		respond.Fail(c, respond.CodeTempError, "Could not buffer the upload")
		return
	}
	tmpPath := tmp.Name()
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeInvalidCSV, "Could not read the upload")
		return
	}
	tmp.Close()

	// Count data rows up front (streaming) so the client's progress bar has a
	// denominator without holding the file in memory.
	total, err := countCSVRows{{Pascal}}(tmpPath)
	if err != nil || total < 0 {
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeInvalidCSV, "Could not read the CSV file")
		return
	}

	job, err := h.service().StartImport(h.ctx(c), total)
	if err != nil {
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeJobError, "Could not start import")
		return
	}

	// In the background, so a large file never blocks the request. The context
	// keeps the caller and the organization and drops the cancellation: the
	// request is over as soon as this returns, and the import has only begun.
	go h.service().ImportCSV(context.WithoutCancel(h.ctx(c)), job.ID, tmpPath)

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

// Template returns a ready-to-fill CSV template (header row) for importing {{Plural}}.
// belongs_to columns use the related record's natural key (e.g. "category"), or
// its id column ("<relation>_id") when the related model has no natural key.
func (h *{{Pascal}}Handler) Template(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{PluralKebab}}-template.csv"` + "`" + `)
	c.String(http.StatusOK, "{{HEADERS}}\n")
}
`

// importServiceTemplate is the rest of the import: the job, and the rows.
const importServiceTemplate = `package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"{{STRCONV}}
	"strings"

	"gorm.io/gorm/clause"{{DATATYPES}}

	{{AUTHZ_IMPORT}}"{{MODULE}}/internal/models"
)

// StartImport records a CSV import of {{Plural}} about to run: the job a client
// polls for progress.
func (s *{{Pascal}}Service) StartImport(ctx context.Context, total int) (*models.ImportJob, error) {
	job := models.ImportJob{Resource: "{{Plural}}", Status: "processing", Total: total}
	if err := s.db(ctx).Create(&job).Error; err != nil {
		return nil, fmt.Errorf("starting the {{Plural}} import: %w", err)
	}
	return &job, nil
}

// ImportCSV streams the CSV at path, creating {{Plural}} in batches and
// updating the ImportJob jobID as it goes. belongs_to columns are resolved by
// their natural key (or id); unique-conflict rows are skipped; per-row failures
// are recorded. The file is removed when done.
//
// It outlives the request that started it, so give it a context that is not
// cancelled with the response. context.WithoutCancel of the request's keeps the
// caller, whom an owned resource's rows belong to, and the organization the
// multitenant plugin stamps them with.
func (s *{{Pascal}}Service) ImportCSV(ctx context.Context, jobID, path string) {
	db := s.db(ctx)
{{OWNER_SETUP}}	defer os.Remove(path)
	// This runs in a bare goroutine, so gin.Recovery() does NOT cover it: an
	// unrecovered panic here would crash the whole server. Recover, and mark
	// the job failed so the client's poll terminates instead of hanging.
	defer func() {
		if r := recover(); r != nil {
			db.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
				"status":  "failed",
				"message": fmt.Sprintf("import crashed: %v", r),
			})
		}
	}()

	f, err := os.Open(path)
	if err != nil {
		db.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
			"status": "failed", "message": "could not reopen upload",
		})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		db.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
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
		db.Model(&models.ImportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
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
		res := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(items, len(items))
		if res.Error == nil {
			created += int(res.RowsAffected)
			skipped += len(items) - int(res.RowsAffected)
		} else {
			for i := range batch {
				one := batch[i].item
				r := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&one)
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
`

// importUserLookup resolves a CSV cell to an existing user and fails the row
// when there is none. It never creates one: a users row is an account, and an
// importer that made one per unrecognised email let any signed-in caller mint
// accounts around registration, with no password and nothing verified.
func importUserLookup(col, where, fkGo, cond string) string {
	return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\"%s {\n"+
		"\t\t\tvar rel models.User\n"+
		"\t\t\tif err := db.Where(%q, v).First(&rel).Error; err != nil {\n"+
		"\t\t\t\tfailed++\n"+
		"\t\t\t\tif len(rowErrors) < 50 {\n"+
		"\t\t\t\t\trowErrors = append(rowErrors, map[string]interface{}{\"row\": rowNum, \"message\": fmt.Sprintf(\"no user with %s %%q\", v)})\n"+
		"\t\t\t\t}\n"+
		"\t\t\t\tcontinue\n"+
		"\t\t\t}\n"+
		"\t\t\titem.%s = rel.ID\n"+
		"\t\t}\n", col, cond, where+" = ?", where, fkGo)
}
