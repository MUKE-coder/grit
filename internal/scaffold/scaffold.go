package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
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
	Style        string
	InPlace      bool // Scaffold into current directory (grit new .)

	// Deprecated: use Architecture instead. Kept for backward compatibility.
	APIOnly     bool
	IncludeExpo bool
	MobileOnly  bool
	Full        bool
}

// Normalize maps legacy boolean flags to the new Architecture enum.
// Call this after constructing Options from CLI flags.
func (o *Options) Normalize() {
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
var ValidStyles = []string{"default", "modern", "minimal", "glass"}

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

// ShouldIncludeWeb returns true if a web frontend app should be scaffolded (Turborepo web app).
func (o Options) ShouldIncludeWeb() bool {
	return o.Architecture == ArchDouble || o.Architecture == ArchTriple
}

// ShouldIncludeAdmin returns true if the admin panel should be scaffolded.
func (o Options) ShouldIncludeAdmin() bool {
	return o.Architecture == ArchTriple
}

// ShouldIncludeSingleSPA returns true if this is a single-app embedded SPA.
func (o Options) ShouldIncludeSingleSPA() bool {
	return o.Architecture == ArchSingle
}

// ShouldUseTurborepo returns true if the project uses a Turborepo monorepo.
func (o Options) ShouldUseTurborepo() bool {
	return o.Architecture == ArchDouble || o.Architecture == ArchTriple || o.Architecture == ArchMobile
}

// ShouldIncludeShared returns true if the shared package should be scaffolded.
func (o Options) ShouldIncludeShared() bool {
	return o.ShouldUseTurborepo()
}

// ShouldIncludeFrontend returns true if any frontend (web, admin, or SPA) is included.
func (o Options) ShouldIncludeFrontend() bool {
	return o.Architecture != ArchAPI
}

// ShouldIncludeExpo returns true if the Expo app should be scaffolded.
func (o Options) ShouldIncludeExpo() bool {
	return o.Architecture == ArchMobile || o.IncludeExpo || o.Full
}

// ShouldIncludeDocs returns true if the docs site should be scaffolded.
func (o Options) ShouldIncludeDocs() bool {
	return o.Full
}

// UseTanStack returns true if the frontend uses TanStack Router (Vite).
func (o Options) UseTanStack() bool {
	return o.Frontend == FrontendTanStack
}

// APIRoot returns the base directory for Go API files.
// Single app: project root. Monorepo: apps/api/.
func (o Options) APIRoot(root string) string {
	if o.Architecture == ArchSingle {
		return root
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

	spinner := color.New(color.FgHiBlack)

	// Create directory structure
	spinner.Printf("  → Creating directory structure...\n")
	if err := createDirectories(root, opts); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	// Write root config files
	spinner.Printf("  → Writing configuration files...\n")
	if err := writeRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing root files: %w", err)
	}

	// Write Go API files
	spinner.Printf("  → Scaffolding Go API...\n")
	if err := writeAPIFiles(root, opts); err != nil {
		return fmt.Errorf("writing API files: %w", err)
	}

	// Write migrate and seed entrypoints
	spinner.Printf("  → Adding migration and seed tools...\n")
	if err := writeMigrateSeedFiles(root, opts); err != nil {
		return fmt.Errorf("writing migrate/seed files: %w", err)
	}

	// Write Phase 4 service files (cache, storage, mail, jobs, cron, AI)
	spinner.Printf("  → Adding batteries (cache, storage, mail, jobs, cron, AI)...\n")
	if err := writeCacheFiles(root, opts); err != nil {
		return fmt.Errorf("writing cache files: %w", err)
	}
	if err := writeStorageFiles(root, opts); err != nil {
		return fmt.Errorf("writing storage files: %w", err)
	}
	if err := writeMailFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail files: %w", err)
	}
	if err := writeJobsFiles(root, opts); err != nil {
		return fmt.Errorf("writing jobs files: %w", err)
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

	// Write blog example files
	spinner.Printf("  → Adding blog example...\n")
	if err := writeAPIBlogFiles(root, opts); err != nil {
		return fmt.Errorf("writing blog files: %w", err)
	}

	// Run go mod tidy to resolve dependencies and generate go.sum
	spinner.Printf("  → Resolving Go dependencies...\n")
	apiDir := filepath.Join(root, "apps", "api")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = apiDir
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running go mod tidy: %w\n%s", err, string(out))
	}

	// Write Docker files
	spinner.Printf("  → Creating Docker setup...\n")
	if err := writeDockerFiles(root, opts); err != nil {
		return fmt.Errorf("writing Docker files: %w", err)
	}

	if opts.ShouldIncludeShared() {
		// Write shared package
		spinner.Printf("  → Creating shared package...\n")
		if err := writeSharedFiles(root, opts); err != nil {
			return fmt.Errorf("writing shared files: %w", err)
		}

		// Write Grit UI component registry
		spinner.Printf("  → Creating Grit UI component registry...\n")
		if err := writeGritUIFiles(root, opts); err != nil {
			return fmt.Errorf("writing Grit UI files: %w", err)
		}
	}

	if opts.ShouldIncludeWeb() {
		if opts.UseTanStack() {
			spinner.Printf("  → Scaffolding TanStack Router web app (Vite)...\n")
			if err := writeWebTanStackFiles(root, opts); err != nil {
				return fmt.Errorf("writing TanStack web files: %w", err)
			}
		} else {
			spinner.Printf("  → Scaffolding Next.js web app...\n")
			if err := writeWebFiles(root, opts); err != nil {
				return fmt.Errorf("writing web files: %w", err)
			}
		}
	}

	if opts.ShouldIncludeAdmin() {
		if opts.UseTanStack() {
			spinner.Printf("  → Scaffolding TanStack Router admin panel (Vite)...\n")
			if err := writeAdminTanStackFiles(root, opts); err != nil {
				return fmt.Errorf("writing TanStack admin files: %w", err)
			}
		} else {
			spinner.Printf("  → Scaffolding admin panel...\n")
			if err := writeAdminFiles(root, opts); err != nil {
				return fmt.Errorf("writing admin files: %w", err)
			}
		}
	}

	if opts.ShouldIncludeExpo() {
		// Write Expo mobile app
		spinner.Printf("  → Scaffolding Expo mobile app...\n")
		if err := writeExpoFiles(root, opts); err != nil {
			return fmt.Errorf("writing Expo files: %w", err)
		}
	}

	if opts.ShouldIncludeDocs() {
		// Write docs site
		spinner.Printf("  → Scaffolding documentation site...\n")
		if err := writeDocsFiles(root, opts); err != nil {
			return fmt.Errorf("writing docs files: %w", err)
		}
	}

	// Write frontend test files (Vitest + Playwright)
	if opts.ShouldIncludeWeb() || opts.ShouldIncludeAdmin() {
		spinner.Printf("  → Scaffolding frontend tests (Vitest + Playwright)...\n")
		if err := writeFrontendTestFiles(root, opts); err != nil {
			return fmt.Errorf("writing frontend test files: %w", err)
		}
	}

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

	spinner := color.New(color.FgHiBlack)

	// Create directory structure
	spinner.Printf("  → Creating directory structure...\n")
	if err := createSingleDirectories(root); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	// Write root config files (.env, .gitignore, skill file, Makefile)
	spinner.Printf("  → Writing configuration files...\n")
	if err := writeSingleRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing root files: %w", err)
	}

	// Write .env file (reuse from root_files but with single-app paths)
	if err := writeRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing env files: %w", err)
	}

	// Write Go API files (uses opts.APIRoot which returns root for single)
	spinner.Printf("  → Scaffolding Go API...\n")
	if err := writeAPIFiles(root, opts); err != nil {
		return fmt.Errorf("writing API files: %w", err)
	}

	// Write migrate/seed tools
	spinner.Printf("  → Adding migration and seed tools...\n")
	if err := writeMigrateSeedFiles(root, opts); err != nil {
		return fmt.Errorf("writing migrate/seed files: %w", err)
	}

	// Write batteries
	spinner.Printf("  → Adding batteries (cache, storage, mail, jobs, cron, AI, TOTP)...\n")
	if err := writeCacheFiles(root, opts); err != nil {
		return fmt.Errorf("writing cache files: %w", err)
	}
	if err := writeStorageFiles(root, opts); err != nil {
		return fmt.Errorf("writing storage files: %w", err)
	}
	if err := writeMailFiles(root, opts); err != nil {
		return fmt.Errorf("writing mail files: %w", err)
	}
	if err := writeJobsFiles(root, opts); err != nil {
		return fmt.Errorf("writing jobs files: %w", err)
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

	// Write blog example
	spinner.Printf("  → Adding blog example...\n")
	if err := writeAPIBlogFiles(root, opts); err != nil {
		return fmt.Errorf("writing blog files: %w", err)
	}

	// Write embed-aware main.go (replaces the standard cmd/server/main.go)
	spinner.Printf("  → Writing single-app main.go with go:embed...\n")
	if err := writeSingleMainGo(root, opts); err != nil {
		return fmt.Errorf("writing single main.go: %w", err)
	}

	// Run go mod tidy at project root
	spinner.Printf("  → Resolving Go dependencies...\n")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = root
	if out, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running go mod tidy: %w\n%s", err, string(out))
	}

	// Write Docker files
	spinner.Printf("  → Creating Docker setup...\n")
	if err := writeDockerFiles(root, opts); err != nil {
		return fmt.Errorf("writing Docker files: %w", err)
	}

	// Write frontend
	spinner.Printf("  → Scaffolding React frontend (Vite + TanStack Router)...\n")
	if err := writeSingleFrontendFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend files: %w", err)
	}

	return nil
}

