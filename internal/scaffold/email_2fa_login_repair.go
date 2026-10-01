package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The sign-in half of the email second factor, for a project that already
// exists.
//
// A project can turn codes by email on from the admin and then find that
// signing in still asks for an authenticator, because the challenge that runs
// at sign-in lives in auth.go, which an existing project already has its own
// copy of. Turning on a factor the sign-in does not know about is worse than
// not having it: the account is one step from locked out.
const (
	emailChallengeOld = `	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Select("id").
		Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err != nil {
		return false
	}`

	emailChallengeTokenOld = `	// The token is stored hashed, so somebody who can read the table cannot
	// finish another account's half-completed sign-in with what they find.
	if err := h.DB.WithContext(c.Request.Context()).Create(&models.TOTPPendingToken{
		UserID:    user.ID,
		TokenHash: totp.HashToken(pendingToken),
		ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
	}).Error; err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to create verification session")
		return true
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"totp_required": true,
			"pending_token": pendingToken,
		},
		"message": "Two-factor authentication required",
	})
	return true
}`
)

// emailChallengeTokenNew is the block that records a pending two-factor
// challenge, taken from the auth template rather than transcribed here. It is
// nested inside the login block, so the two cannot disagree about it.
//
// Markers are whole lines, so a block always ends with the newline that closed
// its last line. The text it replaces here ends at a brace, and leaving the
// newline on would put a blank line after it that gofmt has no reason to remove.
var emailChallengeTokenNew = strings.TrimSuffix(repairBlock("api/handlers/auth.go", "email-challenge"), "\n")

// emailChallengeNew is the two-factor lookup, from the same template. It moved
// onto the service in the handlers-run-no-queries work, so the repair has to
// write what the generator writes and not a copy of what it used to.
var emailChallengeNew = strings.TrimSuffix(repairBlock("api/handlers/auth.go", "two-factor-lookup"), "\n")

func repairEmailTwoFactorLoginSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "startTOTPChallenge") {
		return src, nil, nil
	}
	if strings.Contains(src, "models.TwoFactorMethodEmail") {
		return src, nil, nil
	}
	if strings.Count(src, emailChallengeTokenOld) != 1 || strings.Count(src, emailChallengeOld) != 1 {
		return src, nil, []string{"the two-factor challenge in auth.go is not the one Grit wrote, so a second factor set to email would still ask for an authenticator: compare startTOTPChallenge with a fresh project"}
	}

	out := strings.Replace(src, emailChallengeOld, emailChallengeNew, 1)
	out = strings.Replace(out, emailChallengeTokenOld, emailChallengeTokenNew, 1)

	// No import juggling: auth.go already imports mail, jobs and config for
	// the handler struct itself, which is why the mailer is reachable here.
	return out, []string{"signing in can ask for a code by email, not only an authenticator"}, nil
}

// repairEmailTwoFactorLogin applies it to an existing project.
func repairEmailTwoFactorLogin(root string) error {
	path := filepath.Join(root, "apps", "api", "internal", "handlers", "auth.go")
	if !fileExists(path) {
		path = filepath.Join(root, "api", "internal", "handlers", "auth.go")
		if !fileExists(path) {
			return nil
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairEmailTwoFactorLoginSource)
}
