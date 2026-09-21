package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The package.json lines Grit wrote before Biome, for reconstructing old files.
const (
	oldPrettierDeps = "    \"prettier\": \"^3.3.0\",\n    \"prettier-plugin-tailwindcss\": \"^0.6.0\",\n"
	oldNextScripts  = "    \"lint\": \"next lint\",\n    \"format\": \"prettier --write .\",\n"
)

var jsoncComment = regexp.MustCompile(`(?m)^\s*//.*$`)

func TestBiomeConfigIsValidAndExplainsEveryRuleItTurnsOff(t *testing.T) {
	for _, inFrontend := range []bool{false, true} {
		cfg := biomeConfig(inFrontend)
		plain := jsoncComment.ReplaceAllString(cfg, "")
		var parsed map[string]any
		if err := json.Unmarshal([]byte(plain), &parsed); err != nil {
			t.Fatalf("biome.jsonc (inFrontend=%v) is not valid JSON once comments are removed: %v", inFrontend, err)
		}
		if !strings.Contains(cfg, `"$schema": "https://biomejs.dev/schemas/`+biomeVersion+`/schema.json"`) {
			t.Error("the schema URL is not the pinned Biome version")
		}
		if got := strings.Contains(cfg, `"root": ".."`); got != inFrontend {
			t.Errorf("inFrontend=%v: vcs root set = %v", inFrontend, got)
		}
		// Every rule switched off has a comment on the line above it, or shares
		// the comment of the rule directly above it.
		lines := strings.Split(cfg, "\n")
		for i, line := range lines {
			if !strings.Contains(line, `": "off"`) {
				continue
			}
			prev := strings.TrimSpace(lines[i-1])
			if !strings.HasPrefix(prev, "//") && !strings.Contains(prev, `": "off"`) {
				t.Errorf("rule switched off without a reason: %s", strings.TrimSpace(line))
			}
		}
	}
}

func TestNewProjectsUseBiomeNotESLintOrPrettier(t *testing.T) {
	root := t.TempDir()
	opts := tripleOptions()
	files := map[string]string{
		"root":          rootPackageJSON(opts),
		"web":           webPackageJSON(opts),
		"admin":         adminPackageJSON(opts),
		"admin (Vite)":  adminTanStackPackageJSON(opts),
		"web (Vite)":    webTanStackPackageJSON(opts),
		"single (Vite)": singleFrontendPackageJSON(opts),
		"desktop":       desktopClientPackageJSON(opts),
	}
	for name, src := range files {
		if !json.Valid([]byte(src)) {
			t.Errorf("%s package.json is not valid JSON", name)
		}
		for _, stale := range []string{"prettier", "eslint", "next lint"} {
			if strings.Contains(src, stale) {
				t.Errorf("%s package.json still mentions %q", name, stale)
			}
		}
		if name == "root" {
			if !strings.Contains(src, biomeDevDependency) || !strings.Contains(src, `"format": "`+biomeFormatScript+`"`) {
				t.Errorf("root package.json: want %s and a Biome format script", biomeDevDependency)
			}
			continue
		}
		if !strings.Contains(src, `"lint": "`+biomeLintScript+`"`) || !strings.Contains(src, `"format": "`+biomeFormatScript+`"`) {
			t.Errorf("%s package.json: lint and format do not run Biome", name)
		}
	}
	if !strings.Contains(files["single (Vite)"], biomeDevDependency) {
		t.Error("a single app installs from frontend/, so Biome belongs in its devDependencies")
	}

	// No eslint-disable comment survives in a frontend template: Biome reads its
	// own biome-ignore comments and ignores those.
	sources := map[string]map[string]string{
		"admin (Next)": adminFileMap(root, opts),
		"admin (Vite)": adminTanStackFileMap(root, viteTripleOptions()),
		"web (Next)":   webFileMap(root, opts),
	}
	for app, fileMap := range sources {
		for path, body := range fileMap {
			if strings.Contains(body, "eslint-disable") {
				t.Errorf("%s %s still has an eslint-disable comment", app, filepath.Base(path))
			}
		}
	}

	if !strings.Contains(ciYAML(opts), ciLintStep) {
		t.Error("ci.yml does not run pnpm lint")
	}
	if strings.Contains(turboJSON(), `"lint": {
      "dependsOn"`) {
		t.Error("turbo lint still builds every package first; Biome does not need them built")
	}
}

