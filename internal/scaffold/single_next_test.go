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

// A single project's frontend owns the project root: the app is the project,
// with its Go API in api/. One scaffolded before that keeps its frontend where
// it was put, because an upgrade that moved it would break every import path in
// it at once.
func TestSingleKeepsItsFrontendBesideTheGoModule(t *testing.T) {
	root := filepath.FromSlash("/p")
	for _, tc := range []struct {
		name string
		opts Options
		want string
	}{
		{"single next", singleNextOptions(), root},
		{"single vite", Options{ProjectName: "a", Architecture: ArchSingle}, root},
		{"legacy single", Options{ProjectName: "a", Architecture: ArchSingle, LegacySingleFlat: true}, filepath.Join(root, "frontend")},
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

	frontend := webAppRoot(root, opts)
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

// Both singles ship as one binary with their frontend inside it.
//
// The Vite one always did. The Next one does now: it builds with
// output: "export", which is static HTML, CSS and JS, and the export is copied
// into api/web where //go:embed reads it.
//
// It did not, for a while, and the reason was a real constraint rather than an
// oversight: a static export cannot build a dynamic segment, because the values
// are rows in a database the build never sees. The pages that had one read an id
// from the query string instead, which is what made the export possible.
func TestBothSinglesEmbedTheirFrontend(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts Options
	}{
		{"vite", Options{ProjectName: "app", Architecture: ArchSingle}},
		{"next", singleNextOptions()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writeSingleMainGo(root, tc.opts); err != nil {
				t.Fatal(err)
			}

			main := readTestFile(t, filepath.Join(root, "api", "main.go"))
			if !strings.Contains(main, "//go:embed all:web") {
				t.Error("the binary does not embed api/web, so it is not one file")
			}

			// The placeholder lets go build work on a fresh clone: //go:embed
			// fails outright on a pattern that matches nothing.
			if !fileExists(filepath.Join(root, "api", "web", "index.html")) {
				t.Error("no placeholder, so go build fails before the frontend is built once")
			}
		})
	}
}

// And the Next single really is exported rather than served by a Node process.
func TestTheNextSingleIsAStaticExport(t *testing.T) {
	opts := singleNextOptions()
	if !opts.StaticExport() {
		t.Fatal("a Next single is not marked as a static export, so nothing below it applies")
	}

	cfg := webNextConfig(opts)
	if !strings.Contains(cfg, `output: "export"`) {
		t.Error("the Next config does not export, so the build produces a server to run")
	}
	if !strings.Contains(cfg, "unoptimized: true") {
		t.Error("the image optimizer is left on, and an export has no server to run one: every image 404s")
	}
	if strings.Contains(cfg, "async headers()") {
		t.Error("headers() survives into an export, where Next warns it does nothing")
	}

	// The build has to put the export where the binary reads it.
	pkg := webPackageJSON(opts)
	if !strings.Contains(pkg, "node scripts/embed-frontend.mjs") {
		t.Error("the build does not copy the export into api/web, so the binary embeds the placeholder")
	}

	// A monorepo is untouched by all of it.
	mono := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}
	if mono.StaticExport() {
		t.Error("a monorepo was marked as a static export")
	}
	if !strings.Contains(webNextConfig(mono), `output: "standalone"`) {
		t.Error("a monorepo stopped building a server")
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
		{"single", Options{ProjectName: "a", Architecture: ArchSingle}, []string{"api", "internal", "routes"}},
		{"single next", singleNextOptions(), []string{"api", "internal", "routes"}},
		{"legacy single", Options{ProjectName: "a", Architecture: ArchSingle, LegacySingleFlat: true}, []string{"internal", "routes"}},
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

// The Vite single builds into the Go tree, and into this project's.
//
// The binary embeds api/web, so that is where Vite has to write. The path is
// relative to vite.config.ts, which sits at the project root: "../api/web"
// resolves to the directory ABOVE the project, where the built app is read by
// nothing. Vite reported success, the assets landed outside the repository, and
// the binary embedded the placeholder index.html instead, so the one-file deploy
// served an empty page and said nothing about it.
func TestViteSingleBuildsIntoTheEmbeddedDirectory(t *testing.T) {
	cfg := singleFrontendViteConfig()

	if !strings.Contains(cfg, "outDir: './api/web'") {
		t.Error("vite does not build into ./api/web, which is the directory the binary embeds")
	}
	if strings.Contains(cfg, "'../api/web'") {
		t.Error("outDir escapes the project: relative to the config at the root, ../ is the directory above it")
	}
	// Emptied, or a rename leaves the previous build's assets behind and the
	// binary carries both.
	if !strings.Contains(cfg, "emptyOutDir: true") {
		t.Error("outDir is never emptied, so stale assets accumulate inside the binary")
	}

	// And the directive reads the same directory.
	main := singleMainGo(Options{ProjectName: "app", Architecture: ArchSingle})
	if !strings.Contains(main, "//go:embed all:web") {
		t.Error("main.go does not embed web/, so it cannot be serving what Vite builds")
	}
}

// An existing single project keeps the layout it was built with.
//
// Moving the Go module into api/ and giving the project root to the frontend is
// right for a new project and catastrophic for an old one: an upgrade that
// started writing into api/ would leave two halves of an API, and neither would
// compile. LegacySingleFlat is read from disk by the upgrade and pins every path
// to where that project already has them.
func TestAnExistingSingleKeepsItsOwnLayout(t *testing.T) {
	root := filepath.FromSlash("/p")
	legacy := Options{ProjectName: "a", Architecture: ArchSingle, LegacySingleFlat: true}
	fresh := Options{ProjectName: "a", Architecture: ArchSingle}

	if got := legacy.APIRoot(root); got != root {
		t.Errorf("an existing project's Go code moved to %s", got)
	}
	if got := fresh.APIRoot(root); got != filepath.Join(root, "api") {
		t.Errorf("a new project's Go code is at %s, not api/", got)
	}

	if got := webAppRoot(root, legacy); got != filepath.Join(root, "frontend") {
		t.Errorf("an existing project's frontend moved to %s, which breaks every relative path in it", got)
	}
	if got := webAppRoot(root, fresh); got != root {
		t.Errorf("a new project's frontend is at %s, not the project root", got)
	}

	// The CI workflows point at whichever module directory this project has.
	if got := ciLayoutFor(legacy); got.apiDir != "." || got.jsDir != "frontend" {
		t.Errorf("CI points at %s/%s for an existing project", got.apiDir, got.jsDir)
	}
	if got := ciLayoutFor(fresh); got.apiDir != "api" || got.jsDir != "." {
		t.Errorf("CI points at %s/%s for a new project", got.apiDir, got.jsDir)
	}
}

// And the generated config finds .env from wherever the binary runs.
//
// The project root holds one .env. The binary runs from the root in an existing
// single, api/ in a new one, apps/api in a monorepo. Only two of those three
// were tried, so `grit migrate` in a new single found no .env at all and refused
// to start under APP_ENV=production, naming five secrets that were sitting one
// directory above it.
func TestGeneratedConfigFindsTheEnvFromEveryModuleDirectory(t *testing.T) {
	cfg := apiConfigGo()
	for _, want := range []string{
		`godotenv.Load()`,
		`godotenv.Load("../.env")`,
		`godotenv.Load("../../.env")`,
	} {
		if !strings.Contains(cfg, want) {
			t.Errorf("config does not try %s, so a binary run from that depth gets the built-in defaults", want)
		}
	}
}
