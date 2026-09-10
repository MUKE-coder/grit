package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// `grit generate field` used to stop at the model, the Zod schemas, the
// TypeScript type and the admin. Never the handler.
//
// Generated handlers do not bind into the model. They bind into
// Create<X>Request and Update<X>Request and copy fields across one by one, and
// PATCH writes only what an allow-list names. So after adding notes:text to an
// Account, the admin form showed the field and sent it, and:
//
//	POST  {"notes": "supplier balances"}  201, notes came back empty
//	PUT   {"notes": "put"}                200, ignored
//	PATCH {"notes": "patched"}            422, not a writable field
//
// with the column sitting empty in Postgres. Found on a ledger, adding a field
// the way the CLI's own help recommends. The same hole was in the advice to
// edit the model by hand and run grit sync, which never touches a handler
// either.
//
// So a field is added everywhere the generator would have put it: both request
// structs, the create literal, the update map, the PATCH and bulk allow-lists,
// the list's sortable and filterable sets, and the CSV importer.

// handlerFieldParts is what one scalar field contributes to a handler.
type handlerFieldParts struct {
	createField  string
	createAssign string
	patchKey     string
	updateField  string
	updateSet    string
}

// scalarHandlerParts mirrors the scalar branch of writeGoHandler, rule for
// rule. TestAddFieldMatchesTheGenerator fails when the two drift, which is the
// point: a partial copy of a mapping is a copy that falls behind.
func scalarHandlerParts(f Field) handlerFieldParts {
	goName := toPascalCase(f.Name)
	goType := f.GoType()
	jsonTag := toSnakeCase(f.Name)

	bindingTag := ""
	if f.Required {
		bindingTag = ` binding:"required"`
	}

	p := handlerFieldParts{
		createField:  fmt.Sprintf("\t%s %s `json:\"%s\"%s`", goName, goType, jsonTag, bindingTag),
		createAssign: fmt.Sprintf("\t\t%s: req.%s,", goName, goName),
		patchKey:     fmt.Sprintf("\"%s\": true,", jsonTag),
	}
	switch goType {
	case "bool":
		p.updateField = fmt.Sprintf("\t%s *%s `json:\"%s\"`", goName, goType, jsonTag)
		p.updateSet = fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}", goName, jsonTag, goName)
	case "string":
		p.updateField = fmt.Sprintf("\t%s %s `json:\"%s\"`", goName, goType, jsonTag)
		p.updateSet = fmt.Sprintf("\tif req.%s != \"\" {\n\t\tupdates[\"%s\"] = req.%s\n\t}", goName, jsonTag, goName)
	case "*time.Time":
		p.updateField = fmt.Sprintf("\t%s %s `json:\"%s\"`", goName, goType, jsonTag)
		p.updateSet = fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = req.%s\n\t}", goName, jsonTag, goName)
	default:
		p.updateField = fmt.Sprintf("\t%s *%s `json:\"%s\"`", goName, goType, jsonTag)
		p.updateSet = fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}", goName, jsonTag, goName)
	}
	return p
}

// importAssign mirrors writeGoImportHandler for the types grit g field
// supports. ok is false for the types the importer skips (date, datetime).
func importAssign(f Field) (code string, needStrconv, ok bool) {
	goName := toPascalCase(f.Name)
	jsonName := toSnakeCase(f.Name)
	switch FieldType(f.Type) {
	case FieldDate, FieldDatetime:
		return "", false, false
	case FieldInt:
		return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.Atoi(v)\n\t\t\titem.%s = n\n\t\t}", jsonName, goName), true, true
	case FieldUint:
		return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.Atoi(v)\n\t\t\titem.%s = uint(n)\n\t\t}", jsonName, goName), true, true
	case FieldFloat:
		return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\tn, _ := strconv.ParseFloat(v, 64)\n\t\t\titem.%s = n\n\t\t}", jsonName, goName), true, true
	case FieldBool, FieldToggle:
		return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v == \"true\" || v == \"1\" || v == \"yes\"\n\t\t}", jsonName, goName), false, true
	default: // string, text, richtext, select
		return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok {\n\t\t\titem.%s = v\n\t\t}", jsonName, goName), false, true
	}
}

// sortableInList mirrors the handler's list Sortable set: strings and whole
// numbers. A bool is filterable and pointless to sort by.
func sortableInList(f Field) bool {
	switch f.GoType() {
	case "string", "int", "uint":
		return true
	}
	return false
}

// ── line editing ────────────────────────────────────────────────────────────

type goLines struct {
	lines []string
	crlf  bool
}

func readGoLines(path string) (*goLines, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := string(raw)
	crlf := strings.Contains(s, "\r\n")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return &goLines{lines: strings.Split(s, "\n"), crlf: crlf}, nil
}

