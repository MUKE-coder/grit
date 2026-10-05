package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"commerce/apps/api/internal/mail"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// Signing in with a link.
//
// These hang off AuthHandler rather than a handler of their own, and that is
// the point: consuming a link has to end exactly where a password sign-in
// ends, through the same second-factor challenge, the same session record and
// the same activity log. A second handler with its own copy of that ending is
// how one of them quietly stops matching the other.

// MagicLinkRequest asks for a sign-in link.
type MagicLinkRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestMagicLink emails a sign-in link.
//
//	POST /api/v1/auth/magic-link
//
// Always answers the same, whether or not the address has an account. A form
// that says "no such user" is a form for finding out who is registered here.
func (h *AuthHandler) RequestMagicLink(c *gin.Context) {
	var req MagicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	const sameAnswer = "If that address has an account, a sign-in link is on its way."

	if h.Mailer == nil && h.Jobs == nil {
		// Nothing can be sent, and "on its way" would be a lie that leaves
		// somebody watching an inbox.
		respond.Fail(c, respond.CodeMailFailed, "This deployment cannot send email, so sign-in links are unavailable.")
		return
	}

	// Matched exactly as typed, like every other lookup in this package.
	// Lowercasing here and nowhere else would mean an address that signs in
	// with a password cannot be found by a link, and the form answers the same
	// either way, so nobody would ever see why.
	var user models.User
	err := h.DB.WithContext(c.Request.Context()).
		Where("email = ?", req.Email).First(&user).Error
	if err != nil || !user.Active {
		// Same answer, same shape, no work done.
		c.JSON(http.StatusOK, gin.H{"message": sameAnswer})
		return
	}

	token, err := services.IssueMagicLink(c.Request.Context(), h.DB, user.ID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		if errors.Is(err, services.ErrMagicLinkTooSoon) {
			// The one case worth saying out loud: it is their own address, they
			// are staring at their inbox, and "on its way" would have them
			// waiting for a second email that is not coming.
			respond.Fail(c, respond.CodeRateLimited, "A link was sent a moment ago. Check your email, or try again in a minute.")
			return
		}
		log.Printf("magic link: issuing for %s: %v", user.ID, err)
		c.JSON(http.StatusOK, gin.H{"message": sameAnswer})
		return
	}

	link := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/magic-link?token=" + token
	if err := dispatchMail(c.Request.Context(), h.Mailer, h.Jobs, "magic-link:"+services.HashMagicLink(token), mail.SendOptions{
		To:       user.Email,
		Subject:  "Your sign-in link",
		Template: "magic-link",
		Data: map[string]interface{}{
			"AppName": h.Config.AppName,
			"Title":   "Your sign-in link",
			"Link":    link,
			"Minutes": int(services.MagicLinkExpiry.Minutes()),
			"Year":    time.Now().Year(),
		},
	}); err != nil {
		log.Printf("magic link: emailing %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": sameAnswer})
}

// MagicLinkConsumeRequest spends one.
type MagicLinkConsumeRequest struct {
	Token string `json:"token" binding:"required"`
}

// ConsumeMagicLink spends a link and signs the person in.
//
//	POST /api/v1/auth/magic-link/consume
//
// A POST, not the GET that opened the link: mail scanners follow every URL in
// a message, and a token spent by a GET is burned before the person has read
// the email.
func (h *AuthHandler) ConsumeMagicLink(c *gin.Context) {
	var req MagicLinkConsumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	userID, err := services.ConsumeMagicLink(c.Request.Context(), h.DB, req.Token)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMagicLinkUsed):
			respond.Fail(c, respond.CodeInvalidLink, "This link has already been used. Ask for a new one.")
		case errors.Is(err, services.ErrMagicLinkExpired):
			respond.Fail(c, respond.CodeInvalidLink, "This link has expired. Ask for a new one.")
		default:
			respond.Fail(c, respond.CodeInvalidLink, "This link is not valid. Ask for a new one.")
		}
		return
	}

	var user models.User
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
		respond.Fail(c, respond.CodeInvalidLink, "This link is not valid. Ask for a new one.")
		return
	}
	if !user.Active {
		respond.Fail(c, respond.CodeAccountDisabled, "Your account has been disabled.")
		return
	}

	// A link replaces the password, not the second factor: a link sitting in a
	// mailbox is exactly what a second factor exists to survive.
	if h.startTOTPChallenge(c, &user) {
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		respond.Fail(c, respond.CodeTokenError, "Failed to generate tokens")
		return
	}
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("magic link: failed to record session for %s: %v", user.ID, err)
	}
	h.AuthService.SetAuthCookies(c, tokens)
	services.LogLogin(h.DB, c, user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":   user,
			"tokens": tokens,
		},
		"message": "Signed in",
	})
}

// MagicLinkActivity is one row of "who has been asking for links to my account".
type MagicLinkActivity struct {
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	ExpiresAt time.Time  `json:"expires_at"`
	IPAddress string     `json:"ip_address,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
}

// RecentMagicLinks lists the last few sign-in links issued for this account.
//
//	GET /api/v1/auth/magic-link/recent
//
// A link is a bearer credential sitting in a mailbox, and the person whose
// mailbox it is has no other way to notice that somebody keeps asking for one.
// The request form cannot tell them: it answers the same to everybody by
// design. This is where they find out.
//
// No token or hash is returned, only when and from where.
func (h *AuthHandler) RecentMagicLinks(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		respond.Fail(c, respond.CodeUnauthorized, "Sign in first")
		return
	}

	var rows []models.MagicLinkToken
	if err := h.DB.WithContext(c.Request.Context()).
		Where("user_id = ?", userID).
		Order("created_at DESC").Limit(5).Find(&rows).Error; err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not read your sign-in links")
		return
	}

	out := make([]MagicLinkActivity, 0, len(rows))
	for _, row := range rows {
		out = append(out, MagicLinkActivity{
			CreatedAt: row.CreatedAt,
			UsedAt:    row.UsedAt,
			ExpiresAt: row.ExpiresAt,
			IPAddress: row.IPAddress,
			UserAgent: row.UserAgent,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
