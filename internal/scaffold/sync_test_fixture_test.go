package scaffold

import (
	"strings"
	"testing"
)

// The sync engine's tests used models.Blog as their fixture, because it was
// there and it had a Version field. That coupled the sync tests to a resource
// a project is allowed to delete: after `grit remove resource Blog`,
// `go build ./...` still passed and `go test ./internal/handlers/` did not,
// so a project looked healthy until CI ran (grit#74).
//
// Nothing in those tests is about blogs. They own their fixture now.
func TestSyncTestsDoNotDependOnTheBlogResource(t *testing.T) {
	for name, body := range map[string]string{
		"sync_version_test.go": apiSyncVersionTestGo(),
		"sync_pull_test.go":    apiSyncPullTestGo(),
	} {
		for _, line := range strings.Split(body, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue // prose explaining the history is fine
			}
			if strings.Contains(line, "models.Blog") {
				t.Errorf("%s still uses models.Blog: %s", name, trimmed)
			}
			if strings.Contains(line, `"blogs"`) {
				t.Errorf("%s still registers the blogs model: %s", name, trimmed)
			}
		}
	}
}

// The fixture has to carry the hook, because the hook is the thing under test:
// a full-struct Save silently overwrites what SetColumn put in the statement,
// which is the bug these tests exist to catch. A fixture without it would pass
// while testing nothing.
func TestSyncFixtureKeepsTheVersionHook(t *testing.T) {
	body := apiSyncVersionTestGo()
	for _, want := range []string{
		"type syncFixture struct",
		"func (f *syncFixture) BeforeUpdate(",
		`tx.Statement.SetColumn("version", gorm.Expr("version + 1"))`,
		"Version   int",
		"DeletedAt gorm.DeletedAt",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the sync fixture is missing %q", want)
		}
	}
}
