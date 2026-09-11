package scaffold

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestVersionLess(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		less bool
	}{
		{"v2.2.1", "v2.5.0", true},
		{"v2.5.0", "v2.5.0", false},
		{"v2.6.0", "v2.5.0", false},
		{"v1.0.1", "v1.1.0", true},
		// Numbers, not strings: v1.10.0 is the newer one.
		{"v1.10.0", "v1.9.0", false},
		{"v1.9.0", "v1.10.0", true},
		// A pseudo-version is a pre-release of v0.0.0.
		{"v0.0.0-20260529025319-478cdfa8ce5f", "v1.0.0", true},
		{"v1.0.0-rc.1", "v1.0.0", true},
		{"v1.0.0", "v1.0.0-rc.1", false},
		{"v2.0.0+incompatible", "v2.0.0", false},
	} {
		if got := versionLess(tc.a, tc.b); got != tc.less {
			t.Errorf("versionLess(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.less)
		}
	}
}

// A project scaffolded on an older Sentinel stayed on it: grit upgrade rewrote
// framework files and never go.mod, and v2.2.2 is a security release. It is
// raised now, and a library already past its floor is left where it is.
func TestUpgradeRaisesFrameworkDepsAndNeverLowersThem(t *testing.T) {
	dir := t.TempDir()
	write := func(body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var fetched []string
	calls := 0
	orig := goGet
	goGet = func(_ string, specs ...string) error {
		calls++
		fetched = specs
		return nil
	}
	defer func() { goGet = orig }()

	write(`module shop/apps/api

go 1.24

require (
	github.com/MUKE-coder/gorm-studio v1.0.1
	github.com/MUKE-coder/pulse v0.0.0-20260529025319-478cdfa8ce5f
	github.com/MUKE-coder/sentinel/v2 v2.9.0
	github.com/gin-gonic/gin v1.11.0
)
`)
	raised, err := raiseFrameworkDeps(dir)
	if err != nil {
		t.Fatalf("raise: %v", err)
	}
	want := []string{"github.com/MUKE-coder/gorm-studio@v1.1.0", "github.com/MUKE-coder/pulse@v1.0.0"}
	if !reflect.DeepEqual(fetched, want) {
		t.Errorf("fetched %v, want %v (and never a lower Sentinel than the project has)", fetched, want)
	}
	if len(raised) != 2 {
		t.Errorf("reported %d raises, want 2: %v", len(raised), raised)
	}

	// Everything at or past its floor: go get is not run at all.
	calls = 0
	write(`module shop/apps/api

go 1.24

require (
	github.com/MUKE-coder/gorm-studio v1.1.0
	github.com/MUKE-coder/pulse v1.0.0
	github.com/MUKE-coder/sentinel/v2 v2.5.0
)
`)
	if raised, err := raiseFrameworkDeps(dir); err != nil || raised != nil || calls != 0 {
		t.Errorf("an up-to-date project ran go get %d time(s): %v %v", calls, raised, err)
	}

	// No go.mod: nothing to do, and no error.
	if raised, err := raiseFrameworkDeps(t.TempDir()); err != nil || raised != nil {
		t.Errorf("a directory with no go.mod: %v %v", raised, err)
	}
}
