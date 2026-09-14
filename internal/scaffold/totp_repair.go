package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairTwoFactor brings a project scaffolded before v3.245.0 up to the fix for
// H6 in the contact-app review.
//
// A pending 2FA token took unlimited guesses for five minutes, and a correct
// password cleared the account's failure count, so signing in again handed out
// fresh guesses. A code could be used again for its whole window. Enabling 2FA
// replaced a secret that was already enabled, and the second step skipped the
// disabled and locked checks the password step makes.
//
// internal/totp, the 2FA models and handlers/totp.go are framework code and
// arrive whole behind the manifest guard. Login is the developer's, so the one
// change there, when the failure count is cleared, anchors on Grit's text.
func repairTwoFactor(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	handler := filepath.Join(apiRoot, "internal", "handlers", "totp.go")
	if !fileExists(handler) {
		return nil // 2FA was removed from this project
	}
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "totp", "totp.go"):         totpServiceGo(),
		filepath.Join(apiRoot, "internal", "models", "two_factor.go"): twoFactorModelsGo(),
		handler: totpHandlerGo(),
	}
	for path, content := range files {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	if !fileContains(handler, "func (h *TOTPHandler) failSecondFactor(") {
		fmt.Println("  ⚠ handlers/totp.go has been edited, so 2FA codes can still be guessed without limit and replayed.\n" +
			"    See what the new version changes with grit upgrade --diff.")
		return nil
	}

	auth := filepath.Join(apiRoot, "internal", "handlers", "auth.go")
	if !fileExists(auth) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, auth, repairLoginFailureCountSource)
}

var earlyCounterResetRe = regexp.MustCompile(`(?s)\n\t// Any successful password check clears the counter[^\n]*\n(?:\t//[^\n]*\n)*\tif user\.FailedLoginCount > 0 \|\| user\.LockedUntil != nil \{\n.*?\n\t\}\n`)

const lateCounterReset = `
	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}
`

// repairLoginFailureCountSource moves the failure-count reset in Login from the
// password check to after the 2FA branch.
func repairLoginFailureCountSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *AuthHandler) Login(") || strings.Contains(src, "The failure count is cleared when the sign-in completes") {
		return src, nil, nil
	}
	loc := earlyCounterResetRe.FindStringIndex(src)
	if loc == nil {
		return src, nil, nil
	}
	const branch = `"message": "Two-factor authentication required",`
	b := strings.Index(src, branch)
	if b < loc[1] {
		return src, nil, []string{"Login clears the failure count before 2FA and is not the handler Grit wrote: clear it only after the second factor, or signing in again resets the guesses"}
	}
	// The end of the 2FA branch: the first line holding only a tab and a brace.
	end := strings.Index(src[b:], "\n\t}\n")
	if end < 0 {
		return src, nil, []string{"could not find the end of Login's 2FA branch, so the failure count is still cleared before it"}
	}
	at := b + end + len("\n\t}\n")
	out := src[:at] + lateCounterReset + src[at:]
	out = out[:loc[0]] + out[loc[1]:]
	if !strings.Contains(out, "\t\"log\"\n") {
		var ok bool
		if out, ok = addImportGroup(out, "log"); !ok {
			return src, nil, []string{"could not add the log import to handlers/auth.go"}
		}
	}
	return out, []string{"the failure count is cleared only when sign-in completes, so signing in again does not reset 2FA guesses"}, nil
}
