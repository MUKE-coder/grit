package prompt

import (
	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
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
	)

	if err := form.Run(); err != nil {
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
	return nil
}
