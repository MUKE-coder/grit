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

	needsFrontend := func() bool {
		a := scaffold.Architecture(arch)
		return a == scaffold.ArchSingle || a == scaffold.ArchDouble || a == scaffold.ArchTriple
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("arch").
				Title("Select architecture").
				Options(
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
		// Themes ship since v3.28: atlas (default), aurora (centered), pulse
		// (carousel). The choice writes THEME=<name> to .env so the dashboard
		// + auth pages render with matching tokens, fonts, and layouts.
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("theme").
				Title("Select visual theme").
				Description("Drives auth layout, dashboard tokens, fonts, and brand colors.").
				Options(
					huh.NewOption("Atlas — split-screen, blue/white, team/organisation (Inter)", "atlas"),
					huh.NewOption("Aurora — centered Clerk-style, pastel, consumer SaaS (Geist)", "aurora"),
					huh.NewOption("Pulse — split + hero carousel, bold, ecommerce/brand (Onest + DM Serif)", "pulse"),
				).
				Value(&theme),
		).WithHideFunc(func() bool {
			return opts.Theme != "" || !needsFrontend()
		}),
	)

	if err := form.Run(); err != nil {
		return err
	}

	opts.Architecture = scaffold.Architecture(arch)
	opts.Frontend = scaffold.Frontend(frontend)
	opts.Theme = theme
	return nil
}
