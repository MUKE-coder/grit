package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Review 2.27: `grit init --force` replaced the whole of CLAUDE.md.
//
// CLAUDE.md is where a project keeps its own instructions for whoever works on
// it. Refreshing a list of framework rules is not a reason to delete them, and
// --force was the documented way to refresh.
//
// Grit owns one marked section now. These are the four cases, and the only one
// that loses anything is the one whose flag says it will.
func TestInitKeepsWhatTheProjectWrote(t *testing.T) {
	const theirs = "# My project\n\nNever touch the billing module without asking Sarah.\n"

	t.Run("a file with their text keeps it", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "CLAUDE.md")
		if err := os.WriteFile(path, []byte(theirs), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := WriteAgentsDoc(dir, false); err != nil {
			t.Fatal(err)
		}
		got := read(t, path)
		if !strings.Contains(got, "billing module") {
			t.Error("their instructions were deleted")
		}
		if !strings.Contains(got, agentsDocStart) || !strings.Contains(got, agentsDocEnd) {
			t.Error("the conventions were not written")
		}
		// Theirs stays at the top, where they left it.
		if strings.Index(got, "billing module") > strings.Index(got, agentsDocStart) {
			t.Error("the conventions were put above their own instructions")
		}
	})

	t.Run("a refresh keeps text on both sides", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "CLAUDE.md")
		stale := theirs + "\n" + agentsDocStart + "\nRULES FROM AN ANCIENT VERSION\n" + agentsDocEnd +
			"\n\nAlso: the staging database is shared, do not truncate it.\n"
		if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := WriteAgentsDoc(dir, false); err != nil {
			t.Fatal(err)
		}
		got := read(t, path)
		for _, want := range []string{"billing module", "staging database"} {
			if !strings.Contains(got, want) {
				t.Errorf("a refresh deleted %q, which is on the other side of the markers", want)
			}
		}
		if strings.Contains(got, "ANCIENT VERSION") {
			t.Error("the stale section was not replaced")
		}
	})

	t.Run("running it twice changes nothing", func(t *testing.T) {
		dir := t.TempDir()
		if _, err := WriteAgentsDoc(dir, false); err != nil {
			t.Fatal(err)
		}
		first := read(t, filepath.Join(dir, "CLAUDE.md"))

		written, err := WriteAgentsDoc(dir, false)
		if err != nil {
			t.Fatal(err)
		}
		if len(written) != 0 {
			t.Errorf("the second run rewrote %v", written)
		}
		if read(t, filepath.Join(dir, "CLAUDE.md")) != first {
			t.Error("the second run changed the file")
		}
	})

	t.Run("force replaces the whole file, as its help says", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "CLAUDE.md")
		if err := os.WriteFile(path, []byte(theirs), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := WriteAgentsDoc(dir, true); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(read(t, path), "billing module") {
			t.Error("--force kept their text, so the flag no longer does what it says")
		}
	})
}

// An opening marker with no close is a file somebody edited and left broken.
// Replacing from there to the end of the file would delete whatever follows,
// so it is left alone and the block is appended instead.
func TestAnUnclosedMarkerIsNotTreatedAsASection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	broken := "# Mine\n" + agentsDocStart + "\nhalf a section\n\nAnd this line matters.\n"
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteAgentsDoc(dir, false); err != nil {
		t.Fatal(err)
	}
	got := read(t, path)
	if !strings.Contains(got, "And this line matters.") {
		t.Error("an unclosed marker caused everything after it to be deleted")
	}
}

// Both filenames, because different tooling reads different ones and a project
// that gets one and not the other has two sources of truth.
func TestInitWritesBothFilenames(t *testing.T) {
	dir := t.TempDir()
	written, err := WriteAgentsDoc(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 2 {
		t.Fatalf("wrote %v, want both CLAUDE.md and AGENTS.md", written)
	}
	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		if !strings.Contains(read(t, filepath.Join(dir, name)), agentsDocStart) {
			t.Errorf("%s has no marked section", name)
		}
	}
}
