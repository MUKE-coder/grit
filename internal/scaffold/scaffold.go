package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// Architecture represents the project architecture mode.
type Architecture string

const (
	ArchSingle Architecture = "single" // Go API + embedded React SPA (one binary)
	ArchDouble Architecture = "double" // Turborepo: Web + API
	ArchTriple Architecture = "triple" // Turborepo: Web + Admin + API
	ArchAPI    Architecture = "api"    // Go API only (no frontend)
	ArchMobile Architecture = "mobile" // Turborepo: API + Expo mobile
)

// Frontend represents the frontend framework choice.
type Frontend string

const (
	FrontendNext     Frontend = "next"     // Next.js (App Router, SSR)
	FrontendTanStack Frontend = "tanstack" // TanStack Router (Vite, SPA)
)

// Options holds the scaffolding configuration.
type Options struct {
	ProjectName  string
	Architecture Architecture
	Frontend     Frontend
	// LegacySingleFlat marks a single project whose Go code is at the root
	// rather than in api/, which is every single project scaffolded before
	// v3.380.0. Set by the upgrade from what is on disk, never by `grit new`.
	LegacySingleFlat bool
	Style            string
	Theme            string // Full theme: atlas (default), aurora, pulse — controls auth pages, dashboard tokens, fonts, brand colors
	// DBProvider is the database engine written into .env as DB_PROVIDER:
	// postgres (the default), mysql, sqlite or memory. Connect has understood all
	// of them for a long time; until v3.234.0 the only way to pick one was the
	// prefix of DATABASE_URL, so .env had Postgres settings and no place for
	// anybody else's.
	DBProvider     string
	InPlace        bool // Scaffold into current directory (grit new .)
	Force          bool // Allow scaffolding into non-empty directory (--force)
	IncludeDesktop bool // Add apps/desktop (Wails client that shares the monorepo API)

	// Version is the Grit CLI version that scaffolded this project (e.g. "3.25.2").
	// Injected by cmd/grit/main.go so scaffolded README + docs reflect the real
	// version instead of a hardcoded constant. Falls back to DefaultVersion
	// in Normalize() so callers that forget to set it (tests, library users)
	// still get a sensible value.
	Version string

	// Deprecated: use Architecture instead. Kept for backward compatibility.
	APIOnly     bool
	IncludeExpo bool
	MobileOnly  bool
	Full        bool
}

// DefaultVersion is the fallback string written into scaffolded README/docs
// when Options.Version is empty. Kept in sync with cmd/grit/main.go's
// version variable on release.
const DefaultVersion = "3.394.4"

// Normalize maps legacy boolean flags to the new Architecture enum.
// Call this after constructing Options from CLI flags.
func (o *Options) Normalize() {
	if o.Version == "" {
		o.Version = DefaultVersion
	}

	// If Architecture is already set, it takes priority
	if o.Architecture != "" {
		return
	}

	// Map legacy booleans to Architecture
	switch {
	case o.APIOnly:
		o.Architecture = ArchAPI
	case o.MobileOnly:
		o.Architecture = ArchMobile
	case o.Full:
		o.Architecture = ArchTriple // full includes everything
	default:
		o.Architecture = ArchTriple // default: triple (web + admin + api)
	}

	// Legacy expo flag: triple + expo
	if o.IncludeExpo && o.Architecture == ArchTriple {
		// Keep as triple but also include expo
	}

	// Default frontend to Next.js if not set
	if o.Frontend == "" {
		o.Frontend = FrontendNext
	}
}

// ValidStyles lists all supported admin panel style variants.
var ValidStyles = []string{"default", "modern", "minimal", "glass", "centered"}

// DBProviderOrder is the order the picker offers the engines in, commonest
// first. A map has no order, and a picker whose options move between runs is
// one people stop reading.
var DBProviderOrder = []string{"postgres", "mysql", "sqlite", "memory"}

// DBProviders are the engines .env can name, with what each one means.
var DBProviders = map[string]string{
	"postgres": "PostgreSQL, and what docker-compose.yml starts for you",
	"mysql":    "MySQL 8 or MariaDB, on a server you run",
	"sqlite":   "one file, pure Go, no CGO and no server",
	"memory":   "SQLite in RAM, empty at every boot, for tests and demos",
}

// ValidateDBProvider checks the --db value, and normalises the spellings people
// reach for.
func (o *Options) ValidateDBProvider() error {
	if o.DBProvider == "" {
		o.DBProvider = "postgres"
		return nil
	}
	normalised := map[string]string{
		"postgresql": "postgres", "pg": "postgres", "postgres": "postgres",
		"mariadb": "mysql", "mysql": "mysql",
		"sqlite3": "sqlite", "sqlite": "sqlite", "file": "sqlite",
		"memory": "memory", ":memory:": "memory", "inmemory": "memory",
	}
	want := strings.ToLower(strings.TrimSpace(o.DBProvider))
	if name, ok := normalised[want]; ok {
		o.DBProvider = name
		return nil
	}
	return fmt.Errorf("--db %q is not a database Grit knows: use postgres, mysql, sqlite or memory", o.DBProvider)
}

// ValidateStyle checks that the Style field is a supported value.
// If empty, it defaults to "default".
func (o *Options) ValidateStyle() error {
	if o.Style == "" {
		o.Style = "default"
		return nil
	}
	for _, s := range ValidStyles {
		if o.Style == s {
			return nil
		}
	}
	return fmt.Errorf("invalid style %q: must be one of %s", o.Style, strings.Join(ValidStyles, ", "))
}

