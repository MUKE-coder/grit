package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

// DesktopOptions holds configuration for the grit new-desktop command.
type DesktopOptions struct {
	ProjectName string
}

// RunDesktop scaffolds a standalone Wails desktop application.
func RunDesktop(opts DesktopOptions) error {
	root := opts.ProjectName
	dim := color.New(color.FgHiBlack)

	if _, err := os.Stat(root); err == nil {
		return fmt.Errorf("directory %q already exists", root)
	}

	dim.Printf("  → Creating directory structure...\n")
	if err := createDesktopDirectories(root); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	dim.Printf("  → Writing project configuration...\n")
	if err := writeDesktopRootFiles(root, opts); err != nil {
		return fmt.Errorf("writing root files: %w", err)
	}

	dim.Printf("  → Writing Go backend...\n")
	if err := writeDesktopGoFiles(root, opts); err != nil {
		return fmt.Errorf("writing Go files: %w", err)
	}

	if err := writeDesktopGoInternalFiles(root, opts); err != nil {
		return fmt.Errorf("writing Go internal files: %w", err)
	}

	if err := writeDesktopGoServiceFiles(root, opts); err != nil {
		return fmt.Errorf("writing Go service files: %w", err)
	}

	if err := writeDesktopAPIFiles(root, opts); err != nil {
		return fmt.Errorf("writing API files: %w", err)
	}

	if err := writeDesktopStudioFiles(root, opts); err != nil {
		return fmt.Errorf("writing studio files: %w", err)
	}

	dim.Printf("  → Writing React frontend...\n")
	if err := writeDesktopFrontendConfig(root, opts); err != nil {
		return fmt.Errorf("writing frontend config: %w", err)
	}

	if err := writeDesktopFrontendAppFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend app files: %w", err)
	}

	if err := writeDesktopFrontendLayoutFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend layout files: %w", err)
	}

	if err := writeDesktopFrontendUIFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend UI files: %w", err)
	}

	if err := writeDesktopFrontendPageFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend page files: %w", err)
	}

	if err := writeDesktopFrontendTableFormFiles(root, opts); err != nil {
		return fmt.Errorf("writing frontend table/form files: %w", err)
	}

	dim.Printf("  → Writing in-app auto-updater...\n")
	if err := writeDesktopUpdaterFiles(root, opts); err != nil {
		return fmt.Errorf("writing updater files: %w", err)
	}

	dim.Printf("  → Writing Windows NSIS installers + release script...\n")
	if err := writeDesktopInstallerFiles(root, opts); err != nil {
		return fmt.Errorf("writing installer files: %w", err)
	}

	dim.Printf("  → Running go mod tidy...\n")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = root
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}

func createDesktopDirectories(root string) error {
	dirs := []string{
		// Go backend
		filepath.Join(root, "internal", "config"),
		filepath.Join(root, "internal", "db"),
		filepath.Join(root, "internal", "models"),
		filepath.Join(root, "internal", "service"),
		// Studio
		filepath.Join(root, "cmd", "studio"),
		// Wails build assets
		filepath.Join(root, "build"),
		// Frontend
		filepath.Join(root, "frontend", "src", "components", "ui"),
		filepath.Join(root, "frontend", "src", "components", "layout"),
		filepath.Join(root, "frontend", "src", "routes", "_layout"),
		filepath.Join(root, "frontend", "src", "hooks"),
		filepath.Join(root, "frontend", "src", "lib"),
		// Auto-updater + Windows installer assets
		filepath.Join(root, "scripts"),
		filepath.Join(root, "build", "windows", "installer", "webview2"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	return nil
}
