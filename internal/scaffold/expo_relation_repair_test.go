package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpoRelationRepair(t *testing.T) {
	before := "import { View } from \"react-native\";\n" +
		"const a = (item.sender && (item.sender.name || item.sender.title)) || item.sender_id;\n" +
		"const b = (i.conversation && (i.conversation.name || i.conversation.title)) || \"\";\n"
	got, fixed, warn := repairExpoRelationSource(before)
	if len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{"relationLabel(item.sender) || item.sender_id", "relationLabel(i.conversation) || \"\"", relationLabelImport} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if !strings.HasPrefix(got, "import { View } from \"react-native\";\n"+relationLabelImport) {
		t.Error("the import is not right after the first one")
	}
	if again, fixed, _ := repairExpoRelationSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
	plain := "import x from \"y\";\nconst n = user.name;\n"
	if out, fixed, _ := repairExpoRelationSource(plain); out != plain || len(fixed) != 0 {
		t.Error("a file without the old expression was changed")
	}
}

// The whole repair on a project: screens under app/ and components/ are fixed,
// and the two files they now import are written.
func TestExpoRelationRepairOnAProject(t *testing.T) {
	root := t.TempDir()
	expo := filepath.Join(root, "apps", "expo")
	screen := filepath.Join(expo, "app", "(tabs)", "messages.tsx")
	form := filepath.Join(expo, "components", "resource-forms", "messages-form.tsx")
	old := "import { View } from \"react-native\";\nimport { useUsers } from \"@/hooks/use-users\";\n" +
		"const a = (item.sender && (item.sender.name || item.sender.title)) || item.sender_id;\n"
	for _, p := range []string{screen, form} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(old), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := repairExpoRelations(root); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{screen, form} {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "relationLabel(item.sender)") {
			t.Errorf("%s was not repaired:\n%s", p, b)
		}
	}
	for _, rel := range []string{"lib/relation-label.ts", "hooks/use-users.ts"} {
		if !fileExists(filepath.Join(expo, rel)) {
			t.Errorf("%s was not written", rel)
		}
	}
}

func TestExpoRouteRepair(t *testing.T) {
	before := "onPress={() => router.push(\"/conversations/\" + item.id)}\nrouter.replace(\"/roles/\" + role.id);\n"
	got, fixed, _ := repairExpoRouteSource(before)
	want := "onPress={() => router.push(`/conversations/${item.id}`)}\nrouter.replace(`/roles/${role.id}`);\n"
	if got != want || len(fixed) != 1 {
		t.Fatalf("got %q (fixed %v), want %q", got, fixed, want)
	}
	if again, fixed, _ := repairExpoRouteSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
	if !strings.Contains(expoRolesListScreen(), "router.push(`/roles/${role.id}`)") {
		t.Error("the roles screen template still concatenates its route")
	}
}

func TestDesktopFormRepair(t *testing.T) {
	before := "import type { Message, MessageInput } from \"@/hooks/use-messages\";\n" +
		"  const handleSubmit = async (e: React.FormEvent) => {\n" +
		"    e.preventDefault();\n" +
		"    await onSubmit({\n" +
		"      kind: kind,\n" +
		"      attachment: attachment,\n" +
		"    });\n" +
		"  };\n"
	got, fixed, _ := repairDesktopFormSource(before)
	if len(fixed) != 1 || !strings.Contains(got, "      attachment: attachment,\n    } as MessageInput);") {
		t.Fatalf("not repaired (fixed %v):\n%s", fixed, got)
	}
	if again, fixed, _ := repairDesktopFormSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
}

func TestDesktopListRepair(t *testing.T) {
	before := "const COLUMNS: DataColumn[] = [\n" +
		"  { key: \"conversation\", label: \"Conversation\", format: \"text\" },\n" +
		"  { key: \"body\", label: \"Body\", format: \"text\" },\n];\n" +
		"  const rows = items.map((r) => ({ ...r, conversation: conversationMap.get(String((r as any).conversation_id)) ?? \"\" }));\n"
	got, fixed, _ := repairDesktopListSource(before)
	for _, want := range []string{`{ key: "conversation_label", label: "Conversation"`, `conversation_label: conversationMap.get(`, `{ key: "body", label: "Body"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if len(fixed) != 1 {
		t.Errorf("fixed %v", fixed)
	}
	if again, fixed, _ := repairDesktopListSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
}

func TestDesktopTSConfigNodeHasAnOutDir(t *testing.T) {
	now := desktopClientTSConfigNode()
	if !strings.Contains(now, desktopTSConfigNodeOutDir) {
		t.Fatal("a new desktop app's tsconfig.node.json has no outDir")
	}
	before := strings.Replace(now, desktopTSConfigNodeOutDir, "", 1)
	got, fixed, _ := repairDesktopTSConfigNodeSource(before)
	if got != now || len(fixed) != 1 {
		t.Fatalf("repair did not produce the current config:\n%s", got)
	}
	if again, fixed, _ := repairDesktopTSConfigNodeSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
}