func (l *goLines) write(path string) error {
	s := codefmt.Go(strings.Join(l.lines, "\n"))
	if l.crlf {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return os.WriteFile(path, []byte(s), 0644)
}

// find returns the first line at or after from whose trimmed text is want.
func (l *goLines) find(from int, want string) int {
	if from < 0 {
		return -1
	}
	for i := from; i < len(l.lines); i++ {
		if strings.TrimSpace(l.lines[i]) == want {
			return i
		}
	}
	return -1
}

// findPrefix is find on a prefix of the trimmed line.
func (l *goLines) findPrefix(from int, prefix string) int {
	if from < 0 {
		return -1
	}
	for i := from; i < len(l.lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(l.lines[i]), prefix) {
			return i
		}
	}
	return -1
}

// closeOf returns the line closing the block opened at open: the first later
// line that is the opener's indentation followed by "}".
func (l *goLines) closeOf(open int) int {
	if open < 0 {
		return -1
	}
	line := l.lines[open]
	indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	for i := open + 1; i < len(l.lines); i++ {
		if l.lines[i] == indent+"}" {
			return i
		}
	}
	return -1
}

func (l *goLines) insert(at int, code string) {
	add := strings.Split(code, "\n")
	out := make([]string, 0, len(l.lines)+len(add))
	out = append(out, l.lines[:at]...)
	out = append(out, add...)
	out = append(out, l.lines[at:]...)
	l.lines = out
}

// addToInlineMap puts key into a one-line map literal on line i, before its
// last closing brace.
func (l *goLines) addToInlineMap(i int, key string) {
	if i < 0 {
		return
	}
	line := l.lines[i]
	if strings.Contains(line, fmt.Sprintf("%q: true", key)) {
		return
	}
	at := strings.LastIndex(line, "}")
	if at < 0 {
		return
	}
	l.lines[i] = line[:at] + fmt.Sprintf(", %q: true", key) + line[at:]
}

func (l *goLines) contains(s string) bool {
	return strings.Contains(strings.Join(l.lines, "\n"), s)
}

// ── the injections ──────────────────────────────────────────────────────────

// injectHandlerField adds f to the resource's handler wherever the generator
// would have put it. Refuses, rather than warns, when the create request is
// missing: carrying on would reproduce exactly the silent drop this exists to
// fix.
func (g *Generator) injectHandlerField(names Names, f Field) error {
	rel := filepath.ToSlash(filepath.Join("apps", "api", "internal", "handlers", names.Snake+".go"))
	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+".go")
	l, err := readGoLines(path)
	if err != nil {
		return fmt.Errorf("reading the %s handler: %w", names.Pascal, err)
	}
	p := scalarHandlerParts(f)
	snake := toSnakeCase(f.Name)

	createOpen := l.find(0, fmt.Sprintf("type Create%sRequest struct {", names.Pascal))
	createClose := l.closeOf(createOpen)
	if createClose < 0 {
		return fmt.Errorf("could not find Create%sRequest in %s.\n\n"+
			"The field is in the model, but the API will not accept it until it is added there\n"+
			"and to Update%sRequest, the create assignment and the PATCH allow-list", names.Pascal, rel, names.Pascal)
	}
	// Already added: the whole command is safe to re-run.
	for i := createOpen; i < createClose; i++ {
		if strings.Contains(l.lines[i], fmt.Sprintf("json:\"%s\"", snake)) {
			return nil
		}
	}
	l.insert(createClose, p.createField)

	if open := l.find(0, fmt.Sprintf("type Update%sRequest struct {", names.Pascal)); open >= 0 {
		if end := l.closeOf(open); end >= 0 {
			l.insert(end, p.updateField)
		}
	}

	create := l.findPrefix(0, fmt.Sprintf("func (h *%sHandler) Create(", names.Pascal))
	if lit := l.find(create, fmt.Sprintf("item := models.%s{", names.Pascal)); lit >= 0 {
		if end := l.closeOf(lit); end >= 0 {
			l.insert(end, p.createAssign)
		}
	}

	update := l.findPrefix(0, fmt.Sprintf("func (h *%sHandler) Update(", names.Pascal))
	if m := l.find(update, "updates := map[string]interface{}{}"); m >= 0 {
		l.insert(m+1, p.updateSet)
	}

	// The PATCH allow-list, and bulk edit's copy of it.
	for _, fn := range []string{"Patch", "Bulk"} {
		start := l.findPrefix(0, fmt.Sprintf("func (h *%sHandler) %s(", names.Pascal, fn))
		if open := l.find(start, "allowed := map[string]bool{"); open >= 0 {
			if end := l.closeOf(open); end >= 0 {
				indent := l.lines[end][:len(l.lines[end])-1]
				l.insert(end, indent+"\t"+p.patchKey)
			}
		}
	}

	// The list's sort and filter whitelists.
	if sortableInList(f) {
		l.addToInlineMap(l.findPrefix(0, "Sortable:"), snake)
	}
	l.addToInlineMap(l.findPrefix(0, "Filterable:"), snake)

	if f.NeedsJSONTimeImport() && !l.contains("/internal/jsontime\"") {
		if models := l.findPrefix(0, fmt.Sprintf("%q", g.Module+"/internal/models")); models >= 0 {
			l.insert(models, fmt.Sprintf("\t%q", g.Module+"/internal/jsontime"))
		}
	}

	if err := l.write(path); err != nil {
		return fmt.Errorf("writing %s: %w", rel, err)
	}
	fmt.Printf("  ✓ %s  (create, update, patch)\n", rel)
	return nil
}

