package scaffold

import (
	"fmt"
	"path/filepath"
)

// Signing in with Google from the mobile app.
//
// Every generated Expo app had a "Continue with Google" button that could not
// work. It asked the API to send it back to a deep link and read the tokens out
// of the query string; the API ignored redirect_uri entirely, finished the
// handshake, set HttpOnly cookies and redirected to the web frontend. So the
// in-app browser landed on a web page, the app never got its tokens, and
// cookies are no use to a client that authenticates with a bearer token. It had
// never worked (grit#174).
//
// Three decisions, because this is the feature where the obvious fix is the
// wrong one:
//
//   - The redirect carries a ONE-TIME CODE, never the tokens. Tokens in a
//     redirect end up in browser history, in the Referer of the next request,
//     and in any proxy log on the way. The app exchanges the code over POST.
//   - redirect_uri is validated against an allowlist, and a request naming
//     anything else is refused rather than redirected. A callback that will
//     send a code wherever the query string says is an open redirect, which is
//     worse than the dead button it replaces.
//   - The code is single-use, short-lived, stored only as a SHA-256, and spent
//     by a conditional update, so two exchanges arriving together cannot both
//     succeed. The same shape as the magic link, for the same reasons.

// writeOAuthNativeFiles writes the model, the service and the handler for the
// native code exchange.
func writeOAuthNativeFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "oauth_exchange_code.go"):          oauthExchangeModelGo(),
		filepath.Join(apiRoot, "internal", "services", "oauth_exchange.go"):             oauthExchangeServiceGo(opts.Module()),
		filepath.Join(apiRoot, "internal", "services", "oauth_exchange_test.go"):        oauthExchangeServiceTestGo(opts.Module()),
		filepath.Join(apiRoot, "internal", "handlers", "oauth_exchange.go"):             oauthExchangeHandlerGo(opts.Module()),
		filepath.Join(apiRoot, "internal", "services", "oauth_native_redirect.go"):      oauthNativeRedirectGo(),
		filepath.Join(apiRoot, "internal", "services", "oauth_native_redirect_test.go"): oauthNativeRedirectTestGo(opts.Module()),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
	}
	return nil
}

func oauthExchangeModelGo() string {
	return `package models

import "time"

// OAuthExchangeCode is a one-time code a native client swaps for tokens.
//
// A mobile app cannot use the HttpOnly cookies the web callback sets, and
// putting the tokens in the redirect would put a refresh token in browser
// history. So the callback redirects to the app's deep link with a code, and
// the app POSTs that code to /auth/oauth/exchange.
//
// Stored as a hash, like every other credential in this project: a row in this
// table is not enough to sign in as anybody.
type OAuthExchangeCode struct {
	ID       uint   ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	UserID   string ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	CodeHash string ` + "`" + `gorm:"size:64;uniqueIndex;not null" json:"-"` + "`" + `

	// UsedAt is set the moment the code is spent, and the row is kept rather
	// than deleted, so a second attempt can be told apart from a code that
	// never existed. That difference matters here: the first is a retry, the
	// second is somebody guessing.
	UsedAt *time.Time ` + "`" + `json:"used_at,omitempty"` + "`" + `

	// Where the code was sent, recorded for the audit trail. Not used to
	// validate the exchange: the code is the credential.
	RedirectURI string ` + "`" + `gorm:"size:255" json:"redirect_uri,omitempty"` + "`" + `
	IPAddress   string ` + "`" + `gorm:"size:45" json:"ip_address,omitempty"` + "`" + `

	ExpiresAt time.Time ` + "`" + `gorm:"index;not null" json:"expires_at"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
}
`
}

func oauthExchangeServiceGo(module string) string {
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

// OAuthCodeExpiry is how long a native client has to exchange its code.
//
// Sixty seconds. The code travels from the callback to the app through an
// operating-system deep link and is exchanged immediately; anything longer is
// a window in which a code left in a log is still worth something.
const OAuthCodeExpiry = 60 * time.Second

var (
	// ErrOAuthCodeUsed is a code that has already been exchanged. Told apart
	// from an invalid one deliberately: a retry of a request whose response was
	// lost looks like this, and so does somebody replaying a captured code.
	// Both get refused; only the log can tell them apart.
	ErrOAuthCodeUsed = errors.New("this sign-in code has already been used")
	// ErrOAuthCodeExpired is a code past its minute.
	ErrOAuthCodeExpired = errors.New("this sign-in code has expired")
	// ErrOAuthCodeInvalid is anything else: a code that never existed.
	ErrOAuthCodeInvalid = errors.New("this sign-in code is not valid")
)

// HashOAuthCode is how a code is stored and looked up.
func HashOAuthCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

// IssueOAuthCode records a one-time code for a user and returns it.
//
// The plaintext is returned once, here, and never stored.
func IssueOAuthCode(ctx context.Context, db *gorm.DB, userID, redirectURI, ip string) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating a sign-in code: %w", err)
	}
	code := hex.EncodeToString(raw)

	row := models.OAuthExchangeCode{
		UserID:      userID,
		CodeHash:    HashOAuthCode(code),
		RedirectURI: redirectURI,
		IPAddress:   ip,
		ExpiresAt:   time.Now().Add(OAuthCodeExpiry),
	}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return "", fmt.Errorf("recording a sign-in code: %w", err)
	}
	return code, nil
}

