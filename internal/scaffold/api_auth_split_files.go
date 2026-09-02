package scaffold

// The auth handler, split by concern.
//
// auth.go passed a thousand lines covering sessions, password reset, email
// verification, OAuth and account lockout at once, and it is a file people
// customise. Each flow now has its own file, so changing one means opening one.
// auth.go keeps the session core: register, login, refresh, logout, me.

// apiAuthPasswordResetGo emits handlers/auth_password_reset.go.
//
// Carved out of auth.go, which was a thousand lines holding five unrelated
// flows. Changing the reset email or the token lifetime meant scrolling past
// OAuth and session refresh to find them.
func apiAuthPasswordResetGo() string {
	return `package handlers

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Password reset: request a link, and redeem it.
// Split out of auth.go so changing the email or the token lifetime does not
// mean reading past OAuth and session refresh to find them.

// ForgotPassword initiates a password reset.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// One response for every outcome. Any variation — a different message, a
	// different status, a measurably different latency — turns this endpoint
	// into an oracle for which email addresses hold accounts.
	const genericResponse = "If an account with that email exists, a password reset link has been sent"

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": genericResponse})
		return
	}

	// Everything past the lookup — minting the token, storing it, delivering the
	// link — runs off the request path. Both branches then do the same work
	// before answering (parse, one indexed SELECT), so a registered address does
	// not take measurably longer to respond than an unregistered one. Identical
	// wording with a distinguishable response time is still an oracle.
	//
	// c.ClientIP() is read here: the gin context must not be touched once the
	// handler has returned.
	go h.deliverPasswordReset(user, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"message": genericResponse})
}

// deliverPasswordReset issues a reset token and sends the link. It runs in its
// own goroutine, so it owns its context and reports failures only to the log —
// there is no caller left to tell, and telling the original one would have
// confirmed the address exists.

// deliverPasswordReset issues a reset token and sends the link. It runs in its
// own goroutine, so it owns its context and reports failures only to the log —
// there is no caller left to tell, and telling the original one would have
// confirmed the address exists.
func (h *AuthHandler) deliverPasswordReset(user models.User, clientIP string) {
	token, err := services.GenerateResetToken()
	if err != nil {
		log.Printf("password reset: generating token for %s: %v", user.Email, err)
		return
	}

	if _, err := services.CreatePasswordResetToken(h.DB, user.ID, token, clientIP); err != nil {
		log.Printf("password reset: storing token for %s: %v", user.Email, err)
		return
	}

	resetURL := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/reset-password?token=" + url.QueryEscape(token)

	if h.Mailer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.Mailer.Send(ctx, mail.SendOptions{
			To:       user.Email,
			Subject:  "Reset your password",
			Template: "password-reset",
			Data: map[string]interface{}{
				"AppName":  h.Config.AppName,
				"Title":    "Reset your password",
				"Message":  "We received a request to reset your password. This link expires in one hour and can only be used once. If you didn't ask for this, you can ignore this email.",
				"ResetURL": resetURL,
				"Year":     time.Now().Year(),
			},
		}); err != nil {
			log.Printf("password reset: sending email to %s: %v", user.Email, err)
		}
		return
	}

	if h.Config.AppEnv == "production" {
		// No mailer in production means nobody can complete a reset. Say so
		// loudly rather than printing a working token into the log — a live
		// reset link in a log file is a credential.
		log.Printf("password reset: NO MAILER CONFIGURED: %s cannot receive a reset link. Set RESEND_API_KEY.", user.Email)
		return
	}

	// Dev convenience only, and only outside production.
	log.Printf("password reset link for %s: %s", user.Email, resetURL)
}

// Unlock clears a lockout early. Waiting out the window is the normal path;
// this exists for the support call that follows a user locking themselves out
// five minutes before a demo.

// ResetPassword resets a user's password with a valid token.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Consume first. The token is single-use and burning it before doing any
	// work means a failure later can't leave a still-valid token behind.
	userID, err := services.ConsumePasswordResetToken(h.DB, req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "This reset link is invalid or has expired. Request a new one.",
			},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to hash password",
			},
		})
		return
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update password",
			},
		})
		return
	}

	// The reason someone resets a password is to evict whoever they think is in
	// their account. Leaving that person's session alive would defeat the entire
	// exercise, so every device is signed out — including any the attacker holds.
	if err := services.RevokeAllUserSessions(h.DB, userID, ""); err != nil {
		log.Printf("password reset: revoking sessions for %s: %v", userID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully. Please sign in with your new password.",
	})
}

// OAuthBegin redirects the user to the OAuth provider's consent screen.
`
}

