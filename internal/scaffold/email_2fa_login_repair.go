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

	emailChallengeNew = `	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Select("id", "method").
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

	emailChallengeTokenNew = `	pending := models.TOTPPendingToken{
		UserID: user.ID,
		// The token is stored hashed, so somebody who can read the table cannot
		// finish another account's half-completed sign-in with what they find.
		TokenHash: totp.HashToken(pendingToken),
		ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
	}

	// A code by email instead of an authenticator.
	code := ""
	if totpConfig.Method == models.TwoFactorMethodEmail && h.Mailer == nil && h.Jobs == nil {
		// The account is set to a factor this deployment can no longer deliver,
		// which happens when mail is switched off after somebody turned it on.
		// Say so: the alternative is a code box waiting for a code that will
		// never arrive.
		log.Printf("two-factor: %s uses codes by email and no mailer is configured", user.ID)
		respond.Fail(c, respond.CodeMailFailed, "Your sign-in code cannot be sent: this deployment has no mail configured. Ask an administrator.")
		return true
	}
	if totpConfig.Method == models.TwoFactorMethodEmail {
		var err error
		code, err = totp.GenerateEmailCode()
		if err != nil {
			respond.Fail(c, respond.CodeTokenError, "Failed to create verification session")
			return true
		}
		pending.CodeHash = totp.HashToken(code)
	}

	if err := h.DB.WithContext(c.Request.Context()).Create(&pending).Error; err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to create verification session")
		return true
	}

	if code != "" {
		// Sent after the row exists, so a code can never arrive for a challenge
		// that was not recorded.
		if err := dispatchMail(c.Request.Context(), h.Mailer, h.Jobs, "two-factor:"+user.ID+":"+totp.HashToken(code), mail.SendOptions{
			To:       user.Email,
			Subject:  "Your sign-in code",
			Template: "two-factor-code",
			Data: map[string]interface{}{
				"AppName": h.Config.AppName,
				"Title":   "Your sign-in code",
				"Code":    code,
				"Minutes": int(totp.PendingTokenExpiry.Minutes()),
				"Year":    time.Now().Year(),
			},
		}); err != nil {
			log.Printf("two-factor: emailing a code to %s: %v", user.ID, err)
			respond.Fail(c, respond.CodeMailFailed, "We could not email your code. Try again in a moment.")
			return true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"totp_required": true,
			"pending_token": pendingToken,
			// Which of the two the client should ask for. The older shape had
			// no method and meant an authenticator, so that is what an empty
			// value means here.
			"method": totpConfig.Method,
		},
		"message": "Two-factor authentication required",
	})
	return true
}`
)

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