// ValidThemes lists the full themes shipped by Grit v3.28+.
// A theme controls auth pages, dashboard tokens, fonts, sidebar treatment,
// card styling, and the Pulse + Sentinel widget palette — picked once at
// scaffold time and overridable at runtime via THEME=<name> in .env.
var ValidThemes = []string{"atlas", "aurora", "pulse", "coral", "amber", "sky", "mono", "emerald"}

// ValidateTheme checks that the Theme field is a supported value.
// If empty, it defaults to "atlas" — the team/organisation theme,
// chosen as the default because it works for the widest audience.
func (o *Options) ValidateTheme() error {
	if o.Theme == "" {
		o.Theme = "atlas"
		return nil
	}
	for _, t := range ValidThemes {
		if o.Theme == t {
			return nil
		}
	}
	return fmt.Errorf("invalid theme %q: must be one of %s", o.Theme, strings.Join(ValidThemes, ", "))
}

// ShouldIncludeWeb returns true if a web frontend app should be scaffolded (Turborepo web app).
func (o Options) ShouldIncludeWeb() bool {
	return o.Architecture == ArchDouble || o.Architecture == ArchTriple
}

// ShouldIncludeAdmin returns true if the admin panel should be scaffolded.
func (o Options) ShouldIncludeAdmin() bool {
	return o.Architecture == ArchTriple
}

// HasAdminPanel reports whether this project has an admin panel at all, wherever
// it lives: its own app in a triple, or a route group inside the web app in a
// double.
//
// Distinct from ShouldIncludeAdmin, which is about the separate application: a
// container, a port, a dev script, a package.json. Every writer that produces a
// SCREEN wants this one. Asking the other question is how a double ended up with
// an admin panel missing its account-security page while the link to it sat in the
// user menu.
func (o Options) HasAdminPanel() bool {
	return o.ShouldIncludeAdmin() || o.ShouldEmbedAdmin() || o.ShouldEmbedAdminInSPA()
}

// ShouldIncludeSingleSPA returns true if this is a single-app embedded SPA.
func (o Options) ShouldIncludeSingleSPA() bool {
	return o.Architecture == ArchSingle
}

// ShouldUseTurborepo returns true if the project uses a Turborepo monorepo.
func (o Options) ShouldUseTurborepo() bool {
	return o.Architecture == ArchDouble ||
		o.Architecture == ArchTriple ||
		o.Architecture == ArchMobile ||
		o.IncludeDesktop
}

// UsesPnpmWorkspace reports whether pnpm runs from a workspace root here.
//
// Every Turborepo project, and a Next single. The Next single has two JS
// packages that have to see each other, frontend/ and packages/shared, and
// nothing simpler than a workspace makes "@repo/shared" resolve from both the
// type checker and the bundler. It gets no turbo.json: two packages do not need
// a task graph, and Turborepo was the part of a monorepo that --single existed
// to avoid.
func (o Options) UsesPnpmWorkspace() bool {
	return o.ShouldUseTurborepo() || o.SingleUsesNext()
}

// ShouldIncludeShared returns true if the shared package should be scaffolded.
func (o Options) ShouldIncludeShared() bool {
	return o.UsesPnpmWorkspace()
}

// ShouldIncludeFrontend returns true if any frontend (web, admin, or SPA) is included.
func (o Options) ShouldIncludeFrontend() bool {
	return o.Architecture != ArchAPI
}

// ShouldIncludeExpo returns true if the Expo app should be scaffolded.
func (o Options) ShouldIncludeExpo() bool {
	return o.Architecture == ArchMobile || o.IncludeExpo || o.Full
}

// ShouldIncludeDesktop returns true if the desktop (Wails) client should be
// scaffolded as part of the monorepo. This is a separate capability from
// `grit new-desktop` (which is a standalone offline-first app).
// The --desktop flag adds apps/desktop/ to the monorepo — a Wails window
// that calls the shared Go API over HTTP (same as the Expo app does).
func (o Options) ShouldIncludeDesktop() bool {
	return o.IncludeDesktop || o.Full
}

// ShouldIncludeDocs returns true if the docs site should be scaffolded.
func (o Options) ShouldIncludeDocs() bool {
	return o.Full
}

// UseTanStack returns true if the frontend uses TanStack Router (Vite).
func (o Options) UseTanStack() bool {
	return o.Frontend == FrontendTanStack
}

// SingleUsesNext reports whether this is a single project built with Next.js.
//
// A single project is one folder with no monorepo tooling: the Go module at the
// root, the frontend in frontend/. Which frontend was, until now, not a choice:
// --single --next accepted the flag and produced a Vite app anyway, and the only
// sign was grit.json claiming "next" beside a vite.config.ts.
//
// The two singles differ in one further way, and it is not cosmetic. The Vite
// SPA is static, so the Go binary embeds it and a single project ships as one
// file. Next.js is not static here: the admin alone has four [id] routes and the
// web app two more, and `output: "export"` cannot build a dynamic segment whose
// values are rows in a database. So a Next single runs the API and the Next
// server as two processes, which is how every Next.js app runs, and keeps the
// thing that made single worth having: no Turborepo, no workspace, one folder.
func (o Options) SingleUsesNext() bool {
	return o.Architecture == ArchSingle && o.Frontend == FrontendNext
}

