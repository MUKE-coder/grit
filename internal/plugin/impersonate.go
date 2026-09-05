package plugin

func init() { Register(impersonatePlugin()) }

// impersonatePlugin lets an admin sign in AS another user to reproduce a bug or
// verify permissions, with a full audit trail and a one-click return to their
// own account.
//
// Design:
//
//   - The swap is server-side. Auth lives in HttpOnly cookies the browser can't
//     read, so starting impersonation re-issues the auth cookies for the target
//     user and stashes the admin's own token in a separate HttpOnly cookie
//     (grit_impersonator). Stop reads that cookie and restores the admin. The
//     admin never handles a raw token.
//
//   - A second, non-HttpOnly cookie (grit_impersonating) carries only the
//     display name/email so the UI can show a persistent banner — HttpOnly
//     cookies are invisible to JS, so without this the browser couldn't tell it
//     was impersonating.
//
//   - Both start and stop write an activity-log entry, so impersonation is never
//     silent. Starting is logged at "warn".
//
//   - Stop is mounted on the protected group, not the admin group: the caller is
//     currently the impersonated user, who may not be an admin, and must still
//     be able to return.
func impersonatePlugin() Plugin {
	return Plugin{
		Name:    "impersonate",
		Version: "1.0.0",
		Summary: "Admin signs in as another user, with an audit trail",
		Description: `Lets an admin impersonate another user to reproduce a bug or check what that
user can see — then return to their own account in one click.

  • Server-side session swap via HttpOnly cookies — the admin never handles a
    raw token, and the original session is restored on stop
  • A persistent banner while impersonating, with a "Return to your account"
    button
  • Every start and stop is written to the activity log
  • An Impersonate screen under System, gated on users.edit

Requires the roles system (Grit v3.66.0+).`,

		NextSteps: []string{
			"Rebuild the API:  cd apps/api && go build ./...",
			"Open the admin → System → Impersonate to pick a user",
			"Impersonation is logged under System → User Activity",
		},

		Files:      impersonateFiles,
		Injections: impersonateInjections,
	}
}

func impersonateFiles(ctx Context) map[string]string {
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
		p("internal/handlers/impersonate.go"): impersonateHandlerGo(ctx),
	}

	// Frontend surfaces only exist where there's an admin app.
	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		files[adminFile(ctx, "hooks/use-impersonate.ts")] = adminSource(ctx, impersonateHook())
		files[adminFile(ctx, "components/impersonation-banner.tsx")] = adminSource(ctx, impersonateBanner())

		// The page, plus the route shim TanStack Router needs to reach it.
		page, route := adminPageFile(ctx, "system/impersonate", "system-impersonate")
		files[page] = adminSource(ctx, impersonatePage())
		if route != "" {
			files[route] = tanStackRouteShim("system/impersonate", "@/pages/system-impersonate")
		}
	}

	return files
}

func impersonateInjections(ctx Context) []Injection {
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
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\timpersonateHandler := handlers.NewImpersonateHandler(db, authService)",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:routes:admin",
			// The admin group is r.Group("/api"), so the /admin/ segment is part
			// of the path, not the group — this resolves to /api/admin/impersonate/:id.
			Code: "\t\tadmin.POST(\"/admin/impersonate/:id\", impersonateHandler.Start)",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code:   "\t\tprotected.POST(\"/auth/impersonate/stop\", impersonateHandler.Stop)",
		},
	}

	// Wire the banner + nav link into the admin app.
	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		admin := adminDir(ctx)
		injections = append(injections,
			Injection{
				File:   admin + "/components/layout/admin-layout.tsx",
				Marker: "// grit:layout:imports",
				Code:   `import { ImpersonationBanner } from "@/components/impersonation-banner";`,
			},
			Injection{
				File:   admin + "/components/layout/admin-layout.tsx",
				Marker: "{/* grit:layout:banner */}",
				Code:   "        <ImpersonationBanner />",
			},
			Injection{
				File:   admin + "/components/chrome/CollapsibleSidebar.tsx",
				Marker: "// grit:nav:system",
				Code:   `  { href: "/system/impersonate", label: "Impersonate", iconKey: "UserCheck", adminOnly: true, requires: "users.edit" },`,
			},
		)
	}

	return injections
}
