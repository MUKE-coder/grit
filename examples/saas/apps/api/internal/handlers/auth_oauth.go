package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

// Social login. The provider registry lives in services; this is the two
// endpoints the browser actually visits.

// users is the service the account writes below go through. Built per request,
// like every other accessor in this package, so routes.go did not have to change.
//
// The reads go through AuthService instead, which is where signing in already gets
// them: UserByEmail answers an unknown address the same way whatever the caller,
// which is the property this endpoint shares with the password step.
func (h *AuthHandler) users() *services.UserService {
	return &services.UserService{DB: h.DB}
}

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

	// A provider that shares no address cannot be matched to an account.
	if gothUser.Email == "" {
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your provider did not share an email address."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Find or create user by email
	var user models.User
	existing, lookupErr := h.AuthService.UserByEmail(c.Request.Context(), gothUser.Email)
	if existing != nil {
		user = *existing
	}

	if lookupErr != nil {
		if errors.Is(lookupErr, gorm.ErrRecordNotFound) {
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

			// Through CreateUser, so two provider sign-ins for one new address in the
			// same instant are settled by the unique index rather than by which read
			// finished first.
			if err := services.CreateUser(c.Request.Context(), h.DB, &user); err != nil {
				log.Printf("OAuth: failed to create user: %v", err)
				redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Failed to create account."))
				c.Redirect(http.StatusTemporaryRedirect, redirectURL)
				return
			}
		} else {
			log.Printf("OAuth: database error: %v", lookupErr)
			redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Something went wrong."))
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			return
		}
	} else {
		// An account whose address was never confirmed may have been registered by
		// someone else, waiting for the owner to sign in with the provider. The
		// provider has now shown who owns the address, so the unconfirmed password
		// stops working and its sessions end before the provider is linked.
		if user.EmailVerifiedAt == nil {
			verifiedAt := time.Now()
			if err := h.users().Update(c.Request.Context(), &user, map[string]interface{}{"password": "", "email_verified_at": verifiedAt}, ""); err != nil {
				log.Printf("oauth: securing unverified account %s: %v", user.ID, err)
				redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Something went wrong."))
				c.Redirect(http.StatusTemporaryRedirect, redirectURL)
				return
			}
			if err := services.RevokeAllUserSessions(h.DB, user.ID, ""); err != nil {
				log.Printf("oauth: revoking sessions of unverified account %s: %v", user.ID, err)
			}
			user.Password = ""
			user.EmailVerifiedAt = &verifiedAt
		}
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
			if err := h.users().Update(c.Request.Context(), &user, updates, ""); err != nil {
				log.Printf("oauth: linking %s to user %s: %v", provider, user.ID, err)
			}
		}
	}

	if !user.Active {
		redirectURL := fmt.Sprintf("%s/login?error=%s", h.Config.OAuthFrontendURL, url.QueryEscape("Your account has been disabled."))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// What the provider knew, handed to whatever this app wants to do with it:
	// the GitHub login, the avatar, the locale. See services.OnOAuthLogin.
	services.RunOAuthLoginHooks(c.Request.Context(), &user, gothUser)

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
