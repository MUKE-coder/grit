package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L1: sign-in reported the state of an account before its
// password was checked. ACCOUNT_LOCKED, ACCOUNT_DISABLED, EMAIL_NOT_VERIFIED and
// SOCIAL_AUTH_ONLY (which named the provider) answered any address typed into
// the form, and an unknown address answered without any bcrypt work, so both
// the body and the time told an attacker which addresses held accounts.
//
// loginBlockOld283 is the block auth.go carried up to v3.284.0, from
// loginRefusalReason to the end of Login, with hand-built error envelopes.
// loginBlockOld286 is the same block from v3.285.0, answering through
// respond.Fail. loginBlockNew replaces either.

const loginBlockOld283 = `// loginRefusalReason is a decision not to let an account sign in, taken before
// its password is read.
//
// Each of these used to be an if-block in the middle of Login, with its own
// c.JSON and its own activity row, which is most of why that handler was 182
// lines. Answering them here means the order is visible in one place: disabled,
// then locked, then unverified, then social-only.
type loginRefusalReason struct {
	Status  int
	Code    string
	Message string
	// Action is the activity action to record, or "" when the refusal is not
	// worth an audit row. Only the two that suggest an attack are.
	Action  string
	Summary string
}

// loginRefusal reports why user may not sign in, or nil when it may.
//
// Nothing here compares the password. A disabled or locked account answers the
// same way whether or not the password was right, so neither can be probed by
// timing the comparison.
func (h *AuthHandler) loginRefusal(user *models.User) *loginRefusalReason {
	if !user.Active {
		return &loginRefusalReason{
			Status:  http.StatusForbidden,
			Code:    "ACCOUNT_DISABLED",
			Message: "Your account has been disabled",
			Action:  "auth.login_blocked",
			Summary: "Sign-in blocked for disabled account " + user.Email,
		}
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		return &loginRefusalReason{
			Status:  http.StatusTooManyRequests,
			Code:    "ACCOUNT_LOCKED",
			Message: fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			Action:  "auth.login_locked",
			Summary: "Sign-in refused: account is temporarily locked",
		}
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected — the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		return &loginRefusalReason{
			Status:  http.StatusForbidden,
			Code:    "EMAIL_NOT_VERIFIED",
			Message: "Confirm your email address before signing in. Check your inbox for the link.",
		}
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		return &loginRefusalReason{
			Status:  http.StatusBadRequest,
			Code:    "SOCIAL_AUTH_ONLY",
			Message: fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
		}
	}

	return nil
}

// startTOTPChallenge issues the second-factor challenge, and reports whether it
// answered the request.
//
// true means Login is finished: the response carries a short-lived pending
// token the client exchanges at /auth/totp/verify. false means there is no
// second factor to ask for, either because the account has none or because this
// device is already trusted, and Login should carry on and issue tokens.
func (h *AuthHandler) startTOTPChallenge(c *gin.Context, user *models.User) bool {
	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err != nil {
		return false
	}
	if IsTrustedDevice(c, h.DB, user.ID) {
		return false
	}

	pendingToken, err := totp.GeneratePendingToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
		})
		return true
	}

	// The token is stored hashed, so somebody who can read the table cannot
	// finish another account's half-completed sign-in with what they find.
	if err := h.DB.WithContext(c.Request.Context()).Create(&models.TOTPPendingToken{
		UserID:    user.ID,
		TokenHash: totp.HashToken(pendingToken),
		ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
		})
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
}

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Everything that refuses this account before its password is read. Kept out
	// of Login so the handler reads as the sequence it is, and so the rules can
	// be answered as one question rather than four scattered ifs.
	if refusal := h.loginRefusal(&user); refusal != nil {
		if refusal.Action != "" {
			services.LogActivity(h.DB, c, services.ActivityArgs{
				Action:       refusal.Action,
				Severity:     "warn",
				Summary:      refusal.Summary,
				ResourceType: "user",
				ResourceID:   user.ID,
			})
		}
		c.JSON(refusal.Status, gin.H{
			"error": gin.H{"code": refusal.Code, "message": refusal.Message},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// Wrong password on a real account — distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// The second factor, when this account has one and this device is not
	// already trusted. It answers the request itself when it does.
	if h.startTOTPChallenge(c, &user) {
		return
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below — both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

`

