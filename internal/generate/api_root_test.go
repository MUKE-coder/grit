package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// Where a project's Go module is, answered once.
//
// Six places worked this out independently and four were wrong on --single,
// which keeps its module in api/ rather than at the project root:
//
//	grit sync               looked in apps/api and refused to run at all
//	grit generate migration wrote to <root>/migrations, which grit migrate
//	                        does not read: the pair could never run
//	grit generate job       wrote to <root>/internal/jobs
//	grit generate mail      the same
//	grit generate saga      the same
//	grit generate sequence  the same
//
// Generator.APIRoot was the only one that had it right, and it is now the only
// implementation. These tests are over the resolver rather than each command,
// because one wrong answer is what produced six wrong commands.
func writeArch(t *testing.T, root, arch string, withAPIModule bool) {
	t.Helper()
	if withAPIModule {
		api := filepath.Join(root, "api")
		if err := os.MkdirAll(api, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(api, "go.mod"), []byte("module app/api\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "grit.json"),
		[]byte(`{"architecture":"`+arch+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestTheAPIRootIsResolvedPerArchitecture(t *testing.T) {
	t.Run("single keeps its module in api/", func(t *testing.T) {
		root := t.TempDir()
		writeArch(t, root, "single", true)

		want := filepath.Join(root, "api")
		if got := apiRootFrom(root); got != want {
			t.Errorf("apiRootFrom = %s, want %s", got, want)
		}
		// And the method agrees, because it is the same code.
		g := &Generator{Root: root, Architecture: "single"}
		if got := g.APIRoot(); got != want {
			t.Errorf("APIRoot = %s, want %s", got, want)
		}
	})

	t.Run("an older single keeps everything at the root", func(t *testing.T) {
		// The directory on disk decides rather than the version that wrote the
		// project: no api/go.mod means the Go code is at the root.
		root := t.TempDir()
		writeArch(t, root, "single", false)

		if got := apiRootFrom(root); got != root {
			t.Errorf("apiRootFrom = %s, want the project root", got)
		}
	})

	t.Run("a monorepo is apps/api", func(t *testing.T) {
		for _, arch := range []string{"triple", "double", "api", "mobile"} {
			root := t.TempDir()
			writeArch(t, root, arch, false)

			want := filepath.Join(root, "apps", "api")
			if got := apiRootFrom(root); got != want {
				t.Errorf("%s: apiRootFrom = %s, want %s", arch, got, want)
			}
		}
	})
}

// The printed path has to describe the root that was chosen.
//
// migration and sequence carried a display prefix beside the root, and the two
// were decided separately: a single project got an empty prefix while its
// files went to api/, so every line printed named a file one directory up from
// where it was written.
func TestTheDisplayPrefixDescribesTheRootItGot(t *testing.T) {
	root := filepath.Join("x", "project")

	if got := displayPrefix(root, root); got != "" {
		t.Errorf("a root that is the project root rendered as %q, want empty", got)
	}
	if got := displayPrefix(root, filepath.Join(root, "api")); got != "api/" {
		t.Errorf("single rendered as %q, want api/", got)
	}
	if got := displayPrefix(root, filepath.Join(root, "apps", "api")); got != "apps/api/" {
		t.Errorf("a monorepo rendered as %q, want apps/api/", got)
	}
}

// The admin panel, which had the same shape of bug in a different helper.
//
// adminRootFrom knew four layouts and not the one `grit new --single` has
// written since the SPA moved to the project root. Three commands read it, so
// on every single project `grit remove resource` left the admin definition and
// its page behind, still registered and pointing at an endpoint that no longer
// existed, and `grit sync` never added a new field to the admin table or form.
func TestTheAdminRootKnowsTheSingleLayout(t *testing.T) {
	cases := []struct{ dir, name string }{
		{filepath.Join("apps", "admin"), "triple"},
		{filepath.Join("apps", "web", "admin-panel"), "double with Next"},
		{filepath.Join("apps", "web", "src", "admin-panel"), "double with Vite"},
		{filepath.Join("src", "admin-panel"), "single"},
		{filepath.Join("frontend", "src", "admin-panel"), "an older single"},
	}
	for _, tc := range cases {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, tc.dir), 0o755); err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(root, tc.dir)
		if got := adminRootFrom(root); got != want {
			t.Errorf("%s: adminRootFrom = %s, want %s", tc.name, got, want)
		}
	}
}
