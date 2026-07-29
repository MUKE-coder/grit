package testrunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func find(suites []Suite, label string) *Suite {
	for i := range suites {
		if suites[i].Label == label {
			return &suites[i]
		}
	}
	return nil
}

func byKey(suites []Suite, key string) []Suite {
	var out []Suite
	for _, s := range suites {
		if s.Key == key {
			out = append(out, s)
		}
	}
	return out
}

// A monorepo: go.mod under apps/api, a root orchestrator script, e2e available.
func monorepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write(t, filepath.Join(root, "grit.json"), `{"architecture":"triple"}`)
	write(t, filepath.Join(root, "pnpm-lock.yaml"), "")
	write(t, filepath.Join(root, "apps", "api", "go.mod"), "module app/apps/api\n")
	write(t, filepath.Join(root, "package.json"),
		`{"scripts":{"test":"turbo run test","test:e2e":"playwright test"}}`)
	return root
}

// An api-only project: Go at the root, no JavaScript at all.
func apiOnly(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write(t, filepath.Join(root, "grit.json"), `{"architecture":"api"}`)
	write(t, filepath.Join(root, "go.mod"), "module app\n")
	return root
}

func TestDiscoverMonorepoFindsGoAndWorkspace(t *testing.T) {
	root := monorepo(t)
	suites := Discover(root, Options{})

	goSuite := byKey(suites, "go")
	if len(goSuite) != 1 || goSuite[0].SkipReason != "" {
		t.Fatalf("Go suite not discovered: %+v", goSuite)
	}
	if want := filepath.Join(root, "apps", "api"); goSuite[0].Dir != want {
		t.Errorf("Go dir = %q, want %q", goSuite[0].Dir, want)
	}

	// The root script is the project's own definition of "run the tests"; using
	// it avoids running every app twice via turbo AND directly.
	node := byKey(suites, "node")
	if len(node) != 1 {
		t.Fatalf("expected exactly 1 node suite (the root orchestrator), got %d: %+v", len(node), node)
	}
	if node[0].Dir != root || node[0].Cmd != "pnpm" {
		t.Errorf("node suite = %+v, want the root dir run with pnpm", node[0])
	}
}

// The most important property in the package: a suite that cannot run must say
// why. A runner that silently runs nothing looks exactly like one that passed.
func TestSuitesThatCannotRunCarryAReason(t *testing.T) {
	suites := Discover(apiOnly(t), Options{})

	for _, s := range suites {
		if s.Cmd == "" && s.SkipReason == "" {
			t.Errorf("suite %q is neither runnable nor explained", s.Label)
		}
	}

	node := byKey(suites, "node")
	if len(node) != 1 || node[0].SkipReason == "" {
		t.Fatalf("api-only project should report a skipped JS suite with a reason, got %+v", node)
	}
	if !strings.Contains(node[0].SkipReason, "apps/") {
		t.Errorf("skip reason should say what was looked for, got %q", node[0].SkipReason)
	}
}

func TestE2EIsOptInAndExplainsItself(t *testing.T) {
	root := monorepo(t)

	off := find(Discover(root, Options{}), "End-to-end")
	if off == nil || off.SkipReason == "" {
		t.Fatal("e2e should be skipped by default")
	}
	if !strings.Contains(off.SkipReason, "--e2e") {
		t.Errorf("skip reason should tell you how to enable it, got %q", off.SkipReason)
	}

	on := find(Discover(root, Options{E2E: true}), "End-to-end")
	if on == nil || on.SkipReason != "" {
		t.Fatalf("e2e should run when opted in, got %+v", on)
	}
	if on.Args[0] != "test:e2e" {
		t.Errorf("e2e args = %v, want the test:e2e script", on.Args)
	}
}

// Opting into e2e on an architecture that has no Playwright suite must not
// pretend it ran.
func TestE2EOnAProjectWithoutItStillSkips(t *testing.T) {
	s := find(Discover(apiOnly(t), Options{E2E: true}), "End-to-end")
	if s == nil || s.SkipReason == "" {
		t.Fatalf("expected a skip, got %+v", s)
	}
	if !strings.Contains(s.SkipReason, "test:e2e") {
		t.Errorf("reason should name the missing script, got %q", s.SkipReason)
	}
}

