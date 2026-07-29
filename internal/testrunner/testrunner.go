// Package testrunner discovers and runs every test suite in a Grit project.
//
// A Grit project has its tests split across languages and directories — Go in
// the API, Vitest in each frontend app, Playwright at the root — and which of
// those exist depends on the architecture the project was scaffolded with. The
// point of this package is that you should not have to remember which.
//
// Suites that do not apply are reported as SKIPPED WITH A REASON rather than
// omitted. A runner that silently runs nothing looks identical to a runner that
// passed, which is the most expensive kind of green.
package testrunner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Status is the outcome of one suite.
type Status string

const (
	StatusPass Status = "PASS"
	StatusFail Status = "FAIL"
	StatusSkip Status = "SKIP"
)

// Suite is one runnable (or deliberately skipped) group of tests.
type Suite struct {
	Key   string // stable id used by the selection flags: go, node, e2e
	Label string // human label in the report
	Dir   string // working directory
	Cmd   string
	Args  []string

	// SkipReason, when set, means the suite cannot run here. It is always
	// reported, never hidden.
	SkipReason string
}

// Result pairs a suite with what happened when it ran.
type Result struct {
	Suite
	Status   Status
	Duration time.Duration
	Err      error
}

// Options controls discovery.
type Options struct {
	// Race and Cover apply to the Go suite only.
	Race  bool
	Cover bool
	// E2E opts the Playwright suite in. It is off by default because it needs
	// the API and the frontends already running; failing against a server that
	// was never started tells you nothing about your code.
	E2E bool
}

// Discover works out which suites apply to the project rooted at root.
func Discover(root string, opts Options) []Suite {
	var suites []Suite

	suites = append(suites, discoverGo(root, opts))
	suites = append(suites, discoverNode(root)...)
	suites = append(suites, discoverE2E(root, opts))

	return suites
}

func discoverGo(root string, opts Options) Suite {
	s := Suite{Key: "go", Label: "Go", Cmd: "go"}

	dir := goModuleDir(root)
	if dir == "" {
		s.SkipReason = "no go.mod found at the project root or apps/api"
		return s
	}

	s.Dir = dir
	s.Args = []string{"test", "./..."}
	if opts.Race {
		s.Args = append(s.Args, "-race")
	}
	if opts.Cover {
		s.Args = append(s.Args, "-cover")
	}
	return s
}

// discoverNode returns the JavaScript suites.
//
// When the root package.json defines a "test" script, that is the project's own
// answer to "run the tests" — usually turbo fanning out across the workspace —
// so it is used as-is rather than second-guessed by running each app
// separately. Without a root orchestrator, each app under apps/ that defines a
// test script becomes its own suite.
func discoverNode(root string) []Suite {
	pm := packageManager(root)

	if hasScript(filepath.Join(root, "package.json"), "test") {
		return []Suite{{
			Key:   "node",
			Label: "Workspace (" + pm + ")",
			Dir:   root,
			Cmd:   pm,
			Args:  []string{"test"},
		}}
	}

	appsDir := filepath.Join(root, "apps")
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return []Suite{{
			Key:        "node",
			Label:      "JavaScript",
			SkipReason: "no package.json with a test script at the root, and no apps/ directory",
		}}
	}

	var suites []Suite
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		appDir := filepath.Join(appsDir, e.Name())
		if !hasScript(filepath.Join(appDir, "package.json"), "test") {
			continue
		}
		suites = append(suites, Suite{
			Key:   "node",
			Label: e.Name(),
			Dir:   appDir,
			Cmd:   pm,
			Args:  []string{"test"},
		})
	}

	if len(suites) == 0 {
		return []Suite{{
			Key:        "node",
			Label:      "JavaScript",
			SkipReason: "no app under apps/ defines a test script",
		}}
	}
	return suites
}

func discoverE2E(root string, opts Options) Suite {
	s := Suite{Key: "e2e", Label: "End-to-end"}

	if !hasScript(filepath.Join(root, "package.json"), "test:e2e") {
		s.SkipReason = "no test:e2e script (this architecture has no Playwright suite)"
		return s
	}
	if !opts.E2E {
		s.SkipReason = "not requested — pass --e2e (needs the API and frontends running)"
		return s
	}

	s.Dir = root
	s.Cmd = packageManager(root)
	s.Args = []string{"test:e2e"}
	return s
}

// Run executes the suites in order, streaming their output to out, and returns
// what happened to each.
func Run(suites []Suite, out io.Writer, errOut io.Writer) []Result {
	results := make([]Result, 0, len(suites))

	for _, s := range suites {
		if s.SkipReason != "" {
			results = append(results, Result{Suite: s, Status: StatusSkip})
			continue
		}

		fmt.Fprintf(out, "\n  ── %s ──────────────────────────────\n\n", s.Label)

		start := time.Now()
		c := exec.Command(s.Cmd, s.Args...)
		c.Dir = s.Dir
		c.Stdout = out
		c.Stderr = errOut
		err := c.Run()
		elapsed := time.Since(start)

		status := StatusPass
		if err != nil {
			status = StatusFail
		}
		results = append(results, Result{Suite: s, Status: status, Duration: elapsed, Err: err})
	}

	return results
}

// Failed reports whether any suite failed, which is what the command's exit
// code is built from.
func Failed(results []Result) bool {
	for _, r := range results {
		if r.Status == StatusFail {
			return true
		}
	}
	return false
}

// Filter keeps only the suites whose Key is in keys. An empty keys slice keeps
// everything.
func Filter(suites []Suite, keys []string) []Suite {
	if len(keys) == 0 {
		return suites
	}
	want := make(map[string]bool, len(keys))
	for _, k := range keys {
		want[k] = true
	}

	out := make([]Suite, 0, len(suites))
	for _, s := range suites {
		if want[s.Key] {
			out = append(out, s)
		}
	}
	return out
}

// goModuleDir returns the directory holding the API's go.mod — apps/api in a
// monorepo, the project root otherwise — or "" when there is none.
func goModuleDir(root string) string {
	for _, candidate := range []string{filepath.Join(root, "apps", "api"), root} {
		if _, err := os.Stat(filepath.Join(candidate, "go.mod")); err == nil {
			return candidate
		}
	}
	return ""
}

// packageManager picks the tool whose lockfile is present, defaulting to pnpm,
// which is what Grit scaffolds.
func packageManager(root string) string {
	for _, c := range []struct{ lock, cmd string }{
		{"pnpm-lock.yaml", "pnpm"},
		{"bun.lockb", "bun"},
		{"yarn.lock", "yarn"},
		{"package-lock.json", "npm"},
	} {
		if _, err := os.Stat(filepath.Join(root, c.lock)); err == nil {
			return c.cmd
		}
	}
	return "pnpm"
}

// hasScript reports whether a package.json defines a non-empty script.
func hasScript(pkgPath, name string) bool {
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	return pkg.Scripts[name] != ""
}
