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

// grit#160: the command writes into the user's agent configuration, from a
// private repository. Hidden made that unfindable, which is worse than noisy.
func TestSkillsCommandIsFindableAndSaysWhereTheFilesComeFrom(t *testing.T) {
	cmd := skillsCmd()

	if cmd.Hidden {
		t.Error("grit skills is Hidden again. A command that installs into somebody's " +
			"agent configuration has to be in --help, even though the pack is private: " +
			"a clear refusal naming the repository is a better answer than invisibility.")
	}

	// Three artefacts get called a Grit skill. The help has to tell them apart,
	// because "there are two skill definitions, one of which cannot be read"
	// was a fair reading of the old help and is what the issue said.
	//
	// "private repository" was on this list until the pack was published. The
	// test is about telling the three apart, not about the pack's visibility,
	// and TestTheHelpTextDoesNotClaimThePackIsPrivate now asserts the opposite
	// of what that entry did: two tests requiring contradictory things is how
	// a true statement gets reverted to a false one.
	for _, want := range []string{
		expertPack,
		"skills/grit/SKILL.md",
		".claude/skills/grit/SKILL.md",
	} {
		if !strings.Contains(cmd.Long, want) {
			t.Errorf("grit skills --help does not mention %q, so a reader cannot tell "+
				"which of the three Grit skills this installs", want)
		}
	}
}

// The number in the help text is the number of skills in the pack.
//
// `grit skills --help` says how many skills the pack holds, and the list is
// right beside it in the same file, so the two drift the moment somebody adds
// a skill and does not reread the prose. It said twenty-seven while the list
// held twenty-eight, which is the shape of mistake nobody notices because the
// sentence still reads correctly (grit#160).
func TestTheHelpTextCountsThePackCorrectly(t *testing.T) {
	spelled := map[int]string{
		25: "twenty-five", 26: "twenty-six", 27: "twenty-seven",
		28: "twenty-eight", 29: "twenty-nine", 30: "thirty",
		31: "thirty-one", 32: "thirty-two", 33: "thirty-three",
	}
	want, ok := spelled[len(expertSkillNames)]
	if !ok {
		t.Fatalf("the pack has %d skills and this test cannot spell that; add it to the map",
			len(expertSkillNames))
	}

	long := skillsCmd().Long
	if !strings.Contains(long, want+" expert skills") {
		t.Errorf("the pack has %d skills and the help text does not say %q.\n"+
			"Help text says: %s", len(expertSkillNames), want, packCountLine(long))
	}
	// And not some other number, which is what a search for the right one
	// alone would miss.
	for n, word := range spelled {
		if n == len(expertSkillNames) {
			continue
		}
		if strings.Contains(long, word+" expert skills") {
			t.Errorf("the help text also says %q, so it names two different counts", word)
		}
	}
}

// packCountLine is the line of the help text that counts the pack, for an
// error message that shows what it actually says.
func packCountLine(long string) string {
	for _, line := range strings.Split(long, "\n") {
		if strings.Contains(line, "expert skills, private") {
			return strings.TrimSpace(line)
		}
	}
	return "(no line counting the pack)"
}

// The help text does not tell people they cannot read the pack.
//
// It said "a private repository" in three places, which was true until the
// repository was published. A claim about who may read an agent skill is the
// one thing in this help text that has to be right, because it is the whole
// subject of grit#160: a skill tells an agent how to write code in somebody's
// project, so whether they can read it first is the question.
func TestTheHelpTextDoesNotClaimThePackIsPrivate(t *testing.T) {
	long := skillsCmd().Long

	for _, wrong := range []string{"private repository", ", private,", "private pack"} {
		if strings.Contains(long, wrong) {
			t.Errorf("the help text still says %q. The pack at %s is public; if that changes, "+
				"change this test and the text together rather than only one of them",
				wrong, expertPack)
		}
	}
	// And it says the thing that replaced it, so this is not just the absence
	// of a word.
	if !strings.Contains(long, "It is public, so you can read every skill") {
		t.Error("the help text no longer says the pack is readable before installing, which is " +
			"the question grit#160 was about")
	}
}
