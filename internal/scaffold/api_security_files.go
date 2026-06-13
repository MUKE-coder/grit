package scaffold

// This file ships the Go-side security primitives that close the
// OWASP Top 10:2025 gaps identified in the framework's defence audit:
//
//   - internal/safefetch    — SSRF guard (A01: SSRF was folded into
//                             Broken Access Control in 2025).
//   - internal/authz        — ownership/IDOR helpers (A01).
//   - internal/middleware   — CSRF double-submit token; security-event
//                             audit helpers (A07/A09).
//
// Each piece is wired into routes/services where it makes sense, but
// stays opt-in so apps can ignore what they don't need.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeSecurityFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "safefetch", "safefetch.go"):      safefetchGo(),
		filepath.Join(apiRoot, "internal", "safefetch", "safefetch_test.go"): safefetchTestGo(),
		filepath.Join(apiRoot, "internal", "authz", "authz.go"):              authzGo(),
		filepath.Join(apiRoot, "internal", "authz", "authz_test.go"):         authzTestGo(),
		filepath.Join(apiRoot, "internal", "middleware", "csrf.go"):          csrfMiddlewareGo(),
		filepath.Join(apiRoot, "internal", "middleware", "security_log.go"):  securityLogGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// safefetchGo emits internal/safefetch/safefetch.go.
//
// The SSRF defence works at two layers:
//
//  1. Scheme + host validation BEFORE the request — reject anything
//     that isn't http/https; reject literal IPs in private/loopback/
//     link-local ranges; reject cloud-metadata hostnames.
//  2. A custom net.Dialer.Control that re-checks the resolved IP at
//     TCP-connect time. This closes the DNS-rebinding TOCTOU window
//     where validation sees a public IP and connect lands on 169.254.169.254.
func safefetchGo() string {
	return `// Package safefetch performs HTTP requests against user-supplied URLs
// with SSRF defences. Use this any time you fetch a URL the caller
// chose — webhook delivery, "fetch image from URL", PDF render from a
// URL, OEmbed expansion, etc.
//
// Coverage: OWASP Top 10:2025 A01 (SSRF was absorbed into Broken Access
// Control in 2025). The classic SSRF impact — proxying requests to the
// cloud metadata service (169.254.169.254) to steal IAM credentials —
// is blocked at TCP-connect time even if DNS resolution races.
package safefetch

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
)

// Client is the default safefetch HTTP client. Timeouts are short on
// purpose — long-running fetches are the easiest amplification vector.
var Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
			Control:   controlSafeDial,
		}).DialContext,
		MaxIdleConns:        20,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
	// Validate every redirect target — attackers love to 302 you into
	// 169.254.169.254 after a clean validation pass.
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("too many redirects")
		}
		return validateURL(req.URL)
	},
}

// ErrBlocked is returned when a target is rejected by the SSRF guard.
var ErrBlocked = errors.New("safefetch: target blocked")

// Get fetches the URL with SSRF defences. Returns ErrBlocked (wrapped)
// for any disallowed target.
func Get(ctx context.Context, rawURL string) (*http.Response, error) {
	req, err := newRequest(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	return Client.Do(req)
}

// Do runs a pre-built request through the safe client. Use this when
// you need POST/PUT or custom headers.
func Do(req *http.Request) (*http.Response, error) {
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}
	return Client.Do(req)
}

func newRequest(ctx context.Context, method, rawURL string, body interface{ Read(p []byte) (n int, err error) }) (*http.Request, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("safefetch: parsing url: %w", err)
	}
	if err := validateURL(u); err != nil {
		return nil, err
	}
	if body == nil {
		return http.NewRequestWithContext(ctx, method, u.String(), nil)
	}
	return http.NewRequestWithContext(ctx, method, u.String(), body)
}

// validateURL enforces the scheme allowlist and rejects hostnames that
// are obviously private. The IP-level check happens again at dial time
// via controlSafeDial.
func validateURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("%w: nil url", ErrBlocked)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		// ok
	default:
		return fmt.Errorf("%w: scheme %q not allowed", ErrBlocked, u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("%w: empty host", ErrBlocked)
	}
	// Cloud metadata hostnames — rejected by name in case DNS hasn't
	// been consulted yet (e.g. proxy mode).
	lower := strings.ToLower(host)
	if lower == "metadata.google.internal" || lower == "metadata" || lower == "instance-data" {
		return fmt.Errorf("%w: cloud metadata host", ErrBlocked)
	}
	// If the host parses as a literal IP, check it now. (Hostnames are
	// re-checked at dial time after resolution.)
	if ip := net.ParseIP(host); ip != nil {
		if isPrivateOrLoopback(ip) {
			return fmt.Errorf("%w: literal IP %s in private range", ErrBlocked, ip)
		}
	}
	return nil
}

// controlSafeDial runs after DNS resolution, immediately before connect.
// This is where we close the DNS-rebinding hole — even if a hostname
// validated cleanly, the resolved IP at connect time may not have.
func controlSafeDial(network, address string, _ syscall.RawConn) error {
	if !strings.HasPrefix(network, "tcp") {
		return fmt.Errorf("%w: non-tcp network %q", ErrBlocked, network)
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%w: splitting host:port: %v", ErrBlocked, err)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("%w: dial address %q is not an IP", ErrBlocked, host)
	}
	if isPrivateOrLoopback(ip) {
		return fmt.Errorf("%w: resolved IP %s in private range", ErrBlocked, ip)
	}
	return nil
}

// isPrivateOrLoopback returns true for any IP an SSRF guard must reject:
// loopback, link-local, multicast, unspecified, RFC1918 private,
// carrier-grade NAT, and the IMDS endpoints (169.254.169.254 + the IPv6
// AWS metadata fd00:ec2::254).
func isPrivateOrLoopback(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() || ip.IsInterfaceLocalMulticast() {
		return true
	}
	// Go 1.17+ exposes IsPrivate for RFC1918 / RFC4193.
	if ip.IsPrivate() {
		return true
	}
	// Carrier-grade NAT (RFC 6598).
	_, cgnat, _ := net.ParseCIDR("100.64.0.0/10")
	if cgnat.Contains(ip) {
		return true
	}
	// AWS IPv6 metadata endpoint.
	_, awsV6, _ := net.ParseCIDR("fd00:ec2::/32")
	if awsV6.Contains(ip) {
		return true
	}
	return false
}
`
}

