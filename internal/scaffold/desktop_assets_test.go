package scaffold

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Two things a generated desktop app did not have.
//
// It had no icon: build/appicon was created as an empty directory, and Wails
// reads build/appicon.png, a file, which is also what the README describes.
// And it had nothing for //go:embed all:frontend/dist to match, so `go build
// ./...` and `go vet ./...` failed in apps/desktop on a fresh scaffold. `wails
// build` worked, because it builds the frontend first; every other Go tool,
// and every editor, reported a broken package on a project nobody had touched.
//
// Found by generating a project with --desktop and running go build in it.

func desktopProject(t *testing.T) (string, Options) {
	t.Helper()
	root := t.TempDir()
	opts := Options{ProjectName: "app", Theme: "coral", Frontend: FrontendNext}
	if err := writeDesktopClientFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopClientFiles: %v", err)
	}
	return root, opts
}

func TestTheDesktopAppHasAnIcon(t *testing.T) {
	root, opts := desktopProject(t)
	path := filepath.Join(root, "apps", "desktop", "build", "appicon.png")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no build/appicon.png: %v. Wails builds every platform icon from "+
			"that file, and the README says it is there.", err)
	}
	want, err := appIconPNG(opts.Theme)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, want) {
		t.Error("the icon is not the one drawn for this project's theme")
	}
	if bytes.Equal(data, gritLogoPNG) {
		t.Error("the icon is the Grit logo, so this app ships somebody else's brand")
	}

	readme := filepath.Join(root, "apps", "desktop", "build", "README.md")
	body, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("nothing beside the icon says what it is: %v", err)
	}
	if !strings.Contains(string(body), "placeholder") {
		t.Error("the note beside the icon does not say it is a placeholder")
	}
}

