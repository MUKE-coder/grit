package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A resource that belongs to a User and to another resource: every screen names
// the related record through relationLabel, which is written alongside, and the
// use-users hook the screens import exists.
func TestMobile_RelationsToUserCompile(t *testing.T) {
	g, root := mobileGen(t, "Message", "body:text,sender:belongs_to:User,conversation:belongs_to:Conversation")
	if err := g.writeMobileFiles(g.Names()); err != nil {
		t.Fatalf("writeMobileFiles: %v", err)
	}
	expo := filepath.Join(root, "apps", "expo")
	read := func(rel string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(expo, rel))
		if err != nil {
			t.Fatalf("reading %s: %v", rel, err)
		}
		return string(b)
	}
	for _, rel := range []string{"app/messages/index.tsx", "app/messages/[id].tsx"} {
		src := read(rel)
		if strings.Contains(src, ".name || ") {
			t.Errorf("%s still reads .name off a related record", rel)
		}
		if !strings.Contains(src, "relationLabel(item.sender)") || !strings.Contains(src, `from "@/lib/relation-label"`) {
			t.Errorf("%s does not name related records with relationLabel", rel)
		}
	}
	if !strings.Contains(read("lib/relation-label.ts"), "export function relationLabel(") {
		t.Error("lib/relation-label.ts was not written")
	}
	if !strings.Contains(read("app/messages/index.tsx"), `from "@/hooks/use-users"`) {
		t.Fatal("the list screen no longer imports the users hook this test is about")
	}
	if !strings.Contains(read("hooks/use-users.ts"), "export function useUsers(") {
		t.Error("the screens import @/hooks/use-users, but it was not written")
	}
}

// A users hook someone already has is left alone.
func TestMobile_ExistingUsersHookIsKept(t *testing.T) {
	g, root := mobileGen(t, "Note", "body:text,owner:belongs_to:User")
	mine := filepath.Join(root, "apps", "expo", "hooks", "use-users.ts")
	if err := os.MkdirAll(filepath.Dir(mine), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mine, []byte("// mine\nexport function useUsers() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := g.writeMobileFiles(g.Names()); err != nil {
		t.Fatalf("writeMobileFiles: %v", err)
	}
	b, err := os.ReadFile(mine)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(b), "// mine") {
		t.Error("an existing use-users.ts was overwritten")
	}
}
