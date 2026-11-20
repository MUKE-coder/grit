package plugin

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// LoadDir builds a Plugin from a directory on disk.
//
// The authoring docs describe a plugin as a plugin.Plugin value and end by
// pointing at internal/plugin/multitenant.go, which is inside the CLI. So a
// reader who followed them had nowhere to put what they had written: the
// registry was compiled in, and writing a plugin meant forking Grit, adding a
// Go file, and rebuilding the binary. The page never said so.
//
// This is the same Plugin shape, described in JSON, so nothing about install
// or removal changes: the lockfile still records every file and injection, and
// removal still replays that record backwards. What changes is who can write
// one.
//
// A plugin directory:
//
//	my-plugin/
//	  plugin.json          the manifest
//	  files/               file bodies, referenced by "from"
//
// Paths in the manifest are project-relative for "to" and plugin-relative for
// "from". File bodies get {{MODULE}} and {{PROJECT}} substituted, so a plugin
// does not need to know the module path of the project installing it.
func LoadDir(dir string) (Plugin, error) {
	manifestPath := filepath.Join(dir, "plugin.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Plugin{}, fmt.Errorf("%s has no plugin.json\n\n"+
				"A plugin directory holds a plugin.json manifest and a files/ directory.\n"+
				"See https://gritframework.dev/docs/plugins/authoring", dir)
		}
		return Plugin{}, err
	}

	var m localManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return Plugin{}, fmt.Errorf("reading %s: %w", manifestPath, err)
	}
	if err := m.validate(dir); err != nil {
		return Plugin{}, err
	}

	return Plugin{
		Name:        m.Name,
		Version:     m.Version,
		Summary:     m.Summary,
		Description: m.Description,
		Requires:    m.Requires,
		GoDeps:      m.GoDeps,
		NodeDeps:    m.NodeDeps,
		NextSteps:   m.NextSteps,

		Files: func(ctx Context) map[string]string {
			out := map[string]string{}
			_, hasAdmin := adminRootFor(ctx)
			for _, f := range m.Files {
				if !f.When.matches(ctx) {
					continue
				}
				// An admin file on a project with no admin panel. Skipped
				// rather than written to a path with a hole in it, the same
				// way a node dependency for a workspace the project does not
				// have is skipped.
				if !hasAdmin && (needsAdminRoot(f.To) || needsAdminRoot(f.From)) {
					continue
				}
				body, err := os.ReadFile(filepath.Join(dir, f.From))
				if err != nil {
					// Validated at load time, so this is a file that vanished
					// between then and now. Skipping beats a panic mid-install.
					continue
				}
				out[expand(f.To, ctx)] = expand(string(body), ctx)
			}
			for _, d := range m.Directories {
				if !d.When.matches(ctx) {
					continue
				}
				if !hasAdmin && needsAdminRoot(d.To) {
					continue
				}
				for rel, body := range readTree(filepath.Join(dir, d.From)) {
					// The tree's shape is kept under "to", so a plugin's
					// files/components/day/grid.tsx lands at
					// <to>/day/grid.tsx. Slashes, not the host separator: a
					// destination is a project path, and a Windows backslash
					// in one would travel into the lockfile.
					to := path.Join(expand(d.To, ctx), rel)
					out[to] = expand(body, ctx)
				}
			}
			return out
		},

		Injections: func(ctx Context) []Injection {
			var out []Injection
			_, hasAdmin := adminRootFor(ctx)
			for _, in := range m.Injections {
				if !in.When.matches(ctx) {
					continue
				}
				if !hasAdmin && needsAdminRoot(in.File) {
					continue
				}
				code := in.Code
				if in.CodeFrom != "" {
					body, err := os.ReadFile(filepath.Join(dir, in.CodeFrom))
					if err != nil {
						continue
					}
					code = string(body)
				}
				out = append(out, Injection{
					File:     expand(in.File, ctx),
					Marker:   in.Marker,
					Code:     expand(strings.TrimRight(code, "\n"), ctx),
					Optional: in.Optional,
				})
			}
			return out
		},
	}, nil
}