func TestBiomePackageRepairsProduceTheTemplates(t *testing.T) {
	opts := tripleOptions()

	rootNew := rootPackageJSON(opts)
	rootOld := strings.Replace(strings.Replace(rootNew, "    "+biomeDevDependency+",\n", "", 1),
		"    \"format\": \""+biomeFormatScript+"\",\n", "", 1)
	checkRepair(t, "root package.json", func(src string) (string, []string, []string) {
		return biomeRootPackageSource(src, true)
	}, rootOld, rootNew)

	appRepair := func(src string) (string, []string, []string) { return biomeAppPackageSource(src, false, false) }

	for name, tmpl := range map[string]string{"web": webPackageJSON(opts), "admin": adminPackageJSON(opts)} {
		scripts := "    \"lint\": \"" + biomeLintScript + "\",\n    \"format\": \"" + biomeFormatScript + "\",\n"
		old := strings.Replace(tmpl, scripts, oldNextScripts, 1)
		old = strings.Replace(old, "    \"postcss\": \"^8.4.0\",\n", "    \"postcss\": \"^8.4.0\",\n"+oldPrettierDeps, 1)
		checkRepair(t, name+" package.json", appRepair, old, tmpl)
	}

	adminVite := adminTanStackPackageJSON(opts)
	checkRepair(t, "admin (Vite) package.json", appRepair, strings.Replace(adminVite,
		"    \"lint\": \""+biomeLintScript+"\",\n    \"format\": \""+biomeFormatScript+"\",\n",
		"    \"lint\": \"eslint .\",\n", 1), adminVite)

	webVite := webTanStackPackageJSON(opts)
	checkRepair(t, "web (Vite) package.json", appRepair, strings.Replace(webVite,
		"    \"lint\": \""+biomeLintScript+"\",\n    \"format\": \""+biomeFormatScript+"\"\n",
		"    \"lint\": \"eslint .\"\n", 1), webVite)

	desktop := desktopClientPackageJSON(opts)
	oldDesktop := strings.Replace(desktop,
		"    \"lint\": \""+biomeLintScript+"\",\n    \"format\": \""+biomeFormatScript+"\"\n",
		"    \"format\": \"prettier --write .\"\n", 1)
	oldDesktop = strings.Replace(oldDesktop, "    \"postcss\": \"^8.4.49\",\n", "    \"postcss\": \"^8.4.49\",\n"+oldPrettierDeps, 1)
	checkRepair(t, "desktop package.json", appRepair, oldDesktop, desktop)
}