func safefetchTestGo() string {
	return `package safefetch

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestBlocksPrivateHosts(t *testing.T) {
	cases := []string{
		"http://127.0.0.1/",
		"http://localhost/",
		"http://10.0.0.1/",
		"http://192.168.1.1/",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]/",
		"http://metadata.google.internal/computeMetadata/v1/",
		"http://100.64.0.5/",
	}
	for _, raw := range cases {
		_, err := Get(context.Background(), raw)
		if !errors.Is(err, ErrBlocked) {
			t.Errorf("Get(%q): expected ErrBlocked, got %v", raw, err)
		}
	}
}

func TestBlocksUnknownSchemes(t *testing.T) {
	cases := []string{
		"file:///etc/passwd",
		"gopher://evil.example.com/x",
		"ftp://example.com/",
	}
	for _, raw := range cases {
		_, err := Get(context.Background(), raw)
		if err == nil || !strings.Contains(err.Error(), "scheme") {
			t.Errorf("Get(%q): expected scheme rejection, got %v", raw, err)
		}
	}
}
`
}

// authzGo emits the IDOR-prevention helpers.
//
// The fix for A01 / IDOR is checking ownership on every object access.
// authz.MustOwn pulls the row and rejects with 404 if either the row
// doesn't exist OR the owner doesn't match — the same response shape
// either way prevents enumeration through error-message differences.
func authzGo() string {
	return `// Package authz contains the ownership-check helpers Grit uses to
// prevent IDOR (Insecure Direct Object Reference) — OWASP Top 10:2025
// A01 Broken Access Control's most common concrete form.
//
// The cardinal rule (from PHASE 2 §4.3 of the security course): every
// object access must be authorised against the current user, server-side.
// authz.MustOwn enforces that with a single call.
//
// Usage:
//
//	var invoice models.Invoice
//	if err := authz.MustOwn(c, db, &invoice, c.Param("id")); err != nil {
//	    return // helper has already written 404 / 401
//	}
//	// invoice belongs to c.MustGet("user_id"). Safe to use.
package authz

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Ownable is implemented by domain models whose ownership is identified
// by a single user-id column. For models with team/tenant scoping use
// CheckScope instead.
type Ownable interface {
	GetOwnerID() string
}

// ErrNotFound and ErrForbidden are returned by the helpers below so
// callers can branch (e.g. log differently) — but the HTTP responses
// they produce are deliberately identical (404 NOT_FOUND) to avoid
// leaking the existence of rows the caller doesn't own.
var (
	ErrNotFound  = errors.New("authz: not found")
	ErrForbidden = errors.New("authz: forbidden")
)

// MustOwn loads the row by id and verifies that the authenticated user
// is its owner. On any failure it writes a 404 response and returns a
// non-nil error so the caller can return immediately.
//
// The 404 (not 403) is intentional. Returning 403 confirms the row
// exists, which lets attackers enumerate IDs.
func MustOwn(c *gin.Context, db *gorm.DB, dest Ownable, id string) error {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHORIZED", "message": "Authentication required"},
		})
		return ErrForbidden
	}

	if err := db.Where("id = ?", id).First(dest).Error; err != nil {
		writeNotFound(c)
		return ErrNotFound
	}

	if dest.GetOwnerID() != userID {
		writeNotFound(c) // 404, not 403 — see comment above
		return ErrForbidden
	}
	return nil
}

// CheckScope verifies a (column, value) pair matches the current user's
// authoritative scope (e.g. team_id, tenant_id). Use this when ownership
// is by membership rather than a single user_id column.
func CheckScope(c *gin.Context, scopeKey, expectedValue string) bool {
	got, ok := c.Get(scopeKey)
	return ok && got == expectedValue
}

// RequireRoles returns a gin middleware that 403s when the authenticated
// user's role isn't in the allowlist. This is a stricter sibling of the
// generic Auth middleware — use it on admin-only routes.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if _, ok := allowed[asString(role)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{"code": "FORBIDDEN", "message": "Insufficient role"},
			})
			return
		}
		c.Next()
	}
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func writeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": gin.H{"code": "NOT_FOUND", "message": "Resource not found"},
	})
}
`
}