// localManifest is plugin.json.
type localManifest struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Requires    []string     `json:"requires"`
	GoDeps      []Dependency `json:"goDeps"`
	NodeDeps    []Dependency `json:"nodeDeps"`
	NextSteps   []string     `json:"nextSteps"`

	Files       []localFile      `json:"files"`
	Directories []localDirectory `json:"directories"`
	Injections  []localInjection `json:"injections"`
}

type localFile struct {
	From string     `json:"from"` // relative to the plugin directory
	To   string     `json:"to"`   // relative to the project root
	When whenClause `json:"when"`
}

// localDirectory installs a whole tree, keeping its shape.
//
// Without it a plugin lists every file it writes one by one, repeating the
// same "to" prefix and the same "when" clause for each. grit-calendar needed
// nine entries for four groups of files, and adding a component meant editing
// the tree and the manifest and keeping the two in step (grit#227).
//
// Every file still reaches the lockfile as its own entry, so removal replays
// the install exactly and needs no second list.
type localDirectory struct {
	From string     `json:"from"` // a directory, relative to the plugin directory
	To   string     `json:"to"`   // a directory, relative to the project root
	When whenClause `json:"when"`
}

type localInjection struct {
	File     string     `json:"file"`
	Marker   string     `json:"marker"`
	Code     string     `json:"code"`
	CodeFrom string     `json:"codeFrom"` // read the snippet from a file instead
	Optional bool       `json:"optional"`
	When     whenClause `json:"when"`
}

// whenClause limits a file or injection to certain projects.
//
//	"when": {"architecture": ["triple", "full"], "frontend": ["next"]}
//
// Empty means always. Kept as data rather than an expression language: a
// plugin that needs real logic is better written in Go, and a half-implemented
// expression parser is a support burden nobody asked for.
type whenClause struct {
	Architecture []string `json:"architecture"`
	Frontend     []string `json:"frontend"`
}

func (w whenClause) matches(ctx Context) bool {
	return listAllows(w.Architecture, ctx.Architecture) &&
		listAllows(w.Frontend, ctx.Frontend)
}

func listAllows(allowed []string, actual string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, a := range allowed {
		if strings.EqualFold(a, actual) {
			return true
		}
	}
	return false
}

// expand substitutes the placeholders a plugin may use.
func expand(s string, ctx Context) string {
	project := filepath.Base(ctx.Root)
	adminRoot, _ := adminRootFor(ctx)
	return strings.NewReplacer(
		"{{MODULE}}", ctx.Module,
		"{{PROJECT}}", project,
		"{{API_ROOT}}", apiRootFor(ctx),
		"{{ADMIN_ROOT}}", adminRoot,
	).Replace(s)
}

// HasAdminPanel reports where a project's admin panel is, if it has one.
//
// Exported for the one caller outside this package: the install command, which
// has to tell somebody why a plugin wrote nothing.
func HasAdminPanel(ctx Context) (string, bool) { return adminRootFor(ctx) }

// needsAdminRoot reports whether a path cannot be resolved without an admin
// panel, which is how an entry gets skipped on a project that has none.
func needsAdminRoot(s string) bool { return strings.Contains(s, "{{ADMIN_ROOT}}") }

// apiRootFor is where Go code lives, project-relative, so a plugin can write
// one path that works for --single and for a monorepo.
func apiRootFor(ctx Context) string {
	if ctx.Architecture == "single" {
		return "."
	}
	return "apps/api"
}

