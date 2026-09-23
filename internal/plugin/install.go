package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// Install applies a plugin to the project and records what it did.
//
// Ordering matters: files are written before injections, so an injected import
// always refers to a file that exists. If any step fails the install stops and
// reports — it does not attempt a partial rollback, because a half-reverted
// project is harder to reason about than one you can inspect and re-run
// `grit plugin remove` on. The lockfile is written LAST, but records everything
// applied up to that point.
func Install(ctx Context, p Plugin) (*InstalledPlugin, error) {
	lock, err := LoadLock(ctx.Root)
	if err != nil {
		return nil, err
	}

	if _, already := lock.Find(p.Name); already {
		return nil, fmt.Errorf("%s is already installed: run `grit plugin remove %s` first to reinstall", p.Name, p.Name)
	}

	// Dependencies between plugins are checked up front so a half-applied
	// install can't leave a project referencing something that isn't there.
	for _, req := range p.Requires {
		if _, ok := lock.Find(req); !ok {
			return nil, fmt.Errorf("%s requires the %q plugin: install it first: grit plugin add %s", p.Name, req, req)
		}
	}

	record := InstalledPlugin{
		Name:        p.Name,
		Version:     p.Version,
		InstalledAt: time.Now().UTC(),
		Requires:    p.Requires,
		// GoDeps is appended as each one is actually fetched, so the lockfile
		// records what landed rather than what was merely declared.
		NodeDeps: p.NodeDeps,
	}

	// --- Go dependencies ---
	//
	// Run BEFORE any file is written: `go get` loads the module graph, and a
	// tree that already contains generated code importing a module not yet in
	// go.mod is exactly the state that makes it fail. Doing it first also means
	// a network error aborts the install before anything on disk changed.
	if len(p.GoDeps) > 0 {
		apiDir := ctx.APIRoot
		if apiDir == "" {
			apiDir = ctx.Root
		}
		for _, d := range p.GoDeps {
			spec := d.Name
			if d.Version != "" {
				spec += "@" + d.Version
			}
			cmd := exec.Command("go", "get", spec)
			cmd.Dir = apiDir
			cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
			if out, err := cmd.CombinedOutput(); err != nil {
				return nil, fmt.Errorf("go get %s: %w\n%s", spec, err, strings.TrimSpace(string(out)))
			}
			record.GoDeps = append(record.GoDeps, d)
			fmt.Printf("  ✓ go get %s\n", spec)
		}
	}

	// --- Files ---
	if p.Files != nil {
		files := p.Files(ctx)
		paths := make([]string, 0, len(files))
		for rel := range files {
			paths = append(paths, rel)
		}
		sortStrings(paths)

		for _, rel := range paths {
			abs := filepath.Join(ctx.Root, rel)
			// Never clobber. A plugin overwriting a file the user has edited is
			// a data-loss bug, and removal would then delete something the
			// plugin didn't create.
			if _, err := os.Stat(abs); err == nil {
				return nil, fmt.Errorf("refusing to overwrite existing file %s: move it aside and retry", rel)
			}
			if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
				return nil, fmt.Errorf("creating dir for %s: %w", rel, err)
			}
			written := codefmt.File(abs, files[rel])
			if err := os.WriteFile(abs, []byte(written), 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", rel, err)
			}
			record.Files = append(record.Files, rel)
			// What this file looked like when the plugin wrote it, so a later
			// update knows whether anyone has touched it since.
			if record.FileHashes == nil {
				record.FileHashes = map[string]string{}
			}
			record.FileHashes[rel] = hashOf(written)
			fmt.Printf("  ✓ %s\n", rel)
		}
	}

	// --- Injections ---
	if p.Injections != nil {
		for _, inj := range p.Injections(ctx) {
			abs := filepath.Join(ctx.Root, inj.File)
			if !fileExists(abs) {
				if inj.Optional {
					continue
				}
				return nil, fmt.Errorf("injection target %s does not exist", inj.File)
			}
			if err := injectBefore(abs, inj.Marker, inj.Code); err != nil {
				if inj.Optional {
					fmt.Printf("  ⚠ skipped %s (%v)\n", inj.File, err)
					continue
				}
				return nil, fmt.Errorf("injecting into %s: %w", inj.File, err)
			}
			// Record the EXACT text, so removal deletes precisely this.
			record.Injections = append(record.Injections, LockedInjection{
				File:   inj.File,
				Marker: inj.Marker,
				Code:   inj.Code,
			})
			fmt.Printf("  ✓ patched %s\n", inj.File)
		}
	}

	// --- npm dependencies ---
	//
	// Declared per workspace ("apps/expo"). A workspace the project does not
	// have is skipped: a plugin written for every client installs into the ones
	// that exist. These were recorded in the lockfile and never written, which
	// nobody noticed until the first plugin (push) needed one.
	for _, d := range p.NodeDeps {
		if d.Workspace == "" {
			continue
		}
		pkg := filepath.Join(ctx.Root, d.Workspace, "package.json")
		if !fileExists(pkg) {
			continue
		}
		added, err := addNodeDependency(pkg, d.Name, d.Version)
		if err != nil {
			return nil, fmt.Errorf("adding %s to %s: %w", d.Name, d.Workspace, err)
		}
		if added {
			fmt.Printf("  ✓ %s depends on %s %s (run pnpm install)\n", d.Workspace, d.Name, d.Version)
		}
	}

	lock.Add(record)
	if err := lock.Save(ctx.Root); err != nil {
		return nil, err
	}
	return &record, nil
}