func authzTestGo() string {
	return `package authz

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeUserScoped struct {
	ID      string
	OwnerID string
}

func (f *fakeUserScoped) GetOwnerID() string { return f.OwnerID }

func TestRequireRolesAllowsListed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_role", "admin"); c.Next() })
	r.GET("/admin", RequireRoles("admin"), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRequireRolesBlocksOthers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_role", "user"); c.Next() })
	r.GET("/admin", RequireRoles("admin"), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
`
}

// csrfMiddlewareGo emits internal/middleware/csrf.go — a double-submit-
// cookie CSRF defence for any flow that authenticates via cookies.
//
// SPA flows use JWT in Authorization: Bearer which is not auto-sent by
// browsers, so they're not CSRF-vulnerable. The OAuth callback flow
// (cookie-based) IS vulnerable and benefits from this middleware.
func csrfMiddlewareGo() string {
	return `package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CSRF implements double-submit-cookie protection (OWASP Top 10:2025
// A01-adjacent). Use on cookie-authenticated routes. SPA + JWT in
// Authorization header doesn't need this — only cookie sessions do.
//
// Wire it like:
//
//	cookieRoutes := r.Group("/")
//	cookieRoutes.Use(middleware.CSRF())
//	cookieRoutes.POST("/api/auth/oauth/google/callback", h.OAuthCallback)
//
// The middleware sets a "grit_csrf" cookie on safe requests; unsafe
// methods (POST/PUT/PATCH/DELETE) must echo it via the X-CSRF-Token
// header. SameSite=Strict on the cookie is a second layer of defence.
func CSRF() gin.HandlerFunc {
	const (
		cookieName = "grit_csrf"
		headerName = "X-CSRF-Token"
	)
	return func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			// Safe method — issue or refresh the token cookie if missing.
			if existing, err := c.Cookie(cookieName); err != nil || existing == "" {
				token, gerr := newCSRFToken()
				if gerr == nil {
					c.SetSameSite(http.SameSiteStrictMode)
					c.SetCookie(cookieName, token, 86400, "/", "", c.Request.TLS != nil, false)
				}
			}
			c.Next()
			return
		}

		// Unsafe method — require matching header + cookie.
		cookieToken, _ := c.Cookie(cookieName)
		headerToken := c.GetHeader(headerName)
		if cookieToken == "" || headerToken == "" ||
			subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "CSRF_INVALID",
					"message": "CSRF token missing or invalid",
				},
			})
			return
		}
		c.Next()
	}
}

func newCSRFToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", errors.New("generating CSRF token")
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// AutoCSRF is the global, default-on CSRF guard. It only enforces CSRF
// when the request carries the grit_access auth cookie — i.e., a browser
// client that the API previously authenticated via cookies. Native
// mobile / desktop clients use Authorization: Bearer (which the browser
// never auto-sends across origins), so they're immune to CSRF and pass
// through with no header required.
//
// Pair this with the cookie helpers in services/auth.go (SetAuthCookies /
// ClearAuthCookies). Together they close OWASP A01 + A05 for cookie auth
// without forcing every route to opt in.
//
// Behaviour:
//
//   - GET/HEAD/OPTIONS                    → issue grit_csrf cookie if missing.
//   - POST/PUT/PATCH/DELETE with cookie   → require matching X-CSRF-Token.
//   - POST/PUT/PATCH/DELETE bearer-only   → no-op (header auth is CSRF-safe).
//   - Login / register / refresh routes   → skipped (they MINT the cookie).
func AutoCSRF() gin.HandlerFunc {
	const (
		csrfCookie   = "grit_csrf"
		csrfHeader   = "X-CSRF-Token"
		accessCookie = "grit_access"
	)
	// Routes that bootstrap the session (login etc.) can't have a CSRF
	// cookie yet — exempt them so users can sign in on the first try.
	bootstrap := map[string]bool{
		"/api/auth/login":           true,
		"/api/auth/register":        true,
		"/api/auth/refresh":         true,
		"/api/auth/forgot-password": true,
		"/api/auth/reset-password":  true,
		"/api/auth/totp/verify":     true,
		"/api/auth/totp/backup-codes/verify": true,
	}
	return func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)
		path := c.Request.URL.Path

		// Issue / refresh the CSRF cookie on safe methods. We do this even
		// for unauthenticated visitors so SPA bootstrap code can read the
		// token from a sibling /api/auth/csrf call without a chicken-and-egg.
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			if existing, err := c.Cookie(csrfCookie); err != nil || existing == "" {
				token, gerr := newCSRFToken()
				if gerr == nil {
					c.SetSameSite(http.SameSiteLaxMode)
					c.SetCookie(csrfCookie, token, 86400, "/", "", c.Request.TLS != nil, false)
				}
			}
			c.Next()
			return
		}

		// Bootstrap auth endpoints are exempt — they create the session.
		if bootstrap[path] {
			c.Next()
			return
		}

		// State-changing method. If the client did NOT authenticate via
		// cookie, this is a bearer flow (or anonymous) — neither needs CSRF.
		accessVal, _ := c.Cookie(accessCookie)
		if accessVal == "" {
			c.Next()
			return
		}

		// Cookie-authenticated mutation: require the double-submit token.
		cookieToken, _ := c.Cookie(csrfCookie)
		headerToken := c.GetHeader(csrfHeader)
		if cookieToken == "" || headerToken == "" ||
			subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "CSRF_INVALID",
					"message": "CSRF token missing or invalid",
				},
			})
			return
		}
		c.Next()
	}
}
`
}

