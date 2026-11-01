package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// cleanBlankSingle finishes the blank starter on the single layouts.
//
// `grit remove resource Blog` knows the monorepo roots and the shared pieces of
// a single, but the panel and the public site of a single live in roots the
// removal's path list does not include yet. The API, the hook and the shared
// schema go; the panel's resource definition and the links that point at the
// removed pages stay. The project is new and the files are the scaffold's own,
// which is why a targeted sweep is safe here; the alternative is a blank
// starter that does not compile.
func cleanBlankSingle(root, frontend string) error {
	layout := singleBlankPathsFor(root, frontend)

	for _, path := range layout.files {
		if err := removeFile(path); err != nil {
			return err
		}
	}
	for _, path := range layout.dirs {
		if err := removeEmptyDir(path); err != nil {
			return err
		}
	}

	// The resource registry keeps the import and the array entry side by side.
	// Dropping both lines is what stops the panel from asking for a file that is
	// no longer there.
	if err := removeLinesContaining(layout.resourcesIndex, "blogsResource"); err != nil {
		return err
	}

	// The navbar's Blog entries, and the placeholders that borrowed /blog as a
	// destination, would otherwise point at a route this starter just removed.
	if err := removeBlogNavLinks(layout.navMenu); err != nil {
		return err
	}

	// The Vite landing page's primary button says "Read the Blog". Point it at
	// the panel instead of leaving a 404 in the first thing a new project shows.
	// The Next landing page carries the scaffold's removal markers, so the
	// removal already cut its CTA.
	if layout.home != "" {
		if err := replaceInFile(layout.home,
			`to="/blog"`, `to="/admin/dashboard"`,
			"Read the Blog", "Open the admin panel"); err != nil {
			return err
		}
	}

	// A comment in the trash screen names the table the Blog used. It is not a
	// compile error, but the blank starter should not mention a resource it
	// removed.
	if layout.trash != "" {
		if err := removeLinesContaining(layout.trash, "blog_posts -> Blog posts"); err != nil {
			return err
		}
	}

	return nil
}

// singleBlankPaths is where the two single layouts keep the pieces the removal
// does not reach.
type singleBlankPaths struct {
	files          []string
	dirs           []string
	resourcesIndex string
	navMenu        string
	home           string
	trash          string
}

func singleBlankPathsFor(root, frontend string) singleBlankPaths {
	if frontend == string(scaffold.FrontendNext) {
		// A Next single is the whole project: app/ at the root, and the panel
		// beside it in admin-panel/ (no src/). The removal handles the public
		// /blog pages and the landing page CTA; the panel's resource registry
		// and the navbar links are left behind.
		adminRoot := filepath.Join(root, "admin-panel")
		return singleBlankPaths{
			files: []string{
				filepath.Join(adminRoot, "resources", "blogs", "blogs.ts"),
				filepath.Join(adminRoot, "resources", "blogs", "blogs.custom.tsx"),
			},
			dirs:           []string{filepath.Join(adminRoot, "resources", "blogs")},
			resourcesIndex: filepath.Join(adminRoot, "resources", "index.ts"),
			navMenu:        filepath.Join(root, "lib", "nav-menu.ts"),
			trash:          filepath.Join(root, "app", "admin", "(dashboard)", "system", "trash", "page.tsx"),
		}
	}

	// A Vite single puts the panel at src/admin-panel and the public site under
	// src/routes/_site, which is why both the panel's page and the two public
	// blog routes survive the removal.
	adminRoot := filepath.Join(root, "src", "admin-panel")
	webRoot := filepath.Join(root, "src")
	return singleBlankPaths{
		files: []string{
			filepath.Join(adminRoot, "pages", "resources", "blogs.tsx"),
			filepath.Join(adminRoot, "resources", "blogs", "blogs.ts"),
			filepath.Join(adminRoot, "resources", "blogs", "blogs.custom.tsx"),
			filepath.Join(webRoot, "routes", "_site", "blog", "index.tsx"),
			filepath.Join(webRoot, "routes", "_site", "blog", "$slug.tsx"),
		},
		dirs: []string{
			filepath.Join(adminRoot, "resources", "blogs"),
			filepath.Join(webRoot, "routes", "_site", "blog"),
		},
		resourcesIndex: filepath.Join(adminRoot, "resources", "index.ts"),
		navMenu:        filepath.Join(webRoot, "lib", "nav-menu.ts"),
		home:           filepath.Join(webRoot, "routes", "_site", "index.tsx"),
		trash:          filepath.Join(adminRoot, "pages", "system-trash.tsx"),
	}
}

