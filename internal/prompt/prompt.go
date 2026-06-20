package prompt

import (
	"github.com/charmbracelet/huh"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// RunNewProjectPrompt shows an interactive prompt for project configuration.
// Returns the selected architecture and frontend. Skips prompts for fields
// that are already set (via CLI flags).
func RunNewProjectPrompt(opts *scaffold.Options) error {
	// Skip if architecture is already set via flags
	if opts.Architecture == "" {
		var arch string
		err := huh.NewSelect[string]().
			Title("Select architecture").
			Options(
				huh.NewOption("Triple — Web + Admin + API (Turborepo)", string(scaffold.ArchTriple)),
				huh.NewOption("Double — Web + API (Turborepo)", string(scaffold.ArchDouble)),
				huh.NewOption("Single — Go API + embedded React SPA (one binary)", string(scaffold.ArchSingle)),
				huh.NewOption("API Only — Go API (no frontend)", string(scaffold.ArchAPI)),
				huh.NewOption("Mobile — API + Expo (React Native)", string(scaffold.ArchMobile)),
			).
			Value(&arch).
			Run()
		if err != nil {
			return err
		}
		opts.Architecture = scaffold.Architecture(arch)
	}

	// Only ask for frontend if the architecture includes one
	needsFrontend := opts.Architecture == scaffold.ArchSingle ||
		opts.Architecture == scaffold.ArchDouble ||
		opts.Architecture == scaffold.ArchTriple

	if needsFrontend && opts.Frontend == "" {
		var frontend string
		err := huh.NewSelect[string]().
			Title("Select frontend framework").
			Options(
				huh.NewOption("Next.js — SSR, SEO, App Router", string(scaffold.FrontendNext)),
				huh.NewOption("TanStack Router — Vite, fast builds, small bundle (SPA)", string(scaffold.FrontendTanStack)),
			).
			Value(&frontend).
			Run()
		if err != nil {
			return err
		}
		opts.Frontend = scaffold.Frontend(frontend)
	}

	// Theme picker — runs for any architecture that includes a frontend.
	// Themes ship since v3.28: atlas (default), aurora (centered), pulse
	// (carousel). The choice writes THEME=<name> to .env so the dashboard
	// + auth pages render with matching tokens, fonts, and layouts.
	if needsFrontend && opts.Theme == "" {
		var theme string
		err := huh.NewSelect[string]().
			Title("Select visual theme").
			Description("Drives auth layout, dashboard tokens, fonts, and brand colors.").
			Options(
				huh.NewOption("Atlas — split-screen, blue/white, team/organisation (Inter)", "atlas"),
				huh.NewOption("Aurora — centered Clerk-style, pastel, consumer SaaS (Geist)", "aurora"),
				huh.NewOption("Pulse — split + hero carousel, bold, ecommerce/brand (Onest + DM Serif)", "pulse"),
			).
			Value(&theme).
			Run()
		if err != nil {
			return err
		}
		opts.Theme = theme
	}

	return nil
}
