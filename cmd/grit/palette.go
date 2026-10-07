package main

import "github.com/fatih/color"

// The CLI's colours, in one place.
//
// They were per-command `color.New(...)` calls, and the heading colour was
// bright magenta in thirty of them. Magenta reads as decoration rather than
// information, and a purple wordmark is the house style of every AI product
// shipped in the last two years, which is not the association a Go framework
// wants.
//
// What is here instead: one accent (blue, the same family as the default
// atlas theme), bold-but-uncoloured headings, and colour reserved for saying
// something. Green means it worked, yellow means look at this, red means it
// did not, cyan marks a command you can copy, grey is for detail you can skip.
//
// Bold without a colour for headings is deliberate. A heading painted white is
// invisible on a light terminal and one painted a bright colour competes with
// the status lines under it; bold in the terminal's own foreground is readable
// on every background anybody runs.
var (
	// cBrand is the wordmark and nothing else.
	cBrand = color.New(color.FgHiBlue, color.Bold)

	// cHeading introduces a step: "Creating new Grit project", "Seeding
	// database". Bold, in whatever colour the terminal already uses.
	cHeading = color.New(color.Bold)

	// cOK is a thing that finished. cWarn is a thing worth reading. cErr is a
	// thing that stopped.
	cOK   = color.New(color.FgHiGreen, color.Bold)
	cWarn = color.New(color.FgHiYellow)
	cErr  = color.New(color.FgHiRed, color.Bold)

	// cCmd is something to type or paste. cMuted is detail: ports, paths,
	// versions, the lines somebody skims past.
	cCmd   = color.New(color.FgHiCyan)
	cMuted = color.New(color.FgHiBlack)
)