// ConsumeOAuthCode spends a code and returns whose it was.
//
// The update is conditional on the code still being unused, so two exchanges
// arriving together cannot both succeed: exactly one updates a row, and the
// other is told the code has been used.
func ConsumeOAuthCode(ctx context.Context, db *gorm.DB, code string) (string, error) {
	if code == "" {
		return "", ErrOAuthCodeInvalid
	}

	var row models.OAuthExchangeCode
	err := db.WithContext(ctx).Where("code_hash = ?", HashOAuthCode(code)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrOAuthCodeInvalid
	}
	if err != nil {
		return "", fmt.Errorf("looking up a sign-in code: %w", err)
	}
	if row.UsedAt != nil {
		return "", ErrOAuthCodeUsed
	}
	if time.Now().After(row.ExpiresAt) {
		return "", ErrOAuthCodeExpired
	}

	now := time.Now()
	res := db.WithContext(ctx).Model(&models.OAuthExchangeCode{}).
		Where("id = ? AND used_at IS NULL", row.ID).
		Update("used_at", now)
	if res.Error != nil {
		return "", fmt.Errorf("spending a sign-in code: %w", res.Error)
	}
	// Nothing updated means another request took it between the read above and
	// this write. Exactly one caller may win.
	if res.RowsAffected == 0 {
		return "", ErrOAuthCodeUsed
	}
	return row.UserID, nil
}

// CleanupExpiredOAuthCodes removes codes nobody can use any more.
func CleanupExpiredOAuthCodes(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now().Add(-24*time.Hour)).
		Delete(&models.OAuthExchangeCode{}).Error
}
`
}

func oauthNativeRedirectGo() string {
	return `package services

import (
	"net/url"
	"os"
	"strings"
)

// Which deep links the OAuth callback may send a code to.
//
// An allowlist, not a pattern. A callback that redirects to whatever
// redirect_uri says is an open redirect, and one that carries a sign-in code is
// an account takeover: an attacker starts the flow with their own redirect, the
// victim signs in, and the code arrives at the attacker's app.
//
// Read from OAUTH_NATIVE_REDIRECTS, comma separated. A generated project ships
// with its own Expo app's scheme in .env.example and nothing else.
func NativeRedirectAllowed(redirectURI string) bool {
	if redirectURI == "" {
		return false
	}
	for _, allowed := range NativeRedirects() {
		if redirectURI == allowed {
			return true
		}
	}
	return false
}

// NativeRedirects is the configured allowlist.
func NativeRedirects() []string {
	raw := os.Getenv("OAUTH_NATIVE_REDIRECTS")
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// NativeRedirectWithCode is where to send the browser once the handshake is
// done: the app's deep link, carrying the one-time code.
//
// Built with url.Values rather than string concatenation, so a scheme that
// already has a query string keeps it and a code is never spliced into one.
func NativeRedirectWithCode(redirectURI, code string) (string, error) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("code", code)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// NativeRedirectWithError is the same for a handshake that failed. The app
// needs to know, and it cannot read an HTML error page.
func NativeRedirectWithError(redirectURI, message string) (string, error) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("error", message)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
`
}

func oauthExchangeHandlerGo(module string) string {
	return `package handlers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	"` + module + `/internal/respond"
	"` + module + `/internal/services"
)

