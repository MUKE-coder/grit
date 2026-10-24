package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// scoutVersion is the Scout release `grit add scout` installs.
//
// Pinned exactly rather than with a caret, for the reason recorded in
// next_version.go: a range adopted a one-day-old release once and broke every
// generated project. Scout is also pre-1.0, so a minor bump may move more than
// a minor bump usually does.
//
// A tag again. It was a commit for two days: scout.WithRoutes, which the
// whole integration rests on, landed after v0.3.0 was cut, so no tag could
// compile the file this command writes. v0.4.6 contains it.
//
// v0.4.9 is the release that closed Scout issues #2 to #6, all of them
// false positives this integration was producing on a correctly configured
// Grit app: "/api" in the sensitive-word list made every namespaced route
// sensitive, the dashboard-exposure check reported Scout's own three
// dashboards, and the code scanner matched the word "secret" rather than a
// secret-shaped value. A scan of a generated project went from 160 findings
// to 73, and 62 of what is left is the caret-range check on package.json.
//
// Config, WithTarget and WithRoutes are unchanged across the bump: checked
// against the module before raising it, because Scout moves fast.
const scoutVersion = "v0.4.9"

// Scout needs these, and a project generated before now is below them. Go
// would raise them silently on the first build; doing it here means the
// change is in the diff `grit add scout` produces rather than in a lockfile
// nobody reads.
var scoutRequires = []struct{ module, version string }{
	{"github.com/gin-gonic/gin", "v1.12.0"},
	{"gorm.io/gorm", "v1.31.2"},
}

// AddScout installs Scout into a project's API and mounts it at /scout.
//
// Opt-in rather than scaffolded into every project, unlike Sentinel and
// Pulse. Scout is pre-1.0, and a security scanner is a thing you choose to
// point at your application rather than something that should appear in it
// because you ran `grit new`.
//
// What makes the integration worth having is the route table. Scout's
// strongest check reads a route's middleware chain and reports a sensitive
// route that has no auth gate, and it cannot get that chain from gin: a
// RouteInfo carries the last handler and nothing in front of it. Grit already
// solves this for its own Routes screen by parsing the route files and
// generating internal/access, so the table exists and is kept fresh by
// `grit generate`, `grit remove`, `grit upgrade` and checked by `grit doctor`.
// This hands that same table to Scout.
func AddScout(root string, force bool) error {
	apiDir := filepath.Join(root, "apps", "api")
	if !dirExists(apiDir) {
		apiDir = root // an --api-only project is its own api directory
	}
	if !fileExists(filepath.Join(apiDir, "go.mod")) {
		return fmt.Errorf("no Go API found: grit add scout needs apps/api (or an --api project)")
	}
	if !fileExists(filepath.Join(apiDir, "internal", "access", "registry.go")) {
		return fmt.Errorf(
			"this project has no access table, which Scout needs to see which routes are protected\n\n" +
				"Run `grit upgrade` first: it generates internal/access from your route files.")
	}

	dim := color.New(color.FgHiBlack)
	green := color.New(color.FgHiGreen)
	cyan := color.New(color.FgHiCyan)

	modulePath, err := scoutModulePath(apiDir)
	if err != nil {
		return err
	}

	target := filepath.Join(apiDir, "internal", "routes", "scout.go")
	if fileExists(target) && !force {
		dim.Printf("  → internal/routes/scout.go exists, leaving it alone (--force to replace)\n")
	} else {
		dim.Printf("  → Writing internal/routes/scout.go...\n")
		if err := os.WriteFile(target, []byte(scoutRoutesGo(modulePath)), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", target, err)
		}
		green.Printf("  ✓ internal/routes/scout.go\n")
	}

	var wired []string
	routesFile := filepath.Join(apiDir, "internal", "routes", "routes.go")
	if err := injectBeforeLine(routesFile, "mountPulse(r, db, cfg, svc)",
		"\tmountScout(r, cfg)\n", &wired, "routes.go"); err != nil {
		return err
	}
	if len(wired) == 0 && !strings.Contains(readFileOrEmpty(routesFile), "mountScout(") {
		return fmt.Errorf("could not find the mount point in internal/routes/routes.go\n\n" +
			"Add this line to Setup() by hand:\n    mountScout(r, cfg)")
	}

	dim.Printf("  → Adding config fields...\n")
	envFile, err := addScoutConfig(root, apiDir)
	if err != nil {
		return err
	}
	green.Printf("  ✓ internal/config/config.go, %s\n", envFile)

	dim.Printf("  → go get %s@%s...\n", scoutModule, scoutVersion)
	if err := scoutGoGet(apiDir); err != nil {
		return err
	}
	green.Printf("  ✓ go.mod (gin %s, gorm %s)\n", "v1.12.0", "v1.31.2")

	fmt.Println()
	cyan.Printf("  Scout is mounted at /scout.\n\n")
	fmt.Printf("  Set a password before you start it:\n")
	dim.Printf("      SCOUT_PASSWORD=...        in .env, beside grit.json\n\n")
	fmt.Printf("  It runs in safe mode: headers, TLS, DNS, exposed files and config.\n")
	fmt.Printf("  Nothing it sends by default is an attack payload. Raising the mode\n")
	fmt.Printf("  sends real ones, so raise it only against a target you own.\n\n")
	fmt.Printf("  Docs: ")
	cyan.Printf("https://gritframework.dev/docs/security/scout\n\n")
	return nil
}

