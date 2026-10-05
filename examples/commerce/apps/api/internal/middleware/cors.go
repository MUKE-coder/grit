package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// isWailsOrigin reports whether the request came from the Wails desktop
// webview, whose origin is not stably enumerable:
//
//	wails dev   (Windows)     http://wails.localhost:34115   <- port from wails.json
//	wails build (Windows)     http://wails.localhost
//	wails build (mac/linux)   wails://wails
//
// The dev-server port is configurable, so pinning exact origins in
// CORS_ORIGINS is fragile: change the port and the desktop login silently
// starts failing with an opaque "Network Error". Match on the host instead.
//
// Safe by construction: "wails.localhost" is a virtual host the webview
// resolves internally, so a page on the public internet cannot be served
// from it and cannot forge this origin. Every other origin still has to be
// in the explicit CORS_ORIGINS allowlist.
func isWailsOrigin(origin string) bool {
	if origin == "wails://wails" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Hostname() == "wails.localhost"
}

// CORS creates a CORS middleware with a fixed allowlist.
//
// Kept for callers that genuinely have a static list. Anything user-facing
// should use CORSDynamic, so adding a domain does not need a redeploy.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return CORSDynamic(func() []string { return allowedOrigins })
}

// CORSDynamic resolves the allowlist per request.
//
// Per request rather than at construction, because the point of putting
// origins in settings is that somebody adds a domain at 9pm and it works.
// Capturing the slice at boot would mean the setting existed and did nothing
// until the next deploy, which is worse than not offering it.
//
// The cost is a map build per request over a list with single digits of
// entries, against a settings store that is already cached in memory. That is
// not the thing to optimise.
func CORSDynamic(resolve func() []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowed := false
		for _, candidate := range resolve() {
			if candidate == origin && origin != "" {
				allowed = true
				break
			}
		}
		if allowed || isWailsOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// X-CSRF-Token + Idempotency-Key are injected by the web and admin
		// axios clients on every unsafe method. Without them in the allowed
		// list, the browser's preflight strips the headers and the request
		// either fails the AutoCSRF check or replays without an idempotency
		// guarantee. Authorization stays for native bearer clients.
		// X-API-Key is not optional here. A storefront calls the public
		// endpoints with it, cross-origin, and a header missing from this list
		// is stripped by the browser during preflight: the request fails in
		// every browser and works perfectly under curl, which is the worst
		// shape a bug can have.
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-API-Key, X-CSRF-Token, Idempotency-Key, X-Public-IP-Hint")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
