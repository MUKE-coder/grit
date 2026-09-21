package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairLintTooling moves an existing project from ESLint and Prettier to
// Biome, conservatively.
//
// Added: biome.jsonc where there is no Biome config yet, the pinned Biome
// devDependency where pnpm installs, and Biome lint and format scripts in each
// app. An app's lint script is switched only when it is one Grit wrote
// ("next lint", which Next.js 16 removed, or "eslint .", with no ESLint
// installed) and the project has no ESLint setup of its own; an app with its
// own ESLint config or dependency keeps its lint script.
//
// Removed: the .prettierrc and .prettierignore Grit wrote, and the Prettier
// devDependencies, only when those files are exactly what Grit wrote. A
// project that edited them keeps Prettier and its format script.
//
// Grit never wrote an ESLint config or installed ESLint, so there is nothing
// of ESLint's for this to remove.
//
// before is each app's package.json as it was when the upgrade started. An
// unedited package.json is rewritten whole from the template earlier in the
// upgrade, which would take the lint script from an app with its own ESLint,
// and Prettier from a project that kept it; those are put back from here.
func repairLintTooling(root string, opts Options, before map[string]string) error {
	jsRoot, apps := lintToolingLayout(root, opts)
	if jsRoot == "" {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	keepPrettier, prettierFiles := prettierStatus(root, m)
	if !keepPrettier && len(prettierFiles) > 0 {
		for _, path := range prettierFiles {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("removing %s: %w", path, err)
			}
			manifest.Drop(path)
		}
		fmt.Printf("  ✓ Removed Grit's .prettierrc and .prettierignore: Biome formats the project now (pnpm format)\n")
	} else if keepPrettier {
		fmt.Printf("  • .prettierrc or .prettierignore has your edits, so Prettier stays and pnpm format still runs it\n")
	}

	if !fileExists(filepath.Join(jsRoot, "biome.json")) && !fileExists(filepath.Join(jsRoot, "biome.jsonc")) {
		if err := writeFile(filepath.Join(jsRoot, biomeConfigFile), biomeConfig(jsRoot != root)); err != nil {
			return fmt.Errorf("writing %s: %w", biomeConfigFile, err)
		}
		fmt.Printf("  ✓ %s: Biome's lint rules and formatter settings\n", relShown(root, filepath.Join(jsRoot, biomeConfigFile)))
	}

	// Where pnpm installs: the workspace root, or a single app's frontend/.
	rootPkg := filepath.Join(jsRoot, "package.json")
	if err := repairTextFile(root, m, rootPkg, func(src string) (string, []string, []string) {
		return biomeRootPackageSource(src, jsRoot == root)
	}); err != nil {
		return err
	}

	for _, app := range apps {
		pkg := filepath.Join(app, "package.json")
		if !fileExists(pkg) {
			continue
		}
		own := ownESLintSetup(root, app)
		if own != "" {
			fmt.Printf("  • %s keeps its own ESLint setup (%s): its lint script was left alone\n", relShown(root, app), own)
		}
		if err := repairTextFile(root, m, pkg, func(src string) (string, []string, []string) {
			src, restored := restoreOwnTooling(src, before[app], own != "", keepPrettier)
			out, fixed, warnings := biomeAppPackageSource(src, own != "", keepPrettier)
			return out, append(restored, fixed...), warnings
		}); err != nil {
			return err
		}
	}
	return nil
}

// lintToolingSnapshot reads each app's package.json before the upgrade
// rewrites any of them.
func lintToolingSnapshot(root string, opts Options) map[string]string {
	_, apps := lintToolingLayout(root, opts)
	out := map[string]string{}
	for _, app := range apps {
		if raw, err := os.ReadFile(filepath.Join(app, "package.json")); err == nil {
			out[app] = string(raw)
		}
	}
	return out
}

