package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// dompurifyVersion is the browser-side sanitiser the blog pages render through.
const dompurifyVersion = "^3.4.15"

// ensureRichTextSafety brings a project scaffolded before v3.242.0 up to the fix
// for stored XSS (H3 in the contact-app review).
//
// The demo blog stored a post's HTML as it was sent and rendered it with
// dangerouslySetInnerHTML, on the same origin as the admin panel, so anybody
// allowed to edit posts could run script in an admin's browser. Rich text is
// now sanitised in the API on every write, and again in the browser.
func ensureRichTextSafety(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if fileExists(filepath.Join(apiRoot, "go.mod")) {
		gomod, err := os.ReadFile(filepath.Join(apiRoot, "go.mod"))
		if err != nil {
			return fmt.Errorf("reading go.mod: %w", err)
		}
		// The package first, so go get finds an import and records a direct
		// dependency rather than an indirect one.
		for path, content := range sanitizeFiles(apiRoot) {
			if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
		}
		if !strings.Contains(string(gomod), bluemondayModule+" ") {
			if err := goGet(apiRoot, bluemondayModule+"@"+bluemondayVersion); err != nil {
				fmt.Printf("  ⚠ could not add %s, so run `go get %s@%s` in %s: %v\n",
					bluemondayModule, bluemondayModule, bluemondayVersion, apiRoot, err)
			}
		}
		if err := EnsureSanitizeWiring(apiRoot, opts.Module()); err != nil {
			fmt.Printf("  ⚠ %v\n", err)
		}

		if blog := filepath.Join(apiRoot, "internal", "models", "blog.go"); fileExists(blog) {
			m, err := manifest.Load(root)
			if err != nil {
				return err
			}
			if err := repairSourceFile(root, m, blog, repairBlogContentSource); err != nil {
				return err
			}
		}
	}

	for _, path := range []string{
		filepath.Join(root, "apps", "web", "package.json"),
		filepath.Join(root, "frontend", "package.json"),
		filepath.Join(root, "package.json"),
	} {
		added, err := ensureDependency(path, "dompurify", dompurifyVersion)
		if err != nil {
			return err
		}
		if added {
			rel, _ := filepath.Rel(root, path)
			fmt.Printf("  ✓ %s: added dompurify, which the blog page renders through\n", filepath.ToSlash(rel))
		}
	}
	return nil
}

var blogContentRe = regexp.MustCompile("(?m)^(\\s*Content\\s+string\\s+`gorm:\"type:text\" json:\"content\")`")

// repairBlogContentSource tags the demo Blog's Content for the sanitiser.
func repairBlogContentSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "type Blog struct") || strings.Contains(src, `sanitize:"html"`) {
		return src, nil, nil
	}
	out := blogContentRe.ReplaceAllString(src, "${1} sanitize:\"html\"`")
	if out == src {
		return src, nil, []string{"Blog.Content is not the line Grit generated, so it is not sanitised: add sanitize:\"html\" to its tag"}
	}
	return out, []string{"blog content is sanitised on its way into the database"}, nil
}

// ensureDependency adds name at version to the dependencies of a frontend
// package.json that lacks it. A package.json with no React in it (a monorepo
// root) is not a frontend and is left alone.
func ensureDependency(path, name, version string) (bool, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", path, err)
	}
	src := string(raw)
	if !strings.Contains(src, `"react":`) || strings.Contains(src, `"`+name+`":`) {
		return false, nil
	}
	const anchor = `"dependencies": {`
	i := strings.Index(src, anchor)
	if i < 0 {
		return false, nil
	}
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	at := i + len(anchor)
	entry := nl + `    "` + name + `": "` + version + `"`
	if rest := strings.TrimLeft(src[at:], " \t\r\n"); !strings.HasPrefix(rest, "}") {
		entry += ","
	}
	if err := os.WriteFile(path, []byte(src[:at]+entry+src[at:]), 0o644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}
