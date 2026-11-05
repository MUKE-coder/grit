package doctor

import (
	"strconv"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// grit#178: a project ran every local test on go1.27.2, which carries
// GO-2026-6809 in crypto/x509, and grit doctor said nothing. The image pinned
// a safe Go, so what was deployed was fine; what was not fine was every local
// test run, grit start and desktop build.
//
// The version string is the one thing this check cannot get from a fixture, so
// the parsing and the comparison are tested directly and the lookup through
// `go version` is left to the one real run.
func TestAVulnerableToolchainIsNamedWithItsFix(t *testing.T) {
	// The exact case from the report.
	const reported = "go version go1.27.2 windows/amd64"

	m := goVersionRe.FindStringSubmatch(reported)
	if m == nil {
		t.Fatalf("could not read a version out of %q", reported)
	}
	minor := m[1] + "." + m[2]
	patch, _ := strconv.Atoi(m[3])

	floor, known := scaffold.GoPatchFloors[minor]
	if !known {
		t.Fatalf("1.27 is not in GoPatchFloors, so the reported case would be silent")
	}
	if patch >= floor {
		t.Errorf("go%s.%d is not below the floor go%s.%d, so the reported case "+
			"would not warn", minor, patch, minor, floor)
	}
}

// A Go newer than this release knows about is not warned on. Claiming a
// vulnerability in something that did not exist when Grit shipped would be a
// guess, and a check that guesses is one people switch off.
func TestAnUnknownGoMinorIsNotWarnedOn(t *testing.T) {
	if _, known := scaffold.GoPatchFloors["1.99"]; known {
		t.Fatal("1.99 is in the table, so this test proves nothing")
	}
	m := goVersionRe.FindStringSubmatch("go version go1.99.0 linux/amd64")
	if m == nil {
		t.Fatal("could not parse a future version")
	}
	if _, known := scaffold.GoPatchFloors[m[1]+"."+m[2]]; known {
		t.Error("a future minor resolved to a floor")
	}
}

// A project with no Go module is a frontend somebody is working on.
func TestTheToolchainCheckIsQuietWithoutAGoModule(t *testing.T) {
	if findings := checkGoToolchain(&project{}); len(findings) != 0 {
		t.Errorf("a project with no go.mod was warned about a toolchain it does not use: %q",
			findings[0].Message)
	}
}

// The floor the scaffold writes into a generated go.mod has to be in the table,
// or doctor would accept a toolchain older than the project's own pin.
func TestTheGeneratedPinIsInTheFloorTable(t *testing.T) {
	parts := strings.Split(scaffold.GoToolchainPin, ".")
	if len(parts) != 3 {
		t.Fatalf("GoToolchainPin = %q, want major.minor.patch", scaffold.GoToolchainPin)
	}
	minor := parts[0] + "." + parts[1]
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		t.Fatalf("GoToolchainPin patch %q is not a number", parts[2])
	}
	floor, known := scaffold.GoPatchFloors[minor]
	if !known {
		t.Fatalf("the generated go.mod pins Go %s and GoPatchFloors has no entry for %s, "+
			"so doctor would accept a toolchain older than the project's own pin",
			scaffold.GoToolchainPin, minor)
	}
	if floor > patch {
		t.Errorf("GoPatchFloors[%s] = %d is above the pin's patch %d: a generated project "+
			"would be warned about the Go it was told to use", minor, floor, patch)
	}
}