// cleanBlankMonorepoLinks removes the pieces of the Blog example that the
// removal does not reach on a monorepo blank: the navbar links and a comment in
// the trash screen. The files, the API routes and the landing page CTA are the
// removal's job, and it handles them for the monorepo roots.
func cleanBlankMonorepoLinks(root, frontend string, separateAdmin bool) error {
	webRoot := filepath.Join(root, "apps", "web")
	srcRoot := webRoot
	if frontend == string(scaffold.FrontendTanStack) {
		srcRoot = filepath.Join(webRoot, "src")

		// The removal deletes the Blog hook but cannot see the TanStack site's
		// route files, so the pages stay and import a file that is gone. That
		// is a build failure, not a 404: remove the routes with the hook.
		blogDir := filepath.Join(srcRoot, "routes", "_site", "blog")
		for _, name := range []string{"index.tsx", "$slug.tsx"} {
			if err := removeFile(filepath.Join(blogDir, name)); err != nil {
				return err
			}
		}
		if err := removeEmptyDir(blogDir); err != nil {
			return err
		}

		// The landing page's primary button points at the Blog. A double keeps
		// the panel at /admin/dashboard in this app. A triple's panel is a
		// separate app, so the button has to become a plain link and the unused
		// router Link import has to go with it.
		home := filepath.Join(srcRoot, "routes", "_site", "index.tsx")
		if separateAdmin {
			if err := replaceInFile(home,
				"<Link", "<a",
				"</Link>", "</a>",
				`to="/blog"`, `href="http://localhost:3001"`,
				"Read the Blog", "Open the admin panel",
				"import { createFileRoute, Link } from '@tanstack/react-router'",
				"import { createFileRoute } from '@tanstack/react-router'"); err != nil {
				return err
			}
		} else if err := replaceInFile(home,
			`to="/blog"`, `to="/admin/dashboard"`,
			"Read the Blog", "Open the admin panel"); err != nil {
			return err
		}

		// A triple keeps its panel in apps/admin, whose TanStack routes are
		// outside the removal's path list. The pages are gone, so the routes
		// that import them are a build failure.
		if separateAdmin {
			adminBlogRoutes := filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard", "resources", "blogs")
			for _, name := range []string{"index.tsx", "$id.tsx"} {
				if err := removeFile(filepath.Join(adminBlogRoutes, name)); err != nil {
					return err
				}
			}
			if err := removeEmptyDir(adminBlogRoutes); err != nil {
				return err
			}
		}
	}

	// The Expo app has its own Blog screens and two discover entries. The
	// resource removal knows the web and admin roots, not this one.
	expoRoot := filepath.Join(root, "apps", "expo")
	for _, path := range []string{
		filepath.Join(expoRoot, "app", "blogs", "index.tsx"),
		filepath.Join(expoRoot, "app", "blogs", "new.tsx"),
		filepath.Join(expoRoot, "hooks", "use-blogs.ts"),
	} {
		if err := removeFile(path); err != nil {
			return err
		}
	}
	if err := removeEmptyDir(filepath.Join(expoRoot, "app", "blogs")); err != nil {
		return err
	}
	if err := removeLinesContaining(filepath.Join(expoRoot, "app", "(tabs)", "explore.tsx"),
		`title: "Blogs"`, `route: "/blogs"`); err != nil {
		return err
	}
	if err := removeLinesContaining(filepath.Join(expoRoot, "app", "explore", "content.tsx"),
		`title: "Posts"`, `Blog posts`); err != nil {
		return err
	}

	// A full or desktop project has a quick-access card for the Blog's system
	// screen. The screen is removed with the resource, so the card would open a
	// route that does not exist.
	desktopQuickAccess := filepath.Join(root, "apps", "desktop", "frontend", "src", "components", "quick-access.tsx")
	if err := removeLinesContaining(desktopQuickAccess, `"sys:blog"`); err != nil {
		return err
	}
	if err := replaceInFile(desktopQuickAccess,
		"LayoutGrid, X, Home, RefreshCw, MessageSquare, FileText,",
		"LayoutGrid, X, Home, RefreshCw, MessageSquare,",
		"// cards. Support/blogs live under /system, so they get built-in entries.",
		"// cards. Support lives under /system, so it gets a built-in entry."); err != nil {
		return err
	}

	if err := removeBlogNavLinks(filepath.Join(srcRoot, "lib", "nav-menu.ts")); err != nil {
		return err
	}

	// A triple keeps the panel in apps/admin; a double embeds it in the web
	// app. Try every shape: the paths that do not exist are no-ops.
	for _, path := range []string{
		filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "trash", "page.tsx"),
		filepath.Join(webRoot, "app", "admin", "(dashboard)", "system", "trash", "page.tsx"),
		filepath.Join(srcRoot, "admin-panel", "pages", "system-trash.tsx"),
	} {
		if err := removeLinesContaining(path, "blog_posts -> Blog posts"); err != nil {
			return err
		}
	}

	// The scaffold's navbar test is a stub with a Blog link in it. It does not
	// exercise the real navbar, but a blank starter should not ship a test that
	// asserts a page the starter removed.
	testFile := filepath.Join(webRoot, "__tests__", "navbar.test.tsx")
	if err := removeLinesContaining(testFile, `"/blog"`, "name: /blog/i"); err != nil {
		return err
	}

	return nil
}