// loginBlockOld286 is the same block as 3.285.0 and 3.286.0 wrote it, answering
// through respond.Fail.
const loginBlockOld286 = `// loginRefusalReason is a decision not to let an account sign in, taken before
// its password is read.
//
// Each of these used to be an if-block in the middle of Login, with its own
// c.JSON and its own activity row, which is most of why that handler was 182
// lines. Answering them here means the order is visible in one place: disabled,
// then locked, then unverified, then social-only.
type loginRefusalReason struct {
	Status  int
	Code    string
	Message string
	// Action is the activity action to record, or "" when the refusal is not
	// worth an audit row. Only the two that suggest an attack are.
	Action  string
	Summary string
}

// loginRefusal reports why user may not sign in, or nil when it may.
//
// Nothing here compares the password. A disabled or locked account answers the
// same way whether or not the password was right, so neither can be probed by
// timing the comparison.
func (h *AuthHandler) loginRefusal(user *models.User) *loginRefusalReason {
	if !user.Active {
		return &loginRefusalReason{
			Status:  http.StatusForbidden,
			Code:    "ACCOUNT_DISABLED",
			Message: "Your account has been disabled",
			Action:  "auth.login_blocked",
			Summary: "Sign-in blocked for disabled account " + user.Email,
		}
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		return &loginRefusalReason{
			Status:  http.StatusTooManyRequests,
			Code:    "ACCOUNT_LOCKED",
			Message: fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			Action:  "auth.login_locked",
			Summary: "Sign-in refused: account is temporarily locked",
		}
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected — the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		return &loginRefusalReason{
			Status:  http.StatusForbidden,
			Code:    "EMAIL_NOT_VERIFIED",
			Message: "Confirm your email address before signing in. Check your inbox for the link.",
		}
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		return &loginRefusalReason{
			Status:  http.StatusBadRequest,
			Code:    "SOCIAL_AUTH_ONLY",
			Message: fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
		}
	}

	return nil
}

// startTOTPChallenge issues the second-factor challenge, and reports whether it
// answered the request.
//
// true means Login is finished: the response carries a short-lived pending
// token the client exchanges at /auth/totp/verify. false means there is no
// second factor to ask for, either because the account has none or because this
// device is already trusted, and Login should carry on and issue tokens.
func (h *AuthHandler) startTOTPChallenge(c *gin.Context, user *models.User) bool {
	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err != nil {
		return false
	}
	if IsTrustedDevice(c, h.DB, user.ID) {
		return false
	}

	pendingToken, err := totp.GeneratePendingToken()
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to create verification session")
		return true
	}

	// The token is stored hashed, so somebody who can read the table cannot
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
}

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		respond.Fail(c, respond.CodeInvalidCredentials, "Invalid email or password")
		return
	}

	// Everything that refuses this account before its password is read. Kept out
	// of Login so the handler reads as the sequence it is, and so the rules can
	// be answered as one question rather than four scattered ifs.
	if refusal := h.loginRefusal(&user); refusal != nil {
		if refusal.Action != "" {
			services.LogActivity(h.DB, c, services.ActivityArgs{
				Action:       refusal.Action,
				Severity:     "warn",
				Summary:      refusal.Summary,
				ResourceType: "user",
				ResourceID:   user.ID,
			})
		}
		c.JSON(refusal.Status, gin.H{
			"error": gin.H{"code": refusal.Code, "message": refusal.Message},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// Wrong password on a real account — distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		respond.Fail(c, respond.CodeInvalidCredentials, "Invalid email or password")
		return
	}

	// The second factor, when this account has one and this device is not
	// already trusted. It answers the request itself when it does.
	if h.startTOTPChallenge(c, &user) {
		return
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below — both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

`