// OAuthExchange swaps a one-time code for a token pair.
//
// This is how a native client finishes a social login. The web callback sets
// HttpOnly cookies and redirects; a mobile app cannot use those, and putting
// the tokens in the redirect would put a refresh token in browser history. So
// the callback redirects to the app with a code and the app posts it here.
//
// Public by design, like every other step of a login flow: the code IS the
// credential, and it is single-use and lives for a minute.
func (h *AuthHandler) OAuthExchange(c *gin.Context) {
	var body struct {
		Code string ` + "`" + `json:"code" binding:"required"` + "`" + `
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, "A code is required.")
		return
	}

	userID, err := services.ConsumeOAuthCode(c.Request.Context(), h.DB, body.Code)
	if err != nil {
		// One answer for every failure. A client cannot do anything different
		// with "used" than with "never existed", and saying which tells
		// somebody guessing codes that they got one right.
		switch {
		case errors.Is(err, services.ErrOAuthCodeUsed),
			errors.Is(err, services.ErrOAuthCodeExpired),
			errors.Is(err, services.ErrOAuthCodeInvalid):
			log.Printf("oauth exchange refused: %v", err)
		default:
			log.Printf("oauth exchange: %v", err)
		}
		respond.Fail(c, respond.CodeUnauthorized,
			"This sign-in code is not valid. Start the sign-in again.")
		return
	}

	user, err := h.AuthService.UserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		log.Printf("oauth exchange: user %s is gone", userID)
		respond.Fail(c, respond.CodeUnauthorized, "This sign-in code is not valid.")
		return
	}
	if !user.Active {
		respond.Fail(c, respond.CodeForbidden, "Your account has been disabled.")
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("oauth exchange: generating tokens for %s: %v", user.ID, err)
		respond.WriteError(c, err, "Failed to sign in.")
		return
	}

	// Recorded as a session, so a mobile sign-in is listed and revocable
	// exactly like one from a browser.
	if _, err := services.CreateSession(h.DB, c, user.ID, tokens.RefreshToken); err != nil {
		log.Printf("oauth exchange: recording session for %s: %v", user.ID, err)
	}

	respond.OK(c, gin.H{"user": user, "tokens": tokens})
}
`
}

func oauthExchangeServiceTestGo(module string) string {
	return `package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"` + module + `/internal/models"
)

func oauthCodeDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening the test database: %v", err)
	}

	// One connection, and only one.
	//
	// Every connection to ":memory:" gets its OWN empty database, and the pool
	// opens a second one as soon as two queries overlap. The race test below
	// runs two goroutines through ConsumeOAuthCode, so one of them landed on a
	// fresh empty database and answered "no such table: o_auth_exchange_codes".
	// It passed or failed on timing, which reads as a broken feature rather
	// than a test-harness problem: it went green locally and red in CI on the
	// same commit. internal/testkit caps its pool for exactly this reason; a
	// helper that opens its own database has to do it too.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("reaching the connection pool: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&models.OAuthExchangeCode{}); err != nil {
		t.Fatalf("migrating: %v", err)
	}
	return db
}

func TestACodeSignsInOnce(t *testing.T) {
	db := oauthCodeDB(t)
	ctx := context.Background()

	code, err := IssueOAuthCode(ctx, db, "user-1", "myapp://callback", "127.0.0.1")
	if err != nil {
		t.Fatalf("IssueOAuthCode: %v", err)
	}
	if code == "" {
		t.Fatal("no code was returned")
	}

	// The plaintext is never stored. A row in this table must not be enough to
	// sign in as anybody.
	var row models.OAuthExchangeCode
	if err := db.First(&row).Error; err != nil {
		t.Fatalf("reading the row: %v", err)
	}
	if row.CodeHash == code {
		t.Error("the code is stored in plaintext")
	}
	if row.CodeHash != HashOAuthCode(code) {
		t.Error("the stored hash is not the hash of the code")
	}

	userID, err := ConsumeOAuthCode(ctx, db, code)
	if err != nil {
		t.Fatalf("ConsumeOAuthCode: %v", err)
	}
	if userID != "user-1" {
		t.Errorf("ConsumeOAuthCode = %q, want user-1", userID)
	}

	// And only once.
	if _, err := ConsumeOAuthCode(ctx, db, code); err != ErrOAuthCodeUsed {
		t.Errorf("a second exchange returned %v, want ErrOAuthCodeUsed", err)
	}
}

func TestAnExpiredCodeIsRefused(t *testing.T) {
	db := oauthCodeDB(t)
	ctx := context.Background()

	code, err := IssueOAuthCode(ctx, db, "user-1", "myapp://callback", "")
	if err != nil {
		t.Fatalf("IssueOAuthCode: %v", err)
	}
	// Age it past its minute.
	if err := db.Model(&models.OAuthExchangeCode{}).
		Where("1 = 1").
		Update("expires_at", time.Now().Add(-time.Second)).Error; err != nil {
		t.Fatalf("ageing the code: %v", err)
	}

	if _, err := ConsumeOAuthCode(ctx, db, code); err != ErrOAuthCodeExpired {
		t.Errorf("an expired code returned %v, want ErrOAuthCodeExpired", err)
	}
}

