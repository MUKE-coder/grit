package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopRootFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "wails.json"):    desktopWailsJSON(opts),
		filepath.Join(root, "go.mod"):        desktopGoMod(opts),
		filepath.Join(root, ".gitignore"):    desktopGitignore(),
		filepath.Join(root, ".env"):          desktopEnvFile(opts),
		filepath.Join(root, ".env.example"):  desktopEnvExample(opts),
		filepath.Join(root, "README.md"):     desktopReadme(opts),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopWailsJSON(opts DesktopOptions) string {
	return fmt.Sprintf(`{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "%s",
  "outputfilename": "%s",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": ""
  }
}
`, opts.ProjectName, opts.ProjectName)
}

func desktopGoMod(opts DesktopOptions) string {
	return fmt.Sprintf(`module %s

go 1.21

require (
	github.com/wailsapp/wails/v2 v2.9.1
	gorm.io/gorm v1.25.12
	gorm.io/driver/sqlite v1.5.7
	gorm.io/driver/postgres v1.5.11
	golang.org/x/crypto v0.31.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/xuri/excelize/v2 v2.8.1
	github.com/joho/godotenv v1.5.1
)
`, opts.ProjectName)
}

func desktopGitignore() string {
	return `# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/

# Node
node_modules/
frontend/dist/
frontend/src/routeTree.gen.ts

# Wails
build/bin/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Env
.env

# OS
.DS_Store
Thumbs.db
`
}

func desktopEnvFile(opts DesktopOptions) string {
	return fmt.Sprintf(`# Database: sqlite (default) or postgres
DB_DRIVER=sqlite
DB_DSN=%s.db

# For Postgres (uncomment and set):
# DB_DRIVER=postgres
# DB_DSN=host=localhost user=postgres password=postgres dbname=%s port=5432 sslmode=disable

APP_NAME=%s
`, opts.ProjectName, opts.ProjectName, desktopTitleCase(opts.ProjectName))
}

func desktopEnvExample(opts DesktopOptions) string {
	return fmt.Sprintf(`# Database: sqlite (default) or postgres
DB_DRIVER=sqlite
DB_DSN=%s.db

# For Postgres (uncomment and set):
DB_DRIVER=postgres
DB_DSN=host=localhost user=postgres password=postgres dbname=%s port=5432 sslmode=disable

APP_NAME=%s
`, opts.ProjectName, opts.ProjectName, desktopTitleCase(opts.ProjectName))
}

func desktopReadme(opts DesktopOptions) string {
	backtick3 := "```"
	return fmt.Sprintf(`# %s

A desktop application built with [Wails](https://wails.io) and React.

## Prerequisites

- Go 1.21+
- Node.js 18+
- Wails CLI: %sgo install github.com/wailsapp/wails/v2/cmd/wails@latest%s

## Development

%sbash
wails dev
%s

## Build

%sbash
wails build
%s

## Stack

- **Backend**: Go + GORM + SQLite/PostgreSQL
- **Frontend**: React + TypeScript + Tailwind CSS + shadcn-style UI
- **Framework**: Wails v2
`, opts.ProjectName,
		"`", "`",
		backtick3, backtick3,
		backtick3, backtick3)
}

// desktopTitleCase converts a kebab-case or snake_case project name to Title Case.
// For example, "my-cool-app" becomes "My Cool App".
func desktopTitleCase(name string) string {
	// Replace hyphens and underscores with spaces
	s := strings.ReplaceAll(name, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
