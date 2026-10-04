package handlers

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library/apps/api/internal/config"
	"library/apps/api/internal/jobs"
	"library/apps/api/internal/mail"
	"library/apps/api/internal/models"
	"library/apps/api/internal/password"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
	"library/apps/api/internal/totp"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Config      *config.Config
	// Mailer is optional. Without it the reset link is logged instead of sent,
	// which is what you want in dev and must never be what happens in prod —
	// ForgotPassword refuses to log the token when APP_ENV is production.
	Mailer *mail.Mailer
	// Jobs is optional too: it is nil when the project runs without Redis.
	// When it is there, every email this handler sends goes onto the queue
	// instead of a goroutine, so it is retried and survives a restart.
	Jobs *jobs.Client
}

// AuthResponse documents the body returned by register, login and refresh.
//
// The handlers emit gin.H rather than this struct, so nothing enforces the two
// agree — if you change what an auth handler writes, change this with it. It
// exists because a reference that says "No Body" is worse than no reference.
type AuthResponse struct {
	Data struct {
		Tokens struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    int64  `json:"expires_at"`
		} `json:"tokens"`
		User models.User `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

// The types below exist so the API reference can show a body instead of
// "No Body". The handlers emit gin.H maps, so nothing enforces that these stay
// in step — if you change what a handler writes, change its type here too.
// They are documentation with a compiler attached, which is still better than
// prose nobody updates.

// MessageResponse is the plain acknowledgement shape.
type MessageResponse struct {
	Message string `json:"message"`
}

// IssuedKeyResponse is returned once, when an API key is created.
type IssuedKeyResponse struct {
	Data struct {
		Key   models.APIKey `json:"key"`
		Token string        `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
}

// TOTPStatusResponse describes the caller's two-factor state.
type TOTPStatusResponse struct {
	Data struct {
		Enabled              bool  `json:"enabled"`
		BackupCodesRemaining int   `json:"backup_codes_remaining"`
		TrustedDevices       int64 `json:"trusted_devices"`
	} `json:"data"`
}

// PresignResponse carries the URL a browser PUTs to, and the key to send back
// to /uploads/complete afterwards.
type PresignResponse struct {
	Data struct {
		PresignedURL string `json:"presigned_url"`
		Key          string `json:"key"`
	} `json:"data"`
}

// ChainStatusResponse is the activity-log integrity verdict.
type ChainStatusResponse struct {
	Valid        bool   `json:"valid"`
	TotalEntries int    `json:"total_entries"`
	BrokenAt     int    `json:"broken_at,omitempty"`
	BrokenAtID   string `json:"broken_at_id,omitempty"`
	Expected     string `json:"expected,omitempty"`
	Got          string `json:"got,omitempty"`
	Message      string `json:"message,omitempty"`
}

// ErrorResponse is the error envelope every endpoint uses.
type ErrorResponse struct {
	Error struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Details map[string]string `json:"details,omitempty"`
	} `json:"error"`
}

