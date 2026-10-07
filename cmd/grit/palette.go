package main

import "github.com/fatih/color"

// The colours the commands that have not moved to internal/ui yet still use.
//
// internal/ui is the real palette: seven roles, each a dark-terminal value and
// a light-terminal one, from cli-design/cli-design-style-guide.md. These are
// the fatih/color equivalents, kept in step with it, for the thirty-odd call
// sites that still print with .Printf rather than rendering a component.
//
// They were bright magenta. Magenta reads as decoration rather than
// information, and a purple wordmark is the house style of every AI product
// shipped in the last two years.
//
// Moving a command over means replacing these with ui.Step, ui.Result,
// ui.WarningBlock and the rest, which also gets it the light-terminal half.
// Until then: colour says something. Green worked, yellow read this, red
// stopped, cyan is a command to copy, grey is detail.
var (
	// cHeading introduces a step. Bold in the terminal's own foreground, which
	// is readable on a light background and a dark one, where a white heading
	// is invisible on one and a bright one competes with the lines under it.
	cHeading = color.New(color.Bold)

	// cBrand is the one accent, matching ui's cobalt as closely as ANSI-16
	// allows: bright blue is the 12 that ui falls back to on a dark terminal.
	cBrand = color.New(color.FgHiBlue, color.Bold)

	cOK   = color.New(color.FgHiGreen, color.Bold)
	cWarn = color.New(color.FgHiYellow)
	cErr  = color.New(color.FgHiRed, color.Bold)

	cCmd   = color.New(color.FgHiCyan)
	cMuted = color.New(color.FgHiBlack)
)
