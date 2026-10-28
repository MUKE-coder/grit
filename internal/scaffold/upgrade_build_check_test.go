package scaffold

import (
	"strings"
	"testing"
)

// The build check's job is to say what broke, briefly.
func TestBuildErrorLinesDropsPackageHeadersAndCaps(t *testing.T) {
	out := `# whatsapp/apps/api/internal/services
internal\services\security_autoban.go:84:20: cfg.SecurityAutoBanEnabled undefined
internal\services\security_autoban.go:85:20: cfg.SecurityAutoBanMinThreats undefined

# whatsapp/apps/api/internal/jobs
internal\jobs\workers.go:35:11: undefined: services
`
	lines := buildErrorLines(out, 12)
	if len(lines) != 3 {
		t.Fatalf("kept %d lines, wanted the 3 compiler errors: %v", len(lines), lines)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			t.Errorf("a package header survived: %q", line)
		}
	}

	// A project with many errors prints a few and a count, not a wall.
	var many strings.Builder
	for i := 0; i < 40; i++ {
		many.WriteString("internal/x.go:1:1: something undefined\n")
	}
	capped := buildErrorLines(many.String(), 12)
	if len(capped) != 13 {
		t.Fatalf("capped to %d lines, wanted 12 plus the count", len(capped))
	}
	if !strings.Contains(capped[12], "28 more") {
		t.Errorf("the last line should count the rest, got %q", capped[12])
	}
}

// A missing directory or toolchain is not a broken project.
func TestBuildCheckIsSilentWhenThereIsNothingToBuild(t *testing.T) {
	if !reportBuildAfterUpgrade(t.TempDir() + "/not-here") {
		t.Error("reported a failure for a directory that does not exist")
	}
}
