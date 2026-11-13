package scaffold

import (
	"regexp"
	"strings"
)

// RepairGeneratedServiceRowRules teaches an existing generated service to
// apply its row rules.
//
// Without this, a project that upgrades keeps services that never call
// internal/rowrules, while a resource generated after the upgrade does. The
// file still compiles, so nothing says anything: a rule written in
// internal/rowrules/rules.go simply has no effect on the resources the
// project already had, which is the worst shape a security feature can fail
// in.
//
// Four edits, matching what the generator writes today:
//
//	the import
//	the rows helper
//	List      scoped with rowrules.List
//	GetByID   scoped with rowrules.View
//	load      scoped with Update or Delete, by what is about to happen
//
// Idempotent: a service that already calls s.rows is returned untouched, so
// running the upgrade twice does not double anything.
func RepairGeneratedServiceRowRules(src, module string) (string, bool) {
	pascal := serviceTypeName(src)
	if pascal == "" {
		return src, false
	}

	// Already generated against row rules, or already repaired. Guarded on
	// the declaration rather than the name, so the sentence in this comment
	// cannot match itself.
	if strings.Contains(src, ") rows(ctx context.Context, query *gorm.DB") {
		return src, false
	}

	// load's two-argument form is what the policy repair leaves behind, and
	// this repair is written against that shape. A service still on the old
	// one-argument load is left for that repair to reach first.
	loadSig := "func (s *" + pascal + "Service) load(ctx context.Context, id, ability string) (*models." + pascal + ", error) {"
	if !strings.Contains(src, loadSig) {
		return src, false
	}

	out := src

	// 1. The import, added to the block as a line rather than spliced onto an
	//    anchor: a Replace that matches nothing fails silently and ships a
	//    file missing its import.
	importLine := "\t\"" + module + "/internal/rowrules\"\n"
	if !strings.Contains(out, importLine) {
		anchor := "\t\"" + module + "/internal/respond\"\n"
		if !strings.Contains(out, anchor) {
			return src, false
		}
		out = strings.Replace(out, anchor, anchor+importLine, 1)
	}

	// 2. The helper, in front of List where the generator puts it.
	listDoc := "// List returns one page of "
	at := strings.Index(out, listDoc)
	if at < 0 {
		return src, false
	}
	helper := "// rows applies this resource's row rules to a query.\n" +
		"//\n" +
		"// internal/rowrules holds a condition per action, written as an expression\n" +
		"// over this model's own fields, and this is where it reaches the SQL. A\n" +
		"// resource with no rule declared gets its query back untouched, which is how\n" +
		"// the project behaved before it had any.\n" +
		"//\n" +
		"// Applied to the query rather than checked afterwards: for a list that filters\n" +
		"// it, and for one record it means the row is not found. That is a 404 rather\n" +
		"// than a 403, and deliberately so, because a 403 on a row you may not read\n" +
		"// confirms the row exists.\n" +
		"func (s *" + pascal + "Service) rows(ctx context.Context, query *gorm.DB, action rowrules.Action) (*gorm.DB, error) {\n" +
		"\treturn rowrules.Scope(ctx, query, \"" + pascal + "\", action, &models." + pascal + "{})\n" +
		"}\n\n"
	out = out[:at] + helper + out[at:]

	// 3. List.
	listReturn := regexp.MustCompile(`\n\treturn paginate\.List\[models\.` + pascal + `\]\(query, p, (\w+)\)\n`)
	m := listReturn.FindStringSubmatch(out)
	if m == nil {
		return src, false
	}
	scoped := "\n\tquery, err := s.rows(ctx, query, rowrules.List)\n" +
		"\tif err != nil {\n" +
		"\t\treturn paginate.Result[models." + pascal + "]{}, err\n" +
		"\t}\n" +
		"\n\treturn paginate.List[models." + pascal + "](query, p, " + m[1] + ")\n"
	out = strings.Replace(out, m[0], scoped, 1)

	// 4. GetByID.
	getOld := regexp.MustCompile(
		`\tvar item models\.` + pascal + `\n\tif err := s\.db\(ctx\)((?:\.Preload\("[^"]+"\))*)\.First\(&item, "id = \?", id\)\.Error; err != nil \{`)
	if loc := getOld.FindStringSubmatchIndex(out); loc != nil {
		preloads := out[loc[2]:loc[3]]
		replacement := "\tquery, err := s.rows(ctx, s.db(ctx)" + preloads + ", rowrules.View)\n" +
			"\tif err != nil {\n" +
			"\t\treturn nil, err\n" +
			"\t}\n" +
			"\tvar item models." + pascal + "\n" +
			"\tif err := query.First(&item, \"id = ?\", id).Error; err != nil {"
		out = out[:loc[0]] + replacement + out[loc[1]:]
	} else {
		return src, false
	}

	// 5. load, which both the update and the delete go through.
	loadOld := loadSig + "\n\tvar item models." + pascal + "\n" +
		"\tif err := s.db(ctx).First(&item, \"id = ?\", id).Error; err != nil {"
	if !strings.Contains(out, loadOld) {
		return src, false
	}
	loadNew := loadSig + "\n" +
		"\t// The row rule for the action about to happen, so a row this caller may\n" +
		"\t// not change is not found rather than refused.\n" +
		"\taction := rowrules.Update\n" +
		"\tif ability == \"delete\" {\n" +
		"\t\taction = rowrules.Delete\n" +
		"\t}\n" +
		"\tquery, err := s.rows(ctx, s.db(ctx), action)\n" +
		"\tif err != nil {\n" +
		"\t\treturn nil, err\n" +
		"\t}\n" +
		"\tvar item models." + pascal + "\n" +
		"\tif err := query.First(&item, \"id = ?\", id).Error; err != nil {"
	out = strings.Replace(out, loadOld, loadNew, 1)

	return out, true
}

// serviceTypeName reads the resource a generated service is for, from the
// type it declares.
var serviceTypeRe = regexp.MustCompile(`(?m)^type (\w+)Service struct \{`)

func serviceTypeName(src string) string {
	m := serviceTypeRe.FindStringSubmatch(src)
	if m == nil {
		return ""
	}
	return m[1]
}
