package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func planFor(t *testing.T, skipped []SkippedFile) string {
	t.Helper()
	root := t.TempDir()
	path, err := writeUpgradePlan(root, "3.350.0", "3.358.1", skipped, 498)
	if err != nil {
		t.Fatalf("writeUpgradePlan: %v", err)
	}
	if filepath.Base(path) != PlanFile {
		t.Errorf("wrote %s, want %s", filepath.Base(path), PlanFile)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the plan: %v", err)
	}
	return string(body)
}

func TestPlanNamesBothVersions(t *testing.T) {
	plan := planFor(t, []SkippedFile{{
		Rel:      "docker-compose.yml",
		Current:  "services:\n  postgres:\n    image: postgres:15\n",
		Proposed: "services:\n  postgres:\n    image: postgres:16\n",
	}})

	// Which version the project is coming from is the single most useful line
	// in the file: it is what tells the reader where to look at what changed.
	for _, want := range []string{"v3.358.1", "v3.350.0", "498"} {
		if !strings.Contains(plan, want) {
			t.Errorf("the plan does not mention %q", want)
		}
	}
}

// Detect, change, verify. An instruction with no verification step is one an
// agent will report as done whether or not it is.
func TestEveryEntryHasAllThreeSteps(t *testing.T) {
	plan := planFor(t, []SkippedFile{{
		Rel:      "apps/api/internal/config/config.go",
		Current:  "package config\n\nvar A = 1\n",
		Proposed: "package config\n\nvar A = 2\n",
	}})

	for _, want := range []string{
		"**Detect.**",
		"**Change.**",
		"**Verify.**",
		"apps/api/internal/config/config.go",
		"```diff",
		// A Go file verifies from its own module directory, not the repo root.
		"cd apps/api && gofmt -l . && go build ./... && go vet ./...",
	} {
		if !strings.Contains(plan, want) {
			t.Errorf("the plan is missing %q", want)
		}
	}
}

func TestVerifyCommandFollowsTheFileType(t *testing.T) {
	for _, tc := range []struct{ rel, want string }{
		{"apps/api/internal/config/config.go", "cd apps/api && gofmt"},
		{"apps/admin/app/page.tsx", "cd apps/admin && npx tsc --noEmit"},
		{"apps/web/lib/api.ts", "cd apps/web && npx tsc --noEmit"},
		{"docker-compose.yml", "docker compose -f docker-compose.yml config"},
		{"package.json", "JSON.parse"},
		{"turbo.json", "JSON.parse"},
	} {
		if got := verifyCommandFor(tc.rel); !strings.Contains(got, tc.want) {
			t.Errorf("%s: verify = %q, want it to contain %q", tc.rel, got, tc.want)
		}
	}

	// A file type with no useful check gets no command, rather than one that
	// always passes and teaches the reader to skip the step.
	if got := verifyCommandFor("README.md"); got != "" {
		t.Errorf("README.md: verify = %q, want none", got)
	}
}

// Most upgrades leave nothing outstanding, and a runbook that lists no work
// should say so in one line rather than open with three paragraphs of
// instructions for work that does not exist.
func TestNothingOutstandingSaysSo(t *testing.T) {
	plan := planFor(t, nil)

	if !strings.Contains(plan, "Nothing is outstanding") {
		t.Error("the plan does not say that there is nothing to do")
	}
	if strings.Contains(plan, "**Detect.**") {
		t.Error("the plan has instructions for work it has not got")
	}
}

// A runbook has to stay something a person or a model reads rather than scrolls
// past, and a file whose diff runs to thousands of lines is one where the diff
// was never going to be the useful part.
func TestALongDiffIsTrimmedAndSaysSo(t *testing.T) {
	var current, proposed strings.Builder
	for i := 0; i < 1000; i++ {
		current.WriteString("old line\n")
		proposed.WriteString("new line\n")
	}

	plan := planFor(t, []SkippedFile{{
		Rel:      "apps/web/app/page.tsx",
		Current:  current.String(),
		Proposed: proposed.String(),
	}})

	if !strings.Contains(plan, "more lines. Run `grit upgrade --diff`") {
		t.Error("a long diff was not trimmed, or does not say where the rest is")
	}
	if lines := strings.Count(plan, "\n"); lines > 400 {
		t.Errorf("the plan is %d lines; the trim did not hold", lines)
	}
}