const loginBlockNew = `// invalidCredentials is the one answer a sign-in gets until its password has
// been checked and found right. An unknown address, a wrong password, a locked
// account and an account with no password all get it, with the same status,
// code and message, so a sign-in form cannot be used to learn which addresses
// hold accounts or what state those accounts are in.
func invalidCredentials(c *gin.Context) {
	respond.Fail(c, respond.CodeInvalidCredentials, "Invalid email or password")
}

var (
	timingHashOnce sync.Once
	timingHash     []byte
)

// spendPasswordCheck does the bcrypt work of a password comparison when there is
// no hash to compare against. Without it an unknown address answered in a
// fraction of the time a real one took, which gave away the same thing the
// response body no longer does.
func spendPasswordCheck(password string) {
	timingHashOnce.Do(func() {
		hash, err := bcrypt.GenerateFromPassword([]byte("grit: no account to compare"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("auth: preparing the stand-in password hash: %v", err)
			return
		}
		timingHash = hash
	})
	if bcrypt.CompareHashAndPassword(timingHash, []byte(password)) == nil {
		return // a match here signs nobody in: only the time spent matters
	}
}

// loginRefusalReason is a decision not to let an account sign in, taken once
// its password is known to be right.
type loginRefusalReason struct {
	Code    respond.Code
	Message string
	// Action is the activity action to record, or "" when the refusal is not
	// worth an audit row.
	Action  string
	Summary string
}

// loginRefusal reports why user may not sign in, or nil when it may.
//
// It runs only after the password matched. These answers name the state of the
// account, disabled or unverified, and before they waited for the password
// anybody could read them for any address they typed.
func (h *AuthHandler) loginRefusal(user *models.User) *loginRefusalReason {
	if !user.Active {
		return &loginRefusalReason{
			Code:    respond.CodeAccountDisabled,
			Message: "Your account has been disabled",
			Action:  "auth.login_blocked",
			Summary: "Sign-in blocked for disabled account " + user.Email,
		}
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected: the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil {
		return &loginRefusalReason{
			Code:    respond.CodeEmailNotVerified,
			Message: "Confirm your email address before signing in. Check your inbox for the link.",
		}
	}

	return nil
}

// checkCredentials reports whether password signs in as user, and answers the
// request with invalidCredentials when it does not.
//
// A locked account is not compared at all: a lock that still checked the
// password would tell whoever guessed right that they had, and let them keep
// guessing through it. An account with no password (it signs in with a social
// provider) has nothing to match. Both still spend a comparison's worth of work,
// so neither is quicker to answer than a wrong password.
func (h *AuthHandler) checkCredentials(c *gin.Context, user *models.User, password string) bool {
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		spendPasswordCheck(password)
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_locked",
			Severity:     "warn",
			Summary:      "Sign-in refused: account is temporarily locked",
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		invalidCredentials(c)
		return false
	}

	if user.Password == "" {
		spendPasswordCheck(password)
		services.LogLoginFailed(h.DB, c, user.Email)
		invalidCredentials(c)
		return false
	}

	if !user.CheckPassword(password) {
		// Wrong password on a real account: distinct in the activity log from
		// an unknown email, because Sentinel's brute-force heuristics weight
		// these higher. Not distinct in the response.
		services.LogLoginFailed(h.DB, c, user.Email)
		h.registerFailedLogin(user)
		invalidCredentials(c)
		return false
	}
	return true
}

// startTOTPChallenge issues the second-factor challenge, and reports whether it
// answered the request.
//
// true means Login is finished: the response carries a short-lived pending
// token the client exchanges at /auth/totp/verify, or an error. false means
// there is no second factor to ask for, either because the account has none or
// because this device is already trusted, and Login should carry on and issue
// tokens.
func (h *AuthHandler) startTOTPChallenge(c *gin.Context, user *models.User) bool {
	// Only the id is read. The secret is encrypted when FIELD_ENCRYPTION_KEY is
	// set, and loading it here would let a secret that cannot be decrypted read
	// as "this account has no second factor".
	// The id and the method. The secret is encrypted when FIELD_ENCRYPTION_KEY
	// is set, and loading it here would let a secret that cannot be decrypted
	// read as "this account has no second factor".
	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Select("id", "method").
		Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err != nil {
		return false
	}
	if IsTrustedDevice(c, h.DB, user.ID) {
		return false
	}

	pendingToken, err := totp.GeneratePendingToken()
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to create verification session")
		return true
	}

` + emailChallengeTokenNew + `

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		spendPasswordCheck(req.Password)
		services.LogLoginFailed(h.DB, c, req.Email)
		invalidCredentials(c)
		return
	}

	// Nothing about the account is reported until its password is right.
	if !h.checkCredentials(c, &user, req.Password) {
		return
	}

	// The password matched, so the caller may be told why the account still
	// cannot sign in.
	if refusal := h.loginRefusal(&user); refusal != nil {
		if refusal.Action != "" {
			services.LogActivity(h.DB, c, services.ActivityArgs{
				Action:       refusal.Action,
				Severity:     "warn",
				Summary:      refusal.Summary,
				ResourceType: "user",
				ResourceID:   user.ID,
			})
		}
		respond.Fail(c, refusal.Code, refusal.Message)
		return
	}

	// The second factor, when this account has one and this device is not
	// already trusted. It answers the request itself when it does.
	if h.startTOTPChallenge(c, &user) {
		return
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below. Both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}

`

// ─── L2: two-factor secrets are encrypted at rest ────────────────────────────

// twoFactorSecretField is TwoFactorConfig.Secret in models/two_factor.go.
const twoFactorSecretField = `	// Secret is encrypted with FIELD_ENCRYPTION_KEY when one is set. It was
	// stored in the clear, so anybody who could read the table, a backup or a
	// GORM Studio session could generate every user's codes. A secret written
	// before the key reads back as it is and is encrypted by grit migrate, or
	// the next time its code verifies. Without a key it is stored as before.
	Secret crypto.EncryptedString ` + "`" + `gorm:"type:text;not null" json:"-"` + "`" + `
`

