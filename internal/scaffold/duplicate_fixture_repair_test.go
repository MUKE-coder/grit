package scaffold

import (
	"strings"
	"testing"
)

func TestDropRedeclaredFixtures(t *testing.T) {
	// What v3.483.0 to v3.487.0 generated for a resource with a belongs_to.
	broken := `package handlers

func TestAttachment_AdminSeesWhatWasCreated(t *testing.T) {
	client, db, auth := mountAttachment(t)
	file := factory.CreateUpload(t, db)
	record := factory.CreateAttachment(t, db)
	_ = record
}

func TestAttachment_BulkEditRefusesAStaleRow(t *testing.T) {
	client, db, auth := mountAttachment(t)
	file := factory.CreateUpload(t, db)
	first := factory.CreateAttachment(t, db)
	file := factory.CreateUpload(t, db)
	second := factory.CreateAttachment(t, db)
	_, _ = first, second
}
`

	out, removed := dropRedeclaredFixtures(broken)
	if removed != 1 {
		t.Fatalf("removed %d lines, want 1", removed)
	}
	// The second function keeps one, not none: the first row needs the parent.
	second := out[strings.Index(out, "func TestAttachment_BulkEdit"):]
	if n := strings.Count(second, "file := factory.CreateUpload(t, db)"); n != 1 {
		t.Errorf("the bulk-edit test has %d fixture declarations, want 1:\n%s", n, second)
	}
	// And the first function is untouched: the same name in another function is
	// correct, and a repair that took it out would break a test that compiles.
	first := out[strings.Index(out, "func TestAttachment_AdminSees"):strings.Index(out, "func TestAttachment_BulkEdit")]
	if !strings.Contains(first, "file := factory.CreateUpload(t, db)") {
		t.Errorf("the repair removed a fixture from a function that only declared it once:\n%s", first)
	}
}

// TestDropRedeclaredFixturesLeavesAGoodFileAlone: the repair runs on every
// handlers test in every project, so being a no-op on a correct file is the
// property that matters most.
func TestDropRedeclaredFixturesLeavesAGoodFileAlone(t *testing.T) {
	fine := `package handlers

func TestTask_One(t *testing.T) {
	project := factory.CreateProject(t, db)
	assignee := factory.CreateUser(t, db)
	_, _ = project, assignee
}

func TestTask_Two(t *testing.T) {
	project := factory.CreateProject(t, db)
	_ = project
}
`
	out, removed := dropRedeclaredFixtures(fine)
	if removed != 0 {
		t.Errorf("removed %d lines from a file that compiles", removed)
	}
	if out != fine {
		t.Error("the file was rewritten with nothing to fix")
	}
}

// TestDropRedeclaredFixturesIgnoresHandWrittenLines: only the generator's exact
// shape. A line somebody wrote themselves, with an override or a different
// receiver, is theirs.
func TestDropRedeclaredFixturesIgnoresHandWrittenLines(t *testing.T) {
	mine := `package handlers

func TestThing(t *testing.T) {
	file := factory.CreateUpload(t, db)
	file := factory.CreateUpload(t, db, func(u *models.Upload) { u.Size = 1 })
	_ = file
}
`
	// The second line is not the generator's shape, so it is left as it is:
	// somebody wrote it, it does not compile, and guessing at what they meant
	// is worse than saying nothing.
	_, removed := dropRedeclaredFixtures(mine)
	if removed != 0 {
		t.Errorf("removed %d hand-written lines", removed)
	}
}
