package scaffold

import (
	"strings"
	"testing"
)

// The browser-origins field is offered exactly where it is sent.
//
// Origins are a browser-only control, and only a publishable key is meant for
// a browser, so the field is hidden for a secret key. That makes two guards
// that have to agree: the one that draws the field and the one that puts it in
// the request body. A field shown and not sent is a setting somebody will
// swear they configured; a field sent and not shown is a restriction nobody
// can see or remove.
func TestOriginsAreOfferedExactlyWhereTheyAreSent(t *testing.T) {
	page := adminAPIKeysPage()

	if !strings.Contains(page, `if (kind === "publishable" && lines(origins).length) body.origins = lines(origins);`) {
		t.Error("submit() does not guard the origins field on the key type; a value typed " +
			"before switching to Secret would be sent with a key that cannot use it")
	}

	// The field itself, behind the same condition.
	if !strings.Contains(page, `{kind === "publishable" && (`) {
		t.Error("the browser-origins field is not behind a publishable-only condition")
	}

	// And nothing else in the form reads origins outside that block, which is
	// what would make the two drift apart again.
	if strings.Count(page, `setOrigins(e.target.value)`) != 1 {
		t.Error("the origins textarea appears more than once, so one copy can be shown " +
			"for a key type that never sends it")
	}
}

// A secret key is shown in a dialog, and the dialog can be left.
//
// It is a dialog rather than a panel because of where the eye is: the Create
// button sits at the bottom of a tall form, so a panel rendered above that
// form arrives outside the viewport and the key is lost without anybody
// seeing it. There is no second chance, because the server stores a hash.
//
// But a dialog you cannot leave is a keyboard trap, and the people it traps
// are exactly the people using a keyboard. So Escape works, and until the key
// has been copied it asks once rather than refusing.
func TestTheSecretKeyDialogIsNotAKeyboardTrap(t *testing.T) {
	page := adminAPIKeysPage()

	for _, want := range []string{
		`role="dialog"`,
		`aria-modal="true"`,
		`aria-labelledby="reveal-title"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("the reveal dialog is missing %s", want)
		}
	}

	if !strings.Contains(page, `if (e.key === "Escape")`) {
		t.Error("Escape does not close the reveal dialog, which makes it a keyboard trap")
	}
	if !strings.Contains(page, `if (e.key !== "Tab") return;`) {
		t.Error("the reveal dialog does not keep Tab inside itself, which aria-modal promises")
	}
	// The guard before losing the key for good.
	if !strings.Contains(page, "Close without saving it? The key is gone for good.") {
		t.Error("closing the dialog with the key uncopied does not warn, and there is no " +
			"way to recover a secret key")
	}
	// Focus has to land somewhere inside, or a screen reader stays on the page
	// behind the dialog.
	if !strings.Contains(page, "copyButton.current?.focus();") {
		t.Error("the reveal dialog does not move focus into itself when it opens")
	}
}

// Every key in the list is reachable from this page.
//
// "Revoked" was the obvious name for the second filter and it would have been
// a lie: an expired key is just as dead and was never revoked, so under an
// Active/Revoked pair it would have belonged to neither tab and been
// unreachable. The filter is Active / Inactive / All, and the predicate that
// decides the tabs is the same one that decides the row's status pill.
func TestEveryAPIKeyIsReachableFromSomeFilter(t *testing.T) {
	page := adminAPIKeysPage()

	if !strings.Contains(page, `const dead = (k: APIKey) =>`) {
		t.Fatal("the dead-key predicate is gone; the filter counts and the status column " +
			"are about to disagree")
	}
	if !strings.Contains(page, `!!k.revoked_at || (!!k.expires_at && new Date(k.expires_at) < new Date())`) {
		t.Error("the dead-key predicate no longer counts an expired key as dead, so an " +
			"expired key shows as Active or falls out of both filters")
	}
	for _, want := range []string{
		`{ value: "active", label: "Active", count: counts.active },`,
		`{ value: "inactive", label: "Inactive", count: counts.inactive },`,
		`{ value: "all", label: "All", count: counts.all },`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("the keys filter is missing %s", want)
		}
	}
}

// The audit log's method filter and the methods the API actually writes.
//
// The filter is a fixed list of values sent straight to the API, which accepts
// any string and matches it exactly. So a value the API never writes is a tab
// that is always empty, and a value the API writes but the tab list omits is a
// class of entry nobody can isolate. SECURITY is the one that is not an HTTP
// method: the prune job and a reseal both append one under that name.
func TestTheAuditMethodFilterMatchesWhatTheAPIWrites(t *testing.T) {
	page := adminAuditPage()
	api := apiAuditGo()

	if !strings.Contains(page, `{ value: "SECURITY", label: "Security" },`) {
		t.Error("the audit log has no filter for the SECURITY entries the prune job and " +
			"a reseal append, so they cannot be isolated")
	}
	if !strings.Contains(api, `Method:      "SECURITY",`) {
		t.Error("the API no longer writes SECURITY entries, so the audit log's Security " +
			"filter is a tab that can never have anything in it")
	}

	// Every other value in the list is a method the audit middleware records.
	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
		if !strings.Contains(page, `{ value: "`+method+`", label: "`+method+`" },`) {
			t.Errorf("the audit log cannot filter to %s", method)
		}
	}
}

// Both pages speak in the shared primitives rather than in opacities.
//
// Each of them had its own set of hand-written status chips: bg-success/15
// here, bg-success/[0.06] there, border-danger/40 somewhere else. That is how
// an Active key ended up a different green from a healthy component, and it is
// invisible until the two are screenshotted side by side.
func TestTheKeyAndAuditPagesUseTheSharedPrimitives(t *testing.T) {
	pages := map[string]string{
		"API keys":  adminAPIKeysPage(),
		"Audit log": adminAuditPage(),
	}
	// A fraction after a semantic colour is the tell.
	banned := []string{
		"bg-success/", "bg-danger/", "bg-warning/", "bg-info/", "bg-accent/",
		"border-success/", "border-danger/", "border-warning/",
	}
	for name, page := range pages {
		for _, bad := range banned {
			if strings.Contains(page, bad) {
				t.Errorf("%s still mixes its own tint with %q; use the soft and line "+
					"tokens, or the Pill component", name, bad)
			}
		}
		if !strings.Contains(page, `from "@/components/ui/pill"`) {
			t.Errorf("%s does not use the Pill primitive for its statuses", name)
		}
		if !strings.Contains(page, `from "@/components/ui/note"`) {
			t.Errorf("%s does not use the Note primitive for its footnote", name)
		}
		// The old double padding: the layout already pads the content column,
		// and these two pages added 24px more inside it, so their cards sat
		// narrower than every other page's.
		if strings.Contains(page, `className="px-6 py-6`) {
			t.Errorf("%s still wraps its content in its own padding, which insets it "+
				"from every other page in the panel", name)
		}
	}
}

