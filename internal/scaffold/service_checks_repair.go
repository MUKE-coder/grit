package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Bringing an existing service up to the v3.451.0 write checks.
//
// Enterprise review 2.4: Patch, BulkEdit and Bulk applied a column whitelist
// and wrote the map. The whitelist answers "may this column be written at
// all" and nothing else, so binding:"required" never fired on those paths and
// a wrong type reached the database as a 500; and the bulk paths loaded their
// rows with Find, so any rule in internal/policies was evadable by selecting
// the rows and using a bulk action.
//
// The generator writes the fix for new resources. This is how it reaches a
// service that already exists, which is every resource anybody has generated.
// TestGeneratedWritesDoNotReadTheRowBack compares what this produces against
// what the generator produces, byte for byte, so the two cannot drift.

var checksServiceRe = regexp.MustCompile(`func \(s \*(\w+)Service\) Patch\(ctx context\.Context`)

// RepairGeneratedServiceChecks adds the validation and policy checks to a
// generated service that predates them.
func RepairGeneratedServiceChecks(src, module string) (string, bool) {
	m := checksServiceRe.FindStringSubmatch(src)
	if m == nil {
		return src, false
	}
	pascal := m[1]

	// Already done, by this repair or by a newer generator.
	if strings.Contains(src, "func (s *"+pascal+"Service) checked(") {
		return src, false
	}

	plural := serviceOwnsPlural(src, pascal)
	if plural == "" {
		return src, false
	}

	out := src
	var ok bool

	if out, ok = addChecksImports(out); !ok {
		return src, false
	}
	if out, ok = addChecksHelpers(out, pascal, plural); !ok {
		return src, false
	}
	if out, ok = checkPatch(out, pascal); !ok {
		return src, false
	}
	if out, ok = checkBulk(out, pascal, plural); !ok {
		return src, false
	}
	// permitted calls authz.Inspect. The policy repair runs before this one and
	// brings the import in, but a repair that depends on another having run is
	// one bad ordering away from writing a project that does not compile.
	if out, ok = addAuthzImport(out, module); !ok {
		return src, false
	}
	// BulkEdit is absent from a service generated before the grid shipped.
	// Not having one is not a failure; having one and leaving it unchecked
	// would be.
	out = checkBulkEdit(out, pascal)

	return out, true
}

// addAuthzImport brings in internal/authz if nothing else has.
func addAuthzImport(src, module string) (string, bool) {
	line := "\t\"" + module + "/internal/authz\"\n"
	if strings.Contains(src, line) {
		return src, true
	}
	anchor := "\t\"" + module + "/internal/concurrency\"\n"
	if strings.Count(src, anchor) != 1 {
		return src, false
	}
	return strings.Replace(src, anchor, line+anchor, 1), true
}

// addChecksImports adds the two stdlib packages the helpers need.
func addChecksImports(src string) (string, bool) {
	const anchor = "\t\"encoding/json\"\n"
	if strings.Count(src, anchor) != 1 {
		return src, false
	}
	add := "\t// checked reads the decode failure to name the field and the type it\n" +
		"\t// wanted, rather than letting a wrong type reach the database as a 500.\n" +
		"\t\"errors\"\n"
	out := strings.Replace(src, anchor, anchor+add, 1)

	// strings goes after fmt, where gofmt would put it.
	const fmtLine = "\t\"fmt\"\n"
	if strings.Count(out, fmtLine) != 1 {
		return src, false
	}
	return strings.Replace(out, fmtLine, fmtLine+"\t\"strings\"\n", 1), true
}

// addChecksHelpers inserts checked, readableType and permitted.
func addChecksHelpers(src, pascal, plural string) (string, bool) {
	const anchor = "// conflict is the answer to a write whose precondition failed: the version\n"
	if strings.Count(src, anchor) != 1 {
		return src, false
	}
	return strings.Replace(src, anchor, serviceChecksHelpers(pascal, plural)+anchor, 1), true
}

