package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A group emptied by removing its only feature used to stay, showing an empty
// heading in the role editor. Reported after removing the blog resource from a
// real project (grit#74).
func TestPruneEmptyPermissionGroups(t *testing.T) {
	const src = `package authz

var Catalog = []Module{
	{
		Key:  "content",
		Name: "Content",
		Groups: []Group{
			{
				Key:  "publishing",
				Name: "Publishing",
				Features: []Feature{
				},
			},
			{
				Key:  "media",
				Name: "Media",
				Features: []Feature{
					{Key: "uploads", Name: "Uploads", Actions: AllActions},
				},
			},
		},
	},
}
`
	p := filepath.Join(t.TempDir(), "permissions.go")
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if !pruneEmptyPermissionGroups(p) {
		t.Fatal("an empty group should have been pruned")
	}
	out := readBack(t, p)
	if strings.Contains(out, "publishing") {
		t.Error("the empty group survived")
	}
	if !strings.Contains(out, `Key:  "media"`) {
		t.Error("a populated group was removed")
	}
}

// The scaffold explains its built-in features with a comment above them, and
// that comment outlives the feature. A group holding only comments is still
// empty, and the first version of this check did not think so.
func TestPruneEmptyPermissionGroupsSeesPastComments(t *testing.T) {
	const src = `package authz

var Catalog = []Module{
	{
		Key:  "content",
		Groups: []Group{
			{
				Key:  "publishing",
				Features: []Feature{
					// Blog ships as a built-in resource.
					// It belongs in the catalog like any other.
				},
			},
		},
	},
}
`
	p := filepath.Join(t.TempDir(), "permissions.go")
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if !pruneEmptyPermissionGroups(p) {
		t.Fatal("a group holding only comments is empty and should be pruned")
	}
	if strings.Contains(readBack(t, p), "publishing") {
		t.Error("the comment-only group survived")
	}
}

// And it must never touch a group that still has a feature, or removing one
// resource takes another resource's permissions with it.
func TestPruneEmptyPermissionGroupsLeavesPopulatedOnesAlone(t *testing.T) {
	const src = `package authz

var Catalog = []Module{
	{
		Groups: []Group{
			{
				Key:  "media",
				Features: []Feature{
					{Key: "uploads", Name: "Uploads", Actions: AllActions},
				},
			},
		},
	},
}
`
	p := filepath.Join(t.TempDir(), "permissions.go")
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if pruneEmptyPermissionGroups(p) {
		t.Error("nothing was empty, so nothing should have been pruned")
	}
	if readBack(t, p) != src {
		t.Error("the file was rewritten when it should not have been")
	}
}

func readBack(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