// totpSealSecretFunc is TOTPHandler.sealSecret in handlers/totp.go.
const totpSealSecretFunc = `// sealSecret stores a secret that is still in the clear encrypted, once it has
// verified a code.
//
// EncryptedString reads a value without its enc:v1: prefix as it is, so a
// secret enabled before FIELD_ENCRYPTION_KEY was set keeps working. grit migrate
// encrypts all of them at once; this covers a project that has not run it. The
// update matches only a row that is still plaintext, and a failure is logged
// and changes nothing: the secret keeps verifying either way.
func (h *TOTPHandler) sealSecret(c *gin.Context, config *models.TwoFactorConfig) {
	if !crypto.EncryptionEnabled() {
		return
	}
	if err := h.DB.WithContext(c.Request.Context()).Model(&models.TwoFactorConfig{}).
		Where("id = ? AND secret NOT LIKE ?", config.ID, "enc:v1:%").
		UpdateColumn("secret", config.Secret).Error; err != nil {
		log.Printf("totp: encrypting the stored secret for %s: %v", config.UserID, err)
	}
}
`

// migrateEncryptExistingHook follows the sync backfill in cmd/migrate.
const migrateEncryptExistingHook = `	// Values in encrypted columns written before FIELD_ENCRYPTION_KEY was set,
	// two-factor secrets among them, are encrypted now. Without a key this
	// changes nothing.
	if n, err := crypto.EncryptExisting(db, models.Models()...); err != nil {
		log.Printf("Some values in encrypted columns are still stored in the clear: %v", err)
	} else if n > 0 {
		fmt.Printf("Encrypted %d value(s) stored before field encryption was on.\n", n)
	}
`

// ─── L3: the trusted-device cookie is Secure and SameSite ────────────────────

// trustedDeviceCookieNew sets totp_trusted in TOTPHandler.createTrustedDevice.
const trustedDeviceCookieNew = `	// Secure whenever the request came over HTTPS, directly or through a proxy
	// that says so, as the auth cookies are, and SameSite=Lax like them. The
	// cookie that lets a device skip the second factor was neither, so it went
	// out over plain HTTP and rode along with requests other sites started.
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("totp_trusted", deviceToken, int(totp.TrustedDeviceDuration.Seconds()), "/", "", secure, true)
`

// ─── L4: a refused socket handshake does not echo the token error ────────────

const (
	// realtimeInvalidTokenOld286 is the refusal as 3.285.0 and 3.286.0 wrote it.
	realtimeInvalidTokenOld286 = `		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "INVALID_TOKEN", "message": err.Error()},
		})
		return
	}
`
	realtimeInvalidTokenOld283 = `		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_TOKEN", "message": err.Error()}})
		return
	}
`
	realtimeInvalidTokenNew = `		// The parser's reason (expired, bad signature, malformed) is for the log,
		// not for whoever sent the token.
		log.Printf("[ws] refused a handshake token: %v", err)
		respond.Fail(c, respond.CodeInvalidToken, "Invalid or expired token")
		return
	}
`
)

// ─── L5: an import job is read only by whoever started it ────────────────────

// importJobCreatedByField is ImportJob.CreatedBy in models/import_job.go.
const importJobCreatedByField = `	// CreatedBy is the user who started the import. Only they, or an admin, may
	// read the job: its row errors quote the file's contents.
	CreatedBy string ` + "`" + `gorm:"size:36;index" json:"-"` + "`" + `
`

// importJobMessageLine is the field CreatedBy follows.
const importJobMessageLine = "\tMessage   string    `gorm:\"size:500\" json:\"message\"`\n"

const (
	importJobGetOld = `func (h *ImportJobHandler) GetByID(c *gin.Context) {
	var job models.ImportJob
	if err := h.DB.WithContext(c.Request.Context()).First(&job, "id = ?", c.Param("id")).Error; err != nil {
`
	importJobGetNew = `func (h *ImportJobHandler) GetByID(c *gin.Context) {
	var job models.ImportJob
	// Whoever started the import, or an admin. Anybody else, and a job from
	// before its starter was recorded, is told it does not exist: any signed-in
	// user holding the id read another user's counts and row errors.
	q := h.DB.WithContext(c.Request.Context()).Where("id = ?", c.Param("id"))
	if actor := authz.ActorOf(c); !actor.Admin {
		q = q.Where("created_by = ? AND created_by <> ''", actor.UserID)
	}
	if err := q.First(&job).Error; err != nil {
`
)

