package scaffold

import (
	"fmt"
	"path/filepath"
)

// Signing in with a link, for the people who will never remember a password.
//
// Three decisions worth stating, because each of them is a way this feature is
// usually got wrong:
//
//   - Asking for a link answers the same whether or not the address has an
//     account. Anything else turns the form into a way to find out who is
//     registered here.
//   - The token is spent by a POST from the page the link opens, never by the
//     GET that opens it. Corporate mail scanners follow every link in a
//     message, and a token spent by a GET is a token the scanner burns before
//     the person has read the email.
//   - A magic link replaces the password, not the second factor. An account
//     with two-factor on gets the same challenge it always does.

// writeMagicLinkFiles writes the model, the service and the handler.
func writeMagicLinkFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "magic_link.go"):        magicLinkModelGo(),
		filepath.Join(apiRoot, "internal", "services", "magic_link.go"):      magicLinkServiceGo(opts.Module()),
		filepath.Join(apiRoot, "internal", "services", "magic_link_test.go"): magicLinkServiceTestGo(opts.Module()),
		filepath.Join(apiRoot, "internal", "handlers", "magic_link.go"):      magicLinkHandlerGo(opts.Module()),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
	}
	return nil
}

func magicLinkModelGo() string {
	return `package models

import (
	"time"

	"gorm.io/gorm"
)

// MagicLinkToken is one emailed sign-in link.
//
// Only the hash is stored. A table of live sign-in tokens in the clear is a
// table of passwords: anyone who can read it can sign in as anybody who has
// asked for a link in the last quarter of an hour.
type MagicLinkToken struct {
	ID        uint   ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID    string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	TokenHash string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// UsedAt is set the moment the token is spent, and the row is kept rather
	// than deleted: a second click should be able to say "this link has already
	// been used" instead of "this link is invalid", which is what somebody sees
	// when they click the same email twice and needs a different answer.
	UsedAt *time.Time ` + "`" + `json:"used_at,omitempty"` + "`" + `

	// What asked for it, for the admin's activity log and for anybody looking
	// at a link they did not request.
	IPAddress string ` + "`" + `gorm:"size:45" json:"ip_address,omitempty"` + "`" + `
	UserAgent string ` + "`" + `gorm:"size:500" json:"user_agent,omitempty"` + "`" + `

	ExpiresAt time.Time ` + "`" + `gorm:"index;not null" json:"expires_at"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}

// CleanupExpiredMagicLinks removes links nobody can use any more.
func CleanupExpiredMagicLinks(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now().Add(-24*time.Hour)).
		Delete(&MagicLinkToken{}).Error
}
`
}

func magicLinkServiceGo(module string) string {
	return `package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"` + module + `/internal/models"
)

// How long a sign-in link lives.
//
// Fifteen minutes is the compromise everybody lands on: long enough to switch
// to a phone and find the email, short enough that a link forwarded or left in
// a mailbox is not a standing key to the account.
const MagicLinkExpiry = 15 * time.Minute

// MagicLinkRate is how often one account may ask for a link. Without it the
// form is a way to send somebody a hundred emails.
const MagicLinkRate = 60 * time.Second

var (
	// ErrMagicLinkUsed is a link that has already signed somebody in. Told
	// apart from an invalid one on purpose: clicking the same email twice is
	// an ordinary mistake and deserves an answer that says so.
	ErrMagicLinkUsed = errors.New("this sign-in link has already been used")
	// ErrMagicLinkExpired is a link past its life.
	ErrMagicLinkExpired = errors.New("this sign-in link has expired")
	// ErrMagicLinkInvalid is anything else: a token that never existed, or one
	// whose row is gone.
	ErrMagicLinkInvalid = errors.New("this sign-in link is not valid")
	// ErrMagicLinkTooSoon is the rate limit.
	ErrMagicLinkTooSoon = errors.New("a link was sent a moment ago: check your email")
)

// IssueMagicLink records a new link for a user and returns the token to email.
//
// The plaintext is returned once, here, and never stored.
func IssueMagicLink(ctx context.Context, db *gorm.DB, userID, ip, agent string) (string, error) {
	var recent int64
	if err := db.WithContext(ctx).Model(&models.MagicLinkToken{}).
		Where("user_id = ? AND created_at > ?", userID, time.Now().Add(-MagicLinkRate)).
		Count(&recent).Error; err != nil {
		return "", fmt.Errorf("counting recent links: %w", err)
	}
	if recent > 0 {
		return "", ErrMagicLinkTooSoon
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating a sign-in link: %w", err)
	}
	token := hex.EncodeToString(raw)

	row := models.MagicLinkToken{
		UserID:    userID,
		TokenHash: HashMagicLink(token),
		IPAddress: ip,
		UserAgent: agent,
		ExpiresAt: time.Now().Add(MagicLinkExpiry),
	}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return "", fmt.Errorf("recording a sign-in link: %w", err)
	}
	return token, nil
}

// HashMagicLink is how a token is stored and looked up.
func HashMagicLink(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// ConsumeMagicLink spends a token and returns whose it was.
//
// The update is conditional on the token still being unused, so two clicks
// arriving together cannot both succeed: exactly one updates a row, and the
// other is told the link has been used.
func ConsumeMagicLink(ctx context.Context, db *gorm.DB, token string) (string, error) {
	if token == "" {
		return "", ErrMagicLinkInvalid
	}

	var row models.MagicLinkToken
	err := db.WithContext(ctx).Where("token_hash = ?", HashMagicLink(token)).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrMagicLinkInvalid
		}
		return "", fmt.Errorf("reading the sign-in link: %w", err)
	}
	if row.UsedAt != nil {
		return "", ErrMagicLinkUsed
	}
	if time.Now().After(row.ExpiresAt) {
		return "", ErrMagicLinkExpired
	}

	now := time.Now()
	res := db.WithContext(ctx).Model(&models.MagicLinkToken{}).
		Where("id = ? AND used_at IS NULL", row.ID).
		Update("used_at", now)
	if res.Error != nil {
		return "", fmt.Errorf("spending the sign-in link: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		// Somebody else spent it between the read and the update.
		return "", ErrMagicLinkUsed
	}
	return row.UserID, nil
}
`
}

