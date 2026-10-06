package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// `grit new x --api` is offered on the home page as "the Go API alone".
//
// It grew an apps/admin/ holding seven .tsx files and nothing to build them
// with: no package.json, no tsconfig, no layout. Nothing failed, because
// nothing compiled them. An orphan tree that never builds and nobody asked for
// is worse than the feature being absent, because it reads as a scaffold that
// stopped halfway.
//
// writeAdminSecurityFiles already carried the HasAdminPanel guard, with a
// comment describing this exact failure. writeAdminAccountFiles was written
// later and did not, and every screen added to it since inherited the gap. So
// this checks the outcome rather than any one writer: scaffold an API project
// and assert there is no frontend in it.
func TestAPIOnlyProjectHasNoFrontendDirectories(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "apionly", Architecture: ArchAPI, Theme: "atlas"}

	// Every writer that has ever put a file under apps/admin or apps/web. Calling
	// them directly rather than running the whole scaffold keeps this a unit test
	// (no go mod tidy, no network) while still covering the writers that matter.
	for name, write := range map[string]func(string, Options) error{
		"account":  writeAdminAccountFiles,
		"security": writeAdminSecurityFiles,
		"passkey":  writeAdminPasskeyFiles,
	} {
		if err := write(root, opts); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}

	for _, dir := range []string{
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "frontend"),
	} {
		var found []string
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil //nolint:nilerr // a missing directory is the pass
			}
			rel, _ := filepath.Rel(root, path)
			found = append(found, filepath.ToSlash(rel))
			return nil
		})
		if len(found) > 0 {
			t.Errorf("an API-only project has %d frontend file(s) nothing can build: %s",
				len(found), strings.Join(found, ", "))
		}
	}
}

// And the guard belongs to the panel, not to the architecture.
//
// A single project has no apps/admin either, but it does have a panel, inside
// its frontend. A guard written as "is this a triple" would delete the panel
// from both singles and the double.
func TestEveryArchitectureWithAPanelStillGetsTheAccountScreens(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts Options
	}{
		{"triple", Options{ProjectName: "a", Architecture: ArchTriple, Frontend: FrontendNext, Theme: "atlas"}},
		{"double", Options{ProjectName: "a", Architecture: ArchDouble, Frontend: FrontendNext, Theme: "atlas"}},
		{"single vite", Options{ProjectName: "a", Architecture: ArchSingle, Theme: "atlas"}},
		{"single next", Options{ProjectName: "a", Architecture: ArchSingle, Frontend: FrontendNext, Theme: "atlas"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writeAdminAccountFiles(root, tc.opts); err != nil {
				t.Fatal(err)
			}
			var count int
			_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					count++
				}
				return nil
			})
			if count == 0 {
				t.Error("the account, trash and deleted-account screens were not written, so this project has a panel with no account page")
			}
		})
	}
}