// apiAuthEmailVerificationGo emits handlers/auth_email_verification.go.
func apiAuthEmailVerificationGo() string {
	return `package handlers

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Email verification: send the link, and redeem it.

// SendVerificationEmail issues a fresh verification link for the signed-in
// user. Authenticated on purpose: an unauthenticated "send a link to this
// address" endpoint is a spam cannon aimed at whoever you name.
func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	if user.EmailVerifiedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "ALREADY_VERIFIED", "message": "This email is already verified"},
		})
		return
	}

	go h.deliverVerificationEmail(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent. The link is valid for 48 hours.",
	})
}

// The token from a verification link.

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if _, err := services.ConsumeEmailVerificationToken(h.DB, req.Token); err != nil {
		// One message for expired, spent, unknown and address-changed. Telling
		// them apart tells an attacker which tokens once existed.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "That verification link is invalid or has expired. Request a new one.",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}

// deliverVerificationEmail mints a token and sends the link, off the request
// path so a slow SMTP call cannot hold the response open.

// deliverVerificationEmail mints a token and sends the link, off the request
// path so a slow SMTP call cannot hold the response open.
func (h *AuthHandler) deliverVerificationEmail(user models.User) {
	token, err := services.GenerateVerificationToken()
	if err != nil {
		log.Printf("email verification: generating token for %s: %v", user.Email, err)
		return
	}

	if _, err := services.CreateEmailVerificationToken(h.DB, user.ID, user.Email, token); err != nil {
		log.Printf("email verification: storing token for %s: %v", user.Email, err)
		return
	}

	verifyURL := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/verify-email?token=" + url.QueryEscape(token)

	if h.Mailer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.Mailer.Send(ctx, mail.SendOptions{
			To:       user.Email,
			Subject:  "Confirm your email address",
			Template: "email-verification",
			Data: map[string]interface{}{
				"AppName":   h.Config.AppName,
				"Title":     "Confirm your email address",
				"Message":   "Click the button below to confirm this address. The link expires in 48 hours and can only be used once.",
				"VerifyURL": verifyURL,
				"Year":      time.Now().Year(),
			},
		}); err != nil {
			log.Printf("email verification: sending to %s: %v", user.Email, err)
		}
		return
	}

	if h.Config.AppEnv == "production" {
		log.Printf("email verification: NO MAILER CONFIGURED: %s cannot receive a link. Set RESEND_API_KEY.", user.Email)
		return
	}

	log.Printf("email verification link for %s: %s", user.Email, verifyURL)
}

// ResetPassword resets a user's password with a valid token.
`
}

// apiAuthOAuthGo emits handlers/auth_oauth.go. The provider registry lives in
// services; these are the two endpoints a browser actually visits.
func apiAuthOAuthGo() string {
	return `package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

// Social login. The provider registry lives in services; this is the two
// endpoints the browser actually visits.

// OAuthBegin redirects the user to the OAuth provider's consent screen.
func (h *AuthHandler) OAuthBegin(c *gin.Context) {
	provider := c.Param("provider")

	// Gothic reads provider from query string, not URL params
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// OAuthCallback completes the OAuth flow, finds or creates the user, and redirects with JWT tokens.

// OAuthCallback completes the OAuth flow, finds or creates the user, and redirects with JWT tokens.
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")

	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("OAuth callback error: %v", err)
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Authentication failed. Please try again."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Find or create user by email
	var user models.User
	result := h.DB.Where("email = ?", gothUser.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user from OAuth data
			now := time.Now()
			user = models.User{
				FirstName:       gothUser.FirstName,
				LastName:        gothUser.LastName,
				Email:           gothUser.Email,
				Avatar:          gothUser.AvatarURL,
				Provider:        provider,
				Active:          true,
				EmailVerifiedAt: &now,
				IPAddress:       c.ClientIP(),
			}

			if provider == "google" {
				user.GoogleID = gothUser.UserID
			} else if provider == "github" {
				user.GithubID = gothUser.UserID
			}

			// If name is empty, try to use NickName
			if user.FirstName == "" && gothUser.NickName != "" {
				user.FirstName = gothUser.NickName
			}
			if user.FirstName == "" {
				user.FirstName = "User"
			}
			if user.LastName == "" {
				user.LastName = ""
			}

			if err := h.DB.Create(&user).Error; err != nil {
				log.Printf("OAuth: failed to create user: %v", err)
				redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Failed to create account."))
				c.Redirect(http.StatusTemporaryRedirect, redirectURL)
				return
			}
		} else {
			log.Printf("OAuth: database error: %v", result.Error)
			redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Something went wrong."))
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		}
	} else {
		// Link OAuth provider to existing account
		updates := map[string]interface{}{}
		if provider == "google" && user.GoogleID == "" {
			updates["google_id"] = gothUser.UserID
		} else if provider == "github" && user.GithubID == "" {
			updates["github_id"] = gothUser.UserID
		}
		if user.Avatar == "" && gothUser.AvatarURL != "" {
			updates["avatar"] = gothUser.AvatarURL
		}
		if user.Provider == "local" {
			updates["provider"] = provider
		}

		if len(updates) > 0 {
			h.DB.Model(&user).Updates(updates)
		}
	}

	if !user.Active {
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your account has been disabled."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Generate JWT tokens
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("OAuth: failed to generate tokens: %v", err)
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Failed to sign in."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Record the refresh token as a server-side session, so an OAuth login is
	// listed and revocable exactly like a password login.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("OAuth: failed to record session for %s: %v", user.ID, err)
	}

	// Set HttpOnly auth cookies BEFORE redirecting so the browser stores
	// them as part of this same response. The callback page then just
	// navigates — no tokens in URL, no tokens in JS, no XSS exposure.
	h.AuthService.SetAuthCookies(c, tokens)

	// Redirect to frontend callback. No query params — tokens travel as
	// HttpOnly Set-Cookie headers on this 307 response.
	redirectURL := fmt.Sprintf("%s/auth/callback", h.Config.OAuthFrontendURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
`
}

