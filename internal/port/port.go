// Package port copies a component out of the admin panel into another app.
//
// The admin is where Grit's interface work has gone: a data table with
// sorting, filtering, bulk actions, export and import, a form builder, the
// dashboard widgets. A customer-facing portal usually wants some of that, and
// until now the only way to have it was to build it again. `grit expose table`
// writes a page with its own eighty-line table rather than the admin's, and
// anyone who wanted the real one rebuilt it by hand (grit#113).
//
// The copy is deliberate rather than a shared package.
//
// A package would be the tidier answer on a whiteboard, and it would mean
// decoupling every component from next/link, the icon map and the translation
// hook first, adding transpilePackages, and rewriting the imports of every
// project that already exists. The components are also not finished: they are
// edited per project, which is the whole premise of the admin being yours
// rather than a dependency. So this vendors them, the way the admin's own
// components arrived in the first place, and the copy belongs to the app it
// lands in.
//
// What makes that cheap here: both apps use the same "@/*" alias and the same
// design tokens, so a file copied to the same relative path compiles and is
// painted correctly with no rewriting at all. The work is finding the whole
// set of files, not translating them.
package port

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Bundle is a named starting point, so nobody has to know the file layout.
type Bundle struct {
	Name    string
	Summary string
	Entries []string // relative to the admin app root, forward slashes
}

// Bundles are the sets worth asking for by name. A path can be given instead.
var Bundles = []Bundle{
	{
		Name:    "table",
		Summary: "The data table: sorting, filters, pagination, selection, empty and loading states",
		Entries: []string{"components/tables/data-table.tsx"},
	},
	{
		Name:    "table-full",
		Summary: "The table with its toolbar, bulk actions, export and import",
		Entries: []string{
			"components/tables/data-table.tsx",
			"components/tables/table-toolbar.tsx",
			"components/tables/bulk-action-bar.tsx",
			"components/tables/export-menu.tsx",
			"components/tables/import-modal.tsx",
		},
	},
	{
		Name:    "form",
		Summary: "The form builder and its field types",
		Entries: []string{"components/forms/form-builder.tsx"},
	},
	{
		Name:    "widgets",
		Summary: "The dashboard widgets: stat cards, charts, activity",
		Entries: []string{"components/dashboard"},
	},
	{
		Name:    "ui",
		Summary: "The primitives: button, input, modal, dropzone, skeleton",
		Entries: []string{"components/ui"},
	},
}

// Result is what a run did, for the caller to print.
type Result struct {
	Copied     []string // paths relative to the target app
	Skipped    []string // already there, and not --force
	Unresolved []string // an import that could not be followed, with its file
	External   []string // every package these files import
	Missing    []string // the ones not in the target's package.json
}

// Options is one run.
type Options struct {
	AdminRoot  string // the admin app's directory
	TargetRoot string // the app to copy into
	Entries    []string
	Force      bool
	DryRun     bool
}

// importRe finds every module specifier in a TypeScript file: the `from "x"`
// of an import or a re-export, a bare `import "x"`, and a dynamic `import("x")`.
//
// A regex rather than a parser because the alternative is a TypeScript
// toolchain in a Go binary, and the failure mode here is visible: an import
// this misses is reported as unresolved rather than silently dropped.
var importRe = regexp.MustCompile(`(?m)(?:from\s*|import\s*\(?\s*)["']([^"']+)["']`)

