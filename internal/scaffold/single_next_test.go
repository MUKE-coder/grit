package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func singleNextOptions() Options {
	return Options{ProjectName: "app", Architecture: ArchSingle, Frontend: FrontendNext, Theme: "atlas"}
}

// --single --next used to accept the flag and build a Vite app anyway.
//
// Architecture and Frontend were treated as one choice: ShouldEmbedAdminInSPA
// returned true for every single, so the panel was always written in the
// TanStack dialect into a Vite SPA, and the only trace of --next was grit.json
// claiming "next" beside a vite.config.ts.
//
// They are separate choices now, and these are the predicates that decide
// everything downstream: which admin shape is written, which dialect it is
// written in, and where it lands.
func TestSingleNextPicksTheNextShapeNotTheSPAShape(t *testing.T) {
	next := singleNextOptions()
	vite := Options{ProjectName: "app", Architecture: ArchSingle, Theme: "atlas"}

	if !next.SingleUsesNext() {
		t.Fatal("--single --next is not recognised as a Next single, so every branch below takes the Vite path")
	}
	if vite.SingleUsesNext() {
		t.Error("--single on its own was taken for a Next single")
	}

	// The panel is a Next route group, like a double's, not a section of an SPA.
	if !next.ShouldEmbedAdmin() {
		t.Error("the panel is not embedded as a route group, so a Next single has no admin at all")
	}
	if next.ShouldEmbedAdminInSPA() {
		t.Error("the panel is being written as an SPA section into a Next app")
	}
	if next.AdminIsTanStack() {
		t.Error("the panel would be written in the TanStack dialect, so every page imports a router the app does not have")
	}

	// And the Vite single is untouched by all of it.
	if !vite.ShouldEmbedAdminInSPA() || !vite.AdminIsTanStack() {
		t.Error("the Vite single stopped getting the SPA panel")
	}
}

// Both frontends live in frontend/, which is what makes a project "single".
func TestSingleKeepsItsFrontendBesideTheGoModule(t *testing.T) {
	root := filepath.FromSlash("/p")
	for _, tc := range []struct {
		name string
		opts Options
		want string
	}{
		{"single next", singleNextOptions(), filepath.Join(root, "frontend")},
		{"single vite", Options{ProjectName: "a", Architecture: ArchSingle}, filepath.Join(root, "frontend")},
		{"double", Options{ProjectName: "a", Architecture: ArchDouble, Frontend: FrontendNext}, filepath.Join(root, "apps", "web")},
		{"triple", Options{ProjectName: "a", Architecture: ArchTriple, Frontend: FrontendNext}, filepath.Join(root, "apps", "web")},
	} {
		if got := webAppRoot(root, tc.opts); got != tc.want {
			t.Errorf("%s: web app at %s, want %s", tc.name, got, tc.want)
		}
	}
}

// The panel's files land inside that frontend, and its imports are repointed.
//
// embeddedAdminContent was keyed on the literal path "/apps/web/", so a panel
// written into frontend/ kept the standalone app's "@/components/..." imports.
// Those resolve to the host app's own components directory, which does not have
// them: 173 module-not-found errors on the first build, every one a real file
// one directory away behind a different alias.
func TestSingleNextPanelLandsInTheFrontendWithItsImportsRepointed(t *testing.T) {
	root := t.TempDir()
	opts := singleNextOptions()
	files := embeddedAdminFileMap(root, opts)
	if len(files) < 100 {
		t.Fatalf("the panel has only %d files", len(files))
	}

	frontend := filepath.Join(root, "frontend")
	for path := range files {
		rel, err := filepath.Rel(frontend, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			t.Fatalf("%s is outside the frontend a single project has", path)
		}
	}

	// A page under app/admin, and the content rewrite that has to reach it.
	probe := filepath.Join(frontend, "app", "admin", "(dashboard)", "x", "page.tsx")
	got := embeddedAdminContent(probe, `import { X } from "@/components/x";`)
	if !strings.Contains(got, `"@admin/components/x"`) {
		t.Errorf("an admin page in frontend/ kept its standalone import, which resolves to the host app: %s", got)
	}

	// And the panel's own code, under admin-panel/.
	probe2 := filepath.Join(frontend, "admin-panel", "components", "x.tsx")
	got2 := embeddedAdminContent(probe2, `import { Y } from "@/lib/y";`)
	if !strings.Contains(got2, `"@admin/lib/y"`) {
		t.Errorf("a panel component in frontend/ kept its standalone import: %s", got2)
	}

	// A file in the host app's own tree is left alone: it means its own "@/".
	probe3 := filepath.Join(frontend, "app", "blog", "page.tsx")
	got3 := embeddedAdminContent(probe3, `import { Z } from "@/components/z";`)
	if strings.Contains(got3, "@admin/") {
		t.Errorf("the web app's own import was repointed at the admin panel: %s", got3)
	}
}

