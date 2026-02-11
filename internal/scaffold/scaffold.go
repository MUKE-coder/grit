package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Options holds the scaffolding configuration.
type Options struct {
	ProjectName string
	APIOnly     bool
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
	if err := createDirectories(root, opts.APIOnly); err != nil {
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

	// Write Docker files
	spinner.Printf("  → Creating Docker setup...\n")
	if err := writeDockerFiles(root, opts); err != nil {
		return fmt.Errorf("writing Docker files: %w", err)
	}

	if !opts.APIOnly {
		// Write shared package
		spinner.Printf("  → Creating shared package...\n")
		if err := writeSharedFiles(root, opts); err != nil {
			return fmt.Errorf("writing shared files: %w", err)
		}

		// Write Next.js web app
		spinner.Printf("  → Scaffolding Next.js web app...\n")
		if err := writeWebFiles(root, opts); err != nil {
			return fmt.Errorf("writing web files: %w", err)
		}

		// Write admin panel
		spinner.Printf("  → Scaffolding admin panel...\n")
		if err := writeAdminFiles(root, opts); err != nil {
			return fmt.Errorf("writing admin files: %w", err)
		}
	}

	return nil
}

// createDirectories creates the full monorepo folder structure.
func createDirectories(root string, apiOnly bool) error {
	dirs := []string{
		// Go API
		filepath.Join(root, "apps", "api", "cmd", "server"),
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
	}

	if !apiOnly {
		dirs = append(dirs,
			// Next.js web app
			filepath.Join(root, "apps", "web", "app", "(auth)", "login"),
			filepath.Join(root, "apps", "web", "app", "(auth)", "register"),
			filepath.Join(root, "apps", "web", "app", "(auth)", "forgot-password"),
			filepath.Join(root, "apps", "web", "app", "(dashboard)", "dashboard"),
			filepath.Join(root, "apps", "web", "components", "ui"),
			filepath.Join(root, "apps", "web", "components", "shared"),
			filepath.Join(root, "apps", "web", "hooks"),
			filepath.Join(root, "apps", "web", "lib"),

			// Admin panel
			filepath.Join(root, "apps", "admin", "app", "resources", "users"),
			filepath.Join(root, "apps", "admin", "components", "layout"),
			filepath.Join(root, "apps", "admin", "components", "tables"),
			filepath.Join(root, "apps", "admin", "components", "forms", "fields"),
			filepath.Join(root, "apps", "admin", "components", "widgets"),
			filepath.Join(root, "apps", "admin", "hooks"),
			filepath.Join(root, "apps", "admin", "lib"),
			filepath.Join(root, "apps", "admin", "resources"),

			// Shared package
			filepath.Join(root, "packages", "shared", "schemas"),
			filepath.Join(root, "packages", "shared", "types"),
			filepath.Join(root, "packages", "shared", "constants"),
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