// oldStartImportRe is the job grit generate created before it recorded who
// started it.
var oldStartImportRe = regexp.MustCompile(`job := models\.ImportJob\{Resource: "(\w+)", Status: "processing", Total: total\}`)

// repairAuthDisclosure brings a project up to L1 to L5 of the contact-app
// review, in the files upgrade does not deliver whole.
//
// handlers/totp.go, models/two_factor.go (L2, L3), internal/crypto and
// cmd/migrate (L2) arrive whole. auth.go (L1), the realtime handler (L4), the
// import job model and handler, and generated importers (L5) are the
// developer's, and are changed only where they still read as Grit wrote them.
func repairAuthDisclosure(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	handlers := filepath.Join(apiRoot, "internal", "handlers")
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if auth := filepath.Join(handlers, "auth.go"); fileExists(auth) {
		if err := repairSourceFile(root, m, auth, repairLoginDisclosureSource); err != nil {
			return err
		}
	}
	if rt := filepath.Join(handlers, "realtime.go"); fileExists(rt) {
		if err := repairSourceFile(root, m, rt, repairRealtimeTokenErrorSource); err != nil {
			return err
		}
	}
	return repairImportJobOwner(root, opts, m)
}

// repairLoginDisclosureSource checks the password before Login says anything
// about the account.
func repairLoginDisclosureSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *AuthHandler) Login(") || strings.Contains(src, "func (h *AuthHandler) checkCredentials(") {
		return src, nil, nil
	}
	const refused = "Login is not the one Grit wrote, so it still tells anybody whether an address is locked, disabled or social-only before checking the password: answer INVALID_CREDENTIALS until the password matches (see grit upgrade --diff)"
	var out string
	switch {
	case strings.Count(src, loginBlockOld286) == 1:
		out = strings.Replace(src, loginBlockOld286, loginBlockNew, 1)
	case strings.Count(src, loginBlockOld283) == 1:
		out = strings.Replace(src, loginBlockOld283, loginBlockNew, 1)
	default:
		replaced, ok := replaceMonolithicLogin(src)
		if !ok {
			return src, nil, []string{refused}
		}
		out = replaced
	}
	imports := []struct{ path, anchor string }{
		{"sync", "\t\"time\"\n"},
		{"golang.org/x/crypto/bcrypt", "\t\"gorm.io/gorm\"\n"},
		// A project from before 3.285.0 answered by hand and has no respond import.
		{authModule(src) + "/internal/respond", "\t\"" + authModule(src) + "/internal/services\"\n"},
	}
	for _, imp := range imports {
		var ok bool
		if out, ok = addImportBefore(out, imp.path, imp.anchor); !ok {
			return src, nil, []string{"could not add the " + imp.path + " import, so Login was left as it was: answer INVALID_CREDENTIALS until the password matches"}
		}
	}
	out, _ = dropUnusedImport(out, "fmt")
	return out, []string{"sign-in answers INVALID_CREDENTIALS, in the same time, until the password matches, so it no longer reveals which addresses hold accounts or whether they are locked, disabled or social-only"}, nil
}

// authModule is the module path a handler file imports its packages from, read
// from its import of internal/models.
func authModule(src string) string {
	const suffix = "/internal/models\"\n"
	end := strings.Index(src, suffix)
	if end < 0 {
		return ""
	}
	start := strings.LastIndex(src[:end], "\"") + 1
	return src[start:end]
}

