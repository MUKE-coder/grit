package handlers

import (
	"commerce/apps/api/internal/mail"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/password"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
	"context"
	"log"
	"net/http"
	"net/url"
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
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// One response for every outcome. Any variation — a different message, a
	// different status, a measurably different latency — turns this endpoint
	// into an oracle for which email addresses hold accounts.
	const genericResponse = "If an account with that email exists, a password reset link has been sent"

	user, err := h.AuthService.UserByEmail(c.Request.Context(), req.Email)
	if err != nil {
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
	// handler has returned. detach caps how many of these can be in flight at
	// once; it used to be one unbounded goroutine per request.
	clientIP := c.ClientIP()
	detach(func(ctx context.Context) { h.deliverPasswordReset(ctx, *user, clientIP) })

	c.JSON(http.StatusOK, gin.H{"message": genericResponse})
}

// deliverPasswordReset issues a reset token and queues the link.
//
// It runs off the request path, through detach, so it owns its context and
// reports failures only to the log: there is no caller left to tell, and
// telling the original one how long the work took would confirm the address
// exists. The mail itself goes on the background queue, so a provider that is
// down for a minute no longer loses the reset link.
func (h *AuthHandler) deliverPasswordReset(ctx context.Context, user models.User, clientIP string) {
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

	if h.Jobs != nil || h.Mailer != nil {
		if err := dispatchMail(ctx, h.Mailer, h.Jobs, "password-reset:"+user.ID+":"+token, mail.SendOptions{
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
		log.Printf("password reset: NO MAILER CONFIGURED: %s cannot receive a reset link. Set MAIL_MAILER (see .env.example).", user.Email)
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
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	// Consume first. The token is single-use and burning it before doing any
	// work means a failure later can't leave a still-valid token behind.
	userID, err := services.ConsumePasswordResetToken(h.DB, req.Token)
	if err != nil {
		respond.Fail(c, respond.CodeInvalidLink, "This reset link is invalid or has expired. Request a new one.")
		return
	}

	// The same rules as registering and changing. A reset is where a weak
	// password most often gets in: the person is locked out, in a hurry, and
	// picks something they will remember.
	if owner, err := h.AuthService.UserByID(c.Request.Context(), userID); err == nil {
		if failed := password.Check(req.Password, owner.Email, owner.FirstName, owner.LastName); len(failed) > 0 {
			respond.Fail(c, respond.CodeValidationError, password.Message(failed))
			return
		}
	} else if failed := password.Check(req.Password); len(failed) > 0 {
		respond.Fail(c, respond.CodeValidationError, password.Message(failed))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to hash password")
		return
	}

	if err := h.users().SetPassword(c.Request.Context(), userID, string(hashedPassword)); err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to update password")
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
