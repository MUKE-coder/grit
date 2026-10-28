package plugin

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// runOfCapitals finds three or more consecutive uppercase letters in a field
// name.
var runOfCapitals = regexp.MustCompile(`[A-Z]{3,}`)

// gormKnownInitialisms is the list GORM's naming strategy recognises. A run
// of capitals that is one of these is handled correctly and needs no tag:
// PosterURL becomes poster_url, IPAddress becomes ip_address.
//
// The danger is a run that is NOT on this list but CONTAINS one that is,
// because GORM finds the inner one and splits around it. That is exactly how
// WIPLimit became w_ip_limit: WIP is not here, IP is.
var gormKnownInitialisms = map[string]bool{
	"API": true, "ASCII": true, "CPU": true, "CSS": true, "DNS": true,
	"EOF": true, "GUID": true, "HTML": true, "HTTP": true, "HTTPS": true,
	"ID": true, "IP": true, "JSON": true, "LHS": true, "QPS": true,
	"RAM": true, "RHS": true, "RPC": true, "SLA": true, "SMTP": true,
	"SSH": true, "TLS": true, "TTL": true, "UID": true, "UI": true,
	"UUID": true, "URI": true, "URL": true, "UTF8": true, "VM": true,
	"XML": true, "XSRF": true, "XSS": true,
}

// suspiciousInitialisms returns the capital runs in a field name that GORM
// will not resolve the way a reader expects.
func suspiciousInitialisms(name string) []string {
	var out []string
	for _, loc := range runOfCapitals.FindAllStringIndex(name, -1) {
		run := name[loc[0]:loc[1]]
		// A run followed by lowercase ends with the first letter of the next
		// word, not the initialism: in URLValue the run is URLV and the
		// initialism is URL.
		if loc[1] < len(name) && name[loc[1]] >= 'a' && name[loc[1]] <= 'z' {
			run = run[:len(run)-1]
		}
		if len(run) < 2 || gormKnownInitialisms[run] {
			continue
		}
		out = append(out, run)
	}
	return out
}

// TestEmittedModelsPinUnusualInitialisms guards a bug that costs an afternoon
// and leaves no trace until runtime.
//
// GORM derives a column name from the field name, and it does so by looking
// for known initialisms anywhere in the word. WIPLimit contains IP, which is
// on that list, so the column is w_ip_limit, not wip_limit. Nothing complains:
// AutoMigrate creates w_ip_limit happily, reading the struct works, and the
// only thing that breaks is a write by column name:
//
//	db.Model(&col).Updates(map[string]any{"wip_limit": 3})
//
// which fails in a handler, in production, against a column that does not
// exist. The kanban plugin shipped with exactly this bug until a generated
// test caught it.
//
// The rule is therefore: a field name with a run of three or more capitals
// must say what its column is called rather than let GORM decide.
func TestEmittedModelsPinUnusualInitialisms(t *testing.T) {
	ctx := Context{
		Root:         "/project",
		Module:       "example/apps/api",
		APIRoot:      "/project/apps/api",
		Architecture: "triple",
		Frontend:     "next",
	}

	checked := 0
	for _, p := range All() {
		if p.Files == nil {
			continue
		}
		for path, src := range p.Files(ctx) {
			if !strings.Contains(path, "internal/models/") || !strings.HasSuffix(path, ".go") {
				continue
			}
			checked++
			for _, f := range unpinnedInitialismFields(t, p.Name, path, src) {
				t.Errorf("%s emits %s with field %s: %q is not an initialism GORM knows, "+
					"so it picks the column name and its pick is not the obvious one.\n"+
					"Add an explicit column to the gorm tag, e.g. gorm:\"column:%s\".",
					p.Name, path, f.name, strings.Join(suspiciousInitialisms(f.name), ", "), f.jsonName)
			}
		}
	}

	if checked == 0 {
		t.Fatal("no emitted model files were checked; this test has stopped testing anything")
	}
	t.Logf("checked %d emitted model file(s)", checked)
}