// createSingleDirectories creates the flat directory structure for a single app.
func createSingleDirectories(root string) error {
	dirs := []string{
		filepath.Join(root, "cmd", "server"),
		filepath.Join(root, "cmd", "migrate"),
		filepath.Join(root, "cmd", "seed"),
		filepath.Join(root, "internal", "config"),
		filepath.Join(root, "internal", "database"),
		filepath.Join(root, "internal", "models"),
		filepath.Join(root, "internal", "handlers"),
		filepath.Join(root, "internal", "middleware"),
		filepath.Join(root, "internal", "services"),
		filepath.Join(root, "internal", "routes"),
		filepath.Join(root, "internal", "mail", "templates"),
		filepath.Join(root, "internal", "storage"),
		filepath.Join(root, "internal", "jobs"),
		filepath.Join(root, "internal", "cron"),
		filepath.Join(root, "internal", "cache"),
		filepath.Join(root, "internal", "ai"),
		filepath.Join(root, "internal", "totp"),
		filepath.Join(root, "frontend", "src", "routes"),
		filepath.Join(root, "frontend", "src", "components"),
		filepath.Join(root, "frontend", "src", "hooks"),
		filepath.Join(root, "frontend", "src", "lib"),
		filepath.Join(root, "frontend", "public"),
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
		filepath.Join(root, "apps", "api", "internal", "routes"),
		filepath.Join(root, "apps", "api", "internal", "mail", "templates"),
		filepath.Join(root, "apps", "api", "internal", "storage"),
		filepath.Join(root, "apps", "api", "internal", "jobs"),
		filepath.Join(root, "apps", "api", "internal", "cron"),
		filepath.Join(root, "apps", "api", "internal", "cache"),
		filepath.Join(root, "apps", "api", "internal", "ai"),
		filepath.Join(root, "apps", "api", "internal", "docs"),
	}

	if opts.ShouldIncludeWeb() {
		if opts.UseTanStack() {
			dirs = append(dirs,
				filepath.Join(root, "apps", "web", "src", "routes"),
				filepath.Join(root, "apps", "web", "src", "components"),
				filepath.Join(root, "apps", "web", "src", "hooks"),
				filepath.Join(root, "apps", "web", "src", "lib"),
				filepath.Join(root, "apps", "web", "public"),
			)
		} else {
			dirs = append(dirs,
				filepath.Join(root, "apps", "web", "app"),
				filepath.Join(root, "apps", "web", "lib"),
				filepath.Join(root, "apps", "web", "__tests__"),
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
				filepath.Join(root, "apps", "admin", "src", "components", "widgets"),
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
				filepath.Join(root, "apps", "admin", "components", "widgets"),
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
			filepath.Join(root, "packages", "grit-ui", "registry"),
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

// writeFile creates a file with the given content.
func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}
