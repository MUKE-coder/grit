package main

import (
	"os"
	"path/filepath"
	"testing"
)

// grit#62: everything --multitenant needs has existed since the plugin did.
// What was missing is the flag, so the way in was to scaffold and then
// remember a second command.
//
// The flag installs through plugin.Install, the same call `grit plugin add
// multitenant` makes, the way --i18n goes through scaffold.AddI18n. That is
// the property worth a test: two code paths that write the same project are
// how they start writing different ones.

// pluginContextAt must not need the working directory, because `grit new`
// installs into the project it has just written while the process is still
// sitting in the parent.
func TestPluginContextAtDoesNotWalkUpFromTheWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grit.json"),
		[]byte(`{"architecture":"triple","frontend":"next"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, err := pluginContextAt(root)
	if err != nil {
		t.Fatalf("a project root that is not the working directory: %v", err)
	}
	if ctx.Root != root {
		t.Errorf("Root is %q, want %q", ctx.Root, root)
	}
	if ctx.Architecture != "triple" || ctx.Frontend != "next" {
		t.Errorf("the shape was not read from grit.json: %+v", ctx)
	}
}

// And it reads the shape rather than assuming one, because a plugin writes
// admin pages a double project has nowhere to put.
func TestPluginContextAtReadsTheShape(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "grit.json"),
		[]byte(`{"architecture":"api","frontend":"vite"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, err := pluginContextAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Architecture != "api" {
		t.Errorf("architecture is %q, want api", ctx.Architecture)
	}
	if ctx.Frontend != "vite" {
		t.Errorf("frontend is %q, want vite", ctx.Frontend)
	}
}

// A project with no grit.json still gets a usable context: a hand-made one
// should not fail to install a plugin because it lacks a file Grit wrote.
func TestPluginContextAtFallsBackWithoutGritJSON(t *testing.T) {
	ctx, err := pluginContextAt(t.TempDir())
	if err != nil {
		t.Fatalf("a project without grit.json: %v", err)
	}
	if ctx.Architecture == "" || ctx.Frontend == "" {
		t.Errorf("no defaults were applied: %+v", ctx)
	}
}
