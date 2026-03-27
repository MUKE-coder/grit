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

	// Only ask for style if architecture includes admin panel
	if opts.Architecture == scaffold.ArchTriple && opts.Style == "" {
		var style string
		err := huh.NewSelect[string]().
			Title("Select admin panel style").
			Options(
				huh.NewOption("Default — Clean dark theme", "default"),
				huh.NewOption("Modern — Gradient accents", "modern"),
				huh.NewOption("Minimal — Ultra clean", "minimal"),
				huh.NewOption("Glass — Glassmorphism", "glass"),
			).
			Value(&style).
			Run()
		if err != nil {
			return err
		}
		opts.Style = style
	}

	return nil
}