// securityLogGo emits internal/middleware/security_log.go — typed
// helpers for the "audit log of security events" pattern (OWASP A09).
//
// The pattern is the one called out in PHASE 2 §4.5 and PHASE 5 §11.3:
// log every authN/authZ-relevant event with who/what/when, and ensure
// it lands in a tamper-evident store. Grit's existing activity log
// already handles the storage; this file just gives strongly-typed
// helpers for the security event types that matter for compliance
// (SOC 2 / ISO 27001 / GDPR).
func securityLogGo() string {
	return `package middleware

import (
	"context"
	"log"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// SecurityEvent is the discriminator that goes into ActivityLog.Action
// for OWASP A09 "security event" records. Use these constants — never
// raw strings — so the audit dashboard can filter on a stable enum.
const (
	EventLoginSuccess        = "security.login.success"
	EventLoginFailure        = "security.login.failure"
	EventLogout              = "security.logout"
	EventPasswordChanged     = "security.password.changed"
	EventPasswordResetReq    = "security.password.reset_requested"
	EventTOTPEnabled         = "security.totp.enabled"
	EventTOTPDisabled        = "security.totp.disabled"
	EventTOTPChallengeFail   = "security.totp.challenge_failed"
	EventRoleChanged         = "security.role.changed"
	EventAccountLocked       = "security.account.locked"
	EventAuthZDenied         = "security.authz.denied"
	EventSuspiciousRequest   = "security.suspicious_request"
)

// LogSecurityEvent records a security-relevant event to the activity log.
//
// We piggyback on the existing tamper-evident ActivityLog table — Method
// "SECURITY" + Path = event-constant is the convention. This means every
// security event automatically picks up the hash-chain integrity the
// activity log already provides (a deleted security event breaks the
// chain at verify time), without a separate model to maintain.
//
// Database errors are logged but never returned — failing a login or
// logout because the audit log had a hiccup turns audit infrastructure
// into a DoS amplifier. The request itself was usually fine.
func LogSecurityEvent(ctx context.Context, db *gorm.DB, userID, event, ip, userAgent string) {
	if db == nil {
		return
	}
	entry := models.ActivityLog{
		UserID:    userID,
		Method:    "SECURITY",
		Path:      event,
		IPAddress: ip,
		UserAgent: userAgent,
	}
	if err := db.WithContext(ctx).Create(&entry).Error; err != nil {
		log.Printf("security_log: failed to record %s: %v", event, err)
	}
}
`
}
