package scaffold

import (
	"strings"
	"testing"
)

// authShells is every sign-in layout, by the file it is written to.
//
// Keyed on the file name so a shell added without being listed here is a shell
// nothing below checks, and the count assertion catches that.
func authShells() map[string]string {
	return map[string]string{
		"AtlasAuthShell.tsx":   adminAtlasAuthShell(),
		"AuroraAuthShell.tsx":  adminAuroraAuthShell(),
		"PulseAuthShell.tsx":   adminPulseAuthShell(),
		"CoralAuthShell.tsx":   adminCoralAuthShell(),
		"AmberAuthShell.tsx":   adminAmberAuthShell(),
		"SkyAuthShell.tsx":     adminSkyAuthShell(),
		"MonoAuthShell.tsx":    adminMonoAuthShell(),
		"EmeraldAuthShell.tsx": adminEmeraldAuthShell(),
	}
}

// Every sign-in layout announces its error.
//
// grit#149 asked for one shell instead of nine, on the grounds that the same
// markup is copied per theme. The eight are genuinely eight layouts, so
// collapsing them is not the fix; what was real is that five of them were built
// from a shared preamble and three predated it, and the shared preamble's
// ErrorBanner carries role="alert" while the three inline copies did not.
//
// So a failed sign-in was not announced to a screen reader on Atlas, which is
// the DEFAULT theme. Exactly the failure the issue predicted: the one that gets
// missed is whichever theme nobody is looking at.
func TestEverySignInLayoutAnnouncesItsError(t *testing.T) {
	shells := authShells()
	if len(shells) != 8 {
		t.Fatalf("%d shells listed, want 8; a new one needs adding to authShells()", len(shells))
	}

	for name, src := range shells {
		// Either through the shared component, or inline with the role on it.
		shared := strings.Contains(src, "<ErrorBanner message={errorMessage} />")
		inline := strings.Contains(src, `role="alert"`)
		if !shared && !inline {
			t.Errorf("%s renders its sign-in error without role=\"alert\", so a screen "+
				"reader does not announce a failed sign-in. Use <ErrorBanner /> from the "+
				"shared preamble, or put the role on the element.", name)
		}
		// And the error has to be rendered at all.
		if !strings.Contains(src, "errorMessage") {
			t.Errorf("%s never renders errorMessage, so a failed sign-in shows nothing", name)
		}
	}
}

// The shared pieces live in one place for the shells that use them, which is
// what makes a change to the cross-links or the token list one edit.
func TestTheSharedShellPartsAreNotCopiedBack(t *testing.T) {
	for name, src := range authShells() {
		// One BrandMark per file, from the preamble. Two means a local copy
		// was added back beside the shared one.
		if n := strings.Count(src, "function BrandMark("); n > 1 {
			t.Errorf("%s declares BrandMark %d times", name, n)
		}
		if n := strings.Count(src, "function ErrorBanner("); n > 1 {
			t.Errorf("%s declares ErrorBanner %d times", name, n)
		}
		// Every layout must publish the full token list, or a form control
		// renders with no border rather than an obvious error.
		if !strings.Contains(src, `["--auth-radius" as string]: theme.radius,`) {
			t.Errorf("%s does not publish --auth-radius, so its inputs lose their corners", name)
		}
		if !strings.Contains(src, `["--auth-primary" as string]: t.primary,`) {
			t.Errorf("%s does not publish --auth-primary", name)
		}
	}
}

// Pulse keeps its own cross-link copy, and that is the reason it was not moved
// onto the shared preamble. Asserted so the difference is a decision somebody
// can read rather than an oversight somebody later "fixes".
func TestPulseKeepsItsOwnCrossLinkCopy(t *testing.T) {
	src := adminPulseAuthShell()
	if !strings.Contains(src, `"Create an account"`) {
		t.Error("Pulse no longer has its own cross-link copy. If that was deliberate, " +
			"it can move onto authShell() like Atlas and Aurora did; if not, its " +
			"sign-up link now reads like every other theme's.")
	}
	// It also needs its own imports for the carousel, which is the other half
	// of why the shared preamble does not fit it.
	if !strings.Contains(src, `import { useEffect, useState } from "react";`) {
		t.Error("Pulse lost the hooks its rotating hero needs")
	}
}
