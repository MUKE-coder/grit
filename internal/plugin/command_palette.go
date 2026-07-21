package plugin

func init() { Register(commandPalettePlugin()) }

// commandPalettePlugin adds a ⌘K / Ctrl-K command palette to the admin: a
// keyboard-first, fuzzy-searchable overlay to jump to any resource or system
// page and run quick actions ("New <Resource>").
//
// Deliberately frontend-only — it touches no Go at all. That's the point: it
// proves a Grit plugin can be pure client code, mounted via the layout markers,
// with nothing to migrate and nothing on the server.
//
// The admin already ships a floating QuickAccess button (a grid of cards); this
// is the complementary keyboard path power users reach for, built from the same
// resource registry so it stays in sync automatically.
func commandPalettePlugin() Plugin {
	return Plugin{
		Name:    "command-palette",
		Version: "1.0.0",
		Summary: "⌘K command palette to jump anywhere in the admin",
		Description: `A keyboard-first command palette for the admin panel.

  • Press ⌘K (macOS) or Ctrl+K to open from anywhere
  • Fuzzy-search every resource and system page, then Enter to jump
  • Quick "New <Resource>" actions for each registered resource
  • Frontend only — no Go, no migrations, no server changes

Built from the resource registry, so new resources appear automatically.`,

		NextSteps: []string{
			"Rebuild the admin:  cd apps/admin && pnpm build",
			"Press Ctrl+K (or ⌘K) anywhere in the admin to open it",
		},

		Files:      commandPaletteFiles,
		Injections: commandPaletteInjections,
	}
}

func commandPaletteFiles(ctx Context) map[string]string {
	// Admin-only: nothing to install where there's no admin app.
	if ctx.Architecture != "triple" && ctx.Architecture != "full" {
		return map[string]string{}
	}
	return map[string]string{
		"apps/admin/components/command-palette.tsx": commandPaletteComponent(),
	}
}

func commandPaletteInjections(ctx Context) []Injection {
	if ctx.Architecture != "triple" && ctx.Architecture != "full" {
		return []Injection{}
	}
	admin := "apps/admin"
	return []Injection{
		{
			File:   admin + "/components/layout/admin-layout.tsx",
			Marker: "// grit:layout:imports",
			Code:   `import { CommandPalette } from "@/components/command-palette";`,
		},
		{
			File:   admin + "/components/layout/admin-layout.tsx",
			Marker: "{/* grit:layout:banner */}",
			Code:   "        <CommandPalette />",
		},
	}
}
