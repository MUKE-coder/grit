package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// AddOfflineOptions configures `grit add offline`.
type AddOfflineOptions struct {
	// Models to mirror, in plural snake_case. Empty means every model the API
	// has registered with its sync registry.
	Models []string
	// Version of the CLI doing the install, recorded in the manifest.
	Version string
}

// AddOffline installs packages/sync and wires it into whichever clients the
// project has.
//
// The engine is a port of the Go one in apps/desktop, which has been the only
// offline-capable client since v3.60. Nothing about a local mirror, an outbox
// and a version check is desktop-specific, so the second implementation
// targets the interface rather than the platform: the same engine runs on
// IndexedDB in a browser and SQLite on a phone.
func AddOffline(opts AddOfflineOptions) error {
	root, err := FindProjectRoot()
	if err != nil {
		return err
	}

	projectName, err := readProjectName(root)
	if err != nil {
		return err
	}

	hasWeb := dirExists(filepath.Join(root, "apps", "web"))
	hasExpo := dirExists(filepath.Join(root, "apps", "expo"))
	hasAdmin := dirExists(filepath.Join(root, "apps", "admin"))
	if !hasWeb && !hasExpo && !hasAdmin {
		return fmt.Errorf("no client app found: grit add offline needs apps/web, apps/expo or apps/admin")
	}

	spinner := color.New(color.FgHiBlack)
	green := color.New(color.FgHiGreen)
	cyan := color.New(color.FgHiCyan)

	models := opts.Models
	if len(models) == 0 {
		models = discoverSyncModels(root)
		if len(models) == 0 {
			return fmt.Errorf(
				"no syncable models found in internal/routes/routes.go\n\n" +
					"Generate a resource first, or name the models explicitly:\n" +
					"  grit add offline --models products,orders")
		}
	}

	spinner.Printf("  → Writing packages/sync...\n")
	if err := writeSyncPackageFiles(root, Options{ProjectName: projectName, Version: opts.Version}); err != nil {
		return err
	}
	green.Printf("  ✓ packages/sync (engine, 3 storage adapters, React hooks)\n")

	dep := "@" + projectName + "/sync"
	for _, app := range []struct {
		name    string
		present bool
	}{
		{"web", hasWeb},
		{"admin", hasAdmin},
		{"expo", hasExpo},
	} {
		if !app.present {
			continue
		}
		added, err := addWorkspaceDependency(
			filepath.Join(root, "apps", app.name, "package.json"), dep)
		if err != nil {
			spinner.Printf("  Could not add %s to apps/%s: %v\n", dep, app.name, err)
			continue
		}
		if added {
			green.Printf("  ✓ apps/%s depends on %s\n", app.name, dep)
		}
	}

	if hasWeb {
		path := filepath.Join(root, "apps", "web", "lib", "sync.ts")
		created, err := createIfMissing(path, webSyncSetupTS(projectName, models))
		if err != nil {
			return err
		}
		if created {
			green.Printf("  ✓ apps/web/lib/sync.ts\n")
		} else {
			spinner.Printf("  • apps/web/lib/sync.ts already exists, left alone\n")
		}
	}

	if hasExpo {
		path := filepath.Join(root, "apps", "expo", "lib", "sync.ts")
		created, err := createIfMissing(path, expoSyncSetupTS(projectName, models))
		if err != nil {
			return err
		}
		if created {
			green.Printf("  ✓ apps/expo/lib/sync.ts\n")
		} else {
			spinner.Printf("  • apps/expo/lib/sync.ts already exists, left alone\n")
		}
	}

	fmt.Println()
	cyan.Printf("  Mirroring %d model(s): %s\n", len(models), strings.Join(models, ", "))
	fmt.Println()
	cyan.Println("  Next steps:")
	cyan.Println("    pnpm install")
	if hasExpo {
		cyan.Println("    pnpm --filter expo add expo-sqlite")
	}
	cyan.Println()
	cyan.Println("    Wrap your app in <SyncProvider engine={syncEngine}>, then read a")
	cyan.Println("    resource with useOfflineResource(\"products\"). Writes go to the")
	cyan.Println("    outbox and reach the server on the next sync.")
	fmt.Println()

	return nil
}

// discoverSyncModels reads the models the API already registered for sync.
//
// The generator injects syncRegistry.Register("products", &models.Product{})
// for every resource, so the server side is already complete and the client
// only has to name what it wants mirrored. Reading it from routes.go means
// the default is "everything the server can serve" rather than a list that
// silently goes stale.
func discoverSyncModels(root string) []string {
	candidates := []string{
		filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"),
		filepath.Join(root, "internal", "routes", "routes.go"),
	}
	seen := map[string]bool{}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, match := range syncRegisterRe.FindAllStringSubmatch(string(data), -1) {
			if len(match) > 1 && match[1] != "" {
				seen[match[1]] = true
			}
		}
	}
	models := make([]string, 0, len(seen))
	for name := range seen {
		models = append(models, name)
	}
	sort.Strings(models)
	return models
}

// addWorkspaceDependency adds a workspace dependency to a package.json without
// disturbing anything else in it.
//
// Encoding the file back from a decoded map would reorder every key and drop
// the formatting, so this inserts one line into the existing dependencies
// block instead. Reports whether it changed anything.
func addWorkspaceDependency(path, dep string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(data)
	if strings.Contains(content, "\""+dep+"\"") {
		return false, nil
	}

	// Confirm the file parses before editing it. A malformed package.json is
	// the user's to fix, and appending to it would only make it harder.
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return false, fmt.Errorf("parsing %s: %w", path, err)
	}

	const anchor = "\"dependencies\": {"
	idx := strings.Index(content, anchor)
	if idx < 0 {
		return false, fmt.Errorf("no dependencies block in %s", path)
	}
	insertAt := idx + len(anchor)
	line := "\n    \"" + dep + "\": \"workspace:*\","
	updated := content[:insertAt] + line + content[insertAt:]

	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return false, err
	}
	return true, nil
}