// loginMonolithic is Login as projects created before v3.283.0 hold it: one
// function answering every refusal before the password, with no loginRefusal
// or startTOTPChallenge beside it. That split shipped without an upgrade repair,
// so a project created earlier and upgraded since still has this. This copy is
// the function as a v3.282.0 project holds it after upgrading to v3.286.0.
//
// Older copies differ only in details Login's answers do not depend on: whether
// queries carry the request context, whether the pending 2FA token insert is
// checked, and where and how the failure count is cleared. canonicalLogin
// removes exactly those, and nothing else, before comparing.
const loginMonolithic = `func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		services.LogLoginFailed(h.DB, c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if !user.Active {
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_blocked",
			Severity:     "warn",
			Summary:      "Sign-in blocked for disabled account " + user.Email,
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_DISABLED",
				"message": "Your account has been disabled",
			},
		})
		return
	}

	// Locked accounts are refused before the password is even compared, so a
	// lockout cannot be probed by timing the comparison.
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Round(time.Minute)
		if remaining < time.Minute {
			remaining = time.Minute
		}
		services.LogActivity(h.DB, c, services.ActivityArgs{
			Action:       "auth.login_locked",
			Severity:     "warn",
			Summary:      "Sign-in refused: account is temporarily locked",
			ResourceType: "user",
			ResourceID:   user.ID,
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_LOCKED",
				"message": fmt.Sprintf("Too many failed attempts. Try again in about %d minute(s), or reset your password.", int(remaining.Minutes())),
			},
		})
		return
	}

	// Opt-in gate. Social and SSO sign-ins are unaffected: the IdP already
	// proved the address, and those paths set EmailVerifiedAt on first login.
	if h.Config.RequireEmailVerification && user.EmailVerifiedAt == nil && user.Password != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "EMAIL_NOT_VERIFIED",
				"message": "Confirm your email address before signing in. Check your inbox for the link.",
			},
		})
		return
	}

	if user.Password == "" {
		provider := user.Provider
		if provider == "" || provider == "local" {
			provider = "social login"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "SOCIAL_AUTH_ONLY",
				"message": fmt.Sprintf("This account uses %s. Please sign in with your social account.", provider),
			},
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// Wrong password on a real account: distinct from "unknown email"
		// because Sentinel's brute-force heuristics weight these higher.
		services.LogLoginFailed(h.DB, c, req.Email)
		h.registerFailedLogin(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Check if user has TOTP enabled
	var totpConfig models.TwoFactorConfig
	if err := h.DB.WithContext(c.Request.Context()).Where("user_id = ? AND enabled = ?", user.ID, true).First(&totpConfig).Error; err == nil {
		// TOTP is enabled: check for trusted device
		if !IsTrustedDevice(c, h.DB, user.ID) {
			// Generate a short-lived pending token for TOTP verification
			pendingToken, err := totp.GeneratePendingToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
				})
				return
			}

			// Store hashed pending token in DB
			if err := h.DB.WithContext(c.Request.Context()).Create(&models.TOTPPendingToken{
				UserID:    user.ID,
				TokenHash: totp.HashToken(pendingToken),
				ExpiresAt: time.Now().Add(totp.PendingTokenExpiry),
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "TOKEN_ERROR", "message": "Failed to create verification session"},
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"totp_required": true,
					"pending_token": pendingToken,
				},
				"message": "Two-factor authentication required",
			})
			return
		}
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.DB.WithContext(c.Request.Context()).Model(&models.User{}).Where("id = ?", user.ID).
			Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil}).Error; err != nil {
			log.Printf("lockout: clearing the failure count for %s: %v", user.ID, err)
		}
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate tokens",
			},
		})
		return
	}

	// Set HttpOnly auth cookies for browser clients. Native mobile/desktop
	// clients ignore them and continue to use the Bearer header from the
	// tokens object below: both flows work.
	//
	// Record the refresh token as a server-side session so this device can be
	// listed and revoked later.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: successful sign-in lands in /system/activity at info
	// severity. IP + user-agent come from the request context inside
	// LogLogin so brute-force investigation has the full pair.
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Logged in successfully",
	})
}
`

var (
	// A pending 2FA token insert, checked or not.
	loginPendingInsertRe = regexp.MustCompile(`(?s)if err := h\.DB\.Create\(&models\.TOTPPendingToken\{\n(.*?)\n\}\)\.Error; err != nil \{\n.*?\nreturn\n\}`)
	// The failure-count reset, at the password or after 2FA, checked or not.
	loginCounterResetRe = regexp.MustCompile(`if user\.FailedLoginCount > 0 \|\| user\.LockedUntil != nil \{\n(?:if err := )?h\.DB\.Model\(&models\.User\{\}\)\.Where\("id = \?", user\.ID\)\.\nUpdates\(map\[string\]interface\{\}\{"failed_login_count": 0, "locked_until": nil\}\)(?:\.Error; err != nil \{\nlog\.Printf\([^\n]*\)\n\})?\n\}\n`)
	loginSpaceRunRe     = regexp.MustCompile(`\s+`)
)

const loginSignature = "func (h *AuthHandler) Login(c *gin.Context) {"

// canonicalLogin reduces a Login function to what decides its answers: no
// comments, blank lines or indentation, and none of the details in which the
// releases before v3.283.0 differ from one another.
func canonicalLogin(fn string) string {
	var lines []string
	for _, line := range strings.Split(fn, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		lines = append(lines, loginSpaceRunRe.ReplaceAllString(t, " "))
	}
	text := strings.Join(lines, "\n") + "\n"
	text = strings.ReplaceAll(text, ".WithContext(c.Request.Context())", "")
	text = loginPendingInsertRe.ReplaceAllString(text, "h.DB.Create(&models.TOTPPendingToken{\n$1\n})")
	return loginCounterResetRe.ReplaceAllString(text, "")
}

