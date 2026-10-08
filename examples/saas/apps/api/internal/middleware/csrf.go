package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"saas/apps/api/internal/respond"
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
			respond.Fail(c, respond.CodeCSRFInvalid, "CSRF token missing or invalid")
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
//
// isSAMLACSPath matches /api/auth/saml/<slug>/acs on both the versioned and
// unversioned forms of the path, since the IdP posts to whichever one the
// customer registered.
func isSAMLACSPath(path string) bool {
	if !strings.HasSuffix(path, "/acs") {
		return false
	}
	return strings.HasPrefix(path, "/api/auth/saml/") ||
		(strings.HasPrefix(path, "/api/") && strings.Contains(path, "/auth/saml/"))
}

// authRoute is what follows /api/auth/ or /api/<version>/auth/ in path, or ""
// for anything else. The bootstrap list was written as /api/auth/login while the
// router mounts /api/v1, so not one of its entries ever matched a request.
func authRoute(path string) string {
	rest, ok := strings.CutPrefix(path, "/api/")
	if !ok {
		return ""
	}
	if i := strings.IndexByte(rest, '/'); i > 1 && rest[0] == 'v' && strings.Trim(rest[1:i], "0123456789") == "" {
		rest = rest[i+1:]
	}
	name, ok := strings.CutPrefix(rest, "auth/")
	if !ok {
		return ""
	}
	return name
}

func AutoCSRF() gin.HandlerFunc {
	const (
		csrfCookie   = "grit_csrf"
		csrfHeader   = "X-CSRF-Token"
		accessCookie = "grit_access"
	)
	// Routes that bootstrap the session (login etc.) can't have a CSRF
	// cookie yet — exempt them so users can sign in on the first try.
	bootstrap := map[string]bool{
		"login":                    true,
		"register":                 true,
		"refresh":                  true,
		"forgot-password":          true,
		"reset-password":           true,
		"totp/verify":              true,
		"totp/backup-codes/verify": true,
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
		if bootstrap[authRoute(path)] {
			c.Next()
			return
		}

		// SAML assertion consumer. The identity provider POSTs the signed
		// assertion here from its own origin, so there is no CSRF token to
		// present and never will be — the assertion's signature, audience
		// restriction and validity window are what authenticate this request.
		//
		// It usually slips through the cookie check below anyway (a cross-site
		// POST carries no SameSite=Lax cookie), but only by accident: an
		// already-signed-in user coming back through their IdP's app tile can
		// arrive holding grit_access and would then be rejected. Exempting the
		// path states the intent instead of depending on that.
		if isSAMLACSPath(path) {
			c.Next()
			return
		}

		// Bearer-authenticated request → CSRF-immune. A request that carries
		// an "Authorization: Bearer" header authenticates explicitly, not via
		// ambient cookies, so it can't be forged cross-site. We check this
		// BEFORE the cookie check because native clients (React Native fetch,
		// OkHttp) transparently store and resend the grit_access cookie set at
		// login — that stray cookie must not trip CSRF on a bearer mutation.
		if authz := c.GetHeader("Authorization"); strings.HasPrefix(authz, "Bearer ") {
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
			respond.Fail(c, respond.CodeCSRFInvalid, "CSRF token missing or invalid")
			return
		}
		c.Next()
	}
}
