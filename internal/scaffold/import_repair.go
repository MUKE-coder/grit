package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// oldImportNameLookupRe is a belongs_to column resolved by natural key as grit
// generate wrote it before v3.257.0: a query per row, and any error, not only
// "not found", taken as a reason to create the record, unchecked.
var oldImportNameLookupRe = regexp.MustCompile(`\t\tif v, ok := get\(rec, "(\w+)"\); ok && v != "" \{\n` +
	`\t\t\tvar rel models\.(\w+)\n` +
	`\t\t\tif err := db\.Where\("(\w+) = \?", v\)\.First\(&rel\)\.Error; err != nil \{\n` +
	`\t\t\t\trel = models\.(\w+)\{(\w+): v\}\n` +
	`\t\t\t\tdb\.Create\(&rel\)\n` +
	`\t\t\t\}\n` +
	`\t\t\titem\.(\w+) = (&?)rel\.ID\n` +
	`\t\t\}\n`)

const (
	importCountersAnchor = "\tcreated, skipped, failed := 0, 0, 0\n"
	importOpenAnchor     = "\tf, err := os.Open(path)\n"
)

// importLimitGateDirect is ImportLimitGate for an importer generated before
// v3.251.0, which has no record helper to report the wait through.
const importLimitGateDirect = `	// Imports take turns: each holds a database connection and writes in
	// batches for its whole run, and a burst of them drained the pool every
	// request shares. See internal/imports.
	release := imports.Wait(func() {
		if err := db.Model(&models.ImportJob{}).Where("id = ?", jobID).Update("message", "Waiting for another import to finish").Error; err != nil {
			log.Printf("import job %s: recording that it is waiting: %v", jobID, err)
		}
	})
	defer release()

`

// repairImportServices brings generated CSV importers up to the fix for H19 in
// the contact-app review: a belongs_to column cost a query per row, a lookup
// error created a record without checking the create, and any number of imports
// ran at once.
//
// internal/imports is framework code and arrives whole. The importers are
// generated but the developer's, so each change anchors on what grit generate
// wrote and is skipped, with a note, where that text has changed.
func repairImportServices(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if err := repairWAFImportExclusion(root, opts); err != nil {
		return err
	}
	if !fileContains(filepath.Join(apiRoot, "internal", "imports", "limit.go"), "func Wait(") {
		return nil
	}
	services, err := goSources(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	for _, path := range services {
		if !strings.HasSuffix(path, "_import.go") {
			continue
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			return repairImportServiceSource(src, module)
		}); err != nil {
			return err
		}
	}
	return nil
}

// wafImportExclusion steps Sentinel's WAF aside for every generated import, in
// wafExcludedRoutes. Its 1 MB body cap refused any CSV larger than that with 413
// WAF_BODY_TOO_LARGE before the importer ran, which the body limit an import
// route gets (100 MB, v3.252.0) never reached.
const wafImportExclusion = `		// Every generated resource's CSV import. A spreadsheet is over the WAF's
		// body cap, and its cells trip the injection heuristics. "*" spans the
		// one segment naming the resource.
		"/*/import",
`

const wafUploadsEntry = "\t\t\"/uploads\", \"/uploads/*\",\n"

// repairWAFImportExclusion adds the imports to routes.go's WAF exclusions.
func repairWAFImportExclusion(root string, opts Options) error {
	routes := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairWAFImportExclusionSource)
}

func repairWAFImportExclusionSource(src string) (string, []string, []string) {
	const decl = "func wafExcludedRoutes() []string {"
	start := strings.Index(src, decl)
	if start < 0 || strings.Contains(src, `"/*/import"`) {
		return src, nil, nil
	}
	end := strings.Index(src[start:], "\n}\n")
	pos := strings.Index(src[start:], wafUploadsEntry)
	if end < 0 || pos < 0 || pos > end {
		return src, nil, []string{`Sentinel's WAF refuses any CSV import over 1 MB, and wafExcludedRoutes is not the one Grit wrote: add "/*/import" to its paths`}
	}
	at := start + pos + len(wafUploadsEntry)
	return src[:at] + wafImportExclusion + src[at:],
		[]string{"CSV imports over 1 MB are no longer refused by the WAF's body cap"}, nil
}

