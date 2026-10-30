package scaffold

import "strings"

// Review 2.9: a provider sign-in trusted an address the provider never verified.
//
// goth puts the provider's userinfo response in RawData, email_verified is in
// there, and nothing read it. An address a provider has not verified is one
// anybody can claim, and the callback treated the address as proof of who was
// signing in: it created accounts with EmailVerifiedAt set, it linked the
// provider to an account that already existed, and for an account whose address
// was never confirmed it removed the password and revoked every session.
//
// Every constant below is the text the template holds, generated from it, and
// TestOAuthVerifiedRepairMatchesTheTemplate compares a repaired handler to a
// freshly written one byte for byte.

const (
	oauthVerifiedPredicate = `// providerVerifiedTheEmail is whether the provider says it owns this address.
//
// An address a provider has not verified is an address anybody can claim, and
// claiming one is enough to be matched to somebody else's account by email.
// goth puts the provider's own userinfo response in RawData, which is where
// the answer is, and review 2.9 found nothing reading it.
//
// Per provider, because they do not agree on how to say it:
//
//   - Google returns email_verified, and older responses verified_email. Both
//     arrive as a bool through encoding/json, but a provider that hands back
//     "true" as a string is one string compare away, so both are handled.
//   - GitHub does not say, because it does not need to: the public profile
//     email must be a verified address, and when that is empty goth falls back
//     to /user/emails and takes the one that is primary AND verified. There is
//     no path through goth's github provider that yields an unverified
//     address.
//   - Anything else: not asserted. A provider added here without a rule is
//     treated as unverified rather than trusted by default, so adding one
//     cannot quietly open the hole this closed.
func providerVerifiedTheEmail(provider string, raw map[string]interface{}) bool {
	if provider == "github" {
		return true
	}
	for _, key := range []string{"email_verified", "verified_email"} {
		switch v := raw[key].(type) {
		case bool:
			return v
		case string:
			return v == "true"
		}
	}
	return false
}

`

	oauthVerifiedPredicateAnchor = "// OAuthCallback completes the OAuth flow"

	oauthVerifiedCompute = `	// Whether the provider will stand behind the address, which decides what
	// this callback is allowed to do with it. See providerVerifiedTheEmail.
	emailVerified := providerVerifiedTheEmail(provider, gothUser.RawData)

`

	oauthVerifiedComputeAnchor = "\t// Find or create user by email\n"

	oauthVerifiedNewAccount = `				Active:          true,
				IPAddress:       c.ClientIP(),
			}

			// Only if the provider verified it. Setting this from an
			// unverified claim turns "somebody typed this address into a
			// provider" into "this address is confirmed", and the address is
			// what every later match is made on. Left nil, the account goes
			// through the same verification email a local sign-up does.
			if emailVerified {
				user.EmailVerifiedAt = &now
			}

`

	oauthVerifiedNewAccountOld = "\t\t\t\tActive:          true,\n\t\t\t\tEmailVerifiedAt: &now,\n\t\t\t\tIPAddress:       c.ClientIP(),\n\t\t\t}\n\n"

	oauthVerifiedRefusal = `		// This provider account has signed in to this user before, so the link
		// was established deliberately and is not being claimed now. Checked
		// first: a returning user whose provider does not assert verification
		// must not be locked out by the refusal below.
		alreadyLinked := (provider == "google" && user.GoogleID != "" && user.GoogleID == gothUser.UserID) ||
			(provider == "github" && user.GithubID != "" && user.GithubID == gothUser.UserID)

		// An unverified claim does not reach an account that already exists.
		// Anybody can put somebody else's address into a provider that does
		// not check it, and everything below this line treats the address as
		// proof of who is signing in: it links the provider to the account, and
		// for an unconfirmed account it removes the password outright.
		//
		// Refused with the way forward rather than silently creating a second
		// account, because a second account with the same address is its own
		// kind of mess.
		if !emailVerified && !alreadyLinked {
			log.Printf("oauth: %s did not verify %s, and an account with that address exists", provider, gothUser.Email)
			redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your provider has not confirmed that email address. Sign in with your password, then connect "+provider+" from your profile."))
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		}

`

	oauthVerifiedRefusalAnchor = "\t\t// An account whose address was never confirmed may have been registered by\n"

	oauthVerifiedNotice = `
			// And say so. In the case this exists for, the person who had the
			// password was not the owner and should hear nothing useful; in the
			// benign case somebody registered, never confirmed, and later used
			// the provider, and their password silently stopped working. The
			// second case is far more common, and the mail is how they find out
			// why their password no longer works.
			//
			// SendRaw rather than a template: this is one paragraph, and a
			// project that has not configured a mailer gets a log line instead
			// of a failed sign-in.
			if h.Mailer != nil {
				body := fmt.Sprintf(
					"<p>You signed in to %s with %s.</p>"+
						"<p>This account's email address had never been confirmed, so the "+
						"password on it has been removed and any existing sessions were "+
						"signed out. Sign in with %s from now on, or set a new password "+
						"from your profile.</p>"+
						"<p>If this was not you, reset your password immediately.</p>",
					h.Config.AppName, provider, provider)
				if err := h.Mailer.SendRaw(c.Request.Context(), user.Email,
					"Your password was removed", body); err != nil {
					log.Printf("oauth: telling %s their password was removed: %v", user.ID, err)
				}
			} else {
				log.Printf("oauth: password removed for %s and no mailer is configured to say so", user.ID)
			}
`

	oauthVerifiedNoticeAnchor = "\t\t\tuser.Password = \"\"\n\t\t\tuser.EmailVerifiedAt = &verifiedAt\n"
)

// repairOAuthVerifiedEmailSource makes the callback ask whether the provider
// verified the address before it treats the address as proof of anything.
func repairOAuthVerifiedEmailSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func providerVerifiedTheEmail(") {
		return src, nil, nil
	}
	// The repair that clears an unconfirmed password has to be in place first:
	// this one decides when that block may run, and adding the condition to a
	// handler that does not have the block would be a condition on nothing.
	if !strings.Contains(src, "user.EmailVerifiedAt == nil") {
		return src, nil, nil
	}

	const warn = "auth_oauth.go is not the file Grit wrote: make the callback read email_verified from the provider's RawData and refuse to match an existing account on an address the provider has not verified, or anybody who can put that address into a provider can take the account"

	edits := [][2]string{
		{oauthVerifiedPredicateAnchor, oauthVerifiedPredicate + oauthVerifiedPredicateAnchor},
		{oauthVerifiedComputeAnchor, oauthVerifiedCompute + oauthVerifiedComputeAnchor},
		{oauthVerifiedNewAccountOld, oauthVerifiedNewAccount},
		{oauthVerifiedRefusalAnchor, oauthVerifiedRefusal + oauthVerifiedRefusalAnchor},
		{oauthVerifiedNoticeAnchor, oauthVerifiedNoticeAnchor + oauthVerifiedNotice},
	}
	// All of them or none. A handler with the predicate and not the refusal
	// reads email_verified and then ignores it, which is worse than before
	// because it looks done.
	out := src
	for _, e := range edits {
		if strings.Count(out, e[0]) != 1 {
			return src, nil, []string{warn}
		}
		out = strings.Replace(out, e[0], e[1], 1)
	}
	return out, []string{"a provider sign-in is only matched to an existing account on an address the provider has verified, and a password removed by one is now notified"}, nil
}
