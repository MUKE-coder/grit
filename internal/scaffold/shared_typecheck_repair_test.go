package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readSharedPkg(t *testing.T, root string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, "packages", "shared", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pkg map[string]any
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatalf("not valid JSON after the repair: %v\n%s", err, raw)
	}
	return pkg
}

func writeSharedPkg(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "packages", "shared")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFileT(t, filepath.Join(dir, "package.json"), body)
	return root
}

// The ordinary case: an edited package.json with no scripts block at all,
// which is every project generated before this release.
func TestTheRepairAddsAScriptsBlock(t *testing.T) {
	root := writeSharedPkg(t, `{
  "name": "@repo/shared",
  "version": "0.1.0",
  "private": true,
  "exports": {
    "./schemas": "./schemas/index.ts"
  },
  "dependencies": {
    "zod": "^4.6.7",
    "my-own-thing": "^1.0.0"
  },
  "devDependencies": {
    "typescript": "^5.3.0"
  }
}
`)
	if err := addSharedTypeCheckScript(root); err != nil {
		t.Fatal(err)
	}
	pkg := readSharedPkg(t, root)
	scripts, ok := pkg["scripts"].(map[string]any)
	if !ok {
		t.Fatalf("no scripts block was added: %v", pkg)
	}
	if scripts["type-check"] != sharedTypeCheckScript {
		t.Errorf("type-check is %v, want %q", scripts["type-check"], sharedTypeCheckScript)
	}
	// The project's own dependency has to survive, which is the whole reason
	// this is an edit rather than a rewrite from the template.
	deps, _ := pkg["dependencies"].(map[string]any)
	if deps["my-own-thing"] != "^1.0.0" {
		t.Errorf("the project's own dependency was lost: %v", deps)
	}
}

// A project that has a scripts block already: the entry goes in that block,
// and the scripts it has keep working.
func TestTheRepairAddsToAnExistingScriptsBlock(t *testing.T) {
	root := writeSharedPkg(t, `{
  "name": "@repo/shared",
  "scripts": {
    "build": "tsc"
  },
  "dependencies": {
    "zod": "^4.6.7"
  }
}
`)
	if err := addSharedTypeCheckScript(root); err != nil {
		t.Fatal(err)
	}
	scripts := readSharedPkg(t, root)["scripts"].(map[string]any)
	if scripts["type-check"] != sharedTypeCheckScript {
		t.Errorf("type-check is %v", scripts["type-check"])
	}
	if scripts["build"] != "tsc" {
		t.Errorf("the existing build script was lost: %v", scripts)
	}
}

// Already there: nothing changes, including the bytes. A repair that rewrites
// an already-correct file shows up as a change on every upgrade forever.
func TestTheRepairLeavesAProjectThatHasItAlone(t *testing.T) {
	body := `{
  "name": "@repo/shared",
  "scripts": {
    "type-check": "tsc --noEmit --pretty"
  }
}
`
	root := writeSharedPkg(t, body)
	if err := addSharedTypeCheckScript(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "packages", "shared", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != body {
		t.Errorf("a project with its own type-check script was rewritten:\n%s", raw)
	}
}

// A project with no shared package at all, which is every --api and --single
// shape. Silence, not an error.
func TestTheRepairIsSilentWithNoSharedPackage(t *testing.T) {
	if err := addSharedTypeCheckScript(t.TempDir()); err != nil {
		t.Errorf("a project with no packages/shared returned %v", err)
	}
}

// And the template itself has to be idempotent under the repair: delivering a
// fresh package.json and then running the repair over it must change nothing.
func TestTheTemplateAlreadySatisfiesTheRepair(t *testing.T) {
	root := writeSharedPkg(t, sharedPackageJSON(tripleOptions()))
	if err := addSharedTypeCheckScript(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "packages", "shared", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != sharedPackageJSON(tripleOptions()) {
		t.Error("the repair edits the file the template just wrote, so every upgrade reports a change")
	}
	if !strings.Contains(string(raw), "type-check") {
		t.Error("the template does not declare the script the repair exists to add")
	}
}
