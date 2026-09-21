package plugin

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPushPluginIsRegistered(t *testing.T) {
	p, err := Get("push")
	if err != nil {
		t.Fatalf("the push plugin is not registered: %v", err)
	}
	if len(p.NodeDeps) != 1 || p.NodeDeps[0].Workspace != "apps/expo" {
		t.Errorf("NodeDeps = %+v, want expo-notifications in apps/expo", p.NodeDeps)
	}
}

// The API half always installs; the Expo client only where there is an app.
func TestPushFilesFollowTheProject(t *testing.T) {
	root := t.TempDir()
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple"}
	without := pushFiles(ctx)
	if _, ok := without["apps/expo/lib/push.ts"]; ok {
		t.Error("the Expo client was written to a project with no Expo app")
	}
	for _, want := range []string{"apps/api/internal/models/push_token.go", "apps/api/internal/services/push.go", "apps/api/internal/services/push_test.go", "apps/api/internal/handlers/push.go"} {
		if _, ok := without[want]; !ok {
			t.Errorf("missing %s", want)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "apps", "expo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, ok := pushFiles(ctx)["apps/expo/lib/push.ts"]; !ok {
		t.Error("the Expo client was not written to a project with an Expo app")
	}
}

// Every Go template parses once the module is filled in, and none keeps the
// placeholder.
func TestPushTemplatesAreGo(t *testing.T) {
	ctx := Context{Root: t.TempDir(), Module: "shop/apps/api", Architecture: "triple"}
	for path, src := range pushFiles(ctx) {
		if strings.Contains(src, "{{MODULE}}") {
			t.Errorf("%s still has {{MODULE}}", path)
		}
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		if _, err := parser.ParseFile(token.NewFileSet(), path, src, parser.AllErrors); err != nil {
			t.Errorf("%s does not parse: %v", path, err)
		}
	}
}

func TestAddNodeDependencyKeepsTheFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	orig := "{\n  \"name\": \"expo\",\n  \"dependencies\": {\n    \"expo\": \"~54.0.0\"\n  }\n}\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	added, err := addNodeDependency(path, "expo-notifications", "~0.32.17")
	if err != nil || !added {
		t.Fatalf("added %v, err %v", added, err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"name\": \"expo\",\n  \"dependencies\": {\n    \"expo-notifications\": \"~0.32.17\",\n    \"expo\": \"~54.0.0\"\n  }\n}\n"
	if string(got) != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
	if again, _ := addNodeDependency(path, "expo-notifications", "~0.32.17"); again {
		t.Error("a dependency already listed was added again")
	}

	empty := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(empty, []byte("{\n  \"dependencies\": {}\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := addNodeDependency(empty, "a", "1"); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(empty)
	if strings.Contains(string(b), "\"1\",") {
		t.Errorf("the only dependency got a trailing comma:\n%s", b)
	}
}