// StaticExport reports whether this project's Next.js app is built to static
// files the Go binary embeds, rather than served by a Node process.
//
// A Next single is: it builds with output: "export", which writes plain HTML,
// JS and CSS, and every byte of data comes from the Go API the same binary
// serves. That is what makes a single project one file.
//
// The cost is that a dynamic segment cannot be built without knowing its
// values, and the values are rows in a database. So the pages that had one read
// an id from the query string instead: /resources/users/view?id=... rather than
// /resources/users/123. Same screen, same component, a URL the build can
// produce.
//
// Not the Vite single, which is also static and does not have this problem: a
// SPA serves index.html for every path and resolves the route in the browser,
// so /resources/users/123 works there without a file to match it.
func (o Options) StaticExport() bool {
	return o.SingleUsesNext()
}

// AdminIsTanStack reports whether the admin panel is a TanStack Router app,
// which decides the shape of every screen written for it: a page plus a route
// shim rather than one file at its route path, and no next/link.
//
// Not the same question as UseTanStack. Frontend picks the web app's framework
// in a monorepo and is empty for a single project, whose one SPA is a TanStack
// app whatever that field says. Writers that asked UseTanStack emitted Next code
// for it, and one of them put a page in an apps/admin directory a single project
// does not have.
func (o Options) AdminIsTanStack() bool {
	return o.UseTanStack() || o.ShouldEmbedAdminInSPA()
}

// APIRoot returns the base directory for Go API files.
//
// A single project looks like the frontend it is: the app at the root, the way
// a Next.js or Vite project is laid out everywhere else, and the Go API in
// api/. Before this the Go module owned the root and the frontend was pushed
// into frontend/, which reads as a Go repository that happens to contain a web
// app. It is the other way round for the people who use it: they open the
// project to work on a page, and the API is the part they reach for less often.
//
// Monorepo: apps/api, unchanged.
//
// LegacySingleFlat keeps an existing project working. A project scaffolded
// before this has its Go code at the root, and an upgrade that started writing
// into api/ would leave it with two halves of an API and compile neither.
func (o Options) APIRoot(root string) string {
	if o.Architecture == ArchSingle {
		if o.LegacySingleFlat {
			return root
		}
		return filepath.Join(root, "api")
	}
	return filepath.Join(root, "apps", "api")
}

// Module returns the Go module path for the API.
// Single app: project-name. Monorepo: project-name/apps/api.
func (o Options) Module() string {
	if o.Architecture == ArchSingle {
		return o.ProjectName
	}
	return o.ProjectName + "/apps/api"
}

// ValidateProjectName ensures the project name is lowercase, alphanumeric, and hyphens only.
func ValidateProjectName(name string) error {
	re := regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	if !re.MatchString(name) {
		return fmt.Errorf("invalid project name %q: must be lowercase, alphanumeric, and hyphens only (start with a letter)", name)
	}
	if strings.HasSuffix(name, "-") {
		return fmt.Errorf("invalid project name %q: must not end with a hyphen", name)
	}
	return nil
}

func resolveScaffoldRoot(opts Options) (string, bool, error) {
	if opts.InPlace {
		return ".", true, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("getting current directory: %w", err)
	}

	// Quality-of-life behavior: if the user is already inside a directory whose
	// name matches the project name, scaffold in place instead of nesting.
	if filepath.Base(cwd) == opts.ProjectName {
		return ".", true, nil
	}

	return opts.ProjectName, false, nil
}

func ensureTargetDirectory(root string, inPlace bool, force bool) error {
	if !inPlace {
		if _, err := os.Stat(root); err == nil {
			return fmt.Errorf("directory %q already exists", root)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("checking directory %q: %w", root, err)
		}
		return nil
	}

	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(root, 0755); err != nil {
				return fmt.Errorf("creating directory %q: %w", root, err)
			}
			return nil
		}
		return fmt.Errorf("checking directory %q: %w", root, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("target %q is not a directory", root)
	}
	if force {
		return nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("reading directory %q: %w", root, err)
	}
	if len(entries) > 0 {
		return fmt.Errorf("current directory is not empty; rerun with --force to scaffold in place")
	}

	return nil
}