// addImportLikeGenerate puts path in an importer's import block where grit
// generate puts it, so an upgraded importer reads like a fresh one. A block
// without those neighbours gets the import as a group of its own.
func addImportLikeGenerate(content, path, module string) (string, bool) {
	start := strings.Index(content, "\nimport (\n")
	if start < 0 {
		return content, false
	}
	end := strings.Index(content[start:], "\n)\n")
	if end < 0 {
		return content, false
	}
	block := content[start : start+end+1]
	line := "\t\"" + path + "\"\n"
	if strings.Contains(block, "\n"+line) {
		return content, true
	}
	anchor, after := "", false
	switch path {
	case "errors":
		anchor, after = "\t\"encoding/json\"\n", true
	case "gorm.io/gorm":
		anchor = "\t\"gorm.io/gorm/clause\"\n"
	case module + "/internal/imports":
		anchor = "\t\"" + module + "/internal/models\"\n"
	}
	if i := strings.Index(block, "\n"+anchor); anchor != "" && i >= 0 {
		at := start + i + 1
		if after {
			at += len(anchor)
		}
		return content[:at] + line + content[at:], true
	}
	return addImportGroup(content, path)
}

func repairImportServiceSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, ") ImportCSV(ctx context.Context, jobID, path string) {") {
		return src, nil, nil
	}
	out := src
	var fixed, warn []string
	imports := []string{}

	if matches := oldImportNameLookupRe.FindAllStringSubmatchIndex(out, -1); len(matches) > 0 {
		if strings.Count(out, importCountersAnchor) != 1 {
			warn = append(warn, "a belongs_to column here still costs a query per row and ignores lookup errors, and the importer is not the one grit generate wrote: resolve each value once and fail the row on an error")
		} else {
			resolvers := make([]string, 0, len(matches))
			for i := len(matches) - 1; i >= 0; i-- {
				loc := matches[i]
				group := func(n int) string { return out[loc[2*n]:loc[2*n+1]] }
				base, rel, key, rel2, keyGo, fk, amp := group(1), group(2), group(3), group(4), group(5), group(6), group(7)
				if rel != rel2 {
					continue
				}
				out = out[:loc[0]] + ImportNameAssign(base, fk, amp == "&") + out[loc[1]:]
				resolvers = append([]string{ImportNameResolver(base, rel, key, keyGo)}, resolvers...)
			}
			if len(resolvers) > 0 {
				out = strings.Replace(out, importCountersAnchor, strings.Join(resolvers, "")+importCountersAnchor, 1)
				imports = append(imports, "errors", "gorm.io/gorm")
				fixed = append(fixed, "each value in a belongs_to column is looked up once per import, and a lookup error fails the row")
			}
		}
	}

	if !strings.Contains(out, "imports.Wait(") {
		switch {
		case strings.Count(out, importOpenAnchor) != 1:
			warn = append(warn, "imports here are not limited, and the importer is not the one grit generate wrote: call imports.Wait before it opens the file")
		case strings.Contains(out, "record := func("):
			out = strings.Replace(out, importOpenAnchor, ImportLimitGate+importOpenAnchor, 1)
			imports = append(imports, module+"/internal/imports")
			fixed = append(fixed, "imports take turns, IMPORT_CONCURRENCY at a time (default 2)")
		case strings.Contains(out, "\tdb := s.db(ctx)\n"):
			// An importer generated before v3.251.0 has no record helper, so
			// the waiting message is written to the job directly.
			out = strings.Replace(out, importOpenAnchor, importLimitGateDirect+importOpenAnchor, 1)
			imports = append(imports, "log", module+"/internal/imports")
			fixed = append(fixed, "imports take turns, IMPORT_CONCURRENCY at a time (default 2)")
		default:
			warn = append(warn, "imports here are not limited, and the importer is not the one grit generate wrote: call imports.Wait before it opens the file")
		}
	}

	for _, path := range imports {
		var ok bool
		if out, ok = addImportLikeGenerate(out, path, module); !ok {
			return src, nil, append(warn, "could not add the "+path+" import to the importer, so it was left as it was")
		}
	}
	if len(fixed) == 0 {
		return src, nil, warn
	}
	return out, fixed, warn
}
