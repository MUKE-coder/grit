package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// adminSources returns every admin file the scaffold writes, keyed by path.
func adminSources(t *testing.T) map[string]string {
	t.Helper()
	root := t.TempDir()
	return adminFileMap(root, tripleOptions())
}

// Reduced motion takes the movement, not the information.
//
// The blanket rule runs every animation once for a hundredth of a
// millisecond, which for a spinner means it completes one rotation and stops.
// The panel has ninety of them, so with reduced motion on, every saving
// button, every loading table and every refetching list looked exactly like
// an idle one. That is not a gentler interface, it is one that has stopped
// telling you things.
func TestReducedMotionKeepsTheSpinnersTurning(t *testing.T) {
	css := adminGlobalCSS()

	if !strings.Contains(css, "@media (prefers-reduced-motion: reduce)") {
		t.Fatal("the admin has no reduced-motion block at all")
	}
	// The exemption, and the two it covers.
	for _, want := range []string{
		".animate-spin,",
		".animate-pulse,",
		"[data-grit-progress]",
		"animation-duration: 1.6s !important;",
		"animation-iteration-count: infinite !important;",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("reduced motion does not keep progress indicators running: missing %q", want)
		}
	}
	// And the decoration that should stop.
	if !strings.Contains(css, `button[aria-busy="true"]::after,`) {
		t.Error("the busy sweep is not switched off under reduced motion")
	}
}

// One curve, and the durations the standard allows.
func TestTheAdminHasOneEasingCurve(t *testing.T) {
	css := adminGlobalCSS()

	const curve = "cubic-bezier(0.17, 0.55, 0.55, 1)"
	if !strings.Contains(css, "--default-transition-timing-function: "+curve) {
		t.Error("the default transition curve is not the house curve")
	}
	if !strings.Contains(css, "--ease-grit: "+curve) {
		t.Error("--ease-grit is not the house curve")
	}

	// Any other cubic-bezier in the admin's own CSS is a second curve.
	others := regexp.MustCompile(`cubic-bezier\([^)]*\)`).FindAllString(css, -1)
	for _, c := range others {
		if c != curve {
			t.Errorf("a second easing curve is in the admin CSS: %s", c)
		}
	}

	// Feedback under 300ms, movement 4-8px, press scale 0.97.
	for _, want := range []string{
		"--default-transition-duration: 180ms;",
		"--t-press: 120ms;",
		"--t-enter: 260ms;",
		"translateY(8px) scale(0.97)",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("the motion scale is missing %q", want)
		}
	}
}

// Nothing animates a property that relayouts the page.
//
// transition-all asks the browser to watch every animatable property, so the
// day a utility adds a width the element starts relayouting every frame.
// Nobody writes that bug, it arrives. Width, height, top and left are the
// same cost paid deliberately.
func TestNoAdminFileAnimatesLayout(t *testing.T) {
	banned := regexp.MustCompile(`transition-\[(?:[^\]]*\b(?:width|height|top|left|right|bottom|margin|padding)\b[^\]]*)\]`)

	for path, raw := range adminSources(t) {
		name := filepath.Base(path)
		// Comments explain the rule and would otherwise trip it: globals.css
		// carries a paragraph on why transition-all was replaced.
		body := stripComments(raw)
		if strings.Contains(body, "transition-all") {
			t.Errorf("%s uses transition-all; name the properties instead", name)
		}
		for _, hit := range banned.FindAllString(body, -1) {
			t.Errorf("%s animates a layout property with %q; use a transform", name, hit)
		}
	}
}

// The button covers all eight states.
//
// Error and success are the two everybody skips, and busy is the one that
// gets done wrong: greying out a saving button says "you cannot do this"
// when the truth is "I am doing it".
func TestTheButtonCoversEveryState(t *testing.T) {
	btn := adminButtonPrimitive()

	states := map[string]string{
		"hover":         "hover:bg-accent-hover",
		"focus-visible": "focus-visible:ring-2",
		"active/press":  "motion-safe:active:scale-[0.97]",
		"disabled":      "disabled:opacity-50",
		"busy":          `aria-busy:cursor-progress`,
		"error":         "animate-grit-shake",
		"success":       "!bg-success",
	}
	for name, marker := range states {
		if !strings.Contains(btn, marker) {
			t.Errorf("the button has no %s state (looking for %q)", name, marker)
		}
	}

	// Busy keeps its colour and its place in the tab order.
	if strings.Contains(btn, "aria-busy:opacity-50") {
		t.Error("a busy button is being faded out; it should keep its colour")
	}
	if !strings.Contains(btn, "busyLabel") {
		t.Error("a busy button cannot change its label, so the only signal is a 16px glyph")
	}
}

// A busy button hides the idle icon, without every call site remembering to.
//
// A spinner beside a plus sign is two answers to one question. Doing it in
// the base layer is what lets a call site write <Button loading><Plus/>…</Button>
// and get the right thing, instead of thirty call sites hand-rolling
// {!loading && <Plus/>} and some of them forgetting.
func TestABusyButtonHidesItsIdleIcon(t *testing.T) {
	css := adminGlobalCSS()
	if !strings.Contains(css, `button[aria-busy="true"] > svg:not(.animate-spin)`) {
		t.Error("nothing hides the idle icon on a busy button, so it shows a spinner " +
			"and the original icon side by side")
	}
}

