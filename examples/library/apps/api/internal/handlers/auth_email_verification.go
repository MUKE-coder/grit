package handlers

import (
	"context"
	"library/apps/api/internal/mail"
	"library/apps/api/internal/models"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
	"log"
	"net/http"
	"net/url"
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

	user, err := h.AuthService.UserByID(c.Request.Context(), userID)
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "User not found")
		return
	}

	if user.EmailVerifiedAt != nil {
		respond.Fail(c, respond.CodeAlreadyVerified, "This email is already verified")
		return
	}

	link := h.deliverVerificationEmail(c.Request.Context(), *user)

	body := gin.H{"message": "Verification email sent. The link is valid for 48 hours."}

	// In development the link comes back with the response, so the banner can
	// offer it as something to click.
	//
	// Verifying an email in development otherwise means opening the log or
	// digging a .html file out of storage/mail, which is enough friction that
	// the warning banner simply stays up forever in every dev project. Guarded
	// on APP_ENV: anywhere else this would hand a working verification link to
	// whoever asked, which is the whole secret.
	if h.Config.AppEnv == "development" && link != "" {
		body["verify_url"] = link
		body["message"] = "Development: the link is below rather than in your inbox."
	}

	c.JSON(http.StatusOK, body)
}

// The token from a verification link.

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}

	if _, err := services.ConsumeEmailVerificationToken(h.DB, req.Token); err != nil {
		// One message for expired, spent, unknown and address-changed. Telling
		// them apart tells an attacker which tokens once existed.
		respond.Fail(c, respond.CodeInvalidLink, "That verification link is invalid or has expired. Request a new one.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}

// deliverVerificationEmail mints a token and queues the link.
//
// It used to run in a goroutine the request started, so a slow SMTP call could
// not hold the response open. The mail goes on the background queue now, which
// keeps the response just as short and, unlike a goroutine, retries a provider
// that is briefly down and survives a deploy. Without a queue the send is
// inline, which is what development without Redis does.
func (h *AuthHandler) deliverVerificationEmail(ctx context.Context, user models.User) string {
	token, err := services.GenerateVerificationToken()
	if err != nil {
		log.Printf("email verification: generating token for %s: %v", user.Email, err)
		return ""
	}

	if _, err := services.CreateEmailVerificationToken(h.DB, user.ID, user.Email, token); err != nil {
		log.Printf("email verification: storing token for %s: %v", user.Email, err)
		return ""
	}

	verifyURL := strings.TrimSuffix(h.Config.OAuthFrontendURL, "/") + "/verify-email?token=" + url.QueryEscape(token)

	if h.Jobs != nil || h.Mailer != nil {
		// Keyed on the token, so one enqueue per link: a client that retries
		// the request mints a new token and gets a new key, and a proxy that
		// replays the same one does not send the mail twice.
		if err := dispatchMail(ctx, h.Mailer, h.Jobs, "verify:"+user.ID+":"+token, mail.SendOptions{
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
		return verifyURL
	}

	if h.Config.AppEnv == "production" {
		log.Printf("email verification: NO MAILER CONFIGURED: %s cannot receive a link. Set MAIL_MAILER (see .env.example).", user.Email)
		return ""
	}

	log.Printf("email verification link for %s: %s", user.Email, verifyURL)
	return verifyURL
}

// ResetPassword resets a user's password with a valid token.