const scoutModule = "github.com/MUKE-coder/scout"

// scoutGoGet adds the module and the two versions it raises.
func scoutGoGet(apiDir string) error {
	args := [][]string{{"get", scoutModule + "@" + scoutVersion}}
	for _, r := range scoutRequires {
		args = append(args, []string{"get", r.module + "@" + r.version})
	}
	args = append(args, []string{"mod", "tidy"})
	for _, a := range args {
		cmd := exec.Command("go", a...)
		cmd.Dir = apiDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go %s: %w\n%s", strings.Join(a, " "), err, out)
		}
	}
	return nil
}

// addScoutConfig adds the four settings to the config struct, its loader and
// the project's .env.example. Each edit is skipped when it is already there,
// so running the command twice changes nothing. It returns the env file it
// wrote, for the line the command prints.
func addScoutConfig(root, apiDir string) (string, error) {
	cfgPath := filepath.Join(apiDir, "internal", "config", "config.go")
	content := readFileOrEmpty(cfgPath)
	if content == "" {
		return "", fmt.Errorf("reading %s", cfgPath)
	}

	if !strings.Contains(content, "ScoutEnabled") {
		const anchor = "\t// Observability (Pulse"
		fields := "\t// Security scanning (Scout)\n" +
			"\t//\n" +
			"\t// Off unless asked for. A scanner that starts itself is a scanner\n" +
			"\t// nobody decided to run.\n" +
			"\tScoutEnabled  bool\n" +
			"\tScoutUsername string\n" +
			"\tScoutPassword string\n" +
			"\t// Where findings are kept. They outlive a restart on purpose:\n" +
			"\t// an accepted risk has an expiry, and a scanner that forgets what\n" +
			"\t// you triaged re-reports it every run until you stop reading it.\n" +
			"\tScoutDBPath string\n\n"
		if !strings.Contains(content, anchor) {
			return "", fmt.Errorf("could not find the config struct's Pulse block in %s", cfgPath)
		}
		content = strings.Replace(content, anchor, fields+anchor, 1)
	}

	if !strings.Contains(content, "SCOUT_ENABLED") {
		const anchor = "\t\tPulseEnabled:"
		load := "\t\tScoutEnabled:  getEnv(\"SCOUT_ENABLED\", \"false\") == \"true\",\n" +
			"\t\tScoutUsername: getEnv(\"SCOUT_USERNAME\", \"admin\"),\n" +
			"\t\tScoutPassword: getEnv(\"SCOUT_PASSWORD\", \"\"),\n" +
			"\t\tScoutDBPath:   resolveSQLitePath(getEnv(\"SCOUT_DB\", \"scout.db\")),\n\n"
		if !strings.Contains(content, anchor) {
			return "", fmt.Errorf("could not find the config loader's Pulse block in %s", cfgPath)
		}
		content = strings.Replace(content, anchor, load+anchor, 1)
	}

	if err := writeEdit(cfgPath, content); err != nil {
		return "", fmt.Errorf("writing %s: %w", cfgPath, err)
	}

	// The project root, not apps/api. That is where `grit new` writes
	// .env.example, where the developer's .env sits, and where `grit doctor`
	// reads its settings from; the API only ever finds one by walking up
	// (".env", "../.env", "../../.env"). Writing to apps/api/.env.example
	// instead meant writing to a file that does not exist, which this
	// function did silently while the command printed a tick beside it.
	envPath := filepath.Join(root, ".env.example")
	env := readFileOrEmpty(envPath)
	if env == "" {
		return "", fmt.Errorf("no .env.example at %s: Scout's settings have nowhere to go", envPath)
	}
	label := ".env.example"
	if !strings.Contains(env, "SCOUT_ENABLED") {
		block := "\n# Scout security scanning (grit add scout)\n" +
			"# Off by default. SCOUT_PASSWORD has no default: Scout refuses to\n" +
			"# serve its dashboard without one rather than picking a weak one.\n" +
			"SCOUT_ENABLED=false\n" +
			"SCOUT_USERNAME=admin\n" +
			"SCOUT_PASSWORD=\n" +
			"SCOUT_DB=scout.db\n"
		if err := writeEdit(envPath, strings.TrimRight(env, "\n")+"\n"+block); err != nil {
			return "", fmt.Errorf("writing %s: %w", envPath, err)
		}
	}
	return label, nil
}

// scoutModulePath reads the API module's import path out of its go.mod, which
// is what the generated file's own imports are written against.
func scoutModulePath(apiDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(apiDir, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if module, ok := strings.CutPrefix(line, "module "); ok {
			if module = strings.TrimSpace(module); module != "" {
				return module, nil
			}
		}
	}
	return "", fmt.Errorf("no module line in %s/go.mod", apiDir)
}

func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