// Copying says so where it happened.
//
// It was a toast in the corner of the screen, which is the furthest point
// from where the person was looking, and a toast is for something that
// happened elsewhere. Copying a key happened right here, because they
// pressed this.
func TestCopyingConfirmsItselfOnTheButton(t *testing.T) {
	copyBtn := adminCopyButtonPrimitive()

	for _, want := range []string{
		`aria-live="polite"`, // the icon swap is invisible to a screen reader
		"Copied",
		"1600", // and it goes back
		"navigator.clipboard.writeText",
	} {
		if !strings.Contains(copyBtn, want) {
			t.Errorf("the copy button is missing %q", want)
		}
	}
	// A clipboard call can be refused by the browser, and a tick that lied is
	// worse than no tick.
	if !strings.Contains(copyBtn, "setFailed(true)") {
		t.Error("the copy button shows success even when the clipboard refused it")
	}

	// And the pages that copy something use it rather than a toast.
	for name, page := range map[string]string{
		"API keys":     adminAPIKeysPage(),
		"Backups":      adminBackupsPage(),
		"Form sharing": adminFormSharesPage(),
	} {
		if !strings.Contains(page, "<CopyButton") {
			t.Errorf("%s still rolls its own copy button", name)
		}
	}
}

// Everything that arrives uses the one entrance.
func TestDialogsShareOneEntrance(t *testing.T) {
	css := adminGlobalCSS()
	if !strings.Contains(css, "--animate-grit-enter: grit-enter var(--t-enter) var(--ease-grit);") {
		t.Fatal("there is no shared entrance animation")
	}

	for name, body := range map[string]string{
		"confirm modal": adminConfirmModal(),
		"key reveal":    adminAPIKeysPage(),
	} {
		if !strings.Contains(body, "animate-grit-enter") {
			t.Errorf("the %s does not use the shared entrance", name)
		}
		// The old per-surface combinations, each with its own duration.
		if strings.Contains(body, "zoom-in-95") {
			t.Errorf("the %s still has its own ad-hoc entrance", name)
		}
	}
}

// stripComments removes /* ... */ and // ... so a rule cannot match the prose
// that explains it. The admin's own CSS documents why transition-all is gone,
// which is exactly the string a check for transition-all looks for.
func stripComments(src string) string {
	block := regexp.MustCompile(`(?s)/\*.*?\*/`)
	line := regexp.MustCompile(`(?m)^\s*//.*$`)
	return line.ReplaceAllString(block.ReplaceAllString(src, ""), "")
}

// The flourish lands on one button, not on six.
//
// The sheen rides on the primary variant, and a page can hold several
// primary buttons: the account page has six, one per card. Lighting all six
// on hover is not a flourish, it is wallpaper. The animation is scoped to the
// page header's action slot, which the panel already treats as the page's one
// primary action.
func TestTheSheenIsScopedToOneButtonPerPage(t *testing.T) {
	css := adminGlobalCSS()
	chrome := adminPageHeaderComponent()
	btn := adminButtonPrimitive()

	if !strings.Contains(css, "[data-page-actions] .grit-sheen:hover::before") {
		t.Error("the sheen is not scoped to the page's action slot, so every primary " +
			"button on a page lights up at once")
	}
	if !strings.Contains(chrome, `data-page-actions=""`) {
		t.Error("the page header no longer marks its action slot, so the sheen has " +
			"nothing to be scoped to and never fires")
	}
	// And the variant must not animate it on its own.
	if strings.Contains(btn, "hover:[&::before]:animate-grit-sheen") {
		t.Error("the primary variant animates the sheen itself, which puts it on every " +
			"primary button on the page")
	}
	if !strings.Contains(btn, `"grit-sheen "`) {
		t.Error("the primary variant no longer positions the sheen element")
	}
}

// A control a thumb has to hit is big enough for a thumb.
//
// The drawn control can stay small: it is the box that grows, and only where
// the pointer is coarse. The copy button was 30x22, which is comfortable with
// a mouse and a miss with a thumb.
func TestSmallControlsGrowOnATouchScreen(t *testing.T) {
	for name, body := range map[string]string{
		"copy button":        adminCopyButtonPrimitive(),
		"dashboard settings": adminDashboardSettingsPageTS(),
		"button (icon size)": adminButtonPrimitive(),
	} {
		if !strings.Contains(body, "pointer-coarse:") {
			t.Errorf("%s has no coarse-pointer floor, so its smallest control stays "+
				"under the touch target size on a phone", name)
		}
	}
}

// The page itself never scrolls sideways.
//
// A wide table scrolls inside its own box; the document does not move. clip
// rather than hidden, because hidden on the root turns position:sticky inside
// the page into position:static and quietly breaks the top bar.
func TestTheRootCannotScrollSideways(t *testing.T) {
	css := adminGlobalCSS()
	if !strings.Contains(css, "overflow-x: clip;") {
		t.Error("the root is not clipped, so anything wider than the viewport can drag " +
			"the whole page sideways on a phone")
	}
	if strings.Contains(css, "html,\nbody {\n  overflow-x: hidden;") {
		t.Error("the root uses overflow-x: hidden, which makes every sticky element " +
			"inside the page static")
	}
}
