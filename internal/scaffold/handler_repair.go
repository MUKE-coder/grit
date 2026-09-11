package scaffold

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairResourceHandlers patches the code of resources generated before
// v3.214.0, where the generated text is still recognisable.
//
// Upgrade does not regenerate resource code: it is the developer's, and the
// definition it came from is not kept anywhere. These defects are the
// exception, because one is a hole an ordinary account could walk through and
// the other broke a button in every project:
//
//   - A resource generated with --owned-by checked ownership on list, read,
//     update and delete, and nowhere else. Export returned every row, PDF
//     printed anyone's record, PATCH edited it, Bulk acted on it, the importer
//     filed rows under whichever user the CSV named, and a workflow transition
//     moved anyone's row.
//   - Export read each batch through tx.Scan on the fresh session
//     FindInBatches hands its callback, which has no query on it, so every CSV
//     export was an empty 200. The importer also created a user for every
//     email it did not recognise.
//
// Each edit anchors on text the generator wrote. Where that text has been
// changed the edit is skipped and named, so a handler somebody rewrote is never
// half-patched in silence. A file the manifest records as untouched has its
// hash refreshed; an edited one keeps its old hash, so the edit stays visible.
func repairResourceHandlers(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	handlers, err := goSources(filepath.Join(apiRoot, "internal", "handlers"))
	if err != nil {
		return err
	}
	// Handlers first: they say which resources are owned, and the importer and
	// workflow passes need to know.
	owned := map[string]string{}
	for _, path := range handlers {
		if strings.HasSuffix(path, "_import.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, fixed, warn, pascal, col := repairHandlerSource(src)
			if col != "" {
				owned[pascal] = col
			}
			return out, fixed, warn
		}); err != nil {
			return err
		}
	}
	for _, path := range handlers {
		if !strings.HasSuffix(path, "_import.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairImportSource(src, owned, opts.Module())
		}); err != nil {
			return err
		}
	}

	services, err := goSources(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	for _, path := range services {
		if !strings.HasSuffix(path, "_workflow.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairWorkflowSource(src, owned, opts.Module())
		}); err != nil {
			return err
		}
	}
	return nil
}

// goSources lists the non-test Go files in dir, or nothing when it is absent.
func goSources(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	var out []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		out = append(out, filepath.Join(dir, name))
	}
	return out, nil
}

// repairSourceFile applies fn to one file, keeping its line endings, and
// writes the result only when it is valid Go.
func repairSourceFile(root string, m *manifest.Manifest, path string, fn func(string) (string, []string, []string)) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	src := strings.ReplaceAll(string(raw), "\r\n", "\n")

	out, fixed, warnings := fn(src)
	shown := path
	if rel, err := filepath.Rel(root, path); err == nil {
		shown = filepath.ToSlash(rel)
	}
	for _, w := range warnings {
		fmt.Printf("  ⚠ %s: %s\n", shown, w)
	}
	if out == src {
		return nil
	}

	formatted, err := format.Source([]byte(out))
	if err != nil {
		fmt.Printf("  ⚠ %s: left alone, the repair did not produce valid Go (%v)\n", shown, err)
		return nil
	}
	out = string(formatted)

	pristine := false
	if key, inside := manifest.Rel(root, path); inside {
		pristine = m.StatusOf(root, key) == manifest.Unchanged
	}
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if pristine {
		manifest.Refresh(path)
	}
	fmt.Printf("  ✓ %s: %s\n", shown, strings.Join(fixed, "; "))
	return nil
}

