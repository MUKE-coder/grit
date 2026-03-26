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

// Options holds the scaffolding configuration.
type Options struct {
	ProjectName string
	APIOnly     bool
	IncludeExpo bool
	MobileOnly  bool
	Full        bool
	Style       string
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

// ShouldIncludeWeb returns true if the web app should be scaffolded.
func (o Options) ShouldIncludeWeb() bool {
	return !o.APIOnly && !o.MobileOnly
}

// ShouldIncludeAdmin returns true if the admin panel should be scaffolded.
func (o Options) ShouldIncludeAdmin() bool {
	return !o.APIOnly && !o.MobileOnly
}

// ShouldIncludeShared returns true if the shared package should be scaffolded.
func (o Options) ShouldIncludeShared() bool {
	return !o.APIOnly
}

// ShouldIncludeExpo returns true if the Expo app should be scaffolded.
func (o Options) ShouldIncludeExpo() bool {
	return o.IncludeExpo || o.MobileOnly || o.Full
}

// ShouldIncludeDocs returns true if the docs site should be scaffolded.
func (o Options) ShouldIncludeDocs() bool {
	return o.Full
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

// Run executes the full scaffolding process.
func Run(opts Options) error {
	root := opts.ProjectName

	if _, err := os.Stat(root); err == nil {
		return fmt.Errorf("directory %q already exists", root)
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
		// Write Next.js web app
		spinner.Printf("  → Scaffolding Next.js web app...\n")
		if err := writeWebFiles(root, opts); err != nil {
			return fmt.Errorf("writing web files: %w", err)
		}
	}

	if opts.ShouldIncludeAdmin() {
		// Write admin panel
		spinner.Printf("  → Scaffolding admin panel...\n")
		if err := writeAdminFiles(root, opts); err != nil {
			return fmt.Errorf("writing admin files: %w", err)
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
		dirs = append(dirs,
			filepath.Join(root, "apps", "web", "app"),
			filepath.Join(root, "apps", "web", "lib"),
			filepath.Join(root, "apps", "web", "__tests__"),
		)
	}

	if opts.ShouldIncludeAdmin() {
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
