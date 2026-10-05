package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Repairing a generated service that cannot update a money column.
//
// The bug: money.Money is embedded, so the table holds <field>_amount and
// <field>_currency rather than <field>. GORM expands an embedded struct when it
// writes a MODEL, which is why Create always worked, and does not when it
// writes a MAP: there a key is a column name and its value is one bind
// argument, and a struct is not one. Every path that updates through a map
// (Update, Patch, Bulk and BulkEdit) answered
//
//	sql: converting argument $N type: unsupported type money.Money, a struct
//
// and so a record with a price could not be edited at all, from the admin or
// from the API.
//
// Why this is a repair and not a template fix alone: internal/services/<x>.go
// is generated resource code, which `grit upgrade` deliberately preserves. A
// fix to the generator reaches new projects and new resources, and leaves every
// existing shop broken. Regenerating the resource would fix it and throw away
// whatever the project has added to that service since.
//
// So this walks the generated services, works out which columns are money by
// reading the model beside them, and inserts exactly what the generator now
// emits. A service that has already been repaired, or whose text is not what
// the generator wrote, is left alone and reported.

var (
	// `var writableProduct = map[string]bool{` names the resource.
	writableDecl = regexp.MustCompile(`(?m)^var writable([A-Za-z0-9_]+) = map\[string\]bool\{`)
	// `Price money.Money ` + "`" + `gorm:"embedded..." json:"price"` in the model.
	moneyFieldDecl = regexp.MustCompile("(?m)^\\s*([A-Za-z0-9_]+)\\s+money\\.Money\\b[^`]*`([^`]*)`")
	jsonTagValue   = regexp.MustCompile(`json:"([^",]+)`)
)

// RepairMoneyUpdates fixes every generated service in apiRoot that writes a
// money column through a map. It returns one line per service changed, and one
// per service it could not change and why.
func RepairMoneyUpdates(apiRoot string) (changed []string, blocked []string, err error) {
	dir := filepath.Join(apiRoot, "internal", "services")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// No services directory is not a failure: an upgrade runs against
		// projects of every shape.
		return nil, nil, nil
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		src := string(raw)

		match := writableDecl.FindStringSubmatch(src)
		if match == nil {
			continue // not a generated resource service
		}
		pascal := match[1]

		fields, ferr := moneyFieldsOf(apiRoot, pascal)
		if ferr != nil {
			continue
		}

		next, note := repairMoneyUpdateSource(src, pascal, fields)
		if note != "" {
			blocked = append(blocked, "internal/services/"+name+": "+note)
			continue
		}
		if next == src {
			continue
		}
		if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
			return changed, blocked, fmt.Errorf("writing %s: %w", path, err)
		}
		what := "no money columns, but every map update now goes through one helper"
		if len(fields) > 0 {
			what = strings.Join(fields, ", ") + " can be updated again"
		}
		changed = append(changed, "internal/services/"+name+": "+what)

		// The handler and the routes that reach these methods. Without them
		// the service grows two methods nothing calls, and the admin's Bulk
		// Create button posts to a path that does not exist.
		snake := strings.TrimSuffix(name, ".go")
		lower := strings.ToLower(pascal)
		plural := Pluralize(snake)
		if done, herr := repairGridHandler(apiRoot, snake, pascal, lower, plural); herr != nil {
			return changed, blocked, herr
		} else if done {
			changed = append(changed, "internal/handlers/"+name+": the bulk-create and bulk-edit endpoints")
		}
		if done, rerr := repairGridRoutes(apiRoot, snake, plural); rerr != nil {
			return changed, blocked, rerr
		} else if done {
			changed = append(changed, "internal/routes/"+snake+"_routes.go: the two grid routes")
		}
	}
	return changed, blocked, nil
}

// RepairMoneyUpdatesIn is RepairMoneyUpdates for one service, as text: the
// generated service and the model beside it. Returns the service unchanged when
// there is nothing to do, and the reason when it will not touch it.
//
// Exported so the repair can be tested in the same chain an upgrade runs it in.
// A repair proved only against the filesystem is a repair whose output nothing
// compares with what the generator writes.
func RepairMoneyUpdatesIn(serviceSrc, modelSrc string) (string, string) {
	match := writableDecl.FindStringSubmatch(serviceSrc)
	if match == nil {
		return serviceSrc, "" // not a generated resource service
	}
	pascal := match[1]
	return repairMoneyUpdateSource(serviceSrc, pascal, moneyFieldsIn(modelSrc))
}