var (
	repairExportFuncRe = regexp.MustCompile(`func \(h \*(\w+)Handler\) Export\(c \*gin\.Context\) \{`)
	repairScanBatchRe  = regexp.MustCompile(`query\.FindInBatches\(&\[\]models\.(\w+)\{\}, exportBatchSize, func\(tx \*gorm\.DB, batch int\) error \{\n\s*var rows \[\]models\.\w+\n\s*if err := tx\.Scan\(&rows\)\.Error; err != nil \{\n\s*return err\n\s*\}\n`)
	repairOrderRe      = regexp.MustCompile(`(query := h\.(?:scoped\(c\)|DB)\.Model\(&models\.\w+\{\}\)[^\n]*?)\.Order\("created_at desc"\)`)
	repairOwnerColRe   = regexp.MustCompile(`authz\.ScopeToOwner\(c, query, "(\w+)"\)`)
	repairLoadedRowRe  = regexp.MustCompile(`(?s)First\(&item, "id = \?", id\)\.Error; err != nil \{.*?\n\t\treturn\n\t\}\n`)
)

const repairOwnerGuard = `
	// --owned-by: 404 rather than 403 on somebody else's row, so a wrong guess
	// cannot be told from a right one.
	if !authz.OwnsOr404(c, &item) {
		return
	}
`

// repairHandlerSource fixes one resource handler. It reports the resource and,
// for an owned one, the owner column.
func repairHandlerSource(src string) (out string, fixed, warn []string, pascal, col string) {
	head := repairExportFuncRe.FindStringSubmatch(src)
	if head == nil {
		return src, nil, nil, "", ""
	}
	pascal = head[1]
	out = src

	if body, ok := handlerFunc(out, pascal, "Export"); ok && repairScanBatchRe.MatchString(body) {
		model := repairScanBatchRe.FindStringSubmatch(body)[1]
		const anchor = "\n\tif format == \"xlsx\" {"
		if strings.Contains(body, anchor) {
			nb := repairScanBatchRe.ReplaceAllString(body,
				"query.FindInBatches(&rows, exportBatchSize, func(tx *gorm.DB, batch int) error {\n")
			nb = repairOrderRe.ReplaceAllString(nb, "$1")
			nb = strings.Replace(nb, anchor, "\n\tvar rows []models."+model+anchor, 1)
			// Nothing matched: still a CSV with its header row, not an empty body.
			nb = strings.Replace(nb, "\t\t_ = err\n\t}\n}\n",
				"\t\t_ = err\n\t} else if !headerWritten {\n"+
					"\t\t// Nothing matched: still a CSV, with its header row.\n"+
					"\t\t_ = export.CSV(c.Writer, rows, opts)\n\t}\n}\n", 1)
			out = replaceHandlerFunc(out, pascal, "Export", nb)
			fixed = append(fixed, "CSV export returns its rows")
		} else {
			warn = append(warn, "Export reads its batches through tx.Scan, which returns an empty file, "+
				"and it is too changed to patch: read each batch from the slice passed to FindInBatches")
		}
	}

	owner := repairOwnerColRe.FindStringSubmatch(out)
	if owner == nil {
		return out, fixed, warn, pascal, ""
	}
	col = owner[1]
	var guarded, missed []string

	if body, ok := handlerFunc(out, pascal, "Export"); ok && !strings.Contains(body, "ScopeToOwner") {
		const anchor = "\n\topts := export.Options{"
		if strings.Contains(body, anchor) {
			scope := "\n\t// --owned-by: an export is the list without pages, scoped the same way.\n" +
				"\t// ADMIN is exempt.\n" +
				"\tquery = authz.ScopeToOwner(c, query, \"" + col + "\")\n"
			out = replaceHandlerFunc(out, pascal, "Export", strings.Replace(body, anchor, scope+anchor, 1))
			guarded = append(guarded, "Export")
		} else {
			missed = append(missed, "Export")
		}
	}
	for _, fn := range []string{"PDF", "Patch"} {
		body, ok := handlerFunc(out, pascal, fn)
		if !ok || strings.Contains(body, "OwnsOr404") {
			continue
		}
		loc := repairLoadedRowRe.FindStringIndex(body)
		if loc == nil {
			missed = append(missed, fn)
			continue
		}
		out = replaceHandlerFunc(out, pascal, fn, body[:loc[1]]+repairOwnerGuard+body[loc[1]:])
		guarded = append(guarded, fn)
	}
	if body, ok := handlerFunc(out, pascal, "Bulk"); ok && !strings.Contains(body, "ScopeToOwner") {
		const anchor = "\tif err := scope.Find(&items).Error"
		if strings.Contains(body, anchor) {
			scope := "\t// --owned-by: only the caller's own rows are acted on, and somebody\n" +
				"\t// else's id drops out as if it did not exist. ADMIN is exempt.\n" +
				"\tscope = authz.ScopeToOwner(c, scope, \"" + col + "\")\n"
			out = replaceHandlerFunc(out, pascal, "Bulk", strings.Replace(body, anchor, scope+anchor, 1))
			guarded = append(guarded, "Bulk")
		} else {
			missed = append(missed, "Bulk")
		}
	}

	if len(guarded) > 0 {
		fixed = append(fixed, "owner checks on "+strings.Join(guarded, ", "))
	}
	if len(missed) > 0 {
		warn = append(warn, "generated with --owned-by, but "+strings.Join(missed, ", ")+
			" could not be patched and still reach other users' rows: add authz.OwnsOr404 "+
			"or authz.ScopeToOwner there by hand")
	}
	return out, fixed, warn, pascal, col
}

