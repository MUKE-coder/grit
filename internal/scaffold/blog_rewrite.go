package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

// RewriteWebFilesWithoutTheBlog re-renders the web files whose templates
// mention the Blog example, for a project that has just removed the model.
//
// Called by `grit remove resource Blog`. Without it the removal leaves five
// /blog links in the navbar, because the nav menu is not on the removal's file
// list: it belongs to the web app rather than to the resource, and the resource
// removal has never known about it. `grit setup --starter blank` cleaned it up
// afterwards with a line-editing pass, so a project built that way was fine and
// a project where somebody ran `grit remove resource Blog` themselves was not.
//
// Rendered from the template rather than edited, which is the point. The
// template now omits its Blog parts when the model is gone, so the bytes this
// writes are the bytes the template produces, and `grit upgrade` adopts a file
// whose contents already match what it would write. The manifest stays
// truthful and the file keeps receiving framework fixes (grit#173).
//
// Deliberately not recording anything in the manifest. The recorded hash is the
// one from before the removal, so the file reads as edited until the next
// upgrade adopts it, and that is the safe order: a recorded hash that matched
// would hand the file to a template, and templates are how a removal gets
// undone.
func RewriteWebFilesWithoutTheBlog(root string) error {
	opts := projectOptions(root)

	// Only for a project that really has removed it. Called on one that still
	// has the model, this would write the file back with the Blog in it, which
	// is correct but pointless, and on a project mid-scaffold it would be
	// actively wrong.
	if !blogRemoved(root, opts) {
		return nil
	}

	// Every shape's web root, with the file names each uses. A path that is not
	// there is not an error: a project with no web app has nothing to rewrite.
	type target struct {
		path string
		body string
	}
	var targets []target

	webRoot := filepath.Join(root, "apps", "web")
	if dirExists(webRoot) {
		if opts.Frontend == FrontendTanStack {
			src := filepath.Join(webRoot, "src")
			targets = append(targets, target{
				filepath.Join(src, "lib", "nav-menu.ts"),
				nextToTanStack(webNavMenuConfig(opts, false)),
			})
		} else {
			targets = append(targets,
				target{
					filepath.Join(webRoot, "lib", "nav-menu.ts"),
					webNavMenuConfig(opts, false),
				},
				target{
					filepath.Join(webRoot, "app", "(marketing)", "page.tsx"),
					webLandingPage(opts, false),
				},
			)
		}
	}

	// A single keeps the frontend at the project root.
	if single := filepath.Join(root, "src", "lib", "nav-menu.ts"); fileExists(single) {
		targets = append(targets, target{single, nextToTanStack(webNavMenuConfig(opts, false))})
	}

	for _, t := range targets {
		if !fileExists(t.path) {
			continue
		}
		body := t.body
		// The web writer substitutes this; the Vite writers do not, and a
		// literal placeholder in a URL is what grit#173's sibling bug was.
		body = adminHref(body, opts)
		if err := os.WriteFile(t.path, []byte(body), 0o644); err != nil {
			return fmt.Errorf("rewriting %s without the blog: %w", t.path, err)
		}
	}
	return nil
}
