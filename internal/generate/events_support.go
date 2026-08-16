package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// ensureEventsSupport gives a project the event bus its generated handlers
// now emit on.
//
// Same shape as the dialect helper and the sync policy: grit upgrade does not
// refresh API code, so a resource generated into a project from last month
// would otherwise get a handler calling a package that project has never
// seen. The two files are new, so writing them is safe; the boot wiring is an
// edit to routes.go and is anchored on a line that must already exist.
func (g *Generator) ensureEventsSupport() error {
	apiRoot := g.APIRoot()

	for _, f := range []struct {
		path string
		body string
	}{
		{filepath.Join(apiRoot, "internal", "events", "events.go"), scaffold.APIEventsGo()},
		{filepath.Join(apiRoot, "internal", "events", "subscribers.go"), scaffold.APIEventsSubscribersGo()},
	} {
		if fileExists(f.path) {
			continue
		}
		if err := writeFileWithDirs(f.path, strings.ReplaceAll(f.body, "{{MODULE}}", g.Module)); err != nil {
			return fmt.Errorf("writing event bus: %w", err)
		}
		fmt.Printf("  ✓ Added internal/events/%s\n", filepath.Base(f.path))
	}

	return g.ensureEventsBoot()
}

// ensureEventsBoot starts the bus in routes.Setup.
//
// Without this the package compiles, Emit is a no-op, and nothing subscribes:
// the audit rows that generated handlers used to write would simply stop
// appearing, with no error anywhere. That is a worse outcome than the build
// failure it replaces, so it is worth an anchored edit rather than a warning.
func (g *Generator) ensureEventsBoot() error {
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, "events.Init(") {
		return nil
	}

	const anchor = "realtimeHub := realtime.NewHub()"
	idx := strings.Index(content, anchor)
	if idx < 0 {
		yellow := color.New(color.FgHiYellow)
		yellow.Printf("\n  ⚠ Could not find where routes.go creates the realtime hub, so the\n")
		fmt.Println("    event bus was not started. Generated handlers will emit into nothing")
		fmt.Println("    and activity-feed rows will stop appearing. Add this after the hub:")
		fmt.Println()
		fmt.Println("      events.Init(4)")
		fmt.Println("      events.RegisterDefaults(db, realtimeHub, nil)")
		fmt.Println()
		return nil
	}

	insertAt := idx + len(anchor)
	boot := "\n\n\t// Domain event bus: generated handlers emit here, and the audit log,\n" +
		"\t// realtime and (when installed) webhooks subscribe.\n" +
		"\tevents.Init(4)\n" +
		"\tevents.RegisterDefaults(db, realtimeHub, nil)"
	content = content[:insertAt] + boot + content[insertAt:]

	// The import has to come with it, or the file stops compiling.
	if !strings.Contains(content, "/internal/events\"") {
		importAnchor := "/internal/handlers\""
		if i := strings.Index(content, importAnchor); i >= 0 {
			lineStart := strings.LastIndexByte(content[:i], '\n') + 1
			indent := content[lineStart : strings.Index(content[lineStart:], "\"")+lineStart]
			content = content[:lineStart] +
				indent + "\"" + g.Module + "/internal/events\"\n" +
				content[lineStart:]
		}
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("wiring the event bus into routes.go: %w", err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Started the event bus in routes.Setup")
	return nil
}