// serviceChecksHelpers is the same text the generator emits. Kept here rather
// than imported from internal/generate because that package imports this one,
// and the comparator test is what keeps the two identical.
func serviceChecksHelpers(pascal, plural string) string {
	return fmt.Sprintf(`// checked applies a whitelisted patch onto a copy of the row and runs the
// model's own rules over the result.
//
// The whitelist answers "may this column be written at all". It does not
// answer "is this a legal value", and before v3.450.0 nothing did on this
// path: binding:"required" fired on Create and Update and nowhere else, so
// {"seats":"x"} reached the database as a type error and came back a 500
// rather than a 422 naming the field. See review 2.4.
//
// Onto a COPY of the current row, not onto an empty struct, so a rule about a
// field the patch does not mention still sees the value already stored. The
// copy is thrown away; the write still goes through Updates with the map, so
// a column the caller did not send is not overwritten with a zero.
//
// Before expand%sEmbedded, deliberately: that rewrites a money field
// into its two columns, and those names are not json tags, so a decode after
// it would silently drop them.
func (s *%sService) checked(current *models.%s, updates map[string]interface{}) error {
	raw, err := json.Marshal(updates)
	if err != nil {
		return respond.Rule("could not read the fields to write")
	}
	patched := *current
	if err := json.Unmarshal(raw, &patched); err != nil {
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) && typeErr.Field != "" {
			return respond.Rule(fmt.Sprintf("%%s must be %%s", typeErr.Field, s.readableType(typeErr.Type.String())))
		}
		return respond.Rule("one of the fields is the wrong type")
	}
	return respond.ValidateStruct(&patched)
}

// readableType turns a Go type into something worth showing a person.
//
// A method rather than a function, because the generator writes one of these
// per resource: a package-level helper of the same name twice over is a
// project with two generated resources that does not compile.
func (s *%sService) readableType(t string) string {
	switch {
	case strings.HasPrefix(t, "int"), strings.HasPrefix(t, "uint"):
		return "a whole number"
	case strings.HasPrefix(t, "float"):
		return "a number"
	case t == "bool":
		return "true or false"
	case t == "string":
		return "text"
	case strings.HasPrefix(t, "time."):
		return "a date"
	default:
		return "of type " + t
	}
}

// permitted keeps only the rows a policy allows this caller to act on.
//
// The single-row paths go through s.load, which runs authz.Inspect. The bulk
// paths loaded with Find and wrote, so a rule in internal/policies was
// evadable by selecting the rows and using a bulk action: the route's
// permission was the only check (review 2.4).
//
// Denied rows are reported rather than silently skipped. A bulk action that
// says it changed twelve rows and changed nine is worse than one that refuses
// and says which.
func (s *%sService) permitted(ctx context.Context, rows []models.%s, ability string) ([]models.%s, []string) {
	allowed := make([]models.%s, 0, len(rows))
	var denied []string
	for i := range rows {
		if d := authz.Inspect(ctx, "%s."+ability, &rows[i]); !d.Allowed() {
			denied = append(denied, rows[i].ID)
			continue
		}
		allowed = append(allowed, rows[i])
	}
	return allowed, denied
}

`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, plural)
}

// checkPatch validates the patch before anything is written, inside Patch only.
//
// Anchored on the whitelist guard rather than on expand<X>Embedded: the money
// expansion is itself a repair, so a project that has not had it yet has no
// such line, and the check would be skipped without a word.
func checkPatch(src, pascal string) (string, bool) {
	sig := "func (s *" + pascal + "Service) Patch(ctx context.Context"
	body, ok := methodBody(src, sig)
	if !ok {
		return src, false
	}
	const anchor = "\tif len(updates) == 0 && !touchedRelations {\n" +
		"\t\treturn nil, nil, respond.Rule(\"No writable fields in request body\")\n" +
		"\t}\n\n"
	if strings.Count(body, anchor) != 1 {
		return src, false
	}
	patched := strings.Replace(body, anchor, anchor+
		"\t// The model's own rules, before anything is written. See checked.\n"+
		"\tif err := s.checked(item, updates); err != nil {\n"+
		"\t\treturn nil, nil, err\n"+
		"\t}\n\n", 1)
	return strings.Replace(src, body, patched, 1), true
}

