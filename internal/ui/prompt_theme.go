package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// PromptTheme is how every grit picker looks.
//
// huh's default draws a coloured left border down the whole form and paints
// option text in its own palette, which is a second design living beside this
// one. This strips the border, uses the two-space gutter the rest of the output
// uses, and spends colour on exactly one thing: the row you are on.
//
// Unselected options are the terminal's own foreground, because a list of
// eight options in a colour is a list nobody reads.
func PromptTheme() *huh.Theme {
	t := huh.ThemeBase()

	gutter := lipgloss.NewStyle().PaddingLeft(2)

	t.Focused.Base = gutter
	t.Focused.Title = Bold
	t.Focused.Description = Muted
	t.Focused.SelectSelector = Brand.Bold(true).SetString(GlyphPointer + " ")
	t.Focused.SelectedOption = Brand.Bold(true)
	t.Focused.UnselectedOption = Text
	t.Focused.SelectedPrefix = lipgloss.NewStyle()
	t.Focused.UnselectedPrefix = lipgloss.NewStyle()

	t.Focused.ErrorIndicator = Error.SetString(" " + GlyphFail)
	t.Focused.ErrorMessage = Error

	t.Help.ShortKey = Muted
	t.Help.ShortDesc = Muted
	t.Help.ShortSeparator = Faint
	t.Help.FullKey = Muted
	t.Help.FullDesc = Muted
	t.Help.FullSeparator = Faint

	// A blurred group is one that has been answered. It is collapsed by the
	// caller rather than redrawn, so it only has to not shout.
	t.Blurred = t.Focused
	t.Blurred.Base = gutter

	return t
}
