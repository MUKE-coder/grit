package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// brand.css: the one stylesheet in a Grit project that Grit will never write
// over.
//
// Before this existed, an app with its own identity had nowhere to put it.
// globals.css holds the token palette and is framework-owned, so giving an
// app its own accent meant editing a file the upgrade delivers. The manifest
// guard then held that file back for ever: the app kept its colours and
// stopped receiving every later fix to the same file, and was warned about it
// on every upgrade. Ten blueprints were about to hit this at once, which is
// what made it worth solving rather than working around.
//
// So: globals.css stays framework-owned and keeps arriving, brand.css is
// yours, and it is imported last so it wins at equal specificity without
// anybody having to count selectors. The only rule to respect is that the
// overrides must out-specify the theme blocks they are overriding, and the
// template below is written so that uncommenting a block is enough.
//
// Written once. An upgrade creates it if it is missing and never touches it
// again, which is why it is not in the manifest: a file Grit does not write
// cannot be a file Grit holds back.

// brandCSSPath returns where brand.css belongs for a frontend rooted at
// appRoot, or "" when the app has no stylesheet entry point to sit beside.
//
// Decided by the file that is actually there rather than by a flag, the same
// way ensureShadcnConfigs decides: an App Router app keeps its CSS in app/, a
// Vite app in src/, and a project that switched frontends still gets the
// right answer.
func brandCSSPath(appRoot string) string {
	for _, dir := range []string{"app", "src"} {
		if fileExists(filepath.Join(appRoot, dir, "globals.css")) {
			return filepath.Join(appRoot, dir, "brand.css")
		}
	}
	return ""
}

// brandCSS is the seed file. It sets nothing: every rule is commented out, so
// a project that never opens it renders exactly as it did before.
//
// The selectors are not decoration. globals.css sets its tokens on
// [data-theme="atlas"] and its dark variants on a descendant pair, so a bare
// :root here would tie on specificity and lose or win by source order
// depending on how the bundler happened to emit the files. html[data-theme]
// beats the first (0,1,1 over 0,1,0) and the dark selector beats the second
// (0,2,1 over 0,2,0). Uncommenting a block is then enough, which is the whole
// point of shipping it filled in.
func brandCSS(appName string) string {
	return `/* brand.css: your app's colours. Grit never writes this file again.
 *
 * globals.css carries Grit's token palette and is replaced on every upgrade.
 * This file is imported after it and is yours, so anything set here survives
 * ` + "`grit upgrade`" + ` and still lets globals.css keep arriving with fixes.
 *
 * Everything below is commented out. Uncomment what you want to change and
 * leave the rest: a token you do not set keeps the value your theme gives it.
 *
 * The selectors look heavier than they need to be, and they do need to be:
 * they have to out-specify the [data-theme] blocks in globals.css. Change the
 * values, not the selectors.
 *
 *   ` + appName + `
 */

/* ── Light ──────────────────────────────────────────────────────────── */
/*
html[data-theme] {
  --accent: #5b47e0;
  --accent-hover: #4a38c9;
  --accent-fg: #ffffff;

  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --bg-tertiary: #f1f5f9;

  --success: #0b7a45;
  --warning: #a35f00;
  --danger: #c93636;
}
*/

/* ── Dark ───────────────────────────────────────────────────────────── */
/* Both selectors, because the two frontends mark dark mode differently:
 * the admin writes data-theme-mode on <html>, Tailwind writes .dark. */
/*
html[data-theme][data-theme-mode="dark"],
html[data-theme].dark {
  --accent: #a99bff;
  --accent-hover: #b9aeff;
  --accent-fg: #08090b;

  --bg-primary: #08090b;
  --bg-secondary: #0f1114;
  --bg-tertiary: #15181c;

  --success: #3dd68c;
  --warning: #f5a524;
  --danger: #ff6b6b;
}
*/

/* ── Anything else ──────────────────────────────────────────────────── */
/* Ordinary CSS from here down. It loads after Tailwind's layers, so a plain
 * selector here beats a utility class without !important. */
`
}

// writeBrandCSS creates brand.css for every frontend under root that does not
// already have one.
//
// Returns how many it created, for the upgrade's change count. Never an
// error: a frontend without a globals.css is a frontend without a palette to
// override, and that is not a failure.
//
// It does NOT go through writeFile. writeFile records a hash in the manifest,
// and a hash is how a file becomes "you have edited this, so I am leaving it
// alone and telling you so on every upgrade" — which is the exact outcome
// this file exists to prevent.
func writeBrandCSS(root string, opts Options) int {
	created := 0
	for _, app := range frontendRoots(root) {
		path := brandCSSPath(app)
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			continue // yours now
		}
		if err := os.WriteFile(path, []byte(brandCSS(opts.ProjectName)), 0o644); err != nil {
			continue
		}
		created++
	}
	return created
}

// ensureBrandCSSImport makes sure the app actually loads brand.css.
//
// Writing the file is half the job. It only wins over globals.css because it
// is imported AFTER it, and that import lives in the layout, which is
// framework-owned: on an upgrade, a project that has edited its layout keeps
// its own copy and never gets the line. The file then exists, contains the
// right rules, and changes nothing.
//
// That is what happened to the WhatsApp blueprint: painted green, still blue,
// with no warning anywhere. So the import is repaired in place instead of
// being left to the template.
//
// Anchored on the globals.css import rather than on a line number or the
// top of the file, because the anchor is the thing that has to be true for
// the ordering to work in the first place.
func ensureBrandCSSImport(root string) int {
	fixed := 0
	for _, app := range frontendRoots(root) {
		css := brandCSSPath(app)
		if css == "" {
			continue
		}
		// The entry file is whichever one pulls in globals.css.
		for _, rel := range []string{
			filepath.Join("app", "layout.tsx"),
			filepath.Join("src", "main.tsx"),
			filepath.Join("src", "index.tsx"),
			filepath.Join("src", "App.tsx"),
		} {
			path := filepath.Join(app, rel)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			body := string(data)
			if strings.Contains(body, `"./brand.css"`) {
				break // already there
			}
			const anchor = `import "./globals.css";`
			if !strings.Contains(body, anchor) {
				continue
			}
			body = strings.Replace(body, anchor,
				anchor+"\n// Yours, and never written over. Imported last so what it sets wins.\n"+
					`import "./brand.css";`, 1)
			if err := os.WriteFile(path, []byte(body), 0o644); err == nil {
				fixed++
			}
			break
		}
	}
	return fixed
}

// frontendRoots lists the directories that could hold a frontend, whatever
// the architecture: the two apps of a triple, the one of a double, the
// embedded SPA of a single, and the desktop client.
//
// Probed rather than derived from Options.Architecture, because an upgrade
// reads the architecture off the directory layout anyway and a project that
// has grown an app since it was scaffolded should still get the file.
func frontendRoots(root string) []string {
	candidates := []string{
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "desktop", "frontend"),
		filepath.Join(root, "desktop", "frontend"),
		filepath.Join(root, "frontend"),
		root,
	}
	var found []string
	for _, c := range candidates {
		if dirExists(c) {
			found = append(found, c)
		}
	}
	return found
}
