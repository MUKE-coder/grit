package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The plugin must be reachable by the name the docs and NextSteps use.
func TestDevicePairingIsRegistered(t *testing.T) {
	p, err := Get("device-pairing")
	if err != nil {
		t.Fatalf("device-pairing is not registered: %v", err)
	}
	if p.Summary == "" || p.Description == "" {
		t.Error("a plugin with no summary or description is invisible in " +
			"`grit plugin list` and `grit plugin info`")
	}
}

// Every injection has to name a marker the scaffold actually emits.
//
// A marker typo is not a compile error and not a test failure anywhere else:
// install fails at the user's machine, halfway through, with files already
// written. Checking the marker strings against the scaffold is the only place
// this gets caught before then.
func TestDevicePairingMarkersExist(t *testing.T) {
	known := map[string]bool{
		"// grit:models":           true,
		"// grit:handlers":         true,
		"// grit:routes:custom":    true,
		"// grit:routes:protected": true,
		"// grit:routes:public":    true,
		"// grit:routes:admin":     true,
		"// grit:nav:system":       true,
		"// grit:icons:import":     true,
		"// grit:icons:map":        true,
	}
	ctx := Context{Root: ".", Module: "demo/apps/api", Architecture: "triple", Frontend: "next"}
	for _, inj := range devicePairingInjections(ctx) {
		if !known[inj.Marker] {
			t.Errorf("injection into %s uses marker %q, which the scaffold does "+
				"not emit; install would fail half-applied", inj.File, inj.Marker)
		}
	}
}

// The anonymous half has to go somewhere without auth middleware.
//
// The browser being paired has no session and no API key, so neither the
// protected group nor the public group (which wants a key) can serve it. Only
// grit:routes:custom sits at function level with v1 in scope.
func TestPairingAnonymousRoutesAreNotBehindAuth(t *testing.T) {
	ctx := Context{Root: ".", Module: "demo/apps/api", Architecture: "triple", Frontend: "next"}

	var anonymous string
	for _, inj := range devicePairingInjections(ctx) {
		if inj.Marker == "// grit:routes:custom" {
			anonymous = inj.Code
		}
		if inj.Marker == "// grit:routes:protected" {
			if strings.Contains(inj.Code, "/pair/start") || strings.Contains(inj.Code, `GET("/pair/:code"`) {
				t.Error("the browser-facing pairing routes are on the protected " +
					"group, where an anonymous browser can never reach them")
			}
		}
	}
	for _, want := range []string{`v1.POST("/pair/start"`, `v1.GET("/pair/:code"`} {
		if !strings.Contains(anonymous, want) {
			t.Errorf("missing %s on the unauthenticated group", want)
		}
	}
}

// A nav entry missing adminOnly does not compile the admin.
//
// NavEntry requires it, so leaving it off turns a working install into a
// TypeScript error in the user's project, which is the worst place to find it.
func TestPairingNavEntryIsWellFormed(t *testing.T) {
	ctx := Context{Root: ".", Module: "demo/apps/api", Architecture: "triple", Frontend: "next"}
	for _, inj := range devicePairingInjections(ctx) {
		if inj.Marker != "// grit:nav:system" {
			continue
		}
		for _, want := range []string{"href:", "label:", "iconKey:", "adminOnly:"} {
			if !strings.Contains(inj.Code, want) {
				t.Errorf("the sidebar entry has no %s: NavEntry requires it and "+
					"the admin will not typecheck without it", want)
			}
		}
		return
	}
	t.Error("no sidebar entry is injected")
}

// An --api project has no web or admin app to write into.
func TestPairingWritesOnlyWhatTheProjectHas(t *testing.T) {
	for _, tc := range []struct {
		arch      string
		wantWeb   bool
		wantAdmin bool
	}{
		{"api", false, false},
		{"single", false, false},
		{"double", true, false},
		{"triple", true, true},
	} {
		ctx := Context{Root: ".", Module: "demo/apps/api", Architecture: tc.arch, Frontend: "next"}
		files := devicePairingFiles(ctx)

		var sawWeb, sawAdmin bool
		for path := range files {
			if strings.HasPrefix(path, "apps/web/") {
				sawWeb = true
			}
			if strings.HasPrefix(path, "apps/admin/") {
				sawAdmin = true
			}
		}
		if sawWeb != tc.wantWeb {
			t.Errorf("%s: web files present = %v, want %v", tc.arch, sawWeb, tc.wantWeb)
		}
		if sawAdmin != tc.wantAdmin {
			t.Errorf("%s: admin files present = %v, want %v", tc.arch, sawAdmin, tc.wantAdmin)
		}
		// The API half is unconditional, and --single puts it at the root.
		wantModel := "apps/api/internal/models/pairing_request.go"
		if tc.arch == "single" {
			wantModel = "internal/models/pairing_request.go"
		}
		if _, ok := files[wantModel]; !ok {
			t.Errorf("%s: no model at %s", tc.arch, wantModel)
		}
	}
}

// Two injections carrying identical code at different markers must both apply.
//
// The idempotency guard used to be "does this text appear anywhere in the
// file", which is a different question. An icon needs the same token in the
// import list and in the map, and the second injection was silently dropped:
// the nav entry then rendered getIcon's FileText fallback, with nothing
// anywhere reporting a problem. Found by installing this plugin and looking at
// the file.
func TestInjectBeforeAppliesTheSameCodeAtTwoMarkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "icons.ts")
	src := strings.Join([]string{
		"import {",
		"  Users,",
		"  // grit:icons:import",
		"} from \"lucide-react\";",
		"",
		"export const iconMap = {",
		"  Users,",
		"  // grit:icons:map",
		"};",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	for _, marker := range []string{"// grit:icons:import", "// grit:icons:map"} {
		if err := injectBefore(path, marker, "  QrCode,"); err != nil {
			t.Fatalf("%s: %v", marker, err)
		}
	}

	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(out), "  QrCode,"); n != 2 {
		t.Errorf("QrCode appears %d time(s), want 2: an icon needs both the "+
			"import and the map entry, and half of it renders the wrong icon "+
			"with no error\n%s", n, out)
	}
}

// Applying the same injection twice at one marker is still a no-op.
func TestInjectBeforeStaysIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "routes.go")
	if err := os.WriteFile(path, []byte("func Setup() {\n\t// grit:routes\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := injectBefore(path, "// grit:routes", "\tr.GET(\"/x\", h.X)"); err != nil {
			t.Fatalf("pass %d: %v", i, err)
		}
	}
	out, _ := os.ReadFile(path)
	if n := strings.Count(string(out), "r.GET"); n != 1 {
		t.Errorf("the route was injected %d times; installing twice must be a no-op", n)
	}
}
