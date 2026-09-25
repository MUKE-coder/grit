package scaffold

import (
	"strings"
	"testing"
)

// The two account jobs nobody does unprompted, on the screen everybody lands on.
//
// There are four dashboard styles and it is easy to wire three of them. The
// person who picked the fourth style gets an account with no second factor and
// nothing ever telling them, which is exactly the failure this feature exists
// to prevent.
func TestEveryDashboardStyleNudgesTheAccount(t *testing.T) {
	for style, page := range map[string]string{
		"captivating": adminCaptivatingDashboard(),
		"modern":      modernDashboardPage(),
		"minimal":     minimalDashboardPage(),
		"glass":       glassDashboardPage(),
	} {
		if !strings.Contains(page, "<DashboardSecurityNudges />") {
			t.Errorf("the %s dashboard never renders the nudges", style)
		}
		if !strings.Contains(page, `from "@/components/dashboard/security-nudges"`) {
			t.Errorf("the %s dashboard renders the nudges but does not import them", style)
		}
	}
}

func TestTheNudgeSaysNothingWhenThereIsNothingToSay(t *testing.T) {
	card := adminDashboardNudgesTSX()

	// A dashboard that carries a permanent banner is a dashboard whose banners
	// nobody reads, including the next one that matters.
	if !strings.Contains(card, "if (!needsTwoFactor || dismissed.includes(\"two-factor\")) return null;") {
		t.Error("the nudge renders something even when two-factor is already on")
	}
	if !strings.Contains(card, "onDismiss") {
		t.Error("the nudge cannot be dismissed")
	}
	// Waiting for the server before deciding: flashing "turn on two-factor" at
	// somebody who already has it is worse than showing nothing.
	if !strings.Contains(card, "totp !== undefined") {
		t.Error("the nudge decides before the server has answered, so it flashes at accounts that are already fine")
	}
	// localStorage throws outright in some browsers, and a dashboard that
	// crashes over a dismissed banner is a bad trade.
	if strings.Count(card, "} catch {") < 2 {
		t.Error("the localStorage reads and writes are not both guarded")
	}
	if !strings.Contains(card, `"/system/account#security"`) {
		t.Error("the nudge does not lead anywhere useful")
	}

	// Both are load-bearing and both were deleted by accident once, by an edit
	// that removed the card next to them. A string test cannot type-check the
	// emitted TSX, but it can notice a name being used and never defined.
	for _, name := range []string{"const dismiss = ", "useEffect(() => setDismissed(readDismissed()), [])"} {
		if !strings.Contains(card, name) {
			t.Errorf("the component uses %q but never defines it", name)
		}
	}
}

// Verifying an email address is EmailVerifiedBanner's job, on every page.
//
// v3.320.0 shipped a second copy of it on the dashboard, directly under the
// first. Two prompts for one job teach people to ignore both, and the one that
// survives should be the one with the wider reach.
func TestOnlyOneThingAsksYouToConfirmYourEmail(t *testing.T) {
	if !strings.Contains(adminEmailVerifiedBanner(), "/api/auth/verify-email/send") {
		t.Fatal("the banner that owns email verification no longer sends the email")
	}
	if !strings.Contains(adminLayoutComponent(), "<EmailVerifiedBanner />") {
		t.Error("the banner is not in the layout, so it no longer covers every page")
	}
	if strings.Contains(adminDashboardNudgesTSX(), "verify-email") {
		t.Error("the dashboard nudge asks for email verification again, under the banner that already does")
	}

	// House style, and this is the copy a reader sees rather than a comment.
	for _, line := range strings.Split(adminEmailVerifiedBanner(), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}
		if strings.ContainsRune(line, '—') {
			t.Errorf("an em dash is in the banner's copy: %s", trimmed)
		}
	}
}