// restoreOwnTooling puts back the lint script of an app with its own ESLint,
// and the Prettier format script and devDependencies of a project that kept
// Prettier, when the template replaced them earlier in the upgrade.
func restoreOwnTooling(src, before string, ownESLint, keepPrettier bool) (string, []string) {
	if before == "" {
		return src, nil
	}
	var restored []string
	put := func(name string) {
		was, had := scriptValue(before, name)
		now, has := scriptValue(src, name)
		if !had || (has && now == was) {
			return
		}
		if has {
			src = setScriptValue(src, name, was)
		} else if next, ok := addPackageEntry(src, "scripts", `"`+name+`": "`+was+`"`); ok {
			src = next
		}
		restored = append(restored, name+` is "`+was+`" again, as you had it`)
	}
	if ownESLint {
		put("lint")
	}
	if keepPrettier {
		if was, _ := scriptValue(before, "format"); strings.Contains(was, "prettier") {
			put("format")
		}
		var pkg struct {
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal([]byte(before), &pkg) == nil {
			for _, dep := range []string{"prettier", "prettier-plugin-tailwindcss"} {
				version, had := pkg.DevDependencies[dep]
				if !had || strings.Contains(src, `"`+dep+`":`) {
					continue
				}
				if next, ok := addPackageEntry(src, "devDependencies", `"`+dep+`": "`+version+`"`); ok {
					src = next
					restored = append(restored, dep+" kept in devDependencies")
				}
			}
		}
	}
	return src, restored
}

// lintToolingLayout is where pnpm installs and which apps have lint scripts.
// An API-only project has neither.
func lintToolingLayout(root string, opts Options) (string, []string) {
	if opts.Architecture == ArchSingle || (!fileExists(filepath.Join(root, "pnpm-workspace.yaml")) && fileExists(filepath.Join(root, "frontend", "package.json"))) {
		fe := filepath.Join(root, "frontend")
		if !fileExists(filepath.Join(fe, "package.json")) {
			return "", nil
		}
		return fe, []string{fe}
	}
	if !fileExists(filepath.Join(root, "pnpm-workspace.yaml")) || !fileExists(filepath.Join(root, "package.json")) {
		return "", nil
	}
	return root, []string{
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "desktop", "frontend"),
	}
}

// prettierStatus reports whether Prettier stays, and the Prettier files Grit
// wrote that can go when it does not. A file is Grit's when the manifest says
// it is unedited or, in a project older than the manifest, when it is byte for
// byte what Grit wrote.
func prettierStatus(root string, m *manifest.Manifest) (keep bool, files []string) {
	for _, f := range []struct {
		name string
		body string
	}{{".prettierrc", prettierConfig()}, {".prettierignore", prettierIgnore()}} {
		path := filepath.Join(root, f.name)
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rel, _ := manifest.Rel(root, path)
		status := m.StatusOf(root, rel)
		unedited := status == manifest.Unchanged ||
			(status == manifest.Untracked && manifest.Hash(string(raw)) == manifest.Hash(f.body))
		if !unedited {
			keep = true
		}
		files = append(files, path)
	}
	return keep, files
}

// ownESLintSetup names the ESLint config or dependency an app has, or "" when
// it has none. Grit never wrote either, so any it finds is the project's.
func ownESLintSetup(root, app string) string {
	names := []string{
		"eslint.config.js", "eslint.config.mjs", "eslint.config.cjs",
		"eslint.config.ts", "eslint.config.mts", "eslint.config.cts",
		".eslintrc", ".eslintrc.js", ".eslintrc.cjs", ".eslintrc.json",
		".eslintrc.yml", ".eslintrc.yaml",
	}
	for _, dir := range []string{app, root} {
		for _, n := range names {
			if fileExists(filepath.Join(dir, n)) {
				return relShown(root, filepath.Join(dir, n))
			}
		}
	}
	for _, dir := range []string{app, root} {
		raw, err := os.ReadFile(filepath.Join(dir, "package.json"))
		if err != nil {
			continue
		}
		src := string(raw)
		for _, dep := range []string{`"eslint":`, `"eslint-config-next":`, `"eslintConfig":`} {
			if strings.Contains(src, dep) {
				return strings.Trim(strings.TrimSuffix(dep, ":"), `"`) + " in " + relShown(root, filepath.Join(dir, "package.json"))
			}
		}
	}
	return ""
}

// The scripts Grit wrote before Biome, which are the only ones replaced.
var gritLintScripts = map[string]bool{"next lint": true, "eslint .": true}

const gritPrettierFormatScript = "prettier --write ."