// addNodeDependency adds "name": "version" to a package.json's dependencies,
// as text, so the file keeps its formatting and key order. A package already
// listed there, or in devDependencies, is left as it is.
func addNodeDependency(path, name, version string) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	src := string(raw)
	if strings.Contains(src, "\""+name+"\":") {
		return false, nil
	}
	const anchor = "\"dependencies\": {"
	i := strings.Index(src, anchor)
	if i < 0 {
		return false, fmt.Errorf("no \"dependencies\" object")
	}
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	after := i + len(anchor)
	// Indent like the first existing entry, or four spaces for an empty object.
	indent := "    "
	rest := strings.TrimLeft(src[after:], "\r\n")
	if j := strings.IndexFunc(rest, func(r rune) bool { return r != ' ' && r != '\t' }); j > 0 {
		indent = rest[:j]
	}
	entry := nl + indent + "\"" + name + "\": \"" + version + "\""
	if !strings.HasPrefix(strings.TrimSpace(src[after:]), "}") {
		entry += ","
	}
	return true, os.WriteFile(path, []byte(src[:after]+entry+src[after:]), 0o644)
}

// Remove undoes an install using the lockfile record.
//
// Everything it does is derived from what Install recorded — there is no
// separate removal list to fall out of sync. Anything it can't cleanly reverse
// is reported rather than forced, so the user is told what to check instead of
// being handed a project that doesn't compile.
func Remove(root, name string) ([]string, error) {
	lock, err := LoadLock(root)
	if err != nil {
		return nil, err
	}

	record, ok := lock.Find(name)
	if !ok {
		return nil, fmt.Errorf("%s is not installed", name)
	}

	// Refuse if another installed plugin depends on this one.
	// Checked against the LOCKFILE, not the registry: a plugin installed by a
	// newer CLI, or fetched from outside the built-in set, still has its
	// requirements honoured. Consulting the registry here meant an unknown
	// plugin's dependencies were silently ignored.
	for _, other := range lock.Plugins {
		if other.Name == name {
			continue
		}
		for _, req := range other.Requires {
			if req == name {
				return nil, fmt.Errorf("%s is required by %s: remove that first", name, other.Name)
			}
		}
	}

	var warnings []string

	// Injections first: they reference the files, so undo them before the files
	// disappear.
	for _, inj := range record.Injections {
		abs := filepath.Join(root, inj.File)
		if !fileExists(abs) {
			continue // file already gone; nothing to undo
		}
		data, err := os.ReadFile(abs)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("could not read %s: %v", inj.File, err))
			continue
		}
		content := string(data)
		if !strings.Contains(content, inj.Code) {
			// The user edited the injected block. Removing a guess could delete
			// their work, so say so and move on.
			warnings = append(warnings, fmt.Sprintf(
				"%s: the injected block was modified — remove it by hand (anchor: %s)", inj.File, inj.Marker))
			continue
		}
		// injectBefore inserts the snippet as its own line, so remove the
		// trailing newline with it — otherwise every install/remove cycle
		// leaves a blank line behind and the file never returns to its
		// original state.
		if strings.Contains(content, inj.Code+"\n") {
			content = strings.Replace(content, inj.Code+"\n", "", 1)
		} else {
			content = strings.Replace(content, inj.Code, "", 1)
		}
		if err := os.WriteFile(abs, []byte(content), 0644); err != nil {
			warnings = append(warnings, fmt.Sprintf("could not write %s: %v", inj.File, err))
			continue
		}
		fmt.Printf("  ✗ reverted %s\n", inj.File)
	}

	// Files, then any directories the plugin left empty.
	var dirs []string
	for _, rel := range record.Files {
		abs := filepath.Join(root, rel)
		if !fileExists(abs) {
			continue
		}
		if err := os.Remove(abs); err != nil {
			warnings = append(warnings, fmt.Sprintf("could not delete %s: %v", rel, err))
			continue
		}
		dirs = append(dirs, filepath.Dir(abs))
		fmt.Printf("  ✗ %s\n", rel)
	}
	pruneEmptyDirs(root, dirs)

	lock.Remove(name)
	if err := lock.Save(root); err != nil {
		return warnings, err
	}
	return warnings, nil
}