// handlerFunc returns one method of a generated handler, from its signature to
// its closing brace.
func handlerFunc(src, pascal, name string) (string, bool) {
	start, end := handlerFuncBounds(src, pascal, name)
	if start < 0 {
		return "", false
	}
	return src[start:end], true
}

func replaceHandlerFunc(src, pascal, name, body string) string {
	start, end := handlerFuncBounds(src, pascal, name)
	if start < 0 {
		return src
	}
	return src[:start] + body + src[end:]
}

func handlerFuncBounds(src, pascal, name string) (int, int) {
	start := strings.Index(src, "func (h *"+pascal+"Handler) "+name+"(c *gin.Context) {")
	if start < 0 {
		return -1, -1
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return -1, -1
	}
	return start, start + end + len("\n}\n")
}

var (
	repairImportFuncRe = regexp.MustCompile(`func \(h \*(\w+)Handler\) Import\(c \*gin\.Context\) \{`)
	repairUserCreateRe = regexp.MustCompile(`(?m)^(\t+)if err := h\.DB\.Where\("(\w+) = \?", v\)\.First\(&rel\)\.Error; err != nil \{\n\t+rel = models\.User\{\w+: v\}\n\t+h\.DB\.Create\(&rel\)\n\t+\}\n`)
)

// repairImportSource stops an importer creating users, and makes an owned
// resource's imported rows belong to the importer.
func repairImportSource(src string, owned map[string]string, module string) (string, []string, []string) {
	head := repairImportFuncRe.FindStringSubmatch(src)
	if head == nil {
		return src, nil, nil
	}
	pascal := head[1]
	out := src
	var fixed, warn []string

	if repairUserCreateRe.MatchString(out) {
		out = repairUserCreateRe.ReplaceAllStringFunc(out, func(block string) string {
			sm := repairUserCreateRe.FindStringSubmatch(block)
			ind, key := sm[1], sm[2]
			return ind + "if err := h.DB.Where(\"" + key + " = ?\", v).First(&rel).Error; err != nil {\n" +
				ind + "\tfailed++\n" +
				ind + "\tif len(rowErrors) < 50 {\n" +
				ind + "\t\trowErrors = append(rowErrors, map[string]interface{}{\"row\": rowNum, \"message\": fmt.Sprintf(\"no user with " + key + " %q\", v)})\n" +
				ind + "\t}\n" +
				ind + "\tcontinue\n" +
				ind + "}\n"
		})
		fixed = append(fixed, "the importer no longer creates users")
	}

	col, isOwned := owned[pascal]
	if !isOwned || strings.Contains(out, "canAssignOwner") {
		return out, fixed, warn
	}
	call := "go h.runImport" + pascal + "(job.ID, tmpPath)"
	sig := "func (h *" + pascal + "Handler) runImport" + pascal + "(jobID, tmpPath string) {"
	build := "\t\titem := models." + pascal + "{}\n"
	if !strings.Contains(out, call) || !strings.Contains(out, sig) || !strings.Contains(out, build) {
		warn = append(warn, "generated with --owned-by, but the importer is too changed to patch: "+
			"it still takes the owner from the CSV")
		return out, fixed, warn
	}
	field := ownerGoField(col)
	next := strings.Replace(out, call, "go h.runImport"+pascal+"(job.ID, tmpPath, authz.CurrentUserID(c), authz.IsAdmin(c))", 1)
	next = strings.Replace(next, sig, "func (h *"+pascal+"Handler) runImport"+pascal+"(jobID, tmpPath string, ownerID string, canAssignOwner bool) {", 1)
	next = strings.Replace(next, build, build+
		"\t\t// --owned-by: a row belongs to whoever imports it. An ADMIN may name\n"+
		"\t\t// another owner in the CSV; for anyone else that column is ignored.\n"+
		"\t\titem."+field+" = ownerID\n", 1)
	// The CSV's own owner column now only speaks for an ADMIN.
	if k := strings.Index(next, "item."+field+" = rel.ID"); k >= 0 {
		if j := strings.LastIndex(next[:k], "if v, ok := get(rec, "); j >= 0 {
			const cond = `ok && v != "" {`
			if c := strings.Index(next[j:k], cond); c >= 0 {
				p := j + c
				next = next[:p] + `ok && v != "" && canAssignOwner {` + next[p+len(cond):]
			}
		}
	}
	next, ok := addImportGroup(next, module+"/internal/authz")
	if !ok {
		warn = append(warn, "generated with --owned-by, but the authz import could not be added: "+
			"the importer still takes the owner from the CSV")
		return out, fixed, warn
	}
	return next, append(fixed, "imported rows belong to the importer"), warn
}

