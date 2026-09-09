package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
			for _, f := range m.Files {
				if !f.When.matches(ctx) {
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
			return out
		},

		Injections: func(ctx Context) []Injection {
			var out []Injection
			for _, in := range m.Injections {
				if !in.When.matches(ctx) {
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

	Files      []localFile      `json:"files"`
	Injections []localInjection `json:"injections"`
}

type localFile struct {
	From string     `json:"from"` // relative to the plugin directory
	To   string     `json:"to"`   // relative to the project root
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
	return strings.NewReplacer(
		"{{MODULE}}", ctx.Module,
		"{{PROJECT}}", project,
		"{{API_ROOT}}", apiRootFor(ctx),
	).Replace(s)
}

// apiRootFor is where Go code lives, project-relative, so a plugin can write
// one path that works for --single and for a monorepo.
func apiRootFor(ctx Context) string {
	if ctx.Architecture == "single" {
		return "."
	}
	return "apps/api"
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
	if len(m.Files) == 0 && len(m.Injections) == 0 {
		return fmt.Errorf("plugin.json: a plugin with no files and no injections does nothing")
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
		if _, err := os.Stat(filepath.Join(dir, f.From)); err != nil {
			return fmt.Errorf("plugin.json: files[%d].from %q is not in the plugin directory", i, f.From)
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