// Run executes the full scaffolding process.
func Run(opts Options) error {
	opts.Normalize()

	// Dispatch to single app scaffold if applicable
	if opts.Architecture == ArchSingle {
		return RunSingle(opts)
	}

	root := opts.ProjectName
	if opts.InPlace {
		root = "."
	}

	if !opts.InPlace {
		if _, err := os.Stat(root); err == nil {
			return fmt.Errorf("directory %q already exists", root)
		}
	}

	// Record what this scaffold writes, so a later upgrade can tell an
	// untouched framework file from one the developer has edited.
	release, err := manifest.Start(root, opts.Version, "scaffold")
	if err != nil {
		return err
	}
	defer func() {
		if saveErr := release(); saveErr != nil {
			color.New(color.FgHiBlack).Printf("  Could not write .grit/manifest.json: %v\n", saveErr)
		}
	}()

	// Results, not activity: the stage that is running shows as running and
	// is rewritten to a check when it finishes. In a pipe or under NO_COLOR
	// there is no cursor control, just one plain line per finished stage.
	steps := ui.NewStepper(46)
	spinner := color.New(color.FgHiBlack)

	// Create directory structure
	steps.Start("Directory structure")
	if err := createDirectories(root, opts); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	// Write root config files
	steps.Start("Configuration files")
	if err := writeRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing root files: %w", err)
	}

	// Write Go API files
	steps.Start("Go API")
	if err := writeAPIFiles(root, opts); err != nil {
		return fmt.Errorf("writing API files: %w", err)
	}

	// Write migrate and seed entrypoints
	steps.Start("Migration and seed tools")
	if err := writeMigrateSeedFiles(root, opts); err != nil {
		return fmt.Errorf("writing migrate/seed files: %w", err)
	}

	// Write Phase 4 service files (cache, storage, mail, jobs, cron, AI)
	steps.Start("Batteries")
	if err := writeCacheFiles(root, opts); err != nil {
		return fmt.Errorf("writing cache files: %w", err)
	}
	if err := writeUploadPackageFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminSecurityFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminAccountFiles(root, opts); err != nil {
		return err
	}
	if err := writePasswordStrengthFiles(root, opts); err != nil {
		return err
	}
	if err := writeMagicLinkFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminPasskeyFiles(root, opts); err != nil {
		return err
	}
	if err := writePasskeyFiles(root, opts); err != nil {
		return err
	}
	if err := writeRecoveryFiles(root, opts); err != nil {
		return err
	}
	if err := writeMoneyFiles(root, opts); err != nil {
		return err
	}
	if err := writeAppendOnlyFiles(root, opts); err != nil {
		return err
	}
	if err := writeOutboxFiles(root, opts); err != nil {
		return err
	}
	if err := writeStockFiles(root, opts); err != nil {
		return err
	}
	if err := writeJSONTimeFiles(root, opts); err != nil {
		return err
	}
	if err := writeMediaFiles(root, opts); err != nil {
		return err
	}
	// The browser half of realtime: one shared socket, reconnection, and
	// the React Query bindings. The hub has shipped for a long time with
	// nothing on the client able to consume it.
	if err := writeRealtimeClientFiles(root, opts); err != nil {
		return err
	}
	if err := writeStorageFiles(root, opts); err != nil {
		return fmt.Errorf("writing storage files: %w", err)
	}
	if err := writeBackupFiles(root, opts); err != nil {
		return fmt.Errorf("writing backup files: %w", err)
	}
	if err := writeMailFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail files: %w", err)
	}
	if err := writeJobsFiles(root, opts); err != nil {
		return fmt.Errorf("writing jobs files: %w", err)
	}
	if err := writeMailDispatchFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail dispatch files: %w", err)
	}
	if err := writeRequestMetaFiles(root, opts); err != nil {
		return fmt.Errorf("writing request meta files: %w", err)
	}
	if err := writeCronFiles(root, opts); err != nil {
		return fmt.Errorf("writing cron files: %w", err)
	}
	if err := writeAIFiles(root, opts); err != nil {
		return fmt.Errorf("writing AI files: %w", err)
	}
	if err := writeTOTPFiles(root, opts); err != nil {
		return fmt.Errorf("writing TOTP files: %w", err)
	}
	if err := writeSecurityFiles(root, opts); err != nil {
		return fmt.Errorf("writing security files: %w", err)
	}
	if err := writeTestingFiles(root, opts); err != nil {
		return fmt.Errorf("writing testing files: %w", err)
	}
	if err := writeSecurityObservabilityFiles(root, opts); err != nil {
		return fmt.Errorf("writing security/observability files: %w", err)
	}
	if err := writeFormShareFiles(root, opts); err != nil {
		return fmt.Errorf("writing form-share files: %w", err)
	}
	if err := writeResourceStatsFiles(root, opts); err != nil {
		return fmt.Errorf("writing resource-stats files: %w", err)
	}
	if err := writeChartFiles(root, opts); err != nil {
		return fmt.Errorf("writing chart files: %w", err)
	}

	// Write blog example files
	steps.Start("Blog example")
	if err := writeAPIBlogFiles(root, opts); err != nil {
		return fmt.Errorf("writing blog files: %w", err)
	}

	// Run go mod tidy to resolve dependencies and generate go.sum
	steps.Start("Go dependencies")
	apiDir := filepath.Join(root, "apps", "api")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = apiDir
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running go mod tidy: %w\n%s", err, string(out))
	}

	// Write Docker files
	steps.Start("Docker setup")
	if err := writeDockerFiles(root, opts); err != nil {
		return fmt.Errorf("writing Docker files: %w", err)
	}

	if opts.ShouldIncludeShared() {
		// Write shared package
		steps.Start("Shared package")
		if err := writeSharedFiles(root, opts); err != nil {
			return fmt.Errorf("writing shared files: %w", err)
		}
	}

	if opts.ShouldIncludeWeb() {
		if opts.UseTanStack() {
			steps.Start("Web app")
			if err := writeWebTanStackFiles(root, opts); err != nil {
				return fmt.Errorf("writing TanStack web files: %w", err)
			}
		} else {
			steps.Start("Web app")
			if err := writeWebFiles(root, opts); err != nil {
				return fmt.Errorf("writing web files: %w", err)
			}
		}
	}

	if opts.ShouldIncludeAdmin() {
		if opts.UseTanStack() {
			steps.Start("Admin panel")
			if err := writeAdminTanStackFiles(root, opts); err != nil {
				return fmt.Errorf("writing TanStack admin files: %w", err)
			}
		} else {
			steps.Start("Admin panel")
			if err := writeAdminFiles(root, opts); err != nil {
				return fmt.Errorf("writing admin files: %w", err)
			}
		}
	}

	// A double has no admin app, so the admin panel goes inside the web app as a
	// route group at /admin: the same screens and the same file map, moved and
	// repointed (see admin_embedded.go). Before this, a double shipped with the
	// admin link in its navbar and nothing behind it.
	// A double on Next.js. A double built with --vite takes the Vite path below,
	// which is the same panel in the dialect that app speaks.
	if opts.ShouldEmbedAdmin() && !opts.ShouldEmbedAdminInSPA() {
		steps.Start("Admin panel at /admin")
		if err := writeEmbeddedAdminFiles(root, opts); err != nil {
			return fmt.Errorf("writing the embedded admin panel: %w", err)
		}
		// The web app's middleware keeps signed-out visitors out of /admin.
		if err := writeAdminEdgeGuard(root, opts); err != nil {
			return err
		}
	}

	// A Vite app hosting the panel: a single project's SPA, or a double's web app
	// when it was scaffolded with --vite. The Vite admin is the source, because
	// both are TanStack Router apps and the panel is already written that way.
	if opts.ShouldEmbedAdminInSPA() {
		steps.Start("Admin panel at /admin")
		if err := writeEmbeddedSingleAdminFiles(root, opts); err != nil {
			return fmt.Errorf("writing the embedded admin panel: %w", err)
		}
	}

	if opts.ShouldIncludeExpo() {
		// Write Expo mobile app
		steps.Start("Expo mobile app")
		if err := writeExpoFiles(root, opts); err != nil {
			return fmt.Errorf("writing Expo files: %w", err)
		}
	}

	if opts.ShouldIncludeDesktop() {
		// Write desktop client (Wails + Vite + TanStack Router, shares the API)
		steps.Start("Desktop app")
		if err := writeDesktopClientFiles(root, opts); err != nil {
			return fmt.Errorf("writing desktop files: %w", err)
		}

		// The desktop app is its own Go module, so tidying apps/api does not
		// reach it. Without this a fresh project cannot build or test anything
		// under apps/desktop: four missing go.sum entries, including the
		// SQLite driver the offline engine is built on.
		//
		// A failure here is not fatal. The API is already resolved and the
		// project works; losing the whole scaffold to a network hiccup would
		// be the worse outcome, so say what to run and carry on.
		steps.Start("Desktop Go dependencies")
		desktopTidy := exec.Command("go", "mod", "tidy")
		desktopTidy.Dir = filepath.Join(root, "apps", "desktop")
		if out, err := desktopTidy.CombinedOutput(); err != nil {
			spinner.Printf("  ⚠ could not resolve desktop dependencies: %v\n", err)
			spinner.Printf("    run: cd apps/desktop && go mod tidy\n%s\n", string(out))
		}
	}

	if opts.ShouldIncludeDocs() {
		// Write docs site
		steps.Start("Documentation site")
		if err := writeDocsFiles(root, opts); err != nil {
			return fmt.Errorf("writing docs files: %w", err)
		}
	}

	// Write frontend test files (Vitest + Playwright)
	if opts.ShouldIncludeWeb() || opts.ShouldIncludeAdmin() {
		steps.Start("Frontend tests")
		if err := writeFrontendTestFiles(root, opts); err != nil {
			return fmt.Errorf("writing frontend test files: %w", err)
		}
	}

	// One React across the monorepo when it has the Expo app.
	alignReactVersions(root)

	// Close the last stage, or its running line is the last thing on screen.
	steps.Finish()
	return nil
}