// removeFile removes a path, treating "already gone" as success so the cleanup
// stays idempotent.
func removeFile(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", path, err)
	}
	return nil
}

// removeEmptyDir removes a directory only when the cleanup left it empty. A
// non-empty directory means a file the cleanup did not expect, and leaving it
// alone is safer than deleting work the starter did not write.
func removeEmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	if len(entries) != 0 {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing %s: %w", path, err)
	}
	return nil
}

// removeLinesContaining drops every line that contains any of the needles. It
// is deliberately line based: these are generated files whose formatting the
// scaffold owns, and a parser would be more code than the problem deserves.
func removeLinesContaining(path string, needles ...string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	lines := strings.Split(string(data), "\n")
	kept := lines[:0]
	for _, line := range lines {
		drop := false
		for _, needle := range needles {
			if strings.Contains(line, needle) {
				drop = true
				break
			}
		}
		if !drop {
			kept = append(kept, line)
		}
	}

	if err := os.WriteFile(path, []byte(strings.Join(kept, "\n")), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// removeBlogNavLinks deletes the navbar link objects whose title is "Blog" and
// repoints the remaining /blog hrefs at the panel. The objects are small and
// flat, so matching the title line and walking out to its braces is enough.
func removeBlogNavLinks(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	lines := strings.Split(string(data), "\n")
	kept := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != `title: "Blog",` {
			kept = append(kept, lines[i])
			continue
		}

		start := i
		for start >= 0 && strings.TrimSpace(lines[start]) != "{" {
			start--
		}
		end := i
		for end < len(lines) &&
			strings.TrimSpace(lines[end]) != "}," &&
			strings.TrimSpace(lines[end]) != "}" {
			end++
		}

		// No braces means the shape changed. Keep the line rather than guess.
		if start < 0 || end >= len(lines) {
			kept = append(kept, lines[i])
			continue
		}

		// The lines from the opening brace up to the title are already in kept.
		drop := i - start
		if drop > len(kept) {
			drop = len(kept)
		}
		kept = kept[:len(kept)-drop]
		i = end
	}

	content := strings.Join(kept, "\n")
	content = strings.ReplaceAll(content, `href: "/blog",`, `href: ADMIN_URL,`)
	content = strings.ReplaceAll(content, `href: "/blog"`, `href: ADMIN_URL`)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return pruneUnusedLucideImports(path)
}

// pruneUnusedLucideImports removes identifiers from the lucide-react import
// block that no longer appear anywhere else in the file. TypeScript runs with
// noUnusedLocals, so an orphaned icon is a build failure.
func pruneUnusedLucideImports(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	content := string(data)

	start := strings.Index(content, "import {")
	if start < 0 {
		return nil
	}
	const endMarker = `} from "lucide-react";`
	relEnd := strings.Index(content[start:], endMarker)
	if relEnd < 0 {
		return nil
	}
	end := start + relEnd + len(endMarker)

	block := content[start:end]
	rest := content[:start] + content[end:]

	lines := strings.Split(block, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "import") || strings.HasPrefix(trimmed, "}") {
			kept = append(kept, line)
			continue
		}

		name := strings.TrimSuffix(trimmed, ",")
		name = strings.TrimPrefix(name, "type ")
		if identifierUsed(rest, name) {
			kept = append(kept, line)
		}
	}

	newContent := content[:start] + strings.Join(kept, "\n") + content[end:]
	if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// identifierUsed reports whether name appears as a whole word in src.
func identifierUsed(src, name string) bool {
	if name == "" {
		return false
	}
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
	return re.MatchString(src)
}

// replaceInFile applies string replacements in order, leaving the file alone
// when it does not exist. Each old value is replaced once.
func replaceInFile(path string, replacements ...string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	content := string(data)
	for i := 0; i+1 < len(replacements); i += 2 {
		content = strings.Replace(content, replacements[i], replacements[i+1], 1)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
