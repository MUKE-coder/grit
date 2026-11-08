package plugin

import (
	"strings"
	"testing"
)

// The invitation flow, modelled on the Filament reference the plugin was asked
// to follow: a link that works before anybody is signed in, a page that offers
// sign-in or account creation for the invited address, an acceptance that is
// idempotent for somebody who is already a member, and an explicit answer for
// the wrong account.

func TestInvitationPreviewIsPublic(t *testing.T) {
	ctx := Context{Module: "example/apps/api", Architecture: "triple", Frontend: "next"}
	var public, protected string
	for _, injection := range multitenantInjections(ctx) {
		switch injection.Marker {
		case "// grit:routes:custom":
			public += injection.Code
		case "// grit:routes:protected":
			protected += injection.Code
		}
	}
	if !strings.Contains(public, `/organization-invitations/preview`) {
		t.Fatal("the preview route is not registered outside the authenticated group: " +
			"the page has to ask what a link is before anybody signs in")
	}
	if strings.Contains(protected, `/organization-invitations/preview`) {
		t.Error("the preview route is behind auth, so an invited person who is not " +
			"signed in cannot even see which organization invited them")
	}
}

func TestInvitationPreviewReturnsTheLinkAndItsState(t *testing.T) {
	body := invitationBody(t, mtInvitationHandler(Context{Module: "example/apps/api"}), "Preview")

	for _, want := range []string{
		`"organization_name"`,
		`"email"`,
		`"role"`,
		`"status"`,
		"invitation.AcceptedAt",
		"invitation.RevokedAt",
		"invitation.ExpiresAt",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("Preview does not answer %s, so the page cannot decide what to "+
				"show before a sign-in", want)
		}
	}
	if !strings.Contains(body, `token_hash = ?`) {
		t.Error("Preview does not look the token up by hash")
	}
}

func TestAcceptingIsIdempotentForAMemberAndExplicitForAnotherEmail(t *testing.T) {
	body := invitationBody(t, mtInvitationHandler(Context{Module: "example/apps/api"}), "Accept")

	if !strings.Contains(body, "alreadyMember") {
		t.Error("Accept has no already-a-member path, so the second click of a link " +
			"that already worked answers invalid")
	}
	if !strings.Contains(body, "invitation.AcceptedAt != nil") {
		t.Error("Accept does not handle a spent invitation for the person it added")
	}
	if !strings.Contains(body, "EMAIL_MISMATCH") {
		t.Error("Accept does not say when the invitation belongs to another address, " +
			"so the reader only gets a generic invalid link")
	}
	if !strings.Contains(body, "errInvitationExpired") {
		t.Error("Accept does not distinguish an expired invitation from a forged one")
	}
	// The membership is created before the invitation is marked accepted, in
	// one transaction: a crash between them either leaves both undone or the
	// member already in place, never a spent link with no membership.
	if !strings.Contains(body, "h.DB.Transaction(func(tx *gorm.DB) error {") {
		t.Error("Accept is not transactional, so a failure after creating the " +
			"membership could leave the invitation pending for a second attempt")
	}
}

func TestInvitationPageOffersEveryWayIn(t *testing.T) {
	src := mtWebInvitationPage()

	for _, want := range []string{
		`"/api/organization-invitations/preview"`,
		`"/api/auth/me"`,
		`"/api/auth/login"`,
		`"/api/auth/register"`,
		`"/api/auth/verify-email/send"`,
		`"/api/organization-invitations/accept"`,
		"setActiveOrg(organizationID)",
		`window.location.assign("/account/organization")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the invitation page never calls %s, so a reader without a "+
				"session or without an account has no way forward", want)
		}
	}
	if strings.Contains(src, "next/navigation") || strings.Contains(src, `"next/link"`) {
		t.Error("the invitation page imports Next-specific navigation; the Vite " +
			"web app renders the same template and cannot resolve it")
	}
	if strings.Contains(src, "apiClient") {
		t.Error("the invitation page calls apiClient; the Vite client exports api " +
			"and the shared module is the one place that difference lives")
	}
	// Sign-in, account creation and verification each need a visible control.
	for _, state := range []string{`mode === "signin"`, `mode === "register"`, `mode === "verify"`} {
		if !strings.Contains(src, state) {
			t.Errorf("the invitation page has no %s state", state)
		}
	}
}
