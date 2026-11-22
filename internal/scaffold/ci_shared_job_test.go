package scaffold

import (
	"strings"
	"testing"
)

// A project with a shared package gets the job that checks it.
//
// Everything else in the generated CI is happy with a shared type a release
// behind its model: the Go compiles, the TypeScript type-checks, and the
// field is simply absent (grit#255).
func TestTheGeneratedCIChecksTheSharedTypes(t *testing.T) {
	opts := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}
	opts.Normalize()

	ci := ciYAML(opts)
	for _, want := range []string{
		"name: Shared types match the models",
		"grit sync --check",
		"go install github.com/MUKE-coder/grit/v3/cmd/grit@v",
	} {
		if !strings.Contains(ci, want) {
			t.Errorf("the generated CI does not contain %q", want)
		}
	}

	// Pinned to the version that scaffolded the project, not to latest: a
	// release should not be able to change what somebody's CI asserts.
	if !strings.Contains(ci, "cmd/grit@v"+opts.Version) {
		t.Errorf("the CLI is not pinned to the project's own version %q", opts.Version)
	}
	if strings.Contains(ci, "cmd/grit@latest") {
		t.Error("the generated CI installs the CLI at latest")
	}
}

// A project with no shared package does not get a job that would fail on it.
func TestTheSharedJobIsAbsentWithoutASharedPackage(t *testing.T) {
	opts := Options{ProjectName: "app", Architecture: ArchAPI}
	opts.Normalize()
	if opts.ShouldIncludeShared() {
		t.Skip("an --api project has a shared package now, so this guards nothing")
	}

	if ci := ciYAML(opts); strings.Contains(ci, "grit sync --check") {
		t.Error("an API-only project's CI runs a check against a package it does not have")
	}
}

// The job is valid YAML in the shape the rest of the file uses: two spaces of
// indent for a job, four for its keys. A broken workflow file is rejected by
// GitHub on the first push, which is a bad way to find out.
func TestTheSharedJobIndentsLikeTheOthers(t *testing.T) {
	opts := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}
	opts.Normalize()

	lines := strings.Split(ciYAML(opts), "\n")
	found := false
	for i, line := range lines {
		if line != "  shared:" {
			continue
		}
		found = true
		if next := lines[i+1]; !strings.HasPrefix(next, "    name:") {
			t.Errorf("the line after the job key is %q, want a four-space name:", next)
		}
	}
	if !found {
		t.Error("no `  shared:` job key at two spaces of indent")
	}
}
