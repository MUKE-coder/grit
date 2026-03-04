package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType identifies the kind of Grit project.
type ProjectType string

const (
	ProjectWeb     ProjectType = "web"
	ProjectDesktop ProjectType = "desktop"
)

// ProjectInfo holds detected project metadata.
type ProjectInfo struct {
	Root   string      // absolute path to project root
	Type   ProjectType // web or desktop
	Module string      // Go module path
}

// DetectProject walks up from the current working directory to find a Grit project.
func DetectProject() (*ProjectInfo, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	return DetectProjectFrom(dir)
}

// DetectProjectFrom walks up from the given directory to find a Grit project.
func DetectProjectFrom(startDir string) (*ProjectInfo, error) {
	dir := startDir

	for {
		if IsDesktop(dir) {
			mod, err := readModule(filepath.Join(dir, "go.mod"))
			if err != nil {
				return nil, fmt.Errorf("reading module in desktop project: %w", err)
			}
			return &ProjectInfo{Root: dir, Type: ProjectDesktop, Module: mod}, nil
		}

		if IsWeb(dir) {
			mod, err := readModule(filepath.Join(dir, "apps", "api", "go.mod"))
			if err != nil {
				return nil, fmt.Errorf("reading module in web project: %w", err)
			}
			return &ProjectInfo{Root: dir, Type: ProjectWeb, Module: mod}, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil, fmt.Errorf("not inside a Grit project (no wails.json or turbo.json found)\n\nRun this command from inside a Grit project directory")
}

// IsDesktop returns true if the directory contains a wails.json file.
func IsDesktop(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "wails.json"))
	return err == nil && !info.IsDir()
}

// IsWeb returns true if the directory contains turbo.json and apps/api.
func IsWeb(dir string) bool {
	turbo, err := os.Stat(filepath.Join(dir, "turbo.json"))
	if err != nil || turbo.IsDir() {
		return false
	}
	api, err := os.Stat(filepath.Join(dir, "apps", "api"))
	return err == nil && api.IsDir()
}

func readModule(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", goModPath, err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("no module directive found in %s", goModPath)
}