// The embed has something to match before anybody builds the frontend, and
// the thing it matches survives a clone.
func TestTheDesktopEmbedHasSomethingToMatch(t *testing.T) {
	root, _ := desktopProject(t)
	dist := filepath.Join(root, "apps", "desktop", "frontend", "dist")

	entries, err := os.ReadDir(dist)
	if err != nil {
		t.Fatalf("frontend/dist does not exist, so //go:embed all:frontend/dist "+
			"fails and this module does not build: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("frontend/dist is empty, so the embed matches nothing")
	}

	// A real build overwrites index.html, so index.html cannot be the file
	// that is committed.
	gitignore := readTestFileAt(t, filepath.Join(root, "apps", "desktop", ".gitignore"))
	if !strings.Contains(gitignore, "!frontend/dist/.gitkeep") {
		t.Error("the .gitignore excludes all of frontend/dist, so a fresh clone has " +
			"nothing for the embed to match and go build fails again")
	}
	if _, err := os.Stat(filepath.Join(dist, ".gitkeep")); err != nil {
		t.Errorf("the file the .gitignore keeps is not there: %v", err)
	}

	// And the page somebody sees if they run the binary unbuilt says so.
	page := readTestFileAt(t, filepath.Join(dist, "index.html"))
	if !strings.Contains(page, "frontend has not been built") {
		t.Error("the placeholder page does not explain itself, so an unbuilt binary " +
			"opens a blank window")
	}
}

// readTestFileAt reads a generated file or fails.
func readTestFileAt(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}

// The Wails module version and the CLI version the release workflow installs
// are one setting in two files. They drifted: go.mod said 2.9.2 while the
// workflow and anybody's local CLI were several minor versions ahead, and
// Wails printed a warning about it on every single build.
func TestTheWailsVersionsAgree(t *testing.T) {
	gomod := desktopClientGoMod("app/apps/desktop")
	workflow := releaseCIYAML(Options{ProjectName: "app", IncludeDesktop: true})

	version := ""
	for _, line := range strings.Split(gomod, "\n") {
		if strings.Contains(line, "github.com/wailsapp/wails/v2 v") {
			fields := strings.Fields(line)
			version = fields[len(fields)-1]
			break
		}
	}
	if version == "" {
		t.Fatal("the desktop go.mod pins no Wails version")
	}
	if !strings.Contains(workflow, "cmd/wails@"+version) {
		t.Errorf("go.mod pins Wails %s and the release workflow installs a different CLI. "+
			"Wails warns on every build when they differ.", version)
	}
}

// The desktop module is its own module, so the API's dependency floors do not
// reach it. Wails asks for an x/net below the one that carries the fix for
// GO-2026-6617 and three others.
func TestTheDesktopModulePinsAPatchedXNet(t *testing.T) {
	gomod := desktopClientGoMod("app/apps/desktop")
	if !strings.Contains(gomod, "golang.org/x/net v0.60.0") {
		t.Error("the desktop go.mod does not raise golang.org/x/net, so `go mod tidy` " +
			"takes the version Wails asks for, which is below the patched one")
	}
}

// The fifth place the API's address was written down, and the third that did
// not follow APP_PORT.
//
// api-client.ts reads import.meta.env.VITE_API_URL and fell back to the
// literal localhost:8080, and nothing set VITE_API_URL for the desktop app.
// The dev proxy named the same literal. So on a project that had moved
// APP_PORT the app called a port nothing was listening on, and when something
// else happened to be there it answered: sign-in came back 500 from a server
// belonging to a different project.
func TestTheDesktopClientFollowsAppPort(t *testing.T) {
	src := desktopClientViteConfig()

	if !strings.Contains(src, "viteEnv.APP_PORT") {
		t.Error("the desktop Vite config does not read APP_PORT, so the dev proxy cannot " +
			"follow the API when it moves")
	}
	if strings.Contains(src, `target: "http://localhost:8080"`) {
		t.Error("the dev proxy still names a literal port")
	}
	if !strings.Contains(src, "define: clientEnv") {
		t.Error("the config does not hand the resolved address to the browser bundle, so " +
			"the client and the proxy can disagree about where the API is")
	}
	if !strings.Contains(src, `"import.meta.env.VITE_API_URL": JSON.stringify(apiTarget`) {
		t.Error("the client's VITE_API_URL is not derived from the resolved target")
	}

	// envDir is the monorepo root: the app runs from apps/desktop/frontend and
	// there is no .env there, so reading from cwd returns undefined for
	// everything and every default silently wins.
	if !strings.Contains(src, `path.resolve(__dirname, "../../..")`) {
		t.Error("loadEnv does not point at the monorepo root, so it reads no .env at all")
	}
}

// And the route tree the type-check needs is generated by the check itself.
//
// routeTree.gen.ts is written by the router plugin during dev or build and is
// gitignored, so `pnpm run type-check` on a project nobody has run reported an
// error in every route file. Twenty errors that are all one missing file is
// worse than no check, because it buries the real one: the desktop profile
// route was rendering a component it never imported, and the only reason that
// was visible was that wails build runs vite before tsc.
func TestEveryTanStackTypeCheckGeneratesItsRouteTree(t *testing.T) {
	viteOpts := Options{ProjectName: "app", Frontend: FrontendTanStack}
	for name, src := range map[string]string{
		"apps/web":              webTanStackPackageJSON(viteOpts),
		"apps/admin":            adminTanStackPackageJSON(viteOpts),
		"apps/desktop/frontend": desktopClientPackageJSON(Options{ProjectName: "app"}),
	} {
		scripts := scriptsOf(t, name, src)
		check, ok := scripts["type-check"]
		if !ok {
			t.Errorf("%s declares no type-check script", name)
			continue
		}
		if !strings.Contains(check, "tsr generate") {
			t.Errorf("%s runs %q, which reads routeTree.gen.ts without generating it. On a "+
				"project nobody has built, that is an error in every route file.", name, check)
		}
		if !strings.Contains(src, "@tanstack/router-cli") {
			t.Errorf("%s has no @tanstack/router-cli, so `tsr` is not on PATH", name)
		}
	}
}