// loginFunc returns the Login function in src, from its signature to its closing
// brace, or "".
func loginFunc(src string) string {
	at := strings.Index(src, loginSignature)
	if at < 0 {
		return ""
	}
	end := goBraceEnd(src, at+len(loginSignature)-1)
	if end < 0 {
		return ""
	}
	return src[at:end]
}

// replaceMonolithicLogin swaps a pre-v3.283.0 Login for loginBlockNew, and
// reports false, changing nothing, unless that Login is one Grit wrote.
func replaceMonolithicLogin(src string) (string, bool) {
	fn := loginFunc(src)
	if fn == "" || canonicalLogin(fn) != canonicalLogin(loginMonolithic) {
		return src, false
	}
	// loginBlockNew declares these. A file that already has any of them was
	// changed by hand, and two declarations would not compile.
	for _, name := range []string{"loginRefusalReason", "func (h *AuthHandler) loginRefusal(", "func (h *AuthHandler) checkCredentials(",
		"func (h *AuthHandler) startTOTPChallenge(", "func invalidCredentials(", "func spendPasswordCheck(", "timingHashOnce"} {
		if strings.Contains(src, name) {
			return src, false
		}
	}
	out, ok := replaceGoFunc(src, loginSignature, loginBlockNew)
	if !ok {
		return src, false
	}
	// Grit wrote Login's comment twice, and replaceGoFunc takes only the copy
	// against the function. The other would now sit above invalidCredentials.
	first := loginBlockNew[:strings.Index(loginBlockNew, "\n")+1]
	out = strings.Replace(out, "// Login authenticates a user and returns tokens.\n\n"+first, first, 1)
	return out, true
}

// addImportBefore puts path in the import block before the line anchor, or as a
// group of its own when the anchor is not there. It reports false only when the
// file has no import block.
func addImportBefore(src, path, anchor string) (string, bool) {
	start := strings.Index(src, "\nimport (\n")
	if start < 0 {
		return src, false
	}
	end := strings.Index(src[start:], "\n)\n")
	if end < 0 {
		return src, false
	}
	block := src[start : start+end+1]
	line := "\t\"" + path + "\"\n"
	if strings.Contains(block, "\n"+line) {
		return src, true
	}
	if i := strings.Index(block, "\n"+anchor); i >= 0 {
		at := start + i + 1
		return src[:at] + line + src[at:], true
	}
	return addImportGroup(src, path)
}

// repairRealtimeTokenErrorSource stops the socket handshake echoing why a token
// was refused.
func repairRealtimeTokenErrorSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "RealtimeHandler") || !strings.Contains(src, `"message": err.Error()`) {
		return src, nil, nil
	}
	var out string
	switch {
	case strings.Count(src, realtimeInvalidTokenOld286) == 1:
		out = strings.Replace(src, realtimeInvalidTokenOld286, realtimeInvalidTokenNew, 1)
	case strings.Count(src, realtimeInvalidTokenOld283) == 1:
		out = strings.Replace(src, realtimeInvalidTokenOld283, realtimeInvalidTokenNew, 1)
	default:
		return src, nil, []string{"is not the handler Grit wrote: send a fixed message with INVALID_TOKEN rather than err.Error(), which tells the caller why their token was refused"}
	}
	realtimeImport := "/internal/realtime\"\n"
	end := strings.Index(out, realtimeImport)
	if end < 0 {
		return src, nil, []string{"could not find the realtime import: send a fixed message with INVALID_TOKEN rather than err.Error()"}
	}
	module := out[strings.LastIndex(out[:end], "\"")+1 : end]
	var ok bool
	if out, ok = addImportBefore(out, module+"/internal/respond", "\t\""+module+"/internal/services\"\n"); !ok {
		return src, nil, []string{"could not add the respond import: send a fixed message with INVALID_TOKEN rather than err.Error()"}
	}
	return out, []string{"a refused handshake token gets a fixed message, and the reason goes to the log"}, nil
}