type RegisterRequest struct {
	FirstName  string `json:"first_name" binding:"required,min=2"`
	LastName   string `json:"last_name" binding:"required,min=2"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	MACAddress string `json:"mac_address"` // optional — provided by client if available
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register creates a new user account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// The same four rules the admin shows while somebody types, checked here
	// because a checklist the server ignores is decoration. Their own name and
	// address are passed in: "ada@example.com" with the password "ada2026" is
	// the commonest weak password there is, and only the server knows both.
	if failed := password.Check(req.Password, req.Email, req.FirstName, req.LastName); len(failed) > 0 {
		respond.Fail(c, respond.CodeValidationError, password.Message(failed))
		return
	}

	user := models.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   req.Password,
		Role:       models.RoleUser,
		Active:     true,
		IPAddress:  c.ClientIP(),
		MACAddress: req.MACAddress,
	}

	// Same insert as the admin's Create User, through the same helper. The
	// check-then-insert this replaces was written out twice, and both copies
	// had the same gap between the SELECT and the INSERT: two signups for one
	// address in the same instant both passed the check, and the loser was
	// told "Failed to create user" with a 500.
	if err := services.CreateUser(c.Request.Context(), h.DB, &user); err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "EMAIL_EXISTS",
					"message": "A user with this email already exists",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create user",
			},
		})
		return
	}

	// Off the request path: signup should not wait on SMTP, and a mail failure
	// must not fail an account that was created successfully.
	_ = h.deliverVerificationEmail(c.Request.Context(), user)

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}

	// Record the refresh token as a server-side session so it can be revoked.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("auth: failed to record session for %s: %v", user.ID, err)
	}

	// Set HttpOnly auth cookies for browser clients.
	h.AuthService.SetAuthCookies(c, tokens)

	// v3.30.1: emit a semantic activity row so /system/activity reflects
	// the signup. Non-blocking — a logging failure won't fail the
	// register request.
	services.LogRegister(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "User registered successfully",
	})
}

// invalidCredentials is the one answer a sign-in gets until its password has
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
	// Which columns are safe to read here is a decision about the data, so it
	// is the service's: see EnabledTwoFactor.
	totpConfig, err := h.AuthService.EnabledTwoFactor(c.Request.Context(), user.ID)
	if err != nil {
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

	pending := models.TOTPPendingToken{
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

	if err := h.AuthService.StartPendingTOTP(c.Request.Context(), &pending); err != nil {
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
}

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	user, err := h.AuthService.UserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// v3.30.1: unknown email is the most common brute-force fingerprint;
		// surface it in /system/activity as "warn" severity so operators
		// can spot credential-stuffing spikes.
		spendPasswordCheck(req.Password)
		services.LogLoginFailed(h.DB, c, req.Email)
		invalidCredentials(c)
		return
	}

	// Nothing about the account is reported until its password is right.
	if !h.checkCredentials(c, user, req.Password) {
		return
	}

	// The password matched, so the caller may be told why the account still
	// cannot sign in.
	if refusal := h.loginRefusal(user); refusal != nil {
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
	if h.startTOTPChallenge(c, user) {
		return
	}

	// The failure count is cleared when the sign-in completes. Cleared at the
	// password, with 2FA still ahead, it handed a fresh set of guesses at the
	// code to anyone who knew the password and simply signed in again.
	if user.FailedLoginCount > 0 || user.LockedUntil != nil {
		if err := h.AuthService.ClearLoginFailures(c.Request.Context(), user.ID); err != nil {
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

// Refresh generates a new access token from a refresh token. The token is
// read from the grit_refresh cookie first (web client) and falls back to
// the JSON body (mobile/desktop bearer clients) — so a single endpoint
// supports both flows.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshToken string
	if cookieValue, err := c.Cookie("grit_refresh"); err == nil && cookieValue != "" {
		refreshToken = cookieValue
	} else {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respond.Fail(c, respond.CodeValidationError, err.Error())
			return
		}
		refreshToken = req.RefreshToken
	}

	claims, err := h.AuthService.ValidateRefreshToken(refreshToken)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidToken, "Invalid or expired refresh token")
		return
	}

	// Re-verify the account on every refresh. A stateless refresh token is
	// otherwise valid for its full lifetime even after the user is deleted or
	// deactivated; re-loading the user closes that window and lets a role
	// change take effect on the next refresh (partial revocation without a
	// server-side token store).
	user, err := h.AuthService.UserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidToken, "Account no longer exists")
		return
	}
	if !user.Active {
		respond.Fail(c, respond.CodeAccountDisabled, "This account has been disabled")
		return
	}

	// Keep the session id through rotation, so the new access token names the
	// row the old one did. A refresh token from before token types carries none,
	// and its row is found by the token instead.
	sessionID := claims.SessionID
	if sessionID == "" {
		sessionID = services.SessionIDForToken(h.DB, refreshToken)
	}
	tokens, err := h.AuthService.GenerateSessionTokenPair(user.ID, user.Email, user.Role, sessionID)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}

	// Refresh the HttpOnly cookies so the new access token lands in the
	// browser without any JS handling. The bearer JSON path is unchanged
	// for native clients.
	//
	// Rotate the session. This is where revocation actually bites: a session
	// that was revoked, idled out, aged past its absolute limit, or whose token
	// was replayed after rotation has no live row, and the refresh is refused.
	if _, err := services.RotateSession(h.DB, c, refreshToken, tokens.RefreshToken); err != nil {
		h.AuthService.ClearAuthCookies(c)
		respond.Fail(c, respond.CodeSessionRevoked, "This session is no longer valid. Please sign in again.")
		return
	}
	h.AuthService.SetAuthCookies(c, tokens)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tokens": tokens,
		},
		"message": "Token refreshed successfully",
	})
}

// Logout invalidates the user's session. Cookies are cleared immediately;
// native bearer clients should also drop their stored tokens client-side.
func (h *AuthHandler) Logout(c *gin.Context) {
	// v3.30.1: read the user out of context BEFORE clearing cookies so
	// the activity row carries the right email. The auth middleware set
	// "user" on the gin context when the request came in.
	var actorID, actorEmail string
	if v, ok := c.Get("user"); ok {
		if u, ok := v.(models.User); ok {
			actorID = u.ID
			actorEmail = u.Email
		}
	}

	// Revoke the server-side session BEFORE clearing cookies — once they're
	// gone we can no longer identify which session to kill. This is what makes
	// logout real: the refresh token is dead immediately, not merely forgotten
	// by this browser.
	if rt, err := c.Cookie("grit_refresh"); err == nil && rt != "" {
		if err := services.RevokeSessionByToken(h.DB, rt); err != nil {
			log.Printf("auth: failed to revoke session on logout: %v", err)
		}
	}

	h.AuthService.ClearAuthCookies(c)

	if actorID != "" {
		services.LogLogout(h.DB, c, actorID, actorEmail)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user.
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		respond.Fail(c, respond.CodeUnauthorized, "Not authenticated")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// ForgotPassword initiates a password reset.

// The token from a verification link.
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmail consumes a verification token. Public — the user clicks this
// from their mail client, where they are usually not signed in.