func magicLinkServiceTestGo(module string) string {
	return `package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"` + module + `/internal/models"
	"` + module + `/internal/services"
)

func magicLinkDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.MagicLinkToken{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestAMagicLinkSignsInOnce(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, err := services.IssueMagicLink(ctx, db, "user-1", "127.0.0.1", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(token) < 32 {
		t.Errorf("token %q is too short to be unguessable", token)
	}

	// The plaintext is never stored: a table of live tokens is a table of
	// passwords.
	var row models.MagicLinkToken
	if err := db.First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.TokenHash == token {
		t.Error("the token is stored in the clear")
	}

	userID, err := services.ConsumeMagicLink(ctx, db, token)
	if err != nil || userID != "user-1" {
		t.Fatalf("consume = %q, %v", userID, err)
	}

	if _, err := services.ConsumeMagicLink(ctx, db, token); !errors.Is(err, services.ErrMagicLinkUsed) {
		t.Errorf("a spent link signed somebody in again: %v", err)
	}
}

// Clicking the same email twice is an ordinary mistake, and deserves an answer
// that says so rather than "invalid link".
func TestASpentLinkSaysItIsSpent(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, _ := services.IssueMagicLink(ctx, db, "user-1", "", "")
	if _, err := services.ConsumeMagicLink(ctx, db, token); err != nil {
		t.Fatal(err)
	}
	_, err := services.ConsumeMagicLink(ctx, db, token)
	if !errors.Is(err, services.ErrMagicLinkUsed) {
		t.Errorf("err = %v, want ErrMagicLinkUsed", err)
	}
	if _, err := services.ConsumeMagicLink(ctx, db, "never-existed"); !errors.Is(err, services.ErrMagicLinkInvalid) {
		t.Errorf("an unknown token should be invalid, got %v", err)
	}
}

func TestAnExpiredLinkIsRefused(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	token, _ := services.IssueMagicLink(ctx, db, "user-1", "", "")
	if err := db.Model(&models.MagicLinkToken{}).Where("1 = 1").
		Update("expires_at", time.Now().Add(-time.Minute)).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := services.ConsumeMagicLink(ctx, db, token); !errors.Is(err, services.ErrMagicLinkExpired) {
		t.Errorf("an expired link was accepted: %v", err)
	}
}

// Without a rate limit the form is a way to send somebody a hundred emails.
func TestLinksAreRateLimitedPerAccount(t *testing.T) {
	db := magicLinkDB(t)
	ctx := context.Background()

	if _, err := services.IssueMagicLink(ctx, db, "user-1", "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := services.IssueMagicLink(ctx, db, "user-1", "", ""); !errors.Is(err, services.ErrMagicLinkTooSoon) {
		t.Errorf("a second link was issued immediately: %v", err)
	}
	// A different account is unaffected.
	if _, err := services.IssueMagicLink(ctx, db, "user-2", "", ""); err != nil {
		t.Errorf("one account's link blocked another's: %v", err)
	}
}
`
}

func magicLinkHandlerGo(module string) string {
	return `package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"` + module + `/internal/mail"
	"` + module + `/internal/models"
	"` + module + `/internal/respond"
	"` + module + `/internal/services"
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
	Email string ` + "`" + `json:"email" binding:"required,email"` + "`" + `
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
	Token string ` + "`" + `json:"token" binding:"required"` + "`" + `
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
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UsedAt    *time.Time ` + "`" + `json:"used_at,omitempty"` + "`" + `
	ExpiresAt time.Time  ` + "`" + `json:"expires_at"` + "`" + `
	IPAddress string     ` + "`" + `json:"ip_address,omitempty"` + "`" + `
	UserAgent string     ` + "`" + `json:"user_agent,omitempty"` + "`" + `
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
`
}
