package scaffold

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The admin moves the same way everywhere, and the way to keep it that way is
// to make drifting back cost a failing test.
//
// Every one of these was absent before v3.528.0. The panel had no easing
// curve at all, so each transition ran on Tailwind's default; it had
// transition-all in 37 places, which asks the browser to watch every
// animatable property including the layout ones; and prefers-reduced-motion
// appeared nowhere in the entire scaffold.

// One curve and one default duration, set on the variables every Tailwind
// transition-* utility reads. This is what makes the panel feel like one hand
// made it without a duration on two hundred class strings.
func TestTheAdminHasOneEasingCurve(t *testing.T) {
	css := adminScreenCSS()

	const curve = "cubic-bezier(0.17, 0.55, 0.55, 1)"
	if !strings.Contains(css, "--default-transition-timing-function: "+curve) {
		t.Error("the house curve is not on --default-transition-timing-function, " +
			"so every transition-* utility in the admin falls back to Tailwind's own")
	}
	if !strings.Contains(css, "--default-transition-duration:") {
		t.Error("no default duration, so every transition-* utility runs at Tailwind's 150ms")
	}

	// A second curve is the thing this exists to prevent: two is not a system.
	others := regexp.MustCompile(`cubic-bezier\([^)]*\)`).FindAllString(css, -1)
	for _, found := range others {
		if found != curve {
			t.Errorf("a second easing curve, %s. One curve, or the panel is a collection "+
				"of components that each chose a number", found)
		}
	}
}

// The utility that replaced transition-all has to set all three properties.
//
// Setting only transition-property looks right and is not: the duration falls
// back to the CSS initial 0s, so every call site without its own duration-*
// stops transitioning altogether. That was a real regression, caught in a
// browser rather than here, which is why it is here now.
func TestTheTransitionUtilitySetsDurationAndEasing(t *testing.T) {
	css := adminScreenCSS()

	block := betweenBraces(t, css, "@utility transition-grit {")
	for _, want := range []string{
		"transition-property:",
		"transition-timing-function:",
		"transition-duration:",
	} {
		if !strings.Contains(block, want) {
			t.Errorf("transition-grit does not set %s. Tailwind's transition-all set all "+
				"three, so a utility that sets fewer silently removes the animation it "+
				"was meant to fix", strings.TrimSuffix(want, ":"))
		}
	}
	// And it must not animate anything that triggers layout.
	for _, banned := range []string{"width", "height", " top", " left", "margin", "padding"} {
		if strings.Contains(block, banned) {
			t.Errorf("transition-grit animates %q, which relayouts the page every frame",
				strings.TrimSpace(banned))
		}
	}
}

// transition-all, anywhere in the admin's own sources, is the thing the
// utility replaced.
//
// Read off disk rather than by calling every writer, because the writers are
// a hundred functions with different signatures and the class strings are
// what this is about. The files outside the admin are left alone on purpose:
// transition-grit is declared in the admin stylesheet, and apps/web and the
// desktop client include neither, so the class would name a utility Tailwind
// has never heard of and the element would silently stop transitioning.
func TestTheAdminDoesNotUseTransitionAll(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "admin") || !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		checked++
		if strings.Contains(withoutCSSComments(string(raw)), "transition-all") {
			t.Errorf("%s still has transition-all. It watches every animatable property, "+
				"so the day a utility adds a width to that element it animates it too. "+
				"Use transition-grit", name)
		}
	}
	if checked < 20 {
		t.Fatalf("only read %d admin files; this test is looking in the wrong place, "+
			"and a sweep that finds nothing passes", checked)
	}
}

// The off switch, which the scaffold did not have anywhere at all.
//
// Movement goes and colour stays: a hover that changes colour is information
// rather than motion, and the usual blanket "every duration to 0.01ms" takes
// the feedback away along with the discomfort somebody asked to be spared.
func TestTheAdminRespectsReducedMotion(t *testing.T) {
	css := adminScreenCSS()
	if !strings.Contains(css, "@media (prefers-reduced-motion: reduce)") {
		t.Fatal("no reduced-motion block: the admin animates and offers no way to stop it")
	}

	block := betweenBraces(t, css, "@media (prefers-reduced-motion: reduce) {")
	if !strings.Contains(block, "transition-property:") {
		t.Error("reduced motion does not narrow transition-property, so transforms still run")
	}
	if strings.Contains(block, "transform") {
		t.Error("transform is still in the reduced-motion transition list, which is the " +
			"one property that actually moves")
	}
	if !strings.Contains(block, "color") {
		t.Error("reduced motion drops colour transitions too, which removes the feedback " +
			"rather than the movement")
	}
	if !strings.Contains(block, "animation-duration: 0.01ms") {
		t.Error("reduced motion does not stop animations")
	}
}

// Eight states, and the two that are always missing are error and success.
func TestTheButtonHasEveryState(t *testing.T) {
	src := adminButtonPrimitive()

	for what, want := range map[string]string{
		"hover":             "hover:",
		"focus-visible":     "focus-visible:ring-2",
		"the press":         "active:scale-[0.97]",
		"disabled":          "disabled:opacity-50",
		"loading":           "aria-busy",
		"an error outcome":  `error:`,
		"a success outcome": `success:`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the button has no %s state (looked for %q)", what, want)
		}
	}

	// The press has to be motion-safe, or reduced motion still gets a scale.
	if !strings.Contains(src, "motion-safe:active:scale-") {
		t.Error("the press is not motion-safe, so it still moves under reduced motion")
	}

	// Busy is not disabled. Greying a busy button says "you cannot do this"
	// when what happened is "I am doing it".
	if strings.Contains(src, "disabled={disabled || loading}") {
		t.Error("a busy button is still being disabled, which takes it out of the tab " +
			"order: a keyboard user standing on it loses their place and is told nothing")
	}
	if !strings.Contains(src, "busyLabel") {
		t.Error("no busyLabel, so the word never changes and a spinner is the only " +
			"signal. A spinner says nothing to a screen reader")
	}
	// The guard has to cover the keyboard: Enter on a focused button arrives
	// as a click, and pointer-events has no opinion about it.
	if !strings.Contains(src, "event.preventDefault()") {
		t.Error("the busy guard does not preventDefault, so Enter on a busy submit " +
			"button submits the form a second time")
	}
}

// betweenBraces returns the body of the block that opens with head.
func betweenBraces(t *testing.T, css, head string) string {
	t.Helper()
	at := strings.Index(css, head)
	if at < 0 {
		t.Fatalf("no %q in the stylesheet", head)
	}
	rest := css[at+len(head):]
	depth := 1
	for i, r := range rest {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return rest[:i]
			}
		}
	}
	t.Fatalf("%q is never closed", head)
	return ""
}

// withoutCSSComments drops /* */ blocks.
//
// These files are Go holding CSS and JSX, so a class name and the sentence
// explaining it live in the same file. The comment above @utility
// transition-grit says what it replaces, and the first version of this sweep
// read that as a violation: the guard was checking the prose rather than the
// code, which is the same mistake as a guard that passes because of its own
// prose.
//
// Named apart from withoutComments in tailwind_version_guard_test.go, which
// drops // lines instead: the occurrence this has to ignore is a CSS comment
// inside a Go raw string, which that one leaves alone.
func withoutCSSComments(src string) string {
	var out strings.Builder
	for len(src) > 0 {
		if strings.HasPrefix(src, "/*") {
			end := strings.Index(src[2:], "*/")
			if end < 0 {
				break
			}
			src = src[2+end+2:]
			continue
		}
		out.WriteByte(src[0])
		src = src[1:]
	}
	return out.String()
}