// apiAuthLockoutGo emits handlers/auth_lockout.go: the failed-login counter
// and the admin unlock that clears it.
func apiAuthLockoutGo() string {
	return `package handlers

import (
	"log"
	"net/http"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Failed-login counting and the admin unlock that clears it.
//
// Unlock hangs off UserHandler rather than AuthHandler because it is an
// administrative action on a user, not something a signed-out person does.

// Unlock clears a lockout early. Waiting out the window is the normal path;
// this exists for the support call that follows a user locking themselves out
// five minutes before a demo.
func (h *UserHandler) Unlock(c *gin.Context) {
	id := c.Param("id")

	res := h.DB.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"locked_until": nil, "failed_login_count": 0})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to unlock the account"},
		})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "user.unlock",
		Severity:     "warn",
		Summary:      "Account lockout cleared by an administrator",
		ResourceType: "user",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Account unlocked"})
}

// registerFailedLogin counts a wrong password against the account and locks it
// once the threshold is reached.
//
// Only wrong-password-on-a-real-account is counted. Counting unknown emails
// would let anyone lock an address they can guess, which turns a defence into
// a denial-of-service tool.
//
// The increment is a single UPDATE rather than read-modify-write, so parallel
// attempts cannot each read the same count and overwrite one another.

// registerFailedLogin counts a wrong password against the account and locks it
// once the threshold is reached.
//
// Only wrong-password-on-a-real-account is counted. Counting unknown emails
// would let anyone lock an address they can guess, which turns a defence into
// a denial-of-service tool.
//
// The increment is a single UPDATE rather than read-modify-write, so parallel
// attempts cannot each read the same count and overwrite one another.
func (h *AuthHandler) registerFailedLogin(user *models.User) {
	max := h.Config.LoginMaxAttempts
	if max <= 0 {
		return // lockout disabled
	}

	if err := h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		UpdateColumn("failed_login_count", gorm.Expr("failed_login_count + 1")).Error; err != nil {
		log.Printf("lockout: incrementing failed_login_count for %s: %v", user.ID, err)
		return
	}

	var fresh models.User
	if err := h.DB.Select("id", "failed_login_count").First(&fresh, "id = ?", user.ID).Error; err != nil {
		return
	}
	if fresh.FailedLoginCount < max {
		return
	}

	until := time.Now().Add(h.Config.LoginLockoutWindow)
	if err := h.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{"locked_until": until, "failed_login_count": 0}).Error; err != nil {
		log.Printf("lockout: locking %s: %v", user.ID, err)
		return
	}
	log.Printf("lockout: %s locked until %s after %d failed attempts", user.Email, until.Format(time.RFC3339), max)
}

// SendVerificationEmail issues a fresh verification link for the signed-in
// user. Authenticated on purpose: an unauthenticated "send a link to this
// address" endpoint is a spam cannon aimed at whoever you name.
`
}