// Without a root orchestrator, each app that defines a test script is its own
// suite — and apps that define none are not invented.
func TestPerAppDiscoveryWhenThereIsNoRootScript(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "grit.json"), `{"architecture":"triple"}`)
	write(t, filepath.Join(root, "apps", "api", "go.mod"), "module app/apps/api\n")
	write(t, filepath.Join(root, "apps", "web", "package.json"), `{"scripts":{"test":"vitest run"}}`)
	write(t, filepath.Join(root, "apps", "admin", "package.json"), `{"scripts":{"test":"vitest run"}}`)
	// A docs app with no test script must not become a suite.
	write(t, filepath.Join(root, "apps", "docs", "package.json"), `{"scripts":{"build":"next build"}}`)

	node := byKey(Discover(root, Options{}), "node")
	if len(node) != 2 {
		t.Fatalf("expected web + admin, got %d: %+v", len(node), node)
	}
	for _, s := range node {
		if s.Label != "web" && s.Label != "admin" {
			t.Errorf("unexpected suite %q", s.Label)
		}
	}
}

func TestGoFlagsAreApplied(t *testing.T) {
	suites := Discover(monorepo(t), Options{Race: true, Cover: true})
	args := strings.Join(byKey(suites, "go")[0].Args, " ")

	for _, want := range []string{"test", "./...", "-race", "-cover"} {
		if !strings.Contains(args, want) {
			t.Errorf("go args %q missing %q", args, want)
		}
	}
}

func TestPackageManagerFollowsTheLockfile(t *testing.T) {
	for _, tc := range []struct{ lock, want string }{
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
		{"package-lock.json", "npm"},
		{"bun.lockb", "bun"},
	} {
		root := t.TempDir()
		write(t, filepath.Join(root, tc.lock), "")
		if got := packageManager(root); got != tc.want {
			t.Errorf("%s → %q, want %q", tc.lock, got, tc.want)
		}
	}

	// No lockfile: Grit scaffolds pnpm, so that is the honest default.
	if got := packageManager(t.TempDir()); got != "pnpm" {
		t.Errorf("no lockfile → %q, want pnpm", got)
	}
}

func TestFilterSelectsByKeyAndEmptyKeepsEverything(t *testing.T) {
	suites := Discover(monorepo(t), Options{E2E: true})

	if got := len(Filter(suites, nil)); got != len(suites) {
		t.Errorf("empty filter dropped suites: %d → %d", len(suites), got)
	}

	only := Filter(suites, []string{"go"})
	if len(only) != 1 || only[0].Key != "go" {
		t.Errorf("filter go = %+v", only)
	}
}

func TestFailedReflectsAnyFailure(t *testing.T) {
	if Failed([]Result{{Status: StatusPass}, {Status: StatusSkip}}) {
		t.Error("pass + skip should not be a failure")
	}
	if !Failed([]Result{{Status: StatusPass}, {Status: StatusFail}}) {
		t.Error("a failure must be reported")
	}
	// A run where everything skipped is not a failure, but it is also not
	// evidence of anything — the report says so per row.
	if Failed([]Result{{Status: StatusSkip}}) {
		t.Error("all-skipped should not be a failure")
	}
}

// Skipped suites must never be executed.
func TestRunDoesNotExecuteSkippedSuites(t *testing.T) {
	suites := []Suite{
		{Key: "go", Label: "Skipped", SkipReason: "not applicable", Cmd: "definitely-not-a-real-binary"},
	}
	var out strings.Builder
	results := Run(suites, &out, &out)

	if len(results) != 1 || results[0].Status != StatusSkip {
		t.Fatalf("expected a single skip, got %+v", results)
	}
	if out.Len() != 0 {
		t.Errorf("a skipped suite produced output: %q", out.String())
	}
}