// repairImportJobOwner records who started an import and lets only them, or an
// admin, read it. The model goes first: the handler and the importers need its
// field, and neither is changed without it.
func repairImportJobOwner(root string, opts Options, m *manifest.Manifest) error {
	apiRoot := opts.APIRoot(root)
	model := filepath.Join(apiRoot, "internal", "models", "import_job.go")
	handler := filepath.Join(apiRoot, "internal", "handlers", "import_job.go")
	if !fileExists(model) || !fileExists(handler) {
		return nil
	}
	if !fileContains(filepath.Join(apiRoot, "internal", "authz", "actor.go"), "func ActorOf(") {
		fmt.Println("  ⚠ internal/authz/actor.go is missing, so any signed-in user can still read another user's import job by its id.")
		return nil
	}
	if err := repairSourceFile(root, m, model, repairImportJobModelSource); err != nil {
		return err
	}
	if !fileContains(model, "CreatedBy string") {
		return nil
	}
	module := opts.Module()
	if err := repairSourceFile(root, m, handler, func(src string) (string, []string, []string) {
		return repairImportJobHandlerSource(src, module)
	}); err != nil {
		return err
	}
	// Importers generated before services owned the queries create the job in
	// the handler, so both directories are read.
	for _, dir := range []string{"services", "handlers"} {
		sources, err := goSources(filepath.Join(apiRoot, "internal", dir))
		if err != nil {
			return err
		}
		for _, path := range sources {
			if !strings.HasSuffix(path, "_import.go") {
				continue
			}
			if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
				return repairImportStarterSource(src, module)
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func repairImportJobModelSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "type ImportJob struct {") || strings.Contains(src, "CreatedBy ") {
		return src, nil, nil
	}
	if strings.Count(src, importJobMessageLine) != 1 {
		return src, nil, []string{"ImportJob is not the model Grit wrote: add CreatedBy string `gorm:\"size:36;index\"`, set it when an import starts, and scope GET /imports/:id to it, or any signed-in user can read another user's import"}
	}
	return strings.Replace(src, importJobMessageLine, importJobMessageLine+importJobCreatedByField, 1),
		[]string{"an import job records who started it"}, nil
}

func repairImportJobHandlerSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "ImportJobHandler") || strings.Contains(src, "created_by = ?") {
		return src, nil, nil
	}
	if strings.Count(src, importJobGetOld) != 1 {
		return src, nil, []string{"GetByID is not the one Grit wrote: scope the query to created_by unless the caller is an admin, or any signed-in user can read another user's import"}
	}
	out := strings.Replace(src, importJobGetOld, importJobGetNew, 1)
	var ok bool
	if out, ok = addImportBefore(out, module+"/internal/authz", "\t\""+module+"/internal/models\"\n"); !ok {
		return src, nil, []string{"could not add the authz import: scope GetByID to created_by by hand"}
	}
	return out, []string{"an import job is read only by whoever started it, or an admin"}, nil
}

func repairImportStarterSource(src, module string) (string, []string, []string) {
	if !oldStartImportRe.MatchString(src) {
		return src, nil, nil
	}
	if !strings.Contains(src, "ctx context.Context, total int) (*models.ImportJob, error) {") {
		return src, nil, []string{"this importer creates its job outside StartImport(ctx, total), so the job does not record who started it and only an admin can follow its progress: set CreatedBy: authz.UserIDFrom(ctx)"}
	}
	out := oldStartImportRe.ReplaceAllString(src, `job := models.ImportJob{Resource: "$1", Status: "processing", Total: total, CreatedBy: authz.UserIDFrom(ctx)}`)
	var ok bool
	if out, ok = addImportBefore(out, module+"/internal/authz", "\t\""+module+"/internal/imports\"\n"); !ok {
		return src, nil, []string{"could not add the authz import: set CreatedBy: authz.UserIDFrom(ctx) on the import job by hand"}
	}
	return out, []string{"the import job records who started it"}, nil
}

// ─── A field-encryption key for projects that have none ──────────────────────

// fieldKeyAdvice is the line grit upgrade prints for a project with two-factor
// sign-in and no FIELD_ENCRYPTION_KEY in .env.
const fieldKeyAdvice = "  ⚠ FIELD_ENCRYPTION_KEY is not set in .env, so two-factor secrets are stored in the clear. " +
	"Generate one with openssl rand -base64 32, back it up apart from the database, and run grit migrate to encrypt the secrets already stored."

// fieldKeyAdviceFor returns fieldKeyAdvice when it applies to the project at
// root, or "".
//
// grit new writes a key. An upgrade does not: a key that appears in .env
// without anybody choosing it is a key nobody backs up, and losing it locks out
// every two-factor user. A deployment may set the key outside .env, which this
// cannot see, so the line is advice rather than an error.
func fieldKeyAdviceFor(root string, opts Options) string {
	env := filepath.Join(root, ".env")
	if !fileExists(filepath.Join(opts.APIRoot(root), "internal", "models", "two_factor.go")) || !fileExists(env) {
		return ""
	}
	if strings.TrimSpace(readDotEnv(env)["FIELD_ENCRYPTION_KEY"]) != "" {
		return ""
	}
	return fieldKeyAdvice
}
