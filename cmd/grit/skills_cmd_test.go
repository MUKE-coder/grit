package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The command exists because an agent cannot install these any other way:
// the skills ship as a Claude Code plugin, `/plugin` is a slash command a
// person types, and in the VS Code extension it does not exist at all. So the
// one thing that must stay true is that it is reachable as a plain command.
func TestSkillsCommandShape(t *testing.T) {
	c := skillsCmd()
	if c.Name() != "skills" {
		t.Fatalf("command is named %q", c.Name())
	}

	var sub []string
	for _, s := range c.Commands() {
		sub = append(sub, s.Name())
	}
	for _, want := range []string{"install", "status"} {
		var ok bool
		for _, s := range sub {
			if s == want {
				ok = true
			}
		}
		if !ok {
			t.Errorf("grit skills has no %q subcommand, got %v", want, sub)
		}
	}

	// The two flags the project briefs rely on.
	install, _, err := c.Find([]string{"install"})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"all", "project"} {
		if install.Flags().Lookup(f) == nil {
			t.Errorf("grit skills install has no --%s", f)
		}
	}
}

// The list drives both the status report and the post-install check, so a
// duplicate would overcount and a missing router would mean the one skill
// everything routes through is never verified.
func TestExpertSkillNamesIsSane(t *testing.T) {
	seen := map[string]bool{}
	for _, n := range expertSkillNames {
		if seen[n] {
			t.Errorf("%q listed twice", n)
		}
		seen[n] = true
		if !strings.HasPrefix(n, "grit-") {
			t.Errorf("%q is not a grit skill", n)
		}
	}
	for _, must := range []string{"grit-build", "grit-references"} {
		if !seen[must] {
			t.Errorf("%q missing: install would report success without it", must)
		}
	}
}

// Lstat, not Stat: the skills CLI installs a symlink, and a broken one is
// still something somebody installed and should be reported rather than
// silently counted as absent.
func TestSkillInstalledSeesASymlink(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "real")
	if err := os.MkdirAll(real, 0o755); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(dir, "grit-build")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks unavailable here: %v", err)
	}
	if !skillInstalled([]string{dir}, "grit-build") {
		t.Error("a symlinked skill should count as installed")
	}

	broken := filepath.Join(dir, "grit-core")
	if err := os.Symlink(filepath.Join(dir, "gone"), broken); err == nil {
		if !skillInstalled([]string{dir}, "grit-core") {
			t.Error("a broken symlink is still something to report, not absence")
		}
	}

	if skillInstalled([]string{dir}, "grit-nothing") {
		t.Error("an absent skill should not count as installed")
	}
}

// "npx not found" is the one failure a person can act on, so it has to say
// what to do rather than surfacing exec's own message.
func TestRequireNpxExplainsItself(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	err := requireNpx()
	if err == nil {
		t.Skip("npx resolved anyway on this platform")
	}
	for _, want := range []string{"nodejs.org", "npx skills add", expertPack} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error should mention %q, got: %v", want, err)
		}
	}
}
