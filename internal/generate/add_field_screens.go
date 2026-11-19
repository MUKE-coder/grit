package generate

import (
	"fmt"
	"path/filepath"
)

// The screens, when a field is added.
//
// grit#87 is two gaps, and only one of them was ever the hard one.
//
// The first: grit sync already brings the desktop list's COLUMNS up to the
// model, through the marker fences SyncDesktopScreen writes into, and grit
// generate field did not call it. So the command somebody is told to use for
// adding a field did less than the command for noticing one was missing. That
// is a one-line fix and it is the whole of the desktop half.
//
// The second: nothing updates the Expo form, ever. That is the case in the
// report, found on an Android emulator, where the field was in the API
// response and absent from the form on screen.
//
// # Why the form is rebuilt and the list is not
//
// A field in the Expo form is three things in three places: a useState, an
// entry in the payload, and a control chosen by type. Fencing each would be
// three regions to maintain, and the emission there is one loop over six
// string builders rather than a function that can be called for one field. So
// the form is rebuilt from the definition, which the generator already does
// perfectly.
//
// Whether that is allowed is the manifest's answer, not a new rule: it holds
// the hash of what Grit last wrote, so a form nobody has touched is Grit's to
// rewrite and one somebody has edited is theirs to keep. The same rule grit
// upgrade applies to every framework-owned file.
//
// The Expo list stays reported rather than fixed, and the reason is in
// sync_desktop_screen.go: it renders a header block and a row block that must
// line up, with a TABLE_WIDTH summed from the column widths. Three things that
// have to agree, where an insert that gets one wrong produces a table whose
// headers sit over the wrong columns. A rebuild would be safe by the same
// argument as the form, but only for an untouched file, and the list is the
// screen people customise most. It is left to grit generate resource --force,
// which rebuilds it from the model.

// fieldScreenResult is what happened, for the report.
type fieldScreenResult struct {
	// Updated are screens that now carry the field.
	Updated []string
	// Kept are screens left alone because somebody has edited them. These are
	// the point of the report: the field is in the API and not on them.
	Kept []string
}

// updateFieldScreens brings the screens up to the definition after a field is
// added.
//
// The desktop goes through SyncDesktopScreen, the same code path grit sync
// uses, rather than a second implementation of the same fence insert. Two
// writers for one array is how the two would come to disagree about what a
// column looks like.
func (g *Generator) updateFieldScreens(names Names, f Field) fieldScreenResult {
	var out fieldScreenResult

	// The generator's definition is the resource as it was BEFORE this field,
	// because that is what DefinitionFromModel read at the top of AddField.
	// Rendering a screen from it rebuilds the screen exactly as it was: the
	// report said "rebuilt" and the field was not on it, which is worse than
	// not trying, because it claims the job is done.
	//
	// A copy with the field appended, rather than mutating the generator's
	// own definition: every injection above has already run against that
	// definition, and a later one reading a list that grew underneath it is a
	// bug nobody would find.
	withField := *g.Definition
	withField.Fields = append(append([]Field{}, g.Definition.Fields...), f)
	rendered := *g
	rendered.Definition = &withField

	// The desktop list. Re-parsed from the model, because the field was added
	// to it moments ago and the struct this needs is the one on disk.
	modelPath := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	if structs, err := parseGoStructs(modelPath); err == nil {
		for _, s := range structs {
			if s.Name != names.Pascal {
				continue
			}
			added, warnings, syncErr := SyncDesktopScreen(g.Root, s)
			switch {
			case syncErr != nil:
				out.Kept = append(out.Kept, fmt.Sprintf("apps/desktop: %v", syncErr))
			case added > 0:
				out.Updated = append(out.Updated,
					fmt.Sprintf("apps/desktop: %d column%s added to the %s screen",
						added, plural(added), s.Name))
			}
			for _, warning := range warnings {
				out.Kept = append(out.Kept, warning)
			}
		}
	}

	// The Expo form, through the same writer grit sync uses.
	if updated, kept := syncMobileForm(g.Root, &withField, names); updated != "" {
		out.Updated = append(out.Updated, updated+"  (rebuilt from the model)")
	} else if kept != "" {
		out.Kept = append(out.Kept, kept)
	}

	return out
}

// reportFieldScreens prints what happened, and what did not.
//
// The Kept list is the useful half. Before this, every screen was stale and
// nothing said which, so the output worth having is the short list of the ones
// that still are and why.
func reportFieldScreens(names Names, result fieldScreenResult) {
	for _, label := range result.Updated {
		fmt.Printf("  ✓ %s\n", label)
	}
	if len(result.Kept) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("  ⚠ These screens do not have the field:")
	for _, label := range result.Kept {
		fmt.Printf("      %s\n", label)
	}
	fmt.Println("    Add it by hand, or take Grit's version of the screen with:")
	fmt.Printf("      grit generate resource %s --force\n", names.Pascal)
}
