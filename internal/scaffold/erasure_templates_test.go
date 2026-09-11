package scaffold

import (
	"strings"
	"testing"
)

// Right-to-erasure deleted from a fixed list of framework tables, so records a
// user owned in a generated resource survived their erasure.
func TestEraseUsesTheRegistry(t *testing.T) {
	src := apiGDPRServiceGo()
	if !strings.Contains(src, "erasure.Scrub(tx, targetID)") {
		t.Error("EraseUser does not go through the erasure registry")
	}
	if strings.Contains(src, "func erasableTargets()") {
		t.Error("the fixed list of erasable tables is still there")
	}
	if !strings.Contains(apiErasableModelsGo(), `erasure.Register(m, "user_id")`) {
		t.Error("the framework's own tables do not register for erasure")
	}
}

// Restoring a backup taken before an erasure brought the person back and
// deleted the journal entry that proved they had been erased.
func TestRestoreReappliesLaterErasures(t *testing.T) {
	src := backupRestoreGo()
	for _, want := range []string{
		"tx.Order(\"created_at asc, id asc\").Find(&journal)",
		"reapplied, err := reapplyErasures(tx, journal)",
		"erasure.Scrub(tx, j.DeletedUserID)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("restore is missing %s", want)
		}
	}
	// The journal is read before the tables are cleared, or there is nothing
	// left to read.
	if strings.Index(src, "Find(&journal)") > strings.Index(src, "RESTART IDENTITY CASCADE") {
		t.Error("the journal is read after the tables are truncated")
	}
}
