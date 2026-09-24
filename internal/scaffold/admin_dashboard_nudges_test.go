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

func TestTheNudgesSayNothingWhenThereIsNothingToSay(t *testing.T) {
	card := adminDashboardNudgesTSX()

	// A dashboard that carries two permanent banners is a dashboard whose
	// banners nobody reads, including the next one that matters.
	if !strings.Contains(card, "if (show.length === 0) return null;") {
		t.Error("the nudges render something even when both jobs are done")
	}
	if !strings.Contains(card, "onDismiss") {
		t.Error("neither nudge can be dismissed")
	}
	// Waiting for the server before deciding: flashing "turn on two-factor" at
	// somebody who already has it is worse than showing nothing.
	if !strings.Contains(card, "totp !== undefined") || !strings.Contains(card, "overview !== undefined") {
		t.Error("the nudges decide before the server has answered, so they flash at accounts that are already fine")
	}
	// localStorage throws outright in some browsers, and a dashboard that
	// crashes over a dismissed banner is a bad trade.
	if strings.Count(card, "} catch {") < 2 {
		t.Error("the localStorage reads and writes are not both guarded")
	}
	if !strings.Contains(card, `"/system/account?tab=security"`) {
		t.Error("the two-factor nudge does not lead anywhere useful")
	}
	if !strings.Contains(card, `"/api/auth/verify-email/send"`) {
		t.Error("the email nudge cannot actually send the email")
	}
}