// ownerGoField turns an owner column into its Go field: user_id is UserID.
func ownerGoField(col string) string {
	var b strings.Builder
	for _, part := range strings.Split(strings.TrimSuffix(col, "_id"), "_") {
		if part != "" {
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String() + "ID"
}

var repairTransitionRe = regexp.MustCompile(`func Transition(\w+)\(db \*gorm\.DB, c \*gin\.Context, id, action string`)

const repairWorkflowOwnerCheck = `
	// --owned-by: only the owner may move a row, ADMIN excepted, and somebody
	// else's row is not found rather than forbidden. c is nil when a job or a
	// command makes the move, and those are trusted.
	if c != nil && !authz.IsAdmin(c) && item.GetOwnerID() != authz.CurrentUserID(c) {
		return nil, gorm.ErrRecordNotFound
	}
`

// repairWorkflowSource makes an owned resource's transitions check the owner.
func repairWorkflowSource(src string, owned map[string]string, module string) (string, []string, []string) {
	head := repairTransitionRe.FindStringSubmatch(src)
	if head == nil {
		return src, nil, nil
	}
	if _, isOwned := owned[head[1]]; !isOwned || strings.Contains(src, "GetOwnerID()") {
		return src, nil, nil
	}
	const anchor = "\tif err := db.First(&item, \"id = ?\", id).Error; err != nil {\n\t\treturn nil, err\n\t}\n"
	if !strings.Contains(src, anchor) {
		return src, nil, []string{"generated with --owned-by, but Transition is too changed to patch: " +
			"any account can still move another user's row"}
	}
	out, ok := addImportGroup(strings.Replace(src, anchor, anchor+repairWorkflowOwnerCheck, 1), module+"/internal/authz")
	if !ok {
		return src, nil, []string{"generated with --owned-by, but the authz import could not be added to Transition"}
	}
	return out, []string{"transitions check the owner"}, nil
}
