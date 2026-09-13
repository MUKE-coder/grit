package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoVersionLess(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"1.25.0", "1.26.6", true},
		{"1.26", "1.26.6", true},
		{"1.26.6", "1.26.6", false},
		{"1.27", "1.26.6", false},
		{"1.21", "1.26.6", true},
	} {
		if got := goVersionLess(c.a, c.b); got != c.want {
			t.Errorf("goVersionLess(%s, %s) = %v", c.a, c.b, got)
		}
	}
}

// The scaffold builds on the floor everywhere it names a Go version.
func TestTemplatesBuildOnTheGoFloor(t *testing.T) {
	opts := Options{ProjectName: "shop"}
	if got := goDirective(apiGoMod(opts)); got != goToolchainFloor {
		t.Errorf("the API go.mod says go %s, want %s", got, goToolchainFloor)
	}
	for name, src := range map[string]string{
		"Dockerfile (single)": dockerfileSingle(),
		"docker files":        dockerComposeProd(opts),
	} {
		if out, fixed, _ := repairGoVersionPinsSource(src); out != src || len(fixed) > 0 {
			t.Errorf("%s still pins a Go below %s", name, goToolchainFloor)
		}
	}
	for _, s := range []string{"golang.org/x/crypto v0.57.0", "filippo.io/edwards25519 v1.2.0"} {
		if !strings.Contains(apiGoMod(opts), s) {
			t.Errorf("the API go.mod is missing the floor %s", s)
		}
	}
}

func TestRepairGoVersionPins(t *testing.T) {
	src := "FROM golang:1.26-alpine AS builder\nFROM golang:1.27-alpine AS later\n" +
		"      - uses: actions/setup-go@v5\n        with:\n          go-version: '1.24'\n"
	out, fixed, _ := repairGoVersionPinsSource(src)
	for _, want := range []string{"golang:1.26.6-alpine AS builder", "golang:1.27-alpine AS later", "go-version: '1.26.6'"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if len(fixed) != 1 || !strings.Contains(fixed[0], "2 pins") {
		t.Errorf("reported %v, want 2 pins raised", fixed)
	}
}

func TestRepairGoToolchainRaisesTheDirective(t *testing.T) {
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")
	if err := os.MkdirAll(api, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(api, "go.mod"), []byte("module shop/apps/api\n\ngo 1.25.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var got []string
	saved := goModEdit
	goModEdit = func(dir string, args ...string) error { got = append([]string{dir}, args...); return nil }
	t.Cleanup(func() { goModEdit = saved })

	if err := repairGoToolchain(root, Options{ProjectName: "shop"}); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1] != "-go="+goToolchainFloor {
		t.Errorf("go mod edit called with %v", got)
	}
}