func TestBiomeAppRepairLeavesTheProjectsOwnTooling(t *testing.T) {
	opts := tripleOptions()
	scripts := "    \"lint\": \"" + biomeLintScript + "\",\n    \"format\": \"" + biomeFormatScript + "\",\n"
	old := strings.Replace(webPackageJSON(opts), scripts, oldNextScripts, 1)
	old = strings.Replace(old, "    \"postcss\": \"^8.4.0\",\n", "    \"postcss\": \"^8.4.0\",\n"+oldPrettierDeps, 1)

	// Its own ESLint, and an edited .prettierrc: nothing of either is touched.
	out, _, _ := biomeAppPackageSource(old, true, true)
	if out != old {
		t.Errorf("a project with its own ESLint and Prettier was changed:\n%s", out)
	}

	// A lint script of its own is reported, not replaced.
	custom := strings.Replace(old, `"lint": "next lint"`, `"lint": "tsc --noEmit && oxlint"`, 1)
	out, _, warnings := biomeAppPackageSource(custom, false, true)
	if !strings.Contains(out, `"lint": "tsc --noEmit && oxlint"`) || len(warnings) != 1 {
		t.Errorf("a custom lint script was replaced or not reported: %v", warnings)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// oldTripleProject writes the lint and format tooling of a project scaffolded
// before Biome.
func oldTripleProject(t *testing.T) (string, Options) {
	root := t.TempDir()
	opts := tripleOptions()
	rootNew := rootPackageJSON(opts)
	rootOld := strings.Replace(strings.Replace(rootNew, "    "+biomeDevDependency+",\n", "", 1),
		"    \"format\": \""+biomeFormatScript+"\",\n", "", 1)
	scripts := "    \"lint\": \"" + biomeLintScript + "\",\n    \"format\": \"" + biomeFormatScript + "\",\n"
	web := strings.Replace(webPackageJSON(opts), scripts, oldNextScripts, 1)
	web = strings.Replace(web, "    \"postcss\": \"^8.4.0\",\n", "    \"postcss\": \"^8.4.0\",\n"+oldPrettierDeps, 1)
	admin := strings.Replace(adminPackageJSON(opts), scripts, oldNextScripts, 1)
	admin = strings.Replace(admin, "    \"postcss\": \"^8.4.0\",\n", "    \"postcss\": \"^8.4.0\",\n"+oldPrettierDeps, 1)

	writeTestFile(t, filepath.Join(root, "pnpm-workspace.yaml"), pnpmWorkspace(false))
	writeTestFile(t, filepath.Join(root, "package.json"), rootOld)
	writeTestFile(t, filepath.Join(root, "apps", "web", "package.json"), web)
	writeTestFile(t, filepath.Join(root, "apps", "admin", "package.json"), admin)
	writeTestFile(t, filepath.Join(root, ".prettierrc"), prettierConfig())
	writeTestFile(t, filepath.Join(root, ".prettierignore"), prettierIgnore())
	return root, opts
}

func TestLintToolingRepairMovesAnOldProjectToBiome(t *testing.T) {
	root, opts := oldTripleProject(t)
	if err := repairLintTooling(root, opts, lintToolingSnapshot(root, opts)); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, filepath.Join(root, biomeConfigFile)); got != biomeConfig(false) {
		t.Error("biome.jsonc is not the template")
	}
	for _, f := range []string{".prettierrc", ".prettierignore"} {
		if fileExists(filepath.Join(root, f)) {
			t.Errorf("Grit's unedited %s was left behind", f)
		}
	}
	want := map[string]string{
		"package.json":            rootPackageJSON(opts),
		"apps/web/package.json":   webPackageJSON(opts),
		"apps/admin/package.json": adminPackageJSON(opts),
	}
	for rel, body := range want {
		if got := readTestFile(t, filepath.Join(root, filepath.FromSlash(rel))); got != body {
			t.Errorf("%s is not the template after the repair:\n%s", rel, got)
		}
	}

	// A second upgrade changes nothing.
	if err := repairLintTooling(root, opts, lintToolingSnapshot(root, opts)); err != nil {
		t.Fatal(err)
	}
	for rel, body := range want {
		if got := readTestFile(t, filepath.Join(root, filepath.FromSlash(rel))); got != body {
			t.Errorf("%s changed on a second run", rel)
		}
	}
}

func TestLintToolingRepairKeepsCustomisedESLintAndPrettier(t *testing.T) {
	root, opts := oldTripleProject(t)
	editedPrettier := strings.Replace(prettierConfig(), `"printWidth": 100`, `"printWidth": 120`, 1)
	writeTestFile(t, filepath.Join(root, ".prettierrc"), editedPrettier)
	eslintConfig := "export default [];\n"
	writeTestFile(t, filepath.Join(root, "apps", "web", "eslint.config.mjs"), eslintConfig)
	before := lintToolingSnapshot(root, opts)
	// Upgrade rewrites an unedited package.json from the template before the
	// repair runs, so the repair sees the template, not the project's scripts.
	writeTestFile(t, filepath.Join(root, "apps", "web", "package.json"), webPackageJSON(opts))
	writeTestFile(t, filepath.Join(root, "apps", "admin", "package.json"), adminPackageJSON(opts))

	if err := repairLintTooling(root, opts, before); err != nil {
		t.Fatal(err)
	}
	if readTestFile(t, filepath.Join(root, ".prettierrc")) != editedPrettier || !fileExists(filepath.Join(root, ".prettierignore")) {
		t.Error("an edited .prettierrc, or the .prettierignore beside it, was touched")
	}
	if readTestFile(t, filepath.Join(root, "apps", "web", "eslint.config.mjs")) != eslintConfig {
		t.Error("the project's own ESLint config was touched")
	}
	web := readTestFile(t, filepath.Join(root, "apps", "web", "package.json"))
	for _, want := range []string{`"lint": "next lint"`, `"format": "prettier --write ."`, `"prettier": "^3.3.0"`, `"prettier-plugin-tailwindcss": "^0.6.0"`} {
		if !strings.Contains(web, want) {
			t.Errorf("the web app keeps its ESLint and Prettier, but %s is gone:\n%s", want, web)
		}
	}
	if !json.Valid([]byte(web)) {
		t.Error("the web package.json is not valid JSON")
	}

	// The admin has no ESLint of its own: its lint could not run, so it moves to
	// Biome. Prettier stays, because the project edited its config.
	admin := readTestFile(t, filepath.Join(root, "apps", "admin", "package.json"))
	if !strings.Contains(admin, `"lint": "`+biomeLintScript+`"`) {
		t.Error("the admin's next lint, which Next.js 16 removed, was not replaced")
	}
	if !strings.Contains(admin, `"format": "prettier --write ."`) || !strings.Contains(admin, `"prettier": "^3.3.0"`) {
		t.Error("Prettier was taken from a project that edited its config")
	}
	// Whatever else happened, Biome is installed for the lint that now runs it.
	if !strings.Contains(readTestFile(t, filepath.Join(root, "package.json")), biomeDevDependency) {
		t.Error("lint runs Biome, and Biome is not installed")
	}
}

func TestLintToolingRepairOnASingleApp(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchSingle}
	tmpl := singleFrontendPackageJSON(opts)
	old := strings.Replace(tmpl, "    "+biomeDevDependency+",\n", "", 1)
	old = strings.Replace(old, ",\n    \"lint\": \""+biomeLintScript+"\",\n    \"format\": \""+biomeFormatScript+"\"\n", "\n", 1)
	if old == tmpl {
		t.Fatal("the reconstructed old package.json is the template")
	}
	writeTestFile(t, filepath.Join(root, "frontend", "package.json"), old)
	writeTestFile(t, filepath.Join(root, ".prettierrc"), prettierConfig())

	if err := repairLintTooling(root, opts, lintToolingSnapshot(root, opts)); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, filepath.Join(root, "frontend", biomeConfigFile)); got != biomeConfig(true) {
		t.Error("a single app's biome.jsonc belongs in frontend/, pointed at the .gitignore above it")
	}
	got := readTestFile(t, filepath.Join(root, "frontend", "package.json"))
	for _, want := range []string{biomeDevDependency, `"lint": "` + biomeLintScript + `"`, `"format": "` + biomeFormatScript + `"`} {
		if !strings.Contains(got, want) {
			t.Errorf("frontend/package.json lacks %s:\n%s", want, got)
		}
	}
	if !json.Valid([]byte(got)) {
		t.Error("frontend/package.json is not valid JSON after the repair")
	}
	if fileExists(filepath.Join(root, ".prettierrc")) {
		t.Error("Grit's unedited .prettierrc was left behind")
	}
}