// Run copies the entries and everything they reach.
func Run(opts Options) (*Result, error) {
	if opts.AdminRoot == "" || opts.TargetRoot == "" {
		return nil, fmt.Errorf("both the admin and the target app are needed")
	}
	if opts.AdminRoot == opts.TargetRoot {
		return nil, fmt.Errorf("the admin and the target are the same directory: " +
			"on this architecture the admin is already inside that app, so there is " +
			"nothing to copy")
	}

	res := &Result{}
	seen := map[string]bool{}
	external := map[string]bool{}

	var queue []string
	for _, e := range opts.Entries {
		expanded, err := expand(opts.AdminRoot, e)
		if err != nil {
			return nil, err
		}
		queue = append(queue, expanded...)
	}
	if len(queue) == 0 {
		return nil, fmt.Errorf("nothing to copy: no file matched the entries given")
	}

	for len(queue) > 0 {
		rel := queue[0]
		queue = queue[1:]
		if seen[rel] {
			continue
		}
		seen[rel] = true

		src := filepath.Join(opts.AdminRoot, filepath.FromSlash(rel))
		data, err := os.ReadFile(src)
		if err != nil {
			res.Unresolved = append(res.Unresolved, fmt.Sprintf("%s (could not read it)", rel))
			continue
		}

		for _, spec := range specifiers(string(data)) {
			switch {
			case strings.HasPrefix(spec, "@/"):
				next, ok := resolve(opts.AdminRoot, strings.TrimPrefix(spec, "@/"))
				if !ok {
					// Not a file in the admin: an alias the target app also
					// has, or something generated. Named rather than guessed.
					res.Unresolved = append(res.Unresolved, fmt.Sprintf("%s imports %s", rel, spec))
					continue
				}
				queue = append(queue, next)
			case strings.HasPrefix(spec, "./"), strings.HasPrefix(spec, "../"):
				joined := filepath.ToSlash(filepath.Join(filepath.Dir(rel), spec))
				next, ok := resolve(opts.AdminRoot, joined)
				if !ok {
					res.Unresolved = append(res.Unresolved, fmt.Sprintf("%s imports %s", rel, spec))
					continue
				}
				queue = append(queue, next)
			default:
				// A package. The target needs it in its package.json, which is
				// the caller's to check: writing a dependency into somebody's
				// manifest is not something a copy should do quietly.
				external[packageOf(spec)] = true
			}
		}

		dest := filepath.Join(opts.TargetRoot, filepath.FromSlash(rel))
		if _, err := os.Stat(dest); err == nil && !opts.Force {
			res.Skipped = append(res.Skipped, rel)
			continue
		}
		if !opts.DryRun {
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return nil, fmt.Errorf("creating %s: %w", filepath.Dir(dest), err)
			}
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", dest, err)
			}
		}
		res.Copied = append(res.Copied, rel)
	}

	declared := declaredPackages(opts.TargetRoot)
	for p := range external {
		res.External = append(res.External, p)
		// Only the ones the target does not declare. Listing every package,
		// react and next included, buries the one line that matters under
		// things everybody already has.
		//
		// A workspace hoists its dependencies, so a missing one still type
		// checks here and fails when that app is built on its own. Which is
		// exactly why this is worth saying out loud rather than leaving to a
		// compiler that will not complain today.
		if !declared[p] {
			res.Missing = append(res.Missing, p)
		}
	}
	sort.Strings(res.Missing)
	sort.Strings(res.Copied)
	sort.Strings(res.Skipped)
	sort.Strings(res.External)
	sort.Strings(res.Unresolved)
	return res, nil
}

// expand turns an entry into the files it names: one file, or every source
// file directly inside a directory.
func expand(root, entry string) ([]string, error) {
	clean := filepath.ToSlash(strings.TrimPrefix(entry, "./"))
	full := filepath.Join(root, filepath.FromSlash(clean))

	info, err := os.Stat(full)
	if err == nil && info.IsDir() {
		entries, err := os.ReadDir(full)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", clean, err)
		}
		var out []string
		for _, e := range entries {
			if e.IsDir() || !isSource(e.Name()) {
				continue
			}
			out = append(out, clean+"/"+e.Name())
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("%s holds no .ts or .tsx files", clean)
		}
		return out, nil
	}

	if resolved, ok := resolve(root, clean); ok {
		return []string{resolved}, nil
	}
	return nil, fmt.Errorf("%s is not in the admin panel at %s", clean, root)
}

// resolve finds the real file for a path with no extension, the way a
// bundler does.
func resolve(root, rel string) (string, bool) {
	rel = filepath.ToSlash(rel)
	candidates := []string{rel}
	if filepath.Ext(rel) == "" {
		candidates = []string{
			rel + ".tsx", rel + ".ts",
			rel + "/index.tsx", rel + "/index.ts",
		}
	}
	for _, c := range candidates {
		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(c))); err == nil && !info.IsDir() {
			return c, true
		}
	}
	return "", false
}

// declaredPackages reads the target app's package.json.
//
// Read as text rather than decoded, because all that is wanted is whether a
// name appears as a key in one of the dependency blocks, and a manifest with
// a field this does not know about should not stop a copy.
func declaredPackages(appRoot string) map[string]bool {
	out := map[string]bool{}
	data, err := os.ReadFile(filepath.Join(appRoot, "package.json"))
	if err != nil {
		return out
	}
	for _, m := range depRe.FindAllStringSubmatch(string(data), -1) {
		out[m[1]] = true
	}
	return out
}

// depRe matches a quoted key followed by a quoted version, which is what a
// dependency entry looks like and what a "name" or "main" field does not.
var depRe = regexp.MustCompile(`"([^"]+)"\s*:\s*"(?:\^|~|>=|<|\*|workspace:|catalog:|npm:|file:|link:|\d)[^"]*"`)

func specifiers(src string) []string {
	var out []string
	for _, m := range importRe.FindAllStringSubmatch(src, -1) {
		out = append(out, m[1])
	}
	return out
}

// packageOf is the installable name of a specifier: "lucide-react/icons/x" is
// the lucide-react package, "@tanstack/react-query/foo" is @tanstack/react-query.
func packageOf(spec string) string {
	parts := strings.Split(spec, "/")
	if strings.HasPrefix(spec, "@") && len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return parts[0]
}

func isSource(name string) bool {
	return (strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".tsx")) &&
		!strings.HasSuffix(name, ".d.ts") &&
		!strings.HasSuffix(name, ".test.ts") &&
		!strings.HasSuffix(name, ".test.tsx")
}
