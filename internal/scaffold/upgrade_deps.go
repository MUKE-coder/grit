package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// frameworkDeps are the libraries every generated API mounts, at the lowest
// version a project should run. grit upgrade raises a project below one and
// leaves one at or above it alone: it never downgrades.
//
// Nothing moved them before. grit upgrade rewrites framework files, not
// go.mod, so a project scaffolded on Sentinel v2.2.1 stayed there after the
// scaffold moved on, and v2.2.2 is a security release.
// FrameworkDep is one library every generated API mounts, with the lowest
// version a project should run and why that version.
type FrameworkDep struct {
	Path, Floor, Why string
}

// FrameworkDepFloors returns those libraries, for a caller that reports on
// them rather than changes them: grit doctor names a project left behind.
func FrameworkDepFloors() []FrameworkDep { return frameworkDeps }

var frameworkDeps = []FrameworkDep{
	{"github.com/MUKE-coder/sentinel/v2", "v2.5.0",
		"security fixes from v2.2.2 (client-IP spoofing behind a proxy, a sort_by SQL injection, SSRF bypasses)"},
	{"github.com/MUKE-coder/gorm-studio", "v1.1.0",
		"read-only SQL enforced on the read path, fail-closed imports, composite keys matched in full"},
	{"github.com/MUKE-coder/pulse", "v1.0.0", "the tagged release"},
	{"golang.org/x/crypto", "v0.57.0", "govulncheck flags releases before v0.56.0"},
	{"filippo.io/edwards25519", "v1.2.0", "govulncheck flags releases before v1.1.1"},
	{"github.com/gorilla/mux", "v1.8.1", "goth's gothic links it in and still asks for v1.6.2, from 2018"},
}

// frameworkRequires are modules a framework-owned file imports, which an
// upgrade therefore has to put in go.mod.
//
// Separate from frameworkDeps because the two answer different questions. A
// floor says "if you have this, do not run it below this version", and adding
// one a project never had would pull in a library it does not use. A
// requirement says "a file this upgrade just wrote imports this", and leaving
// it out produces a project that does not compile.
//
// NeededBy is the file, relative to the API root, whose presence means the
// module is wanted. Checked rather than assumed, so a shape that does not get
// that file does not get the dependency either.
var frameworkRequires = []struct {
	Path, Version, NeededBy, Why string
}{
	{"go.opentelemetry.io/otel", "v1.45.0", "internal/tracing/tracing.go",
		"internal/tracing imports it; inert until OTEL_EXPORTER_OTLP_ENDPOINT is set"},
	{"go.opentelemetry.io/otel/sdk", "v1.45.0", "internal/tracing/tracing.go", "the tracer provider"},
	{"go.opentelemetry.io/otel/trace", "v1.45.0", "internal/tracing/tracing.go", "the span API"},
	{"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp", "v1.45.0",
		"internal/tracing/tracing.go", "the OTLP/HTTP exporter"},
}

// addFrameworkRequires puts a module in go.mod when the framework file that
// imports it is present and the module is not.
func addFrameworkRequires(apiRoot string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(apiRoot, "go.mod"))
	if err != nil {
		return nil, nil // no go.mod, so no API to add to
	}

	var specs, added []string
	for _, r := range frameworkRequires {
		if _, err := os.Stat(filepath.Join(apiRoot, r.NeededBy)); err != nil {
			continue // this project does not have the file that needs it
		}
		if regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(r.Path) + `\s+v`).Match(data) {
			continue // already required, and raiseFrameworkDeps owns the version
		}
		specs = append(specs, r.Path+"@"+r.Version)
		added = append(added, fmt.Sprintf("%s %s: %s", r.Path, r.Version, r.Why))
	}
	if len(specs) == 0 {
		return nil, nil
	}
	if err := goGet(apiRoot, specs...); err != nil {
		return nil, fmt.Errorf("could not add %s, so run `go get %s` in %s: %v",
			strings.Join(specs, ", "), strings.Join(specs, " "), apiRoot, err)
	}
	return added, nil
}

// goGet runs go get in dir. A variable, so a test can see what would be fetched
// without the network.
var goGet = func(dir string, specs ...string) error {
	cmd := exec.Command("go", append([]string{"get"}, specs...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// raiseFrameworkDeps raises every framework library in apiRoot's go.mod that
// is below its floor, and describes each one it raised.
func raiseFrameworkDeps(apiRoot string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(apiRoot, "go.mod"))
	if err != nil {
		return nil, nil // no go.mod, so no API to raise
	}
	var specs, raised []string
	for _, d := range frameworkDeps {
		m := regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(d.Path) + `\s+(v\S+)`).FindSubmatch(data)
		if m == nil {
			continue // not a dependency of this project
		}
		current := string(m[1])
		if !versionLess(current, d.Floor) {
			continue
		}
		specs = append(specs, d.Path+"@"+d.Floor)
		raised = append(raised, fmt.Sprintf("%s %s → %s: %s", d.Path, current, d.Floor, d.Why))
	}
	if len(specs) == 0 {
		return nil, nil
	}
	if err := goGet(apiRoot, specs...); err != nil {
		return nil, fmt.Errorf("could not raise %s, so run `go get %s` in %s: %v",
			strings.Join(specs, ", "), strings.Join(specs, " "), apiRoot, err)
	}
	return raised, nil
}

// VersionBelow reports whether current is lower than floor, by the same rule
// grit upgrade raises a dependency by. Exported for grit doctor, which reports
// on the same floors without changing anything.
func VersionBelow(current, floor string) bool { return versionLess(current, floor) }

// versionLess reports whether semantic version a is lower than b. As much of
// semver as go.mod needs: numbers compare as numbers, and a pre-release sorts
// below its release, which puts a pseudo-version such as
// v0.0.0-20260529025319-478cdfa8ce5f below every tagged release.
func versionLess(a, b string) bool {
	an, ap := splitVersion(a)
	bn, bp := splitVersion(b)
	for i := range an {
		if an[i] != bn[i] {
			return an[i] < bn[i]
		}
	}
	switch {
	case ap != "" && bp == "":
		return true
	case ap == "" && bp != "":
		return false
	}
	return ap < bp
}

// splitVersion returns the three numbers of a version and its pre-release, if
// any. Build metadata such as +incompatible is dropped.
func splitVersion(v string) ([3]int, string) {
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexByte(v, '+'); i >= 0 {
		v = v[:i]
	}
	pre := ""
	if i := strings.IndexByte(v, '-'); i >= 0 {
		v, pre = v[:i], v[i+1:]
	}
	var n [3]int
	for i, part := range strings.SplitN(v, ".", 3) {
		n[i], _ = strconv.Atoi(part)
	}
	return n, pre
}