// injectImportField teaches the CSV importer the new column. Best effort: the
// importer is a convenience, and a hand-reshaped one is reported, not fatal.
func (g *Generator) injectImportField(names Names, f Field) {
	code, needStrconv, ok := importAssign(f)
	if !ok {
		return
	}
	rel := filepath.ToSlash(filepath.Join("apps", "api", "internal", "handlers", names.Snake+"_import.go"))
	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_import.go")
	l, err := readGoLines(path)
	if err != nil {
		return
	}
	snake := toSnakeCase(f.Name)
	if l.contains(fmt.Sprintf("get(rec, %q)", snake)) {
		return
	}

	lit := l.find(0, fmt.Sprintf("item := models.%s{}", names.Pascal))
	batch := l.findPrefix(lit, "batch = append(batch, pendingRow{")
	if lit < 0 || batch < 0 {
		fmt.Printf("  ⚠ %s: add %q to the CSV import by hand\n", rel, snake)
		return
	}
	l.insert(batch, code)

	// The downloadable template's header row.
	tmpl := l.findPrefix(0, fmt.Sprintf("func (h *%sHandler) Template(", names.Pascal))
	if hdr := l.findPrefix(tmpl, "c.String(http.StatusOK, \""); hdr >= 0 {
		l.lines[hdr] = strings.Replace(l.lines[hdr], `\n")`, ","+snake+`\n")`, 1)
	}

	if needStrconv && !l.contains("\"strconv\"") {
		if imp := l.find(0, "import ("); imp >= 0 {
			l.insert(imp+1, "\t\"strconv\"")
		}
	}
	if err := l.write(path); err != nil {
		fmt.Printf("  ⚠ %s: %v\n", rel, err)
		return
	}
	fmt.Printf("  ✓ %s  (CSV import)\n", rel)
}

// injectServiceSort lets the list sort by the new column. The service keeps
// its own whitelist, because sort_by goes into ORDER BY.
func (g *Generator) injectServiceSort(names Names, f Field) {
	path := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+".go")
	l, err := readGoLines(path)
	if err != nil {
		return
	}
	i := l.findPrefix(0, fmt.Sprintf("sortable%s := map[string]bool{", names.Pascal))
	if i < 0 {
		return
	}
	l.addToInlineMap(i, toSnakeCase(f.Name))
	if err := l.write(path); err != nil {
		fmt.Printf("  ⚠ services/%s.go: %v\n", names.Snake, err)
	}
}

// ensureModelImport gives the model the import a new field's type lives in.
//
// injectModelField adds the struct line and nothing else, so a date or
// datetime field put *jsontime.Date into a model that did not import jsontime
// and the project stopped building with "undefined: jsontime". Found adding
// opened_on:date to a ledger account. Every other type grit g field supports is
// a builtin.
func (g *Generator) ensureModelImport(names Names, f Field) error {
	if !f.NeedsJSONTimeImport() {
		return nil
	}
	path := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	l, err := readGoLines(path)
	if err != nil {
		return fmt.Errorf("reading the %s model: %w", names.Pascal, err)
	}
	if l.contains("/internal/jsontime\"") {
		return nil
	}
	imp := fmt.Sprintf("\t%q", g.Module+"/internal/jsontime")
	switch ids, open := l.findPrefix(0, fmt.Sprintf("%q", g.Module+"/internal/ids")), l.find(0, "import ("); {
	case ids >= 0:
		l.insert(ids+1, imp)
	case open >= 0:
		l.insert(open+1, imp)
	default:
		return fmt.Errorf("could not add the jsontime import to models/%s.go: add %s by hand", names.Snake, imp)
	}
	return l.write(path)
}
