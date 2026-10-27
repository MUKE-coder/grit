package port

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// write makes a file and the directories above it.
func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func has(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// A component is useless without the files it imports, so the copy follows
// the graph rather than taking the one file asked for.
func TestCopiesWhatTheEntryReaches(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()

	write(t, admin, "components/tables/data-table.tsx", `
import { useT } from "@/lib/i18n";
import { ColumnHeader } from "./column-header";
import type { ColumnDefinition } from "@/lib/resource";
import Link from "next/link";
`)
	write(t, admin, "components/tables/column-header.tsx", `
import { cn } from "@/lib/utils";
`)
	write(t, admin, "lib/i18n.tsx", "export const useT = () => {};")
	write(t, admin, "lib/resource.ts", "export type ColumnDefinition = {};")
	write(t, admin, "lib/utils.ts", "export const cn = () => {};")
	// Not reachable from the entry, so not copied.
	write(t, admin, "components/tables/import-modal.tsx", "export const X = 1;")
	write(t, target, "package.json", `{"dependencies":{"next":"15.0.0"}}`)

	res, err := Run(Options{
		AdminRoot: admin, TargetRoot: target,
		Entries: []string{"components/tables/data-table.tsx"},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"components/tables/data-table.tsx",
		"components/tables/column-header.tsx",
		"lib/i18n.tsx", // found despite the entry importing it without an extension
		"lib/resource.ts",
		"lib/utils.ts", // two levels down, through column-header
	} {
		if !has(res.Copied, want) {
			t.Errorf("did not copy %s; copied %v", want, res.Copied)
		}
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(want))); err != nil {
			t.Errorf("%s was reported copied but is not on disk", want)
		}
	}
	if has(res.Copied, "components/tables/import-modal.tsx") {
		t.Error("copied a file nothing imported")
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("unresolved: %v", res.Unresolved)
	}
}

// A package is not a file to copy, and only the ones the target does not
// already declare are worth a line in the output.
func TestReportsOnlyTheUndeclaredPackages(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()
	write(t, admin, "a.tsx", `
import Link from "next/link";
import { parsePhoneNumber } from "libphonenumber-js";
import { useQuery } from "@tanstack/react-query";
`)
	write(t, target, "package.json", `{
  "name": "web",
  "main": "index.js",
  "dependencies": { "next": "15.0.0", "@tanstack/react-query": "^5.0.0" }
}`)

	res, err := Run(Options{AdminRoot: admin, TargetRoot: target, Entries: []string{"a.tsx"}})
	if err != nil {
		t.Fatal(err)
	}
	if !has(res.External, "libphonenumber-js") || !has(res.External, "@tanstack/react-query") {
		t.Errorf("external packages wrong: %v", res.External)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "libphonenumber-js" {
		t.Errorf("missing should be only libphonenumber-js, got %v", res.Missing)
	}
	// "name" and "main" are keys in package.json too, and neither is a
	// dependency. A reader that treats them as declared packages would hide a
	// genuinely missing one called "main".
	if has(res.Missing, "name") || has(res.External, "name") {
		t.Error("a manifest field was read as a package")
	}
}

// An import that cannot be followed is said out loud. Dropping it silently
// produces a copy that does not compile and gives no clue why.
func TestUnfollowableImportsAreReported(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()
	write(t, admin, "a.tsx", `import { gone } from "@/lib/not-here";`)
	write(t, target, "package.json", "{}")

	res, err := Run(Options{AdminRoot: admin, TargetRoot: target, Entries: []string{"a.tsx"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unresolved) != 1 || !strings.Contains(res.Unresolved[0], "@/lib/not-here") {
		t.Errorf("expected the missing import to be reported, got %v", res.Unresolved)
	}
}

// Somebody has edited the copy by now. Replacing it silently is how that work
// disappears.
func TestExistingFilesAreKeptUnlessForced(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()
	write(t, admin, "a.tsx", "export const A = 2;")
	write(t, target, "a.tsx", "export const A = 1; // edited here")
	write(t, target, "package.json", "{}")

	res, err := Run(Options{AdminRoot: admin, TargetRoot: target, Entries: []string{"a.tsx"}})
	if err != nil {
		t.Fatal(err)
	}
	if !has(res.Skipped, "a.tsx") || len(res.Copied) != 0 {
		t.Errorf("an existing file was not skipped: copied %v skipped %v", res.Copied, res.Skipped)
	}
	data, _ := os.ReadFile(filepath.Join(target, "a.tsx"))
	if !strings.Contains(string(data), "edited here") {
		t.Fatal("the edited file was overwritten")
	}

	res, err = Run(Options{AdminRoot: admin, TargetRoot: target,
		Entries: []string{"a.tsx"}, Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if !has(res.Copied, "a.tsx") {
		t.Error("--force did not replace the file")
	}
}

// A dry run writes nothing, which is the only thing that makes it a dry run.
func TestDryRunWritesNothing(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()
	write(t, admin, "a.tsx", "export const A = 1;")
	write(t, target, "package.json", "{}")

	res, err := Run(Options{AdminRoot: admin, TargetRoot: target,
		Entries: []string{"a.tsx"}, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Copied) != 1 {
		t.Fatalf("a dry run should still report what it would copy, got %v", res.Copied)
	}
	if _, err := os.Stat(filepath.Join(target, "a.tsx")); err == nil {
		t.Error("the dry run wrote the file")
	}
}

// On a double or a single the admin lives inside the app you would be copying
// into, and copying a file onto itself is not a useful thing to do quietly.
func TestRefusesToCopyIntoItself(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.tsx", "export const A = 1;")
	_, err := Run(Options{AdminRoot: dir, TargetRoot: dir, Entries: []string{"a.tsx"}})
	if err == nil {
		t.Fatal("copying an app into itself was allowed")
	}
	if !strings.Contains(err.Error(), "already inside") {
		t.Errorf("the error does not explain why: %v", err)
	}
}

// A directory entry takes the source files in it, and not the tests beside
// them.
func TestADirectoryEntryTakesItsSources(t *testing.T) {
	admin, target := t.TempDir(), t.TempDir()
	write(t, admin, "components/ui/button.tsx", "export const B = 1;")
	write(t, admin, "components/ui/input.tsx", "export const I = 1;")
	write(t, admin, "components/ui/button.test.tsx", "test")
	write(t, admin, "components/ui/types.d.ts", "declare const x: number;")
	write(t, target, "package.json", "{}")

	res, err := Run(Options{AdminRoot: admin, TargetRoot: target,
		Entries: []string{"components/ui"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected the two components, got %v", res.Copied)
	}
	for _, unwanted := range []string{"components/ui/button.test.tsx", "components/ui/types.d.ts"} {
		if has(res.Copied, unwanted) {
			t.Errorf("copied %s", unwanted)
		}
	}
}