// adminRootFor is where the admin panel's code lives, project-relative.
//
// {{API_ROOT}} exists so a Go path works on any architecture, and there was no
// equivalent for the admin, which is why grit-calendar's manifest needed 21
// file entries for 11 files: every one written twice, "apps/admin/..." and
// "apps/admin/src/...", each with its own when clause. And double, single and
// api projects could not be targeted at all (grit#228).
//
// Read off the disk rather than derived from the architecture, which is what
// apiRootFor does. The panel lives in five different places depending on shape
// and frontend, and exactly one of them exists in any project that has one, so
// looking is both shorter and right for a shape nobody enumerated here. The
// same approach as internal/generate's own resolver, which these two packages
// cannot share because neither imports the other.
//
// Returns false when the project has no admin panel at all: an --api or
// --mobile project, where a plugin's admin files should be skipped rather than
// written to a path with a hole in it.
func adminRootFor(ctx Context) (string, bool) {
	for _, candidate := range []string{
		// triple: its own application. The Vite shape first, because
		// apps/admin exists in both and only the Vite one has src/.
		"apps/admin/src",
		"apps/admin",
		// double: a section of the web app, Vite first for the same reason.
		"apps/web/src/admin-panel",
		"apps/web/admin-panel",
		// single: the frontend owns the project root.
		"src/admin-panel",
		"frontend/src/admin-panel",
	} {
		if dirExists(filepath.Join(ctx.Root, filepath.FromSlash(candidate))) {
			// The code root, not the application root. A Vite admin keeps its
			// components under src/, so a plugin writing
			// {{ADMIN_ROOT}}/components/x.tsx has to land in the right one of
			// the two, and the deeper candidate is checked first for exactly
			// that reason.
			return candidate, true
		}
	}
	return "", false
}

// validate catches the mistakes worth catching before anything is written.
//
// A plugin that fails halfway leaves files on disk and a lockfile that does
// not mention them, which is the state removal cannot recover from. Every
// check here is one that would otherwise surface mid-install.
func (m localManifest) validate(dir string) error {
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("plugin.json: name is required")
	}
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("plugin.json: version is required")
	}
	if strings.TrimSpace(m.Summary) == "" {
		return fmt.Errorf("plugin.json: summary is required (it is what `grit plugin list` shows)")
	}
	if len(m.Files) == 0 && len(m.Directories) == 0 && len(m.Injections) == 0 {
		return fmt.Errorf("plugin.json: a plugin with no files, no directories and no injections does nothing")
	}

	for i, f := range m.Files {
		if f.From == "" || f.To == "" {
			return fmt.Errorf("plugin.json: files[%d] needs both \"from\" and \"to\"", i)
		}
		if err := safeRelative(f.From); err != nil {
			return fmt.Errorf("plugin.json: files[%d].from %q: %w", i, f.From, err)
		}
		if err := safeRelative(f.To); err != nil {
			return fmt.Errorf("plugin.json: files[%d].to %q: %w", i, f.To, err)
		}
		info, err := os.Stat(filepath.Join(dir, f.From))
		if err != nil {
			return fmt.Errorf("plugin.json: files[%d].from %q is not in the plugin directory", i, f.From)
		}
		// A directory used to pass this stat and then fail the read at install
		// time, where the error was swallowed as "a file that vanished": the
		// entry installed nothing and the install reported success (grit#227).
		if info.IsDir() {
			return fmt.Errorf("plugin.json: files[%d].from %q is a directory. "+
				"Use a \"directories\" entry for a tree, or name a file here", i, f.From)
		}
	}

	for i, d := range m.Directories {
		if d.From == "" || d.To == "" {
			return fmt.Errorf("plugin.json: directories[%d] needs both \"from\" and \"to\"", i)
		}
		if err := safeRelative(d.From); err != nil {
			return fmt.Errorf("plugin.json: directories[%d].from %q: %w", i, d.From, err)
		}
		if err := safeRelative(d.To); err != nil {
			return fmt.Errorf("plugin.json: directories[%d].to %q: %w", i, d.To, err)
		}
		info, err := os.Stat(filepath.Join(dir, d.From))
		if err != nil {
			return fmt.Errorf("plugin.json: directories[%d].from %q is not in the plugin directory", i, d.From)
		}
		if !info.IsDir() {
			return fmt.Errorf("plugin.json: directories[%d].from %q is a file. "+
				"Use a \"files\" entry for one file", i, d.From)
		}
		// Checked here rather than at install, where there is nothing to
		// return an error to: the whole point of this pass is that a manifest
		// which installs nothing says so before it installs nothing.
		tree, err := treeFiles(filepath.Join(dir, d.From))
		if err != nil {
			return fmt.Errorf("plugin.json: directories[%d].from %q: %w", i, d.From, err)
		}
		if len(tree) == 0 {
			return fmt.Errorf("plugin.json: directories[%d].from %q holds no files, so it installs nothing", i, d.From)
		}
		// Every destination the tree produces, not just the prefix. A file
		// named ../x inside the tree would otherwise climb out of the project
		// after the prefix was joined.
		for _, rel := range tree {
			if err := safeRelative(path.Join(d.To, rel)); err != nil {
				return fmt.Errorf("plugin.json: directories[%d] would install %q: %w", i, rel, err)
			}
		}
	}

	for i, in := range m.Injections {
		if in.File == "" || in.Marker == "" {
			return fmt.Errorf("plugin.json: injections[%d] needs both \"file\" and \"marker\"", i)
		}
		if in.Code == "" && in.CodeFrom == "" {
			return fmt.Errorf("plugin.json: injections[%d] needs \"code\" or \"codeFrom\"", i)
		}
		if in.Code != "" && in.CodeFrom != "" {
			return fmt.Errorf("plugin.json: injections[%d] sets both \"code\" and \"codeFrom\"; pick one", i)
		}
		if err := safeRelative(in.File); err != nil {
			return fmt.Errorf("plugin.json: injections[%d].file %q: %w", i, in.File, err)
		}
		if in.CodeFrom != "" {
			if err := safeRelative(in.CodeFrom); err != nil {
				return fmt.Errorf("plugin.json: injections[%d].codeFrom %q: %w", i, in.CodeFrom, err)
			}
			if _, err := os.Stat(filepath.Join(dir, in.CodeFrom)); err != nil {
				return fmt.Errorf("plugin.json: injections[%d].codeFrom %q is not in the plugin directory", i, in.CodeFrom)
			}
		}
	}
	return nil
}

