package scaffold

import (
	"encoding/json"
	"strings"
	"testing"
)

// The root `type-check` script runs `turbo type-check`, and turbo matches
// tasks by name. Four app templates declared `typecheck` without the hyphen,
// so turbo found nothing to run in them, reported "No tasks were executed" and
// exited 0.
//
// That is the worst shape a check can take. It is not a missing script, which
// somebody notices: it is a command that looks like it ran, says nothing, and
// is green. A project generated this way could ship a type error through a CI
// step whose whole purpose was to catch it.
//
// Found by running `pnpm run type-check` on a freshly generated triple-tier
// project, which is the only way it could have been found: every unit test in
// this package passed throughout.

// scriptsOf pulls the scripts map out of a rendered package.json.
func scriptsOf(t *testing.T, name, body string) map[string]string {
	t.Helper()
	var parsed struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("%s is not valid JSON: %v", name, err)
	}
	return parsed.Scripts
}

// taskNamesOf pulls the task names out of turbo.json.
func taskNamesOf(t *testing.T, body string) map[string]bool {
	t.Helper()
	var parsed struct {
		Tasks    map[string]interface{} `json:"tasks"`
		Pipeline map[string]interface{} `json:"pipeline"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("turbo.json is not valid JSON: %v", err)
	}
	names := map[string]bool{}
	for name := range parsed.Tasks {
		names[name] = true
	}
	for name := range parsed.Pipeline {
		names[name] = true
	}
	return names
}

// Every app declares the type-check task turbo is told to run.
func TestEveryAppDeclaresTheTypeCheckTask(t *testing.T) {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}
	tasks := taskNamesOf(t, turboJSON())

	const task = "type-check"
	if !tasks[task] {
		t.Fatalf("turbo.json declares no %q task; the root script runs `turbo %s`", task, task)
	}

	apps := map[string]string{
		"apps/web/package.json":        webPackageJSON(opts),
		"apps/admin/package.json":      adminPackageJSON(opts),
		"apps/web (tanstack)":          webTanStackPackageJSON(Options{ProjectName: "app", Frontend: FrontendTanStack}),
		"apps/admin (tanstack)":        adminTanStackPackageJSON(Options{ProjectName: "app", Frontend: FrontendTanStack}),
		"apps/desktop/frontend":        desktopClientPackageJSON(opts),
		"packages/upload/package.json": uploadPackageJSON(),
	}

	for name, body := range apps {
		scripts := scriptsOf(t, name, body)
		if _, ok := scripts[task]; !ok {
			t.Errorf("%s has no %q script, so `turbo %s` skips it silently and still exits 0 "+
				"(it has: %s)", name, task, task, strings.Join(keysOf(scripts), ", "))
		}
	}
}

// And nothing declares the unhyphenated spelling, which is the one that went
// unrun for as long as it existed.
func TestNoAppUsesTheUnhyphenatedSpelling(t *testing.T) {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}
	viteOpts := Options{ProjectName: "app", Frontend: FrontendTanStack}

	for name, body := range map[string]string{
		"apps/web/package.json":        webPackageJSON(opts),
		"apps/admin/package.json":      adminPackageJSON(opts),
		"apps/web (tanstack)":          webTanStackPackageJSON(viteOpts),
		"apps/admin (tanstack)":        adminTanStackPackageJSON(viteOpts),
		"apps/desktop/frontend":        desktopClientPackageJSON(opts),
		"packages/upload/package.json": uploadPackageJSON(),
		"turbo.json":                   turboJSON(),
		"package.json":                 rootPackageJSON(opts),
	} {
		if strings.Contains(body, `"typecheck"`) {
			t.Errorf("%s declares \"typecheck\"; the task turbo runs is \"type-check\", so "+
				"this one is never executed", name)
		}
	}
}

// The root script has to point at the task that exists.
func TestTheRootScriptRunsTheRealTask(t *testing.T) {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}
	scripts := scriptsOf(t, "package.json", rootPackageJSON(opts))

	script, ok := scripts["type-check"]
	if !ok {
		t.Fatal("the root package.json has no type-check script")
	}
	tasks := taskNamesOf(t, turboJSON())
	// "turbo type-check" or "turbo run type-check"
	fields := strings.Fields(script)
	named := fields[len(fields)-1]
	if !tasks[named] {
		t.Fatalf("the root script runs %q and turbo.json declares no such task", named)
	}
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
