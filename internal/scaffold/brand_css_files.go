package scaffold

import (
	"fmt"
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
func brandCSS(appName string, light, dark map[string]string) string {
	return `/* brand.css: your app's colours. Grit never writes this file again.
 *
 * globals.css carries Grit's token palette and is replaced on every upgrade.
 * This file is imported after it and is yours, so anything set here survives
 * ` + "`grit upgrade`" + ` and still lets globals.css keep arriving with fixes.
 *
 * Everything below is commented out, and the values shown are the ones this
 * app renders today, read from its own globals.css. Uncommenting a block
 * changes nothing until you edit a value, which is the point: it is a
 * starting position, not a repaint.
 *
 * The selectors look heavier than they need to be, and they do need to be:
 * they have to out-specify the [data-theme] blocks in globals.css. Change the
 * values, not the selectors.
 *
 *   ` + appName + `
 */

/* ── Light ──────────────────────────────────────────────────────────── */
` + brandBlock("html[data-theme]", light) + `

/* ── Dark ───────────────────────────────────────────────────────────── */
/* Both selectors, because the two frontends mark dark mode differently:
 * the admin writes data-theme-mode on <html>, Tailwind writes .dark. */
` + brandBlock("html[data-theme][data-theme-mode=\"dark\"],\nhtml[data-theme].dark", dark) + `

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
		light, dark := readBrandTokens(app)
		if err := os.WriteFile(path, []byte(brandCSS(opts.ProjectName, light, dark)), 0o644); err != nil {
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

// The tokens the example offers, in the order they read best: what an app
// changes first, then its surfaces, then its status colours.
var brandExampleTokens = []string{
	"accent", "accent-hover", "accent-fg",
	"", // a blank line in the output
	"bg-primary", "bg-secondary", "bg-tertiary",
	"",
	"success", "warning", "danger",
}

// brandTokenLines renders the example's declarations from what the app's own
// globals.css sets, and shows nothing else.
//
// A token the app does not set is left out rather than filled in from Grit's
// palette. The web app has no dark palette at all, so a dark block seeded
// from defaults showed a colour that app has never rendered, which is the
// same fault this is fixing in a smaller place (grit#135).
func brandTokenLines(found map[string]string) string {
	var b strings.Builder
	pending := false
	for _, name := range brandExampleTokens {
		if name == "" {
			pending = true
			continue
		}
		value := found[name]
		if value == "" {
			continue
		}
		if pending && b.Len() > 0 {
			b.WriteString("\n")
		}
		pending = false
		fmt.Fprintf(&b, "  --%s: %s;\n", name, value)
	}
	return b.String()
}

// readBrandTokens reads the light and dark values an app's globals.css sets
// for the tokens the example offers.
//
// Read from the project rather than from the theme palette, because a project
// that edited globals.css has colours no palette knows. That is exactly the
// project this file is for, and the one that reads "your colours" over Grit's
// purple and concludes an upgrade is about to repaint it (grit#135).
//
// The dark values are whichever the file sets last under a dark selector:
// globals.css writes its dark block after its light one, so a second value
// for the same token is the dark one. Crude, and right for the file it reads.
func readBrandTokens(appRoot string) (light, dark map[string]string) {
	light, dark = map[string]string{}, map[string]string{}

	var path string
	for _, dir := range []string{"app", "src"} {
		candidate := filepath.Join(appRoot, dir, "globals.css")
		if fileExists(candidate) {
			path = candidate
			break
		}
	}
	if path == "" {
		return light, dark
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return light, dark
	}

	wanted := map[string]bool{}
	for _, name := range brandExampleTokens {
		if name != "" {
			wanted[name] = true
		}
	}

	inDark := false
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.Contains(trimmed, "data-theme-mode=\"dark\""),
			strings.Contains(trimmed, ".dark"),
			strings.Contains(trimmed, "prefers-color-scheme: dark"):
			inDark = true
		case strings.HasPrefix(trimmed, "html[data-theme"),
			strings.HasPrefix(trimmed, ":root"):
			inDark = false
		}

		if !strings.HasPrefix(trimmed, "--") {
			continue
		}
		name, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			continue
		}
		name = strings.TrimPrefix(strings.TrimSpace(name), "--")
		value = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), ";"))
		if !wanted[name] || value == "" || strings.HasPrefix(value, "var(") {
			continue
		}
		if inDark {
			dark[name] = value
		} else if _, seen := light[name]; !seen {
			light[name] = value
		}
	}
	return light, dark
}

// brandBlock renders one commented example block, or a note when the app's
// stylesheet sets none of these tokens.
//
// The note is a comment of its own rather than a line inside the block,
// because CSS comments do not nest: an inner */ ends the outer comment, and
// the block's closing brace becomes live CSS in a file whose whole promise
// is that it declares nothing.
func brandBlock(selector string, found map[string]string) string {
	lines := brandTokenLines(found)
	if strings.TrimSpace(lines) == "" {
		return "/* This app's stylesheet sets none of these tokens here, so there is\n" +
			" * nothing to copy. Write the ones you want, on " + selector + ". */\n"
	}
	return "/*\n" + selector + " {\n" + lines + "}\n*/\n"
}
