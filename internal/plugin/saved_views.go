package plugin

func init() { Register(savedViewsPlugin()) }

// savedViewsPlugin lets each user save a table's current filters, sort, search
// and date range as a named view they can return to in one click.
//
// The elegant part is what it does NOT need. Grit's resource tables already sync
// their entire state to the URL query string, so a "saved view" is just a saved
// query string, per user, per resource. Applying one is a navigation. No changes
// to the DataTable, no per-column serialization — the plugin reads and writes the
// URL the table already reads.
//
// This is the equivalent of what some admin ecosystems sell as a paid "advanced
// tables" add-on; here it's a small plugin.
func savedViewsPlugin() Plugin {
	return Plugin{
		Name:    "saved-views",
		Version: "1.0.0",
		Summary: "Save a table's filters, sort and columns as named views",
		Description: `Save the current state of any resource table — filters, sort, search, date
range — as a named view, and switch between views with one click.

  • Per-user, per-resource: your views are yours
  • A "Save view" button and a row of view chips above every resource table
  • Built on the URL state the tables already use, so nothing in the table
    itself changes

Requires an admin app (triple or full).`,

		NextSteps: []string{
			"Run migrations:  cd apps/api && go run cmd/migrate/main.go",
			"Open any resource, set some filters, and click \"Save view\"",
		},

		Files:      savedViewsFiles,
		Injections: savedViewsInjections,
	}
}

func savedViewsFiles(ctx Context) map[string]string {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	files := map[string]string{
		p("internal/models/saved_view.go"):   savedViewModelGo(ctx),
		p("internal/handlers/saved_view.go"): savedViewHandlerGo(ctx),
	}

	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		files[adminFile(ctx, "hooks/use-saved-views.ts")] = adminSource(ctx, savedViewsHook())
		files[adminFile(ctx, "components/saved-views.tsx")] = adminSource(ctx, savedViewsComponent())
	}

	return files
}

func savedViewsInjections(ctx Context) []Injection {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	injections := []Injection{
		{
			File:   p("internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&SavedView{},",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\tsavedViewHandler := handlers.NewSavedViewHandler(db)",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: `		// Saved table views — scoped to the signed-in user.
		protected.GET("/saved-views", savedViewHandler.List)
		protected.POST("/saved-views", savedViewHandler.Create)
		protected.DELETE("/saved-views/:id", savedViewHandler.Delete)`,
		},
	}

	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		admin := adminDir(ctx)
		injections = append(injections,
			Injection{
				File:   admin + "/components/resource/resource-page.tsx",
				Marker: "// grit:resource:imports",
				Code:   `import { SavedViews } from "@/components/saved-views";`,
			},
			Injection{
				File:   admin + "/components/resource/resource-page.tsx",
				Marker: "{/* grit:table:toolbar */}",
				Code:   "        <SavedViews resource={resource.slug} />",
			},
		)
	}

	return injections
}
