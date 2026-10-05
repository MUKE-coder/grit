package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
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

	unlocked, err := h.users().ClearLockout(c.Request.Context(), id)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Failed to unlock the account")
		return
	}
	if !unlocked {
		respond.Fail(c, respond.CodeNotFound, "User not found")
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
// The threshold and the window are this project's configuration, so the decision
// is here; the counting and the locking are AuthService's, which is also what the
// second factor calls. There were two copies of those writes before, one per
// factor.
func (h *AuthHandler) registerFailedLogin(user *models.User) {
	max := h.Config.LoginMaxAttempts
	if max <= 0 {
		return // lockout disabled
	}
	ctx := context.Background()

	count, err := h.AuthService.CountLoginFailure(ctx, user.ID)
	if err != nil {
		log.Printf("lockout: counting a failure for %s: %v", user.ID, err)
		return
	}
	if count < max {
		return
	}

	until := time.Now().Add(h.Config.LoginLockoutWindow)
	if err := h.AuthService.LockAccount(ctx, user.ID, until); err != nil {
		log.Printf("lockout: locking %s: %v", user.ID, err)
		return
	}
	log.Printf("lockout: %s locked until %s after %d failed attempts", user.Email, until.Format(time.RFC3339), max)
}

// SendVerificationEmail issues a fresh verification link for the signed-in
// user. Authenticated on purpose: an unauthenticated "send a link to this
// address" endpoint is a spam cannon aimed at whoever you name.
