package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The account screen, and the one thing that can quietly rot about it: the
// checklist somebody reads while typing a password is written twice, once in
// Go and once in TSX, and if the two drift the page tells people a rule that
// is not enforced or hides one that is.
func TestPasswordChecklistMatchesTheServerRules(t *testing.T) {
	goIDs := idsIn(apiPasswordRulesGo(), regexp.MustCompile(`\{ID: "([a-z-]+)", Label: "([^"]+)"`))
	tsIDs := idsIn(adminAccountPasswordFormTSX(), regexp.MustCompile(`\{ id: "([a-z-]+)", label: "([^"]+)"`))

	if len(goIDs) == 0 || len(tsIDs) == 0 {
		t.Fatalf("found %d server rules and %d checklist rules; one of the lists moved", len(goIDs), len(tsIDs))
	}
	if len(goIDs) != len(tsIDs) {
		t.Fatalf("the server enforces %d rules and the checklist shows %d", len(goIDs), len(tsIDs))
	}
	for id, goLabel := range goIDs {
		tsLabel, ok := tsIDs[id]
		if !ok {
			t.Errorf("the server enforces %q and the checklist never mentions it", id)
			continue
		}
		if goLabel != tsLabel {
			t.Errorf("rule %q reads %q on the server and %q in the admin", id, goLabel, tsLabel)
		}
	}
}

func idsIn(src string, re *regexp.Regexp) map[string]string {
	out := map[string]string{}
	for _, match := range re.FindAllStringSubmatch(src, -1) {
		out[match[1]] = match[2]
	}
	return out
}

// The rules have to be enforced everywhere a password is set, or they are a
// rule about one door in a building with three.
func TestPasswordRulesAreEnforcedOnEveryPath(t *testing.T) {
	register := apiAuthHandlerGo()
	if !strings.Contains(register, "password.Check(req.Password, req.Email, req.FirstName, req.LastName)") {
		t.Error("registering does not check the password rules")
	}

	profile := apiUserHandlerGo()
	if !strings.Contains(profile, "password.Check(req.Password, user.Email, user.FirstName, user.LastName)") {
		t.Error("changing a password from the account page does not check the rules")
	}

	reset := apiAuthPasswordResetGo()
	if !strings.Contains(reset, "password.Check(req.Password") {
		t.Error("resetting a password does not check the rules, which is where a weak one most often gets in")
	}
}

// The page is an assembly of what already existed, not a second copy of it.
func TestAccountPageComposesTheExistingCards(t *testing.T) {
	page := adminAccountPageTSX()

	for _, part := range []string{
		`import { TwoFactorCard } from "@/components/profile/two-factor-card"`,
		`import { PasskeysCard } from "@/components/security/passkeys"`,
		`import { ActiveSessions } from "@/components/profile/active-sessions"`,
	} {
		if !strings.Contains(page, part) {
			t.Errorf("the account page does not reuse an existing card: %s", part)
		}
	}

	for _, card := range []string{
		"<ProfileForm />", "<PasswordForm />", "<TwoFactorCard />",
		"<PasskeysCard />", "<SignInLinksCard />", "<ActiveSessions />",
		"<CloseAccountCard />",
	} {
		if !strings.Contains(page, card) {
			t.Errorf("the account page does not render %s", card)
		}
	}
}

// One column, not tabs.
//
// Tabs hid six cards behind four labels, so answering "where am I signed in"
// meant knowing that devices were under Devices rather than under Security.
// The page this replaced showed all of it at once, and six cards is a scroll
// rather than a navigation problem.
func TestAccountPageIsOneColumn(t *testing.T) {
	page := adminAccountPageTSX()

	for _, gone := range []string{"useSearchParams", `role="tablist"`, "?tab=", "TABS"} {
		if strings.Contains(page, gone) {
			t.Errorf("the account page still has tabs: found %q", gone)
		}
	}
	// Every link that used to point at a tab now points at a section, so the
	// anchors have to exist.
	for _, id := range []string{`id="security"`, `id="devices"`} {
		if !strings.Contains(page, id) {
			t.Errorf("no anchor for an existing deep link: %s", id)
		}
	}
	if !strings.Contains(page, "scroll-mt-24") {
		t.Error("an anchored section has no scroll margin, so the sticky header covers its heading")
	}
	// Destructive last. It used to sit in the middle of the page because it
	// lived inside the component that draws the first card.
	if strings.Index(page, "<CloseAccountCard />") < strings.Index(page, "<ActiveSessions />") {
		t.Error("closing the account is drawn before the other cards")
	}
}

// The heading comes from PageHeader, not a hand-written <h1>.
//
// PageHeader derives the back link for every /system/* route and carries the
// refresh, theme and notification controls. Writing the heading by hand is how
// this page ended up the only screen in the admin with no way back to the hub.
func TestAccountPageUsesTheSharedHeader(t *testing.T) {
	page := adminAccountPageTSX()

	if !strings.Contains(page, `import { PageHeader } from "@/components/chrome/PageHeader"`) {
		t.Fatal("the account page does not use PageHeader, so it has no back link and no chrome")
	}
	if !strings.Contains(page, `title="Account"`) {
		t.Error("PageHeader is imported but not given the title")
	}
	// The doc comment above the component mentions <h1> to explain the rule, so
	// the check is on markup rather than on the file.
	for _, line := range strings.Split(page, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if strings.Contains(line, "<h1") {
			t.Error("the page writes its own <h1>, which is what lost it the back link")
		}
	}
	// And PageHeader still derives the link for this route.
	if !strings.Contains(adminPageHeaderComponent(), `pathname.startsWith("/system/")`) {
		t.Error("PageHeader no longer derives the back link for /system/* pages")
	}
}

// Every field on the page names itself, which is the defect this admin has
// shipped more than once.
func TestAccountFormsAreLabelled(t *testing.T) {
	for name, form := range map[string]string{
		"password": adminAccountPasswordFormTSX(),
		"profile":  adminAccountProfileFormTSX(),
	} {
		labels := regexp.MustCompile(`htmlFor="([a-z-]+)"`).FindAllStringSubmatch(form, -1)
		ids := regexp.MustCompile(`\bid="([a-z-]+)"`).FindAllStringSubmatch(form, -1)
		if len(labels) == 0 {
			t.Errorf("the %s form has no associated labels", name)
			continue
		}
		have := map[string]bool{}
		for _, id := range ids {
			have[id[1]] = true
		}
		for _, label := range labels {
			if !have[label[1]] {
				t.Errorf("the %s form labels %q, which no field has as its id", name, label[1])
			}
		}
	}
}

// The hub is where somebody goes looking for it.
func TestAccountIsOnTheSystemHub(t *testing.T) {
	hub := adminSystemHubPageV2()
	if !strings.Contains(hub, `href: "/system/account"`) {
		t.Error("the account page is not on the System hub, so nothing links to it")
	}
	if !strings.Contains(hub, `{ href: "/system/account",        category: "Security & Access"`) {
		t.Error("the account tile is not under Security & Access")
	}
}