// checkBulk adds the policy pass and the per-row validation to Bulk.
func checkBulk(src, pascal, plural string) (string, bool) {
	sig := "func (s *" + pascal + "Service) Bulk(ctx context.Context"
	body, ok := methodBody(src, sig)
	if !ok {
		return src, false
	}

	const loop = "\tfor _, item := range items {\n"
	if strings.Count(body, loop) != 1 {
		return src, false
	}
	gate := fmt.Sprintf(`	// A policy rule narrows this further than the route's permission did.
	// Without it, any rule in internal/policies was evadable by selecting the
	// rows and using a bulk action (review 2.4).
	//
	// The ability is the one the action performs, so a policy that allows
	// archiving and refuses deleting is honoured by both.
	ability := "update"
	switch action {
	case "delete":
		ability = "delete"
	case "archive", "restore":
		ability = "update"
	}
	items, denied := s.permitted(ctx, items, ability)
	if len(denied) > 0 {
		// Refused whole rather than partly applied. A bulk action that
		// reports twelve and changed nine is worse than one that says no.
		//
		// As a denial, not a rule violation: this is 403, the same answer the
		// single-row path gives, and a client that branches on the status
		// should not see one code for one row and another for twelve.
		return %sBulkResult{}, authz.Denied("%s."+ability, authz.Deny(fmt.Sprintf(
			"not allowed to %%s %%d of the selected %s", action, len(denied))))
	}

`, pascal, plural, plural)
	patched := strings.Replace(body, loop, gate+loop, 1)

	rule := "\t\t\treturn " + pascal + "BulkResult{}, respond.Rule(\"No writable fields in patch\")\n\t\t}\n"
	if strings.Count(patched, rule) != 1 {
		return src, false
	}
	perRow := rule[:len(rule)-len("\t\t}\n")] +
		fmt.Sprintf(`		}
		// And the model's own rules, against each row as it will be. One
		// patch over many rows can still be invalid for some of them: a
		// conditional field's requirement depends on the row it lands on.
		for i := range items {
			if err := s.checked(&items[i], result.Updates); err != nil {
				return %sBulkResult{}, err
			}
		}
`, pascal)
	patched = strings.Replace(patched, rule, perRow, 1)

	return strings.Replace(src, body, patched, 1), true
}

