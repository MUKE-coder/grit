package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// The Grit skill exists twice: skills/grit/SKILL.md, which `npx skills add`
// installs, and docs/public/skill.md, which any agent can fetch over HTTP when
// it has no skills CLI. Two copies of a teaching document drift, and the
// failure is silent: an agent following the stale one builds something subtly
// wrong and nothing reports it.
func TestThePublishedSkillMatchesTheInstallableOne(t *testing.T) {
	root := repoRoot(t)

	installable, err := os.ReadFile(filepath.Join(root, "skills", "grit", "SKILL.md"))
	if err != nil {
		t.Fatalf("reading the installable skill: %v", err)
	}
	published, err := os.ReadFile(filepath.Join(root, "docs", "public", "skill.md"))
	if err != nil {
		t.Fatalf("reading the published skill: %v", err)
	}
	if string(installable) != string(published) {
		t.Error("skills/grit/SKILL.md and docs/public/skill.md have drifted; copy one over the other")
	}
}

// repoRoot walks up from the test's directory until it finds go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not find the repository root")
	return ""
}