// RunSingle executes the single-app scaffolding process.
// Single app: Go API + embedded React SPA, one binary, no Turborepo.
func RunSingle(opts Options) error {
	root := opts.ProjectName
	if opts.InPlace {
		root = "."
	}

	if !opts.InPlace {
		if _, err := os.Stat(root); err == nil {
			return fmt.Errorf("directory %q already exists", root)
		}
	}

	// Record what this scaffold writes, so a later upgrade can tell an
	// untouched framework file from one the developer has edited.
	release, err := manifest.Start(root, opts.Version, "scaffold")
	if err != nil {
		return err
	}
	defer func() {
		if saveErr := release(); saveErr != nil {
			color.New(color.FgHiBlack).Printf("  Could not write .grit/manifest.json: %v\n", saveErr)
		}
	}()

	steps := ui.NewStepper(46)

	// Create directory structure
	steps.Start("Directory structure")
	if err := createSingleDirectories(root, opts); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	// Write root config files (.env, .gitignore, skill file, Makefile)
	steps.Start("Configuration files")
	if err := writeSingleRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing root files: %w", err)
	}

	// Write .env file (reuse from root_files but with single-app paths)
	if err := writeRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing env files: %w", err)
	}

	// Write Go API files (uses opts.APIRoot which returns root for single)
	steps.Start("Go API")
	if err := writeAPIFiles(root, opts); err != nil {
		return fmt.Errorf("writing API files: %w", err)
	}

	// Write migrate/seed tools
	steps.Start("Migration and seed tools")
	if err := writeMigrateSeedFiles(root, opts); err != nil {
		return fmt.Errorf("writing migrate/seed files: %w", err)
	}

	// Write batteries
	steps.Start("Batteries")
	if err := writeCacheFiles(root, opts); err != nil {
		return fmt.Errorf("writing cache files: %w", err)
	}
	if err := writeUploadPackageFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminSecurityFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminAccountFiles(root, opts); err != nil {
		return err
	}
	if err := writePasswordStrengthFiles(root, opts); err != nil {
		return err
	}
	if err := writeMagicLinkFiles(root, opts); err != nil {
		return err
	}
	if err := writeAdminPasskeyFiles(root, opts); err != nil {
		return err
	}
	if err := writePasskeyFiles(root, opts); err != nil {
		return err
	}
	if err := writeRecoveryFiles(root, opts); err != nil {
		return err
	}
	if err := writeMoneyFiles(root, opts); err != nil {
		return err
	}
	if err := writeAppendOnlyFiles(root, opts); err != nil {
		return err
	}
	if err := writeOutboxFiles(root, opts); err != nil {
		return err
	}
	if err := writeStockFiles(root, opts); err != nil {
		return err
	}
	if err := writeJSONTimeFiles(root, opts); err != nil {
		return err
	}
	if err := writeMediaFiles(root, opts); err != nil {
		return err
	}
	if err := writeStorageFiles(root, opts); err != nil {
		return fmt.Errorf("writing storage files: %w", err)
	}
	if err := writeBackupFiles(root, opts); err != nil {
		return fmt.Errorf("writing backup files: %w", err)
	}
	if err := writeMailFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail files: %w", err)
	}
	if err := writeJobsFiles(root, opts); err != nil {
		return fmt.Errorf("writing jobs files: %w", err)
	}
	if err := writeMailDispatchFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail dispatch files: %w", err)
	}
	if err := writeRequestMetaFiles(root, opts); err != nil {
		return fmt.Errorf("writing request meta files: %w", err)
	}
	if err := writeCronFiles(root, opts); err != nil {
		return fmt.Errorf("writing cron files: %w", err)
	}
	if err := writeAIFiles(root, opts); err != nil {
		return fmt.Errorf("writing AI files: %w", err)
	}
	if err := writeTOTPFiles(root, opts); err != nil {
		return fmt.Errorf("writing TOTP files: %w", err)
	}
	if err := writeSecurityFiles(root, opts); err != nil {
		return fmt.Errorf("writing security files: %w", err)
	}
	if err := writeTestingFiles(root, opts); err != nil {
		return fmt.Errorf("writing testing files: %w", err)
	}
	if err := writeSecurityObservabilityFiles(root, opts); err != nil {
		return fmt.Errorf("writing security/observability files: %w", err)
	}
	if err := writeFormShareFiles(root, opts); err != nil {
		return fmt.Errorf("writing form-share files: %w", err)
	}
	if err := writeResourceStatsFiles(root, opts); err != nil {
		return fmt.Errorf("writing resource-stats files: %w", err)
	}
	if err := writeChartFiles(root, opts); err != nil {
		return fmt.Errorf("writing chart files: %w", err)
	}

	// Write blog example
	steps.Start("Blog example")
	if err := writeAPIBlogFiles(root, opts); err != nil {
		return fmt.Errorf("writing blog files: %w", err)
	}

	// Write embed-aware main.go (replaces the standard cmd/server/main.go)
	steps.Start("Single binary entry point")
	if err := writeSingleMainGo(root, opts); err != nil {
		return fmt.Errorf("writing single main.go: %w", err)
	}

	// Run go mod tidy at project root
	steps.Start("Go dependencies")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = opts.APIRoot(root)
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running go mod tidy: %w\n%s", err, string(out))
	}

	// Write Docker files
	steps.Start("Docker setup")
	if err := writeDockerFiles(root, opts); err != nil {
		return fmt.Errorf("writing Docker files: %w", err)
	}

	// Write frontend, in the framework the flags asked for.
	//
	// Both land in frontend/, which is what makes this a single project: one
	// folder, one Go module at its root, no workspace and no Turborepo. They
	// differ in how they are served, and that difference is forced rather than
	// chosen. The Vite SPA is static files, so the binary embeds it and the whole
	// project ships as one executable. Next.js is not static here: the panel
	// alone has four [id] routes, so `output: "export"` cannot build it, and the
	// app needs its own server like every other Next app.
	if opts.SingleUsesNext() {
		// The shared package, which the Vite single mirrors into its SPA and this
		// one resolves as a workspace member. The frontend's package.json asks for
		// "@repo/shared": "workspace:*", so without it pnpm install fails on a
		// dependency that is not there, before anything is built.
		steps.Start("Shared package")
		if err := writeSharedFiles(root, opts); err != nil {
			return fmt.Errorf("writing shared files: %w", err)
		}

		steps.Start("Next.js frontend")
		if err := writeWebFiles(root, opts); err != nil {
			return fmt.Errorf("writing frontend files: %w", err)
		}

		// The panel as a route group at /admin, which is the shape a double uses:
		// the same screens, moved and repointed. See admin_embedded.go.
		steps.Start("Admin panel at /admin")
		if err := writeEmbeddedAdminFiles(root, opts); err != nil {
			return fmt.Errorf("writing the embedded admin panel: %w", err)
		}
		// No middleware.ts. A static export has no server to run it, and Next
		// fails the build outright rather than ignoring the file. The guard it
		// provided was defence in depth: the panel's own auth already redirects a
		// signed-out visitor, and every endpoint behind it is checked by the API,
		// which is the check that actually matters.
		if !opts.StaticExport() {
			if err := writeAdminEdgeGuard(root, opts); err != nil {
				return err
			}
		}
		return nil
	}

	steps.Start("React frontend")
	if err := writeSingleFrontendFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend files: %w", err)
	}

	// And the admin panel inside it, at /admin. The Vite admin is the source,
	// because this SPA is a TanStack Router app and the panel already exists in
	// that dialect: see admin_embedded_single.go.
	steps.Start("Admin panel at /admin")
	if err := writeEmbeddedSingleAdminFiles(root, opts); err != nil {
		return fmt.Errorf("writing the embedded admin panel: %w", err)
	}

	// Close the last stage, or its running line is the last thing on screen.
	steps.Finish()
	return nil
}