// biomeRootPackageSource adds the Biome devDependency where pnpm installs and,
// at a workspace root, a format script.
func biomeRootPackageSource(src string, workspaceRoot bool) (string, []string, []string) {
	out := src
	var fixed, warnings []string
	if !strings.Contains(out, `"@biomejs/biome"`) {
		next, ok := addPackageEntry(out, "devDependencies", biomeDevDependency)
		if !ok {
			return src, nil, []string{`no devDependencies block: add ` + biomeDevDependency + ` to it, then pnpm install`}
		}
		out = next
		fixed = append(fixed, "Biome "+biomeVersion+" added to devDependencies (run pnpm install)")
	}
	if workspaceRoot {
		if _, has := scriptValue(out, "format"); !has {
			if next, ok := addFormatScript(out); ok {
				out = next
				fixed = append(fixed, "pnpm format runs Biome")
			}
		}
	}
	if !json.Valid([]byte(out)) {
		return src, nil, []string{"could not add Biome without breaking the JSON: add " + biomeDevDependency + " to devDependencies by hand"}
	}
	return out, fixed, warnings
}

// biomeAppPackageSource points an app's lint and format scripts at Biome and
// drops the Prettier devDependencies once nothing runs Prettier.
func biomeAppPackageSource(src string, ownESLint, keepPrettier bool) (string, []string, []string) {
	out := src
	var fixed, warnings []string

	lint, hasLint := scriptValue(out, "lint")
	switch {
	case ownESLint && lint == "next lint":
		warnings = append(warnings, `lint runs "next lint", which Next.js 16 removed: point it at "eslint ." to use your ESLint setup`)
	case ownESLint:
		// Reported by the caller; the project's lint is its own business.
	case !hasLint:
		if next, ok := addLintScript(out); ok {
			out = next
			fixed = append(fixed, "pnpm lint runs Biome")
		}
	case gritLintScripts[lint]:
		out = setScriptValue(out, "lint", biomeLintScript)
		fixed = append(fixed, "pnpm lint runs Biome instead of "+lint+", which could not run")
	case lint != biomeLintScript && !strings.HasPrefix(lint, "biome "):
		warnings = append(warnings, `lint runs "`+lint+`", which is yours and was left alone`)
	}

	format, hasFormat := scriptValue(out, "format")
	switch {
	case !hasFormat:
		if next, ok := addFormatScript(out); ok {
			out = next
			fixed = append(fixed, "pnpm format runs Biome")
		}
	case format == gritPrettierFormatScript && !keepPrettier:
		out = setScriptValue(out, "format", biomeFormatScript)
		fixed = append(fixed, "pnpm format runs Biome instead of Prettier")
	}

	// Only once no script runs Prettier, including any the project added.
	if scripts, _, _ := packageBlock(out, "scripts"); !keepPrettier && !strings.Contains(scripts, "prettier") {
		for _, dep := range []string{"prettier", "prettier-plugin-tailwindcss"} {
			if next, ok := removePackageEntry(out, dep); ok {
				out = next
				fixed = append(fixed, dep+" removed from devDependencies")
			}
		}
	}
	if !json.Valid([]byte(out)) {
		return src, nil, []string{"could not switch lint and format to Biome without breaking the JSON: point them at biome by hand"}
	}
	return out, fixed, warnings
}