type modelField struct {
	name     string
	jsonName string
}

func unpinnedInitialismFields(t *testing.T, plugin, path, src string) []modelField {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("%s emits %s, which does not parse as Go: %v", plugin, path, err)
	}

	var out []modelField
	ast.Inspect(file, func(n ast.Node) bool {
		st, ok := n.(*ast.StructType)
		if !ok || st.Fields == nil {
			return true
		}
		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue // embedded
			}
			name := field.Names[0].Name
			if len(suspiciousInitialisms(name)) == 0 {
				continue
			}

			gormTag, jsonTag := tagsOf(field)
			if strings.Contains(gormTag, "column:") {
				continue // pinned, which is the point
			}
			if jsonTag == "-" || strings.Contains(gormTag, "-") && gormTag == "-" {
				continue // not persisted
			}

			out = append(out, modelField{name: name, jsonName: strings.TrimSuffix(jsonTag, ",omitempty")})
		}
		return true
	})
	return out
}

func tagsOf(field *ast.Field) (gormTag, jsonTag string) {
	if field.Tag == nil {
		return "", ""
	}
	raw, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		raw = strings.Trim(field.Tag.Value, "`")
	}
	return structTag(raw, "gorm"), structTag(raw, "json")
}

// structTag reads one key out of a struct tag. reflect.StructTag would do
// this, but reaching for reflect to read a string that is already in hand is
// more machinery than the job needs.
func structTag(raw, key string) string {
	for _, part := range strings.Fields(raw) {
		prefix := key + ":\""
		if strings.HasPrefix(part, prefix) {
			return strings.TrimSuffix(strings.TrimPrefix(part, prefix), "\"")
		}
	}
	// Tags can contain spaces inside the quoted value, so fall back to a
	// scan when the simple split misses.
	at := strings.Index(raw, key+":\"")
	if at < 0 {
		return ""
	}
	rest := raw[at+len(key)+2:]
	end := strings.Index(rest, "\"")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// TestInitialismRuleCatchesTheRealCase proves the rule above is not vacuous:
// it fires on the exact field that shipped broken, and stays quiet once the
// column is pinned.
func TestInitialismRuleCatchesTheRealCase(t *testing.T) {
	broken := "package models\n\ntype C struct {\n\tWIPLimit int `gorm:\"default:0\" json:\"wip_limit\"`\n}\n"
	if got := unpinnedInitialismFields(t, "test", "internal/models/c.go", broken); len(got) != 1 {
		t.Fatalf("the rule did not flag an unpinned WIPLimit: %v", got)
	}

	fixed := "package models\n\ntype C struct {\n\tWIPLimit int `gorm:\"column:wip_limit;default:0\" json:\"wip_limit\"`\n}\n"
	if got := unpinnedInitialismFields(t, "test", "internal/models/c.go", fixed); len(got) != 0 {
		t.Fatalf("the rule flagged a pinned WIPLimit: %v", got)
	}

	ordinary := "package models\n\ntype C struct {\n\tBoardID string `gorm:\"size:36\" json:\"board_id\"`\n}\n"
	if got := unpinnedInitialismFields(t, "test", "internal/models/c.go", ordinary); len(got) != 0 {
		t.Fatalf("the rule flagged BoardID, whose trailing ID GORM handles: %v", got)
	}

	// The first version of this rule flagged these, which was noise: GORM
	// resolves all three correctly because URL and IP are on its list.
	known := "package models\n\ntype C struct {\n\tURL string `json:\"url\"`\n\tPosterURL string `json:\"poster_url\"`\n\tIPAddress string `json:\"ip_address\"`\n}\n"
	if got := unpinnedInitialismFields(t, "test", "internal/models/c.go", known); len(got) != 0 {
		t.Fatalf("the rule flagged initialisms GORM already knows: %v", got)
	}
}
