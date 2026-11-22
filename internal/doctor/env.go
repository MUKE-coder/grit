package doctor

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/clierr"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The environment check: is this machine ready to run a Grit project?
//
// Separate from the project checks in doctor.go, and deliberately so. Those
// read a project's files and answer "is this code right". This one runs before
// there is a project at all, and answers "will grit new work here". It is the
// check a beginner needs most and the only one they can run before they have
// anything to run it against.
//
// It shells out to go, node, pnpm and docker, and it opens local sockets to see
// whether a port is free. It still touches no network, which is the property
// that makes `grit doctor` usable in CI, and it is worth keeping: a check that
// reaches the internet fails on a build machine for reasons that have nothing
// to do with the thing being checked.

// EnvLevel is how much a finding matters.
type EnvLevel string

const (
	// EnvOK is a tool that is present and new enough.
	EnvOK EnvLevel = "ok"
	// EnvWarn is something that will bite later, or only for some projects.
	EnvWarn EnvLevel = "warning"
	// EnvFail is something that stops grit new or grit start working.
	EnvFail EnvLevel = "error"
)

// EnvFinding is one line of the environment report.
type EnvFinding struct {
	Level EnvLevel
	// Name is what was checked: "Go", "Docker", "port 5434".
	Name string
	// Detail is what was found: a version, or why it could not be.
	Detail string
	// Problem is the catalogue entry to read, when there is one. Empty on a
	// finding that is fine.
	Problem *clierr.Problem
}

// EnvReport is the whole environment check.
type EnvReport struct {
	Findings []EnvFinding
}

// Failures counts the findings that stop something working.
func (r EnvReport) Failures() int { return r.count(EnvFail) }

// Warnings counts the findings worth reading.
func (r EnvReport) Warnings() int { return r.count(EnvWarn) }

func (r EnvReport) count(l EnvLevel) int {
	n := 0
	for _, f := range r.Findings {
		if f.Level == l {
			n++
		}
	}
	return n
}

// EnvOptions narrows what is checked.
type EnvOptions struct {
	// SkipNode leaves out the frontend toolchain, for an API-only project.
	SkipNode bool
	// SkipDocker leaves out Docker, for a project on SQLite and local
	// storage, which needs no containers at all.
	SkipDocker bool
	// Ports to test. Empty means the defaults a generated project uses.
	Ports []PortCheck
}

// PortCheck is one port and what wants it.
type PortCheck struct {
	Port int
	What string
}

// DefaultPorts are the ports a generated project listens on, with the
// non-obvious choices explained where they are made: Postgres is on 5434 and
// Redis on 6380 so that a native install of either does not collide.
var DefaultPorts = []PortCheck{
	{8080, "the Go API"},
	{3000, "the web app"},
	{3001, "the admin"},
	{5434, "Postgres, in Docker"},
	{6380, "Redis, in Docker"},
	{9002, "MinIO, in Docker"},
	{8025, "the Mailhog inbox"},
}

// CheckEnv runs the environment checks and returns what it found.
func CheckEnv(ctx context.Context, opts EnvOptions) EnvReport {
	var r EnvReport

	r.Findings = append(r.Findings, checkGo(ctx))
	if !opts.SkipNode {
		r.Findings = append(r.Findings, checkNode(ctx), checkPnpm(ctx))
	}
	if !opts.SkipDocker {
		r.Findings = append(r.Findings, checkDocker(ctx))
	}

	ports := opts.Ports
	if ports == nil {
		ports = DefaultPorts
	}
	for _, p := range ports {
		r.Findings = append(r.Findings, checkPort(p))
	}
	return r
}

// run executes a command with a timeout and returns its combined output.
//
// The timeout matters: `docker info` against a daemon that is starting can
// hang for a long time, and a preflight check that hangs is worse than one
// that fails.
func runTool(ctx context.Context, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func problem(slug string) *clierr.Problem {
	if p, ok := clierr.BySlug(slug); ok {
		return &p
	}
	return nil
}

var envGoVersionRe = regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)

