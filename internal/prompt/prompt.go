package prompt

import (
	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// RunNewProjectPrompt shows an interactive prompt for project configuration.
// Returns the selected architecture and frontend. Skips prompts for fields
// that are already set (via CLI flags).
//
// All three selects are wrapped in a single huh.Form so the TUI renders
// exactly once per step instead of three separate render cycles. On
// terminals with weak ANSI support (Git Bash / MINGW64), individual
// .Run() calls leave stale frames in the scrollback because cursor-up
// codes don't fully apply. A Form runs one inline render that cleans up
// after itself, producing a single tidy block of output.
func RunNewProjectPrompt(opts *scaffold.Options) error {
	arch := string(opts.Architecture)
	frontend := string(opts.Frontend)
	theme := opts.Theme
	db := opts.DBProvider

	// "full" is a sentinel, not an Architecture enum: it maps to opts.Full,
	// which Normalize() expands to triple + web + admin + docs + expo + desktop.
	needsFrontend := func() bool {
		a := scaffold.Architecture(arch)
		return arch == "full" || a == scaffold.ArchSingle || a == scaffold.ArchDouble || a == scaffold.ArchTriple
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("arch").
				Title("Select architecture").
				Options(
					huh.NewOption("Full — Web + Admin + API + Docs + Expo + Desktop (everything)", "full"),
					huh.NewOption("Triple — Web + Admin + API (Turborepo)", string(scaffold.ArchTriple)),
					huh.NewOption("Double — Web + API (Turborepo)", string(scaffold.ArchDouble)),
					huh.NewOption("Single — Go API + embedded React SPA (one binary)", string(scaffold.ArchSingle)),
					huh.NewOption("API Only — Go API (no frontend)", string(scaffold.ArchAPI)),
					huh.NewOption("Mobile — API + Expo (React Native)", string(scaffold.ArchMobile)),
				).
				Value(&arch),
		).WithHideFunc(func() bool { return opts.Architecture != "" }),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("frontend").
				Title("Select frontend framework").
				Options(
					huh.NewOption("Next.js — SSR, SEO, App Router", string(scaffold.FrontendNext)),
					huh.NewOption("TanStack Router — Vite, fast builds, small bundle (SPA)", string(scaffold.FrontendTanStack)),
				).
				Value(&frontend),
		).WithHideFunc(func() bool {
			return opts.Frontend != "" || !needsFrontend()
		}),

		// Theme picker — runs for any architecture that includes a frontend.
		// The choice writes THEME=<name> to .env so the dashboard and auth
		// pages render with matching tokens, fonts and layout.
		//
		// Every theme in scaffold.ValidThemes belongs here. The list started
		// at three and grew to eight, and the five that arrived in v3.322.0
		// were reachable only with --theme: the picker is what most people
		// see, so for them the new layouts may as well not have shipped.
		// TestThePickerOffersEveryTheme keeps the two lists together.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("theme").
				Title("Select visual theme").
				Description("Drives auth layout, dashboard tokens, fonts, and brand colors.").
				Options(
					huh.NewOption("Atlas — split-screen, blue/white, team/organisation (Inter)", "atlas"),
					huh.NewOption("Aurora — Apple-inspired, monochrome black/white/grey (Geist)", "aurora"),
					huh.NewOption("Pulse — Cloudflare-inspired, premium blue, elevated cards (Onest)", "pulse"),
					huh.NewOption("Coral — a card over a blurred glimpse of the app, rose (Inter)", "coral"),
					huh.NewOption("Amber — a plain boxed form under a wordmark, amber (Inter)", "amber"),
					huh.NewOption("Sky — one bold heading, social sign-in first, crisp blue (Inter)", "sky"),
					huh.NewOption("Mono — a fine grid beside a panel of proof, black/white (Inter)", "mono"),
					huh.NewOption("Emerald — a form beside a customer quote, green (Inter)", "emerald"),
				).
				Value(&theme),
		).WithHideFunc(func() bool {
			return opts.Theme != "" || !needsFrontend()
		}),

		// The database. Asked of every architecture, because the API is the one
		// thing every shape has.
		//
		// Postgres stays the default and the first option, but it needs a
		// server: without one, a new project migrates into nothing and the
		// first thing anybody sees is a connection refused. SQLite is the
		// shortest path to a running app and it was reachable only by reading
		// the help text for --db.
		//
		// Options come from scaffold.DBProviders, so this list and the flag
		// cannot drift. TestThePickerOffersEveryDatabase keeps them together.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("db").
				Title("Select database").
				Description("Written to .env as DB_PROVIDER. Change it later by editing that line.").
				Options(databaseOptions()...).
				Value(&db),
		).WithHideFunc(func() bool { return opts.DBProvider != "" }),
	)

	// One theme for every picker, from internal/ui: huh's default draws a
	// coloured border down the whole form and paints the options in its own
	// palette, which is a second design living beside the one the rest of the
	// output follows.
	if err := form.WithTheme(ui.PromptTheme()).Run(); err != nil {
		return err
	}

	if arch == "full" {
		// Full is a shorthand, not an architecture. Normalize() turns opts.Full
		// into triple + web + admin + docs + expo + desktop.
		opts.Full = true
	} else {
		opts.Architecture = scaffold.Architecture(arch)
	}
	opts.Frontend = scaffold.Frontend(frontend)
	opts.Theme = theme
	opts.DBProvider = db
	return nil
}

// databaseOptions builds the picker's list from the engines the flag accepts.
//
// Labelled "<name> — <what it means>", with the description scaffold.DBProviders
// already holds, so the two say the same thing in both places.
func databaseOptions() []huh.Option[string] {
	labels := map[string]string{
		"postgres": "PostgreSQL",
		"mysql":    "MySQL / MariaDB",
		"sqlite":   "SQLite",
		"memory":   "In-memory",
	}
	options := make([]huh.Option[string], 0, len(scaffold.DBProviderOrder))
	for _, name := range scaffold.DBProviderOrder {
		label := labels[name]
		if label == "" {
			label = name
		}
		options = append(options, huh.NewOption(label+" — "+scaffold.DBProviders[name], name))
	}
	return options
}