func TestAnUnknownCodeIsRefused(t *testing.T) {
	db := oauthCodeDB(t)
	for _, code := range []string{"", "not-a-code", HashOAuthCode("guess")} {
		if _, err := ConsumeOAuthCode(context.Background(), db, code); err != ErrOAuthCodeInvalid {
			t.Errorf("ConsumeOAuthCode(%q) = %v, want ErrOAuthCodeInvalid", code, err)
		}
	}
}

// Two exchanges arriving together must not both succeed. The app retries on a
// lost response, so this is an ordinary race rather than an attack, and it has
// to resolve to one session rather than two.
func TestTwoExchangesRaceAndOnlyOneWins(t *testing.T) {
	db := oauthCodeDB(t)
	ctx := context.Background()

	code, err := IssueOAuthCode(ctx, db, "user-1", "myapp://callback", "")
	if err != nil {
		t.Fatalf("IssueOAuthCode: %v", err)
	}

	var wg sync.WaitGroup
	results := make([]error, 2)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, results[i] = ConsumeOAuthCode(ctx, db, code)
		}(i)
	}
	wg.Wait()

	won := 0
	for _, err := range results {
		if err == nil {
			won++
		}
	}
	if won != 1 {
		t.Errorf("%d of 2 exchanges succeeded, want exactly 1 (errors: %v)", won, results)
	}
}
`
}

func oauthNativeRedirectTestGo(module string) string {
	return `package services

import (
	"os"
	"strings"
	"testing"
)

// An allowlist, and the reason it is not a pattern: a callback that redirects
// wherever the query string says, carrying a sign-in code, is an account
// takeover. The attacker starts the flow with their own redirect, the victim
// signs in, and the code arrives at the attacker's app.
func TestOnlyAllowlistedRedirectsAreAccepted(t *testing.T) {
	t.Setenv("OAUTH_NATIVE_REDIRECTS", "myapp://callback, myapp://oauth")

	for _, allowed := range []string{"myapp://callback", "myapp://oauth"} {
		if !NativeRedirectAllowed(allowed) {
			t.Errorf("%q is in the allowlist and was refused", allowed)
		}
	}

	for _, refused := range []string{
		"",
		"evil://callback",
		"https://evil.example/callback",
		"myapp://callback/../evil",
		"myapp://callback?next=https://evil.example",
		// A prefix of an allowed value, which a strings.HasPrefix check would
		// let through along with everything after it.
		"myapp://callbackevil",
		// And the allowed value as a prefix of something longer.
		"myapp://callback.evil.example",
	} {
		if NativeRedirectAllowed(refused) {
			t.Errorf("%q is not in the allowlist and was accepted", refused)
		}
	}
}

// With nothing configured, nothing is allowed. A project that has not opted in
// does not have an open redirect by default.
func TestNoAllowlistMeansNoNativeRedirects(t *testing.T) {
	t.Setenv("OAUTH_NATIVE_REDIRECTS", "")
	if NativeRedirectAllowed("myapp://callback") {
		t.Error("a redirect was allowed with no allowlist configured")
	}
	if got := NativeRedirects(); len(got) != 0 {
		t.Errorf("NativeRedirects() = %v, want nothing", got)
	}
	os.Unsetenv("OAUTH_NATIVE_REDIRECTS")
	if NativeRedirectAllowed("myapp://callback") {
		t.Error("a redirect was allowed with OAUTH_NATIVE_REDIRECTS unset")
	}
}

// The code goes on as a query parameter, so a deep link that already has one
// keeps it.
func TestTheCodeIsAddedAsAQueryParameter(t *testing.T) {
	got, err := NativeRedirectWithCode("myapp://callback?from=login", "abc123")
	if err != nil {
		t.Fatalf("NativeRedirectWithCode: %v", err)
	}
	if !strings.Contains(got, "code=abc123") {
		t.Errorf("%q does not carry the code", got)
	}
	if !strings.Contains(got, "from=login") {
		t.Errorf("%q lost the query it already had", got)
	}

	failed, err := NativeRedirectWithError("myapp://callback", "Authentication failed")
	if err != nil {
		t.Fatalf("NativeRedirectWithError: %v", err)
	}
	if !strings.Contains(failed, "error=Authentication+failed") &&
		!strings.Contains(failed, "error=Authentication%20failed") {
		t.Errorf("%q does not carry the message", failed)
	}
}
`
}
