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

	for _, tab := range []string{"profile", "password", "security", "devices"} {
		if !strings.Contains(page, `id: "`+tab+`"`) {
			t.Errorf("the %s tab is missing", tab)
		}
	}
}

// Tabs as links, because the pattern that fails a keyboard user is always the
// hand-rolled one.
func TestAccountTabsAreLinks(t *testing.T) {
	page := adminAccountPageTSX()
	if strings.Contains(page, `role="tablist"`) || strings.Contains(page, `role="tab"`) {
		t.Error("the tabs are an ARIA tablist, which needs roving focus to be correct; links need nothing")
	}
	if !strings.Contains(page, `href={"/system/account?tab=" + tab.id}`) {
		t.Error("the tabs do not carry the section in the URL, so a refresh loses the tab and nobody can link to one")
	}
	if !strings.Contains(page, `aria-current={on ? "page" : undefined}`) {
		t.Error("the active tab is not announced")
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