// An empty chain is not a verified chain.
//
// Walking zero entries finds no mismatch, so the endpoint answers valid, and
// the card turned green under a "Verified" badge reading "0 entries checked
// just now". That reads as evidence, and on a log nothing has ever been
// written to there is none. It is the same fault as a chart drawing a slice so
// that it has something to draw: an interface filling a shape it was given
// rather than reporting what it was given.
func TestAnEmptyChainIsNotReportedAsVerified(t *testing.T) {
	page := adminAuditPage()

	if !strings.Contains(page, "const empty = !!result && result.valid && result.total_entries === 0;") {
		t.Fatal("the audit page no longer separates an empty log from a verified chain")
	}
	if !strings.Contains(page, `{empty ? "Nothing to check" : result.valid ? "Verified" : "Broken"}`) {
		t.Error("an empty log still claims the badge of a chain that was checked")
	}
	if !strings.Contains(page, "The log is empty, so there is no chain to verify yet.") {
		t.Error("an empty log does not say why there was nothing to check")
	}
	// And the card itself stays neutral: a green card is the loudest claim on
	// the page and it would be making it about nothing.
	for _, want := range []string{
		`(!result || empty
              ? "border-border bg-bg-elevated"`,
		`(!result || empty
                  ? "bg-accent-soft text-accent"`,
	} {
		if !strings.Contains(page, want) {
			t.Error("the chain card turns green for an empty log")
		}
	}
}
