package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// countFiles returns how many files exist under dir.
func countFiles(t *testing.T, dir string) int {
	t.Helper()
	n := 0
	_ = filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			n++
		}
		return nil
	})
	return n
}

// Admin-only writers must write nothing when the project has no admin app.
//
// The account security screen renders inside the dashboard chrome and its route
// lives under the dashboard layout, so it only means anything in a triple-tier
// project. Without the guard, --single and --double grew an apps/admin/
// directory holding five files and nothing else: no package.json, no layout,
// nothing that builds. An orphan tree is worse than an absent feature, because
// it reads as something half-finished.
func TestAdminOnlyWritersSkipProjectsWithNoAdmin(t *testing.T) {
	for _, arch := range []Architecture{ArchSingle, ArchDouble, ArchAPI, ArchMobile} {
		for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
			root := t.TempDir()
			opts := Options{
				ProjectName:  "shop",
				Architecture: arch,
				Frontend:     frontend,
			}

			if err := writeAdminSecurityFiles(root, opts); err != nil {
				t.Fatalf("%s/%s security: %v", arch, frontend, err)
			}
			if err := writeAdminPasskeyFiles(root, opts); err != nil {
				t.Fatalf("%s/%s passkeys: %v", arch, frontend, err)
			}

			adminDir := filepath.Join(root, "apps", "admin")
			if n := countFiles(t, adminDir); n != 0 {
				t.Errorf("%s/%s: wrote %d file(s) into apps/admin, which this "+
					"architecture does not have", arch, frontend, n)
			}
		}
	}
}

// And they must still write in a triple-tier project, into the layout that
// project actually uses.
func TestAdminOnlyWritersFollowTheLayout(t *testing.T) {
	cases := []struct {
		frontend Frontend
		want     []string
		unwanted string
	}{
		{
			frontend: FrontendNext,
			want: []string{
				"apps/admin/components/security/passkeys.tsx",
				"apps/admin/components/security/recovery-contacts.tsx",
				"apps/admin/hooks/use-security.ts",
				"apps/admin/lib/webauthn.ts",
				"apps/admin/app/(dashboard)/account/security/page.tsx",
			},
			unwanted: "apps/admin/src",
		},
		{
			frontend: FrontendTanStack,
			want: []string{
				"apps/admin/src/components/security/passkeys.tsx",
				"apps/admin/src/components/security/recovery-contacts.tsx",
				"apps/admin/src/hooks/use-security.ts",
				"apps/admin/src/lib/webauthn.ts",
				"apps/admin/src/pages/account-security.tsx",
				// Routing is the file tree here, so the page needs a route or it
				// has no URL.
				"apps/admin/src/routes/_dashboard/account/security.tsx",
			},
			unwanted: "apps/admin/components",
		},
	}

	for _, tc := range cases {
		root := t.TempDir()
		opts := Options{
			ProjectName:  "shop",
			Architecture: ArchTriple,
			Frontend:     tc.frontend,
		}

		if err := writeAdminSecurityFiles(root, opts); err != nil {
			t.Fatalf("%s security: %v", tc.frontend, err)
		}
		if err := writeAdminPasskeyFiles(root, opts); err != nil {
			t.Fatalf("%s passkeys: %v", tc.frontend, err)
		}

		for _, rel := range tc.want {
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
				t.Errorf("%s: missing %s", tc.frontend, rel)
			}
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(tc.unwanted))); err == nil {
			t.Errorf("%s: wrote into %s, which is the other frontend's layout",
				tc.frontend, tc.unwanted)
		}
	}
}