func checkGo(ctx context.Context) EnvFinding {
	out, err := runTool(ctx, 10*time.Second, "go", "version")
	if err != nil {
		return EnvFinding{EnvFail, "Go", "not found on PATH", problem("go-not-found")}
	}
	m := envGoVersionRe.FindStringSubmatch(out)
	if m == nil {
		// It ran and said something unexpected. Report what it said rather
		// than guessing: a wrong version comparison is worse than none.
		return EnvFinding{EnvWarn, "Go", "could not read the version from: " + out, nil}
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch := 0
	if m[3] != "" {
		patch, _ = strconv.Atoi(m[3])
	}
	found := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	wantMajor, wantMinor, wantPatch := parseVersion(scaffold.GoToolchainPin)
	if older(major, minor, patch, wantMajor, wantMinor, wantPatch) {
		return EnvFinding{
			EnvFail, "Go",
			fmt.Sprintf("%s, and a generated project asks for %s", found, scaffold.GoToolchainPin),
			problem("go-too-old"),
		}
	}
	return EnvFinding{EnvOK, "Go", found, nil}
}

var nodeVersionRe = regexp.MustCompile(`v?(\d+)\.(\d+)\.(\d+)`)

// nodeFloorMajor is Next's own engines floor, as a major version.
//
// Not read from a package.json here: this runs before a project exists, so
// there is nothing to read. docs/config/site.ts states the same number for
// readers, and docs/scripts/check-toolchain-versions.py checks that against
// Next's installed engines field.
const nodeFloorMajor = 20

func checkNode(ctx context.Context) EnvFinding {
	out, err := runTool(ctx, 10*time.Second, "node", "--version")
	if err != nil {
		return EnvFinding{EnvFail, "Node", "not found on PATH", problem("node-not-found")}
	}
	m := nodeVersionRe.FindStringSubmatch(out)
	if m == nil {
		return EnvFinding{EnvWarn, "Node", "could not read the version from: " + out, nil}
	}
	major, _ := strconv.Atoi(m[1])
	found := strings.TrimPrefix(out, "v")
	if major < nodeFloorMajor {
		return EnvFinding{
			EnvFail, "Node",
			fmt.Sprintf("%s, and the frontend needs %d or newer", found, nodeFloorMajor),
			problem("node-too-old"),
		}
	}
	return EnvFinding{EnvOK, "Node", found, nil}
}

func checkPnpm(ctx context.Context) EnvFinding {
	out, err := runTool(ctx, 10*time.Second, "pnpm", "--version")
	if err != nil {
		// A warning rather than a failure: corepack can fetch it on first
		// use, so this is a thing to know rather than a thing that stops you.
		return EnvFinding{EnvWarn, "pnpm", "not found on PATH", problem("pnpm-not-found")}
	}
	return EnvFinding{EnvOK, "pnpm", out, nil}
}

func checkDocker(ctx context.Context) EnvFinding {
	if _, err := exec.LookPath("docker"); err != nil {
		return EnvFinding{
			EnvWarn, "Docker", "not found on PATH",
			problem("docker-not-running"),
		}
	}
	// Installed is not running. `docker info` is the cheapest question that
	// actually reaches the daemon; `docker --version` answers from the client
	// and says yes while the daemon is down.
	if _, err := runTool(ctx, 15*time.Second, "docker", "info", "--format", "{{.ServerVersion}}"); err != nil {
		return EnvFinding{
			EnvWarn, "Docker", "installed, but the daemon did not answer",
			problem("docker-not-running"),
		}
	}
	out, _ := runTool(ctx, 15*time.Second, "docker", "info", "--format", "{{.ServerVersion}}")
	return EnvFinding{EnvOK, "Docker", out, nil}
}

func checkPort(p PortCheck) EnvFinding {
	name := fmt.Sprintf("port %d", p.Port)
	// Listening is the question, so try to listen. Dialling it instead would
	// report a port as busy whenever something refuses the connection, and
	// free whenever a server is slow to accept.
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p.Port))
	if err != nil {
		return EnvFinding{
			EnvWarn, name,
			fmt.Sprintf("in use, and %s wants it", p.What),
			problem("port-in-use"),
		}
	}
	_ = ln.Close()
	return EnvFinding{EnvOK, name, "free, for " + p.What, nil}
}

func parseVersion(s string) (int, int, int) {
	parts := strings.SplitN(s, ".", 3)
	get := func(i int) int {
		if i >= len(parts) {
			return 0
		}
		n, _ := strconv.Atoi(parts[i])
		return n
	}
	return get(0), get(1), get(2)
}

func older(aMaj, aMin, aPatch, bMaj, bMin, bPatch int) bool {
	if aMaj != bMaj {
		return aMaj < bMaj
	}
	if aMin != bMin {
		return aMin < bMin
	}
	return aPatch < bPatch
}