// checkBulkEdit brings the grid's save up to the same two checks.
//
// Four edits, not two: the policy pass returns the rows it refused, and those
// have to be reported by index so the grid can put the message on the row.
// Adding the call and ignoring its second return is a service that does not
// compile, which is how this was found.
//
// A service generated before the grid shipped has no BulkEdit at all. That is
// not a failure; having one and leaving it unchecked would be.
func checkBulkEdit(src, pascal string) string {
	sig := "func (s *" + pascal + "Service) BulkEdit(ctx context.Context"
	body, ok := methodBody(src, sig)
	if !ok {
		return src
	}
	if strings.Contains(body, "s.permitted(ctx, allowed,") {
		return src
	}

	// The singular the generator uses in a row error, which is Names.Lower:
	// the Pascal name lower-cased. NOT the plural with its "s" trimmed, which
	// turns companies into "companie" and leaves the repair matching nothing.
	lower := strings.ToLower(pascal)

	// 1. The money expansion leaves the build loop: s.checked decodes these
	// keys as json tags, and a money field's two columns are not that.
	edits := [][2]string{{
		"\t\texpand" + pascal + "Embedded(updates)\n" +
			"\t\twrites = append(writes, write{id: edit.ID, updates: updates})\n",
		"\t\t// Not expanded yet: s.checked below decodes these keys as json tags,\n" +
			"\t\t// and a money field's two columns are not that.\n" +
			"\t\twrites = append(writes, write{id: edit.ID, updates: updates})\n",
	}}

	// 2. The policy pass, right after the rows are loaded.
	edits = append(edits, [2]string{
		"\texists := make(map[string]bool, len(allowed))\n",
		"\t// A policy rule narrows this further than the route's permission did.\n" +
			"\t// Without it, a rule in internal/policies was evadable by selecting the\n" +
			"\t// rows in the grid and editing them here: Patch consults the policy\n" +
			"\t// through s.load and this path did not (review 2.4).\n" +
			"\tpermittedRows, denied := s.permitted(ctx, allowed, \"update\")\n" +
			"\tallowed = permittedRows\n\n" +
			"\texists := make(map[string]bool, len(allowed))\n",
	})

	// 3. A refused row is reported by index, beside a missing one. Without
	// this the policy pass would drop rows and the grid would be told they
	// saved.
	edits = append(edits, [2]string{
		"\tfor i, w := range writes {\n" +
			"\t\tif !exists[w.id] {\n" +
			"\t\t\trowErrors = append(rowErrors, " + pascal + "RowError{Index: i, Message: \"no such " + lower + "\"})\n" +
			"\t\t}\n\t}\n",
		"\tdeniedSet := make(map[string]bool, len(denied))\n" +
			"\tfor _, id := range denied {\n" +
			"\t\tdeniedSet[id] = true\n\t}\n" +
			"\tfor i, w := range writes {\n" +
			"\t\tswitch {\n" +
			"\t\tcase deniedSet[w.id]:\n" +
			"\t\t\trowErrors = append(rowErrors, " + pascal + "RowError{Index: i, Message: \"not allowed to edit this " + lower + "\"})\n" +
			"\t\tcase !exists[w.id]:\n" +
			"\t\t\trowErrors = append(rowErrors, " + pascal + "RowError{Index: i, Message: \"no such " + lower + "\"})\n" +
			"\t\t}\n\t}\n",
	})

	// 4. The model's own rules per row, then the expansion the build loop no
	// longer does.
	edits = append(edits, [2]string{
		"\tsaved := make([]string, 0, len(writes))\n",
		"\t// And the model's own rules, per row, against the row as it will be.\n" +
			"\t// Reported by index so a grid can put the message on the cell.\n" +
			"\tcurrent := make(map[string]*models." + pascal + ", len(allowed))\n" +
			"\tfor i := range allowed {\n" +
			"\t\tcurrent[allowed[i].ID] = &allowed[i]\n\t}\n" +
			"\tfor i, w := range writes {\n" +
			"\t\trow, ok := current[w.id]\n" +
			"\t\tif !ok {\n\t\t\tcontinue\n\t\t}\n" +
			"\t\tif err := s.checked(row, w.updates); err != nil {\n" +
			"\t\t\trowErrors = append(rowErrors, " + pascal + "RowError{Index: i, Message: err.Error()})\n" +
			"\t\t}\n\t}\n" +
			"\tif len(rowErrors) > 0 {\n\t\treturn nil, rowErrors, nil\n\t}\n\n" +
			"\t// Expanded after validation: a money field becomes two columns whose\n" +
			"\t// names are not json tags, so a decode after this drops them.\n" +
			"\tfor i := range writes {\n" +
			"\t\texpand" + pascal + "Embedded(writes[i].updates)\n\t}\n\n" +
			"\tsaved := make([]string, 0, len(writes))\n",
	})

	// All four or none. A BulkEdit that has drifted from the template gets the
	// policy pass alone otherwise, and that is the shape that does not build.
	patched := body
	for _, e := range edits {
		if strings.Count(patched, e[0]) != 1 {
			return src
		}
		patched = strings.Replace(patched, e[0], e[1], 1)
	}
	return strings.Replace(src, body, patched, 1)
}

// RepairServiceChecks applies the write checks to every generated service in
// the project the command is run in.
//
// Called from the CLI after the money repair rather than from inside
// scaffold.Upgrade, because that repair is what adds the grid methods and the
// encoding/json import these checks are written against. Run first, this would
// check a Patch and a Bulk and then watch a BulkEdit arrive unchecked.
func RepairServiceChecks() error {
	root, err := FindProjectRoot()
	if err != nil {
		return err
	}
	opts := Options{
		Architecture:     detectArchitecture(root),
		LegacySingleFlat: dirExists(filepath.Join(root, "internal", "routes")),
	}
	opts.Normalize()
	return repairGeneratedServiceChecks(root, opts)
}

// repairGeneratedServiceChecks applies the write checks to every generated
// service under root.
func repairGeneratedServiceChecks(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module, err := detectModule(apiRoot)
	if err != nil || module == "" {
		// No module path means the authz import cannot be written, and a repair
		// that adds a call without its import is a project that does not build.
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	services, err := goSources(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	for _, path := range services {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, changed := RepairGeneratedServiceChecks(src, module)
			if !changed {
				return src, nil, nil
			}
			return out, []string{"a patch or a bulk action is checked against the model's own rules and the row policy, instead of writing whatever the whitelist let through"}, nil
		}); err != nil {
			return err
		}
	}
	return nil
}
