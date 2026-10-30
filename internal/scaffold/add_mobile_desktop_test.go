package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#125 was reported through a workaround: there is no way to add the Expo
// app to a project that exists, so the way in is to copy apps/expo from
// another project, and a copied directory has no manifest records. These
// commands write it properly, which is the difference between an upgrade that
// knows whose files those are and one that guesses.

func TestSetAppIncludedFlipsTheFlag(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "grit.json")
	const before = `{
  "name": "demo",
  "apps": {
    "expo": false,
    "desktop": false
  },
  "somethingThisVersionDoesNotKnow": 42
}`
	if err := os.WriteFile(path, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := setAppIncluded(root, "expo", true); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out := string(body)

	if !strings.Contains(out, `"expo": true`) {
		t.Error("the flag was not set, so grit upgrade would not maintain the app")
	}
	if !strings.Contains(out, `"desktop": false`) {
		t.Error("the other app's flag was changed")
	}
	// Edited as text on purpose: a round trip through a struct drops every key
	// this version of Grit does not know about, including a project's own.
	if !strings.Contains(out, "somethingThisVersionDoesNotKnow") {
		t.Error("a key Grit does not know about was dropped from grit.json")
	}

	// Setting it twice changes nothing and is not an error.
	if err := setAppIncluded(root, "expo", true); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(path)
	if string(again) != out {
		t.Error("setting the flag twice rewrote the file")
	}
}

// A project whose grit.json has no flag for the app is told so, rather than
// getting an app the upgrade will not maintain.
func TestSetAppIncludedRefusesWhenThereIsNoFlag(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grit.json"), []byte(`{"name":"demo"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	err := setAppIncluded(root, "expo", true)
	if err == nil {
		t.Fatal("no error for a grit.json with no expo flag")
	}
	if !strings.Contains(err.Error(), "maintain") {
		t.Errorf("the error does not say what the consequence is: %v", err)
	}
}

// Adding an app that is already there leaves it alone and says how to insist.
func TestAddAppLeavesAnExistingAppAlone(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grit.json"),
		[]byte(`{"name":"demo","apps":{"expo":false}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "apps", "expo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	const mine = "// the project's own\n"
	kept := filepath.Join(dir, "app.json")
	if err := os.WriteFile(kept, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	called := false
	err := addApp(root, false, appAddition{
		Key: "expo", Name: "Expo app", Dir: filepath.Join("apps", "expo"),
		Write: func(string, Options) error { called = true; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Error("it wrote over an app that was already there without --force")
	}
	body, _ := os.ReadFile(kept)
	if string(body) != mine {
		t.Error("the existing file was changed")
	}
}

// Both apps are self-contained, so both work in an --api project, which has no
// root package.json and no pnpm workspace at all. Telling somebody there to
// run pnpm install at the root sends them to a file that is not there.
func TestInstallStepDependsOnTheProject(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := installStep(workspace, filepath.Join("apps", "expo")); got != "pnpm install" {
		t.Errorf("in a monorepo: got %q", got)
	}

	apiOnly := t.TempDir()
	want := "cd apps/expo && pnpm install"
	if got := installStep(apiOnly, filepath.Join("apps", "expo")); got != want {
		t.Errorf("in an --api project: got %q, want %q", got, want)
	}
}
