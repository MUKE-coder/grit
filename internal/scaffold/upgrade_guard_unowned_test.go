package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// grit#125: the upgrade overwrote files Grit had no record of writing, and
// wrote back files the project had deleted. "No record" and "mine to
// overwrite" were the same answer, and an Expo app copied in from another
// project (the only way in, since grit add has no expo) has no records at
// all. The reporter lost a tab bar, an API client, an auth provider and a
// query client, and the upgrade said "Updated 603 files".
//
// The discriminator is whether the PROJECT has a manifest, not whether the
// file has an entry, because a project generated before manifests existed has
// no entries either and freezing all of it would be worse than the problem.

func guardProject(t *testing.T, withManifest bool) string {
	t.Helper()
	root := t.TempDir()
	if withManifest {
		m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
		if err := m.Save(root, "3.443.0"); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func guardFile(t *testing.T, root, rel, body string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// The case that destroyed work: a file in a project that HAS a manifest, with
// no entry of its own, holding something other than the template.
func TestAFileGritNeverWroteIsNotOverwritten(t *testing.T) {
	root := guardProject(t, true)
	path := guardFile(t, root, "apps/expo/app/(tabs)/_layout.tsx",
		"// MY EDIT: this file is the project's own\n")

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// the template\n")
	_, _, unowned, _ := stopGuardFull()

	if allowed {
		t.Error("the upgrade overwrote a file Grit has no record of writing, which " +
			"is what destroyed the reporter's Expo app (grit#125)")
	}
	if len(unowned) != 1 || unowned[0].Rel != "apps/expo/app/(tabs)/_layout.tsx" {
		t.Errorf("it was not reported either, which was the other half of the " +
			"complaint: the upgrade said nothing at all")
	}
}

// And the reason the Untracked case could not simply be protected: a project
// from before manifests existed has no records for anything, so protecting
// every unrecorded file would refuse to update the whole project.
func TestAProjectWithNoManifestStillUpgrades(t *testing.T) {
	root := guardProject(t, false)
	path := guardFile(t, root, "apps/web/app/page.tsx", "// whatever is here\n")

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// the template\n")
	written, _, unowned, _ := stopGuardFull()

	if !allowed {
		t.Error("a pre-manifest project was refused an upgrade. With no provenance " +
			"anywhere there is no basis for claiming a file was edited, and " +
			"freezing the whole project is worse than the problem")
	}
	if written != 1 || len(unowned) != 0 {
		t.Errorf("written=%d unowned=%d", written, len(unowned))
	}
}

// A file with no entry that already matches the template byte for byte has
// nothing to lose, so it is written and adopted rather than reported. Without
// this, every file of a copied-in app that happens to be unmodified would be
// listed, and the report would be useless.
func TestAnUnrecordedFileThatMatchesIsAdopted(t *testing.T) {
	root := guardProject(t, true)
	const body = "// the template\n"
	path := guardFile(t, root, "apps/expo/app/index.tsx", body)

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, body)
	_, _, unowned, _ := stopGuardFull()

	if !allowed || len(unowned) != 0 {
		t.Errorf("an identical file was reported as a conflict: allowed=%v unowned=%d",
			allowed, len(unowned))
	}
}

// A whole directory the project removed stays removed. Six deleted Expo
// screens returned on every upgrade because "missing" and "never written"
// were the same answer (grit#125).
//
// The directory is what distinguishes this from the case below. Deleting one
// framework file out of a directory full of them is how people ask for the
// stock version back, and has been the documented way since v3.147.0;
// removing the directory is not that. Both behaviours are wanted, and this is
// the only thing that tells them apart.
func TestADeletedDirectoryIsNotWrittenBack(t *testing.T) {
	root := guardProject(t, true)

	path := guardFile(t, root, "apps/expo/app/blogs/index.tsx", "// the template\n")
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	m.Record("apps/expo/app/blogs/index.tsx", "scaffold", "3.444.0", "// the template\n")
	if err := m.Save(root, "3.444.0"); err != nil {
		t.Fatal(err)
	}
	// The whole screen, the way the reporter removed it.
	if err := os.RemoveAll(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// the template\n")
	_, _, _, deleted := stopGuardFull()

	if allowed {
		t.Error("a screen the project deleted was written back (grit#125)")
	}
	if len(deleted) != 1 {
		t.Errorf("it was not reported: %v", deleted)
	}

	// --force is the way back, so nothing is unrecoverable.
	if err := startGuard(root, true); err != nil {
		t.Fatal(err)
	}
	forced := guardAllows(path, "// the template\n")
	stopGuardFull()
	if !forced {
		t.Error("--force no longer restores a deleted file, so there is no way back")
	}
}

// And one file deleted out of a directory that is still there comes back,
// which is the idiom for resetting a framework component. TestDeletedFileIsRestored
// in upgrade_guard_test.go is the same assertion from the other side; this one
// exists so the pair sits together with the reason.
func TestADeletedFileInALiveDirectoryStillComesBack(t *testing.T) {
	root := guardProject(t, true)

	guardFile(t, root, "apps/admin/components/ui/badge.tsx", "// stock\n")
	path := guardFile(t, root, "apps/admin/components/ui/card.tsx", "// stock\n")
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	m.Record("apps/admin/components/ui/card.tsx", "scaffold", "3.444.0", "// stock\n")
	if err := m.Save(root, "3.444.0"); err != nil {
		t.Fatal(err)
	}
	// Just the one file. badge.tsx keeps the directory alive.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// stock v2\n")
	_, _, _, deleted := stopGuardFull()

	if !allowed {
		t.Error("deleting one framework file no longer brings the stock one back. " +
			"That has been how you reset a component since v3.147.0, and the " +
			"grit#125 fix must not take it away")
	}
	if len(deleted) != 0 {
		t.Errorf("it was reported as deliberately deleted: %v", deleted)
	}
}

// An edited file Grit did write keeps the behaviour it has had since v3.147.0,
// under the heading that says you edited it.
func TestAnEditedFileStillReportsAsYours(t *testing.T) {
	root := guardProject(t, true)
	path := guardFile(t, root, "apps/web/app/page.tsx", "// mine now\n")
	m, err := manifest.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	m.Record("apps/web/app/page.tsx", "scaffold", "3.443.0", "// the old template\n")
	if err := m.Save(root, "3.443.0"); err != nil {
		t.Fatal(err)
	}

	if err := startGuard(root, false); err != nil {
		t.Fatal(err)
	}
	allowed := guardAllows(path, "// the new template\n")
	_, skipped, unowned, _ := stopGuardFull()

	if allowed {
		t.Error("an edited file was overwritten")
	}
	if len(skipped) != 1 || len(unowned) != 0 {
		t.Errorf("it landed under the wrong heading: skipped=%d unowned=%d",
			len(skipped), len(unowned))
	}
}