// createSingleDirectories creates the flat directory structure for a single app.
// Note: no cmd/server/ — the single-app main.go lives at the project root so
// the //go:embed all:frontend/dist directive resolves correctly. cmd/migrate
// and cmd/seed remain as separate binaries.
func createSingleDirectories(root string, opts Options) error {
	api := opts.APIRoot(root)
	web := webAppRoot(root, opts)
	dirs := []string{
		filepath.Join(api, "cmd", "migrate"),
		filepath.Join(api, "cmd", "seed"),
		filepath.Join(api, "internal", "config"),
		filepath.Join(api, "internal", "database"),
		filepath.Join(api, "internal", "models"),
		filepath.Join(api, "internal", "handlers"),
		filepath.Join(api, "internal", "middleware"),
		filepath.Join(api, "internal", "services"),
		filepath.Join(api, "internal", "routes"),
		filepath.Join(api, "internal", "mail", "templates"),
		filepath.Join(api, "internal", "storage"),
		filepath.Join(api, "internal", "jobs"),
		filepath.Join(api, "internal", "cron"),
		filepath.Join(api, "internal", "cache"),
		filepath.Join(api, "internal", "ai"),
		filepath.Join(api, "internal", "totp"),
		filepath.Join(api, "internal", "safefetch"),
		filepath.Join(api, "internal", "authz"),
		filepath.Join(web, "public"),
	}

	// src/ is the Vite app's tree. A Next single keeps app/, components/, hooks/
	// and lib/ at the top instead, and creating src/ for it leaves four empty
	// directories that suggest a layout the project does not use.
	if !opts.SingleUsesNext() {
		dirs = append(dirs,
			filepath.Join(web, "src", "routes"),
			filepath.Join(web, "src", "components"),
			filepath.Join(web, "src", "hooks"),
			filepath.Join(web, "src", "lib"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// createDirectories creates the full monorepo folder structure.
func createDirectories(root string, opts Options) error {
	dirs := []string{
		// Go API
		filepath.Join(root, "apps", "api", "cmd", "server"),
		filepath.Join(root, "apps", "api", "cmd", "migrate"),
		filepath.Join(root, "apps", "api", "cmd", "seed"),
		filepath.Join(root, "apps", "api", "internal", "config"),
		filepath.Join(root, "apps", "api", "internal", "database"),
		filepath.Join(root, "apps", "api", "internal", "models"),
		filepath.Join(root, "apps", "api", "internal", "handlers"),
		filepath.Join(root, "apps", "api", "internal", "middleware"),
		filepath.Join(root, "apps", "api", "internal", "services"),
		filepath.Join(root, "apps", "api", "internal", "paginate"),
		filepath.Join(root, "apps", "api", "internal", "realtime"),
		filepath.Join(root, "apps", "api", "internal", "sync"),
		filepath.Join(root, "apps", "api", "internal", "export"),
		filepath.Join(root, "apps", "api", "internal", "respond"),
		filepath.Join(root, "apps", "api", "internal", "pdf"),
		filepath.Join(root, "apps", "api", "internal", "audit"),
		filepath.Join(root, "apps", "api", "internal", "webhooks"),
		filepath.Join(root, "apps", "api", "internal", "flags"),
		filepath.Join(root, "apps", "api", "internal", "routes"),
		filepath.Join(root, "apps", "api", "internal", "mail", "templates"),
		filepath.Join(root, "apps", "api", "internal", "storage"),
		filepath.Join(root, "apps", "api", "internal", "jobs"),
		filepath.Join(root, "apps", "api", "internal", "cron"),
		filepath.Join(root, "apps", "api", "internal", "cache"),
		filepath.Join(root, "apps", "api", "internal", "ai"),
		filepath.Join(root, "apps", "api", "internal", "docs"),
		filepath.Join(root, "apps", "api", "internal", "safefetch"),
		filepath.Join(root, "apps", "api", "internal", "authz"),
	}

	if opts.ShouldIncludeWeb() {
		if opts.UseTanStack() {
			dirs = append(dirs,
				filepath.Join(webAppRoot(root, opts), "src", "routes"),
				filepath.Join(webAppRoot(root, opts), "src", "components"),
				filepath.Join(webAppRoot(root, opts), "src", "hooks"),
				filepath.Join(webAppRoot(root, opts), "src", "lib"),
				filepath.Join(webAppRoot(root, opts), "public"),
			)
		} else {
			dirs = append(dirs,
				// No (auth) folders: the web app's sign-in pages come from
				// grit add web-auth, which creates them with the pages. Created
				// here they were five empty folders in every new project.
				filepath.Join(webAppRoot(root, opts), "app"),
				filepath.Join(webAppRoot(root, opts), "lib"),
				filepath.Join(webAppRoot(root, opts), "__tests__"),
			)
		}
	}

	if opts.ShouldIncludeAdmin() {
		if opts.UseTanStack() {
			dirs = append(dirs,
				filepath.Join(root, "apps", "admin", "src", "routes", "_auth"),
				filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard", "resources"),
				filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard", "system"),
				filepath.Join(root, "apps", "admin", "src", "components", "layout"),
				filepath.Join(root, "apps", "admin", "src", "components", "tables"),
				filepath.Join(root, "apps", "admin", "src", "components", "forms", "fields"),
				filepath.Join(root, "apps", "admin", "src", "components", "resource"),
				filepath.Join(root, "apps", "admin", "src", "components", "shared"),
				filepath.Join(root, "apps", "admin", "src", "components", "ui"),
				filepath.Join(root, "apps", "admin", "src", "components", "profile"),
				filepath.Join(root, "apps", "admin", "src", "hooks"),
				filepath.Join(root, "apps", "admin", "src", "lib"),
				filepath.Join(root, "apps", "admin", "src", "resources"),
				filepath.Join(root, "apps", "admin", "public"),
			)
		} else {
			dirs = append(dirs,
				filepath.Join(root, "apps", "admin", "app", "(auth)", "login"),
				filepath.Join(root, "apps", "admin", "app", "(auth)", "sign-up"),
				filepath.Join(root, "apps", "admin", "app", "(auth)", "forgot-password"),
				filepath.Join(root, "apps", "admin", "app", "(auth)", "reset-password"),
				filepath.Join(root, "apps", "admin", "app", "(auth)", "callback"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "dashboard"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "profile"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "resources", "users"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "jobs"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "files"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "cron"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "mail"),
				filepath.Join(root, "apps", "admin", "app", "(dashboard)", "system", "security"),
				filepath.Join(root, "apps", "admin", "components", "layout"),
				filepath.Join(root, "apps", "admin", "components", "tables"),
				filepath.Join(root, "apps", "admin", "components", "forms", "fields"),
				filepath.Join(root, "apps", "admin", "components", "resource"),
				filepath.Join(root, "apps", "admin", "components", "shared"),
				filepath.Join(root, "apps", "admin", "components", "ui"),
				filepath.Join(root, "apps", "admin", "components", "profile"),
				filepath.Join(root, "apps", "admin", "hooks"),
				filepath.Join(root, "apps", "admin", "lib"),
				filepath.Join(root, "apps", "admin", "resources"),
			)
		}
	}

	if opts.ShouldIncludeShared() {
		dirs = append(dirs,
			filepath.Join(root, "packages", "shared", "schemas"),
			filepath.Join(root, "packages", "shared", "types"),
			filepath.Join(root, "packages", "shared", "constants"),
		)
	}

	if opts.ShouldIncludeWeb() || opts.ShouldIncludeAdmin() {
		dirs = append(dirs,
			filepath.Join(root, "e2e"),
		)
	}

	if opts.ShouldIncludeAdmin() {
		dirs = append(dirs,
			filepath.Join(root, "apps", "admin", "__tests__"),
		)
	}

	if opts.ShouldIncludeExpo() {
		dirs = append(dirs,
			filepath.Join(root, "apps", "expo", "app", "(auth)"),
			filepath.Join(root, "apps", "expo", "app", "(tabs)"),
			filepath.Join(root, "apps", "expo", "lib"),
			filepath.Join(root, "apps", "expo", "components"),
			filepath.Join(root, "apps", "expo", "assets"),
		)
	}

	if opts.ShouldIncludeDesktop() {
		dirs = append(dirs,
			// Real path segments, NOT "_auth"/"_app": a leading underscore makes
			// a TanStack pathless layout, so routes/_app/index.tsx would resolve
			// to "/" and collide with routes/index.tsx (the route generator then
			// errors and never writes routeTree.gen.ts). The app links to
			// /app/... and /auth/login, so these must be real segments.
			filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "auth"),
			filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app"),
			filepath.Join(root, "apps", "desktop", "frontend", "src", "components", "layout"),
			filepath.Join(root, "apps", "desktop", "frontend", "src", "components", "ui"),
			filepath.Join(root, "apps", "desktop", "frontend", "src", "lib"),
			filepath.Join(root, "apps", "desktop", "frontend", "src", "hooks"),
			// build/ holds appicon.png, which writeDesktopClientFiles draws.
			// This used to create build/appicon as a directory, which is not
			// what Wails reads and is not what the README describes.
			filepath.Join(root, "apps", "desktop", "build"),
		)
	}

	if opts.ShouldIncludeDocs() {
		dirs = append(dirs,
			filepath.Join(root, "apps", "docs", "app", "api", "search"),
			filepath.Join(root, "apps", "docs", "app", "docs", "[[...slug]]"),
			filepath.Join(root, "apps", "docs", "content", "docs", "api"),
			filepath.Join(root, "apps", "docs", "public"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// writeFile creates a file with the given content. Go files are gofmt'd on the
// way out (see internal/codefmt) so the templates only have to be correct, not
// aligned; anything else is written byte-for-byte.
func writeFile(path, content string) error {
	// An admin file landing inside the web app has its imports and links
	// repointed. Keyed on the destination rather than on a flag, because the
	// panel's screens are written from a dozen different places and every one of
	// them would have had to remember.
	content = embeddedAdminContent(path, content)
	content = embeddedSingleAdminContent(path, content)

	// Record what actually lands on disk, not what the template produced:
	// codefmt reformats Go on the way out, and a hash of the pre-gofmt text
	// would read as an edit the moment anyone looked at the file.
	final := codefmt.File(path, content)
	if _, err := guardedWrite(path, final); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
