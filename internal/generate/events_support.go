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

	eventsPath := filepath.Join(apiRoot, "internal", "events", "events.go")
	for _, f := range []struct {
		path string
		body string
		// needs is a symbol the generated code depends on. A file from an
		// earlier release can exist and still be missing it, which fails the
		// build with an undefined reference rather than anything explanatory.
		needs string
		// what names the missing piece in the warning.
		what string
		// requires, when set, is a symbol events.go must have before this file
		// is written: durable.go uses what a current events.go declares, so
		// writing it beside an old one would break the build instead of fixing it.
		requires string
	}{
		{eventsPath, scaffold.APIEventsGo(), "wantDurable", "Durable delivery", ""},
		{filepath.Join(apiRoot, "internal", "events", "durable.go"), scaffold.APIEventsDurableGo(), "func EmitTx(", "EmitTx", "wantDurable"},
		{filepath.Join(apiRoot, "internal", "outbox", "outbox.go"), scaffold.APIOutboxGo(), "func Enqueue(", "Enqueue", ""},
		{filepath.Join(apiRoot, "internal", "outbox", "relay.go"), scaffold.APIOutboxRelayGo(), "TopicPrefix", "the relay's TopicPrefix", ""},
		{filepath.Join(apiRoot, "internal", "models", "outbox_message.go"), scaffold.APIOutboxModelGo(), "type OutboxMessage struct", "OutboxMessage", ""},
		{filepath.Join(apiRoot, "internal", "services", "event_subscribers.go"), scaffold.APIEventsSubscribersGo(), "func RegisterEventSubscribers(", "RegisterEventSubscribers", ""},
	} {
		body := strings.ReplaceAll(f.body, "{{MODULE}}", g.Module)
		rel := strings.TrimPrefix(filepath.ToSlash(f.path), filepath.ToSlash(apiRoot)+"/")
		if f.requires != "" && !fileContainsText(eventsPath, f.requires) {
			continue
		}

		if !fileExists(f.path) {
			if err := writeFileWithDirs(f.path, body); err != nil {
				return fmt.Errorf("writing event bus: %w", err)
			}
			fmt.Printf("  ✓ Added %s\n", rel)
			continue
		}

		current, err := os.ReadFile(f.path)
		if err != nil || strings.Contains(string(current), f.needs) {
			continue
		}

		// Present but too old. Refreshed only when the manifest proves nobody
		// has edited it, on the same rule the upgrade guard follows: no
		// evidence is not permission.
		if !g.refreshIfUnchanged(f.path, body) {
			color.New(color.FgHiYellow).Printf("\n  ⚠ %s is from an earlier release and is missing %s\n", rel, f.what)
			fmt.Println("    The build will fail with an undefined reference. Replace that file")
			fmt.Println("    with the current template, or delete it and generate again.")
			fmt.Println()
			continue
		}
		fmt.Printf("  ✓ Updated %s\n", rel)
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
		// Started before Durable delivery existed: add the relay.
		if changed, err := scaffold.EnsureEventRelay(g.APIRoot()); err != nil {
			color.New(color.FgHiYellow).Printf("\n  ⚠ %v\n\n", err)
		} else if changed {
			fmt.Println("  ✓ Started the event relay in routes.Setup")
		}
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
		fmt.Println("      services.RegisterEventSubscribers(db, realtimeHub, nil)")
		fmt.Println()
		return nil
	}

	insertAt := idx + len(anchor)
	boot := "\n\n\t// Domain event bus: generated handlers emit here, and the audit log,\n" +
		"\t// realtime and (when installed) webhooks subscribe.\n" +
		"\tevents.Init(4)\n" +
		"\tservices.RegisterEventSubscribers(db, realtimeHub, nil)"
	if fileContainsText(filepath.Join(g.APIRoot(), "internal", "events", "durable.go"), "func StartRelay(") {
		boot += "\n\tevents.StartRelay(db)"
	}
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

// fileContainsText reports whether a file exists and contains a string.
func fileContainsText(path, needle string) bool {
	data, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(data), needle)
}