// scriptValue is the value of a script in package.json text.
func scriptValue(src, name string) (string, bool) {
	block, _, ok := packageBlock(src, "scripts")
	if !ok {
		return "", false
	}
	m := regexp.MustCompile(`"` + regexp.QuoteMeta(name) + `"\s*:\s*"((?:[^"\\]|\\.)*)"`).FindStringSubmatch(block)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// setScriptValue replaces the value of a script that exists.
func setScriptValue(src, name, value string) string {
	block, start, ok := packageBlock(src, "scripts")
	if !ok {
		return src
	}
	re := regexp.MustCompile(`("` + regexp.QuoteMeta(name) + `"\s*:\s*)"(?:[^"\\]|\\.)*"`)
	loc := re.FindStringSubmatchIndex(block)
	if loc == nil {
		return src
	}
	replaced := block[:loc[0]] + block[loc[2]:loc[3]] + `"` + value + `"` + block[loc[1]:]
	return src[:start] + replaced + src[start+len(block):]
}

// packageBlock is the text of a top-level object such as "scripts", from its
// opening brace to its closing one, and where it starts.
func packageBlock(src, key string) (string, int, bool) {
	anchor := regexp.MustCompile(`"` + regexp.QuoteMeta(key) + `"\s*:\s*\{`).FindStringIndex(src)
	if anchor == nil {
		return "", 0, false
	}
	open := anchor[1] - 1
	depth, inString := 0, false
	for i := open; i < len(src); i++ {
		switch c := src[i]; {
		case inString && c == '\\':
			i++
		case c == '"':
			inString = !inString
		case inString:
		case c == '{':
			depth++
		case c == '}':
			depth--
			if depth == 0 {
				return src[open : i+1], open, true
			}
		}
	}
	return "", 0, false
}

// addPackageEntry puts entry first in a top-level object, in the file's own
// indentation and line endings.
func addPackageEntry(src, key, entry string) (string, bool) {
	block, start, ok := packageBlock(src, key)
	if !ok {
		return src, false
	}
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	indent := "    "
	if m := regexp.MustCompile(`\{\r?\n([ \t]+)"`).FindStringSubmatch(block); m != nil {
		indent = m[1]
	}
	inner := strings.TrimSpace(block[1 : len(block)-1])
	line := nl + indent + entry
	if inner != "" {
		line += ","
	} else {
		// An empty object: close it on its own line again.
		outer := strings.TrimSuffix(indent, "  ")
		return src[:start] + "{" + line + nl + outer + "}" + src[start+len(block):], true
	}
	return src[:start+1] + line + src[start+1:], true
}

// lintScriptLine is a lint script line: its indent, the value's end, whether a
// comma follows, and the line break.
var lintScriptLine = regexp.MustCompile(`([ \t]*)"lint"\s*:\s*"(?:[^"\\]|\\.)*"()(,?)[ \t]*(\r?\n)`)

// addFormatScript adds the Biome format script right after lint, where the
// templates have it, or first in scripts when there is no lint to sit beside.
func addFormatScript(src string) (string, bool) {
	block, start, ok := packageBlock(src, "scripts")
	if !ok {
		return src, false
	}
	m := lintScriptLine.FindStringSubmatchIndex(block)
	if m == nil {
		return addPackageEntry(src, "scripts", `"format": "`+biomeFormatScript+`"`)
	}
	indent, nl := block[m[2]:m[3]], block[m[8]:m[9]]
	entry := indent + `"format": "` + biomeFormatScript + `"`
	if m[6] != m[7] {
		// lint has a comma, so format goes on the next line with one of its own.
		at := start + m[1]
		return src[:at] + entry + "," + nl + src[at:], true
	}
	// lint was last: it takes a comma, and format becomes last.
	at := start + m[4]
	return src[:at] + "," + nl + entry + src[at:], true
}

var formatScriptLine = regexp.MustCompile(`(?m)^([ \t]*)"format"\s*:`)

// addLintScript adds the Biome lint script just before format, where the
// templates have it, or first in scripts.
func addLintScript(src string) (string, bool) {
	block, start, ok := packageBlock(src, "scripts")
	if !ok {
		return src, false
	}
	entry := `"lint": "` + biomeLintScript + `"`
	m := formatScriptLine.FindStringSubmatchIndex(block)
	if m == nil {
		return addPackageEntry(src, "scripts", entry)
	}
	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	at := start + m[0]
	return src[:at] + block[m[2]:m[3]] + entry + "," + nl + src[at:], true
}

// removePackageEntry deletes the line of a dependency, and the comma before it
// when it was the last entry of its object.
func removePackageEntry(src, name string) (string, bool) {
	re := regexp.MustCompile(`(?m)^[ \t]*"` + regexp.QuoteMeta(name) + `"\s*:\s*"[^"]*"(,?)[ \t]*\r?\n`)
	loc := re.FindStringSubmatchIndex(src)
	if loc == nil {
		return src, false
	}
	out := src[:loc[0]] + src[loc[1]:]
	if loc[2] == loc[3] {
		// It was the last entry, so the one before it now is, and keeps no comma.
		before := strings.TrimRight(out[:loc[0]], " \t\r\n")
		if strings.HasSuffix(before, ",") {
			out = before[:len(before)-1] + out[len(before):]
		}
	}
	return out, true
}

// relShown is a path as upgrade prints it: relative to the project, slashes.
func relShown(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return path
}
