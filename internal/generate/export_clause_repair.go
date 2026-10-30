package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Review 2.19: Export carried its own copy of the paginate search clause.
//
// buildSearchClause was unexported, so the generated Export reimplemented it.
// Two copies of a query builder drift, and this pair had already drifted into
// the same bug twice over: neither escaped LIKE wildcards, so a search for
// "50%" returned every row and "a_b" matched "axb". Fixing one copy would have
// left the other.
//
// The generator now calls paginate.BuildSearchClause. This brings an existing
// service to the same call, which matters more than it looks: the escaping is
// in the shared builder, so a service still carrying the copy keeps answering
// the wrong rows for an export after the screen beside it was fixed.
//
// Matched as a block rather than line by line, because the copy is eleven
// lines of loop and a regexp over each would match the real builder too.

var exportSearchCopyRe = regexp.MustCompile(
	`(?s)\tif search != "" \{\n` +
		`\t\tclause := ""\n` +
		`\t\targs := \[\]any\{\}\n` +
		`\t\twild := "%" \+ search \+ "%"\n` +
		`\t\tfor i, col := range (\w+)ListConfig\.Searchable \{\n` +
		`\t\t\tif i > 0 \{\n` +
		`\t\t\t\tclause \+= " OR "\n` +
		`\t\t\t\}\n` +
		`\t\t\tclause \+= "LOWER\(" \+ col \+ "\) LIKE LOWER\(\?\)"\n` +
		`\t\t\targs = append\(args, wild\)\n` +
		`\t\t\}\n` +
		`\t\tif clause != "" \{\n` +
		`\t\t\tquery = query\.Where\(clause, args\.\.\.\)\n` +
		`\t\t\}\n` +
		`\t\}\n`)

// RepairExportSearchClause points every generated Export at the shared clause
// builder, in every service under apiRoot.
func RepairExportSearchClause(apiRoot string) (changed []string, err error) {
	dir := filepath.Join(apiRoot, "internal", "services")
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		// No services directory is an --api project that has generated nothing
		// yet, not a failure.
		return nil, nil
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return changed, fmt.Errorf("reading %s: %w", path, readErr)
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")
		out, ok := repairExportSearchClauseSource(src)
		if !ok {
			continue
		}
		if err := writeRepairedGo(path, out); err != nil {
			return changed, err
		}
		changed = append(changed, filepath.ToSlash(filepath.Join("internal", "services", e.Name()))+
			": an export searches through the same clause the list screen uses, so a term with a % in it means the same thing to both")
	}
	return changed, nil
}

func repairExportSearchClauseSource(src string) (string, bool) {
	m := exportSearchCopyRe.FindStringSubmatchIndex(src)
	if m == nil {
		return src, false
	}
	camel := src[m[2]:m[3]]
	replacement := "\tif search != \"\" {\n" +
		"\t\t// The clause List uses, not a copy of it. This was a copy, and the\n" +
		"\t\t// copy had already drifted: neither escaped LIKE wildcards, so a\n" +
		"\t\t// search for \"50%\" matched every row, and fixing one would have left\n" +
		"\t\t// the other (review 2.19).\n" +
		"\t\tif clause, args := paginate.BuildSearchClause(" + camel + "ListConfig.Searchable, search); clause != \"\" {\n" +
		"\t\t\tquery = query.Where(clause, args...)\n" +
		"\t\t}\n" +
		"\t}\n"
	out := src[:m[0]] + replacement + src[m[1]:]

	// And the sentence above Export, which claimed the columns matched List's.
	// They do more than match now: it is the same clause, and the comment the
	// generator writes says so. Byte equality with a freshly generated service
	// is what TestGeneratedWritesDoNotReadTheRowBack holds this to.
	const oldDoc = "// table is never in memory at once. search matches the columns List searches.\n"
	const newDoc = "// table is never in memory at once.\n" +
		"//\n" +
		"// search goes through the same clause builder List uses, so the rows an export\n" +
		"// contains are the rows the screen showed. It used to be a second copy of that\n" +
		"// builder, which is how the two came to disagree about escaping a wildcard.\n"
	if strings.Count(out, oldDoc) == 1 {
		out = strings.Replace(out, oldDoc, newDoc, 1)
	}
	return out, true
}