// treeFiles lists a directory's files, as slash-separated paths relative to it.
//
// Symlinks are refused rather than followed. A plugin is somebody else's code,
// and a link named files/components/x.tsx pointing at /etc/passwd would be
// read and installed into the project without comment, which is the same hole
// safeRelative exists to close from the other direction.
//
// Nothing is filtered, deliberately. A rule that skips dotfiles would be a
// rule somebody has to learn, and the install prints every file it writes, so
// a stray .DS_Store is visible in the output where a silent exclusion would
// not be.
func treeFiles(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			rel, _ := filepath.Rel(root, p)
			return fmt.Errorf("%s is a symlink; a plugin's files have to be real files",
				filepath.ToSlash(rel))
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// readTree is treeFiles with the contents, for the install.
//
// A file that cannot be read is skipped, the same way a missing single file is:
// validate has already walked this tree, so anything failing here vanished
// between load and install, and skipping beats a panic mid-install.
func readTree(root string) map[string]string {
	out := map[string]string{}
	names, err := treeFiles(root)
	if err != nil {
		return out
	}
	for _, rel := range names {
		body, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if readErr != nil {
			continue
		}
		out[rel] = string(body)
	}
	return out
}

// safeRelative refuses paths that climb out of their root.
//
// A plugin is somebody else's code, and "to": "../../.ssh/authorized_keys"
// installs it wherever the user happens to be. The installer would otherwise
// write it without comment.
func safeRelative(p string) error {
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") || strings.HasPrefix(p, `\`) {
		return fmt.Errorf("must be relative")
	}
	if len(p) > 1 && p[1] == ':' {
		return fmt.Errorf("must be relative") // C:\... on Windows
	}
	clean := filepath.ToSlash(filepath.Clean(p))
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return fmt.Errorf("must not climb outside the project")
	}
	return nil
}

// IsDirRef reports whether a `grit plugin add` argument names a directory
// rather than a built-in plugin.
//
// Explicit path syntax, so a future built-in cannot be shadowed by a directory
// that happens to share its name.
func IsDirRef(arg string) bool {
	// A leading slash is checked explicitly rather than left to
	// filepath.IsAbs, which is false for "/x" on Windows. A plugin name never
	// starts with a separator, so reading one as a path is safe on either.
	return strings.HasPrefix(arg, "./") ||
		strings.HasPrefix(arg, "../") ||
		strings.HasPrefix(arg, ".\\") ||
		strings.HasPrefix(arg, "..\\") ||
		strings.HasPrefix(arg, "/") ||
		strings.HasPrefix(arg, "\\") ||
		filepath.IsAbs(arg)
}
