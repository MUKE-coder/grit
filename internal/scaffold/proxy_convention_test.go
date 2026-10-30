package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

const oldGate = `import { NextResponse, type NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  return NextResponse.next();
}

export const config = { matcher: ["/account/:path*"] };
`

// nextApp marks a directory as a Next app: the config file is what
// nextAppRoots looks for, because the directory's name does not say which
// framework is in it.
func nextApp(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "next.config.ts"), []byte("export default {};\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runRepair calls the repair the way upgrade does: inside an open recording.
// manifest.IsUnchanged answers false when nothing is recording, so without
// this every file reads as edited and the move never happens.
func runRepair(t *testing.T, root string) string {
	t.Helper()
	return captureStdout(t, func() {
		release, err := manifest.Start(root, "3.459.0", "upgrade")
		if err != nil {
			t.Fatal(err)
		}
		if err := repairProxyConvention(root, tripleOptions()); err != nil {
			t.Fatal(err)
		}
		if err := release(); err != nil {
			t.Fatal(err)
		}
	})
}

// Next 16 deprecated the filename. The file Grit wrote moves to the new one,
// with the export renamed, and the old file goes: Next reads only one of them,
// and leaving both behind is how somebody ends up editing the wrong one.
func TestGritsMiddlewareMovesToProxy(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	nextApp(t, web)
	old := record(t, root, "apps/web/middleware.ts", oldGate)

	captured := runRepair(t, root)

	if fileExists(old) {
		t.Error("middleware.ts is still there, so Next reads it and still warns")
	}
	body := read(t, filepath.Join(web, "proxy.ts"))
	if !strings.Contains(body, "export function proxy(request: NextRequest)") {
		t.Error("the export was not renamed, so Next finds the file and no function in it")
	}
	if strings.Contains(body, "export function middleware") {
		t.Error("the old export is still there")
	}
	// Otherwise the file they had: there are two shapes of it, the web-auth one
	// and the admin-only one, and the file on disk is the only thing that knows
	// which this project has.
	if !strings.Contains(body, `matcher: ["/account/:path*"]`) {
		t.Error("the matcher did not survive the move")
	}
	if !strings.Contains(captured, "proxy.ts") {
		t.Errorf("the move was not reported; it printed: %q", captured)
	}
}

// A developer's own protected paths move with the file. Renaming it preserves
// the content, so there is nothing to lose and a warning to be rid of.
func TestAnEditedMiddlewareMovesWithItsEdits(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	nextApp(t, web)
	old := record(t, root, "apps/web/middleware.ts", oldGate)

	mine := strings.Replace(oldGate, `"/account/:path*"`, `"/account/:path*", "/checkout"`, 1)
	if err := os.WriteFile(old, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	runRepair(t, root)

	if fileExists(old) {
		t.Error("the file did not move, so the deprecation warning stays")
	}
	body := read(t, filepath.Join(web, "proxy.ts"))
	if !strings.Contains(body, `"/checkout"`) {
		t.Error("their protected path did not survive the move")
	}
	if !strings.Contains(body, "export function proxy(request: NextRequest)") {
		t.Error("the export was not renamed")
	}
}

// A file whose export is not the one Grit wrote is one the rename would
// break: the new name would be read and there would be no function behind it.
// Those are left alone with the one line it takes to do it by hand.
func TestAMiddlewareWithAnotherExportIsLeftAlone(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	nextApp(t, web)
	theirs := `import { NextResponse, type NextRequest } from "next/server";

export default function (request: NextRequest) {
  return NextResponse.next();
}
`
	old := record(t, root, "apps/web/middleware.ts", theirs)

	captured := runRepair(t, root)

	if !fileExists(old) {
		t.Fatal("a file the rename could not safely touch was moved anyway")
	}
	if fileExists(filepath.Join(web, "proxy.ts")) {
		t.Error("a proxy.ts was written beside the file they still have, so the project holds two and Next ignores one")
	}
	for _, want := range []string{"proxy.ts", "middleware to proxy"} {
		if !strings.Contains(captured, want) {
			t.Errorf("the message does not mention %q; it printed: %q", want, captured)
		}
	}
}

// Running it twice changes nothing: the second pass finds no middleware.ts.
func TestTheProxyMoveIsIdempotent(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	nextApp(t, web)
	record(t, root, "apps/web/middleware.ts", oldGate)

	runRepair(t, root)
	runRepair(t, root)
	if !fileExists(filepath.Join(web, "proxy.ts")) {
		t.Error("the second pass removed the file the first one wrote")
	}
}

// Both present is somebody migrating by hand. Next reads one and ignores the
// other, which is worth saying out loud rather than resolving for them.
func TestBothFilesPresentIsReported(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "apps", "web")
	nextApp(t, web)
	old := record(t, root, "apps/web/middleware.ts", oldGate)
	if err := os.WriteFile(filepath.Join(web, "proxy.ts"), []byte(oldGate), 0o644); err != nil {
		t.Fatal(err)
	}

	captured := runRepair(t, root)
	if !fileExists(old) {
		t.Error("the file was removed without being read, while a proxy.ts of unknown provenance sat beside it")
	}
	if !strings.Contains(captured, "both exist") {
		t.Errorf("nothing was said about the two files; it printed: %q", captured)
	}
}

// A Vite app has no such file and must not be given one. The marker is the
// Next config, not the directory's name.
func TestOnlyNextAppsAreLookedAt(t *testing.T) {
	root := t.TempDir()
	vite := filepath.Join(root, "apps", "admin")
	if err := os.MkdirAll(vite, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vite, "vite.config.ts"), []byte("export default {};\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, dir := range nextAppRoots(root, tripleOptions()) {
		if dir == vite {
			t.Error("a Vite app was counted as a Next app")
		}
	}
}