// One workspace, no task runner.
//
// The frontend asks for "@repo/shared": "workspace:*", so pnpm needs a
// workspace root or install fails on a dependency that is not there. It does not
// need Turborepo: two packages have no task graph worth the dependency, and
// Turborepo is the part of a monorepo --single exists to avoid.
func TestSingleNextHasAWorkspaceButNoTurborepo(t *testing.T) {
	opts := singleNextOptions()

	if !opts.UsesPnpmWorkspace() {
		t.Error("no workspace root, so pnpm install fails on @repo/shared")
	}
	if opts.ShouldUseTurborepo() {
		t.Error("a Next single would get turbo.json and a turbo dependency it has no use for")
	}
	if !opts.ShouldIncludeShared() {
		t.Error("packages/shared is not written, so the dependency the frontend declares does not exist")
	}

	// The workspace lists the frontend a single project actually has.
	ws := pnpmWorkspace(opts)
	if !strings.Contains(ws, `- "frontend"`) {
		t.Errorf("the workspace does not list frontend/, so it has no frontend in it:\n%s", ws)
	}
	if strings.Contains(ws, `- "apps/*"`) {
		t.Errorf("the workspace lists apps/*, which a single project does not have:\n%s", ws)
	}

	// And the scripts run that frontend rather than a task runner that is absent.
	pkg := rootPackageJSON(opts)
	if strings.Contains(pkg, "turbo") {
		t.Errorf("a root script calls turbo, which is not a dependency here:\n%s", pkg)
	}
	if !strings.Contains(pkg, `"dev": "pnpm --filter ./frontend dev"`) {
		t.Errorf("grit start would bring up no frontend:\n%s", pkg)
	}
}

// The binary embeds a frontend only when there is a static one to embed.
//
// The Vite build is static files, so the binary carries them and the project
// ships as one executable. Next is not static here: the panel alone has four
// [id] routes and the web app two more, and `output: "export"` cannot build a
// dynamic segment whose values are rows in a database. Keeping the embed would
// mean a //go:embed of a directory Next never writes, which does not compile.
func TestOnlyTheViteSingleEmbedsItsFrontend(t *testing.T) {
	viteRoot := t.TempDir()
	if err := writeSingleMainGo(viteRoot, Options{ProjectName: "app", Architecture: ArchSingle}); err != nil {
		t.Fatal(err)
	}
	vite := readTestFile(t, filepath.Join(viteRoot, "main.go"))
	if !strings.Contains(vite, "//go:embed all:frontend/dist") {
		t.Error("the Vite single stopped embedding its SPA, so it is no longer one binary")
	}

	nextRoot := t.TempDir()
	if err := writeSingleMainGo(nextRoot, singleNextOptions()); err != nil {
		t.Fatal(err)
	}
	next := readTestFile(t, filepath.Join(nextRoot, "main.go"))
	if strings.Contains(next, "go:embed") {
		t.Error("a Next single embeds frontend/dist, a directory Next never writes: the project does not compile")
	}
	if !fileExists(filepath.Join(viteRoot, "frontend", "dist", "index.html")) {
		t.Error("the Vite single lost the placeholder that lets go build work before the first pnpm build")
	}
	if fileExists(filepath.Join(nextRoot, "frontend", "dist", "index.html")) {
		t.Error("a Next single got a frontend/dist placeholder it never serves")
	}
}

// The seeder writes an .env.local into whichever frontend this project has.
//
// It listed ../web and ../admin, which are where a monorepo's apps sit relative
// to apps/api. A single project runs the seeder from its root and has its
// frontend at frontend/, so it matched neither and got no file: the app fell
// back to localhost:8080 whatever the API was listening on, and had no
// publishable key, so every public endpoint answered INVALID_API_KEY. Both
// failures surface as a CORS error in a browser console, which names neither
// cause.
func TestSeederWritesAnEnvForEveryFrontendLayout(t *testing.T) {
	seeder := apiAPIKeySeederGo()

	for _, want := range []string{
		`filepath.Join("..", "web", ".env.local")`,
		`filepath.Join("..", "admin", ".env.local")`,
		`filepath.Join("frontend", ".env.local")`,
	} {
		if !strings.Contains(seeder, want) {
			t.Errorf("the seeder does not write %s, so that frontend gets no API URL and no publishable key", want)
		}
	}

	// It has to keep skipping the ones this project does not have, or a triple
	// grows a frontend/ directory holding one env file.
	if !strings.Contains(seeder, "continue // that app is not part of this project") {
		t.Error("the seeder no longer skips a frontend this project does not have")
	}
}

// grit upgrade finds a single project's API.
//
// The API block was gated on dirExists(root/apps/api), a path neither single
// has: the Go module sits at the project root. So every API writer below that
// gate was skipped, and `grit upgrade` on a single project printed "Upgrade
// complete" after updating its root config and Docker files and nothing else.
// No API fix had ever reached one. The Vite single that proved this took 548
// files on the first upgrade that looked in the right place.
//
// Checked against a directory rather than a string, because the bug was a path
// that existed in the author's head and not on disk.
func TestUpgradeFindsTheAPIInBothLayouts(t *testing.T) {
	for _, tc := range []struct {
		name   string
		opts   Options
		apiDir []string
	}{
		{"single", Options{ProjectName: "a", Architecture: ArchSingle}, []string{"internal", "routes"}},
		{"single next", singleNextOptions(), []string{"internal", "routes"}},
		{"triple", Options{ProjectName: "a", Architecture: ArchTriple, Frontend: FrontendNext}, []string{"apps", "api", "internal", "routes"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(append([]string{root}, tc.apiDir...)...), 0o755); err != nil {
				t.Fatal(err)
			}
			// The same expression the upgrade uses.
			if !dirExists(filepath.Join(tc.opts.APIRoot(root), "internal", "routes")) {
				t.Error("the upgrade would skip every API writer for this layout")
			}
		})
	}

	// And a project with no Go API at all is still skipped, or the upgrade
	// writes an API into a frontend-only directory.
	root := t.TempDir()
	opts := Options{ProjectName: "a", Architecture: ArchSingle}
	if dirExists(filepath.Join(opts.APIRoot(root), "internal", "routes")) {
		t.Error("an empty directory was taken for a Go API")
	}
}
