package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// wafRichtextEntry is one entry pair already in wafExcludedRoutes: a resource at
// its public path, which is how every version before v3.335.0 listed one.
var wafRichtextEntry = regexp.MustCompile(`(?m)^\t\t"/([a-z][a-z0-9_-]*)", "/([a-z][a-z0-9_-]*)/\*",$`)

// wafRouteRe finds the path a route is registered at, in either shape the
// generator has used: inline in routes.go, or in a per-resource routes file.
var wafRouteRe = regexp.MustCompile(`\.(?:GET|POST|PUT|PATCH|DELETE)\("/(?:admin/)?([a-z][a-z0-9_-]*)`)

// repairWAFRichtext closes grit#90 on a project already generated.
//
// Sentinel's WAF inspects a request body on the way in, so it is the writes that
// carry the markup, and the admin panel's writes go through the /admin prefix.
// wafExcludedRoutes listed only the public path of each richtext resource, so a
// post whose body quoted a shell command could be published once and then never
// edited: every PUT from the admin was answered 403, with the WAF logging
// PathTraversal for the "../" inside a code block.
//
// Two things are missing on a project generated before v3.335.0 and both are
// repaired here. The /admin twin of everything already listed, and the resources
// the list never mentioned at all: it was hand-written for the scaffolded blog
// and grit generate did not touch it, so every resource generated with a
// richtext field has the same fault at its own path.
func repairWAFRichtext(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	plurals := richtextPlurals(apiRoot)
	return repairSourceFile(root, m, routes, func(src string) (string, []string, []string) {
		return repairWAFRichtextSource(src, plurals)
	})
}

// repairWAFRichtextSource edits the paths literal in wafExcludedRoutes. Nothing
// is reordered and nothing is removed: a list somebody has curated keeps
// everything they put in it.
func repairWAFRichtextSource(src string, plurals []string) (string, []string, []string) {
	const decl = "func wafExcludedRoutes() []string {"
	start := strings.Index(src, decl)
	if start < 0 {
		return src, nil, nil
	}
	const literal = "paths := []string{\n"
	open := strings.Index(src[start:], literal)
	if open < 0 {
		return src, nil, []string{`its wafExcludedRoutes has no paths literal, so the admin panel is still refused 403 on a richtext body: add "/admin/<resource>" for every richtext resource`}
	}
	open += start + len(literal)
	end := strings.Index(src[open:], "\n\t}\n")
	if end < 0 {
		return src, nil, []string{`its wafExcludedRoutes paths literal does not end where Grit wrote it, so it was left alone: add "/admin/<resource>" for every richtext resource`}
	}
	end += open + 1

	block := src[open:end]
	var add []string

	// The /admin twin of everything already listed. Uploads is here for the
	// body cap rather than for markup and has no admin route, and the public
	// form share is a subtree entry that matches neither shape.
	for _, pair := range wafRichtextEntry.FindAllStringSubmatch(block, -1) {
		plural := pair[1]
		if plural != pair[2] || plural == "uploads" {
			continue
		}
		if strings.Contains(block, `"/admin/`+plural+`/*"`) {
			continue
		}
		add = append(add, "\t\t\"/admin/"+plural+"\", \"/admin/"+plural+"/*\",")
	}

	// The richtext resources the list never mentioned.
	var listed []string
	for _, plural := range plurals {
		if strings.Contains(block, `"/`+plural+`/*"`) {
			continue
		}
		listed = append(listed, plural)
		add = append(add, WAFRichtextExclusion(plural))
	}

	// Additions go above the marker, where grit generate puts them, so a list
	// repaired here and a list grown by the generator read the same way.
	at := end
	if i := strings.Index(block, "\t\t"+WAFRichtextMarker); i >= 0 {
		at = open + i
	} else {
		add = append(add, wafRichtextMarkerBlock)
	}
	if len(add) == 0 {
		return src, nil, nil
	}

	out := src[:at] + strings.Join(add, "\n") + "\n" + src[at:]
	what := "the admin panel's writes are no longer refused 403 on a richtext body"
	if len(listed) > 0 {
		what += ", " + strings.Join(listed, " and ") + " included"
	}
	return out, []string{what}, nil
}

// richtextPlurals is the route plural of every model with a richtext field.
//
// A richtext field is the one thing in a generated model that carries
// sanitize:"html", which makes the models the record of which resources have
// one. The plural is read back off the routes rather than derived, because
// pluralization lives in the generator: guessing it wrong here would add an
// exclusion for a path that does not exist and leave the real one inspected,
// which is the failure this repair exists to fix, reported as fixed.
func richtextPlurals(apiRoot string) []string {
	models, err := goSources(filepath.Join(apiRoot, "internal", "models"))
	if err != nil {
		return nil
	}
	var pascals []string
	for _, path := range models {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		pascals = append(pascals, structsWithHTMLSanitize(string(data))...)
	}
	if len(pascals) == 0 {
		return nil
	}

	routeFiles, err := goSources(filepath.Join(apiRoot, "internal", "routes"))
	if err != nil {
		return nil
	}
	sources := make(map[string]string, len(routeFiles))
	for _, path := range routeFiles {
		if data, err := os.ReadFile(path); err == nil {
			sources[filepath.Base(path)] = string(data)
		}
	}

	found := map[string]bool{}
	for _, pascal := range pascals {
		for _, plural := range routePluralsFor(sources, pascal) {
			found[plural] = true
		}
	}
	out := make([]string, 0, len(found))
	for plural := range found {
		out = append(out, plural)
	}
	sort.Strings(out)
	return out
}

// structsWithHTMLSanitize names every struct in one file that has a field
// tagged sanitize:"html".
func structsWithHTMLSanitize(src string) []string {
	var out []string
	current := ""
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if name, ok := structDeclName(trimmed); ok {
			current = name
			continue
		}
		if trimmed == "}" {
			current = ""
			continue
		}
		if current != "" && strings.Contains(trimmed, `sanitize:"html"`) {
			out = append(out, current)
			current = ""
		}
	}
	return out
}

// structDeclName reads the name off a one-line struct declaration.
func structDeclName(line string) (string, bool) {
	if !strings.HasPrefix(line, "type ") || !strings.HasSuffix(line, "struct {") {
		return "", false
	}
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return "", false
	}
	return fields[1], true
}

// routePluralsFor finds the paths a model's routes are registered at. The
// generator writes them one of two ways: inline in routes.go against a
// <camel>Handler variable, or in a per-resource file named <snake>_routes.go
// where every route in the file belongs to the one resource.
func routePluralsFor(sources map[string]string, pascal string) []string {
	handler := lowerFirst(pascal) + "Handler."
	own := toSnake(pascal) + "_routes.go"
	found := map[string]bool{}
	for name, src := range sources {
		mine := name == own
		for _, line := range strings.Split(src, "\n") {
			if !mine && !strings.Contains(line, handler) {
				continue
			}
			if match := wafRouteRe.FindStringSubmatch(line); match != nil {
				found[match[1]] = true
			}
		}
	}
	out := make([]string, 0, len(found))
	for plural := range found {
		out = append(out, plural)
	}
	sort.Strings(out)
	return out
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// toSnake converts a Go struct name to the stem the generator names files with.
// No pluralization: BlogPost is blog_post, whose routes file is
// blog_post_routes.go.
func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
