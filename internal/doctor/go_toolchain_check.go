package doctor

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Is the Go on PATH one with known standard-library vulnerabilities?
//
// A generated project pins a safe Go in its go.mod and its Dockerfile, so what
// is deployed is fine. What is not is every local test run, every `grit start`
// and every desktop build, which use whatever `go` the developer has. A
// project sat on go1.27.2, which carries GO-2026-6809 in crypto/x509, and
// nothing said anything: it was found with govulncheck, by hand (grit#178).
//
// Offline, like every other check here. The floors are compiled in rather than
// fetched from go.dev, which is the same shape checkFrameworkDeps uses: it
// keeps `grit doctor` fast, usable in CI and identical on a machine with no
// network, and this project ships often enough that a compiled-in floor is
// fresher than most projects' upgrade cadence.
//
// The cost of that choice, stated rather than hidden: a Go security release
// published after this version of Grit is one this check does not know about.
// It is a floor, not a vulnerability scanner, and the message says so.

var goVersionRe = regexp.MustCompile(`go(\d+)\.(\d+)(?:\.(\d+))?`)

func checkGoToolchain(p *project) []Finding {
	// A project with no Go module is a frontend somebody is working on, and
	// asking them about a toolchain they do not use is noise.
	if p.goMod == "" {
		return nil
	}

	version, ok := localGoVersion()
	if !ok {
		// No Go on PATH. Not a finding: plenty of people work on the frontend
		// of a Grit project and never build the API.
		return nil
	}

	m := goVersionRe.FindStringSubmatch(version)
	if m == nil {
		return nil
	}
	minor := m[1] + "." + m[2]
	patch, _ := strconv.Atoi(m[3])

	floor, known := scaffold.GoPatchFloors[minor]
	if !known {
		// A minor this release has never heard of, which usually means a Go
		// newer than Grit. Warning there would be inventing a vulnerability in
		// something that did not exist yet.
		return nil
	}
	if patch >= floor {
		return nil
	}

	want := fmt.Sprintf("go%s.%d", minor, floor)
	return []Finding{{
		Level:    "warning",
		Resource: "go toolchain",
		Message: fmt.Sprintf(
			"go%s.%d is on your PATH, and go%s.%d is the oldest patch of %s without known "+
				"standard-library vulnerabilities. Your deployed image pins a safe Go, so this "+
				"is about local test runs, grit start and desktop builds",
			minor, patch, minor, floor, minor),
		Fix: "go env -w GOTOOLCHAIN=" + want + "   (or install it: https://go.dev/dl/)",
	}}
}

// localGoVersion is what `go version` prints, or false when there is no Go.
func localGoVersion() (string, bool) {
	path, err := exec.LookPath("go")
	if err != nil {
		return "", false
	}
	out, err := exec.Command(path, "version").Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), true
}
