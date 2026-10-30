package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Next 16 renamed the middleware file convention. A project with a
// middleware.ts prints this on every dev server start:
//
//	⚠ The "middleware" file convention is deprecated. Please use "proxy"
//	  instead.
//
// The file is otherwise identical: same request object, same matcher, and the
// function it exports is called proxy rather than middleware. So this is a
// rename, and it is applied to the files Grit wrote.
//
// A file the developer has edited is theirs. The deprecated name still works,
// so leaving it costs them a warning rather than a build, and silently
// rewriting a file somebody has put their own protected paths into is worse
// than that. They get the one line it takes to do it themselves.
func repairProxyConvention(root string, opts Options) error {
	for _, app := range nextAppRoots(root, opts) {
		if err := moveMiddlewareToProxy(root, app); err != nil {
			return err
		}
	}
	return nil
}

// nextAppRoots is every directory in the project that Next reads a proxy file
// from: the root of each Next app. A Vite app has no such file, and neither
// has the API.
func nextAppRoots(root string, opts Options) []string {
	var out []string
	for _, dir := range []string{
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "apps", "admin"),
		webAppRoot(root, opts), // a single's frontend, which is the project root
	} {
		if !fileExists(filepath.Join(dir, "next.config.ts")) &&
			!fileExists(filepath.Join(dir, "next.config.js")) &&
			!fileExists(filepath.Join(dir, "next.config.mjs")) {
			continue
		}
		if !contains(out, dir) {
			out = append(out, dir)
		}
	}
	return out
}

func moveMiddlewareToProxy(root, app string) error {
	old := filepath.Join(app, "middleware.ts")
	if !fileExists(old) {
		return nil
	}
	shown := displayPath(root, old)

	// Both present is a project somebody is already migrating by hand. Next
	// reads one of them and ignores the other, which is the sort of thing worth
	// saying out loud rather than resolving on their behalf.
	if fileExists(filepath.Join(app, "proxy.ts")) {
		fmt.Printf("  ⚠ %s and proxy.ts both exist. Next reads proxy.ts and ignores the other; delete %s when you have moved what you need out of it.\n",
			shown, filepath.Base(old))
		return nil
	}

	body, err := os.ReadFile(old)
	if err != nil {
		return fmt.Errorf("reading %s: %w", shown, err)
	}

	// The manifest is not the question here, and asking it was wrong: this file
	// is written by grit add web-auth without a recording open, so every
	// project's copy reads as untracked, which IsUnchanged reports as edited.
	// The repair would have declined to move the one file it exists for, and
	// told people they had edited something they had not.
	//
	// What decides whether the rename is safe is the export. Renaming the file
	// is content-preserving, so a developer's own protected paths move with it
	// untouched; renaming the function is the only edit, and it needs the line
	// to be the one Grit wrote. A file exporting a default, or an arrow
	// function, or nothing at all, would move to a name Next reads and leave it
	// with no function to call, which is a broken build rather than a warning.
	const export = "export function middleware(request: NextRequest)"
	if !strings.Contains(string(body), export) {
		fmt.Printf("  ⚠ %s: Next 16 deprecated this filename, and this file does not export middleware the way Grit wrote it, so it was left alone.\n"+
			"    Rename it to proxy.ts and rename the function it exports from middleware to proxy. Nothing else changes.\n", shown)
		return nil
	}

	// Renamed rather than rewritten from a template: there are two of these
	// files, the web-auth one and the admin-only one, and the copy on disk is
	// the only thing that knows which this project has, let alone what has been
	// added to its list of protected paths.
	moved := strings.Replace(string(body), export,
		"export function proxy(request: NextRequest)", 1)

	if err := writeFile(filepath.Join(app, "proxy.ts"), moved); err != nil {
		return err
	}
	if err := os.Remove(old); err != nil {
		return fmt.Errorf("removing %s: %w", shown, err)
	}
	manifest.Drop(old)
	fmt.Printf("  ✓ %s is now proxy.ts, the name Next 16 reads without a deprecation notice. Same file, and the function it exports is called proxy.\n", shown)
	return nil
}

func displayPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return path
}