// injectBefore inserts code on the line before a marker.
//
// Self-contained rather than shared with internal/generate: a plugin applying
// an edit and a resource generator applying one have different failure
// expectations, and coupling them means a change for one silently alters the
// other.
func injectBefore(path, marker, code string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)

	lines := strings.Split(content, "\n")
	codeLines := strings.Split(code, "\n")

	for i, line := range lines {
		if strings.TrimSpace(line) != strings.TrimSpace(marker) {
			continue
		}
		// Already applied? Only the lines directly above THIS marker can
		// answer that. Searching the whole file asks a different question and
		// gets it wrong twice over: a plugin legitimately injecting the same
		// snippet at two markers loses the second, and a snippet that occurs
		// naturally somewhere else blocks the injection entirely. Both fail
		// silently, which is the worst way for codegen to be wrong.
		if i >= len(codeLines) && slicesEqual(lines[i-len(codeLines):i], codeLines) {
			return nil
		}
		out := append([]string{}, lines[:i]...)
		out = append(out, code)
		out = append(out, lines[i:]...)
		// Formatted, because two plugins injecting at the same marker leave
		// the second one's import out of order: a project with both video and
		// stripe installed stopped being gofmt-clean. No plugin has to know
		// what another injected.
		joined := strings.Join(out, "\n")
		return os.WriteFile(path, []byte(codefmt.File(path, joined)), 0644)
	}
	return fmt.Errorf("marker %q not found", marker)
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// pruneEmptyDirs removes directories a plugin emptied, deepest first, stopping
// at the project root.
func pruneEmptyDirs(root string, dirs []string) {
	sortStrings(dirs)
	for i := len(dirs) - 1; i >= 0; i-- {
		d := dirs[i]
		for strings.HasPrefix(d, root) && d != root {
			entries, err := os.ReadDir(d)
			if err != nil || len(entries) > 0 {
				break
			}
			if os.Remove(d) != nil {
				break
			}
			d = filepath.Dir(d)
		}
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// hashOf is the fingerprint stored in the lockfile for a file the plugin wrote.
func hashOf(content string) string {
	// Normalised, because a line ending is not a change: git rewrites them on
	// checkout on Windows, and a fingerprint that moves when it does would
	// report every file in the project as edited.
	sum := sha256.Sum256([]byte(normalizeEOL(content)))
	return hex.EncodeToString(sum[:])
}
