package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// An upgrade must leave a route file the project deleted deleted.
//
// A shop replaces the scaffold's (marketing) group with its own, because the
// landing page and the shop homepage are both "/" and only one of them can be.
// Writing (marketing)/page.tsx back does not add an unused page: Next.js
// refuses to build two route groups that resolve to the same path, so the
// upgrade leaves the project unable to start.
func TestUpgradeDoesNotResurrectADeletedPage(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "shop"}
	landing := filepath.Join(root, "apps", "web", "app", "(marketing)", "page.tsx")
	utils := filepath.Join(root, "apps", "web", "lib", "utils.ts")

	// A project Grit wrote, recorded as Grit wrote it.
	release, err := manifest.Start(root, "3.300.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeWebFiles(root, opts); err != nil {
		t.Fatalf("the first write failed: %v", err)
	}
	for path, body := range webFileMap(root, opts) {
		manifest.Note(path, adminHref(body, opts))
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(landing); err != nil {
		t.Fatalf("the scaffold did not write the landing page: %v", err)
	}

	// The shop deletes the group and takes over "/" with its own.
	if err := os.RemoveAll(filepath.Dir(landing)); err != nil {
		t.Fatal(err)
	}
	// And something that is not a page goes missing too, by accident.
	if err := os.Remove(utils); err != nil {
		t.Fatal(err)
	}

	release, err = manifest.Start(root, "3.301.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeWebFiles(root, opts); err != nil {
		t.Fatalf("the upgrade failed: %v", err)
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(landing); err == nil {
		t.Error("the upgrade put the deleted landing page back, so a project that " +
			"replaced the (marketing) group now has two route groups resolving to /")
	}

	// The narrowness is the point: a deleted page is a routing decision, a
	// deleted lib file is damage, and repairing damage is what an upgrade does.
	if _, err := os.Stat(utils); err != nil {
		t.Error("the upgrade did not restore lib/utils.ts, which is not a route file")
	}
}

// A project that never had the file, or one upgrading from before the manifest
// existed, is not a project that deleted it.
func TestUpgradeWritesAPageTheProjectNeverHad(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "shop"}
	landing := filepath.Join(root, "apps", "web", "app", "(marketing)", "page.tsx")

	// Nothing recorded: the manifest has no entry, so the status is Untracked
	// rather than Missing and the file is written.
	release, err := manifest.Start(root, "3.301.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeWebFiles(root, opts); err != nil {
		t.Fatalf("the write failed: %v", err)
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(landing); err != nil {
		t.Errorf("a project with no manifest entry for the landing page did not get one: %v", err)
	}
}