// moneyFieldsIn returns the JSON names of a model's money columns, in
// declaration order.
func moneyFieldsIn(modelSrc string) []string {
	var fields []string
	for _, m := range moneyFieldDecl.FindAllStringSubmatch(modelSrc, -1) {
		if tag := jsonTagValue.FindStringSubmatch(m[2]); tag != nil {
			fields = append(fields, tag[1])
		} else {
			fields = append(fields, toSnakeCase(m[1]))
		}
	}
	return fields
}

// repairGridHandler appends the grid's two handler methods to a generated
// handler that predates them.
func repairGridHandler(apiRoot, snake, pascal, lower, plural string) (bool, error) {
	path := filepath.Join(apiRoot, "internal", "handlers", snake+".go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	src := string(raw)
	if !strings.Contains(src, "func (h *"+pascal+"Handler) Bulk(") {
		return false, nil // not a generated resource handler
	}
	if strings.Contains(src, "func (h *"+pascal+"Handler) BulkCreate(") {
		return false, nil
	}
	out := strings.TrimRight(src, "\n") + "\n\n" + gridHandlerMethodsSource(pascal, lower, plural)
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}

// repairGridRoutes adds the two routes, each shaped like the one it belongs
// beside.
//
// bulk-create is a create and goes after Create; bulk-edit is an edit and goes
// after Patch. Copying the neighbouring line carries whatever guard the project
// uses: a role, a permission, or the admin group.
func repairGridRoutes(apiRoot, snake, plural string) (bool, error) {
	path := filepath.Join(apiRoot, "internal", "routes", snake+"_routes.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	src := string(raw)
	if strings.Contains(src, "h.BulkCreate)") {
		return false, nil
	}

	out := src
	for _, pair := range []struct{ from, path, to string }{
		{"h.Create)", "/" + plural + "/bulk-create", "h.BulkCreate)"},
		{"h.Patch)", "/" + plural + "/bulk-edit", "h.BulkEdit)"},
	} {
		line := routeLineFor(out, pair.from)
		if line == "" {
			continue // an append-only resource has no Patch, and wants no bulk-edit
		}
		// The same line with the handler and the path swapped. The verb is
		// already POST for Create; for Patch it has to become POST, because a
		// grid of new values is a post of rows and not a patch of one.
		next := strings.Replace(line, pair.from, pair.to, 1)
		next = routeWithPath(next, pair.path)
		next = strings.Replace(next, ".PATCH(", ".POST(", 1)
		out = strings.Replace(out, line, line+"\n"+next, 1)
	}
	if out == src {
		return false, nil
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}

// routeLineFor returns the whole line registering handler, or "".
func routeLineFor(src, handler string) string {
	for _, line := range strings.Split(src, "\n") {
		if strings.Contains(line, handler) {
			return line
		}
	}
	return ""
}

// routeWithPath swaps the quoted path in a route line for a new one.
var routePath = regexp.MustCompile(`"\/[^"]*"`)

func routeWithPath(line, path string) string {
	return routePath.ReplaceAllString(line, strconv.Quote(path))
}

// moneyFieldsOf reads the model beside a service and returns the JSON names of
// its money columns, in declaration order.
//
// From the model rather than from the service, because the service never names
// the type: the whitelist holds column names and nothing says which of them is
// two columns underneath.
func moneyFieldsOf(apiRoot, pascal string) ([]string, error) {
	path := filepath.Join(apiRoot, "internal", "models", toSnakeCase(pascal)+".go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return moneyFieldsIn(string(raw)), nil
}

// repairMoneyUpdateSource inserts the helper and its four call sites, and
// returns the reason when it cannot.
func repairMoneyUpdateSource(src, pascal string, fields []string) (string, string) {
	call := "expand" + pascal + "Embedded"
	out := src

	// Each insertion guards itself: a service may have had one pass of this
	// and not another, and a resource generated between two releases has the
	// expander already and not the grid methods.
	if !strings.Contains(out, "func "+call+"(") {

		// Above the whitelist, where the generator puts it.
		anchor := "// writable" + pascal + " is every column Patch and Bulk may write."
		if strings.Count(out, anchor) != 1 {
			return src, "the writable-columns comment is not the one Grit wrote"
		}
		out = strings.Replace(out, anchor, expandEmbeddedSource(pascal, fields)+anchor, 1)
	}

	if !strings.Contains(out, call+"(updates)") {

		// Update. Anchored on the method, not on "db := s.db(ctx)": that line is
		// in Create and Delete too, and the first one wins a Replace.
		updateSig := "func (s *" + pascal + "Service) Update(ctx context.Context, id string, updates map[string]interface{}"
		i := strings.Index(out, updateSig)
		if i < 0 {
			return src, "no Update method of the shape Grit generates"
		}
		dbLine := "\tdb := s.db(ctx)\n"
		if j := strings.Index(out[i:], dbLine); j >= 0 {
			at := i + j + len(dbLine)
			out = out[:at] + "\t" + call + "(updates)\n" + out[at:]
		}

		// Patch, found from its own method for the same reason.
		patchSig := "func (s *" + pascal + "Service) Patch(ctx context.Context"
		if pi := strings.Index(out, patchSig); pi >= 0 {
			guard := "\tif len(updates) > 0 {\n"
			if j := strings.Index(out[pi:], guard); j >= 0 {
				at := pi + j
				out = out[:at] + "\t" + call + "(updates)\n" + out[at:]
			}
		}

		// Bulk's patch action.
		bulkAnchor := "\t\tcase \"patch\":\n"
		if strings.Count(out, bulkAnchor) == 1 {
			out = strings.Replace(out, bulkAnchor, bulkAnchor+"\t\t\t"+call+"(result.Updates)\n", 1)
		}

		// BulkEdit, where it exists: the grid editors are newer than some projects.
		gridAnchor := "\t\twrites = append(writes, write{id: edit.ID, updates: updates})\n"
		if strings.Count(out, gridAnchor) == 1 {
			out = strings.Replace(out, gridAnchor, "\t\t"+call+"(updates)\n"+gridAnchor, 1)
		}

	}

	// The grid's two methods, for a service generated before they existed. The
	// admin upgrade delivers the Bulk Create button and the bulk-edit grid to
	// every project, and a button whose endpoint does not exist is worse than
	// no button.
	if !strings.Contains(out, "func (s *"+pascal+"Service) BulkCreate(") {
		lower := strings.ToLower(pascal)
		plural := Pluralize(toSnakeCase(pascal))
		out = strings.TrimRight(out, "\n") + "\n\n" + gridServiceMethodsSource(pascal, lower, plural)
	}

	if len(fields) > 0 && !strings.Contains(out, "/internal/money\"") {
		var ok bool
		out, ok = addMoneyImport(out)
		if !ok {
			return src, "could not add the money import"
		}
	}
	if strings.Contains(out, "json.Marshal(row)") && !strings.Contains(out, "\"encoding/json\"") {
		var ok bool
		out, ok = addJSONImport(out)
		if !ok {
			return src, "could not add the encoding/json import"
		}
	}
	return out, ""
}

// addMoneyImport puts internal/money beside internal/models, where the
// generator puts it.
func addMoneyImport(src string) (string, bool) {
	marker := regexp.MustCompile(`(?m)^\t"([^"]+)/internal/models"$`)
	loc := marker.FindStringSubmatchIndex(src)
	if loc == nil {
		return src, false
	}
	module := src[loc[2]:loc[3]]
	line := "\t\"" + module + "/internal/money\"\n"
	// AFTER models, because that is where it sorts: "models" < "money" on the
	// third letter. Putting it before left every repaired service failing
	// gofmt, which CI checks.
	at := loc[1] + 1
	return src[:at] + line + src[at:], true
}

// addJSONImport puts encoding/json in the standard-library group, where the
// generator puts it, with the generator's comment: a repaired file has to
// come out byte-identical to a freshly generated one.
func addJSONImport(src string) (string, bool) {
	const anchor = "import (\n	\"context\"\n"
	if !strings.Contains(src, anchor) {
		return src, false
	}
	return strings.Replace(src, anchor, anchor+"	// Grid rows arrive as maps and are decoded into the model here, through\n	// the same JSON path a single create takes.\n	\"encoding/json\"\n", 1), true
}
